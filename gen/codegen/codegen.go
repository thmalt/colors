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

	suffix := "_gen.go"
	if !strings.HasSuffix(filename, suffix) {
		filename += suffix
	}

	fmt.Println("  Generate file", filename)
	w.SaveGoFile(filepath.Join(path, filename), pkg)
}
