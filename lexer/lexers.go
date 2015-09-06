package lexer

func lexExpression(l *Lexer) lexerFunc {
	l.skipWhitespace()
	r := l.peek()
	lastType := l.lastTokenType()
	if r == eof {
		l.next()
		l.emit(TokenEOF)
		return nil
	} else if is(r, newLine) {
		l.next()
		l.emit(TokenNewLine)
		return lexExpression
	} else if is(r, digit) || (is(r, decimalPoint) && lastType != TokenName) || (is(r, signSymbol) && lastType != TokenNumber && lastType != TokenName) {
		return lexRealNumber
	} else if is(r, letter, underscore) {
		return lexName
	} else if is(r, punctuation) {
		return lexPunctuation
	} else {
		return lexOperator
	}
}

// grammar: real-number
func lexRealNumber(l *Lexer) lexerFunc {
	if !acceptDecimal(l) {
		if is(l.peek(), signSymbol) {
			return lexOperator
		} else if is(l.peek(), decimalPoint) {
			return lexPunctuation
		}
		return l.errorExpected("number literal")
	}
	pos := l.getPos()
	acceptWhitespace(l)
	if !(l.acceptRune(exponentSymbol) && acceptDecimal(l)) {
		l.restorePos(pos)
	}
	l.emit(TokenNumber)
	return lexExpression
}

// grammar: operator
func lexOperator(l *Lexer) lexerFunc {
	if l.acceptString(operator) {
		l.emit(TokenOperator)
		return lexExpression
	}
	return l.errorExpected("operator")
}

// grammar: punctuation
func lexPunctuation(l *Lexer) lexerFunc {
	if l.acceptRune(punctuation) {
		l.emit(TokenPunctuation)
		return lexExpression
	}
	return l.errorExpected("punctuation")
}

// grammar: name
func lexName(l *Lexer) lexerFunc {
	if l.acceptRune("_", letter) {
		l.acceptRunes("_", letter, digit)
		l.emit(TokenName)
		return lexExpression
	}
	return l.errorExpected("name")
}
