package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	internalSchema string = "auralis"
	tables         string = "aura_tables"
	columns        string = "aura_columns"

	internalDirPath   string = "./data"
	tableSchemeFile   string = internalDirPath + "/" + tables
	columnsSchemeFile string = internalDirPath + "/" + columns
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
	position int
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

	schemeF, err := os.Create("data/auralis.aura_tables")
	if err != nil {
		panic(err)
	}
	defer schemeF.Close()

	schemeF, err = os.Create("data/auralis.aura_columns")
	if err != nil {
		panic(err)
	}
	defer schemeF.Close()
}

func getAuralisTables() (*DataSet, error) {
	return readFromTable(TableDescriptor{
		source: SchemeTable[string, string]{"auralis", tables},
		columnDescriptors: []ColumnDescriptor{
			{
				name:     "database_name",
				dataType: varchar,
				position: 1,
			},
			{
				name:     "table_schema",
				dataType: varchar,
				position: 2,
			},
			{
				name:     "table_name",
				dataType: varchar,
				position: 3,
			},
		},
	}, []string{"database_name", "table_schema", "table_name"})
}

func getAuralisColumns() (*DataSet, error) {
	return readFromTable(TableDescriptor{
		source: SchemeTable[string, string]{"auralis", columns},
		columnDescriptors: []ColumnDescriptor{
			{
				name:     "table_schema",
				dataType: varchar,
			},
			{
				name:     "table_name",
				dataType: varchar,
			},
			{
				name:     "column_name",
				dataType: varchar,
			},
			{
				name:     "data_type",
				dataType: varchar,
			},
			{
				name:     "position",
				dataType: smallint,
			},
		},
	}, []string{"table_schema", "table_name", "column_name", "data_type", "position"})
}

func createInformationSchemaTable() error {
	schemeF, err := os.Create(tableSchemeFile)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	return nil
}

func getTableDescriptor(source SchemeTable[string, string]) (TableDescriptor, error) {
	fileBytes, err := os.ReadFile(tableSchemeFile)
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
	schemeFile, err := os.OpenFile(tableSchemeFile, os.O_RDWR, 0600)
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

	writeIntoTable(SchemeTable[string, string]{internalSchema, tables}, DataSet{
		columnDescriptors: []ColumnDescriptor{
			{
				name:     "database_name",
				dataType: varchar,
				position: 1,
			},
			{
				name:     "table_schema",
				dataType: varchar,
				position: 2,
			},
			{
				name:     "table_name",
				dataType: varchar,
				position: 3,
			},
		},
		rows: []Row{
			{
				cells: []Cell{{"test-database"}, {tableDescriptor.source.scheme}, {tableDescriptor.source.name}},
			},
		},
	})

	rows := []Row{}
	for _, cd := range tableDescriptor.columnDescriptors {
		rows = append(rows, Row{
			cells: []Cell{
				{tableDescriptor.source.scheme}, {tableDescriptor.source.name}, {cd.name}, {string(cd.dataType)}, {cd.position},
			},
		})
	}

	writeIntoTable(SchemeTable[string, string]{internalSchema, columns}, DataSet{
		columnDescriptors: []ColumnDescriptor{
			{
				name:     "table_schema",
				dataType: varchar,
			},
			{
				name:     "table_name",
				dataType: varchar,
			},
			{
				name:     "column_name",
				dataType: varchar,
			},
			{
				name:     "data_type",
				dataType: varchar,
			},
			{
				name:     "position",
				dataType: smallint,
			},
		},
		rows: rows,
	})

	return nil
}
