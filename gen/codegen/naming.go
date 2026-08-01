package codegen

import (
	"strings"
	"unicode"
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
