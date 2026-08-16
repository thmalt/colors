package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorConstructors(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuildSpaces {
		genRootPkgSpaceColorConstructor(ctx, w, space, false)
		genRootPkgSpaceColorConstructor(ctx, w, space, true)
	}
}

func genRootPkgSpaceColorConstructor(ctx *Context, w *writer.GoWriter, space *model.Space, withAlpha bool) {
	w.Separate()
	w.CommentFunc(func(w *writer.GoWriter) {
		w.Write(space.Name, " returns a [Color] from ", space.DisplayName, " components")
		if withAlpha {
			w.Write(" with alpha")
		}
		w.Writeln(".")

		w.Separate()

		for _, c := range space.Channels {
			w.LineWrite(
				'\t',
				c.Ident,
				": [",
				formatNormalizedFloat(c.Min),
				", ",
				formatNormalizedFloat(c.Max),
			)

			if c.Circular {
				w.Write(')')
			} else {
				w.Write(']')
			}

			if c.Unrestricted {
				w.Write(" (typical)")
			}
		}
		if withAlpha {
			w.LineWriteln("\talpha: [0, 1]")
		}
	})

	params := space.ChannelIdent()
	w.Func(space.Name)
	if withAlpha {
		w.Write("Alpha")
		params = append(params, "alpha")
	}
	w.FuncParams(joinIdentsWithType(FloatType, params...))
	w.FuncResults("Color")
	w.FuncBody()

	w.Begin("return Color")
	w.LineWriteln("space: ", ctx.SpacePkg.Join(space.Name), ",")
	for i := range space.ChannelCount() {
		w.LineWriteln("c", i+1, ": ", params[i], ",")
	}
	if withAlpha {
		w.LineWriteln("alpha: ", params[len(params)-1], ",")
	} else {
		w.LineWriteln("alpha: 1,")
	}
	w.End()

	w.End()
}
