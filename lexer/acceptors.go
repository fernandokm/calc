package lexer

func acceptWhitespace(l *Lexer) bool {
	return l.acceptRunes(whitespace)
}

func acceptDecimal(l *Lexer) bool {
	start := l.getPos()
	l.acceptRune(signSymbol)
	acceptWhitespace(l)
	acceptDecimalPart := func() bool {
		pos := l.getPos()
		if l.acceptRune(decimalPoint) && l.acceptRunes(digit) {
			return true
		}
		l.restorePos(pos)
		return false
	}

	if l.acceptRunes(digit) {
		acceptDecimalPart()
		return true
	} else if acceptDecimalPart() {
		return true
	}
	l.restorePos(start)
	return false

}
