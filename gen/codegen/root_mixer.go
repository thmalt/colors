package codegen

import (
	"strconv"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgHueIndexFunc(ctx *Context, w *writer.GoWriter) {
	w.Func("hueIndex")
	w.FuncParams("s space.Space")
	w.FuncResults("int")
	w.FuncBody()

	hueIndexes := make([][]string, ctx.MaxChannelCount+1)
	for _, s := range ctx.BuildSpaces {
		for i, c := range s.Channels {
			if !c.Circular {
				continue
			}
			hueIndexes[i] = append(hueIndexes[i], ctx.SpacePkg.Join(s.Name))
		}
	}

	w.Switch("s")
	for i, h := range hueIndexes {
		if len(h) > 0 {
			w.Case(strings.Join(h, ", "))
			w.Return(i)
		}
	}
	w.Default()
	w.Return(-1)
	w.End()
	w.End()
	w.Separate()
}

func genRootPkgMixerMethod(ctx *Context, w *writer.GoWriter) {
	channels := make([]bool, ctx.MaxChannelCount+1)

	for _, s := range ctx.BuildSpaces {
		channels[s.ChannelCount()] = true
	}

	w.Method("m Mixer", "Mix")
	w.FuncParams("c1, c2 Color, t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()
	w.Switch("m.channels")

	for n, v := range channels {
		if !v {
			continue
		}

		w.Case(n)
		genMixerCase(w, n)
	}
	w.Default()
	w.Return("Color{}")
	w.End()
	w.End()
	w.Separate()
}

func genMixerCase(w *writer.GoWriter, n int) {
	var args []string

	a := genChannels(nil, "a", n, true)
	b := genChannels(nil, "b", n, true)
	c := genChannels(nil, "c", n, true)

	genMixerInputs(w, "a", "c1", n)
	genMixerInputs(w, "b", "c2", n)

	w.LineWriteJoin(c, ", ")
	w.Write(" := m.unsafe.Mix", n, '(')

	args = append(args, a...)
	args = append(args, b...)
	args = append(args, "t")

	w.WriteJoin(args, ", ")
	w.Writeln(')')

	genMixerResult(w, n)
}

func genMixerInputs(w *writer.GoWriter, prefix, color string, n int) {
	args := genChannels(nil, prefix, n, true)
	channels := genChannels(nil, color+".c", n, false)

	w.LineWriteJoin(args, ", ")
	w.Write(" := ")
	w.WriteJoin(channels, ", ")
	w.Write(", ", color, ".alpha")
}

func genMixerResult(w *writer.GoWriter, n int) {
	w.ReturnInline("Color{space: m.space,")

	for i := 1; i <= n; i++ {
		w.Write("c", i, ": c", i, ", ")
	}

	w.Writeln("alpha: ca}")
}

func genChannels(dst []string, prefix string, num int, a bool) []string {
	dst = dst[:0]
	for i := range num {
		dst = append(dst, prefix+strconv.Itoa(i+1))
	}
	if a {
		dst = append(dst, prefix+"a")
	}
	return dst
}
