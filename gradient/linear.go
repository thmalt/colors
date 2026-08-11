package gradient

import (
	"math"

	"github.com/thmalt/colors"
)

type Linear struct {
	gradient colors.Gradient

	angle float64

	width  float64
	height float64

	ax, ay float64 // gradient direction coefficients
	bias   float64 // gradient offset
}

func NewLinear(gradient colors.Gradient, angle float64, width, height float64) *Linear {
	l := &Linear{
		gradient: gradient,
		angle:    angle,
		width:    width,
		height:   height,
	}

	l.updateTransform()

	return l
}

func (l *Linear) SetBounds(width, height float64) {
	l.width = width
	l.height = height
	l.updateTransform()
}

func (l *Linear) At(x, y float64) colors.Color {
	return l.gradient.At(l.ax*x + l.ay*y + l.bias)
}

func (l *Linear) updateTransform() {
	rad := l.angle * math.Pi / 180
	sin, cos := math.Sincos(rad)

	dx, dy := sin, -cos
	hw, hh := l.width/2, l.height/2

	L := (math.Abs(dx)*hw + math.Abs(dy)*hh) * 2

	l.ax = dx / L
	l.ay = dy / L
	l.bias = 0.5 - hw*l.ax - hh*l.ay
}
