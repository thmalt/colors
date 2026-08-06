package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func emitGoFile(w *writer.GoWriter, pkg, path, filename string, fn func(w *writer.GoWriter)) {
	w.Reset()

	fn(w)

	filename = normalizeGenFilename(filename)

	fmt.Println("  Generate file", filename)
	w.SaveGoFile(filepath.Join(path, filename), pkg)
}

func normalizeGenFilename(filename string) string {
	const suffix = "_gen.go"

	if strings.HasSuffix(filename, suffix) {
		return filename
	}

	filename = strings.TrimSuffix(filename, ".go")
	filename = strings.TrimSuffix(filename, "_gen")

	return filename + suffix
}
