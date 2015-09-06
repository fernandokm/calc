package lexer

import (
	"fmt"
	"unicode/utf8"
)

// Lexer is a type that can lex an input string into a chan *Token.
type Lexer struct {
	input  string
	tokens chan *Token
	last   *Token

	start    int
	pos      int
	eofDepth int
}

//TODO(ferandokm): make Lexer thread-safe

// NewLexer creates a new lexer from an input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		tokens: make(chan *Token, 2),
		start:  0,
		pos:    0,
	}
}

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
		for !utf8.RuneStart(l.input[l.pos]) {
			l.pos--
		}
	}
}

func (l *Lexer) emit(tt TokenType) {
	t := &Token{tt, l.getCurrentValue()}
	l.tokens <- t
	l.last = t
	l.setStart()
}

func (l *Lexer) lastToken() *Token {
	return l.last
}

func (l *Lexer) lastTokenType() TokenType {
	if l.last == nil {
		return -1
	}
	fmt.Println(l.last)
	return l.last.Type
}

func (l *Lexer) lastTokenValue() string {
	if l.last == nil {
		return ""
	}
	return l.last.Value
}

func (l *Lexer) getCurrentValue() string {
	return l.input[l.start:l.pos]
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
	l.acceptRunes(func(r rune) bool { return is(r, whitespace) })
	l.setStart()
}

func (l *Lexer) acceptRune(rules ...interface{}) bool {
	if is(l.peek(), rules...) {
		l.next()
		return true
	}
	return false
}

func (l *Lexer) acceptRunes(rules ...interface{}) bool {
	if !l.acceptRune(rules...) {
		return false
	}
	for l.acceptRune(rules...) {
	}
	return true
}

func (l *Lexer) acceptString(rules ...interface{}) bool {
	length := stringRuleLength(l.input, l.pos, rules...)
	if length == 0 {
		return false
	}
	l.pos += length
	return true
}

func (l *Lexer) acceptStrings(rules ...interface{}) bool {
	if !l.acceptString(rules...) {
		return false
	}
	for l.acceptString(rules...) {
	}
	return true
}

func (l *Lexer) getPos() int {
	return l.pos
}

func (l *Lexer) restorePos(pos int) {
	l.pos = pos
}

func (l *Lexer) errorf(format string, args ...interface{}) lexerFunc {
	l.tokens <- &Token{
		TokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *Lexer) errorExpected(what string) lexerFunc {
	return l.errorf("Expected %s, got %q", what, l.peek())
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
	for state := lexExpression; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}
