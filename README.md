# colors

Modern Go library for color space conversions and CSS color serialization.

### Usage

```go
import "github.com/thmalt/colors"
```

### Example

```go
package main

import (
	"fmt"

	"github.com/thmalt/colors"
)

func main() {
	c := colors.Rgb(50, 60, 70).WithAlpha(0.995)

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