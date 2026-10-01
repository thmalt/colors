package tokenizer

func lower(ch byte) byte {
	return ('a' - 'A') | ch
}

func isLetter(ch byte) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z'
}

// 0-9
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isNewline(ch byte) bool {
	return ch == '\n' || ch == '\r' || ch == '\f'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\f' || ch == '\r'
}

// 0-9, A-F, a-f
func isHex(ch byte) bool {
	return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f'
}

// A-Z, a-z, _, Non-ASCII (>= 0x80)
func isIdentStart(ch byte) bool {
	return isLetter(ch) || ch == '_' || ch >= 0x80
}

// A-Za-z_-0-9 >= RuneSelf
func isIdentContinue(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '-'
}
