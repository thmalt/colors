package codegen

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func toSnakeCase(input string) string {
	var b strings.Builder
	for i, r := range input {
		if r >= 'A' && r <= 'Z' && i > 0 {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func toLowerCaseFirstWord(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	n := len(r)

	end := 0
	for end < n && unicode.IsUpper(r[end]) {
		end++
	}

	if end == 0 {
		return s
	}

	if end < n && unicode.IsLower(r[end]) && end > 1 {
		end--
	}

	for i := range end {
		r[i] = unicode.ToLower(r[i])
	}

	return string(r)
}

func toUpperCaseFirstChar(s string) string {
	if s == "" {
		return ""
	}

	c, size := utf8.DecodeRuneInString(s)

	if unicode.IsUpper(c) {
		return s
	}

	return string(unicode.ToUpper(c)) + s[size:]
}
