package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genSpacePkgTables(ctx *Context, w *writer.GoWriter) {
	w.Begin("var spaceInfos = [...]*spaceInfo")
	w.LineWriteln("nil,")
	for _, space := range ctx.BuildSpaces {
		w.LineWriteln("&", spaceInfoName(space), ",")
	}
	w.End()

	next := wrapEvery(w, 8)
	w.Begin("var spaceChannelCounts = [...]uint", smallestUintType(ctx.MaxChannelCount))
	w.LineWrite("0, ")
	next()
	for _, space := range ctx.BuildSpaces {
		w.Write(space.ChannelCount(), ", ")
		next()
	}
	w.End()

	w.Separate()
	next = wrapEvery(w, 8)
	w.Begin("var coordinateSystems = [...]CoordinateSystem")
	w.LineWrite(model.CoordinateSystem(0), ", ")
	next()
	for _, space := range ctx.BuildSpaces {
		w.Write(space.Coordinate, ", ")
		next()
	}
	w.End()

}

func wrapEvery(w *writer.GoWriter, n int) func() {
	count := 0
	return func() {
		count++
		if count%n == 0 {
			w.Indent()
		}
	}
}
