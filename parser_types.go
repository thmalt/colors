package colors

import "github.com/thmalt/colors/space"

type colorCompCtx struct {
	channels [space.MaxChannelCount + 1]float64
	space    space.Space
	chFlags  space.ChIdentFlags
}

func (ctx *colorCompCtx) setSpace(s space.Space) {
	ctx.space = s
	ctx.chFlags = s.ChIdentFlags()
}

func (ctx *colorCompCtx) setChannel3(c1, c2, c3, alpha float64) {
	ctx.channels[0] = c1
	ctx.channels[1] = c2
	ctx.channels[2] = c3
	ctx.channels[3] = alpha
}

func (ctx *colorCompCtx) setChannel4(sp space.Space, c1, c2, c3, c4, alpha float64) {
	ctx.space = sp
	ctx.channels[0] = c1
	ctx.channels[1] = c2
	ctx.channels[2] = c3
	ctx.channels[3] = c4
	ctx.channels[4] = alpha
}

func (ctx *colorCompCtx) channel(ch space.ChIdentFlags) float64 {
	idx := uint8(ctx.space.ChIndex(ch))
	if idx >= uint8(len(ctx.channels)) {
		return 0
	}
	return ctx.channels[idx]
}

type parserNumericKind uint8

const (
	parserKindNumber parserNumericKind = 1 << iota
	parserKindPercentage
	parserKindAngle
)

func (nk parserNumericKind) has(k parserNumericKind) bool {
	return nk&k == k
}

type parserNumericValue struct {
	value  float64
	kind   parserNumericKind
	unit   parserNumericUnit
	isNone bool
}

type parserNumericUnit uint8

const (
	_ parserNumericUnit = iota
)

func new3(sp space.Space, c1, c2, c3, alpha float64) Color {
	return Color{
		space: sp,
		c1:    c1,
		c2:    c2,
		c3:    c3,
		alpha: alpha,
	}
}

func new4(sp space.Space, c1, c2, c3, c4, alpha float64) Color {
	return Color{
		space: sp,
		c1:    c1,
		c2:    c2,
		c3:    c3,
		c4:    c4,
		alpha: alpha,
	}
}
