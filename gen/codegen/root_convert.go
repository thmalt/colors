package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorConvertMethods(ctx *Context, w *writer.GoWriter) {
	var scope VariableScope
	var convertPkg = ctx.ConvertPkg
	var spacePkg = ctx.SpacePkg
	var b strings.Builder
	for _, space := range ctx.BuildSpaces {
		b.Reset()

		scope.Reset()
		scope.Reserve("c")

		names := space.ChannelIdent()

		if scope.ContainsAny(names...) {
			names = scope.ReserveUniqueAll(space.ChannelIdent()...)
		}

		w.Separate()

		w.Comment(space.Name, " returns the color components in the [", ctx.SpacePkg.Join(space.Name), "] color space.")
		w.Method("c Color", space.Name)
		w.FuncResults(joinIdentsWithType(FloatType, names...))
		w.FuncBody()

		w.If("c.space == ", spacePkg.Join(space.Name))
		for i := range names {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "c.c%d", i+1)
		}
		w.Return(b.String())
		w.End()

		w.Separate()

		sub := w.SubWriter()

		sub.Switch("c.space ")

		var foundPath = false
		for _, src := range ctx.BuildSpaces {
			b.Reset()
			if space == src {
				continue
			}

			path := ctx.Graph.FindPath(src, space)
			if len(path) == 0 {
				log.Printf("no conversion path found: %q -> %q\n", src.Name, space.Name)
				continue
			}
			foundPath = true
			sub.Case(spacePkg.Join(src.Name))
			for i := range names {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(&b, "c.c%d", i+1)
			}

			sub.Return(convertPkg.Join(FuncName(src.Name, space.Name)), "(", b.String(), ")")
		}
		sub.Default()
		sub.Return()
		sub.End()

		if foundPath {
			w.Drain(sub)
		} else {
			w.Return()
		}

		w.End()
	}
}
