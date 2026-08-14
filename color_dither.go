package colors

import (
	"github.com/thmalt/colors/dither"
	_ "github.com/thmalt/colors/space"
)

// Dither applies ordered dithering to the sRGB color channels at pixel position (x, y).
// The dithering offset is scaled to the normalized sRGB [0, 1] range;
// the alpha channel is left unchanged. The returned color is always in [space.Srgb].
func (c Color) Dither(x, y int) Color {
	r, g, b := c.Srgb()
	d := dither.Offset(x, y) * (1 / 255.0)
	return SrgbAlpha(clamp01(r+d), clamp01(g+d), clamp01(b+d), c.alpha)
}
