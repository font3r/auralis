package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/google/uuid"
)

type DataSet struct {
	columnDescriptors []ColumnDescriptor
	rows              []Row
}

type Row struct {
	cells []Cell
}

type Cell struct {
	value any
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

	f, err := os.Create(getTableDiskPath(td.source))
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// TODO: dataSet validation against table descriptor
func writeIntoTable(source SchemeTable[string, string], dataSet DataSet) error {
	f, err := os.OpenFile(getTableDiskPath(source), os.O_WRONLY|os.O_APPEND, 0600)
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
			switch dataSet.columnDescriptors[cellIndex].dataType {
			case smallint:
				{
					buf := bytes.NewBuffer(make([]byte, 0, getDataTypeByteSize(smallint)))
					// binary writer is implicitly converting to uint16 with two's complement
					err := binary.Write(buf, binary.BigEndian, uint16(cell.value.(int)))
					if err != nil {
						return err
					}

					val = buf.Bytes()
				}
			case varchar:
				{
					// we don't care about endianness because we support only utf-8 for now
					padded := make([]byte, getDataTypeByteSize(varchar))
					copy(padded, cell.value.(string))
					val = padded
				}
			case uniqueidentifier:
				{
					val, err = cell.value.(uuid.UUID).MarshalBinary()
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

func readAllFromTable(source SchemeTable[string, string]) (*DataSet, error) {
	tableDescriptor, err := getTableDescriptor(source)
	if err != nil {
		return &DataSet{}, err
	}

	columns := []string{}
	for _, v := range tableDescriptor.columnDescriptors {
		columns = append(columns, v.name)
	}

	return readFromTable(tableDescriptor, columns)
}

func readColumnsFromTable(source SchemeTable[string, string], columns []string) (*DataSet, error) {
	tableDescriptor, err := getTableDescriptor(source)
	if err != nil {
		return &DataSet{}, err
	}

	return readFromTable(tableDescriptor, columns)
}

func readFromTable(tableDescriptor TableDescriptor, selectedColumns []string) (*DataSet, error) {
	f, err := os.Open(getTableDiskPath(tableDescriptor.source))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrTableNotFound
		}

		return nil, err
	}
	defer f.Close()

	fmt.Printf("INFO: %s.%s table descriptor %+v\n", tableDescriptor.source.scheme,
		tableDescriptor.source.name, tableDescriptor)

	dataSet := DataSet{}
	dataSet.columnDescriptors = tableDescriptor.columnDescriptors

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
			if !slices.Contains(selectedColumns, cd.name) {
				rowOffset += getDataTypeByteSize(cd.dataType)
				continue
			}

			switch cd.dataType {
			case smallint:
				{
					cellDataSize = getDataTypeByteSize(smallint)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])

					row.cells = append(row.cells, Cell{
						value: data,
					})
				}
			case varchar:
				{
					cellDataSize = getDataTypeByteSize(varchar)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])

					row.cells = append(row.cells, Cell{
						value: data,
					})
				}
			case uniqueidentifier:
				{
					cellDataSize = getDataTypeByteSize(uniqueidentifier)
					data := make([]byte, cellDataSize)
					copy(data, rowBuf[rowOffset:rowOffset+cellDataSize])

					row.cells = append(row.cells, Cell{
						value: data,
					})

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
