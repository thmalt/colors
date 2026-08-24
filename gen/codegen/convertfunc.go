package codegen

import (
	"slices"
	"strconv"
)

type ConvertFunc struct {
	Pair Pair
	Cost int
	Ops  []Op

	Implemented bool
}

type OpType uint

const (
	OpOther OpType = iota
	OpCall
	OpTransfer
	OpCbrt
	OpCube
	OpMatrix
)

func (op OpType) String() string {
	switch op {
	case OpOther:
		return "OpOther"
	case OpCall:
		return "OpCall"
	case OpTransfer:
		return "OpTransfer"
	case OpCbrt:
		return "OpCbrt"
	case OpCube:
		return "OpCube"
	case OpMatrix:
		return "OpMatrix"
	}
	return strconv.FormatUint(uint64(op), 10)
}

type Op struct {
	Type OpType
	Pair Pair

	Transfer string

	Func   Pair
	Matrix *[9]float64
}

type Pair struct {
	From string
	To   string
}

func (p Pair) IsNone() bool {
	return p.From == "" && p.To == ""
}

func (p Pair) FuncName() string {
	return FuncName(p.From, p.To)
}

func FuncName(from, to string) string {
	return from + "To" + to
}

func opCall(from, to string) Op {
	return Op{
		Type: OpCall,
		Func: Pair{from, to},
	}
}

func opCbrt() Op {
	return Op{Type: OpCbrt}
}

func opCube() Op {
	return Op{Type: OpCube}
}

func opTransfer(transfer string) Op {
	return Op{
		Type:     OpTransfer,
		Transfer: transfer,
	}
}

func opMatrix(m [9]float64) Op {
	return Op{
		Type:   OpMatrix,
		Matrix: &m,
	}
}

func transferFunc(from, to string, transfer string) ConvertFunc {
	return ConvertFunc{
		Pair: Pair{from, to},
		Ops:  []Op{opTransfer(transfer)},
	}
}

func implementedFunc(from, to string) ConvertFunc {
	return ConvertFunc{
		Pair:        Pair{from, to},
		Implemented: true,
	}
}

func convertFunc(from, to string, ops ...Op) ConvertFunc {
	return ConvertFunc{
		Pair: Pair{from, to},
		Ops:  slices.Clone(ops),
	}
}
