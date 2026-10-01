package tokenizer

type state struct {
	b, p uint
	line uint
	chw  uint
	ch   byte
}

type State = state

type source struct {
	text string
	state
}

func (s *source) init(text string) {
	s.text = text

	s.b, s.p = 0, 0
	s.line = 1
	s.nextch()
}

func (s *source) State() State {
	return s.state
}

func (s *source) Reset(state State) {
	s.state = state
}

func (s *source) error(msg string) {
}

func (s *source) len() uint { return uint(len(s.text)) }

func (s *source) eof() bool              { return s.chw == 0 }
func (s *source) hasAtLeast(n uint) bool { return s.p-s.chw+n < s.len() }
func (s *source) position() uint         { return s.p - s.chw }

func (s *source) start()                           { s.b = s.p - s.chw }
func (s *source) segment() string                  { return s.text[s.b : s.p-s.chw] }
func (s *source) segmentAt(begin, end uint) string { return s.text[begin:end] }
func (s *source) segmentFrom(begin uint) string    { return s.text[begin : s.p-s.chw] }
func (s *source) segmentTo(end uint) string        { return s.text[s.b:end] }

func (s *source) at(pos uint) byte {
	return s.text[pos]
}

func (s *source) right() string {
	if s.p == 0 {
		return s.text
	}
	return s.text[s.p-s.chw:]
}

func (s *source) rewind(offset uint) {
	s.p = s.b + offset
	s.nextch()
}

func (s *source) rewindTo(position uint) {
	s.p = position
	s.nextch()
}

func (s *source) nextch() {
	if s.p >= s.len() {
		s.ch = ';'
		s.chw = 0
		return
	}

	s.ch = s.text[s.p]
	s.chw = 1
	s.p++
}

func (s *source) advance(n uint) {
	s.p = min(s.len(), s.p-s.chw+n)
	s.nextch()
}

// byteAt returns the byte at offset from the next byte, starting at 1.
func (s *source) byteAt(offset uint) byte {
	return s.text[s.p-s.chw+offset]
}

func (s *source) skipNewline() bool {
	if s.ch == '\r' {
		s.nextch()
		if s.ch == '\n' {
			s.nextch()
		}
		s.line++
		return true
	}

	s.nextch()
	s.line++
	return false
}

func (s *source) skipWhitespace() {
	for {
		switch s.ch {
		case ' ', '\t':
			s.nextch()
		case '\n', '\f', '\r':
			s.skipNewline()
		default:
			return
		}
	}
}

func (s *source) skipDigits() {
	for s.isDigit() {
		s.nextch()
	}
}

func (s *source) skipComment() bool {
	for s.chw > 0 {
		for s.ch == '*' {
			s.nextch()
			if s.ch == '/' {
				s.nextch()
				return true
			}
		}
		s.nextch()
	}
	return false
}

func (s *source) skipCDOAndCDC() {
	for {
		switch s.ch {
		case ' ', '\t':
			s.nextch()
		case '\n', '\f', '\r':
			s.skipNewline()
		case '/':
			if s.hasAtLeast(1) && s.byteAt(1) == '*' {
				s.skipComment()
			} else {
				return
			}
		case '<':
			if s.hasAtLeast(3) && s.byteAt(1) == '!' && s.byteAt(2) == '-' && s.byteAt(3) == '-' {
				s.advance(4)
			} else {
				return
			}
		case '-':
			if s.hasAtLeast(2) && s.byteAt(1) == '-' && s.byteAt(2) == '>' {
				s.advance(3)
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *source) isIdentStart() bool {
	switch s.ch {
	case '-':
		if !s.hasAtLeast(1) {
			return false
		}
		ch := s.byteAt(1)
		if isIdentStartLUT(ch) || ch == '-' {
			return true
		}
		return ch == '\\' && (!s.hasAtLeast(2) || !s.hasNewlineAt(2))

	case '\\':
		if !s.hasAtLeast(1) {
			return false
		}
		return !s.hasNewlineAt(1)

	default:
		return isIdentStartLUT(s.ch)
	}
}

func (s *source) isDigit() bool {
	return '0' <= s.ch && s.ch <= '9'
}

func (s *source) isNewline() bool {
	return s.ch == '\n' || s.ch == '\r' || s.ch == '\f'
}

func (s *source) hasDigitsAt(offset uint) bool {
	ch := s.byteAt(offset)
	return '0' <= ch && ch <= '9'
}

func (s *source) hasNewlineAt(offset uint) bool {
	ch := s.byteAt(offset)
	return ch == '\n' || ch == '\r' || ch == '\f'
}
