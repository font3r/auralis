package main

import (
	"errors"
	"log"
	"strings"
)

type Command struct{}

type SchemaTable[T, U string] struct {
	schema T
	name   U
}

type SelectQuery struct {
	source      SchemaTable[string, string]
	dataColumns []string
	conditions  []Condition
}

type InsertQuery struct {
	source      SchemaTable[string, string]
	dataColumns []string // column names
	values      [][]any  // column values
}

type CreateTableQuery struct {
	source  SchemaTable[string, string]
	columns map[string][]string
}

func ParseTokens(tokens []TokenLiteral) (any, error) {
	valid := hasAnyKeyword(&tokens)
	if !valid {
		return Command{}, errors.New("missing any keyword")
	}

	// TODO: detect which query type to analyze
	switch tokens[0].value {
	case "select":
		return parseSelect(&tokens)
	case "insert":
		return parseInsert(&tokens)
	case "create":
		return parseCreate(&tokens)
	}

	return Command{}, errors.New("unsupported keyword")
}

func hasAnyKeyword(tokens *[]TokenLiteral) bool {
	valid := false
	for _, v := range *tokens {
		if v.kind != keyword {
			continue
		}

		valid = true
		break
	}

	return valid
}

func parseSelect(tokens *[]TokenLiteral) (SelectQuery, error) {
	v := *tokens
	q := SelectQuery{}
	i := 0

	// select
	if i >= len(v) || v[i].kind != keyword || v[i].value != "select" {
		return SelectQuery{}, errors.New("missing select keyword")
	}
	i++

	// * or csv columns
	for ; i < len(v); i++ {
		if v[i].kind == comma {
			continue
		}

		if v[i].kind != symbol {
			break
		}

		q.dataColumns = append(q.dataColumns, v[i].value)
	}

	if len(q.dataColumns) == 0 {
		return SelectQuery{}, errors.New("missing columns")
	}
	// i is incremented by loop

	// from
	if i >= len(v) || v[i].kind != keyword || v[i].value != "from" {
		return SelectQuery{}, errors.New("missing from keyword")
	}
	i++

	// source table
	if i >= len(v) || v[i].kind != symbol {
		return SelectQuery{}, errors.New("missing source table")
	}

	s := strings.Split(v[i].value, ".")
	if len(s) == 1 {
		q.source = SchemaTable[string, string]{defaultScheme, s[0]}
	} else {
		q.source = SchemaTable[string, string]{s[0], s[1]}
	}
	i++

	if i >= len(v) {
		return q, nil
	}

	// where clause
	if v[i].kind == keyword && v[i].value == "where" {
		i++

		condition := Condition{}
		for ; i < len(v); i++ {
			if v[i].kind == symbol {
				if condition.sign == "" {
					condition.target = v[i].value
				} else {
					condition.value = v[i].value
					q.conditions = append(q.conditions, condition)
				}
				continue
			}

			if v[i].kind == equal ||
				v[i].kind == notequal ||
				v[i].kind == greater ||
				v[i].kind == greaterorequal ||
				v[i].kind == less ||
				v[i].kind == lessorequal {
				condition.sign = v[i].value
				continue
			}
		}
	}
	i++

	return q, nil
}

func parseInsert(tokens *[]TokenLiteral) (InsertQuery, error) {
	v := *tokens
	q := InsertQuery{}
	i := 0

	// insert
	if i >= len(v) || v[i].kind != keyword || v[i].value != "insert" {
		return InsertQuery{}, errors.New("missing insert keyword")
	}
	i++

	// into
	if i >= len(v) || v[i].kind != keyword || v[i].value != "into" {
		return InsertQuery{}, errors.New("missing into keyword")
	}
	i++

	// destination table
	if i >= len(v) || v[i].kind != symbol {
		return InsertQuery{}, errors.New("missing destination table")
	} else {
		s := strings.Split(v[i].value, ".")
		if len(s) == 1 {
			q.source = SchemaTable[string, string]{defaultScheme, s[0]}
		} else {
			q.source = SchemaTable[string, string]{s[0], s[1]}
		}
	}
	i++

	// TODO: implement better algorithm for data sets and columns specification, eg. parenthesis stack
	// values or csv columns
	if i >= len(v) || (v[i].kind == keyword && v[i].value != "values") {
		return InsertQuery{}, errors.New("missing values keyword")
	} else if v[i].kind == openingroundbracket {
		// csv columns or skip
		for ; i < len(v); i++ {
			if v[i].kind == comma || v[i].kind == openingroundbracket {
				continue
			}

			if v[i].kind != symbol {
				if v[i].kind == closingroundbracket {
					i++
				}
				break
			}

			q.dataColumns = append(q.dataColumns, v[i].value)
		}

		if len(q.dataColumns) == 0 {
			return InsertQuery{}, errors.New("invalid columns specification")
		}
		// i is incremented by loop

		// values after csv
		if i >= len(v) || v[i].kind != keyword && v[i].value != "values" {
			return InsertQuery{}, errors.New("missing values keyword")
		}
		i++
	}

	// column values
	dataSet := 0
	for ; i < len(v); i++ {
		if v[i].kind == comma || v[i].kind == openingroundbracket {
			continue
		}

		// TODO: after closing bracket we should handle next data set
		if v[i].kind != symbol || v[i].kind == closingroundbracket {
			break
		}

		if len(q.values)-1 < dataSet {
			q.values = append(q.values, []any{})
		}

		q.values[dataSet] = append(q.values[dataSet], v[i].value)
	}

	return q, nil
}

func parseCreate(tokens *[]TokenLiteral) (CreateTableQuery, error) {
	v := *tokens
	q := CreateTableQuery{}
	i := 0

	// create
	if i >= len(v) || v[i].kind != keyword || v[i].value != "create" {
		return CreateTableQuery{}, errors.New("missing create keyword")
	}
	i++

	// table
	if i >= len(v) || v[i].kind != keyword || v[i].value != "table" {
		return CreateTableQuery{}, errors.New("missing table keyword")
	}
	i++

	// table name
	if i >= len(v) || v[i].kind != symbol {
		return CreateTableQuery{}, errors.New("missing table name")
	} else {
		s := strings.Split(v[i].value, ".")
		if len(s) == 1 {
			q.source = SchemaTable[string, string]{defaultScheme, s[0]}
		} else {
			q.source = SchemaTable[string, string]{s[0], s[1]}
		}
	}
	i++

	// csv columns or skip
	q.columns = make(map[string][]string)
	columnsAttributes := []string{}
	for ; i < len(v); i++ {
		if v[i].kind == openingroundbracket {
			continue
		}

		if v[i].kind == comma {
			attributes := make([]string, len(columnsAttributes[1:]))
			copy(attributes, columnsAttributes[1:])
			q.columns[columnsAttributes[0]] = attributes
			columnsAttributes = []string{}
			continue
		}

		if v[i].kind != symbol {
			if v[i].kind == closingroundbracket {
				attributes := make([]string, len(columnsAttributes[1:]))
				copy(attributes, columnsAttributes[1:])
				q.columns[columnsAttributes[0]] = attributes
				i++
			}
			break
		}

		columnsAttributes = append(columnsAttributes, v[i].value)
	}

	log.Println(q)

	if len(q.columns) == 0 {
		return CreateTableQuery{}, errors.New("invalid columns specification")
	}
	// i is incremented by loop

	return q, nil
}
