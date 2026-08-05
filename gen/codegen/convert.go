package codegen

import (
	"log"
	"math"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateConvertPkg(ctx *Context) {
	if ctx.ConvertPkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	pkgPath := filepath.Join(ctx.Directory, ctx.ConvertPkg.Path)
	pkg := ctx.ConvertPkg.Name

	if true {
		// multiple files
		for i, space := range ctx.BuildSpaces {
			if space == nil {
				log.Printf("space at index %d is nil\n", i)
				continue
			}

			if space.Disable {
				continue
			}

			filename := space.Name
			if space.SnakeName != "" {
				filename = space.SnakeName
			}

			emitGoFile(w, pkg, pkgPath, toSnakeCase(filename), func(w *writer.GoWriter) {
				genConvertPkgConversionFile(ctx, w, space)
			})
		}
	} else {
		// single file
		// large file problem
		emitGoFile(w, pkg, pkgPath, pkg, func(w *writer.GoWriter) {
			w.Import("math")

			genConvertPkgConversions(ctx, w)
		})
	}

	emitGoFile(w, pkg, pkgPath, "whitepoint", func(w *writer.GoWriter) {
		genConvertPkgWhitePoint(ctx, w)
	})
}

func genConvertPkgConversions(ctx *Context, w *writer.GoWriter) {
	for i, from := range ctx.BuildSpaces {
		if from == nil {
			log.Printf("space at index %d is nil\n", i)
			continue
		}

		if from.Disable {
			continue
		}

		for j := i + 1; j < len(ctx.BuildSpaces); j++ {
			to := ctx.BuildSpaces[j]
			if to == nil {
				log.Printf("space at index %d is nil\n", j)
				continue
			}

			if to.Disable {
				continue
			}

			if _, ok := ctx.impls[Pair{from.Name, to.Name}]; !ok {
				processPair(ctx, w, from, to)
			}

			if _, ok := ctx.impls[Pair{to.Name, from.Name}]; !ok {
				processPair(ctx, w, to, from)
			}
		}
	}
}

func genConvertPkgConversionFile(ctx *Context, w *writer.GoWriter, space *model.Space) {
	for i, to := range ctx.BuildSpaces {
		if to == nil {
			log.Printf("space at index %d is nil\n", i)
			continue
		}

		if space == to {
			continue
		}

		if to.Disable {
			continue
		}

		if _, ok := ctx.impls[Pair{space.Name, to.Name}]; ok {
			continue
		}

		processPair(ctx, w, space, to)

		w.Separate()
	}
}

func genConvertPkgWhitePoint(ctx *Context, w *writer.GoWriter) bool {
	if len(ctx.WhitePoints) == 0 {
		return false
	}

	for _, whitepoint := range ctx.WhitePoints {
		w.BeginGroup("const ")

		xyz := data.ChromaToXyz(whitepoint.X, whitepoint.Y)
		name := whitepoint.Name
		privateName := toLowerCaseFirstWord(name)

		w.LineWriteln(privateName, "X = ", formatNormalizedFloat(xyz[0]))
		w.LineWriteln(privateName, "Y = ", formatNormalizedFloat(xyz[1]))
		w.LineWriteln(privateName, "Z = ", formatNormalizedFloat(xyz[2]))

		w.Separate()

		w.LineWriteln("inv", name, "X = 1 / ", privateName, "X")
		w.LineWriteln("inv", name, "Y = 1 / ", privateName, "Y")
		w.LineWriteln("inv", name, "Z = 1 / ", privateName, "Z")

		w.End()

		w.Separate()
	}

	return true
}

func processPair(ctx *Context, w *writer.GoWriter, from, to *model.Space) bool {
	path := ctx.Graph.FindPath(from, to)
	if len(path) == 0 {
		log.Printf("Not found conversion path: %q -> %q\n", from.Name, to.Name)
		return false
	}

	ctx.impls[Pair{from.Name, to.Name}] = struct{}{}

	ops := buildOps(ctx, path)
	newOps, changed := combineOps(ops)
	if changed {
		ops = newOps
	}

	paramsVars := from.ChannelSymbols()
	resultsVars := to.ChannelSymbols()

	var scope VariableScope
	scope.ReserveAll(paramsVars...)

	var retString string

	// variable returnResult of params and results
	returnResult := scope.ContainsAny(resultsVars...)
	if returnResult {
		retString = repeatType(FloatType, len(resultsVars))
	} else {
		scope.ReserveAll(resultsVars...)
		retString = joinIdentsWithType(FloatType, resultsVars...)
	}

	funcName := FuncName(from.Name, to.Name)

	w.Comment("Conversion path (", len(path), " steps):")
	w.Comment()
	w.Comment("\t", from.DisplayName)
	for _, node := range path {
		w.Comment("\t-> ", node.To.DisplayName)
	}

	w.Func(funcName)
	w.FuncParams(joinIdentsWithType(FloatType, paramsVars...))
	w.FuncResults(retString)
	w.FuncBody()

	imported := false

	sub := w.SubWriter()
	inputVars := paramsVars

	returned := false
	last := len(ops) - 1

	separate := false

	for idx, op := range ops {
		isLastOp := idx == last

		switch op.Type {
		case OpCall:
			if separate {
				sub.Separate()
			}
			separate = false

			if isLastOp {
				sub.Return(op.Pair.FuncName(), "(", strings.Join(inputVars, ", "), ")")
				returned = true
			} else {
				_, to := ctx.ResolveSpacePair(op.Pair)
				outputVars := to.ChannelSymbols()

				sub.LineWrite(strings.Join(outputVars, ", "))
				if scope.ContainsAll(outputVars...) {
					sub.Write(" = ")
				} else {
					sub.Write(" := ")
				}
				sub.Write(op.Pair.FuncName(), "(")
				sub.Write(strings.Join(inputVars, ", "))
				sub.Writeln(")")

				scope.ReserveAll(outputVars...)
				inputVars = outputVars
			}
		case OpCbrt:
			if idx > 0 {
				sub.Separate()
			}
			separate = true

			if !imported {
				w.Import("math")
				imported = true
			}

			for _, v := range inputVars {

				sub.LineWrite(v, " = math.Cbrt(", v, ")")
			}
		case OpCube:
			if idx > 0 {
				sub.Separate()
			}
			separate = true

			for _, v := range inputVars {
				sub.LineWrite(v, " *= ", v, " * ", v)
			}

		case OpMatrix:
			if idx > 0 {
				sub.Separate()
			}
			separate = true

			outputVars := scope.ReserveUniqueN("f", len(inputVars))
			isNewVar := true
			if isLastOp {
				if !ContainsAny(inputVars, resultsVars) && scope.ContainsAll(resultsVars...) {
					outputVars = resultsVars
					isNewVar = false
				} else {
					resultsVars = outputVars
					returnResult = true
				}
			}

			l := len(inputVars)
			for i, out := range outputVars {
				if isNewVar {
					sub.LineWrite(out, " := ")
				} else {
					sub.LineWrite(out, " = ")
				}

				first := true
				for j, v := range inputVars {
					f := normalizeFloat(op.Matrix[i*l+j])

					switch {
					case f == 0:
					case math.Signbit(f):
						if first {
							sub.Write('-')
						} else {
							sub.Write(" - ")
						}
						f = -f
					case !first:
						sub.Write(" + ")

					}

					switch f {
					case 0:
					case 1:
						sub.Write(v)
						first = false
					default:
						sub.Write(formatFloat(f), '*', v)
						first = false
					}
				}
			}
			inputVars = outputVars
		}
	}

	w.Drain(sub)

	if !returned {
		if last > 0 {
			w.Separate()
		}

		if returnResult {
			w.Return(strings.Join(inputVars, ", "))
		} else {
			w.Return()
		}
	}

	w.End()

	return true
}

func buildOps(ctx *Context, path []*Node) []Op {
	var result []Op
	for _, node := range path {
		if len(node.Fn.Ops) == 0 {
			result = append(result, Op{
				Type: OpCall,
				Pair: node.Fn.Pair,
			})
			continue
		}

		for _, op := range node.Fn.Ops {
			result = append(result, expandOps(ctx, op)...)

		}
	}
	return result
}

func expandOps(ctx *Context, op Op) []Op {
	var out []Op

	switch op.Type {
	case OpCall:
		fn := ctx.ConvertFuncByPair(op.Pair)
		if fn == nil {
			panic("expandOps notfound " + op.Pair.FuncName())
		}

		if len(fn.Ops) == 0 {
			return []Op{op}
		}

		for _, op := range fn.Ops {
			out = append(out, expandOps(ctx, op)...)
		}
	default:
		out = append(out, op)
	}

	return out
}

func combineOps(ops []Op) (out []Op, changed bool) {
	out = make([]Op, 0, len(ops))

	var (
		mat      [9]float64
		hasMat   bool
		lastPair Pair
	)

	flush := func() {
		if !hasMat {
			return
		}

		m := mat
		out = append(out, Op{
			Type:   OpMatrix,
			Pair:   lastPair,
			Matrix: &m,
		})

		hasMat = false
	}

	for _, op := range ops {
		if op.Type == OpMatrix {
			if !hasMat {
				mat = *op.Matrix
				lastPair = op.Pair
				hasMat = true
			} else {
				mat = data.Mat3MulFMA(*op.Matrix, mat)
				changed = true
			}
			continue
		}

		flush()
		out = append(out, op)
	}

	flush()

	return
}

func stringifyPath(path []*Node) string {
	var b strings.Builder

	var last *Node
	for _, node := range path {
		b.WriteString(node.From.Name)
		b.WriteString(" -> ")
		last = node
	}
	b.WriteString(last.To.Name)

	return b.String()
}
