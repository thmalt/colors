package codegen

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

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

func snakeCase(s string) string {
	if s == "" {
		return ""
	}

	const (
		lower = 1 << iota
		upper
		digit
	)

	var (
		prevClass        uint8
		appendUnderscore = false

		currClass uint8

		b strings.Builder
	)

	i := 0
	curr, size := utf8.DecodeRuneInString(s)
	i += size

	switch {
	case unicode.IsLower(curr):
		currClass = lower
	case unicode.IsUpper(curr):
		currClass = upper
	case unicode.IsDigit(curr):
		currClass = digit
	}

	b.Grow(len(s))
	for {
		var (
			next      rune
			nextClass uint8
		)

		if i < len(s) {
			next, size = utf8.DecodeRuneInString(s[i:])

			switch {
			case unicode.IsLower(next):
				nextClass = lower
			case unicode.IsUpper(next):
				nextClass = upper
			case unicode.IsDigit(next):
				nextClass = digit
			}
		}

		if unicode.IsLetter(curr) || unicode.IsDigit(curr) {
			if appendUnderscore {
				if b.Len() != 0 {
					b.WriteByte('_')
				}
				appendUnderscore = false
			} else if currClass&upper != 0 &&
				(prevClass&lower != 0 || prevClass&digit != 0 ||
					prevClass&upper != 0 && nextClass&lower != 0) {
				// currUpper && (prevLower || prevUpper && nextLower)
				b.WriteByte('_')
			}

			if currClass&upper != 0 {
				curr = unicode.ToLower(curr)
			}

			b.WriteRune(curr)
			prevClass = currClass
		} else {
			appendUnderscore = true
		}

		if i >= len(s) {
			break
		}
		i += size

		curr, currClass = next, nextClass
	}

	return b.String()
}

func pascalCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	toUpper := true
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			toUpper = true
			continue
		}

		if toUpper {
			r = unicode.ToUpper(r)
			toUpper = false
		}

		b.WriteRune(r)
	}

	return b.String()
}
