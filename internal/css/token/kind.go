package token

type Kind uint8

//go:generate stringer -type Kind -linecomment kind.go

const (
	_ Kind = iota

	EOF                // <EOF-token>
	Ident              // <ident-token>
	Function           // <function-token>
	AtKeyword          // <at-keyword-token>
	Hash               // <hash-token>
	String             // <string-token>
	BadString          // <bad-string-token>
	Url                // <url-token>
	BadUrl             // <bad-url-token>
	Delim              // <delim-token>
	Number             // <number-token>
	Percentage         // <percentage-token>
	Dimension          // <dimension-token>
	Whitespace         // <whitespace-token>
	CDO                // <CDO-token>
	CDC                // <CDC-token>
	Colon              // <colon-token>
	Semicolon          // <semicolon-token>
	Comma              // <comma-token>
	OpenSquareBracket  // <[-token>
	CloseSquareBracket // <]-token>
	OpenParenthesis    // <(-token>
	CloseParenthesis   // <)-token>
	OpenCurlyBracket   // <{-token>
	CloseCurlyBracket  // <}-token>
	// UnicodeRange       // <unicode-range-token>
	// Comment       // <comment-token>
)

func (k Kind) Is(other Kind) bool {
	return k == other
}
