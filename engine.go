package main

import (
	"errors"
	"log"
)

const defaultScheme = "dbo"

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
	table, err := getTable(query.source)
	if err != nil {
		return &DataSet{}, err
	}

	if len(query.dataColumns) == 1 && query.dataColumns[0] == "*" {
		query.dataColumns = make([]string, 0)
		for _, cd := range table.columns {
			query.dataColumns = append(query.dataColumns, cd.name)
		}
	}

	// TODO: validate conditions, eg. data types
	if len(query.conditions) > 0 {
		for i := range query.conditions {
			err := ConvertConditionType(table, &query.conditions[i])
			if err != nil {
				return &DataSet{}, err
			}
		}
	}

	dataSet, err := readFromTable(table, query)
	if err != nil {
		return &DataSet{}, err
	}

	return dataSet, nil
}

func handleInsertQuery(query InsertQuery) (*DataSet, error) {
	table, err := getTable(query.source)
	if err != nil {
		return &DataSet{}, err
	}

	// TODO: validate provided data against table and cast datatype
	// TODO: handle default values in case of non-null columns (that are also not supported)
	rows := []Row{}
	for _, valueRow := range query.values {
		row := Row{cells: make([]any, 0, len(valueRow))}
		for i, valueCell := range valueRow {
			value, err := ConvertToConcreteType(table.columns[i].dataType, valueCell)
			if err != nil {
				return &DataSet{}, err
			}

			row.cells = append(row.cells, value)
		}

		rows = append(rows, row)
	}

	err = writeIntoTable(table, DataSet{
		columns: table.columns,
		rows:    rows,
	})

	return nil, err
}

func handleCreateTableQuery(query CreateTableQuery) (*DataSet, error) {
	cds := []Column{}
	var i int16 = 1
	for name, attributes := range query.columns {
		if len(attributes) == 0 {
			return nil, errors.New("missing data type for columns")
		}

		// TODO: validate data type
		// if attributes[0] != string(uniqueidentifier) {

		// }

		cds = append(cds, Column{
			name:     name,
			dataType: DataType(attributes[0]),
			position: i,
		})

		i++
	}

	err := cretateTable(Table{
		schemaTable: query.source,
		columns:     cds,
	})

	if err != nil {
		panic(err)
	}

	return nil, nil
}
