package codegen

import (
	"fmt"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateRootPkg(ctx *Context) {
	pkg := ctx.RootPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	emitGoFile(ctx, pkg, w, "color", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgColor(ctx, w)
		genRootPkgColorChannel(ctx, w)
		genRootPkgHexLUT(ctx, w)
	})

	var stats conversionStats
	emitGoFile(ctx, pkg, w, "color_convert", func(w *writer.GoWriter) {
		w.Import(
			ctx.ConvertPkg.Path,
			ctx.SpacePkg.Path,
		)

		w.AddBuildTags("!colors_full")

		stats = genRootPkgColorConvertMethods(ctx, w, false)
	})

	fmt.Printf("    unique conversions: direct=%d, hub=%d\n",
		stats.DirectCounts(),
		stats.HubCounts(),
	)

	emitGoFile(ctx, pkg, w, "color_convert_full", func(w *writer.GoWriter) {
		w.Import(
			ctx.ConvertPkg.Path,
			ctx.SpacePkg.Path,
		)

		w.AddBuildTags("colors_full")

		stats = genRootPkgColorConvertMethods(ctx, w, true)
	})

	fmt.Printf("    unique conversions: direct=%d, hub=%d\n",
		stats.DirectCounts(),
		stats.HubCounts(),
	)

	emitGoFile(ctx, pkg, w, "color_mutable", func(w *writer.GoWriter) {
		w.Import(
			ctx.SpacePkg.Path,
		)

		genRootPkgColorMutTo(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "color_constructors", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgColorConstructors(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "color_string", func(w *writer.GoWriter) {
		w.Import(
			"strconv",
			"unsafe",
			ctx.SpacePkg.Path,
		)

		genRootPkgColorStringMethod(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "mixer", func(w *writer.GoWriter) {
		genRootPkgMixerMethod(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "gamut", func(w *writer.GoWriter) {
		w.Import(ctx.SpacePkg.Path)

		genRootPkgGamut(ctx, w)
		genRootPkgClamp(ctx, w)
	})
}
