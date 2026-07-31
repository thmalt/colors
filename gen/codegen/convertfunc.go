package codegen

import (
	"strconv"
	"strings"
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
	Type   OpType
	Pair   Pair
	Matrix *[9]float64
}

type Pair struct {
	From string
	To   string
}

func (p Pair) FuncName() string {
	return FuncName(string(p.From), string(p.To))
}

// LinearTo // to = ""
func FuncName(from, to string) string {
	if strings.TrimPrefix(from, "Linear") == to {
		return "LinearTo" + to
	}
	if strings.TrimPrefix(to, "Linear") == from {
		return from + "ToLinear"
	}

	return from + "To" + to
}
