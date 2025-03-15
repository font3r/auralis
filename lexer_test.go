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
				{kind: symbol, value: ","},
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
				{kind: symbol, value: ","},
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
				{kind: symbol, value: ","},
				{kind: symbol, value: "id2"},
				{kind: keyword, value: "from"},
				{kind: symbol, value: "users"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.raw, func(t *testing.T) {
			tokens := Analyze(tC.raw)
			if !slices.Equal(tokens, tC.expected) {
				t.Errorf("expect result %+v\n got %+v", tC.expected, tokens)
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
			raw: "INSERT INTO users VALUES ('1', '2')",
			expected: []TokenLiteral{
				{kind: keyword, value: "insert"},
				{kind: keyword, value: "into"},
				{kind: symbol, value: "users"},
				{kind: keyword, value: "values"},
				{kind: openingroundbracket, value: "("},
				{kind: symbol, value: "'1'"},
				{kind: symbol, value: ","},
				{kind: symbol, value: "'2'"},
				{kind: closingroundbracket, value: ")"},
			},
		},
		// {
		// 	raw: "INSERT INTO users (id1, id2) VALUES ('1','2')",
		// 	expected: []TokenLiteral{
		// 		{kind: keyword, value: "insert"},
		// 		{kind: keyword, value: "into"},
		// 		{kind: symbol, value: "users"},
		// 		{kind: keyword, value: "values"},
		// 		{kind: openingroundbracket, value: "("},
		// 		{kind: symbol, value: "'1'"},
		// 		{kind: symbol, value: ","},
		// 		{kind: symbol, value: "'2'"},
		// 		{kind: closingroundbracket, value: ")"},
		// 	},
		// },
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

/*

expect result [{kind:0 value:insert} {kind:0 value:into} {kind:1 value:users} {kind:0 value:values} {kind:2 value:(} {kind:1 value:'1'} {kind:1 value:,} {kind:1 value:'2'} {kind:3 value:)}]
          got [{kind:0 value:insert} {kind:0 value:into} {kind:1 value:users} {kind:0 value:values} {kind:2 value:(} {kind:1 value:'1'} {kind:1 value:,} {kind:2 value:)} {kind:1 value:}]
*/
