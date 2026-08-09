package codegen

import (
	"log"
	"slices"
	"strconv"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorConvertMethods(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuildSpaces {
		genRootPkgColorConvertMethod(ctx, w, space)
		w.Separate()
	}
}

func genRootPkgColorConvertMethod(ctx *Context, w *writer.GoWriter, space *model.Space) {
	names := space.ChannelIdent()
	hasNamedReturn := slices.Contains(names, "c")

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
	for _, src := range ctx.BuildSpaces {
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
