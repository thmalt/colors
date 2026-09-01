package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genRootPkgColorTo(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg
	w.Separate()
	// func (c Color) to(dst space.Space) (Color, bool)
	w.Comment("to converts the color to the destination color space.")
	w.Method("c Color", "to")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, bool")
	w.FuncBody()
	w.If("!c.space.IsValid()")
	w.Return("Color{}, ", false)
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return("c, ", true)
	w.End()

	w.Separate()

	var scope VariableScope

	w.Switch("dst")
	for _, space := range ctx.BuiltSpaces {
		scope.Reset()
		scope.Reserve("c")

		spaceIdent := spacePkg.Join(space.Name)

		w.Case(spaceIdent)
		names := scope.ReserveUniqueAll(space.ChannelIdent()...)
		w.LineWriteJoin(names, ", ")
		w.Writeln(" := c.", space.Name, "()")

		w.ReturnInline("Color{")
		w.Write("space: ", ctx.SpacePkg.Join(space.Name), ", ")
		for i, name := range names {
			w.Write("c", i+1, ": ", name, ", ")
		}
		w.Write("alpha: c.alpha")
		w.Writeln("}, ", true)
	}
	w.Default()
	w.Return("Color{}, ", false)
	w.End()

	w.End()
}
