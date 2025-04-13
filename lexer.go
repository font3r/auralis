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
	notequal
	greater
	greaterorequal
	less
	lessorequal
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
	for r := 0; r < len(raw); r++ {
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
		case byte('>'):
			{
				if raw[r+1] == byte('=') {
					tokens = append(tokens, TokenLiteral{kind: greaterorequal, value: string(">=")})
					// r + 1 incremented by next loop iteration
					l = r + 2
					r += 1
				} else {
					tokens = append(tokens, TokenLiteral{kind: greater, value: string('>')})
					l = r + 1
				}
			}
		case byte('<'):
			{
				if raw[r+1] == byte('=') {
					tokens = append(tokens, TokenLiteral{kind: lessorequal, value: string("<=")})
					// r + 1 incremented by next loop iteration
					l = r + 2
					r += 1
				} else {
					tokens = append(tokens, TokenLiteral{kind: less, value: string('<')})
					l = r + 1
				}
			}
		case byte('!'):
			{
				if raw[r+1] == byte('=') {
					tokens = append(tokens, TokenLiteral{kind: notequal, value: string("!=")})
				}

				// r + 1 incremented by next loop iteration
				l = r + 2
				r += 1
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
		if r == len(raw)-1 && r != l-1 {
			tokens = append(tokens, TokenLiteral{kind: symbol, value: string(raw[l:])})
		}
	}

	return tokens
}
