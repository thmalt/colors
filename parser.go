package colors

import (
	"errors"
	"math"
	"unsafe"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
	"github.com/thmalt/colors/internal/css/tokenizer"
	"github.com/thmalt/colors/space"
)

var (
	errInvalidColorToken = errors.New("invalid color token")
	errInvalidHexColor   = errors.New("invalid hex color")

	errInvalidAlpha          = errors.New("invalid alpha")
	errInvalidColorComponent = errors.New("invalid color component")

	errParserMissingCloseParenthesis = errors.New("missing closing parenthesis")

	errUnknownNamedColor = errors.New("unknown named color")
	errUnknownColorSpace = errors.New("unknown color space")

	errNotImplemented = errors.New("not implemented")
)

// Parse parses a CSS color value from s.
func Parse(s string) (Color, error) {
	var p parser
	p.init(s)
	p.t.Next()
	return p.parseColor()
}

// ParseBytes parses a CSS color value from b.
func ParseBytes(b []byte) (Color, error) {
	var p parser
	p.initBytes(b)
	p.t.Next()
	return p.parseColor()
}

type parser struct {
	t tokenizer.Tokenizer
}

func (p *parser) init(text string) {
	p.t.Init(text)
}

func (p *parser) initBytes(b []byte) {
	p.t.Init(unsafe.String(unsafe.SliceData(b), len(b)))
}

// main
func (p *parser) parseColor() (Color, error) {
	switch p.t.Kind() {
	case token.Hash:
		return p.parseHexColor()

	case token.Ident:
		return p.parseNamedColor()

	case token.Function:
		return p.parseColorFunction()

	default:
		return Color{}, errInvalidColorToken
	}
}

func (p *parser) parseHexColor() (Color, error) {
	c, ok := TryHex(p.t.Text())
	p.t.Next()
	if !ok {
		return Color{}, errInvalidHexColor
	}
	return c, nil
}

func (p *parser) parseNamedColor() (Color, error) {
	if c, ok := lookupNamedColor(p.t.Text()); ok {
		p.t.Next()
		return c, nil
	}

	return Color{}, errUnknownNamedColor
}

func (p *parser) parseColorFunction() (Color, error) {
	var buf [16]byte
	n := copy(buf[:], p.t.Text())
	ascii.ToLowerBytes(buf[:n])

	p.t.Next()

	originColor, err := p.parseOriginColor()
	if err != nil {
		return Color{}, err
	}

	switch string(buf[:n]) {
	case "rgb", "rgba":
		return p.parseRgb(originColor)
	case "hsl", "hsla":
		return p.parseHxxLike(originColor, space.Hsl, space.ChIdentFlagsHSL, true)
	case "hsv":
		return p.parseHxxLike(originColor, space.Hsv, space.ChIdentFlagsHSV, false)
	case "hwb":
		return p.parseHxxLike(originColor, space.Hwb, space.ChIdentFlagsHWB, false)

	case "luv":
		return p.parseLabLike(originColor, space.Luv, 100, 215)
	case "lab":
		return p.parseLabLike(originColor, space.Lab, 100, 125)
	case "oklab":
		return p.parseLabLike(originColor, space.Oklab, 1, .4)

	case "lchuv":
		return p.parseLchLike(originColor, space.Lchuv, 100, 220)
	case "lch":
		return p.parseLchLike(originColor, space.Lch, 100, 150)
	case "oklch":
		return p.parseLchLike(originColor, space.Oklch, 1, .4)

	case "color":
		return p.parseColorSpace(originColor)
	case "alpha":
		return p.parseRelativeAlpha(originColor)

	case "color-mix":
		return p.parseColorMix()
	}

	return Color{}, errInvalidColorToken
}

func (p *parser) parseOriginColor() (Color, error) {
	if !p.t.ExpectIdent("from") {
		return Color{}, nil
	}
	p.t.Next()
	return p.parseColor()
}

// rgb( <percentage>#{3} , <alpha-value>? ) | rgb( <number>#{3} , <alpha-value>? )

func (p *parser) parseRgb(originColor Color) (Color, error) {
	ctx := &colorCompCtx{}
	isLegacy := true

	if originColor.IsValid() {
		const sp = space.Srgb

		originColor.mutTo(sp)
		ctx.setSpace(sp)
		ctx.setChannel3(
			originColor.c1*255,
			originColor.c2*255,
			originColor.c3*255,
			originColor.alpha,
		)

		isLegacy = false
	}

	red, ok := p.parseNumberOrPercentage(ctx, true, 255)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	var green, blue, alpha parserNumericValue

	isLegacy = isLegacy && !red.isNone && p.t.NextExpect(token.Comma)
	if isLegacy {
		if red.kind.has(parserKindPercentage) {
			green, ok = p.parsePercentage(ctx, false, 255)
			if !ok {
				return Color{}, errInvalidColorComponent
			}

			if !p.t.NextExpect(token.Comma) {
				return Color{}, errInvalidColorComponent
			}

			blue, ok = p.parsePercentage(ctx, false, 255)
			if !ok {
				return Color{}, errInvalidColorComponent
			}
		} else {
			green, ok = p.parseNumber(ctx, false)
			if !ok {
				return Color{}, errInvalidColorComponent
			}

			if !p.t.NextExpect(token.Comma) {
				return Color{}, errInvalidColorComponent
			}

			blue, ok = p.parseNumber(ctx, false)
			if !ok {
				return Color{}, errInvalidColorComponent
			}
		}

		alpha, ok = p.parseLegacyAlpha(ctx)
		if !ok {
			return Color{}, errInvalidColorComponent
		}
	} else {
		green, ok = p.parseNumberOrPercentage(ctx, true, 255)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		blue, ok = p.parseNumberOrPercentage(ctx, true, 255)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		alpha, ok = p.parseModernAlpha(ctx)
		if !ok {
			return Color{}, errInvalidAlpha
		}
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	return SrgbAlpha(
		red.value*inv255,
		green.value*inv255,
		blue.value*inv255,
		clamp01(alpha.value),
	), nil
}

// legacy ( <hue>, <percentage>, <percentage>, <alpha-value>? )
// modern ( [<hue> | none] [<percentage> | <number> | none] [<percentage> | <number> | none] [ / [<alpha-value> | none] ]? )

func (p *parser) parseHxxLike(originColor Color, sp space.Space, chIdFlags space.ChIdentFlags, allowLegacy bool) (Color, error) {
	ctx := &colorCompCtx{}
	isLegacy := allowLegacy

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chFlags = chIdFlags
		ctx.setChannel3(
			originColor.c1,
			originColor.c2*100,
			originColor.c3*100,
			originColor.alpha,
		)

		isLegacy = false
	}

	hue, ok := p.parseNumberOrAngle(ctx, true)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	var c2, c3, alpha parserNumericValue

	isLegacy = isLegacy && !hue.isNone && p.t.NextExpect(token.Comma)
	if isLegacy {
		c2, ok = p.parsePercentage(ctx, false, 100)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		if !p.t.NextExpect(token.Comma) {
			return Color{}, errInvalidColorComponent
		}

		c3, ok = p.parsePercentage(ctx, false, 100)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		alpha, ok = p.parseLegacyAlpha(ctx)
		if !ok {
			return Color{}, errInvalidAlpha
		}
	} else {
		c2, ok = p.parseNumberOrPercentage(ctx, true, 100)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		c3, ok = p.parseNumberOrPercentage(ctx, true, 100)
		if !ok {
			return Color{}, errInvalidColorComponent
		}

		alpha, ok = p.parseModernAlpha(ctx)
		if !ok {
			return Color{}, errInvalidAlpha
		}
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	const inv = 1 / 100.0
	return new3(
		sp,
		wrap360(hue.value),
		c2.value*inv,
		c3.value*inv,
		clamp01(alpha.value),
	), nil
}

/*
lab([from <color>]?
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [ / [<alpha-value> | none] ]? )
*/

func (p *parser) parseLabLike(originColor Color, sp space.Space, lightnessRange, abRange float64) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chFlags = space.ChIdentFlagsLAB
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	lightness, ok := p.parseNumberOrPercentage(ctx, true, lightnessRange)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	a, ok := p.parseNumberOrPercentage(ctx, true, abRange)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	b, ok := p.parseNumberOrPercentage(ctx, true, abRange)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		return Color{}, errInvalidAlpha
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	return new3(
		sp,
		lightness.value,
		a.value,
		b.value,
		clamp01(alpha.value),
	), nil
}

/*
lch([from <color>]?
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [<hue> | none]
        [ / [<alpha-value> | none] ]? )
*/

func (p *parser) parseLchLike(originColor Color, sp space.Space, lightnessRange, chromaRange float64) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chFlags = space.ChIdentFlagsLCH
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	lightness, ok := p.parseNumberOrPercentage(ctx, true, lightnessRange)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	chroma, ok := p.parseNumberOrPercentage(ctx, true, chromaRange)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	hue, ok := p.parseNumberOrAngle(ctx, true)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		return Color{}, errInvalidAlpha
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	return new3(
		sp,
		lightness.value,
		chroma.value,
		wrap360(hue.value),
		clamp01(alpha.value),
	), nil
}

// alpha([from <color>] [ / [<alpha-value> | none] ]? )

func (p *parser) parseRelativeAlpha(originColor Color) (Color, error) {
	if !originColor.IsValid() {
		return Color{}, errInvalidAlpha
	}

	idx := uint8(originColor.space.ChIndex(space.ChIdentAlpha))
	if idx > space.MaxChannelCount {
		return Color{}, errInvalidAlpha
	}

	ctx := &colorCompCtx{
		space:   originColor.space,
		chFlags: space.ChIdentAlpha,
	}
	ctx.channels[idx] = originColor.alpha

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		return Color{}, errInvalidAlpha
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	originColor.alpha = clamp01(alpha.value)
	return originColor, nil
}

func (p *parser) parseColorMix() (Color, error) {
	return Color{}, errNotImplemented
}

func (p *parser) parseColorSpace01(originColor Color, sp space.Space, chIdFlags space.ChIdentFlags) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chFlags = chIdFlags
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	c1, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	c2, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	c3, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		return Color{}, errInvalidColorComponent
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		return Color{}, errInvalidAlpha
	}

	if !p.t.NextExpect(token.CloseParenthesis) {
		return Color{}, errParserMissingCloseParenthesis
	}

	return new3(
		sp,
		c1.value,
		c2.value,
		c3.value,
		clamp01(alpha.value),
	), nil
}

// TODO: Implement Math functions.
func (p *parser) parseColorComponent(ctx *colorCompCtx, allowNone bool, allowKind parserNumericKind, percentScale float64) (parserNumericValue, bool) {
	switch p.t.Kind() {
	case token.Ident:
		if allowNone && p.t.ExpectIdent("none") {
			p.t.Next()
			return parserNumericValue{isNone: true}, true
		} else if ctx.chFlags != 0 {
			if chid := space.LookupChIdentFlags(p.t.Text()); ctx.chFlags.Has(chid) {
				p.t.Next()
				return parserNumericValue{value: ctx.channel(chid)}, true
			}
		}

	case token.Function:
		return parserNumericValue{}, false

	case token.Number:
		if allowKind.has(parserKindNumber) {
			n := p.t.Number()
			p.t.Next()
			return parserNumericValue{value: n, kind: parserKindNumber}, true
		}

	case token.Percentage:
		if allowKind.has(parserKindPercentage) {
			n := p.t.Number()
			p.t.Next()
			return parserNumericValue{value: n * percentScale, kind: parserKindPercentage}, true
		}

	case token.Dimension:
		if allowKind.has(parserKindAngle) {
			var buf [8]byte
			n := copy(buf[:], p.t.Unit())
			ascii.ToLowerBytes(buf[:n])

			value := p.t.Number()

			switch string(buf[:n]) {
			case "deg":
			case "grad":
				value *= 0.9
			case "rad":
				value *= 180 / math.Pi
			case "turn":
				value *= 360
			default:
				return parserNumericValue{}, false
			}

			p.t.Next()
			return parserNumericValue{value: value, kind: parserKindAngle}, true
		}
	}

	return parserNumericValue{}, false
}

func (p *parser) parseNumberOrAngle(ctx *colorCompCtx, allowNone bool) (parserNumericValue, bool) {
	return p.parseColorComponent(ctx, allowNone, parserKindNumber|parserKindAngle, 360)
}

func (p *parser) parseNumberOrPercentage(ctx *colorCompCtx, allowNone bool, percentScale float64) (parserNumericValue, bool) {
	return p.parseColorComponent(ctx, allowNone, parserKindNumber|parserKindPercentage, percentScale)
}

func (p *parser) parseNumber(ctx *colorCompCtx, allowNone bool) (parserNumericValue, bool) {
	return p.parseColorComponent(ctx, allowNone, parserKindNumber, 0)
}

func (p *parser) parsePercentage(ctx *colorCompCtx, allowNone bool, percentScale float64) (parserNumericValue, bool) {
	return p.parseColorComponent(ctx, allowNone, parserKindPercentage, percentScale)
}

func (p *parser) parseLegacyAlpha(ctx *colorCompCtx) (parserNumericValue, bool) {
	if p.t.Is(token.CloseParenthesis) {
		return parserNumericValue{value: 1}, true
	}

	if !p.t.NextExpect(token.Comma) {
		return parserNumericValue{}, false
	}

	return p.parseNumberOrPercentage(ctx, false, 1)
}

func (p *parser) parseModernAlpha(ctx *colorCompCtx) (parserNumericValue, bool) {
	if p.t.ExpectCloseParenthesis() {
		if ctx.chFlags.Has(space.ChIdentAlpha) {
			return parserNumericValue{value: ctx.channel(space.ChIdentAlpha)}, true
		}
		return parserNumericValue{value: 1}, true
	}

	if !p.t.ExpectDelim('/') {
		return parserNumericValue{}, false
	}

	p.t.Next()

	value, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if value.isNone && ctx.chFlags.Has(space.ChIdentAlpha) {
		value.value = ctx.channel(space.ChIdentAlpha)
	}
	return value, ok
}
