package main

import (
	"fmt"
)

func ExecuteQuery(raw string) (*DataSet, error) {
	tokens := Analyze(raw)
	if len(tokens) <= 0 {
		return &DataSet{}, AuraError{
			Code:    "INVALID_QUERY",
			Message: "missing query tokens"}
	}

	fmt.Printf("INFO: tokens %v\n", tokens)

	query, err := ParseTokens(tokens)
	if err != nil {
		return &DataSet{}, err
	}

	fmt.Printf("INFO: parsed query %v\n", query)

	switch query := query.(type) {
	case SelectQuery:
		{
			var dataSet *DataSet
			var err error

			if len(query.columns) == 1 && query.columns[0] == "*" {
				dataSet, err = readAllFromTable(query.source)
			} else {
				dataSet, err = readColumnsFromTable(query.source, query.columns)
			}

			if err != nil {
				return &DataSet{}, err
			}

			return dataSet, nil
		}
	default:
		{
			panic("unsupported query")
		}
	}
}
