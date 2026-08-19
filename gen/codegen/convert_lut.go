package codegen

import (
	"math"

	"github.com/thmalt/colors/gen/codegen/writer"
)

// makeLUT generates lookup data for a normalized domain [0, 1].
// size must be in the range [2, math.MaxUint16].
func makeLUT(size int, fn func(float64) float64) (lut, threshold []float64, coarse []uint16) {
	if size < 2 || size > math.MaxUint16 {
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

func genConvertPkgLUT(w *writer.GoWriter, linear, name string, transfer func(float64) float64) {
	lut, threshold, coarse := makeLUT(256, transfer)

	linear = toUpperCaseFirstChar(linear)
	name = toUpperCaseFirstChar(name)

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
		w.Write(formatNormalizedFloat(lut[i]))
		w.Write(',')
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
		w.Write(formatNormalizedFloat(threshold[i]))
		w.Write(',')
	}
	w.End()

	w.Separate()
	w.Begin(privLinearTo, "Coarse = [", len(coarse), "]uint8")
	w.Indent()
	for i := range coarse {
		w.Writef("%d", coarse[i])
		if i > 0 && i%16 == 0 {
			w.Write(',')
			w.Indent()
		} else {
			w.Write(", ")
		}
	}
	w.End()

	w.End()

	w.Separate()
	// func (r, g, b uint8) (float64, float64, float64)
	w.Comment(toLinear, " converts 8-bit components to linear components.")
	w.Func(toLinear)
	w.FuncParams("r, g, b uint8")
	w.FuncResults(joinRepeatN(FloatType, 3))
	w.FuncBody()
	w.Return(
		privToLinear, "LUT[r],",
		privToLinear, "LUT[g],",
		privToLinear, "LUT[b]",
	)
	w.End()

	w.Separate()
	// func (r, g, b float64) (uint8, uint8, uint8)
	w.Comment(linearTo, " converts linear components to 8-bit components.")
	w.Func(linearTo)
	w.FuncParams("r, g, b ", FloatType)
	w.FuncResults(joinRepeatN("uint8", 3))
	w.FuncBody()
	w.Return(
		privLinearTo, "(r),",
		privLinearTo, "(g),",
		privLinearTo, "(b)",
	)
	w.End()

	w.Separate()
	w.Func(privLinearTo)
	w.FuncParams("x ", FloatType)
	w.FuncResults("uint8")
	w.FuncBody()
	w.LineWriteln("x = min(1, max(0, x)) // clamp01")

	w.Separate()
	w.LineWriteln("i := int(x * ", privLinearTo, "CoarseSize)")
	w.LineWriteln("n := ", privLinearTo, "Coarse[i]")

	w.Separate()
	w.Begin("for n < 255 && x >= ", privLinearTo, "Threshold[n+1]")
	w.LineWriteln("n++")
	w.End()

	w.Separate()
	w.Return("n")
	w.End()
}
