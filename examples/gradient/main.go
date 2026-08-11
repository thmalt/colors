package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/gradient"
)

var (
	blue   = colors.Oklab(0.42, -0.02, -0.15)
	purple = colors.Oklab(0.52, 0.08, -0.14)
	orange = colors.Oklab(0.72, 0.12, 0.12)

	linearGradient = colors.NewGradient(
		colors.NewStopAt(blue, 0),
		colors.NewStopAt(blue.Mix(purple, 0.5), 0.25),
		colors.NewStopAt(purple, 0.5),
		colors.NewStopAt(purple.Mix(orange, 0.5), 0.75),
		colors.NewStopAt(orange, 1),
	)

	radialGradient = colors.NewGradient(
		colors.NewStopAt(colors.Rgb(63, 94, 251), 0),
		colors.NewStopAt(colors.Rgb(0, 255, 136), .2),
		colors.NewStopAt(colors.Rgb(206, 67, 109), .4),
		colors.NewStopAt(colors.Rgb(4, 19, 191), .6),
		colors.NewStopAt(colors.Rgb(252, 70, 107), 1),
	)

	conicGradient = colors.NewGradient(
		colors.NewStopAt(colors.Rgb(0xf6, 0x9d, 0x3c), 0),
		colors.NewHint(gradient.DegToTurn(10)),
		colors.NewStopAt(colors.Rgb(0x3f, 0x87, 0xa6), .5),
		colors.NewHint(gradient.DegToTurn(350)),
		colors.NewStopAt(colors.Rgb(0xeb, 0xf8, 0xe1), 1),
	)
)

const width, height = 512, 256

func main() {
	// Create geometric gradients.
	linear := gradient.NewLinear(width, height, gradient.DegToTurn(135))
	radial := gradient.NewRadial(width, height, 0.5, 0.5)
	conic := gradient.NewConic(width, height, 0.5, 0.3, 0.25)

	// Render linear gradient.
	linearImage := renderImage(width, height, func(x, y float64) colors.Color {
		t := linear.PositionAt(x, y)
		return linearGradient.At(t)
	})
	saveImage("output/linear.png", linearImage)

	// Render circular radial gradient.
	radial.SetShape(gradient.RadialCircle)
	circleRadialImage := renderImage(width, height, func(x, y float64) colors.Color {
		t := radial.PositionAt(x, y)
		return radialGradient.At(t)
	})
	saveImage("output/radial-circle.png", circleRadialImage)

	// Render elliptical radial gradient.
	radial.SetShape(gradient.RadialEllipse)
	ellipseRadialImage := renderImage(width, height, func(x, y float64) colors.Color {
		t := radial.PositionAt(x, y)
		return radialGradient.At(t)
	})
	saveImage("output/radial-ellipse.png", ellipseRadialImage)

	// Render conic gradient.
	conicImage := renderImage(width, height, func(x, y float64) colors.Color {
		t := conic.PositionAt(x, y)
		return conicGradient.At(t)
	})
	saveImage("output/conic.png", conicImage)
}

func renderImage(width, height float64, fn func(x, y float64) colors.Color) *image.RGBA64 {
	w, h := int(width), int(height)
	img := image.NewRGBA64(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Sample at pixel centers to avoid sampling at pixel corners. +0.5
			c := fn(float64(x)+0.5, float64(y)+0.5)

			img.SetRGBA64(x, y, c.ToRGBA64())
		}
	}

	return img
}

func saveImage(file string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}
