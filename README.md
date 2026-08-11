# colors

Modern Go library for color space conversions and CSS color serialization.

## Usage

```go
import "github.com/thmalt/colors"
```

## Examples

### Color and Color Spaces

```go
package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func main() {
	c := colors.RgbAlpha(50, 60, 70, 0.995) // or colors.Rgb(50, 60, 70).WithAlpha(0.995)

	// Loop over all available color spaces.
	for s := space.FirstSpace; s < space.SpaceCount; s++ {
		fmt.Println(c.MustTo(s))
	}

	fmt.Println()

	// color space info
	inf := c.Space().Info()
	if inf == nil {
		fmt.Println("invalid color space: ", c.Space())
		return
	}

	fmt.Println("Space display name:", inf.DisplayName())

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

	fmt.Println()

	fmt.Println("hex:", c.WithAlpha(1).Hex())
	// hex with alpha if alpha != 1
	fmt.Println("hex:", c.WithAlpha(0.995).Hex())
}
```

### Linear Gradient

```go
package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/gradient"
	"github.com/thmalt/colors/space"
)

func main() {
	blue := colors.Oklab(0.42, -0.02, -0.15)
	purple := colors.Oklab(0.52, 0.08, -0.14)
	orange := colors.Oklab(0.72, 0.12, 0.12)

	// create gradient
	g, err := colors.NewGradient(
		space.Oklab, true, colors.HueShorter,

		colors.NewStop(0.00, blue),
		colors.NewStop(0.25, blue.Mix(purple, 0.5)),
		colors.NewStop(0.50, purple),
		colors.NewStop(0.75, purple.Mix(orange, 0.5)),
		colors.NewStop(1.00, orange),
	)

	if err != nil {
		log.Fatal(err)
	}

	width, height := 512, 288

	// create linear gradient
	l := gradient.NewLinear(g, 135, float64(width), float64(height))

	// render image
	img := image.NewRGBA64(image.Rect(0, 0, width, height))

	for x := 0; x < int(width); x++ {
		for y := 0; y < int(height); y++ {
			c := l.At(float64(x), float64(y))

			img.SetRGBA64(x, y, c.ToRGBA64())
		}
	}

	// save output.png
	f, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
```

<p align="center">
  <img src="./examples/linear-gradient/output.png" alt="Linear Gradient">
</p>