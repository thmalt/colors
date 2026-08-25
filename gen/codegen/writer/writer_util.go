package writer

import (
	"bytes"
	"fmt"
	"go/format"
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

func formatCode(pkg, tags, header string, b []byte, formatSource bool) ([]byte, error) {
	var src []byte

	header = strings.TrimSpace(header)
	if header != "" {
		src = append(src, header...)
		src = append(src, '\n', '\n')
	}

	tags = strings.TrimSpace(tags)
	if tags != "" {
		src = append(src, buildTags(tags)...)
		src = append(src, '\n', '\n')
	}

	src = append(src, "package "...)
	src = append(src, pkg...)
	src = append(src, '\n', '\n')

	src = append(src, b...)

	src = normalizeGoSource(src)

	if formatSource {
		formatted, err := format.Source(src)
		if err != nil {
			return src, err
		}
		return formatted, nil
	}

	return src, nil
}

func buildTags(tags string) string {
	return "//go:build " + strings.ReplaceAll(tags, ",", " && ")
}

func normalizeGoSource(src []byte) []byte {
	var buf bytes.Buffer
	lastBlank := false

	forEachLine(src, func(line []byte) bool {
		blank := isBlank(line)
		if blank && lastBlank {
			return true
		}

		if !blank {
			buf.Write(bytes.TrimRightFunc(line, unicode.IsSpace))
		}
		buf.WriteByte('\n')

		lastBlank = blank
		return true
	})

	b := buf.Bytes()

	// Remove extra trailing blank lines.
	for len(b) > 1 && b[len(b)-1] == '\n' && b[len(b)-2] == '\n' {
		b = b[:len(b)-1]
	}

	// Ensure exactly one trailing newline.
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return b
}

func forEachLine(p []byte, fn func(line []byte) bool) {
	start := 0

	for start < len(p) {
		end := start
		for end < len(p) && p[end] != '\n' {
			end++
		}

		line := p[start:end]
		if n := len(line); n > 0 && line[n-1] == '\r' {
			line = line[:n-1]
		}

		if !fn(line) {
			return
		}

		if end == len(p) {
			return
		}

		start = end + 1
	}

	// Empty input
	if len(p) == 0 {
		fn(nil)
	}
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
