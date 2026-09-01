package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/gradient"
)

const width, height = 512, 256

var (
	blue   = colors.Oklab(0.42, -0.02, -0.15)
	purple = colors.Oklab(0.52, 0.08, -0.14)
	orange = colors.Oklab(0.72, 0.12, 0.12)

	// CSS:
	// linear-gradient(
	//   to bottom right in oklab,
	//   #1b469e,
	//   #504baa,
	//   #784db6,
	//   #b56c8a,
	//   #f87c3e
	// )
	linearGradient = colors.NewGradientBuilder().
			AddStop(blue).
			AddStop(blue.Mix(purple, 0.5), 0.25).
			AddStop(purple, 0.5).
			AddStop(purple.Mix(orange, 0.5), 0.75).
			AddStop(orange, 1).
			Build()

	// CSS:
	// conic-gradient(
	//   from 0.25turn at 50% 30% in oklab,
	//   #f69d3c,
	//   10deg,
	//   #3f87a6,
	//   350deg,
	//   #ebf8e1
	// )
	conicGradient = colors.NewGradientBuilder().
			AddStop(colors.Hex("#f69d3c")).
			AddHint(gradient.DegToTurn(10)).
			AddStop(colors.Hex("#3f87a6")).
			AddHint(gradient.DegToTurn(350)).
			AddStop(colors.Hex("#ebf8e1")).
			Build()

	// CSS:
	// radial-gradient(
	//   farthest-corner in oklab,
	//   #3f5efb,
	//   #00ff88,
	//   #ce436d,
	//   #0413bf 50%,
	//   #fc466b
	// )
	radialGradient = colors.NewGradientBuilder().
			AddStop(colors.Hex("#3f5efb")).
			AddStop(colors.Hex("#00ff88")).
			AddStop(colors.Hex("#ce436d")).
			AddStop(colors.Hex("#0413bf"), 0.5).
			AddStop(colors.Hex("#fc466b")).
			Build()
)

func main() {
	// CSS: linear-gradient(to bottom right, ...)
	linearSpec := gradient.NewLinearSpec()
	linearSpec.SetSize(width, height)
	linearSpec.SetDirection(gradient.ToBottomRight)

	linear, err := linearSpec.Build()
	if err != nil {
		log.Fatal(err)
	}

	// CSS: conic-gradient(from 0.25turn at 50% 30%, ...)
	conicSpec := gradient.NewConicSpec()
	conicSpec.SetSize(width, height)
	conicSpec.SetCenter(0.5, 0.3)
	conicSpec.SetStartAngle(gradient.ToRight)

	conic, err := conicSpec.Build()
	if err != nil {
		log.Fatal(err)
	}

	// CSS: radial-gradient(farthest-corner, ...)
	radialSpec := gradient.NewRadialSpec()
	radialSpec.SetSize(width, height)

	// CSS: radial-gradient(circle, ...)
	radialSpec.SetShape(gradient.RadialCircle)

	circle, err := radialSpec.Build()
	if err != nil {
		log.Fatal(err)
	}

	// CSS: radial-gradient(ellipse, ...)
	radialSpec.SetShape(gradient.RadialEllipse)

	ellipse, err := radialSpec.Build()
	if err != nil {
		log.Fatal(err)
	}

	// Render the examples.
	saveImage("output/linear.png", renderImage(width, height, func(x, y float64) colors.Color {
		return linearGradient.At(linear.PositionAt(x, y))
	}))

	saveImage("output/conic.png", renderImage(width, height, func(x, y float64) colors.Color {
		return conicGradient.At(conic.PositionAt(x, y))
	}))

	saveImage("output/radial-circle.png", renderImage(width, height, func(x, y float64) colors.Color {
		return radialGradient.At(circle.PositionAt(x, y))
	}))

	saveImage("output/radial-ellipse.png", renderImage(width, height, func(x, y float64) colors.Color {
		return radialGradient.At(ellipse.PositionAt(x, y))
	}))

}

func renderImage(width, height float64, fn func(x, y float64) colors.Color) image.Image {
	w, h := int(width), int(height)
	img := image.NewRGBA64(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			// Sample the gradient at pixel centers rather than pixel corners.
			c := fn(float64(x)+0.5, float64(y)+0.5)

			// Dither before quantizing to RGBA64 to reduce visible banding.
			c = c.Dither(x, y)

			img.SetRGBA64(x, y, c.ToRGBA64())
		}
	}

	return img
}

func saveImage(file string, img image.Image) {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
