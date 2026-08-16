package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgClamp(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	w.Func("Clamp")
	w.FuncParams("c Color")
	w.FuncResults("Color")
	w.FuncBody()

	w.If("!c.space.IsValid()")
	w.Return("c")
	w.End()

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
	for i, space := range ctx.BuildSpaces {
		w.Reset()

		count := 0
		for j, ch := range space.Channels {
			if ch.Unrestricted {
				continue
			}

			w.LineWrite("c.c", j+1, " = ")
			if ch.Circular {
				w.Write("wrap")
			} else {
				w.Write("clamp")
			}
			w.Writeln('(',
				"c.c", j+1,
				", ", formatNormalizedFloat(ch.Min),
				", ", formatNormalizedFloat(ch.Max),
				')',
			)

			count++
		}
		gs.Append(string(w.Bytes()), count, i, space)
	}

	return gs.SortedSlice()
}
