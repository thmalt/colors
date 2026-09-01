package codegen

import (
	"fmt"
	"math"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColor(ctx *Context, w *writer.GoWriter) {
	maxChannelCnt := ctx.MaxChannelCount

	// type Color struct{}
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
	// func (c Color) Channel(index int) (float64, bool)
	w.Comment("Channel returns the channel value at index in the order defined by the")
	w.Commentf("color's [%s]. The returned boolean reports whether index is valid.", spaceIdent)
	w.Method("c Color", "Channel")
	w.FuncParams("index int")
	w.FuncResults(FloatType, ", ", "bool")
	w.FuncBody()

	w.If("index >= c.ChannelCount()")
	w.Return("0, false")
	w.End()

	w.Separate()
	w.Switch("index")

	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i - 1)
		w.Return("c.c", i, ", true")
	}

	w.End()

	w.Separate()
	w.LineWritef("return 0, false\n")

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
	// func (c Color) Channels() []float64
	w.Commentf("Channels returns the channel values of c in the order defined by its [%s].\n", spaceIdent)
	w.Method("c Color", "Channels")
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.Switch("c.ChannelCount()")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("[]", FloatType, "{", cases[i-1], "}")
	}
	w.Default()
	w.Return("nil")
	w.End()

	w.End()

	w.Separate()
	// func (c Color) AppendChannels(dst []float64) []float64
	w.Comment("AppendChannels appends the color channels to dst and returns the resulting slice.")
	w.Method("c Color", "AppendChannels")
	w.FuncParams("dst []", FloatType)
	w.FuncResults("[]", FloatType)
	w.FuncBody()
	w.Switch("c.ChannelCount()")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("append(dst, ", cases[i-1], ")")
	}
	w.Default()
	w.Return("dst")
	w.End()

	w.End()
}

func genRootPkgColorChannel(ctx *Context, w *writer.GoWriter) {
	channels := buildChannelCounts(ctx)

	vars := appendVars(nil, "c", ctx.MaxChannelCount)
	fields := appendVars(nil, "c.c", ctx.MaxChannelCount)
	for count, ok := range channels {
		if !ok {
			continue
		}
		w.Separate()
		// func (c Color) Channel?() (c1, c2, ... float64)
		w.Commentf("Channel%d returns the first %d channels.\n", count, count)
		w.Method("c Color", "Channel", count)
		w.FuncResults(joinIdentsWithType(FloatType, vars[:count]...))
		w.FuncBody()
		w.ReturnInline()
		w.Writeln(strings.Join(fields[:count], ", "))
		w.End()
	}
}

func genRootPkgHexLUT(_ *Context, w *writer.GoWriter) {
	count := math.MaxUint8 + 1
	w.Separate()
	w.Begin("var hexLUT = [", count, "]uint8 ")

	next := wrapEvery(w, 8)
	invalid := true
	for i := range count {
		c := uint8(i)
		if i >= '0' && i <= '9' {
			if invalid {
				w.Separate()
				w.Comment("0 - 9")
			}
			invalid = false
			v := c - '0'
			w.Write(fmt.Sprintf("0x%02x,", v))
		} else if i >= 'A' && i <= 'F' {
			if invalid {
				w.Separate()
				w.Comment("A - F")
			}
			invalid = false

			v := c - 'A' + 10
			w.Write(fmt.Sprintf("0x%02x,", v))
		} else if i >= 'a' && i <= 'f' {
			if invalid {
				w.Separate()
				w.Comment("a - f")
			}
			invalid = false

			v := c - 'a' + 10
			w.Write(fmt.Sprintf("0x%02x,", v))
		} else {
			if !invalid {
				w.Separate()
			}
			invalid = true

			w.Write("maxUint8,")
			next()
		}
	}
	w.End()
}
