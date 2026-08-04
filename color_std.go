package colors

// RGBA implements the [image/color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	red, green, blue := c.Srgb()

	red = clamp01(red)
	green = clamp01(green)
	blue = clamp01(blue)
	alpha := clamp01(c.alpha)

	const max = 65535.0

	r = uint32(red*alpha*max + 0.5)
	g = uint32(green*alpha*max + 0.5)
	b = uint32(blue*alpha*max + 0.5)
	a = uint32(alpha*max + 0.5)

	return
}
