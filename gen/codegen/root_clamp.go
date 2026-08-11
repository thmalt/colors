package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgClamp(ctx *Context, w *writer.GoWriter) {
	w.Func("Clamp")
	w.FuncParams("c Color")
	w.FuncResults("Color")
	w.FuncBody()
	w.If("!c.space.IsValid()")
	w.Return("c")
	w.End()

	w.Separate()
	w.Switch("c.space")
	for _, space := range ctx.BuildSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))
		w.Return("clamp", space.Name, "(c)")
	}
	w.Default()
	w.Return("c")
	w.End()
	w.End()

	w.Separate()
	for _, space := range ctx.BuildSpaces {
		genRootPkgSpaceClamp(w, space)
		w.Separate()
	}
}

func genRootPkgSpaceClamp(w *writer.GoWriter, space *model.Space) {
	w.Func("clamp", space.Name)
	w.FuncParams("c Color")
	w.FuncResults("Color")
	w.FuncBody()
	for i, c := range space.Channels {
		if c.Unrestricted {
			continue
		}

		var fn = "clamp"
		if c.Circular {
			fn = "wrap"
		}

		w.LineWritef("c.c%d = %s(c.c%d, %s, %s)\n",
			i+1,
			fn,
			i+1,
			formatNormalizedFloat(c.Min),
			formatNormalizedFloat(c.Max),
		)
	}
	w.Return("c")
	w.End()
}
