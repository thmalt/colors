package codegen

import (
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
