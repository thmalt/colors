package colors

import (
	"errors"
	"unsafe"
)

// MarshalText encodes the color as a hexadecimal text representation.
func (c Color) MarshalText() ([]byte, error) {
	s := c.Hex()
	return unsafe.Slice(unsafe.StringData(s), len(s)), nil
}

// UnmarshalText decodes a color from its hexadecimal text representation.
func (c *Color) UnmarshalText(text []byte) error {
	s := unsafe.String(unsafe.SliceData(text), len(text))

	v, ok := TryHex(s)
	if !ok {
		return errors.New("invalid color: " + s)
	}

	*c = v
	return nil
}
