package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genRootPkgColorMutTo(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg
	w.Separate()
	// func (c *Color) mutTo(dst space.Space) bool
	w.Comment("mutTo converts the color to the specified color space in place")
	w.Comment("and reports whether the conversion succeeded.")
	w.Method("c *Color", "mutTo")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("bool")
	w.FuncBody()
	w.If("!c.space.IsValid()")
	w.Return(false)
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return(true)
	w.End()

	w.Separate()

	w.Switch("dst")
	temp := make([]string, ctx.MaxChannelCount)
	for _, space := range ctx.BuiltSpaces {
		spaceIdent := spacePkg.Join(space.Name)
		channelCount := space.ChannelCount()

		w.Case(spaceIdent)

		w.LineWriteJoin(appendVars(temp[:0], "c.c", channelCount), ", ")

		w.Writeln(" = c.", space.Name, "()")
		for i := channelCount; i < ctx.MaxChannelCount; i++ {
			w.LineWriteln("c.c", i+1, " = 0")
		}
		w.LineWriteln("c.space = ", spaceIdent)
		w.Return(true)
	}
	w.Default()
	w.Return(false)
	w.End()

	w.End()
}
