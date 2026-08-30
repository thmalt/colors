package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func emitGoFile(ctx *Context, pkg Pkg, w *writer.GoWriter, filename string, emit func(w *writer.GoWriter)) {
	w.Reset()

	w.SetBuildTags(ctx.Opts.BuildTags)

	emit(w)

	pkgName, path := pkgInfo(ctx, pkg)

	filename = normalizeGenFilename(filename)

	written := w.SaveGoFile(filepath.Join(path, filename), pkgName, ctx.Opts.ForceWrite)
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
