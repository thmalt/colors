package colors

import "github.com/thmalt/colors/space"

// MutTo converts the color to dst in place and returns the receiver.
func (c *Color) MutTo(dst space.Space) *Color {
	*c, _ = c.to(dst)
	return c
}
