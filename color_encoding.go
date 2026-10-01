package colors

import (
	"unsafe"
)

// MarshalText encodes the color as a hexadecimal text representation.
func (c Color) MarshalText() ([]byte, error) {
	return []byte(c.Hex()), nil
}

// UnmarshalText decodes a color from its CSS text representation.
func (c *Color) UnmarshalText(text []byte) error {
	if len(text) > 0 && text[0] == '#' {
		v, ok := TryHex(unsafe.String(unsafe.SliceData(text), len(text)))
		if ok {
			*c = v
			return nil
		}
	}

	v, err := ParseBytes(text)
	if err != nil {
		return err
	}
	*c = v
	return nil
}
