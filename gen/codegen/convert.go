package codegen

import (
	"fmt"
	"log"

	"github.com/thmalt/colors/gen/codegen/model"
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

	type spaceAlias struct {
		Alias string
		Space *model.Space
	}

	var spaceAliases []spaceAlias

	for i, space := range ctx.BuiltSpaces {
		if space == nil {
			log.Printf("space at index %d is nil\n", i)
			continue
		}

		if len(space.Aliases) > 0 {
			for _, alias := range space.Aliases {
				spaceAliases = append(spaceAliases, spaceAlias{
					Alias: alias,
					Space: space,
				})
			}
		}

		filename := space.Name
		if space.SnakeName != "" {
			filename = space.SnakeName
		}

		emitGoFile(ctx, pkg, w, snakeCase(filename), func(w *writer.GoWriter) {
			ctx.TotalConversionGenerated += genConvertPkgSpaceConversions(ctx, w, space)
		})
	}

	aliasCount := 0
	emitGoFile(ctx, pkg, w, "aliases", func(w *writer.GoWriter) {
		// alias <-> alias
		for i, spaceAlias := range spaceAliases {
			paramsVars := spaceAlias.Space.ChannelIdent()

			alias := spaceAlias.Alias
			name := spaceAlias.Space.Name

			for _, other := range spaceAliases[i+1:] {
				if name == other.Space.Name {
					continue
				}

				resultsVars := other.Space.ChannelIdent()

				aliasCount += genConvertPkgAliasFunc(
					ctx, w,
					Pair{alias, other.Alias},
					Pair{spaceAlias.Space.Name, other.Space.Name},
					paramsVars, resultsVars,
				)
				aliasCount += genConvertPkgAliasFunc(
					ctx, w,
					Pair{other.Alias, alias},
					Pair{other.Space.Name, spaceAlias.Space.Name},
					resultsVars, paramsVars,
				)
			}
		}

		// alias <-> canonical
		for _, spaceAlias := range spaceAliases {
			paramsVars := spaceAlias.Space.ChannelIdent()

			alias := spaceAlias.Alias
			name := spaceAlias.Space.Name

			for _, space := range ctx.BuiltSpaces {
				if name == space.Name {
					continue
				}

				resultsVars := space.ChannelIdent()

				aliasCount += genConvertPkgAliasFunc(
					ctx, w,
					Pair{alias, space.Name},
					Pair{name, space.Name},
					paramsVars, resultsVars,
				)
				aliasCount += genConvertPkgAliasFunc(
					ctx, w,
					Pair{space.Name, alias},
					Pair{space.Name, name},
					resultsVars, paramsVars,
				)
			}
		}
	})

	ctx.TotalConversionGenerated += aliasCount

	fmt.Println()
}
