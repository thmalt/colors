package codegen

import (
	"fmt"
	"log"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateConvertPkg(ctx *Context) {
	pkg := ctx.ConvertPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	genConvertPkgConversionFiles(ctx, w, pkg)

	if ctx.Opts.genLUT() {
		emitGoFile(ctx, pkg, w, "srgb_lut_static", func(w *writer.GoWriter) {
			w.AddBuildTags("static || st || colors_full")

			genConvertPkgPrecomputedLUT(w, 8)
			genConvertPkgPrecomputedLUT(w, 16)
		})

		emitGoFile(ctx, pkg, w, "srgb_lut_init", func(w *writer.GoWriter) {
			w.AddBuildTags("!static && !st && !colors_full")

			genConvertPkgRuntimeInitLUT(w, 8)
			genConvertPkgRuntimeInitLUT(w, 16)
		})
	}

	emitGoFile(ctx, pkg, w, "whitepoint", func(w *writer.GoWriter) {
		genConvertPkgWhitePoint(ctx, w)
	})
}

func genConvertPkgConversionFiles(ctx *Context, w *writer.GoWriter, pkg Pkg) {
	ctx.TotalConversionGenerated = 0
	for i, space := range ctx.BuiltSpaces {
		if space == nil {
			log.Printf("space at index %d is nil\n", i)
			continue
		}

		if space.Disable {
			continue
		}

		filename := space.Name
		if space.SnakeName != "" {
			filename = space.SnakeName
		}

		emitGoFile(ctx, pkg, w, snakeCase(filename), func(w *writer.GoWriter) {
			ctx.TotalConversionGenerated += genConvertPkgSpaceConversions(ctx, w, space)
		})
	}

	fmt.Println()
}
