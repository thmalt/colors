package codegen

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateInterpPkg(ctx *Context) {
	if ctx.InterpPkg.Name == "" {
		return
	}

	var pkgPath = filepath.Join(ctx.Directory, ctx.InterpPkg.Path)

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	for i, space := range ctx.BuildSpaces {
		if space == nil {
			log.Printf("GenerateInterpPkg: space at index %d is nil\n", i)
			continue
		}

		if space.Disable {
			continue
		}

		genInterp(w, space, false)
		genInterp(w, space, true)
	}

	fileName := ctx.InterpPkg.Name + "_gen.go"
	fmt.Println("  Generate file", fileName)

	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.InterpPkg.Name,
	)

}

func genInterp(w *writer.GoWriter, space *model.Space, withAlpha bool) {
	ident := space.ChannelSymbols()
	p1 := make([]string, 0, len(ident)+1)
	p2 := make([]string, 0, len(ident)+1)
	for _, param := range ident {
		p1 = append(p1, param+"1")
		p2 = append(p2, param+"2")
	}

	if withAlpha {
		p1 = append(p1, "alpha1")
		p2 = append(p2, "alpha2")
		ident = append(ident, "alpha")
	}

	if withAlpha {
		w.Func(space.Name, "Alpha")
	} else {
		w.Func(space.Name)
	}

	if space.Coordinate == model.Polar {
		w.FuncParams(varJoinWithType(p1...), ", ", varJoinWithType(p2...), ", t float64", ", hue HueInterpolation")
	} else {
		w.FuncParams(varJoinWithType(p1...), ", ", varJoinWithType(p2...), ", t float64")
	}

	w.FuncResults(varJoinWithType(ident...))
	w.FuncBody()

	if withAlpha {
		l := len(ident) - 1
		w.LineWritef("%s = lerp(%s, %s, t)\n", ident[l], p1[l], p2[l])
	}

	for i, c := range space.Channels {
		if c.Circular {
			if i > 0 {
				w.Separate()
			}

			w.Switch("hue")
			for _, hue := range HueInterpolation {
				w.Case(hue)
				w.LineWritef("%s = lerp%s(%s, %s, t)", ident[i], hue, p1[i], p2[i])
			}
			w.End()

			w.Separate()

		} else {
			w.LineWritef("%s = lerp(%s, %s, t)\n", ident[i], p1[i], p2[i])
		}
	}

	w.Return()

	w.End()

	w.Separate()
}
