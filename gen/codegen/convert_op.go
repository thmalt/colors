package codegen

import (
	"github.com/thmalt/colors/gen/codegen/data"
)

type GenOp struct {
	Op

	InputVars  []string
	OutputVars []string
}

func buildGenOps(ctx *Context, path []*Node, expand bool) []GenOp {
	var ops []GenOp

	if !expand && len(path) == 1 {
		expand = true
	}

	for _, node := range path {
		nodeOps := buildNodeGenOps(ctx, node, expand)
		assignNodeVars(ctx, node, nodeOps)
		ops = append(ops, nodeOps...)
	}

	clearIntermediateVars(ops)

	return ops
}

func buildNodeGenOps(ctx *Context, node *Node, expand bool) []GenOp {
	if !expand || len(node.Fn.Ops) == 0 {
		return []GenOp{{
			Op: Op{
				Type: OpCall,
				Pair: node.Fn.Pair,
				Func: node.Fn.Pair,
			},
			InputVars:  node.From.ChannelIdent(),
			OutputVars: node.To.ChannelIdent(),
		}}
	}

	var ops []GenOp
	in := node.From.ChannelIdent()

	for i, op := range node.Fn.Ops {
		genOps := expandOps(ctx, op)
		if len(genOps) == 0 {
			continue
		}

		genOps[0].InputVars = in

		if next := nextCall(node.Fn.Ops, i); next != nil {
			if sp := ctx.SpaceByName(next.Pair.From); sp != nil {
				out := sp.ChannelIdent()
				genOps[len(genOps)-1].OutputVars = out
				in = out
			} else {
				genOps[len(genOps)-1].OutputVars = nil
				in = nil
			}
		} else {
			out := node.To.ChannelIdent()
			genOps[len(genOps)-1].OutputVars = out
			in = out
		}

		ops = append(ops, genOps...)
	}

	if ops[0].Pair.IsNone() {
		ops[0].Pair = node.Fn.Pair
	}

	if n := len(ops) - 1; ops[n].Pair.IsNone() {
		ops[n].Pair = node.Fn.Pair
	}

	return ops
}

func assignNodeVars(ctx *Context, node *Node, ops []GenOp) {
	if len(ops) == 0 {
		return
	}

	in := node.From.ChannelIdent()

	for i := range ops {
		if ops[i].Type != OpCall {
			continue
		}

		ops[i].InputVars = in

		next := -1
		for j := i + 1; j < len(ops); j++ {
			if ops[j].Type == OpCall {
				next = j
				break
			}
		}

		if next == -1 {
			ops[i].OutputVars = node.To.ChannelIdent()
			break
		}

		if sp := ctx.SpaceByName(ops[next].Pair.From); sp != nil {
			out := sp.ChannelIdent()
			ops[i].OutputVars = out
			in = out
		} else {
			ops[i].OutputVars = nil
			in = nil
		}
	}
}

func expandOps(ctx *Context, op Op) []GenOp {
	if op.Type != OpCall {
		return []GenOp{{Op: op}}
	}

	fn := ctx.ConvertFuncByPair(op.Pair)
	if fn == nil || len(fn.Ops) == 0 {
		return []GenOp{{Op: op}}
	}

	var out []GenOp
	for _, op := range fn.Ops {
		out = append(out, expandOps(ctx, op)...)
	}
	return out
}

func nextCall(ops []Op, i int) *Op {
	for i++; i < len(ops); i++ {
		if ops[i].Type == OpCall {
			return &ops[i]
		}
	}
	return nil
}

func clearIntermediateVars(ops []GenOp) {
	first := -1
	last := -1

	flush := func() {
		if first < 0 {
			return
		}

		for i := first; i < last; i++ {
			ops[i].OutputVars = nil
			if i+1 < len(ops) {
				ops[i+1].InputVars = nil
			}
		}

		first = -1
		last = -1
	}

	for i := range ops {
		if ops[i].Type == OpCall {
			flush()
			continue
		}

		if first < 0 {
			first = i
		}
		last = i
	}

	flush()
}

func combineOps(ops []GenOp) []GenOp {
	out := make([]GenOp, 0, len(ops))

	var (
		mat       [9]float64
		hasMat    bool
		firstPair Pair
		lastPair  Pair
		input     []string
		output    []string
	)

	flush := func() {
		if !hasMat {
			return
		}

		m := mat
		if lastPair.IsNone() {
			lastPair = firstPair
		}

		out = append(out, GenOp{
			Op: Op{
				Type:   OpMatrix,
				Pair:   Pair{firstPair.From, lastPair.To},
				Matrix: &m,
			},
			InputVars:  input,
			OutputVars: output,
		})

		hasMat = false
		input = nil
		output = nil
	}

	for _, op := range ops {
		if op.Type == OpMatrix {
			if !hasMat {
				mat = *op.Matrix
				firstPair = op.Pair
				input = op.InputVars // first matrix
				output = op.OutputVars
				hasMat = true
			} else {
				mat = data.Mat3MulFMA(*op.Matrix, mat)
				lastPair = op.Pair
				output = op.OutputVars // last matrix
			}
			continue
		}

		flush()
		out = append(out, op)
	}

	flush()

	return out
}

func replaceMatrixWithCall(ops []GenOp) []GenOp {
	var out []GenOp
	var replaced bool

	n := len(ops)
	for i := 0; i < n; {
		op := ops[i]
		if op.Type == OpCall {
			replaced = true
		}

		if i+2 < n &&
			ops[i].Type == OpMatrix &&
			(ops[i+1].Type == OpCbrt || ops[i+1].Type == OpCube) &&
			ops[i+2].Type == OpMatrix {
			if out == nil {
				out = make([]GenOp, 0, n)
				out = append(out, ops[:i]...)
			}

			pair := Pair{ops[i].Pair.From, ops[i+2].Pair.To}
			out = append(out, GenOp{
				Op: Op{
					Type: OpCall,
					Func: pair,
				},
				InputVars:  op.InputVars,
				OutputVars: ops[i+2].OutputVars,
			})

			i += 3
			continue
		}

		if op.Type == OpMatrix {
			if out == nil {
				out = make([]GenOp, 0, n)
				out = append(out, ops[:i]...)
			}

			out = append(out, GenOp{
				Op: Op{
					Type: OpCall,
					Func: op.Pair,
				},
				InputVars:  op.InputVars,
				OutputVars: op.OutputVars,
			})
			i++
			continue
		}

		if out != nil {
			out = append(out, op)
		}
		i++
	}

	if replaced && len(out) > 0 {
		return out
	}
	return ops
}

func isShareableOps(ops []GenOp) bool {
	switch len(ops) {
	case 1:
		return ops[0].Type == OpMatrix
	case 3:
		return ops[0].Type == OpMatrix && (ops[1].Type == OpCbrt || ops[1].Type == OpCube) && ops[2].Type == OpMatrix
	default:
		return false
	}
}
