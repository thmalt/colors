package tokenizer

import (
	"strconv"
	"unicode/utf8"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
	"github.com/thmalt/colors/internal/hexlut"
)

type Tokenizer struct {
	source

	tok token.Token
	buf []byte
}

func New(text string) Tokenizer {
	var t Tokenizer
	t.init(text)
	return t
}

func (t *Tokenizer) Init(text string) {
	t.init(text)
}

func (t *Tokenizer) Token() token.Token {
	return t.tok
}

func (t *Tokenizer) init(text string) {
	t.source.init(text)
	t.tok = token.Token{}
}

func (t *Tokenizer) kind(kind token.Kind) {
	t.tok.Kind = kind
}

func (t *Tokenizer) flags(flags token.Flags) {
	t.tok.Flags |= flags
}

func (t *Tokenizer) text(text string) {
	t.tok.Text = text
}

func (t *Tokenizer) delim() {
	t.tok.Delim = t.at(t.b)
	t.kind(token.Delim)
}

func (t *Tokenizer) Next() {
	t.skipWhitespace()
	t.next()
}

func (t *Tokenizer) NextIncludingWhitespace() {
	t.next()
}

func (t *Tokenizer) next() {
	t.tok = token.Token{}

redo:
	t.start()

	switch t.ch {
	case ' ', '\t':
		t.consumeWhitespace(false)
		t.kind(token.Whitespace)
		return
	case '\n', '\f', '\r':
		t.consumeWhitespace(true)
		t.kind(token.Whitespace)
		return

	case '[':
		t.nextch()
		t.kind(token.OpenSquareBracket)
		return
	case ']':
		t.nextch()
		t.kind(token.CloseSquareBracket)
		return

	case '(':
		t.nextch()
		t.kind(token.OpenParenthesis)
		return
	case ')':
		t.nextch()
		t.kind(token.CloseParenthesis)
		return

	case '{':
		t.nextch()
		t.kind(token.OpenCurlyBracket)
		return
	case '}':
		t.nextch()
		t.kind(token.CloseCurlyBracket)
		return

	case ',':
		t.nextch()
		t.kind(token.Comma)
		return
	case ':':
		t.nextch()
		t.kind(token.Colon)
		return
	case ';':
		if t.eof() {
			t.kind(token.EOF)
			return
		}

		t.nextch()
		t.kind(token.Semicolon)
		return

	case '/':
		t.nextch()

		if t.ch == '*' {
			t.nextch()
			t.skipComment()
			goto redo
		}

		t.delim()
		return

	case '#':
		t.nextch()
		if t.isIdentStart() || isDigit(t.ch) {
			t.start()
			t.text(t.consumeName())
			t.kind(token.Hash)
			return
		}

		t.delim()
		return

	case '@':
		t.nextch()
		if t.isIdentStart() {
			t.start()
			t.text(t.consumeName())
			t.kind(token.AtKeyword)
			return
		}
		t.delim()
		return

	case '\\':
		if t.hasAtLeast(1) && !t.hasNewlineAt(1) {
			t.consumeIdentLike()
			return
		}
		t.nextch()
		t.delim()
		return

	case '.':
		t.nextch()
		if isDigit(t.ch) {
			t.consumeNumeric(true)
			return
		}
		t.delim()
		return

	case '+':
		t.nextch()

		if isDigit(t.ch) {
			t.flags(token.HasSign)
			t.consumeNumeric(false)
			return
		}
		if t.ch == '.' && t.hasAtLeast(1) && isDigit(t.byteAt(1)) {
			t.nextch()

			t.flags(token.HasSign)
			t.consumeNumeric(true)
			return
		}

		t.delim()
		return

	case '-':
		t.nextch()

		switch t.ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.nextch()

			t.flags(token.HasSign)
			t.consumeNumeric(false)
			return

		case '.':
			if t.hasAtLeast(1) && t.hasDigitsAt(1) {
				t.advance(2)

				t.flags(token.HasSign)
				t.consumeNumeric(true)
				return
			}

		case '-':
			t.nextch()

			if t.ch == '>' {
				t.nextch()

				t.kind(token.CDC)
				return
			}

			t.flags(token.IsCustom)
			t.consumeIdentLike()
			return

		default:
			if isIdentStartLUT(t.ch) || t.ch == '\\' && (!t.hasAtLeast(1) || !t.hasNewlineAt(1)) {
				t.consumeIdentLike()
				return
			}
		}

		t.delim()
		return

	case '<':
		t.nextch()
		if t.hasAtLeast(2) && t.ch == '!' && t.byteAt(1) == '-' && t.byteAt(2) == '-' {
			t.advance(3)
			t.kind(token.CDO)
			return
		}

		t.delim()
		return

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t.consumeNumeric(false)
		return

	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'_', '\x00':
		t.consumeIdentLike()
		return

	case '"', '\'':
		t.tok.Text, t.tok.Kind = t.consumeString(t.ch)
		return

	default:
		if t.ch >= utf8.RuneSelf {
			t.consumeIdentLike()
			return
		}

		t.nextch()
		t.delim()
		return
	}
}

func (t *Tokenizer) consumeWhitespace(newline bool) {
	if newline {
		t.skipNewline()
	} else {
		t.nextch()
	}

	for {
		switch t.ch {
		case ' ', '\t':
			t.nextch()
		case '\n', '\f', '\r':
			t.skipNewline()
		default:
			return
		}
	}
}

func (t *Tokenizer) consumeNumeric(seenPoint bool) {
	if !seenPoint {
		t.skipDigits()
		if t.ch == '.' {
			t.nextch()
			seenPoint = true
		}
	}

	if seenPoint {
		t.skipDigits()
	}

	isFloat := seenPoint
	pos := t.position()
	if lower(t.ch) == 'e' {
		t.nextch()
		if isDigit(t.ch) {
			isFloat = true

			t.skipDigits()
		} else if (t.ch == '+' || t.ch == '-') && t.hasAtLeast(1) && t.hasDigitsAt(1) {
			isFloat = true

			t.nextch()
			t.skipDigits()
		} else {
			t.rewindTo(pos)
		}
	}

	t.tok.Number, _ = strconv.ParseFloat(t.segment(), 64)

	if !isFloat {
		t.tok.Flags |= token.IsInteger
	}

	if t.ch == '%' {
		t.nextch()
		t.tok.Number *= 1 / 100.0
		t.kind(token.Percentage)
		return
	}

	if t.isIdentStart() {
		t.start()
		t.text(t.consumeName())
		t.kind(token.Dimension)
		return
	}

	t.kind(token.Number)
}

func (t *Tokenizer) consumeIdentLike() {
	t.text(t.consumeName())
	if t.ch == '(' {
		t.nextch()
		if ascii.EqualFold(t.tok.Text, "url") {
			t.consumeUrl()
			return
		}

		t.kind(token.Function)
		return
	}

	t.kind(token.Ident)
}

// TODO: Implement URL tokenization.
func (t *Tokenizer) consumeUrl() {
	for {
		switch t.ch {
		case ' ', '\t':
			t.nextch()
		case '\n', '\f', '\r':
			t.skipNewline()
		case '"', '\'':
			t.consumeString(t.ch)
		case '\\':
			t.consumeEscape()
		case ')':
			t.nextch()
			t.kind(token.BadUrl)
			return
		case ';':
			if t.eof() {
				t.kind(token.BadUrl)
				return
			}
		}
	}
}

func (t *Tokenizer) consumeName() string {
	for {
		if isIdentContinueBoolLUT(t.ch) {
			t.nextch()
			continue
		}

		switch t.ch {
		case '\\':
			if !t.hasAtLeast(1) {
				t.buf = append(t.buf[:0], t.segment()...)
				t.buf = utf8.AppendRune(t.buf, utf8.RuneError)

				t.nextch()
				return string(t.buf)
			}

			if t.hasNewlineAt(1) {
				return t.segment()
			}

			t.buf = append(t.buf[:0], t.segment()...)
			t.nextch()

			if r := t.consumeEscape(); r >= 0 {
				t.buf = utf8.AppendRune(t.buf, r)
			}
		case '\x00':
			t.buf = append(t.buf[:0], t.segment()...)
			t.buf = utf8.AppendRune(t.buf, utf8.RuneError)
			t.nextch()
		default:
			return t.segment()
		}

		break
	}

	for {
		if isIdentContinueBoolLUT(t.ch) {
			t.buf = append(t.buf, t.ch)

			t.nextch()
			continue
		}

		switch t.ch {
		case '\\':
			if !t.hasAtLeast(1) {
				t.buf = utf8.AppendRune(t.buf, utf8.RuneError)
				return string(t.buf)
			}

			if t.hasNewlineAt(1) {
				return string(t.buf)
			}

			t.nextch()

			if r := t.consumeEscape(); r >= 0 {
				t.buf = utf8.AppendRune(t.buf, r)
			}

		case '\x00':
			t.nextch()
			t.buf = utf8.AppendRune(t.buf, utf8.RuneError)

		default:
			return string(t.buf)
		}
	}
}

func (t *Tokenizer) consumeEscape() rune {
	switch {
	case t.eof():
		return utf8.RuneError

	case isHex(t.ch):
		r := rune(hexlut.Dec[t.ch])
		t.nextch()

		for range 5 {
			if !isHex(t.ch) {
				break
			}

			r = r<<4 | rune(hexlut.Dec[t.ch])

			t.nextch()
		}

		switch t.ch {
		case ' ', '\t':
			t.nextch()
		case '\n', '\f', '\r':
			t.skipNewline()
		}

		switch {
		case r == 0, r > utf8.MaxRune:
			return utf8.RuneError
		case 0xD800 <= r && r <= 0xDFFF:
			return utf8.RuneError
		default:
			return r
		}

	case t.ch == '\x00':
		return utf8.RuneError

	default:
		r := rune(t.ch)
		t.nextch()
		return r
	}
}

var stringSpecialLUT = [256]bool{
	'\x00': true,
	'\\':   true,
	'\n':   true,
	'\r':   true,
	'\f':   true,
	'\'':   true,
	'"':    true,
	';':    true, // eof check
}

func (t *Tokenizer) consumeString(quote byte) (string, token.Kind) {
	t.nextch()
	t.start()

	for {
		if !stringSpecialLUT[t.ch] {
			t.nextch()
			continue
		}

		switch t.ch {
		case '"', '\'':
			if t.ch == quote {
				text := t.segment()
				t.nextch()
				return text, token.String
			}
			t.nextch()
			continue

		case '\n', '\f', '\r':
			text := t.segment()
			t.skipNewline()
			return text, token.BadString

		case '\\':
			pos := t.position()
			t.nextch()

			if t.eof() {
				return t.segmentTo(pos), token.String
			}

			if t.isNewline() {
				t.skipNewline()

				if t.eof() || t.ch == quote {
					return t.segmentTo(pos), token.String
				}

				t.buf = append(t.buf[:0], t.segmentTo(pos)...)
				break
			}

			t.buf = append(t.buf[:0], t.segmentTo(pos)...)
			if r := t.consumeEscape(); r >= 0 {
				t.buf = utf8.AppendRune(t.buf, r)
			}

		case '\x00':
			t.buf = append(t.buf[:0], t.segment()...)
			t.buf = utf8.AppendRune(t.buf, utf8.RuneError)
			t.nextch()

		case ';':
			if t.eof() {
				return t.segment(), token.String
			}
			t.nextch()
			continue
		}

		break
	}

	for {
		if !stringSpecialLUT[t.ch] {
			t.buf = append(t.buf, t.ch)
			t.nextch()
			continue
		}

		switch t.ch {
		case '\n', '\f', '\r':
			t.nextch()
			return string(t.buf), token.BadString

		case '"', '\'':
			if t.ch == quote {
				t.nextch()
				return string(t.buf), token.String
			}

			t.buf = append(t.buf, t.ch)
			t.nextch()

		case '\\':
			t.nextch()
			if t.eof() {
				return string(t.buf), token.String
			}

			if t.isNewline() {
				t.skipNewline()
			} else if r := t.consumeEscape(); r >= 0 {
				t.buf = utf8.AppendRune(t.buf, r)
			}

		case '\x00':
			t.nextch()
			t.buf = utf8.AppendRune(t.buf, utf8.RuneError)

		case ';':
			if t.eof() {
				return string(t.buf), token.String
			}

			t.buf = append(t.buf, t.ch)
			t.nextch()
		}
	}
}
