package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgClamp(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	// func Clamp(c Color) Color
	w.Comment("Clamp clamps the color channels to the valid range of the color space.")
	w.Func("Clamp")
	w.FuncParams("c Color")
	w.FuncResults("Color")
	w.FuncBody()

	w.Switch("c.space")

	sub := w.SubWriter()
	groups := rootPkgClampGroup(ctx, sub)

	for _, group := range groups {
		sub.Reset()

		for i, s := range group.Spaces {
			if i > 0 {
				if i%4 == 0 {
					sub.Write(",")
					sub.Indent()
				} else {
					sub.Write(", ")
				}
			}
			sub.Write(ctx.SpacePkg.Join(s.Name))
		}

		w.Case(sub.Bytes())
		w.Write(group.Key)
		w.Return("c")
	}

	w.Default()
	w.Return("c")
	w.End()
	w.End()
}

func rootPkgClampGroup(ctx *Context, w *writer.GoWriter) []groupSpaceValue {
	defer w.Reset()

	gs := newGroupSpace()
	for i, space := range ctx.BuiltSpaces {
		w.Reset()

		count := 0
		for j, ch := range space.Channels {
			if ch.Unrestricted {
				continue
			}

			w.LineWrite("c.c", j+1, " = ")

			min := normalizeFloat(ch.Min)
			max := normalizeFloat(ch.Max)

			if ch.Circular && min == 0 && max == 360 {
				w.Writeln("wrap360", '(', "c.c", j+1, ')')
			} else {
				fn := "clamp"
				if ch.Circular {
					fn = "wrap"
				}

				w.Writeln(
					fn, '(',
					"c.c", j+1,
					", ", formatFloat(min),
					", ", formatFloat(max),
					')',
				)
			}

			count++
		}
		gs.Append(string(w.Bytes()), count, i, space)
	}

	return gs.SortedSlice()
}
