package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateSpacePkg(ctx *Context) {
	if ctx.SpacePkg.Name == "" {
		return
	}

	var pkgPath = filepath.Join(ctx.Directory, ctx.SpacePkg.Path)
	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	// generate file ctx.SpacePkg.Name + _gen.go
	fileName := ctx.SpacePkg.Name + "_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		"strconv",
		// ctx.ConvertPkg.Path,
	)

	genSpacePkgType(ctx, w)

	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.SpacePkg.Name,
	)

	fileName = "whitepoint_gen.go"
	fmt.Println("  Generate file", fileName)

	genSpacePkgWhitePoint(ctx, w)

	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.SpacePkg.Name,
	)

}

func genSpacePkgType(ctx *Context, w *writer.GoWriter) {
	w.LineWriteln("type Space uint", smallestUintType(len(ctx.BuildSpaces)))

	w.Separate()

	w.BeginGroup("const ")

	for i, space := range ctx.BuildSpaces {
		if i == 0 {
			w.LineWriteln(space.Name, " Space = iota")
			continue
		}
		w.LineWriteln(space.Name)
	}

	w.Separate()

	w.Comment("Available space count")
	w.LineWriteln("SpaceCount")
	w.End()

	w.Separate()

	ws := w.SubWriter()

	w.Begin("var spaceInfos = [...]*spaceInfo")

	for _, space := range ctx.BuildSpaces {
		name := toLowerCaseFirstWord(space.Name) + "Info"
		w.LineWriteln("&", name, ",")

		ws.Begin("var ", name, " = spaceInfo")

		ws.LineWritef("name: %q,\n", space.Name)
		ws.LineWritef("displayName: %q,\n", space.DisplayName)
		ws.LineWritef("cssName: %q,\n", space.CssName)

		if whitePoint := LookupWhitePoint(space.WhitePoint); whitePoint != nil {
			ws.LineWritef("whitePoint: %s,\n", whitePoint.Name)
		}

		ws.LineWritef("coordinate: %s,\n", space.Coordinate)

		ws.Begin("channels: []Channel")
		for _, c := range space.Channels {
			ws.Begin()

			ws.LineWritef("Name: %q,\n", c.Name)
			ws.LineWritef("Symbol: %q,\n", c.Symbol)
			ws.LineWritef("DisplayName: %q,\n", c.DisplayName)

			unit := c.Unit
			switch c.Unit {
			case model.UnitRadian, model.UnitGradian, model.UnitTurn:
				ws.LineWritef("Min: %g,\n", model.AngleToDegree(c.Min, c.Unit))
				ws.LineWritef("Max: %g,\n", model.AngleToDegree(c.Max, c.Unit))
				unit = model.UnitDegree
			default:
				ws.LineWritef("Min: %g,\n", c.Min)
				ws.LineWritef("Max: %g,\n", c.Max)

			}

			if c.Circular {
				ws.LineWritef("Circular: %t,\n", c.Circular)
			}

			if c.Unrestricted {
				ws.LineWritef("Unrestricted: %t,\n", c.Unrestricted)
			}

			if c.Unit != model.UnitNumber {
				ws.LineWritef("Unit: %s,\n", unit.GoString())
			}

			prec := c.Precision
			if prec == 0 {
				prec = DefaultPrecision
			}

			ws.LineWritef("Precision: %d,\n", prec)

			ws.End(',')
		}
		ws.End(',')

		if space.UseColorFunction {
			ws.LineWriteln("useColorFunction: true,")
		}

		ws.End()

		ws.Separate()
	}

	w.End()

	w.Separate()

	w.Drain(ws)

	w.Separate()

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
	w.Begin("type WhitePoint struct ")
	w.LineWriteln("Name string")

	w.Separate()

	w.LineWriteln("X ", FloatType)
	w.LineWriteln("Y ", FloatType)
	w.LineWriteln("Z ", FloatType)
	w.End()

	w.Separate()

	w.BeginGroup("var ")

	for i, whitePoint := range ctx.WhitePoints {
		if i > 0 {
			w.Separate()
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
