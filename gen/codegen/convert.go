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
	var pkgPath = filepath.Join(ctx.Directory, ctx.ConvertPkg.Path)

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))

	w.Import("math")

	genConvertPkgConversion(ctx, w)

	w.WriteGoFile(
		filepath.Join(pkgPath, ctx.ConvertPkg.Name+"_gen.go"),
		ctx.ConvertPkg.Name,
	)

	if convertGenWhitePoint(ctx, w) {
		w.WriteGoFile(
			filepath.Join(pkgPath, "whitepoint_gen.go"),
			ctx.ConvertPkg.Name,
		)
	}
}

func genConvertPkgConversion(ctx *Context, w *writer.GoWriter) {

	impls := map[Pair]struct{}{}

	for _, fn := range ctx.Funcs {
		if fn.Implemented {
			from, to := ctx.ResolveSpacePair(fn.Pair)
			if from == nil || to == nil {
				if from == nil {
					log.Printf("space of %s not found\n", fn.Pair.From)
				}

				if to == nil {
					log.Printf("space of %s not found\n", fn.Pair.To)
				}

				continue
			}
			impls[Pair{from.Name, to.Name}] = struct{}{}
		}
	}

	for i, from := range ctx.Spaces {
		if from == nil {
			log.Printf("space at index %d is nil\n", i)
			continue
		}

		if from.Disable {
			continue
		}

		for j := i + 1; j < len(ctx.Spaces); j++ {
			to := ctx.Spaces[j]
			if to == nil {
				log.Printf("space at index %d is nil\n", i)
				continue
			}

			if to.Disable {
				continue
			}

			if _, ok := impls[Pair{from.Name, to.Name}]; ok {
				continue
			}

			processPair(ctx, w, from, to)
			processPair(ctx, w, to, from)
		}
	}

}

func convertGenWhitePoint(ctx *Context, w *writer.GoWriter) bool {
	if len(ctx.WhitePoints) == 0 {
		return false
	}

	w.BeginGroup("var")
	for _, whitepoint := range ctx.WhitePoints {
		xyz := data.ChromaToXyz(whitepoint.X, whitepoint.Y)
		w.LineWriteln(
			whitepoint.Name,
			" = [3]",
			FloatType,
			"{",
			formatNormalizedFloat(xyz[0]),
			", ",
			formatNormalizedFloat(xyz[1]),
			", ",
			formatNormalizedFloat(xyz[2]),
			"}",
		)
	}
	w.End()

	return true
}

func processPair(ctx *Context, w *writer.GoWriter, from, to *model.Space) {
	path := ctx.Graph.FindPath(from, to)
	if len(path) == 0 {
		log.Printf("Not found conversion path: %q -> %q\n", from.Name, to.Name)
		return
	}

	ops := buildOps(ctx, path)
	newOps, changed := combineOps(ops)
	if changed {
		ops = newOps
	}

	params := from.ChannelSymbols()
	results := to.ChannelSymbols()

	var state VarState
	state.Reserve(params...)

	var retString string

	// variable conflict of params and results
	conflict := state.ContainsAny(results...)
	if conflict {
		retString = valueTypeRepeat(len(results))
	} else {
		state.Reserve(results...)
		retString = varJoinWithType(results...)
	}

	funcName := FuncName(from.Name, to.Name)
	w.LineComment("Conversion path (", len(path), " steps):")
	w.LineComment("")
	w.LineComment("\t", from.DisplayName)
	for _, node := range path {
		w.LineComment("\t-> ", node.To.DisplayName)
	}

	// w.LineComment(stringifyPath(path))
	w.Func(funcName)
	w.FuncParams(varJoinWithType(params...))
	w.FuncResults(retString)
	w.FuncBody()

	wop := w.NewTemp()
	prevVars := params

	firstOp := true
	hasReturn := false
	last := len(ops) - 1

	tempVar := []string{"f", "f", "f"}

	for idx, op := range ops {
		isLastOp := idx == last
		switch op.Type {
		case OpCall:
			firstOp = false
			if isLastOp {
				wop.Return(op.Pair.FuncName(), "(", strings.Join(prevVars, ", "), ")")
				hasReturn = true
			} else {
				_, to := ctx.ResolveSpacePair(op.Pair)
				nextVars := to.ChannelSymbols()
				assign := " := "
				if state.ContainsAll(nextVars...) {
					assign = " = "
				}
				state.Reserve(nextVars...)
				wop.LineWriteln(strings.Join(nextVars, ", ")+assign, op.Pair.FuncName(), "(", strings.Join(prevVars, ", "), ")")
				prevVars = nextVars
			}
		case OpCbrt:
			if !firstOp {
				wop.Writeln()
			}
			firstOp = false
			for _, v := range prevVars {
				wop.NewlineWrite(v, " = math.Cbrt(", v, ")")
			}
			wop.Writeln()
		case OpCube:
			if !firstOp {
				wop.Writeln()
			}
			firstOp = false
			for _, v := range prevVars {
				wop.NewlineWrite(v, " *= ", v, " * ", v)
			}
			wop.Writeln()
		case OpMatrix:
			if !firstOp {
				wop.Writeln()
			}

			firstOp = false
			next := state.ReserveNumberAddNames(tempVar...)
			decl := true
			if isLastOp {
				if !ContainsAny(prevVars, results) && state.ContainsAll(results...) {
					next = results
					decl = false
				} else {
					results = next
					conflict = true
				}
			}
			l := len(prevVars)
			for i := range prevVars {
				if decl {
					wop.LineWrite(next[i], " := ")
				} else {
					wop.LineWrite(next[i], " = ")
				}

				first := true
				for j, v := range prevVars {
					f := normalizeFloat(op.Matrix[i*l+j])
					if f == 0 {

					} else if math.Signbit(f) {
						wop.Write('-')
						f = -f
					} else if !first {
						wop.Write('+')
					}

					switch f {
					case 0:
					case 1:
						wop.Write(v)
						first = false
					default:
						wop.Write(formatFloat(f), '*', v)
						first = false
					}
				}
			}
			prevVars = next
			if idx != 0 {
				wop.NewlineWriteln()
			}
		}
	}

	w.Write(wop.Bytes())

	if !hasReturn {
		if conflict {
			w.Return(strings.Join(prevVars, ", "))
		} else {
			w.Return()
		}
	}

	w.End()

	w.Writeln()
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

	var ()

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
