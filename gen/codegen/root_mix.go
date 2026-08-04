package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMix(ctx *Context, w *writer.GoWriter) {
	w.Comment("Mix returns the [Color] interpolation of c1 and c2.")
	w.Comment()
	w.Comment("The interpolation is performed in opts.Space. If opts.Space is")
	w.Comment("[space.SpaceInvalid], [space.Oklab] is used.")
	w.Comment()
	w.Comment("The interpolation behavior can be customized through opts, including")
	w.Comment("premultiplied alpha and hue interpolation for polar color spaces.")
	w.Func("Mix")
	w.FuncParams("c1, c2 Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	w.If("opts.Space == space.SpaceInvalid")
	w.LineWriteln("opts.Space = space.Oklab")
	w.End()

	w.Separate()
	w.If("c1.space.IsValid() && c2.space.IsValid() && !opts.Space.IsValid()")
	w.Return("Color{}")
	w.End()

	w.Separate()
	w.LineWriteln("c1 = c1.MustTo(opts.Space)")
	w.LineWriteln("c2 = c2.MustTo(opts.Space)")

	w.Separate()
	w.Switch("opts.Space ")
	for _, space := range ctx.BuildSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))
		w.Return("mix", space.Name, "(c1, c2, t, opts)")
	}
	w.End()

	w.Separate()
	w.LineWriteln(PanicUnreachable)
	w.End()

	for _, space := range ctx.BuildSpaces {
		w.Separate()
		genRootPkgMixSpaceFunc(ctx, w, space)
	}
}

func genRootPkgMixSpaceFunc(ctx *Context, w *writer.GoWriter, space *model.Space) {
	w.Func("mix", space.Name)
	w.FuncParams("c1, c2 Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	syms := space.ChannelSymbols()

	w.LineWriteln("var ", varJoinWithType(syms...))

	w.Separate()
	for i, c := range space.Channels {
		if c.Circular {
			if i > 0 {
				w.Separate()
			}
			w.Switch("opts.Hue ")
			for _, hue := range HueInterpolation {
				w.Case(hue)
				w.LineWritef("%s = lerpHueShorter(c1.c%d, c2.c%d, t)\n", syms[i], i+1, i+1)
			}
			w.End()

			w.Separate()
		}
	}

	w.LineWriteln("a1, a2 := c1.alpha, c2.alpha")
	w.LineWriteln("alpha := lerp(a1, a2, t)")

	w.Separate()
	w.If("opts.Premultiplied")
	w.LineWriteln("w1, w2 := a1*(1-t), a2*t")

	w.Separate()
	for i, c := range space.Channels {
		if c.Circular {
			continue
		}

		w.LineWritef("%s = c1.c%d*w1 + c2.c%d*w2\n", syms[i], i+1, i+1)
	}
	w.Separate()
	w.If("alpha != 0")
	w.LineWriteln("inv := 1 / alpha")
	for i, c := range space.Channels {
		if !c.Circular {
			w.LineWriteln(syms[i], " *= inv")
		}
	}
	w.End()

	w.Else()
	for i, c := range space.Channels {
		if !c.Circular {
			w.LineWritef("%s = lerp(c1.c%d, c2.c%d, t)\n", syms[i], i+1, i+1)
		}
	}
	w.End()

	w.Separate()
	w.Begin("return Color")
	w.LineWriteln("space: ", ctx.SpacePkg.Join(space.Name), ",")
	w.Indent()
	for i, c := range syms {
		if i > 0 {
			w.Write(", ")
		}
		w.Writef("c%d: %s", i+1, c)
	}
	w.Write(",")
	w.LineWriteln("alpha: alpha,")
	w.End()

	w.End()
}
