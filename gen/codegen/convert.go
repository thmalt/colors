package codegen

import (
	"fmt"
	"log"
	"math"
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateConvertPkg(ctx *Context) {
	if ctx.ConvertPkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.Opts.FormatSource)

	pkg := ctx.ConvertPkg.Name
	pkgPath := filepath.Join(ctx.Directory, ctx.ConvertPkg.Path)

	genConvertPkgConversionFiles(ctx, w, pkg, pkgPath)

	emitGoFile(ctx, w, pkg, pkgPath, "rgb8_lut", func(w *writer.GoWriter) {
		genConvertPkgLUT(w, math.MaxUint8, "LinearSrgb", "Rgb", srgbToLinearSrgb)
	})

	emitGoFile(ctx, w, pkg, pkgPath, "whitepoint", func(w *writer.GoWriter) {
		genConvertPkgWhitePoint(ctx, w)
	})
}

func genConvertPkgConversionFiles(ctx *Context, w *writer.GoWriter, pkg, pkgPath string) {
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

		emitGoFile(ctx, w, pkg, pkgPath, toSnakeCase(filename), func(w *writer.GoWriter) {
			ctx.TotalConversionGenerated += genConvertPkgSpaceConversions(ctx, w, space)
		})
	}

	fmt.Println()
}
