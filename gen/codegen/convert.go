package codegen

import (
	"fmt"
	"log"
	"math"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateConvertPkg(ctx *Context) {
	pkg := ctx.ConvertPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	genConvertPkgConversionFiles(ctx, w, pkg)

	emitGoFile(ctx, pkg, w, "rgb8_lut", func(w *writer.GoWriter) {
		genConvertPkgLUT(w, math.MaxUint8, "LinearSrgb", "Rgb", srgbToLinearSrgb)
	})

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

		emitGoFile(ctx, pkg, w, toSnakeCase(filename), func(w *writer.GoWriter) {
			ctx.TotalConversionGenerated += genConvertPkgSpaceConversions(ctx, w, space)
		})
	}

	fmt.Println()
}
