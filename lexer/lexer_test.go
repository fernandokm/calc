package lexer

import (
	"bytes"
	"testing"
)

var (
	eofToken = &Token{TokenEOF, ""}
)

func getTokens(input string) <-chan *Token {
	lex := NewLexer(input)
	go lex.Run()
	return lex.Tokens()
}

func checkTokens(input string, expected ...*Token) bool {
	actualTokens := getTokens(input)
	for _, token := range expected {
		actual, ok := <-actualTokens
		if !ok || *actual != *token {
			return false
		}
	}
	_, ok := <-actualTokens
	if ok {
		return false // Do not accept more tokens than expected
	}
	return true
}

func printTokens(input string, t *testing.T) {
	buffer := bytes.NewBufferString("Got: ")
	for token := range getTokens(input) {
		buffer.WriteString(token.String())
		buffer.WriteString(", ")
	}
	msg := buffer.String()
	t.Log(msg[:len(msg)-2])
}

func getChecker(typ TokenType, t *testing.T) func(input, value string) {
	return func(input, value string) {
		expected := &Token{typ, value}
		if !checkTokens(input, expected, eofToken) {
			t.Logf("Expected \"%s\" to yield %s, %s", input, expected, eofToken)
			printTokens(input, t)
			t.Fail()
		}
	}
}

func TestNumbers(t *testing.T) {
	check := getChecker(TokenNumber, t)
	check("-53", "-53")
	check("+  5", "+  5")
	check("-.9 ", "-.9")
	check("  -2e+1 ", "-2e+1")
}

func TestPunctuation(t *testing.T) {
	check := getChecker(TokenPunctuation, t)
	check(".", ".")
	check("  ,  ", ",")
	check("(  ", "(")
}

func TestOperators(t *testing.T) {
	check := getChecker(TokenOperator, t)
	check(" +", "+")
	check("  - ", "-")
	check(" ^  ", "^")
	check("*", "*")
	check(" /   ", "/")
}

func TestNames(t *testing.T) {
	check := getChecker(TokenName, t)
	check("ab ", "ab")
	check("_ad5", "_ad5")
	check("  _av", "_av")
}

func TestExpressions(t *testing.T) {
	check := func(input string, tokens ...*Token) {
		tokens = append(tokens, eofToken)
		if !checkTokens(input, tokens...) {
			t.Fail()
			t.Logf("Expression: %s", input)
			t.Logf("Expected: %s", tokens)
			printTokens(input, t)
		}
	}
	check("1+2", &Token{TokenNumber, "1"}, &Token{TokenOperator, "+"}, &Token{TokenNumber, "2"})
	check("-1 +_abc / 7 ", &Token{TokenNumber, "-1"}, &Token{TokenOperator, "+"}, &Token{TokenName, "_abc"}, &Token{TokenOperator, "/"}, &Token{TokenNumber, "7"})
}
