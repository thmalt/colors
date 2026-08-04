package colors

import "math"

const hexDigits = "0123456789abcdef"

func encodeHexByte(dst []byte, b byte) {
	dst[0] = hexDigits[b>>4]
	dst[1] = hexDigits[b&0x0f]
}

func (c Color) Hex() string {
	const max = 255

	r, g, b := c.Srgb()

	red := byte(math.Round(clamp01(r) * max))
	green := byte(math.Round(clamp01(g) * max))
	blue := byte(math.Round(clamp01(b) * max))
	alpha := byte(math.Round(clamp01(c.alpha) * max))

	var out []byte = make([]byte, 9)
	out[0] = '#'

	encodeHexByte(out[1:], red)
	encodeHexByte(out[3:], green)
	encodeHexByte(out[5:], blue)

	if alpha == max {
		return string(out[:7])
	}

	encodeHexByte(out[7:], alpha)

	return string(out[:])
}
