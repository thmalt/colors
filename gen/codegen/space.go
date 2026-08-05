package codegen

import (
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateSpacePkg(ctx *Context) {
	if ctx.SpacePkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	pkgPath := filepath.Join(ctx.Directory, ctx.SpacePkg.Path)
	pkg := ctx.SpacePkg.Name

	emitGoFile(w, pkg, pkgPath, pkg, func(w *writer.GoWriter) {
		w.Import(
			"strconv",
		)

		genSpacePkgType(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, pkg+"_info", func(w *writer.GoWriter) {
		genSpacePkgSpaceInfo(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "whitepoint", func(w *writer.GoWriter) {
		genSpacePkgWhitePoint(ctx, w)
	})
}

func genSpacePkgType(ctx *Context, w *writer.GoWriter) {
	w.LineWriteln("type Space uint", smallestUintType(len(ctx.BuildSpaces)))

	w.Separate()

	w.BeginGroup("const ")
	w.LineWriteln("SpaceInvalid Space = iota")

	w.Separate()
	var spaceFirst string
	var maxChannelCnt = 0
	for i, space := range ctx.BuildSpaces {
		if i == 0 {
			spaceFirst = space.Name
		}

		if space.Comment != "" {
			if ctx.SeparateAfterComment {
				w.Separate()
			}

			w.Comment(space.Comment)
		}
		w.LineWriteln(space.Name)

		maxChannelCnt = max(maxChannelCnt, space.ChannelCount())
	}

	w.Separate()

	w.Comment("Available space count")
	w.LineWriteln("SpaceCount")
	w.End()
	w.Separate()
	w.LineWriteln("const SpaceFirst = ", spaceFirst)

	w.Separate()

	w.Begin("var spaceInfos = [...]*spaceInfo")
	w.LineWriteln("nil,")
	for _, space := range ctx.BuildSpaces {
		w.LineWriteln("&", spaceInfoName(space), ",")
	}

	w.End()

	w.Separate()
	w.Begin("var spaceChannelCounts = [...]uint", smallestUintType(maxChannelCnt))
	w.LineWriteln("0,")
	for _, space := range ctx.BuildSpaces {
		w.LineWriteln(space.ChannelCount(), ",")
	}
	w.End()

	w.Separate()
	w.Begin("var coordinateSystems = [...]CoordinateSystem")
	w.LineWriteln(model.CoordinateSystem(0), ",")
	for _, space := range ctx.BuildSpaces {
		w.LineWriteln(space.Coordinate, ",")
	}
	w.End()

	w.Separate()
	w.Method("s Space", "String")
	w.FuncResults("string")
	w.FuncBody()

	w.Switch("s ")
	w.Case("SpaceInvalid")
	w.Return(`"Invalid"`)

	for _, space := range ctx.BuildSpaces {
		w.Case(space.Name)
		w.Return('"', space.Name, '"')
	}

	w.Default()

	w.Return(`"Space(" + strconv.FormatUint(uint64(s), 10) + ")"`)

	w.End()

	w.End()
}

func spaceInfoName(space *model.Space) string {
	return toLowerCaseFirstWord(space.Name) + "Info"
}

func genSpacePkgSpaceInfo(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuildSpaces {
		w.Begin("var ", spaceInfoName(space), " = spaceInfo")

		w.LineWritef("name: %q,\n", space.Name)
		w.LineWritef("displayName: %q,\n", space.DisplayName)
		w.LineWritef("cssName: %q,\n", space.CssName)

		if whitePoint := LookupWhitePoint(space.WhitePoint); whitePoint != nil {
			w.LineWritef("whitePoint: %s,\n", whitePoint.Name)
		}

		w.LineWritef("coordinate: %s,\n", space.Coordinate)

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
				w.LineWritef("Unit: %s,\n", unit)
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

		w.End()

		w.Separate()
	}
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
