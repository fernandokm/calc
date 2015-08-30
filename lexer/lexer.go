package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Lexer is a type that can lex an input string into a chan *Token.
type Lexer struct {
	input  string
	tokens chan *Token

	start    int
	pos      int
	eofDepth int
}

// NewLexer creates a new lexer from an input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		tokens: make(chan *Token, 2),
		start:  0,
		pos:    0,
	}
}

// TODO(fernandokm): make next() and backoff() methods better

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.eofDepth++
		return eof
	}

	result, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += width
	return result
}

func (l *Lexer) backoff() {
	if l.eofDepth > 0 {
		l.eofDepth--
	} else {
		l.pos--
	}
}

func (l *Lexer) emit(t TokenType) {
	l.tokens <- &Token{t, l.input[l.start:l.pos]}
	l.setStart()
}

func (l *Lexer) peek() rune {
	pos := l.pos
	r := l.next()
	l.pos = pos
	return r
}

func (l *Lexer) resetToken() {
	l.pos = l.start
}

func (l *Lexer) setStart() {
	l.start = l.pos
}

func (l *Lexer) skipWhitespace() {
	l.acceptWhile(isWhitespace)
	l.setStart()
}

func (l *Lexer) accept(runes string) bool {
	if strings.IndexRune(runes, l.next()) >= 0 {
		return true
	}
	l.backoff()
	return false
}

func (l *Lexer) acceptAll(runes string) {
	for l.accept(runes) {
	}
}

func (l *Lexer) acceptIf(condition func(rune) bool) bool {
	r := l.peek()
	if condition(r) {
		l.next()
		return true
	}
	return false
}

func (l *Lexer) acceptWhile(condition func(rune) bool) {
	r := l.next()
	for condition(r) {
		r = l.next()
	}
	if r != eof {
		l.backoff()
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) lexerFunc {
	l.tokens <- &Token{
		TokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// Tokens returns the token chan for this lexer.
func (l *Lexer) Tokens() <-chan *Token {
	return l.tokens
}

// Input returns the string input of this lexer.
func (l *Lexer) Input() string {
	return l.input
}

// Run runs this lexer and closes the chan returned by Lexer.Tokens() when done.
func (l *Lexer) Run() {
	for state := lexStatement; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}
