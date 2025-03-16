package main

import (
	"errors"
	"strings"
)

type Command struct{}

type SelectQuery struct {
	source  SchemeTable[string, string]
	columns []string
}

type SchemeTable[T, U string] struct {
	scheme T
	name   U
}

func ParseTokens(tokens []TokenLiteral) (any, error) {
	valid := hasAnyKeyword(&tokens)
	if !valid {
		return Command{}, errors.New("missing keyword")
	}

	// for now assume that query starts with keyword
	switch tokens[0].value {
	case "select":
		return parseSelect(&tokens)
	}

	return Command{}, errors.New("missing keyword")
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

		q.columns = append(q.columns, v[i].value)
	}

	if len(q.columns) == 0 {
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
