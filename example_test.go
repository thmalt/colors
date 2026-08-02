package colors_test

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func ExampleColor_To() {
	c := colors.Rgb(50, 60, 70).WithAlpha(0.995)

	c1, err := c.To(space.XyzD65)
	if err != nil {
		return
	}

	fmt.Println(c1.Space())

	c2, err := c.To(space.Oklab)
	if err != nil {
		return
	}

	fmt.Println(c2.Space())

	// Output:
	// XyzD65
	// Oklab
}

func ExampleColor_Space() {
	c := colors.Rgb(50, 60, 70).WithAlpha(0.995)

	// color space info
	inf := c.Space().Info()
	if inf == nil {
		fmt.Println("invalid color space: ", c.Space())
		return
	}

	fmt.Println("Space display name:", inf.DisplayName())

	// Output:
	// Space display name: sRGB
}
