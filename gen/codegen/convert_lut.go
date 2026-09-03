package codegen

import (
	"math"
	"strconv"

	"github.com/thmalt/colors/gen/codegen/writer"
)

// makeLUT generates lookup data for a normalized domain [0, 1].
// size must be in the range [2, math.MaxUint16).
func makeLUT(size int, fn func(float64) float64) (lut, threshold []float64, coarse []uint16) {
	if size < 2 || size > math.MaxUint16+1 {
		panic("invalid LUT size")
	}

	lut = make([]float64, size)
	threshold = make([]float64, size)
	coarse = make([]uint16, size+1)

	div := float64(size - 1)
	coarseDiv := float64(size)

	for i := 1; i < size; i++ {
		x := float64(i)

		lut[i] = fn(x / div)
		threshold[i] = fn((x - 0.5) / div)
	}

	for i := range coarse {
		x := float64(i) / coarseDiv

		n := uint16(0)
		for n < uint16(size-1) && x > threshold[n+1] {
			n++
		}

		coarse[i] = n
	}

	return
}

func genConvertPkgLUT(w *writer.GoWriter, size int, linear, name string, transfer func(float64) float64) {
	lut, threshold, coarse := makeLUT(size+1, transfer)

	bits := strconv.Itoa(smallestUintType(size))
	uintType := "uint" + bits

	linear = toUpperCaseFirstChar(linear)
	name = toUpperCaseFirstChar(name) + bits

	linearToUbits := linear + "ToU" + bits
	linearTo := linear + "To" + name
	toLinear := name + "To" + linear

	privLinearTo := toLowerCaseFirstWord(linearTo)
	privToLinear := toLowerCaseFirstWord(toLinear)

	w.LineWrite("const ", privLinearTo, "CoarseSize = ", len(threshold))
	w.Separate()
	w.BeginGroup("var ")

	w.Begin(privToLinear, "LUT = [", len(lut), "]", FloatType)
	w.Indent()
	for i := range lut {
		if i > 0 {
			if i%8 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Write(formatNormalizedFloat(lut[i]), ',')
	}
	w.End()

	w.Separate()
	w.Begin(privLinearTo, "Threshold = [", len(threshold), "]", FloatType)
	w.Indent()
	for i := range threshold {
		if i > 0 {
			if i%8 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Write(formatNormalizedFloat(threshold[i]), ',')
	}
	w.End()

	w.Separate()
	w.Begin(privLinearTo, "Coarse = [", len(coarse), "]", uintType)
	w.Indent()
	for i := range coarse {
		if i > 0 {
			if i%16 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Writef("%d,", coarse[i])
	}
	w.End()

	w.End()

	w.Separate()
	// func (r, g, b uint?) (float64, float64, float64)
	w.Comment(toLinear, " converts ", bits, "-bit components to linear components.")
	w.Func(toLinear)
	w.FuncParams("r, g, b ", uintType)
	w.FuncResults(joinRepeatN(FloatType, 3))
	w.FuncBody()
	w.Return(
		privToLinear, "LUT[r], ",
		privToLinear, "LUT[g], ",
		privToLinear, "LUT[b]",
	)
	w.End()

	w.Separate()
	// func (r, g, b float64) (uint?, uint?, uint?)
	w.Comment(linearTo, " converts linear components to ", bits, "-bit components.")
	w.Func(linearTo)
	w.FuncParams("r, g, b ", FloatType)
	w.FuncResults(joinRepeatN(uintType, 3))
	w.FuncBody()
	w.Return(
		linearToUbits, "(r), ",
		linearToUbits, "(g), ",
		linearToUbits, "(b)",
	)
	w.End()

	w.Separate()
	w.Comment(linearToUbits, " converts a linear color component to an unsigned integer value.")
	w.Func(linearToUbits)
	w.FuncParams("x ", FloatType)
	w.FuncResults(uintType)
	w.FuncBody()
	w.LineWriteln("x = min(1, max(0, x)) // clamp01")

	w.Separate()
	w.LineWriteln("i := int(x * ", privLinearTo, "CoarseSize)")
	w.LineWriteln("n := ", privLinearTo, "Coarse[i]")

	w.Separate()
	w.Begin("for n < ", size, " && x >= ", privLinearTo, "Threshold[n+1] ")
	w.LineWriteln("n++")
	w.End()

	w.Separate()
	w.Return("n")
	w.End()
}
