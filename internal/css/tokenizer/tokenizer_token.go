package tokenizer

import (
	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
)

func (t *Tokenizer) Is(kind token.Kind) bool {
	return t.tok.Kind == kind
}

func (t *Tokenizer) Kind() token.Kind {
	return t.tok.Kind
}

func (t *Tokenizer) Text() string {
	return t.tok.Text
}

func (t *Tokenizer) Number() float64 {
	return t.tok.Number
}

func (t *Tokenizer) Delim() byte {
	return t.tok.Delim
}

func (t *Tokenizer) Unit() string {
	return t.tok.Unit()
}

func (t *Tokenizer) Flags() token.Flags {
	return t.tok.Flags
}

func (t *Tokenizer) ExpectIdent(value string) bool {
	return t.tok.Is(token.Ident) && ascii.EqualFold(t.Text(), value)
}

func (t *Tokenizer) ExpectDelim(delim byte) bool {
	return t.tok.Is(token.Delim) && t.tok.Delim == delim
}

func (t *Tokenizer) ExpectCloseParenthesis() bool {
	return t.tok.Is(token.CloseParenthesis)
}

func (t *Tokenizer) ExpectComma() bool {
	return t.tok.Is(token.Comma)
}

func (t *Tokenizer) NextIf(ok bool) bool {
	if ok {
		t.Next()
	}
	return ok
}

func (t *Tokenizer) NextExpect(kind token.Kind) bool {
	if !t.Is(kind) {
		return false
	}

	t.Next()
	return true
}
