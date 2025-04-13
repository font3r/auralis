package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/google/uuid"
)

type DataSet struct {
	columnDescriptors []ColumnDescriptor
	rows              []Row
}

type Row struct {
	cells []any
}

type DataType string

var (
	ErrTableNotFound = AuraError{Code: "TABLE_NOT_FOUND", Message: "table not found"}
)

const (
	smallint         DataType = "smallint" // 2B
	integer          DataType = "integer"  // 4
	bigint           DataType = "bigint"   // 8
	varchar          DataType = "varchar"
	uniqueidentifier DataType = "uniqueidentifier" // 16
	boolean          DataType = "boolean"          // 1
)

const terminationByte = byte(10)

func cretateTable(td TableDescriptor) error {
	err := addTableDescriptor(td)
	if err != nil {
		return err
	}

	f, err := os.Create(getTableDiskPath(td.schemeTable))
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func writeIntoTable(tableDescriptor TableDescriptor, dataSet DataSet) error {
	log.Printf("INFO: executing insert query %+v", dataSet)
	f, err := os.OpenFile(getTableDiskPath(tableDescriptor.schemeTable), os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTableNotFound
		}
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, row := range dataSet.rows {
		for cellIndex, cell := range row.cells {
			var val []byte
			switch tableDescriptor.columnDescriptors[cellIndex].dataType {
			case smallint:
				{
					buf := bytes.NewBuffer(make([]byte, 0, getDataTypeByteSize(smallint)))
					// binary writer is implicitly converting to uint16 with two's complement
					err := binary.Write(buf, binary.BigEndian, uint16(cell.(int)))
					if err != nil {
						return err
					}

					val = buf.Bytes()
				}
			case varchar:
				{
					// we don't care about endianness because we support only utf-8 for now
					padded := make([]byte, getDataTypeByteSize(varchar))
					copy(padded, cell.(string))
					val = padded
				}
			case uniqueidentifier:
				{
					val, err = cell.(uuid.UUID).MarshalBinary()
					if err != nil {
						return err
					}
				}
			default:
				return errors.New("unhandled type")
			}

			_, err = w.Write(val)
			if err != nil {
				return err
			}
		}

		err = w.WriteByte(terminationByte)
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func readFromTable(tableDescriptor TableDescriptor, query SelectQuery) (*DataSet, error) {
	log.Printf("INFO: executing select query %+v", query)
	f, err := os.Open(getTableDiskPath(tableDescriptor.schemeTable))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrTableNotFound
		}

		return nil, err
	}
	defer f.Close()

	log.Printf("INFO: %s.%s table descriptor %+v\n", tableDescriptor.schemeTable.scheme,
		tableDescriptor.schemeTable.name, tableDescriptor)

	dataSet := DataSet{}
	for _, v := range tableDescriptor.columnDescriptors {
		if slices.Contains(query.dataColumns, v.name) {
			dataSet.columnDescriptors = append(dataSet.columnDescriptors, v)
		}
	}

	var fileOffset int64 = 0
	rowBuffSize := calculateRowBuffer(tableDescriptor)
	rowBuf := make([]byte, rowBuffSize)

	for {
		n, err := f.ReadAt(rowBuf, fileOffset)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		row := Row{}
		rowOffset := 0

		var cellDataSize int
		for _, cd := range tableDescriptor.columnDescriptors {
			if !slices.Contains(query.dataColumns, cd.name) {
				rowOffset += getDataTypeByteSize(cd.dataType)
				continue
			}

			switch cd.dataType {
			case smallint:
				{
					cellDataSize = getDataTypeByteSize(smallint)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])
					value := int16(binary.BigEndian.Uint16(data))

					conditions := GetMatchingCondition(query.conditions, cd.name)
					if len(conditions) > 0 {
						// for now assume only one condition
						if EvaluateIntCondition(conditions[0], value) {
							row.cells = append(row.cells, value)
							continue
						}
					} else {
						row.cells = append(row.cells, value)
					}
				}
			case varchar:
				{
					cellDataSize = getDataTypeByteSize(varchar)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])

					row.cells = append(row.cells, string(bytes.TrimRight(data, "\x00")))
				}
			case uniqueidentifier:
				{
					cellDataSize = getDataTypeByteSize(uniqueidentifier)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])

					row.cells = append(row.cells, uuid.UUID(data))

				}
			default:
				return &dataSet, errors.New("unhandled type")
			}

			rowOffset += cellDataSize
		}

		dataSet.rows = append(dataSet.rows, row)
		fileOffset += int64(rowBuffSize)
		clear(rowBuf)
	}

	return &dataSet, nil
}

func calculateRowBuffer(td TableDescriptor) int {
	size := 0
	for _, v := range td.columnDescriptors {
		size += getDataTypeByteSize(v.dataType)
	}
	size += 1 // termination byte

	return size
}

func getTableDiskPath(source SchemeTable[string, string]) string {
	return fmt.Sprintf("./data/%s.%s", source.scheme, source.name)
}

func getDataTypeByteSize(dataType DataType) int {
	switch dataType {
	case smallint:
		return 2 // int16
	case integer:
		return 4 // int32
	case bigint:
		return 8 // int64
	case varchar:
		return 16 // TODO: this should be configurable
	case uniqueidentifier:
		return 16
	case boolean:
		return 1
	default:
		panic("unhandled type")
	}
}
