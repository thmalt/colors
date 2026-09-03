package colors

import "github.com/thmalt/colors/space"

// MutTo converts the color to the specified color space in place
// and returns the receiver.
func (c *Color) MutTo(dst space.Space) *Color {
	c.mutTo(dst)
	return c
}

// TryMutTo converts the color to the specified color space in place
// and reports whether the conversion succeeded.
func (c *Color) TryMutTo(dst space.Space) bool {
	return c.mutTo(dst)
}
