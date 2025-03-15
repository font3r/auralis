package main

import (
	"fmt"
	"slices"
	"strings"
)

type TokenKind int

const (
	keyword TokenKind = iota
	symbol
	openingroundbracket
	closingroundbracket
)

type TokenLiteral struct {
	kind  TokenKind
	value string
}

var keywords []string = []string{
	"select",
	"from",
	"insert",
	"into",
	"values",
}

func Analyze(raw string) []TokenLiteral {
	tokens := make([]TokenLiteral, 0)

	l := 0
	len := len(raw)
	for r := range len {
		if raw[r] == byte(' ') {
			// after jumping l can be equal to r
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

			// always r + 1 if space
			l = r + 1
			continue
		}

		if raw[r] == byte('(') {
			tokens = append(tokens, TokenLiteral{kind: openingroundbracket, value: string('(')})
			l = r + 1
			continue
		}

		if raw[r] == byte(')') {
			fmt.Println(string(raw[r]))
			tokens = append(tokens, TokenLiteral{kind: closingroundbracket, value: string(')')})
			l = r + 1
			continue
		}

		if raw[r] == byte(',') {
			if raw[r-1] == byte(' ') {
				tokens = append(tokens, TokenLiteral{kind: symbol, value: string(',')})
			} else {
				frag := strings.ToLower(string(raw[l:r]))
				tokens = append(tokens, TokenLiteral{kind: symbol, value: frag})
				tokens = append(tokens, TokenLiteral{kind: symbol, value: string(',')})
			}
			l = r + 1
			continue
		}

		if r == len-1 {
			tokens = append(tokens, TokenLiteral{kind: symbol, value: string(raw[l:len])})
		}
	}

	return tokens
}
