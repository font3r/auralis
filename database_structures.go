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
	schemaTable       SchemaTable[string, string]
	columnDescriptors []ColumnDescriptor // describes table schema
}

type ColumnDescriptor struct {
	name     string
	dataType DataType
	position int
	// attributes eg. PK
}

var auralisTablesTableDescriptor = TableDescriptor{
	schemaTable: SchemaTable[string, string]{internalSchema, tables},
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
	schemaTable: SchemaTable[string, string]{internalSchema, columns},
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

	schemaF, err := os.Create(tablesPath)
	if err != nil {
		panic(err)
	}
	defer schemaF.Close()

	schemaF, err = os.Create(columnsPath)
	if err != nil {
		panic(err)
	}
	defer schemaF.Close()

	err = addAuralisInternalTables()
	if err != nil {
		panic(err)
	}
}

// TODO: implement WHERE clause for table_name/schema
func getTableDescriptor(source SchemaTable[string, string]) (TableDescriptor, error) {
	dataSet, err := readFromTable(auralisColumnsTableDescriptor, SelectQuery{
		source:      SchemaTable[string, string]{internalSchema, tables},
		dataColumns: []string{"table_schema", "table_name", "column_name", "data_type", "position"},
		conditions: []Condition{
			{target: "table_schema", sign: "=", value: source.schema},
			{target: "table_name", sign: "=", value: source.name},
		},
	})
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
		schemaTable:       source,
		columnDescriptors: sourceColumnDescriptors,
	}, nil
}

func addTableDescriptor(tableDescriptor TableDescriptor) error {
	writeIntoTable(auralisTablesTableDescriptor,
		DataSet{
			columnDescriptors: auralisTablesTableDescriptor.columnDescriptors,
			rows: []Row{
				{
					cells: []any{"test-database", tableDescriptor.schemaTable.schema,
						tableDescriptor.schemaTable.name},
				},
			},
		})

	rows := []Row{}
	for _, cd := range tableDescriptor.columnDescriptors {
		rows = append(rows, Row{
			cells: []any{
				tableDescriptor.schemaTable.schema, tableDescriptor.schemaTable.name,
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

func addAuralisInternalTables() error {
	err := writeIntoTable(auralisTablesTableDescriptor,
		DataSet{
			columnDescriptors: auralisTablesTableDescriptor.columnDescriptors,
			rows: []Row{
				{
					cells: []any{"auralis", internalSchema, tables},
				},
			},
		})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisTablesTableDescriptor,
		DataSet{
			columnDescriptors: auralisTablesTableDescriptor.columnDescriptors,
			rows: []Row{
				{
					cells: []any{"auralis", internalSchema, columns},
				},
			},
		})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisColumnsTableDescriptor, DataSet{
		columnDescriptors: auralisColumnsTableDescriptor.columnDescriptors,
		rows: []Row{
			{
				cells: []any{"auralis", "aura_tables", "database_name", string(varchar), 1},
			},
			{
				cells: []any{"auralis", "aura_tables", "table_schema", string(varchar), 2},
			},
			{
				cells: []any{"auralis", "aura_tables", "table_name", string(varchar), 3},
			},
		},
	})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisColumnsTableDescriptor, DataSet{
		columnDescriptors: auralisColumnsTableDescriptor.columnDescriptors,
		rows: []Row{
			{
				cells: []any{"auralis", "aura_columns", "table_schema", string(varchar), 1},
			},
			{
				cells: []any{"auralis", "aura_columns", "table_name", string(varchar), 2},
			},
			{
				cells: []any{"auralis", "aura_columns", "column_name", string(varchar), 3},
			},
			{
				cells: []any{"auralis", "aura_columns", "data_type", string(varchar), 4},
			},
			{
				cells: []any{"auralis", "aura_columns", "position", string(smallint), 5},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
