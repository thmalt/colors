package codegen

import (
	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
)

var (
	// To arrange color space.
	FamilyOrder = [...]string{
		"XYZ",
		"RGB",
		"Lab",
		"Luv",
		"Oklab",
	}

	HueInterpolation = [...]string{
		"HueShorter",
		"HueLonger",
		"HueIncreasing",
		"HueDecreasing",
	}

	rgbChannels = []model.Channel{
		{Name: "Red", Ident: "r", Symbol: "R", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
		{Name: "Green", Ident: "g", Symbol: "G", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
		{Name: "Blue", Ident: "b", Symbol: "B", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
	}

	hueChannel        = model.Channel{Name: "Hue", Ident: "h", Symbol: "H", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 2}
	saturationChannel = model.Channel{Name: "Saturation", Ident: "s", Symbol: "S", DisplayName: "Saturation", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 4}
	lightnessChannel  = model.Channel{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 4}
	valueChannel      = model.Channel{Name: "Value", Ident: "v", Symbol: "V", DisplayName: "Value", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 4}

	whitenessChannel = model.Channel{Name: "Whiteness", Ident: "w", Symbol: "W", DisplayName: "Whiteness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 4}
	blacknessChannel = model.Channel{Name: "Blackness", Ident: "b", Symbol: "B", DisplayName: "Blackness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 4}

	xyzChannels = []model.Channel{
		{Name: "X", Ident: "x", Symbol: "X", DisplayName: "X", Min: 0, Max: 1, Precision: 8},
		{Name: "Y", Ident: "y", Symbol: "Y", DisplayName: "Y", Min: 0, Max: 1, Precision: 8},
		{Name: "Z", Ident: "z", Symbol: "Z", DisplayName: "Z", Min: 0, Max: 1, Precision: 8},
	}

	Spaces = [...]model.Space{
		{
			Name:             "LinearSrgb",
			Family:           "RGB",
			Base:             "XyzD65",
			DisplayName:      "Linear sRGB",
			CssName:          "srgb-linear",
			WhitePoint:       "D65",
			Channels:         rgbChannels,
			UseColorFunction: true,
		},
		{
			Name:        "Srgb",
			Family:      "RGB",
			Base:        "LinearSrgb",
			DisplayName: "sRGB",
			CssName:     "srgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
		},
		{
			Name:        "LinearA98",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Adobe RGB (1998)",
			CssName:     "a98-rgb-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "A98",
			Family:      "RGB",
			Base:        "LinearA98",
			DisplayName: "Adobe RGB (1998)",
			CssName:     "a98-rgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "LinearDisplayP3",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Display P3",
			CssName:     "display-p3-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "DisplayP3",
			Family:      "RGB",
			Base:        "LinearDisplayP3",
			DisplayName: "Display P3",
			CssName:     "display-p3",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "LinearProPhoto",
			Family:      "RGB",
			Base:        "XyzD50",
			DisplayName: "Linear ProPhoto",
			CssName:     "prophoto-rgb-linear",
			WhitePoint:  "D50",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
			SnakeName: "linear_prophoto",
		},
		{
			Name:        "ProPhoto",
			Family:      "RGB",
			Base:        "LinearProPhoto",
			DisplayName: "ProPhoto",
			CssName:     "prophoto-rgb",
			WhitePoint:  "D50",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,

			SnakeName: "prophoto",
		},
		{
			Name:        "LinearRec2020",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Rec. 2020",
			CssName:     "rec2020-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "Rec2020",
			Family:      "RGB",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020",
			CssName:     "rec2020",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "Hsl",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSL",
			CssName:     "hsl",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    []model.Channel{hueChannel, saturationChannel, lightnessChannel},
		},
		{
			Name:        "Hsv",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSV",
			CssName:     "hsv",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    []model.Channel{hueChannel, saturationChannel, valueChannel},
		},
		{
			Name:        "Hwb",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HWB",
			CssName:     "hwb",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    []model.Channel{hueChannel, whitenessChannel, blacknessChannel},
		},
		{
			Name:        "XyzD65",
			Family:      "XYZ",
			DisplayName: "CIE XYZ D65",
			CssName:     "xyz-d65",
			// Aliases:     []string{"xyz"},
			WhitePoint:       "D65",
			Channels:         xyzChannels,
			UseColorFunction: true,
		},
		{
			Name:             "XyzD50",
			Family:           "XYZ",
			Base:             "XyzD65",
			DisplayName:      "CIE XYZ D50",
			CssName:          "xyz-d50",
			WhitePoint:       "D50",
			Channels:         xyzChannels,
			UseColorFunction: true,
		},
		{
			Name:        "XyY",
			Family:      "XYZ",
			Base:        "XyzD65",
			DisplayName: "CIE xyY",
			CssName:     "xyY",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Name: "Chromaticity x", Ident: "x", Symbol: "x", DisplayName: "x", Min: 0, Max: 1, Precision: 8},
				{Name: "Chromaticity y", Ident: "y", Symbol: "y", DisplayName: "y", Min: 0, Max: 1, Precision: 8},
				{Name: "Luminance", Ident: "luminance", Symbol: "Y", DisplayName: "Y", Min: 0, Max: 1, Unrestricted: true, Precision: 8},
			},
			UseColorFunction: true,
			SnakeName:        "xyy",
			Comment: "XyY is the CIE xyY color space using the D65 reference white.\n" +
				"Conversions involving other reference whites automatically perform\n" +
				"chromatic adaptation.",
		},
		{
			Name:        "Lab",
			Family:      "Lab",
			Base:        "XyzD50",
			DisplayName: "CIE Lab",
			CssName:     "lab",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Name: "A", Ident: "a", Symbol: "a", DisplayName: "Green-Red", Min: -125, Max: 125, Unrestricted: true, Precision: 4},
				{Name: "B", Ident: "b", Symbol: "b", DisplayName: "Blue-Yellow", Min: -125, Max: 125, Unrestricted: true, Precision: 4},
			},
		},
		{
			Name:        "Lch",
			Family:      "Lab",
			Base:        "Lab",
			DisplayName: "CIE LCh",
			CssName:     "lch",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Name: "Chroma", Ident: "c", Symbol: "C", DisplayName: "Chroma", Min: 0, Max: 150, Unrestricted: true, Precision: 4},
				{Name: "Hue", Ident: "h", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 4},
			},
		},
		{
			Name:        "Luv",
			Family:      "Luv",
			Base:        "XyzD50",
			DisplayName: "CIE Luv",
			CssName:     "luv",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Name: "U", Ident: "u", Symbol: "u", DisplayName: "Green-Red Opponent", Min: -134, Max: 220, Unrestricted: true, Precision: 4},
				{Name: "V", Ident: "v", Symbol: "v", DisplayName: "Blue-Yellow Opponent", Min: -140, Max: 122, Unrestricted: true, Precision: 4},
			},
		},
		{
			Name:        "Lchuv",
			Family:      "Luv",
			Base:        "Luv",
			DisplayName: "CIE LChuv",
			CssName:     "lchuv",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Name: "Chroma", Ident: "c", Symbol: "C", DisplayName: "Chroma", Min: 0, Max: 180, Unrestricted: true, Precision: 4},
				{Name: "Hue", Ident: "h", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 4},
			},
		},
		{
			Name:        "Oklab",
			Family:      "Oklab",
			Base:        "XyzD65",
			DisplayName: "Oklab",
			CssName:     "oklab",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 1, Precision: 6},
				{Name: "A", Ident: "a", Symbol: "a", DisplayName: "Green-Red Opponent", Min: -0.4, Max: 0.4, Unrestricted: true, Precision: 6},
				{Name: "B", Ident: "b", Symbol: "b", DisplayName: "Blue-Yellow Opponent", Min: -0.4, Max: 0.4, Unrestricted: true, Precision: 6},
			},
		},
		{
			Name:        "Oklch",
			Family:      "Oklab",
			Base:        "Oklab",
			DisplayName: "Oklch",
			CssName:     "oklch",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels: []model.Channel{
				{Name: "Lightness", Ident: "l", Symbol: "L", DisplayName: "Lightness", Min: 0, Max: 1, Precision: 6},
				{Name: "Chroma", Ident: "c", Symbol: "C", DisplayName: "Chroma", Min: 0, Max: 0.4, Unrestricted: true, Precision: 6},
				{Name: "Hue", Ident: "h", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 4},
			},
		},
	}

	WhitePoints = [...]model.WhitePoint{
		{
			Name: "D50",
			X:    data.D50.X,
			Y:    data.D50.Y,
		},
		{
			Name: "D65",
			X:    data.D65.X,
			Y:    data.D65.Y,
		},
	}

	ConvertFuncs = [...]ConvertFunc{
		// standard transfer converter
		{Pair: Pair{"Srgb", "LinearSrgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"LinearSrgb", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"A98", "LinearA98"}, Cost: 1, Implemented: true},
		{Pair: Pair{"LinearA98", "A98"}, Cost: 1, Implemented: true},
		{Pair: Pair{"DisplayP3", "LinearDisplayP3"}, Cost: 1, Implemented: true},
		{Pair: Pair{"LinearDisplayP3", "DisplayP3"}, Cost: 1, Implemented: true},
		{Pair: Pair{"ProPhoto", "LinearProPhoto"}, Cost: 1, Implemented: true},
		{Pair: Pair{"LinearProPhoto", "ProPhoto"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Rec2020", "LinearRec2020"}, Cost: 1, Implemented: true},
		{Pair: Pair{"LinearRec2020", "Rec2020"}, Cost: 1, Implemented: true},
		// standard converter
		{Pair: Pair{"Srgb", "Hsl"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hsl", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Srgb", "Hsv"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hsv", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Srgb", "Hwb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hwb", "Srgb"}, Cost: 1, Implemented: true},
		// standard converter
		{Pair: Pair{"XyY", "XyzD65"}, Cost: 1, Implemented: true},
		{Pair: Pair{"XyzD65", "XyY"}, Cost: 1, Implemented: true},
		//
		{Pair: Pair{"Lab", "Lch"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Lch", "Lab"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Lab", "XyzD50"}, Cost: 1, Implemented: true},
		{Pair: Pair{"XyzD50", "Lab"}, Cost: 1, Implemented: true},
		//
		{Pair: Pair{"Luv", "Lchuv"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Lchuv", "Luv"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Luv", "XyzD50"}, Cost: 1, Implemented: true},
		{Pair: Pair{"XyzD50", "Luv"}, Cost: 1, Implemented: true},
		//
		{Pair: Pair{"Oklab", "Oklch"}, Cost: 2, Implemented: true},
		{Pair: Pair{"Oklch", "Oklab"}, Cost: 2, Implemented: true},
		// generate with Ops
		{
			Pair: Pair{"Oklab", "XyzD65"},
			Cost: 3,
			Ops: []Op{
				{Type: OpMatrix, Matrix: &data.OklabLabToLms},
				{Type: OpCube},
				{Type: OpMatrix, Matrix: &data.OklabLmsToXyzD65},
			},
		},
		{
			Pair: Pair{"XyzD65", "Oklab"},
			Cost: 3,
			Ops: []Op{
				{Type: OpMatrix, Matrix: &data.OklabXyzD65ToLms},
				{Type: OpCbrt},
				{Type: OpMatrix, Matrix: &data.OklabLmsToLab},
			},
		},
		// Xyz
		{
			Pair: Pair{"XyzD65", "XyzD50"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToXyzD50}},
		},
		{
			Pair: Pair{"XyzD50", "XyzD65"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD50ToXyzD65}},
		},
		// Xyz* -> Linear*
		{
			Pair: Pair{"XyzD65", "LinearSrgb"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearSrgb}},
		},
		{
			Pair: Pair{"LinearSrgb", "XyzD65"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearSrgbToXyzD65}},
		},
		{
			Pair: Pair{"XyzD65", "LinearDisplayP3"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearDisplayP3}},
		},
		{
			Pair: Pair{"LinearDisplayP3", "XyzD65"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearDisplayP3ToXyzD65}},
		},
		{
			Pair: Pair{"XyzD65", "LinearA98"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearA98}},
		},
		{
			Pair: Pair{"LinearA98", "XyzD65"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearA98ToXyzD65}},
		},
		{
			Pair: Pair{"XyzD50", "LinearProPhoto"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD50ToLinearProPhoto}},
		},
		{
			Pair: Pair{"LinearProPhoto", "XyzD50"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearProPhotoToXyzD50}},
		},
		{
			Pair: Pair{"XyzD65", "LinearRec2020"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearRec2020}},
		},
		{
			Pair: Pair{"LinearRec2020", "XyzD65"},
			Cost: 1,
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearRec2020ToXyzD65}},
		},
	}
)

func LookupSpace(name string) *model.Space {
	for _, space := range Spaces {
		if name == space.Name {
			return &space
		}
	}
	return nil
}

func LookupWhitePoint(name string) *model.WhitePoint {
	for _, whitePoint := range WhitePoints {
		if name == whitePoint.Name {
			return &whitePoint
		}
	}
	return nil
}
