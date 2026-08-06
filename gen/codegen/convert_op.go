package codegen

import "github.com/thmalt/colors/gen/codegen/data"

func prepareGenOps(ctx *Context, ops []Op, last *Node) []genOp {
	genOps := make([]genOp, len(ops))
	for i, op := range ops {
		genOps[i].Op = op
	}

	for i := range genOps {
		if genOps[i].Type != OpMatrix {
			continue
		}

		switch {
		case i+1 < len(ops) && ops[i+1].Type == OpCall:
			from := ctx.SpaceByName(genOps[i+1].Pair.From)
			genOps[i].outputVars = from.ChannelIdent()
		case i == len(genOps)-1:
			genOps[i].outputVars = last.To.ChannelIdent()
		}
	}

	return genOps
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

func combineOps(ops []Op) []Op {
	out := make([]Op, 0, len(ops))

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
			}
			continue
		}

		flush()
		out = append(out, op)
	}

	flush()

	return out
}
