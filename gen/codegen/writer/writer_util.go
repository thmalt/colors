package writer

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

type blockKind uint

const (
	blockNone blockKind = iota
	blockFuncName
	blockFuncSignture
	blockIf
	blockSwitch
)

func (s blockKind) String() string {
	switch s {
	case blockNone:
		return "blockNone"
	case blockFuncName:
		return "blockFuncName"
	case blockFuncSignture:
		return "blockFuncSignture"
	case blockIf:
		return "blockIf"
	case blockSwitch:
		return "blockSwitch"
	default:
		return fmt.Sprintf("unknown block %d", s)
	}
}

const (
	bracketBlock = '{'
	bracketGroup = '('
	bracketArray = '['
)

func WriteGo(w io.Writer, pkg, tags string, b []byte) (int, error) {
	var src []byte

	if tags != "" {
		src = append(src, tagLines(tags)...)
		src = append(src, '\n')
	}

	src = append(src, fmt.Sprintf("package %s\n\n", pkg)...)
	src = append(src, b...)

	formatted, err := format.Source(src)

	if err != nil {
		n, _ := w.Write(src)
		return n, err
	}

	return w.Write(formatted)
}

func tagLines(tags string) string {
	return "//go:build " + strings.ReplaceAll(tags, ",", " && ") + "\n"
}

func appendArgsFormat(dst *bytes.Buffer, format string, a ...any) {
	fmt.Fprintf(dst, format, a...)
}

func appendArgs(dst *bytes.Buffer, a ...any) {
	for _, v := range a {
		switch x := v.(type) {
		case string:
			dst.WriteString(x)
		case []byte:
			dst.Write(x)
		case []rune:
			dst.WriteString(string(x))
		case byte:
			dst.WriteByte(x)
		case rune:
			dst.WriteRune(x)
		default:
			fmt.Fprint(dst, x)
		}
	}
}

func hasLeadingSpace(p []byte) bool {
	r, _ := utf8.DecodeRune(p)
	return unicode.IsSpace(r)
}

func isBlank(p []byte) bool {
	for len(p) > 0 {
		r, n := utf8.DecodeRune(p)
		if !unicode.IsSpace(r) {
			return false
		}
		p = p[n:]
	}
	return true
}

func isASCIIBlank(p []byte) bool {
	for _, c := range p {
		switch c {
		case ' ', '\t', '\r', '\n':
		default:
			return false
		}
	}
	return true
}

// for [bytes.Buffer]
func isBlankArgs(a ...any) bool {
	if len(a) == 0 {
		return true
	}

	for _, v := range a {
		switch x := v.(type) {
		case string:
			for _, r := range x {
				if !unicode.IsSpace(r) {
					return false
				}
			}
		case []byte:
			for len(x) > 0 {
				r, size := utf8.DecodeRune(x)
				if !unicode.IsSpace(r) {
					return false
				}
				x = x[size:]
			}
		case byte:
			if !unicode.IsSpace(rune(x)) {
				return false
			}
		case rune:
			if !unicode.IsSpace(x) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func trailingBlankLines(p []byte) int {
	return trailingBlankLinesN(p, math.MaxInt)
}

// Use [unicode.IsSpace] to check.
func trailingBlankLinesN(p []byte, max int) int {
	i := len(p) - 1

	count := 0
	for count < max {
		blank := true

		for i > 0 && p[i] != '\n' {
			r, n := utf8.DecodeLastRune(p[:i+1])

			if n == 0 {
				break
			}

			if !unicode.IsSpace(r) {
				blank = false
			}

			i -= n
		}

		if !blank {
			break
		}

		count++

		if i < 0 {
			break
		}

		i--
	}

	return count
}

func isStdlib(path string) bool {
	before, _, _ := strings.Cut(path, "/")
	return !strings.Contains(before, ".")
}

func splitImports(imports []string) (stdlib []string, others []string) {
	for _, pkg := range imports {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			continue
		}

		if isStdlib(pkg) {
			stdlib = append(stdlib, pkg)
		} else {
			others = append(others, pkg)
		}
	}

	return
}
