package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	tables string = "aura_tables"

	internalDirPath string = "./data/internal"
	schemeInfoFile  string = internalDirPath + "/" + tables
)

var (
	ErrTableDescriptorNotFound = AuraError{Code: "TABLE_DESCRIPTOR_NOT_FOUND", Message: "table descriptor not found"}
)

type TableDescriptor struct {
	source            SchemeTable[string, string]
	columnDescriptors []ColumnDescriptor // describes table schema
}

type ColumnDescriptor struct {
	name     string
	dataType DataType
	// position int // columnd order
	// attributes eg. PK
}

func initDatabaseInternalStructure() {
	err := os.MkdirAll(internalDirPath, os.ModePerm) // TODO: verify perm bits
	if err != nil {
		panic(err)
	}

	err = createInformationSchemaTable()
	if err != nil {
		panic(err)
	}
}

func createInformationSchemaTable() error {
	schemeF, err := os.Create(schemeInfoFile)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	return nil
}

func getTableDescriptor(source SchemeTable[string, string]) (TableDescriptor, error) {
	fileBytes, err := os.ReadFile(schemeInfoFile)
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
			if td.source.name == "" && td.source.scheme == "" {
				td.source = SchemeTable[string, string]{parts[0], parts[1]}
			} else {
				td.columnDescriptors = append(td.columnDescriptors, ColumnDescriptor{
					name:     parts[0],
					dataType: DataType(parts[1]),
				})
			}

			l = i + 1
		}

		if fileBytes[i] == byte(terminationByte) {
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
		if td.source.name == source.name && td.source.scheme == source.scheme {
			exists = td
			return exists, nil
		}
	}

	return TableDescriptor{}, ErrTableDescriptorNotFound
}

func addTableDescriptor(tableDescriptor TableDescriptor) error {
	schemeFile, err := os.OpenFile(schemeInfoFile, os.O_RDWR, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrTableNotFound
		}
		return err
	}
	defer schemeFile.Close()

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s.%s", tableDescriptor.source.scheme, tableDescriptor.source.name))

	for _, cd := range tableDescriptor.columnDescriptors {
		b.WriteString(fmt.Sprintf("|%s.%s", cd.name, cd.dataType))
	}

	b.WriteString("\n")

	_, err = schemeFile.WriteString(b.String())
	if err != nil {
		return err
	}

	return nil
}
