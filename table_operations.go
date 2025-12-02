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
	columns []Column
	rows    []Row
}

type Row struct {
	cells []any
}

var (
	ErrTableNotFound = AuraError{Code: "TABLE_NOT_FOUND", Message: "table not found"}
)

const terminationByte = byte(10)

func cretateTable(table Table) error {
	err := addTable(table)
	if err != nil {
		return err
	}

	f, err := os.Create(getTableDiskPath(table.schemaTable))
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func writeIntoTable(table Table, dataSet DataSet) error {
	log.Printf("INFO: executing insert query %+v", dataSet)
	f, err := os.OpenFile(getTableDiskPath(table.schemaTable), os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTableNotFound
		}
		return err
	}
	defer f.Close()

	log.Printf("INFO: %s.%s table %+v\n", table.schemaTable.schema, table.schemaTable.name, table)

	w := bufio.NewWriter(f)

	for _, row := range dataSet.rows {
		for cellIndex, cell := range row.cells {
			var val []byte
			switch table.columns[cellIndex].dataType {
			case smallint:
				{
					buf := bytes.NewBuffer(make([]byte, 0, getDataTypeByteSize(smallint)))
					// binary writer is implicitly converting to uint16 with two's complement
					err := binary.Write(buf, binary.BigEndian, cell.(int16))
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

func readFromTable(table Table, query SelectQuery) (*DataSet, error) {
	log.Printf("INFO: executing select query %+v", query)
	f, err := os.Open(getTableDiskPath(table.schemaTable))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrTableNotFound
		}

		return nil, err
	}
	defer f.Close()

	log.Printf("INFO: %s.%s table %+v\n", table.schemaTable.schema,
		table.schemaTable.name, table)

	dataSet := DataSet{}
	for _, v := range table.columns {
		if slices.Contains(query.dataColumns, v.name) {
			dataSet.columns = append(dataSet.columns, v)
		}
	}

	var fileOffset int64 = 0
	rowBuffSize := calculateRowBuffer(table)
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
		includeRow := true

		var cellDataSize int
		for _, cd := range table.columns {
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
						} else {
							includeRow = false
							break
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
					value := string(bytes.TrimRight(data, "\x00"))

					conditions := GetMatchingCondition(query.conditions, cd.name)
					if len(conditions) > 0 {
						// for now assume only one condition
						if EvaluateStringCondition(conditions[0], value) {
							row.cells = append(row.cells, value)
						} else {
							includeRow = false
							break
						}
					} else {
						row.cells = append(row.cells, value)
					}
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

		if includeRow {
			dataSet.rows = append(dataSet.rows, row)
		}

		fileOffset += int64(rowBuffSize)
		clear(rowBuf)
	}

	return &dataSet, nil
}

func calculateRowBuffer(table Table) int {
	size := 0
	for _, v := range table.columns {
		size += getDataTypeByteSize(v.dataType)
	}
	size += 1 // termination byte

	return size
}

func getTableDiskPath(source SchemaTable[string, string]) string {
	return fmt.Sprintf("./data/%s.%s", source.schema, source.name)
}
