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
	source            SchemeTable[string, string]
	columnDescriptors []ColumnDescriptor // describes table schema
}

type ColumnDescriptor struct {
	name     string
	dataType DataType
	position int
	// attributes eg. PK
}

var auralisTablesTableDescriptor = TableDescriptor{
	source: SchemeTable[string, string]{internalSchema, tables},
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
	source: SchemeTable[string, string]{internalSchema, columns},
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
	err := os.Mkdir(dataPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	schemeF, err := os.Create(tables)
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
				sourceColumnDescriptor.name = cell.value.(string)
			}

			if dataSet.columnDescriptors[i].name == "data_type" {
				sourceColumnDescriptor.dataType = DataType(cell.value.(string))
			}

			if dataSet.columnDescriptors[i].name == "position" {
				sourceColumnDescriptor.position = int(cell.value.(int16))
			}
		}

		sourceColumnDescriptors = append(sourceColumnDescriptors, sourceColumnDescriptor)
	}

	return TableDescriptor{
		source:            source,
		columnDescriptors: sourceColumnDescriptors,
	}, nil
}

func addTableDescriptor(tableDescriptor TableDescriptor) error {
	writeIntoTable(SchemeTable[string, string]{internalSchema, tables}, DataSet{
		columnDescriptors: auralisTablesTableDescriptor.columnDescriptors,
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
		columnDescriptors: auralisColumnsTableDescriptor.columnDescriptors,
		rows:              rows,
	})

	return nil
}
