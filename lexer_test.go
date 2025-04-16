package main

import (
	"slices"
	"testing"
)

func TestSelectStatementLexer(t *testing.T) {
	testCases := []struct {
		raw      string
		expected []TokenLiteral
	}{
		{
			raw: "SELECT * FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "select * from users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "select * FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "SELECT id1, id2 FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "SELECT id1 , id2 FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "SELECT id1 ,id2 FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "SELECT id1, id2, id3 FROM users",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id3"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id = 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: equal, value: "="},
				{kind: symbol, value: "1"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id != 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: notequal, value: "!="},
				{kind: symbol, value: "1"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id < 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: less, value: "<"},
				{kind: symbol, value: "1"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id > 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: greater, value: ">"},
				{kind: symbol, value: "1"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id <= 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: lessorequal, value: "<="},
				{kind: symbol, value: "1"},
			},
		},
		{
			raw: "SELECT * FROM users WHERE id >= 1",
			expected: []TokenLiteral{
				{kind: keyword, value: "select"},
				{kind: symbol, value: "*"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "where"},
				{kind: symbol, value: "id"},
				{kind: greaterorequal, value: ">="},
				{kind: symbol, value: "1"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.raw, func(t *testing.T) {
			tokens := Analyze(tC.raw)
			if !slices.Equal(tokens, tC.expected) {
				t.Errorf("\nexp %+v\ngot %+v", tC.expected, tokens)
			}
		})
	}
}

func TestInsertStatementLexer(t *testing.T) {
	testCases := []struct {
		raw      string
		expected []TokenLiteral
	}{
		{
			raw: "insert into users values ('1','2')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "insert into users (id1,  id2) values  ('1',   '2'  )",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "INSERT INTO users VALUES ('1', '2')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "INSERT INTO users (id1, id2) VALUES ('1','2')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "INSERT INTO users(id1, id2) VALUES ('1','2')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "INSERT INTO users (id1, id2) VALUES ('1','2'), ('3', '4')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id1"},
				{kind: comma, value: ","},
				{kind: symbol, value: "id2"},
				{kind: closingroundbracket, value: ")"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
				{kind: comma, value: ","},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'3'"},
				{kind: comma, value: ","},
				{kind: symbol, value: "'4'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.raw, func(t *testing.T) {
			tokens := Analyze(tC.raw)
			if !slices.Equal(tokens, tC.expected) {
				t.Errorf("\nexp %+v\ngot %+v", tC.expected, tokens)
			}
		})
	}
}

func TestCreateStatementLexer(t *testing.T) {
	testCases := []struct {
		raw      string
		expected []TokenLiteral
	}{
		{
			raw: "create table users",
			expected: []TokenLiteral{
				{kind: keyword, value: "create"},
				{kind: keyword, value: "table"},
				{kind: symbol, value: "users"},
			},
		},
		{
			raw: "create table users (age smallint)",
			expected: []TokenLiteral{
				{kind: keyword, value: "create"},
				{kind: keyword, value: "table"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "age"},
				{kind: symbol, value: "smallint"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "create table users (age smallint, name varchar)",
			expected: []TokenLiteral{
				{kind: keyword, value: "create"},
				{kind: keyword, value: "table"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "age"},
				{kind: symbol, value: "smallint"},
				{kind: comma, value: ","},
				{kind: symbol, value: "name"},
				{kind: symbol, value: "varchar"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "create table users (id uniqueidentifier, age smallint, name varchar)",
			expected: []TokenLiteral{
				{kind: keyword, value: "create"},
				{kind: keyword, value: "table"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id"},
				{kind: symbol, value: "uniqueidentifier"},
				{kind: comma, value: ","},
				{kind: symbol, value: "age"},
				{kind: symbol, value: "smallint"},
				{kind: comma, value: ","},
				{kind: symbol, value: "name"},
				{kind: symbol, value: "varchar"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		{
			raw: "create table users (id uniqueidentifier, age smallint, name varchar not null)",
			expected: []TokenLiteral{
				{kind: keyword, value: "create"},
				{kind: keyword, value: "table"},
				{kind: symbol, value: "users"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "id"},
				{kind: symbol, value: "uniqueidentifier"},
				{kind: comma, value: ","},
				{kind: symbol, value: "age"},
				{kind: symbol, value: "smallint"},
				{kind: comma, value: ","},
				{kind: symbol, value: "name"},
				{kind: symbol, value: "varchar"},
				{kind: symbol, value: "not"},
				{kind: symbol, value: "null"},
				{kind: closingroundbracket, value: ")"},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.raw, func(t *testing.T) {
			tokens := Analyze(tC.raw)
			if !slices.Equal(tokens, tC.expected) {
				t.Errorf("\nexp %+v\ngot %+v", tC.expected, tokens)
			}
		})
	}
}
