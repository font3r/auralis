package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestSelectParser(t *testing.T) {
	testCases := map[string]struct {
		tokens      []TokenLiteral
		expectedCmd any
		expectedErr error
	}{
		"query without any keyword": {
			tokens: []TokenLiteral{
				{kind: symbol, value: "test"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: Command{},
			expectedErr: errors.New("missing keyword"),
		},
		"query without select keyword": {
			tokens: []TokenLiteral{
				{kind: symbol, value: "test"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: Command{},
			expectedErr: errors.New("missing keyword"),
		},
		"select without specified columns": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: SelectQuery{},
			expectedErr: errors.New("missing columns"),
		},
		"query without from keyword": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: SelectQuery{},
			expectedErr: errors.New("missing from keyword"),
		},
		"select without specified source table": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
			},
			expectedCmd: SelectQuery{},
			expectedErr: errors.New("missing source table"),
		},
		"valid select all columns from table": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: SelectQuery{
				source:      SchemeTable[string, string]{"dbo", "users"},
				dataColumns: []string{"*"},
			},
		},
		"valid select specific columns from table": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: SelectQuery{
				source:      SchemeTable[string, string]{"dbo", "users"},
				dataColumns: []string{"id1", "id2"},
			},
		},
		"valid select specific columns with where clause": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id1"},
				{kind: equal, value: "="},
				{kind: symbol, value: "1"},
			},
			expectedCmd: SelectQuery{
				source:      SchemeTable[string, string]{"dbo", "users"},
				dataColumns: []string{"id1", "id2"},
				conditions: []Condition{
					{target: "id1", sign: "=", value: "1"},
				},
			},
		},
	}
	for test, tC := range testCases {
		t.Run(test, func(t *testing.T) {
			cmd, err := ParseTokens(tC.tokens)
			// TODO: compare constant errors
			if err != nil && err.Error() != tC.expectedErr.Error() {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedErr, err)
			} else if !reflect.DeepEqual(cmd, tC.expectedCmd) {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedCmd, cmd)
			}
		})
	}
}

func TestInsertParser(t *testing.T) {
	testCases := map[string]struct {
		tokens      []TokenLiteral
		expectedCmd any
		expectedErr error
	}{
		"query without any keyword": {
			tokens: []TokenLiteral{
				{kind: symbol, value: "test"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: Command{},
			expectedErr: errors.New("missing keyword"),
		},
		"query without insert keyword": {
			tokens: []TokenLiteral{
				{kind: symbol, value: "test"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: Command{},
			expectedErr: errors.New("missing keyword"),
		},
		"query without into keyword": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: InsertQuery{},
			expectedErr: errors.New("missing into keyword"),
		},
		"insert without specified destination table": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
			},
			expectedCmd: InsertQuery{},
			expectedErr: errors.New("missing destination table"),
		},
		"insert without values keyword": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
			},
			expectedCmd: InsertQuery{},
			expectedErr: errors.New("missing values keyword"),
		},
		"insert with missing values after column specification": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id"},
				{kind: closingroundbracket, value: ")"},
			},
			expectedCmd: InsertQuery{},
			expectedErr: errors.New("missing values keyword"),
		},
		"valid insert with columns specification": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "1"},
				{kind: closingroundbracket, value: ")"},
			},
			expectedCmd: InsertQuery{
				source:      SchemeTable[string, string]{"dbo", "users"},
				dataColumns: []string{"id"},
				values: [][]any{
					{"1"},
				},
			},
			expectedErr: nil,
		},
		"valid insert values with single quotation marks and columns specification": {
			tokens: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "name"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'example'"},
				{kind: closingroundbracket, value: ")"},
			},
			expectedCmd: InsertQuery{
				source:      SchemeTable[string, string]{"dbo", "users"},
				dataColumns: []string{"name"},
				values: [][]any{
					{"'example'"},
				},
			},
			expectedErr: nil,
		},
	}
	for test, tC := range testCases {
		t.Run(test, func(t *testing.T) {
			cmd, err := ParseTokens(tC.tokens)
			// TODO: compare constant errors
			if err != nil && err.Error() != tC.expectedErr.Error() {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedErr, err)
			} else if !reflect.DeepEqual(cmd, tC.expectedCmd) {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedCmd, cmd)
			}
		})
	}
}
