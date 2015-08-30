package lexer

import (
	"strings"
	"unicode"
)

var (
	operators   = "+-*/="
	whitespace  = " \t"
	newLine     = "\r\n"
	punctuation = ".,()[]{}"
)

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func isWhitespace(r rune) bool {
	return is(whitespace, r)
}

func isNewLine(r rune) bool {
	return is(newLine, r)
}

func isOperator(r rune) bool {
	return is(operators, r)
}

func isPunctuation(r rune) bool {
	return is(punctuation, r)
}

func is(options string, r rune) bool {
	return strings.IndexRune(options, r) >= 0
}

func isLetterOrUnderscore(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}
