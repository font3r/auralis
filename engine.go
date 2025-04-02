package main

import (
	"fmt"
	"os"
)

const (
	schemaStore string = "schema_store"
)

func ExecuteQuery(raw string) ([]Cell, error) {
	tokens := Analyze(raw)
	if len(tokens) <= 0 {
		return []Cell{}, AuraError{
			Code:    "INVALID_QUERY",
			Message: "missing query tokens"}
	}

	fmt.Printf("INFO: tokens %v\n", tokens)

	query, err := ParseTokens(tokens)
	if err != nil {
		return []Cell{}, err
	}

	fmt.Printf("INFO: parsed query %v\n", query)

	switch query := query.(type) {
	case SelectQuery:
		{
			var cells []Cell
			var err error

			if len(query.columns) == 1 && query.columns[0] == "*" {
				cells, err = readAllFromTable(query.source)
			} else {
				cells, err = readColumnsFromTable(query.source, query.columns)
			}

			if err != nil {
				return []Cell{}, err
			}

			return cells, nil
		}
	default:
		{
			panic("unsupported query")
		}
	}
}

func CreateSchemaStore() error {
	schemeF, err := os.OpenFile(fmt.Sprintf("./data/%s", schemaStore), os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	return nil
}
