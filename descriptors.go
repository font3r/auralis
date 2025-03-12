package main

import (
	"bufio"
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
	smallint         DataType = "smallint"
	integer          DataType = "integer"
	bigint           DataType = "bigint"
	varchar          DataType = "varchar"
	uniqueidentifier DataType = "uniqueidentifier"
	boolean          DataType = "boolean"
)

func cretateTable(td TableDescriptor) error {
	schemeF, err := os.OpenFile(fmt.Sprintf("./data/%s", schemeStore), os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	_, err = schemeF.WriteString(getTableDescriptorAsString(td))
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
		case uniqueidentifier:
			{
				val, _ = v.value.(uuid.UUID).MarshalBinary()
			}
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

	tableDescriptor := buildTableDescriptor(scheme, name)
	cells := make([]Cell, 10)
	buf := make([]byte, 1024) // buffer should be calculated based od schema row
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		for _, cd := tableDescriptor.columnDescriptors {
			
		}

		// two uuids
		cells[0].value = buf[0:16]
		cells[1].value = buf[16:32]
	}

	return nil
}

func getTableDescriptorAsString(td TableDescriptor) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s.%s|", td.scheme, td.name))

	for _, cd := range td.columnDescriptors {
		b.WriteString(fmt.Sprintf("%s.%s|", cd.name, cd.dataType))
	}

	return b.String()
}

func buildTableDescriptor(scheme, name string) TableDescriptor {
	return TableDescriptor{
		name:   name,
		scheme: scheme,
		columnDescriptors: []ColumnDescriptor{
			{name: "id1", dataType: uniqueidentifier},
			{name: "id2", dataType: uniqueidentifier},
		},
	}
}

func getTableRowByteSize(td TableDescriptor) int {
	x := 0
	for _, cd := range td.columnDescriptors {
		x += getDataTypeSize(cd.dataType)
	}

	return x
}

func getDataTypeSize(dataType DataType) int {
	switch dataType {
	case smallint:
		return 2 // int16
	case integer:
		return 4 // int32
	case bigint:
		return 8 // int64
	case varchar:
		return 8 // <- TOOD: this should be configurable
	case uniqueidentifier:
		return 16
	case boolean:
		return 1
	}

	panic("ERROR: unsupported datatype")
}
