package codegen

import (
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateRootPkg(ctx *Context) {
	if ctx.RootPkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.Opts.FormatSource)

	pkg := ctx.RootPkg.Name
	pkgPath := filepath.Join(ctx.Directory, ctx.RootPkg.Path)

	emitGoFile(w, pkg, pkgPath, "color", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgColor(ctx, w)
		genRootPkgColorChannel(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_convert", func(w *writer.GoWriter) {
		w.Import(
			ctx.ConvertPkg.Path,
			ctx.SpacePkg.Path,
		)

		genRootPkgColorTo(ctx, w)
		genRootPkgColorConvertMethods(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_constructors", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgColorConstructors(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_string", func(w *writer.GoWriter) {
		w.Import(
			"strconv",
			"unsafe",
			ctx.SpacePkg.Path,
		)

		genRootPkgColorStringMethod(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "mixer", func(w *writer.GoWriter) {
		genRootPkgMixerMethod(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "clamp", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgClamp(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "gamut", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgGamut(ctx, w)
	})
}
