package codegen

import (
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

type conversionStats struct {
	directFuncs map[string]struct{}
	hubFuncs    map[string]struct{}
}

func newConversionStats() conversionStats {
	return conversionStats{
		directFuncs: make(map[string]struct{}),
		hubFuncs:    make(map[string]struct{}),
	}
}

func (s *conversionStats) DirectCounts() int { return len(s.directFuncs) }
func (s *conversionStats) HubCounts() int    { return len(s.hubFuncs) }

func genRootPkgColorConvertMethods(ctx *Context, w *writer.GoWriter, directConversion bool) conversionStats {
	stats := newConversionStats()
	for _, space := range ctx.BuiltSpaces {
		genRootPkgColorConvertMethod(ctx, w, space, stats, directConversion)
	}
	return stats
}

func genRootPkgColorConvertMethod(ctx *Context, w *writer.GoWriter, dst *model.Space, stats conversionStats, directConversion bool) {
	names := dst.ChannelIdent()
	hasNamedReturn := slices.Contains(names, "c")

	w.Separate()
	w.Comment(dst.Name, " returns the color components in the [", ctx.SpacePkg.Join(dst.Name), "] color space.")
	w.Method("c Color", dst.Name)
	if !hasNamedReturn {
		w.FuncResults(joinIdentsWithType(FloatType, names...))
	} else {
		w.FuncResults(joinRepeatN(FloatType, len(names)))
	}
	w.FuncBody()

	if eq := ctx.SpaceByName(dst.Equivalent); eq != nil {
		w.Return("c.", eq.Name, "()")
		w.End()
		return
	}

	var args = make([]string, len(names))
	for i := range args {
		args[i] = "c.c" + strconv.FormatInt(int64(i+1), 10)
	}

	sub := w.SubWriter()
	sub.Write("c.space == ", ctx.SpacePkg.Join(dst.Name))
	for _, name := range dst.Equivalents {
		sub.Write(" || ", "c.space == ", ctx.SpacePkg.Join(name))
	}

	w.If(sub.Bytes())
	w.ReturnInline()
	w.WriteJoin(args, ", ")
	w.End()

	w.Separate()

	sub.Reset()

	sub.Switch("c.space")

	var cases []string

	var foundPath = false

	hubD50 := ctx.SpaceByName(defaultHubD50)
	hubD65 := ctx.SpaceByName(defaultHubD65)

	if hubD50 == nil {
		hubD50 = hubD65
	}
	if hubD65 == nil {
		hubD65 = hubD50
	}
	if hubD50 == nil || hubD65 == nil {
		log.Fatalln("hub not found")
	}

	for _, src := range ctx.BuiltSpaces {
		eq := ctx.SpaceByName(src.Equivalent)
		if dst == src || eq != nil {
			continue
		}

		cases = append(cases[:0], ctx.SpacePkg.Join(src.Name))
		for _, name := range src.Equivalents {
			cases = append(cases, ctx.SpacePkg.Join(name))
		}

		path := ctx.Graph.FindPath(src, dst)
		if len(path) == 0 {
			log.Printf("no conversion path found: %q -> %q\n", src.Name, dst.Name)
			continue
		}
		foundPath = true
		sub.Case(strings.Join(cases, ", "))

		sub.ReturnInline()

		hub := hubD50.Name
		if src.WhitePoint == hubD65.WhitePoint {
			hub = hubD65.Name
		}

		srcOrDstIsHub := src.Name == hubD50.Name || src.Name == hubD65.Name || dst.Name == hubD50.Name || dst.Name == hubD65.Name

		// direct conversion
		if directConversion || src.Family == dst.Family || srcOrDstIsHub {
			pair := Pair{src.Name, dst.Name}
			fnName := pair.FuncName()
			if fn := ctx.ConvertFuncByPair(pair); fn != nil {
				fnName = fn.Pair.FuncName()
			}

			sub.WriteCallln(ctx.ConvertPkg.Join(fnName), args...)

			stats.directFuncs[fnName] = struct{}{}
		} else {
			srcToHub := Pair{src.Name, hub}
			hubToDst := Pair{hub, dst.Name}
			srcToHubFn := srcToHub.FuncName()
			hubToDstFn := hubToDst.FuncName()

			if fn := ctx.ConvertFuncByPair(srcToHub); fn != nil {
				srcToHubFn = fn.Pair.FuncName()
			}

			if fn := ctx.ConvertFuncByPair(hubToDst); fn != nil {
				hubToDstFn = fn.Pair.FuncName()
			}

			sub.Write(ctx.ConvertPkg.Join(hubToDstFn), '(')
			sub.WriteCall(ctx.ConvertPkg.Join(srcToHubFn), args...)
			sub.Writeln(')')

			stats.hubFuncs[srcToHubFn] = struct{}{}
			stats.hubFuncs[hubToDstFn] = struct{}{}
		}
	}

	sub.Default()
	if !hasNamedReturn {
		sub.Return()
	} else {
		sub.Return(joinRepeatN("0", len(names)))
	}
	sub.End()

	if foundPath {
		w.Drain(sub)
	} else {
		if !hasNamedReturn {
			w.Return()
		} else {
			w.Return(joinRepeatN("0", len(names)))
		}
	}

	w.End()
}
