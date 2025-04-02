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
			cells, err := readTableRow(query.source, query.columns)
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
