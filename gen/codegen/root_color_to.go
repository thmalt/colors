package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genRootPkgColorTo(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg

	w.Separate()

	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()
	w.If("!c.space.IsValid() || !dst.IsValid()")
	w.Return("Color{}, ErrInvalidSpace")
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return("c, nil")
	w.End()

	var scope VariableScope

	w.Separate()

	w.Switch("dst")
	for _, space := range ctx.BuildSpaces {
		scope.Reset()
		scope.Reserve("c")

		spaceIdent := spacePkg.Join(space.Name)

		w.Case(spaceIdent)
		names := scope.ReserveUniqueAll(space.ChannelIdent()...)
		w.LineWriteJoin(names, ", ")
		w.Writef(" := c.%s()\n", space.Name)

		w.ReturnInline("Color{")
		w.Write("space: ", ctx.SpacePkg.Join(space.Name), ",")
		for i, name := range names {
			w.Write("c", i+1, ": ", name, ",")
		}
		w.Write("alpha: c.alpha")
		w.Writeln("}, nil")
	}
	w.Default()
	w.Return("Color{}, ErrInvalidSpace")
	w.End()

	w.End()
}
