package colors

import "github.com/thmalt/colors/space"

func (c Color) MustTo(dst space.Space) Color {
	to, err := c.To(dst)
	if err != nil {
		panic(err)
	}

	return to
}
