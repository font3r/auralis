package main

import (
	"os"
)

const (
	dataPath string = "./data"

	internalSchema string = "auralis"
	tables         string = "tables"
	columns        string = "columns"

	tablesPath  string = dataPath + "/" + internalSchema + "." + tables
	columnsPath string = dataPath + "/" + internalSchema + "." + columns
)

type Table struct {
	schemaTable SchemaTable[string, string]
	columns     []Column // describes table schema
}

type Column struct {
	name     string
	dataType DataType
	position int16
	// attributes eg. PK
}

var auralisTables = Table{
	schemaTable: SchemaTable[string, string]{internalSchema, tables},
	columns: []Column{
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

var auralisColumnsTable = Table{
	schemaTable: SchemaTable[string, string]{internalSchema, columns},
	columns: []Column{
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

func addAuralisInternalTables() error {
	err := writeIntoTable(auralisTables,
		DataSet{
			columns: auralisTables.columns,
			rows: []Row{
				{
					cells: []any{"auralis", internalSchema, tables},
				},
			},
		})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisTables,
		DataSet{
			columns: auralisTables.columns,
			rows: []Row{
				{
					cells: []any{"auralis", internalSchema, columns},
				},
			},
		})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisColumnsTable, DataSet{
		columns: auralisColumnsTable.columns,
		rows: []Row{
			{
				cells: []any{"auralis", "tables", "database_name", string(varchar), int16(1)},
			},
			{
				cells: []any{"auralis", "tables", "table_schema", string(varchar), int16(2)},
			},
			{
				cells: []any{"auralis", "tables", "table_name", string(varchar), int16(3)},
			},
		},
	})

	if err != nil {
		return err
	}

	err = writeIntoTable(auralisColumnsTable, DataSet{
		columns: auralisColumnsTable.columns,
		rows: []Row{
			{
				cells: []any{"auralis", "columns", "table_schema", string(varchar), int16(1)},
			},
			{
				cells: []any{"auralis", "columns", "table_name", string(varchar), int16(2)},
			},
			{
				cells: []any{"auralis", "columns", "column_name", string(varchar), int16(3)},
			},
			{
				cells: []any{"auralis", "columns", "data_type", string(varchar), int16(4)},
			},
			{
				cells: []any{"auralis", "columns", "position", string(smallint), int16(5)},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func getTable(source SchemaTable[string, string]) (Table, error) {
	dataSet, err := readFromTable(auralisColumnsTable, SelectQuery{
		source:      SchemaTable[string, string]{internalSchema, tables},
		dataColumns: []string{"table_schema", "table_name", "column_name", "data_type", "position"},
		conditions: []Condition{
			{target: "table_schema", sign: "=", value: source.schema},
			{target: "table_name", sign: "=", value: source.name},
		},
	})
	if err != nil {
		return Table{}, err
	}

	sourceColumns := []Column{}
	for _, row := range dataSet.rows {
		sourceColumn := Column{}
		for i, cell := range row.cells {
			if dataSet.columns[i].name == "column_name" {
				sourceColumn.name = cell.(string)
			}

			if dataSet.columns[i].name == "data_type" {
				sourceColumn.dataType = DataType(cell.(string))
			}

			if dataSet.columns[i].name == "position" {
				sourceColumn.position = cell.(int16)
			}
		}

		sourceColumns = append(sourceColumns, sourceColumn)
	}

	return Table{
		schemaTable: source,
		columns:     sourceColumns,
	}, nil
}

func addTable(table Table) error {
	writeIntoTable(auralisTables,
		DataSet{
			columns: auralisTables.columns,
			rows: []Row{
				{
					cells: []any{"test-database", table.schemaTable.schema,
						table.schemaTable.name},
				},
			},
		})

	rows := []Row{}
	for _, cd := range table.columns {
		rows = append(rows, Row{
			cells: []any{
				table.schemaTable.schema, table.schemaTable.name,
				cd.name, string(cd.dataType), cd.position,
			},
		})
	}

	writeIntoTable(auralisColumnsTable, DataSet{
		columns: auralisColumnsTable.columns,
		rows:    rows,
	})

	return nil
}
