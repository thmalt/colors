package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func main() {
	c := colors.Rgb(50, 60, 70)

	// loop through available color spaces.
	for s := range space.SpaceCount {
		to := c.MustTo(s).WithAlpha(0.995)

		fmt.Println(to)
		fmt.Println()
	}

	// color space info
	inf := c.Space().Info()
	if inf == nil {
		fmt.Println("invalid color space: ", c.Space())
		return
	}

	fmt.Println("space name:", inf.DisplayName())

	// get channel info.
	for i, ch := range inf.Channels() {
		fmt.Println("channel:", i, "name:", ch.Name, "min:", ch.Min, "max:", ch.Max)
	}

	fmt.Println()

	fmt.Println("hex:", c.Hex())
	// hex with alpha if alpha != 1
	fmt.Println("hex:", c.WithAlpha(0.998).Hex())
}
