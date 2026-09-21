package codegen

import (
	"math"

	"github.com/thmalt/colors/gen/codegen/writer"
)

const hexEnc = "0123456789abcdef"

func generateHexLUTPkg(ctx *Context) {
	if !ctx.Opts.genLUT() {
		return
	}

	pkg := Pkg{Name: "hexlut", Path: "internal/hexlut"}

	w := newWriter(ctx)

	emitGoFile(ctx, pkg, w, "lut", func(w *writer.GoWriter) {
		genHexLUTPkgDec(w)
	})

	emitGoFile(ctx, pkg, w, "lut_little", func(w *writer.GoWriter) {
		w.AddBuildTags("386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh || wasm")

		genHexLUTPkgEnc(w, true)
	})

	emitGoFile(ctx, pkg, w, "lut_big", func(w *writer.GoWriter) {
		w.AddBuildTags("armbe || arm64be || m68k || mips || mips64 || mips64p32 || ppc || ppc64 || s390 || s390x || shbe || sparc || sparc64")

		genHexLUTPkgEnc(w, false)
	})
}

func genHexLUTPkgDec(w *writer.GoWriter) {
	w.Separate()
	w.BeginGroup("const")

	const nibbleLimit = 16
	w.Comment("NibbleLimit is the exclusive upper bound for valid nibble values.")
	w.LineWriteln("NibbleLimit = ", nibbleLimit)

	w.Separate()
	w.Comment("Dec maps each byte value (0..255) to its 4-bit nibble value (0..15).")
	w.Comment("Invalid characters are represented by 0xFF (>= [NibbleLimit]).")
	w.LineWriteln(`Dec = "" +`)

	w.In()

	const count = math.MaxUint8 + 1
	for i := range count {
		if i%16 == 0 {
			if i > 0 {
				w.Writeln('"', " +")
			}
			w.LineWrite('"')
		}

		switch c := uint8(i); {
		case c >= '0' && c <= '9':
			w.Write(`\x0`, c)
		case c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
			w.Write(`\x0`, c|0x20)
		default:
			w.Write(`\xff`)
		}
	}
	w.Writeln('"')
	w.Out()

	w.End()
}

func genHexLUTPkgEnc(w *writer.GoWriter, littleEndian bool) {
	w.Separate()
	w.Comment("Enc maps each byte value (0..255) to its two-digit hexadecimal representation.")
	w.LineWriteln(`const Enc = "" +`)
	w.In()

	const count = math.MaxUint8 + 1
	for i := range count {
		if i%16 == 0 {
			if i > 0 {
				w.Writeln('"', " +")
			}
			w.LineWrite('"')
		}

		c := uint8(i)
		if littleEndian {
			w.Write(hexEnc[c>>4], hexEnc[c&0x0f])
		} else {
			w.Write(hexEnc[c&0x0f], hexEnc[c>>4])
		}
	}
	w.Writeln('"')

	w.Out()
}
