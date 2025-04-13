package main

import (
	"slices"
	"strings"
)

type TokenKind int

const (
	keyword TokenKind = iota
	symbol
	comma
	openingroundbracket
	closingroundbracket
	equal
)

var keywords []string = []string{
	"select",
	"from",
	"where",

	"insert",
	"into",
	"values",

	"create",
	"table",
}

type TokenLiteral struct {
	kind  TokenKind
	value string
}

func Analyze(raw string) []TokenLiteral {
	tokens := make([]TokenLiteral, 0)

	l := 0
	len := len(raw)
	for r := range len {
		switch raw[r] {
		case byte(' '):
			{
				if r == l {
					l++
					continue
				}

				frag := strings.ToLower(string(raw[l:r]))
				if slices.Contains(keywords, frag) {
					tokens = append(tokens, TokenLiteral{kind: keyword, value: frag})
				} else {
					tokens = append(tokens, TokenLiteral{kind: symbol, value: frag})
				}

				l = r + 1
			}
		case byte('='):
			{
				tokens = append(tokens, TokenLiteral{kind: equal, value: string('=')})
				l = r + 1
			}
		case byte('('):
			{
				if l != r {
					frag := strings.ToLower(string(raw[l:r]))
					tokens = append(tokens, TokenLiteral{kind: symbol, value: frag})
				}

				tokens = append(tokens, TokenLiteral{kind: openingroundbracket, value: string('(')})
				l = r + 1
			}
		case byte(')'):
			{
				if l != r {
					frag := strings.ToLower(string(raw[l:r]))
					tokens = append(tokens, TokenLiteral{kind: symbol, value: frag})
				}

				tokens = append(tokens, TokenLiteral{kind: closingroundbracket, value: string(')')})
				l = r + 1
			}
		case byte(','):
			{
				if raw[r-1] != byte(' ') && r != l {
					frag := strings.ToLower(string(raw[l:r]))
					tokens = append(tokens, TokenLiteral{kind: symbol, value: frag})
				}

				tokens = append(tokens, TokenLiteral{kind: comma, value: string(',')})
				l = r + 1
			}
		}

		// l-1 because we assign l = r + 1
		if r == len-1 && r != l-1 {
			tokens = append(tokens, TokenLiteral{kind: symbol, value: string(raw[l:len])})
		}
	}

	return tokens
}
