package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateSpacePkg(ctx *Context) {
	var pkgPath = filepath.Join(ctx.Directory, ctx.SpacePkg.Path)
	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))

	// generate file ctx.SpacePkg.Name + _gen.go
	fileName := ctx.SpacePkg.Name + "_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		"strconv",
		// ctx.ConvertPkg.Path,
	)

	genSpacePkgType(ctx, w)

	w.WriteGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.SpacePkg.Name,
	)

	fileName = "whitepoint_gen.go"
	fmt.Println("  Generate file", fileName)

	genSpacePkgWhitePoint(ctx, w)

	w.WriteGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.SpacePkg.Name,
	)

}

func genSpacePkgType(ctx *Context, w *writer.GoWriter) {
	w.LineWriteln("type Space uint", smallestUintType(len(ctx.BuildSpaces)))

	w.Writeln()

	w.BeginGroup("const")

	for i, space := range ctx.BuildSpaces {
		if i == 0 {
			w.LineWriteln(space.Name, " Space = iota")
			continue
		}
		w.LineWriteln(space.Name)
	}

	w.Writeln()
	w.LineCommentln("Available space count")
	w.LineWriteln("SpaceCount")
	w.End()

	w.Writeln()

	w.BeginGroup("var")
	w.Begin("spaceInfos = [...]SpaceInfo")
	for _, space := range ctx.BuildSpaces {
		w.Begin()

		w.LineWritef("name: %q,\n", space.Name)
		w.LineWritef("displayName: %q,\n", space.DisplayName)
		w.LineWritef("cssName: %q,\n", space.CssName)

		if whitePoint := LookupWhitePoint(space.WhitePoint); whitePoint != nil {
			w.LineWritef("whitePoint: %s,\n", whitePoint.Name)
		}

		w.Begin("channels: []Channel")
		for _, c := range space.Channels {
			w.Begin()

			w.LineWritef("Name: %q,\n", c.Name)
			w.LineWritef("Symbol: %q,\n", c.Symbol)
			w.LineWritef("DisplayName: %q,\n", c.DisplayName)

			unit := c.Unit
			switch c.Unit {
			case model.UnitRadian, model.UnitGradian, model.UnitTurn:
				w.LineWritef("Min: %g,\n", model.AngleToDegree(c.Min, c.Unit))
				w.LineWritef("Max: %g,\n", model.AngleToDegree(c.Max, c.Unit))
				unit = model.UnitDegree
			default:
				w.LineWritef("Min: %g,\n", c.Min)
				w.LineWritef("Max: %g,\n", c.Max)

			}

			if c.Circular {
				w.LineWritef("Circular: %t,\n", c.Circular)
			}

			if c.Unrestricted {
				w.LineWritef("Unrestricted: %t,\n", c.Unrestricted)
			}

			if c.Unit != model.UnitNumber {
				w.LineWritef("Unit: %s,\n", unit.GoString())
			}

			prec := c.Precision
			if prec == 0 {
				prec = DefaultPrecision
			}

			w.LineWritef("Precision: %d,\n", prec)

			w.End(',')
		}
		w.End(',')

		if space.UseColorFunction {
			w.LineWriteln("useColorFunction: true,")
		}

		w.End(',')
	}
	w.End()

	w.End()

	w.Writeln()

	w.Method("s Space", "String")
	w.FuncResults("string")
	w.FuncBody()

	w.Switch("s")

	for _, space := range ctx.BuildSpaces {
		w.Case(space.Name)
		w.Return('"', space.Name, '"')
	}

	w.Default()

	w.Return(`"Space(" + strconv.FormatUint(uint64(s), 10) + ")"`)

	w.End()

	w.End()
}

func genSpacePkgWhitePoint(ctx *Context, w *writer.GoWriter) {
	w.Begin("type WhitePoint struct")
	w.LineWriteln("Name string")
	w.Writeln()
	w.LineWriteln("X ", FloatType)
	w.LineWriteln("Y ", FloatType)
	w.LineWriteln("Z ", FloatType)
	w.End()

	w.Writeln()

	w.BeginGroup("var")

	for i, whitePoint := range ctx.WhitePoints {
		if i > 0 {
			w.Writeln()
		}

		xyz := data.ChromaToXyz(whitePoint.X, whitePoint.Y)

		w.Begin(whitePoint.Name, " = WhitePoint")

		w.LineWritef("Name: %q,\n", whitePoint.Name)
		w.LineWritef("X: %s,\n", formatNormalizedFloat(xyz[0]))
		w.LineWritef("Y: %s,\n", formatNormalizedFloat(xyz[1]))
		w.LineWritef("Z: %s,\n", formatNormalizedFloat(xyz[2]))

		w.End()
	}

	w.End()
}
