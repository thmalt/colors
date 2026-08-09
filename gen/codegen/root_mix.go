package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMix(ctx *Context, w *writer.GoWriter) {
	spacePkg := ctx.SpacePkg

	w.Comment("MixWith returns the [Color] interpolation of c1 and c2.")
	w.Comment()
	w.Comment("The interpolation is performed in opts.Space. If opts.Space is")
	w.Commentf("[%s], [%s] is used.", spacePkg.Join("InvalidSpace"), spacePkg.Join("Oklab"))
	w.Comment()
	w.Comment("The interpolation behavior can be customized through opts, including")
	w.Comment("premultiplied alpha and hue interpolation for polar color spaces.")
	w.Func("MixWith")
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
	w.LineWrite("unsafeMixer := ")
	w.WriteCallln(ctx.MixerPkg.Join("NewUnsafeMixer"), "opts.Space.HueIndex()", "!opts.Unpremultiplied", "opts.Hue")

	w.Separate()
	w.Switch("opts.Space.ChannelCount()")

	channelCounts := buildChannelCounts(ctx)

	for channelCount, ok := range channelCounts {
		if !ok {
			continue
		}

		w.Case(channelCount)
		genRootPkgMixerMethodCase(w, "unsafeMixer.Mix", "opts.Space", channelCount)
	}
	w.End()

	w.Separate()
	w.Return("Color{}")
	w.End()
}
