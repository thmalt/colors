package codegen

import (
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateMixerPkg(ctx *Context) {
	if ctx.MixerPkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.Opts.FormatSource)

	pkg := ctx.MixerPkg.Name
	pkgPath := filepath.Join(ctx.Directory, ctx.MixerPkg.Path)

	emitGoFile(ctx, w, pkg, pkgPath, "unsafe", func(w *writer.GoWriter) {
		w.Import(ctx.InterpPkg.Path)

		genMixerPkgUnsafe(ctx, w)
	})
}
