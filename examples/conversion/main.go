package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func main() {
	c := colors.RgbAlpha(50, 60, 70, 0.995)

	// Convert a color to different color spaces.
	fmt.Println("sRGB: ", c)
	fmt.Println("Oklab:", c.To(space.Oklab))
	fmt.Println("Oklch:", c.To(space.Oklch))
	fmt.Println("Lab:  ", c.To(space.Lab))
	fmt.Println("LCh:  ", c.To(space.Lch))

	fmt.Println()

	// Serialize the color as hexadecimal.
	fmt.Println("Hex:           ", c.WithAlpha(1).Hex())
	fmt.Println("Hex with alpha:", c.Hex())

	fmt.Println()

	// Inspect color space metadata.
	info := c.Space().Info()
	if info == nil {
		fmt.Println("invalid color space:", c.Space())
		return
	}

	fmt.Println("Color space:", info.DisplayName())

	whitePoint := info.WhitePoint()
	fmt.Println("White point:", whitePoint.Name)
	fmt.Println("  X:", whitePoint.X)
	fmt.Println("  Y:", whitePoint.Y)
	fmt.Println("  Z:", whitePoint.Z)

	fmt.Println("Channels:", info.ChannelCount())
	for i, ch := range info.Channels() {
		fmt.Printf(
			"  %d: %s [%g, %g], unrestricted=%t\n",
			i, ch.Name, ch.Min, ch.Max, ch.Unrestricted,
		)
	}
}
