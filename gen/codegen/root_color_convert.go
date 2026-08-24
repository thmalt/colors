package codegen

import (
	"log"
	"slices"
	"strconv"
	"strings"

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

	if eq := ctx.SpaceByName(space.Equivalent); eq != nil {
		w.Return("c.", eq.Name, "()")
		w.End()
		return
	}

	var args = make([]string, len(names))
	for i := range args {
		args[i] = "c.c" + strconv.FormatInt(int64(i+1), 10)
	}

	sub := w.SubWriter()
	sub.Write("c.space == ", ctx.SpacePkg.Join(space.Name))
	for _, name := range space.Equivalents {
		sub.Write(" || ", "c.space == ", ctx.SpacePkg.Join(name))
	}

	w.If(sub.Bytes())
	w.ReturnInline()
	w.WriteJoin(args, ", ")
	w.End()

	w.Separate()

	sub.Reset()

	sub.Switch("c.space")

	var cases []string

	var foundPath = false
	for _, src := range ctx.BuiltSpaces {
		eq := ctx.SpaceByName(src.Equivalent)
		if space == src || eq != nil {
			continue
		}

		cases = append(cases[:0], ctx.SpacePkg.Join(src.Name))
		for _, name := range src.Equivalents {
			cases = append(cases, ctx.SpacePkg.Join(name))
		}

		path := ctx.Graph.FindPath(src, space)
		if len(path) == 0 {
			log.Printf("no conversion path found: %q -> %q\n", src.Name, space.Name)
			continue
		}
		foundPath = true
		sub.Case(strings.Join(cases, ", "))

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
	// func (c Color) To(dst space.Space) (Color, bool)
	w.Comment("To converts the color to the destination color space.")
	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, bool")
	w.FuncBody()
	w.If("!c.space.IsValid()")
	w.Return("Color{}, ", false)
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return("c, ", true)
	w.End()

	w.Separate()

	var scope VariableScope

	w.Switch("dst")
	for _, space := range ctx.BuiltSpaces {
		scope.Reset()
		scope.Reserve("c")

		spaceIdent := spacePkg.Join(space.Name)

		w.Case(spaceIdent)
		names := scope.ReserveUniqueAll(space.ChannelIdent()...)
		w.LineWriteJoin(names, ", ")
		w.Writeln(" := c.", space.Name, "()")

		w.ReturnInline("Color{")
		w.Write("space: ", ctx.SpacePkg.Join(space.Name), ", ")
		for i, name := range names {
			w.Write("c", i+1, ": ", name, ", ")
		}
		w.Write("alpha: c.alpha")
		w.Writeln("}, ", true)
	}
	w.Default()
	w.Return("Color{}, ", false)
	w.End()

	w.End()
}
