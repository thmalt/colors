# colors

Modern Go library for color space conversions and CSS color serialization.

### Installation

```go
import "github.com/thmalt/colors"
```

### Example

```go
package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

func main() {
	r, g, b := convert.RgbToSrgb(50, 60, 70)

	c := colors.Srgb(r, g, b)
	c = c.WithAlpha(0.999)
	for i := range space.Oklch {
		to := c.MustTo(i)
		fmt.Println(to)
		fmt.Println()
	}

	fmt.Println(colors.Srgb(r, g, b))
}

```