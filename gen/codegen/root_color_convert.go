package codegen

import (
	"log"
	"slices"
	"strconv"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorConvertMethods(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuiltSpaces {
		genRootPkgColorConvertMethod(ctx, w, space)
	}
}

func genRootPkgColorConvertMethod(ctx *Context, w *writer.GoWriter, space *model.Space) {
	names := space.ChannelIdent()
	hasNamedReturn := slices.Contains(names, "c")

	w.Separate()
	w.Comment(space.Name, " returns the color components in the [", ctx.SpacePkg.Join(space.Name), "] color space.")
	w.Method("c Color", space.Name)
	if !hasNamedReturn {
		w.FuncResults(joinIdentsWithType(FloatType, names...))
	} else {
		w.FuncResults(joinRepeatN(FloatType, len(names)))
	}
	w.FuncBody()

	var args = make([]string, len(names))
	for i := range args {
		args[i] = "c.c" + strconv.FormatInt(int64(i+1), 10)
	}

	w.If("c.space == ", ctx.SpacePkg.Join(space.Name))
	w.ReturnInline()
	w.WriteJoin(args, ", ")
	w.End()

	w.Separate()

	sub := w.SubWriter()

	sub.Switch("c.space")

	var foundPath = false
	for _, src := range ctx.BuiltSpaces {
		if space == src {
			continue
		}

		path := ctx.Graph.FindPath(src, space)
		if len(path) == 0 {
			log.Printf("no conversion path found: %q -> %q\n", src.Name, space.Name)
			continue
		}
		foundPath = true
		sub.Case(ctx.SpacePkg.Join(src.Name))

		sub.ReturnInline()

		pair := Pair{src.Name, space.Name}

		fnName := pair.FuncName()
		if fn := ctx.ConvertFuncByPair(pair); fn != nil {
			fnName = fn.Pair.FuncName()
		}

		sub.WriteCallln(ctx.ConvertPkg.Join(fnName), args...)
	}
	sub.Default()
	if !hasNamedReturn {
		sub.Return()
	} else {
		sub.Return(joinRepeatN("0", len(names)))
	}
	sub.End()

	if foundPath {
		w.Drain(sub)
	} else {
		if !hasNamedReturn {
			w.Return()
		} else {
			w.Return(joinRepeatN("0", len(names)))
		}
	}

	w.End()
}

func genRootPkgColorTo(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg
	w.Separate()
	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()
	w.If("!c.space.IsValid() || !dst.IsValid()")
	w.Return("Color{}, ErrInvalidSpace")
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return("c, nil")
	w.End()

	w.Separate()
	w.Return("c.to(dst)")
	w.End()

	w.Separate()
	w.Method("c Color", "to")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()

	var scope VariableScope

	w.Switch("dst")
	for _, space := range ctx.BuiltSpaces {
		scope.Reset()
		scope.Reserve("c")

		spaceIdent := spacePkg.Join(space.Name)

		w.Case(spaceIdent)
		names := scope.ReserveUniqueAll(space.ChannelIdent()...)
		w.LineWriteJoin(names, ", ")
		w.Writef(" := c.%s()\n", space.Name)

		w.ReturnInline("Color{")
		w.Write("space: ", ctx.SpacePkg.Join(space.Name), ", ")
		for i, name := range names {
			w.Write("c", i+1, ": ", name, ", ")
		}
		w.Write("alpha: c.alpha")
		w.Writeln("}, nil")
	}
	w.Default()
	w.Return("Color{}, ErrInvalidSpace")
	w.End()

	w.End()
}
