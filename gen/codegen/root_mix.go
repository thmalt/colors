package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMix(ctx *Context, w *writer.GoWriter) {
	spacePkg := ctx.SpacePkg

	w.Comment("Mix returns the [Color] interpolation of c1 and c2.")
	w.Comment()
	w.Comment("The interpolation is performed in opts.Space. If opts.Space is")
	w.Commentf("[%s], [%s] is used.", spacePkg.Join("InvalidSpace"), spacePkg.Join("Oklab"))
	w.Comment()
	w.Comment("The interpolation behavior can be customized through opts, including")
	w.Comment("premultiplied alpha and hue interpolation for polar color spaces.")
	w.Func("Mix")
	w.FuncParams("c1, c2 Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	w.If("opts.Space == ", spacePkg.Join("InvalidSpace"))
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
		genRootPkgSpaceMix(ctx, w, space)
	}
}

func genRootPkgSpaceMix(ctx *Context, w *writer.GoWriter, space *model.Space) {
	w.Func("mix", space.Name)
	w.FuncParams("c1, c2 Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	w.LineWriteln("a1, a2 := c1.alpha, c2.alpha")
	w.LineWriteln("alpha := lerp(a1, a2, t)")

	w.Separate()

	idents := space.ChannelIdent()

	w.LineWrite("var ")
	sep := ""
	for i, c := range space.Channels {
		if c.Circular {
			continue
		}

		w.Write(sep)
		w.Write(idents[i])
		sep = ", "
	}
	w.Write(" ", FloatType)

	w.Separate()
	w.If("!opts.Unpremultiplied")
	w.LineWriteln("w1, w2 := a1*(1-t), a2*t")

	w.Separate()
	for i, c := range space.Channels {
		if c.Circular {
			continue
		}
		index := i + 1
		w.LineWritef("%s = c1.c%d*w1 + c2.c%d*w2\n", idents[i], index, index)
	}

	w.Separate()
	w.If("alpha != 0")
	w.LineWriteln("inv := 1 / alpha")
	for i, c := range space.Channels {
		if c.Circular {
			continue
		}
		w.LineWriteln(idents[i], " *= inv")
	}
	w.End()

	w.Else()
	for i, c := range space.Channels {
		if c.Circular {
			continue
		}
		index := i + 1
		w.LineWritef("%s = lerp(c1.c%d, c2.c%d, t)\n", idents[i], index, index)
	}
	w.End()

	for i, c := range space.Channels {
		if !c.Circular {
			continue
		}

		if i > 0 {
			w.Separate()
		}

		genRootPkgMixHueInterpolation(ctx, w, i+1, idents[i])
	}

	w.Separate()
	w.Begin("return Color")
	w.LineWriteln("space: ", ctx.SpacePkg.Join(space.Name), ",")
	for i, ident := range idents {
		w.LineWritef("c%d: %s,\n", i+1, ident)
	}
	w.LineWriteln("alpha: alpha,")
	w.End()

	w.End()
}

func genRootPkgMixHueInterpolation(ctx *Context, w *writer.GoWriter, index int, ident string) {
	switch ctx.Optimization {
	case OptimizeSize:
		w.LineWritef("%s := lerpHue(c1.c%d, c2.c%d, t, opts.Hue)\n", ident, index, index)
	default:
		w.LineWrite("var ", ident, " ", FloatType)
		w.Switch("opts.Hue")
		for _, hue := range HueInterpolation {
			w.Case(hue)
			w.LineWritef("%s = lerp%s(c1.c%d, c2.c%d, t)\n", ident, hue, index, index)
		}
		w.Default()
		w.LineWritef("%s = lerpHueShorter(c1.c%d, c2.c%d, t)\n", ident, index, index)
		w.End()
	}
}
