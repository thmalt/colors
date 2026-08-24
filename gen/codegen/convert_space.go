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
	for i, to := range ctx.BuiltSpaces {
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

func buildConversion(ctx *Context, from, to *model.Space) ([]*Node, []GenOp) {
	path := ctx.Graph.FindPath(from, to)
	ops := buildGenOps(ctx, path, true)
	ops = combineOps(ops)

	if len(ops) == 0 {
		return path, ops
	}

	eqFrom := ctx.SpaceByName(from.Equivalent)
	eqTo := ctx.SpaceByName(to.Equivalent)

	if from.Equivalent != "" && (eqFrom == nil || eqFrom.Disable) {
		log.Fatalln("Equivalent space", from.Equivalent, "is unavailable")
	}
	if to.Equivalent != "" && (eqTo == nil || eqTo.Disable) {
		log.Fatalln("Equivalent space", to.Equivalent, "is unavailable")
	}

	if eqFrom != nil || eqTo != nil {
		if eqFrom == to || eqTo == from {
			path = []*Node{{
				From: from,
				To:   to,
			}}
			ops = nil
		} else if ops[0].Type != OpTransfer {
			f := from
			t := to

			if eqFrom != nil {
				f = eqFrom
			}
			if eqTo != nil {
				t = eqTo
			}

			ops = []GenOp{{
				Op: Op{
					Type: OpCall,
					Pair: Pair{from.Name, to.Name},
					Func: Pair{f.Name, t.Name},
				},
			}}
		}
	}

	eqFilter := make([]GenOp, 0, len(ops))
	for i := 0; i < len(ops); i++ {
		if ops[i].Type == OpCall {
			pair := ops[i].Pair

			f := ctx.SpaceByName(pair.From)
			t := ctx.SpaceByName(pair.To)
			if t != nil && f != nil {
				eqFrom := ctx.SpaceByName(t.Equivalent)
				eqTo := ctx.SpaceByName(f.Equivalent)
				if f == eqFrom || t == eqTo {
					continue
				}
			}
		}

		eqFilter = append(eqFilter, ops[i])
	}
	ops = eqFilter

	return path, ops
}

func genConvertPkgSpacePair(ctx *Context, w *writer.GoWriter, from, to *model.Space) bool {
	path, ops := buildConversion(ctx, from, to)

	if len(path) == 0 {
		log.Printf("Not found conversion path: %q -> %q\n", from.Name, to.Name)
		return false
	}

	ctx.impls[Pair{from.Name, to.Name}] = struct{}{}

	// ops := buildGenOps(ctx, path, expand)
	// ops = combineOps(ops)

	if !ctx.Opts.EmbedMatrix {
		ops = replaceMatrixWithCall(ops)
	}

	paramsVars := from.ChannelIdent()
	resultsVars := to.ChannelIdent()

	var scope VariableScope
	scope.ReserveAll(paramsVars...)

	var retString string

	// variable hasReturnVars of params and results
	hasReturnVars := scope.ContainsAny(resultsVars...)
	if hasReturnVars {
		retString = joinRepeatN(FloatType, len(resultsVars))
	} else {
		scope.ReserveAll(resultsVars...)
		retString = joinIdentsWithType(FloatType, resultsVars...)
	}

	funcName := FuncName(from.Name, to.Name)

	w.Separate()

	eqFrom, eqTo := ctx.ResolveSpacePairName(from.Equivalent, to.Equivalent)
	if eqFrom != nil || eqTo != nil {
		f := from
		t := to

		if eqFrom != nil {
			f = eqFrom
		}
		if eqTo != nil {
			t = eqTo
		}

		if f != t && (len(ops) > 0 && ops[0].Type != OpTransfer) {
			w.Comment("Calls [", Pair{f.Name, t.Name}.FuncName(), ']')
			w.Comment()
		}
	}

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

			fn := op.Func.FuncName()

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
		case OpTransfer:
			if separate {
				sub.Separate()
			}
			separate = false

			if isLastOp {
				sub.ReturnInline()
				for i := range inputVars {
					if i > 0 {
						sub.Write(", ")
					}
					sub.Write(op.Transfer, '(', inputVars[i], ')')
				}
				returned = true
			} else {
				outputVars := op.OutputVars
				sub.LineWriteJoin(outputVars, ", ")
				if scope.ContainsAll(outputVars...) {
					sub.Write(" = ")
				} else {
					sub.Write(" := ")
				}
				for i := range inputVars {
					if i > 0 {
						sub.Write(", ")
					}
					sub.Write(op.Transfer, '(', inputVars[i], ')')
				}
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
					hasReturnVars = true
				}
			}

			l := len(inputVars)
			for i, out := range outputVars {
				startLen := sub.Len()

				sub.LineWrite(out)

				if scope.Reserve(out) {
					sub.Write(" := ")
				} else {
					sub.Write(" = ")
				}

				varName := ""
				count := 0

				for j, v := range inputVars {
					f := normalizeFloat(op.Matrix[i*l+j])
					if f == 0 {
						continue
					}
					count++
					neg := math.Signbit(f)
					if neg {
						f = -f
						if count == 1 {
							sub.Write('-')
						} else {
							sub.Write(" - ")
						}
					} else if count > 1 {
						sub.Write(" + ")
					}

					if f == 1 {
						sub.Write(v)
						if count == 1 && !neg {
							varName = v
						} else {
							varName = ""
						}
					} else {
						sub.Write(formatFloat(f), '*', v)
						varName = ""
					}
				}
				sub.Newline()

				if varName != "" && (!isLastOp && ops[idx+1].Type == OpCall || isLastOp && hasReturnVars) {
					sub.Truncate(startLen)
					outputVars[i] = inputVars[i]
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

		if hasReturnVars {
			w.Return(strings.Join(inputVars, ", "))
		} else {
			w.Return()
		}
	}

	w.End()

	return true
}
