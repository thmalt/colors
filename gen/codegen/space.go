package codegen

import (
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateSpacePkg(ctx *Context) {
	var pkgPath = filepath.Join(ctx.Directory, ctx.SpacePkg.Path)
	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))

	w.Import(
		"strconv",
		ctx.ConvertPkg.Path,
	)

	genSpacePkgType(ctx, w)

	w.WriteGoFile(
		filepath.Join(pkgPath, ctx.SpacePkg.Name+"_gen.go"),
		ctx.SpacePkg.Name,
	)
}

func genSpacePkgType(ctx *Context, w *writer.GoWriter) {
	w.LineWriteln("type Space uint", smallestUintType(len(ctx.Spaces)))

	w.Writeln()

	w.BeginGroup("const")
	lastSpace := ""
	for i, space := range ctx.Spaces {
		if i == 0 {
			w.LineWriteln(space.Name, " Space = iota")
			continue
		}
		w.LineWriteln(space.Name)
		lastSpace = space.Name
	}
	w.End()

	w.Writeln()

	w.BeginGroup("var")
	w.Begin("infos = [...]SpaceInfo")
	for _, space := range ctx.Spaces {
		w.Begin()

		w.LineWritef("name: %q,\n", space.Name)
		w.LineWritef("displayName: %q,\n", space.DisplayName)
		w.LineWritef("cssName: %q,\n", space.CssName)
		w.LineWritef("whitePoint: %s,\n", ctx.ConvertPkg.Join(space.WhitePoint))

		w.Begin("channels: []Channel")
		for _, c := range space.Channels {
			w.Begin()

			w.LineWritef("Name: %q,\n", c.Ident)
			w.LineWritef("Symbol: %q,\n", c.Symbol)
			w.LineWritef("DisplayName: %q,\n", c.DisplayName)

			switch c.Unit {
			case model.UnitRadian, model.UnitGradian, model.UnitTurn:
				w.LineWritef("Min: %v,\n", model.AngleToDegree(c.Min, c.Unit))
				w.LineWritef("Max: %v,\n", model.AngleToDegree(c.Max, c.Unit))
				if c.Circular {
					w.LineWritef("Circular: %v,\n", c.Circular)
				}
				w.LineWritef("Unit: %s,\n", model.UnitDegree.GoString())
			default:
				w.LineWritef("Min: %v,\n", c.Min)
				w.LineWritef("Max: %v,\n", c.Max)
				if c.Circular {
					w.LineWritef("Circular: %v,\n", c.Circular)
				}
				if c.Unit != model.UnitNumber {
					w.LineWritef("Unit: %s,\n", c.Unit.GoString())
				}
			}

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

	w.Method("s Space", "Info")
	w.FuncParams()
	w.FuncResults("*SpaceInfo")
	w.FuncBody()
	w.If("s > ", lastSpace)
	w.Return("nil")
	w.End()
	w.Return("&infos[s]")
	w.End()

	w.Writeln()

	w.Method("s Space", "String")
	w.FuncResults("string")
	w.FuncBody()
	w.Switch("s")
	for _, space := range ctx.Spaces {
		w.Case(space.Name)
		w.Return("\"", space.Name, "\"")
	}
	w.Default()
	w.Return("strconv.FormatUint(uint64(s), 10)")
	w.End()
	w.End()
}
