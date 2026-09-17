package codegen

import (
	"strconv"

	"github.com/thmalt/colors/gen/codegen/internal/interp"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genMixerPkgRawMixer(ctx *Context, w *writer.GoWriter) {
	indexes := make([]bool, ctx.MaxChannelCount)
	for i := range indexes {
		indexes[i] = true
	}

	channels := buildChannelCounts(ctx)
	for count, ok := range channels {
		if !ok {
			continue
		}

		genMixerPkgRawMixerMethod(ctx, w, count, indexes[:count])
	}
}

func genMixerPkgRawMixerMethod(ctx *Context, w *writer.GoWriter, channelCount int, indexes []bool) {
	a := toVars(nil, "a", channelCount, "a")
	b := toVars(nil, "b", channelCount, "a")
	c := toVars(nil, "c", channelCount, "a")

	w.Separate()
	w.Comment("Mix", channelCount, " interpolates ", channelCount, " color channels and alpha between two colors.")
	w.Comment("The channel values must match the mixer's channel configuration.")
	w.Method("m RawMixer", "Mix", channelCount)
	w.FuncParams(
		joinIdentsWithType(FloatType, a...),
		", ",
		joinIdentsWithType(FloatType, b...),
		", t ",
		FloatType,
	)
	w.FuncResults(joinIdentsWithType(FloatType, c...))
	w.FuncBody()

	w.Switch("m.hueIndex")

	for index, isHueIndex := range indexes {
		if index >= channelCount {
			break
		}

		if !isHueIndex {
			continue
		}

		w.Case(index)

		genMixerPkgRawMixMethodHue(ctx, w, a[index], b[index], c[index])

		w.Separate()
		genMixerPkgRawMixMethodLinear(w, a, b, c, channelCount, index)

		w.Return()
	}
	w.Default()

	genMixerPkgRawMixMethodLinear(w, a, b, c, channelCount, -1)

	w.Return()

	w.End()

	w.End()
}

func genMixerPkgRawMixMethodHue(ctx *Context, w *writer.GoWriter, a, b, c string) {
	lerpHue := ctx.InterpPkg.Join("LerpHue")
	w.Switch("m.hueInterp")
	for i := range interp.HueDecreasing {
		name := (i + 1).String()

		w.Case(ctx.InterpPkg.Join("Hue" + name))
		w.LineWrite(c, " = ")
		w.WriteCall(lerpHue+name, a, b, "t")
	}
	w.Default()
	w.LineWrite(c, " = ")
	w.WriteCall(lerpHue+interp.HueShorter.String(), a, b, "t")
	w.End()
}

func genMixerPkgRawMixMethodLinear(w *writer.GoWriter, a, b, c []string, count, hueIndex int) {
	w.If("m.premultiplied")

	w.LineWriteln("w1, w2 := ", a[count], "*(1-t), ", b[count], "*t") // w1, w2 := a*(1-t), b*t
	w.LineWriteln(c[count], " = w1 + w2")
	for i := range count {
		if i == hueIndex {
			continue
		}
		w.LineWriteln(c[i], " = ", a[i], "*w1 + ", b[i], "*w2") // c = a*w1 + b*w2
	}

	w.Separate()
	w.If(c[count], " != 0")
	w.LineWriteln("inv := 1 / ", c[count])
	for i := range count {
		if i == hueIndex {
			continue
		}
		w.LineWriteln(c[i], " *= inv")
	}
	w.End()

	w.Else()
	w.LineWriteln(c[count], " = ", a[count], " + (", b[count], "-", a[count], ")*t") // c = a + (b-a)*t
	for i := range count {
		if i == hueIndex {
			continue
		}

		w.LineWriteln(c[i], " = ", a[i], " + (", b[i], "-", a[i], ")*t")
	}

	w.End()
}

func toVars(dst []string, prefix string, count int, a ...string) []string {
	for i := range count {
		dst = append(dst, prefix+strconv.Itoa(i+1))
	}

	for _, c := range a {
		dst = append(dst, prefix+c)
	}

	return dst
}
