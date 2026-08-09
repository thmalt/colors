package codegen

import (
	"fmt"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColor(ctx *Context, w *writer.GoWriter) {
	maxChannelCnt := 0
	for _, space := range ctx.BuildSpaces {
		maxChannelCnt = max(maxChannelCnt, space.ChannelCount())
	}

	spaceIdent := ctx.SpacePkg.Join("Space")
	w.Commentf("Color represents a color in a [%s].", spaceIdent)
	w.Begin("type Color struct ")
	w.LineWriteln("space ", spaceIdent)
	w.Separate()
	w.Comment("channels")

	for i := range maxChannelCnt {
		w.LineWriteln("c", i+1, " ", FloatType)
	}

	w.Separate()
	w.LineWriteln("alpha ", FloatType)
	w.End()

	w.Separate()

	w.Comment("Channel returns the channel value at index in the order defined by the")
	w.Commentf("color's [%s]. The returned boolean reports whether index is valid.", spaceIdent)
	w.Method("c Color", "Channel")
	w.FuncParams("index int")
	w.FuncResults(FloatType, ", ", "bool")
	w.FuncBody()

	w.If("index < 0 || index >= c.ChannelCount()")
	w.Return("0, false")
	w.End()

	w.Separate()
	w.Switch("index ")

	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i - 1)
		w.Return("c.c", i, ", true")
	}

	w.End()

	w.Separate()
	w.LineWritef("return 0, false\n") // panic(%q)", "Color.Channel unreachable

	w.End()

	var b strings.Builder
	var cases []string

	for i := range maxChannelCnt {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "c.c%d", i+1)
		cases = append(cases, b.String())
	}

	w.Separate()

	w.Commentf("Channels returns the channel values of c in the order defined by its [%s].", spaceIdent)
	w.Method("c Color", "Channels")
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.Switch("c.ChannelCount() ")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("[]", FloatType, "{", cases[i-1], "}")
	}
	w.Default()
	w.Return("nil")
	w.End()

	w.End()

	w.Separate()

	// AppendChannels
	w.Method("c Color", "AppendChannels")
	w.FuncParams("dst []", FloatType)
	w.FuncResults("[]", FloatType)
	w.FuncBody()
	w.Switch("c.ChannelCount() ")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("append(dst, ", cases[i-1], ")")
	}
	w.Default()
	w.Return("dst")
	w.End()

	w.End()
}
