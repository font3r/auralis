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
	"strings"

	"github.com/google/uuid"
)

type TableDescriptor struct {
	name              string
	scheme            string             // eg. dbo
	columnDescriptors []ColumnDescriptor // describes table schema
}

type ColumnDescriptor struct {
	name     string
	dataType DataType
	// attributes eg. PK
}

type Cell struct {
	columnDescriptor ColumnDescriptor
	value            any
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

func cretateTable(td TableDescriptor) error {
	schemeF, err := os.OpenFile(fmt.Sprintf("./data/%s", schemaStore), os.O_RDWR, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTableNotFound
		}
		return err
	}
	defer schemeF.Close()

	_, err = schemeF.WriteString(fmt.Sprintf("%s\n", tableDescriptorAsString(td)))
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./data/%s.%s", td.scheme, td.name))
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func writeIntoTable(scheme, name string, data ...Cell) error {
	f, err := os.OpenFile(fmt.Sprintf("./data/%s.%s", scheme, name), os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTableNotFound
		}
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, v := range data {
		var val []byte
		switch v.columnDescriptor.dataType {
		case smallint:
			{
				buf := bytes.NewBuffer(make([]byte, 0, getDataTypeByteSize(smallint)))
				// binary writer is implicitly converting to uint16 with two's complement
				err := binary.Write(buf, binary.BigEndian, uint16(v.value.(int)))
				if err != nil {
					return err
				}

				val = buf.Bytes()
			}
		case uniqueidentifier:
			{
				val, err = v.value.(uuid.UUID).MarshalBinary()
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

	// row termination byte
	err = w.WriteByte(10)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func readAllFromTable(source SchemeTable[string, string]) ([]Cell, error) {
	tableDescriptor, err := getTableDescriptor(source)
	if err != nil {
		return []Cell{}, err
	}

	columns := []string{}
	for _, v := range tableDescriptor.columnDescriptors {
		columns = append(columns, v.name)
	}

	return readFromTable(tableDescriptor, columns)
}

func readColumnsFromTable(source SchemeTable[string, string], columns []string) ([]Cell, error) {
	tableDescriptor, err := getTableDescriptor(source)
	if err != nil {
		return []Cell{}, err
	}

	return readFromTable(tableDescriptor, columns)
}

func readFromTable(tableDescriptor TableDescriptor, columns []string) ([]Cell, error) {
	f, err := os.Open(fmt.Sprintf("./data/%s.%s", tableDescriptor.scheme, tableDescriptor.name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrTableNotFound
		}

		return nil, err
	}
	defer f.Close()

	fmt.Printf("INFO: %s.%s table descriptor %+v\n", tableDescriptor.scheme, tableDescriptor.name, tableDescriptor)
	cells := make([]Cell, 0)
	buf := make([]byte, 1024) // TODO: buffer should be calculated based od schema row + delimeters etc

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		offset := 0
		for _, cd := range tableDescriptor.columnDescriptors {
			if !slices.Contains(columns, cd.name) {
				offset += getDataTypeByteSize(cd.dataType)
				continue
			}

			switch cd.dataType {
			case smallint:
				{
					size := getDataTypeByteSize(smallint)
					cells = append(cells, Cell{
						columnDescriptor: cd,
						value:            buf[offset : offset+size],
					})

					offset += size
				}
			case uniqueidentifier:
				{
					size := getDataTypeByteSize(uniqueidentifier)
					cells = append(cells, Cell{
						columnDescriptor: cd,
						value:            buf[offset : offset+size],
					})

					offset += size
				}
			default:
				return []Cell{}, errors.New("unhandled type")
			}
		}
	}

	return cells, nil
}

func tableDescriptorAsString(td TableDescriptor) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s.%s", td.scheme, td.name))

	for _, cd := range td.columnDescriptors {
		b.WriteString(fmt.Sprintf("|%s.%s", cd.name, cd.dataType))
	}

	return b.String()
}

func getTableDescriptor(schemeTable SchemeTable[string, string]) (TableDescriptor, error) {
	fileBytes, err := os.ReadFile(fmt.Sprintf("./data/%s", schemaStore))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return TableDescriptor{}, ErrTableNotFound
		}
		return TableDescriptor{}, err
	}

	tableDescriptors := make([]TableDescriptor, 0)
	l := 0
	td := TableDescriptor{}

	for i := range fileBytes {
		if fileBytes[i] == byte('|') {
			parts := strings.Split(string(fileBytes[l:i]), ".")

			if td.name == "" && td.scheme == "" {
				td.scheme = parts[0]
				td.name = parts[1]
			} else {
				td.columnDescriptors = append(td.columnDescriptors, ColumnDescriptor{
					name:     parts[0],
					dataType: DataType(parts[1]),
				})
			}

			l = i + 1
		}

		if fileBytes[i] == byte(10) {
			parts := strings.Split(string(fileBytes[l:i]), ".")
			td.columnDescriptors = append(td.columnDescriptors, ColumnDescriptor{
				name:     parts[0],
				dataType: DataType(parts[1]),
			})

			tableDescriptors = append(tableDescriptors, td)
		}
	}

	var exists TableDescriptor
	for _, td := range tableDescriptors {
		if td.name == schemeTable.name && td.scheme == schemeTable.scheme {
			exists = td
			return exists, nil
		}
	}

	return TableDescriptor{}, AuraError{
		Code:    "TABLE_DESCRIPTOR_NOT_FOUND",
		Message: "table descriptor for given scheme name foes not exist"}
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
		return 8 // TODO: this should be configurable
	case uniqueidentifier:
		return 16
	case boolean:
		return 1
	default:
		panic("unhandled type")
	}
}
