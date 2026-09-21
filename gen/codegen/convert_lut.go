package codegen

import (
	"github.com/thmalt/colors/gen/codegen/internal/convert"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genSrgbLUT(bits int) (dec, encThreshold []float64, encCoarse []uint16) {
	if bits > 16 {
		panic("colors: LUT supports at most 16 bits")
	}

	size := 1 << bits

	dec = make([]float64, size)
	encThreshold = make([]float64, size)
	encCoarse = make([]uint16, size+1)

	scale := float64(size - 1)
	coarseScale := float64(size)

	for i := 1; i < size; i++ {
		x := float64(i)

		dec[i] = convert.SrgbDecode(x / scale)
		encThreshold[i] = convert.SrgbDecode((x - 0.5) / scale)
	}

	var n uint16
	for i := range encCoarse {
		x := float64(i) / coarseScale

		for n < uint16(size-1) && x > encThreshold[n+1] {
			n++
		}

		encCoarse[i] = n
	}
	return
}

func genConvertPkgPrecomputedLUT(w *writer.GoWriter, bits int) {
	bitsSize := uintBits(bits)

	dec, encT, encC := genSrgbLUT(bitsSize)

	w.Separate()
	w.Comment(
		len(dec)*8+len(encT)*8+len(encC)*(bitsSize/8),
		" bytes",
	)

	w.Separate()
	w.BeginGroup("var")
	w.Begin("srgb", bitsSize, "DecLUT = [", len(dec), "]", FloatType)
	w.Indent()
	for i, v := range dec {
		if i > 0 {
			if i%8 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Write(formatFloat(v), ',')
	}
	w.End()

	w.Begin("srgb", bitsSize, "EncTLUT = [", len(encT), "]", FloatType)
	w.Indent()
	for i, v := range encT {
		if i > 0 {
			if i%8 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Write(formatFloat(v), ',')
	}
	w.End()

	w.Begin("srgb", bitsSize, "EncCLUT = [", len(encC), "]uint", bitsSize)
	w.Indent()
	for i, v := range encC {
		if i > 0 {
			if i%16 == 0 {
				w.Indent()
			} else {
				w.Write(' ')
			}
		}
		w.Write(v, ',')
	}
	w.End()

	w.Separate()
	w.End()
}

func genConvertPkgRuntimeInitLUT(w *writer.GoWriter, bits int) {
	if bits > 16 {
		panic("colors: LUT supports at most 16 bits")
	}

	bitsSize := uintBits(bits)
	size := 1 << bitsSize

	w.Separate()
	w.BeginGroup("var")
	w.LineWriteln("srgb", bitsSize, "DecLUT[", size, "]", FloatType)
	w.LineWriteln("srgb", bitsSize, "EncTLUT[", size, "]", FloatType)
	w.LineWriteln("srgb", bitsSize, "EncCLUT[", size+1, "]uint", bitsSize)
	w.End()

	w.Separate()
	w.Func("init")
	w.FuncBody()
	w.BeginGroup("const")
	w.LineWriteln("size        = 1 << ", bitsSize)
	w.LineWriteln("scale       = size - 1")
	w.LineWriteln("coarseScale = size")
	w.End()

	w.Separate()
	w.Begin("for i := 1; i < size; i++")
	w.LineWriteln("x := float64(i)")
	w.Separate()
	w.LineWriteln("srgb", bitsSize, "DecLUT[i] = SrgbDecode(x / scale)")
	w.LineWriteln("srgb", bitsSize, "EncTLUT[i] = SrgbDecode((x - 0.5) / scale)")
	w.End()

	w.Separate()
	w.LineWriteln("var n uint", bitsSize)
	w.Begin("for i := range ", "srgb", bitsSize, "EncCLUT")
	w.LineWriteln("x := float64(i) / coarseScale")
	w.Separate()
	w.Begin("for n < scale && x > ", "srgb", bitsSize, "EncTLUT[n+1]")
	w.LineWriteln("n++")
	w.End()

	w.Separate()
	w.LineWriteln("srgb", bitsSize, "EncCLUT[i] = n")
	w.End()

	w.End()
}
