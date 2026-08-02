package colors_test

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func ExampleColor_To() {
	c := colors.Rgb(50, 60, 70).WithAlpha(0.995)

	// loop through available color spaces.
	for s := range space.SpaceCount {
		to, err := c.To(s)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(to)
		fmt.Println()
	}

	fmt.Println("hex:", c.WithAlpha(1).Hex())
	// hex with alpha if alpha != 1
	fmt.Println("hex:", c.WithAlpha(0.995).Hex())

	// Output:
	// color(xyz-d65 0.04036525 0.04351959 0.06421888 / 0.995)
	//
	// color(xyz-d50 0.04007529 0.04320278 0.04856662 / 0.995)
	//
	// color(srgb-linear 0.031896 0.045186 0.061246 / 0.995)
	//
	// color(srgb 0.196078 0.235294 0.27451 / 0.995)
	//
	// color(display-p3-linear 0.034256 0.044745 0.059582 / 0.995)
	//
	// color(display-p3 0.203657 0.23411 0.270749 / 0.995)
	//
	// color(a98-rgb-linear 0.035682 0.045186 0.060585 / 0.995)
	//
	// color(a98-rgb 0.219679 0.244581 0.279469 / 0.995)
	//
	// color(prophoto-rgb-linear 0.040409 0.044331 0.058861 / 0.995)
	//
	// color(prophoto-rgb 0.168199 0.177081 0.207287 / 0.995)
	//
	// color(rec2020-linear 0.037544 0.04445 0.059351 / 0.995)
	//
	// color(rec2020 0.151693 0.17151 0.209135 / 0.995)
	//
	// hsl(210 16.67% 23.53% / 0.995)
	//
	// hsv(210 28.57% 27.45% / 0.995)
	//
	// hwb(210 19.61% 72.55% / 0.995)
	//
	// lab(24.7032 -2.2538 -7.6209 / 0.995)
	//
	// lch(24.7032 7.9472 253.5253 / 0.995)
	//
	// oklab(0.351128 -0.00808 -0.020463 / 0.995)
	//
	// oklch(0.351128 0.022001 248.4537 / 0.995)
	//
	// hex: #323c46
	// hex: #323c46fe
}

func ExampleColor_Space() {
	c := colors.Rgb(50, 60, 70).WithAlpha(0.995)

	// color space info
	inf := c.Space().Info()
	if inf == nil {
		fmt.Println("invalid color space: ", c.Space())
		return
	}

	fmt.Println("Space name:", inf.DisplayName())

	fmt.Println()

	// get whitepoint info.
	whitePoint := inf.WhitePoint()
	fmt.Println("WhitePoint")
	fmt.Println("  Name:", whitePoint.Name)
	fmt.Println("  Z:", whitePoint.X)
	fmt.Println("  Y:", whitePoint.Y)
	fmt.Println("  Z:", whitePoint.Z)

	fmt.Println()

	// get channel info.
	fmt.Println("Channels:", inf.ChannelCount())
	for i, ch := range inf.Channels() {

		fmt.Println(
			" ",
			"Channel:", i,
			"Name:", ch.Name,
			"Min:", ch.Min,
			"Max:", ch.Max,
			"Unrestricted:", ch.Unrestricted,
		)
	}

	// Output:
	// Space name: sRGB
	//
	// WhitePoint
	//   Name: D65
	//   Z: 0.9504559270516716
	//   Y: 1
	//   Z: 1.0890577507598784
	//
	// Channels: 3
	//   Channel: 0 Name: Red Min: 0 Max: 1 Unrestricted: false
	//   Channel: 1 Name: Green Min: 0 Max: 1 Unrestricted: false
	//   Channel: 2 Name: Blue Min: 0 Max: 1 Unrestricted: false
}
