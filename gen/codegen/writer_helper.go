package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func writeColorChannels(w *writer.GoWriter, name string, channelCount int) {
	for i := range channelCount {
		if i > 0 {
			w.Write(", ")
		}
		w.Write(name, ".c", i+1)
	}
	w.Write(", ", name, ".alpha")
}

func writeColorLiteral(w *writer.GoWriter, space string, channelCount int, vars ...string) {
	if len(vars) < channelCount {
		panic("writeColorLiteral: len(vars) < channelCount")
	}

	w.Write("Color{space: ", space, ", ")
	for i := range channelCount {
		w.Write("c", i+1, ": ", vars[i], ", ")
	}
	alpha := "1"
	if len(vars) > channelCount {
		alpha = vars[len(vars)-1]
	}
	w.Write("alpha: ", alpha, "}")
}

func writeColorLiteralMultiLine(w *writer.GoWriter, space string, channelCount int, vars ...string) {
	if len(vars) < channelCount {
		panic("writeColorLiteralMultiLine: len(vars) < channelCount")
	}

	w.Writeln("Color{")
	w.In()
	w.LineWriteln("space: ", space, ",")
	for i := range channelCount {
		w.LineWriteln("c", i+1, ":    ", vars[i], ",")
	}
	alpha := "1"
	if len(vars) > channelCount {
		alpha = vars[len(vars)-1]
	}
	w.LineWriteln("alpha: ", alpha, ",")
	w.Out()
	w.LineWriteln('}')
}
