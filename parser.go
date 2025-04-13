package main

import (
	"errors"
	"strings"
)

type Command struct{}

type SchemeTable[T, U string] struct {
	scheme T
	name   U
}

type Condition struct {
	target string
	sign   string
	value  any
}

type SelectQuery struct {
	source      SchemeTable[string, string]
	dataColumns []string
	conditions  []Condition
}

type InsertQuery struct {
	source      SchemeTable[string, string]
	dataColumns []string // column names
	values      [][]any  // column values
}

type CreateTableQuery struct {
	source  SchemeTable[string, string]
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
		q.source = SchemeTable[string, string]{"dbo", s[0]}
	} else {
		q.source = SchemeTable[string, string]{s[0], s[1]}
	}

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
			q.source = SchemeTable[string, string]{"dbo", s[0]}
		} else {
			q.source = SchemeTable[string, string]{s[0], s[1]}
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
	// TODO: temp passthrough to allow create hardcoded tables
	return CreateTableQuery{}, nil
}
