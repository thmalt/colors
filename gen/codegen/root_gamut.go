package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgGamut(ctx *Context, w *writer.GoWriter) {
	genRootPkgInGamut(ctx, w)
	genRootPkgInGamutSpace(w)
}

func genRootPkgInGamut(ctx *Context, w *writer.GoWriter) {
	w.Separate()
	w.Comment("InGamut reports whether c is within the gamut of its color space.")
	w.Func("InGamut")
	w.FuncParams("c Color")
	w.FuncResults("bool")
	w.FuncBody()

	w.Switch("c.space")

	sub := w.SubWriter()
	groups := rootPkgInGamutGroup(ctx, sub)

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
		w.Return(group.Key)
	}

	w.Default()
	w.Return("false")
	w.End()
	w.End()
}

func genRootPkgInGamutSpace(w *writer.GoWriter) {
	w.Separate()
	w.Comment("InGamutSpace reports whether c is within the gamut of the specified color space.")
	w.Func("InGamutSpace")
	w.FuncParams("c Color, dst space.Space")
	w.FuncResults("bool")
	w.FuncBody()

	w.If("c.space == dst")
	w.Return("InGamut(c)")
	w.End()

	w.Separate()
	w.LineWriteln("converted, err := c.To(dst)")
	w.If("err != nil")
	w.Return("false")
	w.End()

	w.Separate()
	w.Return("InGamut(converted)")
	w.End()

}

func rootPkgInGamutGroup(ctx *Context, w *writer.GoWriter) []groupSpaceValue {
	defer w.Reset()

	gs := newGroupSpace()
	for i, space := range ctx.BuiltSpaces {
		w.Reset()

		count := 0
		for j, ch := range space.Channels {
			if ch.Unrestricted || ch.Circular {
				continue
			}

			if count > 0 {
				w.Write(" && ")
			}

			w.Write(
				"c.c", j+1, " >= ", formatNormalizedFloat(ch.Min), " && ",
				"c.c", j+1, " <= ", formatNormalizedFloat(ch.Max),
			)

			count++
		}

		if count == 0 {
			w.Write("true")
		}
		gs.Append(string(w.Bytes()), count, i, space)
	}

	return gs.SortedSlice()
}
