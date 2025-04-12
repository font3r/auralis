package main

import (
	"os"
)

const (
	dataPath string = "./data"

	internalSchema string = "auralis"
	tables         string = "aura_tables"
	columns        string = "aura_columns"

	tablesPath  string = dataPath + "/" + internalSchema + "." + tables
	columnsPath string = dataPath + "/" + internalSchema + "." + columns
)

var (
	ErrTableDescriptorNotFound = AuraError{Code: "TABLE_DESCRIPTOR_NOT_FOUND", Message: "table descriptor not found"}
)

type TableDescriptor struct {
	schemeTable       SchemeTable[string, string]
	columnDescriptors []ColumnDescriptor // describes table schema
}

type ColumnDescriptor struct {
	name     string
	dataType DataType
	position int
	// attributes eg. PK
}

var auralisTablesTableDescriptor = TableDescriptor{
	schemeTable: SchemeTable[string, string]{internalSchema, tables},
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
}

var auralisColumnsTableDescriptor = TableDescriptor{
	schemeTable: SchemeTable[string, string]{internalSchema, columns},
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
}

func initDatabaseInternalStructure() {
	if _, err := os.Stat(dataPath); !os.IsNotExist(err) {
		return
	}

	err := os.Mkdir(dataPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	schemeF, err := os.Create(tablesPath)
	if err != nil {
		panic(err)
	}
	defer schemeF.Close()

	schemeF, err = os.Create(columnsPath)
	if err != nil {
		panic(err)
	}
	defer schemeF.Close()
}

func getAuralisTables() (*DataSet, error) {
	return readFromTable(auralisTablesTableDescriptor,
		[]string{"database_name", "table_schema", "table_name"})
}

func getAuralisColumns() (*DataSet, error) {
	return readFromTable(auralisColumnsTableDescriptor,
		[]string{"table_schema", "table_name", "column_name", "data_type", "position"})
}

// TODO: implement WHERE clause for table_name/scheme
func getTableDescriptor(source SchemeTable[string, string]) (TableDescriptor, error) {
	dataSet, err := readFromTable(auralisColumnsTableDescriptor,
		[]string{"table_schema", "table_name", "column_name", "data_type", "position"})
	if err != nil {
		return TableDescriptor{}, err
	}

	sourceColumnDescriptors := []ColumnDescriptor{}
	for _, row := range dataSet.rows {
		sourceColumnDescriptor := ColumnDescriptor{}
		for i, cell := range row.cells {
			if dataSet.columnDescriptors[i].name == "column_name" {
				sourceColumnDescriptor.name = cell.(string)
			}

			if dataSet.columnDescriptors[i].name == "data_type" {
				sourceColumnDescriptor.dataType = DataType(cell.(string))
			}

			if dataSet.columnDescriptors[i].name == "position" {
				sourceColumnDescriptor.position = int(cell.(int16))
			}
		}

		sourceColumnDescriptors = append(sourceColumnDescriptors, sourceColumnDescriptor)
	}

	return TableDescriptor{
		schemeTable:       source,
		columnDescriptors: sourceColumnDescriptors,
	}, nil
}

func addTableDescriptor(tableDescriptor TableDescriptor) error {
	writeIntoTable(auralisTablesTableDescriptor,
		DataSet{
			columnDescriptors: auralisTablesTableDescriptor.columnDescriptors,
			rows: []Row{
				{
					cells: []any{"test-database", tableDescriptor.schemeTable.scheme,
						tableDescriptor.schemeTable.name},
				},
			},
		})

	rows := []Row{}
	for _, cd := range tableDescriptor.columnDescriptors {
		rows = append(rows, Row{
			cells: []any{
				tableDescriptor.schemeTable.scheme, tableDescriptor.schemeTable.name,
				cd.name, string(cd.dataType), cd.position,
			},
		})
	}

	writeIntoTable(auralisColumnsTableDescriptor, DataSet{
		columnDescriptors: auralisColumnsTableDescriptor.columnDescriptors,
		rows:              rows,
	})

	return nil
}
