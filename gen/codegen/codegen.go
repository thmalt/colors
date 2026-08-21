package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func emitGoFile(ctx *Context, w *writer.GoWriter, pkg string, path, filename string, fn func(w *writer.GoWriter)) {
	w.Reset()

	fn(w)

	filename = normalizeGenFilename(filename)

	written := w.SaveGoFile(filepath.Join(path, filename), pkg, ctx.Opts.BuildTags, ctx.Opts.ForceWrite)
	if written {
		fmt.Println("  Generate file", filename)
	}
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
