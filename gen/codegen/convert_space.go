package codegen

import (
	"log"
	"math"
	"strings"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genConvertPkgSpaceConversions(ctx *Context, w *writer.GoWriter, space *model.Space) int {
	count := 0
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

		genConvertPkgSpacePair(ctx, w, space, to)

		w.Separate()

		count++
	}

	return count
}

func genConvertPkgSpacePair(ctx *Context, w *writer.GoWriter, from, to *model.Space) bool {
	path := ctx.Graph.FindPath(from, to)
	if len(path) == 0 {
		log.Printf("Not found conversion path: %q -> %q\n", from.Name, to.Name)
		return false
	}

	ctx.impls[Pair{from.Name, to.Name}] = struct{}{}

	ops := buildGenOps(ctx, path, true)
	ops = combineOps(ops)

	paramsVars := from.ChannelIdent()
	resultsVars := to.ChannelIdent()

	var scope VariableScope
	scope.ReserveAll(paramsVars...)

	var retString string

	// variable returnResult of params and results
	returnResult := scope.ContainsAny(resultsVars...)
	if returnResult {
		retString = joinRepeatN(FloatType, len(resultsVars))
	} else {
		scope.ReserveAll(resultsVars...)
		retString = joinIdentsWithType(FloatType, resultsVars...)
	}

	funcName := FuncName(from.Name, to.Name)

	w.Separate()
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

			fn := op.Pair.FuncName()

			if isLastOp {
				sub.ReturnInline()
				sub.WriteCallln(fn, inputVars...)
				returned = true
			} else {
				var outputVars []string
				outputVars = op.OutputVars
				if len(outputVars) == 0 {
					outputVars = scope.CreateIndexedUniqueN("f", len(inputVars))
				}

				sub.LineWriteJoin(outputVars, ", ")
				if scope.ContainsAll(outputVars...) {
					sub.Write(" = ")
				} else {
					sub.Write(" := ")
				}
				sub.WriteCallln(fn, inputVars...)

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
				sub.LineWriteln(v, " = math.Cbrt(", v, ")")
			}
		case OpCube:
			if idx > 0 {
				sub.Separate()
			}
			separate = true

			for _, v := range inputVars {
				sub.LineWriteln(v, " *= ", v, " * ", v)
			}

		case OpMatrix:
			if idx > 0 {
				sub.Separate()
			}
			separate = true

			outputVars := op.OutputVars

			needTemp := len(outputVars) == 0 || ContainsAny(inputVars, outputVars)
			if needTemp {
				outputVars = scope.CreateIndexedUniqueN("f", len(inputVars))
			}

			if isLastOp {
				if !needTemp && scope.ContainsAll(resultsVars...) {
					outputVars = resultsVars
				} else {
					resultsVars = outputVars
					returnResult = true
				}
			}

			l := len(inputVars)
			for i, out := range outputVars {
				sub.LineWrite(out)

				if scope.Reserve(out) {
					sub.Write(" := ")
				} else {
					sub.Write(" = ")
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
				sub.Newline()
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
