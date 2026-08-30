package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorConstructors(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuiltSpaces {
		genRootPkgSpaceColorConstructor(ctx, w, space, false)
		genRootPkgSpaceColorConstructor(ctx, w, space, true)
	}
}

func genRootPkgSpaceColorConstructor(ctx *Context, w *writer.GoWriter, space *model.Space, withAlpha bool) {
	funcName := space.Name
	params := space.ChannelIdent()
	alpha := "1"

	if withAlpha {
		funcName += "Alpha"
		alpha = "alpha"
		params = append(params, alpha)
	}

	w.Separate()
	w.CommentFunc(func(w *writer.GoWriter) {
		w.Write(funcName, " returns a [Color] from ", space.DisplayName, " components")
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

	w.Func(funcName)
	w.FuncParams(joinIdentsWithType(FloatType, params...))
	w.FuncResults("Color")
	w.FuncBody()

	w.Begin("return Color")
	w.LineWriteln("space: ", ctx.SpacePkg.Join(space.Name), ",")
	for i := range space.ChannelCount() {
		w.LineWriteln("c", i+1, ": ", params[i], ",")
	}
	w.LineWriteln("alpha: ", alpha, ",")
	w.End()

	w.End()
}
