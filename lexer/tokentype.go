package lexer

import "fmt"

type TokenType int

const (
	// A TokenError represents any lexical analysis error.
	TokenError TokenType = iota
	// A TokenEOF represents the end of a file.
	TokenEOF
	// A TokenNewLine represents a new line.
	TokenNewLine

	// A TokenPunctuation represents punctuation (parenthesis, dot, comma, ...).
	TokenPunctuation

	// A TokenNumber represents a number literal.
	TokenNumber
	// A TokenOperator represents an operator (addition, subtraction, ...).
	TokenOperator

	// A TokenName represents any name.
	// A name is any sequence of unicode letters and numbers
	// that does not start with a number.
	TokenName
)

var tokenTypeNames = map[TokenType]string{
	TokenError:   "ERROR",
	TokenEOF:     "EOF",
	TokenNewLine: "NEW_LINE",

	TokenPunctuation: "PUNCTUATION",

	TokenNumber:   "NUMBER",
	TokenOperator: "OPERATOR",

	TokenName: "NAME",
}

// String returns a string representation of a TokenType
func (tt TokenType) String() string {
	str := tokenTypeNames[tt]
	if str != "" {
		return str
	}
	return fmt.Sprintf("Token%d", tt)
}
