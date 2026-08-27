package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateMixerPkg(ctx *Context) {
	pkg := ctx.MixerPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	emitGoFile(ctx, pkg, w, "unsafe", func(w *writer.GoWriter) {
		w.Import(ctx.InterpPkg.Path)

		genMixerPkgUnsafe(ctx, w)
	})
}
