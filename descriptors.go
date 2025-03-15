package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
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

const (
	smallint         DataType = "smallint" // 2
	integer          DataType = "integer"  // 4
	bigint           DataType = "bigint"   // 8
	varchar          DataType = "varchar"
	uniqueidentifier DataType = "uniqueidentifier" // 16
	boolean          DataType = "boolean"          // 1
)

func cretateTable(td TableDescriptor) error {
	schemeF, err := os.OpenFile(fmt.Sprintf("./data/%s", schemeStore), os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	_, err = schemeF.WriteString(fmt.Sprintf("%s\n", getTableDescriptorAsString(td)))
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
				err := binary.Write(buf, binary.BigEndian, int16(v.value.(int)))
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

func readTableRow(scheme, name string) ([]Cell, error) {
	f, err := os.Open(fmt.Sprintf("./data/%s.%s", scheme, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	storeTd, _ := getTableDescriptorFromStore(scheme, name)

	fmt.Printf("INFO: %s.%s table descriptor %s\n", scheme, name, storeTd)
	cells := make([]Cell, 0)
	buf := make([]byte, 1024) // buffer should be calculated based od schema row + delimeters etc

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		for i, cd := range storeTd.columnDescriptors {
			switch cd.dataType {
			case smallint:
				{
					size := getDataTypeByteSize(smallint)
					cells = append(cells, Cell{
						columnDescriptor: cd,
						value:            buf[i*size : (i+1)*size],
					})
				}
			case uniqueidentifier:
				{
					size := getDataTypeByteSize(uniqueidentifier)
					cells = append(cells, Cell{
						columnDescriptor: cd,
						value:            buf[i*size : (i+1)*size],
					})
				}
			default:
				return []Cell{}, errors.New("unhandled type")
			}
		}
	}

	return cells, nil
}

func getTableDescriptorAsString(td TableDescriptor) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s.%s|", td.scheme, td.name))

	for _, cd := range td.columnDescriptors {
		b.WriteString(fmt.Sprintf("%s.%s|", cd.name, cd.dataType))
	}

	return b.String()
}

func getTableDescriptorFromStore(scheme, name string) (TableDescriptor, error) {
	// TODO: read table structure from store (maybe should be table as other)

	return TableDescriptor{
		name:   name,
		scheme: scheme,
		columnDescriptors: []ColumnDescriptor{
			{name: "id1", dataType: smallint},
			{name: "id2", dataType: uniqueidentifier},
			{name: "id3", dataType: smallint},
		},
	}, nil
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
