package codegen

import (
	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
)

var (
	Spaces = [...]model.Space{
		{
			Name:        "LinearSrgb",
			Base:        "XyzD65",
			DisplayName: "Linear sRGB",
			CssName:     "srgb-linear",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
		},
		{
			Name:        "Srgb",
			Base:        "LinearSrgb",
			DisplayName: "sRGB",
			CssName:     "srgb",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
		},
		{
			Name:        "LinearA98",
			Base:        "XyzD65",
			DisplayName: "Linear Adobe RGB (1998)",
			CssName:     "a98-rgb-linear",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "A98",
			Base:        "LinearA98",
			DisplayName: "Adobe RGB (1998)",
			CssName:     "a98-rgb",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "LinearDisplayP3",
			Base:        "XyzD65",
			DisplayName: "Linear Display P3",
			CssName:     "display-p3-linear",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "DisplayP3",
			Base:        "LinearDisplayP3",
			DisplayName: "Display P3",
			CssName:     "display-p3",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "LinearProPhoto",
			Base:        "XyzD50",
			DisplayName: "Linear ProPhoto",
			CssName:     "prophoto-rgb-linear",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "ProPhoto",
			Base:        "LinearProPhoto",
			DisplayName: "ProPhoto",
			CssName:     "prophoto-rgb",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "LinearRec2020",
			Base:        "XyzD65",
			DisplayName: "Linear Rec. 2020",
			CssName:     "rec2020-linear",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "Rec2020",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020",
			CssName:     "rec2020",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "red", Symbol: "r", DisplayName: "Red", Min: 0, Max: 1, Precision: 6},
				{Ident: "green", Symbol: "g", DisplayName: "Green", Min: 0, Max: 1, Precision: 6},
				{Ident: "blue", Symbol: "b", DisplayName: "Blue", Min: 0, Max: 1, Precision: 6},
			},
			UseColorFunction: true,
			// Disable: true,
		},
		{
			Name:        "Hsl",
			Base:        "Srgb",
			DisplayName: "HSL",
			CssName:     "hsl",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "hue", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 2},
				{Ident: "sat", Symbol: "s", DisplayName: "Saturation", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
				{Ident: "light", Symbol: "l", DisplayName: "Lightness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
			},
		},
		{
			Name:        "Hsv",
			Base:        "Srgb",
			DisplayName: "HSV",
			CssName:     "hsv",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "hue", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 2},
				{Ident: "sat", Symbol: "s", DisplayName: "Saturation", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
				{Ident: "val", Symbol: "v", DisplayName: "Value", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
			},
			UseColorFunction: true,
		},
		{
			Name:        "Hwb",
			Base:        "Srgb",
			DisplayName: "HWB",
			CssName:     "hwb",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "hue", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 2},
				{Ident: "white", Symbol: "w", DisplayName: "Whiteness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
				{Ident: "black", Symbol: "b", DisplayName: "Blackness", Min: 0, Max: 1, Unit: model.UnitPercent, Precision: 2},
			},
		},
		{
			Name:        "XyzD65",
			DisplayName: "CIE XYZ D65",
			CssName:     "xyz-d65",
			// Aliases:     []string{"xyz"},
			WhitePoint: "D65",
			Channels: []model.Channel{
				{Ident: "x", Symbol: "x", DisplayName: "X", Min: 0, Max: 1, Precision: 8},
				{Ident: "y", Symbol: "y", DisplayName: "Y", Min: 0, Max: 1, Precision: 8},
				{Ident: "z", Symbol: "z", DisplayName: "Z", Min: 0, Max: 1, Precision: 8},
			},
			UseColorFunction: true,
		},
		{
			Name:        "XyzD50",
			Base:        "XyzD65",
			DisplayName: "CIE XYZ D50",
			CssName:     "xyz-d50",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Ident: "x", Symbol: "x", DisplayName: "X", Min: 0, Max: 1, Precision: 8},
				{Ident: "y", Symbol: "y", DisplayName: "Y", Min: 0, Max: 1, Precision: 8},
				{Ident: "z", Symbol: "z", DisplayName: "Z", Min: 0, Max: 1, Precision: 8},
			},
			UseColorFunction: true,
		},
		{
			Name:        "Lab",
			Base:        "XyzD50",
			DisplayName: "CIE Lab",
			CssName:     "lab",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Ident: "light", Symbol: "l", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Ident: "a", Symbol: "a", DisplayName: "Green-Red", Min: -125, Max: 125, Unrestricted: true, Precision: 4},
				{Ident: "b", Symbol: "b", DisplayName: "Blue-Yellow", Min: -125, Max: 125, Unrestricted: true, Precision: 4},
			},
		},
		{
			Name:        "Lch",
			Base:        "Lab",
			DisplayName: "CIE LCh",
			CssName:     "lch",
			WhitePoint:  "D50",
			Channels: []model.Channel{
				{Ident: "light", Symbol: "l", DisplayName: "Lightness", Min: 0, Max: 100, Precision: 4},
				{Ident: "chroma", Symbol: "c", DisplayName: "Chroma", Min: 0, Max: 150, Unrestricted: true, Precision: 4},
				{Ident: "hue", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 4},
			},
		},
		{
			Name:        "Oklab",
			Base:        "XyzD65",
			DisplayName: "Oklab",
			CssName:     "oklab",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "light", Symbol: "l", DisplayName: "Lightness", Min: 0, Max: 1, Precision: 6},
				{Ident: "a", Symbol: "a", DisplayName: "Green-Red", Min: -0.4, Max: 0.4, Unrestricted: true, Precision: 6},
				{Ident: "b", Symbol: "b", DisplayName: "Blue-Yellow", Min: -0.4, Max: 0.4, Unrestricted: true, Precision: 6},
			},
		},
		{
			Name:        "Oklch",
			Base:        "Oklab",
			DisplayName: "Oklch",
			CssName:     "oklch",
			WhitePoint:  "D65",
			Channels: []model.Channel{
				{Ident: "light", Symbol: "l", DisplayName: "Lightness", Min: 0, Max: 1, Precision: 6},
				{Ident: "chroma", Symbol: "c", DisplayName: "Chroma", Min: 0, Max: 0.4, Unrestricted: true, Precision: 6},
				{Ident: "hue", Symbol: "h", DisplayName: "Hue", Min: 0, Max: 360, Circular: true, Unit: model.UnitDegree, Precision: 4},
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
		{Pair: Pair{"Srgb", "Linear"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Linear", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"A98", "Linear"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Linear", "A98"}, Cost: 1, Implemented: true},
		{Pair: Pair{"DisplayP3", "Linear"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Linear", "DisplayP3"}, Cost: 1, Implemented: true},
		{Pair: Pair{"ProPhoto", "Linear"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Linear", "ProPhoto"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Rec2020", "Linear"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Linear", "Rec2020"}, Cost: 1, Implemented: true},
		// standard converter
		{Pair: Pair{"Srgb", "Hsl"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hsl", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Srgb", "Hsv"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hsv", "Srgb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Srgb", "Hwb"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Hwb", "Srgb"}, Cost: 1, Implemented: true},
		// standard converter
		{Pair: Pair{"Lab", "Lch"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Lch", "Lab"}, Cost: 1, Implemented: true},
		{Pair: Pair{"Lab", "XyzD50"}, Cost: 1, Implemented: true},
		{Pair: Pair{"XyzD50", "Lab"}, Cost: 1, Implemented: true},
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
		{
			Pair: Pair{"Oklab", "Oklch"},
			Cost: 2,
			Ops: []Op{
				{Type: OpCall, Pair: Pair{"Lab", "Lch"}},
			},
		},
		{
			Pair: Pair{"Oklch", "Oklab"},
			Cost: 2,
			Ops: []Op{
				{Type: OpCall, Pair: Pair{"Lch", "Lab"}},
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
