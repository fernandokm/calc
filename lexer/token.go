package lexer

import "fmt"

const (
	eof = rune(0)
)

// Token represents a lexical token composed of a type and a value.
type Token struct {
	Type  TokenType
	Value string
}

// String returns a string representation of this Token.
func (t *Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "<EOF>"
	default:
		return fmt.Sprintf("<%s, %q>", t.Type, t.Value)
	}
}
