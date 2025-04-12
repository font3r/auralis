package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ExecuteQuery(raw string) (*DataSet, error) {
	tokens := Analyze(raw)
	if len(tokens) <= 0 {
		return &DataSet{}, AuraError{
			Code:    "INVALID_QUERY",
			Message: "missing query tokens"}
	}

	log.Printf("INFO: tokens %v\n", tokens)

	query, err := ParseTokens(tokens)
	if err != nil {
		return &DataSet{}, err
	}

	log.Printf("INFO: parsed query %v\n", query)

	switch query := query.(type) {
	case SelectQuery:
		{
			tableDescriptor, err := getTableDescriptor(query.source)
			if err != nil {
				return &DataSet{}, err
			}

			var dataSet *DataSet
			if len(query.columns) == 1 && query.columns[0] == "*" {
				dataSet, err = readAllFromTable(tableDescriptor)
			} else {
				dataSet, err = readColumnsFromTable(tableDescriptor, query.columns)
			}

			if err != nil {
				return &DataSet{}, err
			}

			return dataSet, nil
		}
	case InsertQuery:
		{
			tableDescriptor, err := getTableDescriptor(query.destination)
			if err != nil {
				return &DataSet{}, err
			}

			// TODO: validate provided data against table descriptor and cast datatype
			// TODO: handle default values in case of non-null columns (that are also not supported)
			rows := []Row{}
			for _, valueRow := range query.values {
				row := Row{cells: make([]any, 0, len(valueRow))}
				for i, valueCell := range valueRow {
					switch tableDescriptor.columnDescriptors[i].dataType {
					case smallint:
						{
							v, err := strconv.Atoi(valueCell.(string))
							if err != nil {
								return nil, errors.New("invalid smallint cell value")
							}

							row.cells = append(row.cells, v)
						}
					case varchar:
						{
							row.cells = append(row.cells, strings.Trim(valueCell.(string), "'"))
						}
					case uniqueidentifier:
						{
							v, err := uuid.Parse(valueCell.(string))
							if err != nil {
								return nil, errors.New("invalid UUID cell value")
							}

							row.cells = append(row.cells, v)
						}
					default:
						panic("unhandled type during insert query")
					}
				}

				rows = append(rows, row)
			}

			err = writeIntoTable(tableDescriptor, DataSet{
				columnDescriptors: tableDescriptor.columnDescriptors,
				rows:              rows,
			})

			return nil, err
		}
	default:
		{
			panic("unsupported query")
		}
	}
}
