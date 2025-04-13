package main

import (
	"log"
)

func ExecuteQuery(raw string) (*DataSet, error) {
	tokens := Analyze(raw)
	if len(tokens) <= 0 {
		return &DataSet{}, AuraError{
			Code:    "INVALID_QUERY",
			Message: "missing query tokens"}
	}

	log.Printf("INFO: lexer tokens %v\n", tokens)

	query, err := ParseTokens(tokens)
	if err != nil {
		return &DataSet{}, err
	}

	log.Printf("INFO: parsed query %+v\n", query)

	switch query := query.(type) {
	case SelectQuery:
		return handleSelectQuery(query)
	case InsertQuery:
		return handleInsertQuery(query)
	case CreateTableQuery:
		return handleCreateTableQuery(query)
	default:
		panic("unsupported query")
	}
}

func handleSelectQuery(query SelectQuery) (*DataSet, error) {
	tableDescriptor, err := getTableDescriptor(query.source)
	if err != nil {
		return &DataSet{}, err
	}

	if len(query.dataColumns) == 1 && query.dataColumns[0] == "*" {
		query.dataColumns = make([]string, 0)
		for _, cd := range tableDescriptor.columnDescriptors {
			query.dataColumns = append(query.dataColumns, cd.name)
		}
	}

	// TODO: validate conditions, eg. data types
	if len(query.conditions) > 0 {
		for i := range query.conditions {
			err := ConvertConditionType(tableDescriptor, &query.conditions[i])
			if err != nil {
				return &DataSet{}, err
			}
		}
	}

	dataSet, err := readFromTable(tableDescriptor, query)
	if err != nil {
		return &DataSet{}, err
	}

	return dataSet, nil
}

func handleInsertQuery(query InsertQuery) (*DataSet, error) {
	tableDescriptor, err := getTableDescriptor(query.source)
	if err != nil {
		return &DataSet{}, err
	}

	// TODO: validate provided data against table descriptor and cast datatype
	// TODO: handle default values in case of non-null columns (that are also not supported)
	rows := []Row{}
	for _, valueRow := range query.values {
		row := Row{cells: make([]any, 0, len(valueRow))}
		for i, valueCell := range valueRow {
			value, err := ConvertToConcreteType(tableDescriptor.columnDescriptors[i].dataType, valueCell)
			if err != nil {
				return &DataSet{}, err
			}

			row.cells = append(row.cells, value)
		}
	}

	err = writeIntoTable(tableDescriptor, DataSet{
		columnDescriptors: tableDescriptor.columnDescriptors,
		rows:              rows,
	})

	return nil, err
}

func handleCreateTableQuery(query CreateTableQuery) (*DataSet, error) {
	// create test table
	err := cretateTable(TableDescriptor{
		schemeTable: SchemeTable[string, string]{"dbo", "users"},
		columnDescriptors: []ColumnDescriptor{
			{
				name:     "name",
				dataType: varchar,
				position: 1,
			},
			{
				name:     "age",
				dataType: smallint,
				position: 2,
			},
		},
	})

	if err != nil {
		panic(err)
	}

	return nil, nil
}
