package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgClamp(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	// func Clamp(c Color) Color
	w.Comment("Clamp clamps the color channels to their valid ranges and alpha to [0, 1].")
	w.Func("Clamp")
	w.FuncParams("c Color")
	w.FuncResults("Color")
	w.FuncBody()

	w.LineWriteln("c.alpha = clamp01(c.alpha)")

	w.Separate()
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

			if ch.Circular {
				if min == 0 && max == 360 {
					w.Writeln("wrap360", '(', "c.c", j+1, ')')
				} else {
					w.Writeln(
						"wrap", '(',
						"c.c", j+1,
						", ", formatFloat(min),
						", ", formatFloat(max),
						')',
					)
				}
			} else if min == 0 && max == 1 {
				w.Writeln("clamp01", '(', "c.c", j+1, ')')
			} else {
				w.Writeln(
					"clamp", '(',
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
