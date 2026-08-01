package colors

const hexDigits = "0123456789abcdef"

func encodeHexByte(dst []byte, b byte) {
	dst[0] = hexDigits[b>>4]
	dst[1] = hexDigits[b&0x0f]
}
