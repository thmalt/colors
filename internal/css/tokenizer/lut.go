package tokenizer

const (
	chLetter uint8 = 1 << iota
	chDigit
	chHex
	chWhitespace
	chNewline
	chIdentStart
	chIdentContinue
)

var chLUT [1 << 8]uint8

func init() {
	for ch := '0'; ch <= '9'; ch++ {
		chLUT[ch] |= chDigit | chHex | chIdentContinue
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		flags := chLetter | chIdentStart | chIdentContinue
		if ch <= 'F' {
			flags |= chHex
		}

		chLUT[ch] = flags
		chLUT[ch|0x20] = flags
	}

	for ch := 1 << 7; ch < 1<<8; ch++ {
		chLUT[ch] = chIdentStart | chIdentContinue
	}

	chLUT['-'] = chIdentContinue
	chLUT['_'] = chIdentStart | chIdentContinue

	chLUT[' '] = chWhitespace
	chLUT['\t'] = chWhitespace

	chLUT['\n'] = chWhitespace | chNewline
	chLUT['\f'] = chWhitespace | chNewline
	chLUT['\r'] = chWhitespace | chNewline
}

func chFlagLUT(ch byte) uint8 { return chLUT[ch] }

func isLetterLUT(ch byte) bool        { return chLUT[ch]&chLetter != 0 }
func isDigitLUT(ch byte) bool         { return chLUT[ch]&chDigit != 0 }
func isHexLUT(ch byte) bool           { return chLUT[ch]&chHex != 0 }
func isWhitespaceLUT(ch byte) bool    { return chLUT[ch]&chWhitespace != 0 }
func isNewlineLUT(ch byte) bool       { return chLUT[ch]&chNewline != 0 }
func isIdentStartLUT(ch byte) bool    { return chLUT[ch]&chIdentStart != 0 }
func isIdentContinueLUT(ch byte) bool { return chLUT[ch]&chIdentContinue != 0 }

var identContinueLUT [1 << 8]bool

func init() {
	for ch := '0'; ch <= '9'; ch++ {
		identContinueLUT[ch] = true
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		identContinueLUT[ch] = true
		identContinueLUT[ch|0x20] = true
	}

	for ch := 1 << 7; ch < 1<<8; ch++ {
		identContinueLUT[ch] = true
	}

	identContinueLUT['-'] = true
	identContinueLUT['_'] = true
}

func isIdentContinueBoolLUT(ch byte) bool { return identContinueLUT[ch] }
