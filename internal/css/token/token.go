package token

type Token struct {
	Text   string
	Number float64
	Delim  byte

	Kind  Kind
	Flags Flags
}

type Flags uint8

const (
	IsInteger Flags = 1 << iota // integer
	IsCustom                    // --
	HasEscape                   // \
	HasSign                     // + -
)

func (t Token) IsInteger() bool { return t.Flags&IsInteger != 0 }
func (t Token) IsCustom() bool  { return t.Flags&IsCustom != 0 }
func (t Token) HasSign() bool   { return t.Flags&HasSign != 0 }

func (f Flags) IsInteger() bool { return f&IsInteger != 0 }
func (f Flags) IsCustom() bool  { return f&IsCustom != 0 }
func (f Flags) HasSign() bool   { return f&HasSign != 0 }

func (t Token) Is(kind Kind) bool { return t.Kind == kind }

func (t Token) Unit() string {
	if t.Kind != Dimension {
		return ""
	}
	return t.Text
}
