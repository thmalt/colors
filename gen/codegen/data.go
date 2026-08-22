package codegen

import (
	"math"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
)

var (
	HueInterpolation = [...]string{
		"HueShorter",
		"HueLonger",
		"HueIncreasing",
		"HueDecreasing",
	}

	xChannelName = channelName("X", "x", "X", "X")
	yChannelName = channelName("Y", "y", "Y", "Y")
	zChannelName = channelName("Z", "z", "Z", "Z")

	lightnessChannelName = channelName("Lightness", "l", "L", "Lightness")
	chromaChannelName    = channelName("Chroma", "c", "C", "Chroma")
	hueChannelName       = channelName("Hue", "h", "h", "Hue")

	uChannelName = channelName("U", "u", "u", "Green-Red Opponent")
	vChannelName = channelName("V", "v", "v", "Blue-Yellow Opponent")

	aChannelName = channelName("A", "a", "a", "Green-Red Opponent")
	bChannelName = channelName("B", "b", "b", "Blue-Yellow Opponent")

	rgbChannels = []model.Channel{
		extendChannel(channelName("Red", "r", "R", "Red"), numberChannel(0, 1, 6)),
		extendChannel(channelName("Green", "g", "G", "Green"), numberChannel(0, 1, 6)),
		extendChannel(channelName("Blue", "b", "B", "Blue"), numberChannel(0, 1, 6)),
	}

	xyzD50Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, data.D50Xyz[0], 8), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, data.D50Xyz[1], 8), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, data.D50Xyz[2], 8), unrestrictedChannel),
	}

	xyzD65Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, data.D65Xyz[0], 8), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, data.D65Xyz[1], 8), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, data.D65Xyz[2], 8), unrestrictedChannel),
	}

	xyzAbsD65Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, math.Round(data.D65Xyz[0]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, math.Round(data.D65Xyz[1]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, math.Round(data.D65Xyz[2]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
	}

	xyYChannels = []model.Channel{
		extendChannel(channelName("Chromaticity x", "x", "x", "x"), numberChannel(0, 1, 8)),
		extendChannel(channelName("Chromaticity y", "y", "y", "y"), numberChannel(0, 1, 8)),
		extendChannel(channelName("Luminance", "luminance", "Y", "Y"), numberChannel(0, 1, 8), unrestrictedChannel),
	}

	labChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(aChannelName, numberChannel(-125, 125, 4), unrestrictedChannel),
		extendChannel(bChannelName, numberChannel(-125, 125, 4), unrestrictedChannel),
	}

	lchChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(chromaChannelName, numberChannel(0, 150, 4), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	luvChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(uChannelName, numberChannel(-134, 220, 4), unrestrictedChannel),
		extendChannel(vChannelName, numberChannel(-140, 122, 4), unrestrictedChannel),
	}

	lchuvChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(chromaChannelName, numberChannel(0, 180, 4), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	oklabChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 1, 6)),
		extendChannel(aChannelName, numberChannel(-0.4, 0.4, 6), unrestrictedChannel),
		extendChannel(bChannelName, numberChannel(-0.4, 0.4, 6), unrestrictedChannel),
	}

	oklchChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 1, 6)),
		extendChannel(chromaChannelName, numberChannel(0, 0.4, 6), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	hueChannel        = extendChannel(channelName("Hue", "h", "H", "Hue"), degreeChannel, precisionChannel(2))
	saturationChannel = extendChannel(channelName("Saturation", "s", "S", "Saturation"), percentChannel(0, 1, 4))
	lightnessChannel  = extendChannel(channelName("Lightness", "l", "L", "Lightness"), percentChannel(0, 1, 4))
	valueChannel      = extendChannel(channelName("Value", "v", "V", "Value"), percentChannel(0, 1, 4))
	whitenessChannel  = extendChannel(channelName("Whiteness", "w", "W", "Whiteness"), percentChannel(0, 1, 4))
	blacknessChannel  = extendChannel(channelName("Blackness", "b", "B", "Blackness"), percentChannel(0, 1, 4))

	hslChannels = []model.Channel{hueChannel, saturationChannel, lightnessChannel}
	hsvChannels = []model.Channel{hueChannel, saturationChannel, valueChannel}
	hwbChannels = []model.Channel{hueChannel, whitenessChannel, blacknessChannel}

	Spaces = [...]model.Space{
		{
			Name:        "Srgb",
			Family:      "RGB",
			Base:        "LinearSrgb",
			DisplayName: "sRGB",
			CssName:     "srgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearSrgb",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear sRGB",
			CssName:     "srgb-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "DisplayP3",
			Family:      "RGB",
			Base:        "LinearDisplayP3",
			DisplayName: "Display P3",
			CssName:     "display-p3",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearDisplayP3",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Display P3",
			CssName:     "display-p3-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "A98",
			Family:      "RGB",
			Base:        "LinearA98",
			DisplayName: "Adobe RGB (1998)",
			CssName:     "a98-rgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearA98",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Adobe RGB (1998)",
			CssName:     "a98-rgb-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "ProPhoto",
			Family:      "RGB",
			Base:        "LinearProPhoto",
			DisplayName: "ProPhoto",
			CssName:     "prophoto-rgb",
			WhitePoint:  "D50",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,

			SnakeName: "prophoto",
		},
		{
			Name:        "LinearProPhoto",
			Family:      "RGB",
			Base:        "XyzD50",
			DisplayName: "Linear ProPhoto",
			CssName:     "prophoto-rgb-linear",
			WhitePoint:  "D50",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,

			SnakeName: "linear_prophoto",
		},
		{
			Name:        "Rec2020",
			Family:      "RGB",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020",
			CssName:     "rec2020",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "Rec2020OETF",
			Family:      "RGB",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020 Scene Referred",
			CssName:     "rec2020-oetf",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,

			SnakeName: "rec2020_oetf",
		},
		{
			Name:        "LinearRec2020",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Rec. 2020",
			CssName:     "rec2020-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzD50",
			Family:      "XYZ",
			Base:        "XyzD65",
			DisplayName: "CIE XYZ D50",
			CssName:     "xyz-d50",
			WhitePoint:  "D50",
			Channels:    xyzD50Channels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzD65",
			Aliases:     []string{"Xyz"},
			Family:      "XYZ",
			DisplayName: "CIE XYZ D65",
			CssName:     "xyz-d65",
			WhitePoint:  "D65",
			Channels:    xyzD65Channels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzAbsD65",
			Family:      "XYZ",
			DisplayName: "Absolute XYZ D65",
			CssName:     "xyz-abs-d65",
			WhitePoint:  "D65",
			Channels:    xyzAbsD65Channels,

			UseGenericColorFunction: true,

			SnakeName: "xyz_abs_d65",
		},
		{
			Name:        "XyYD50",
			Family:      "XYZ",
			Base:        "XyzD50",
			DisplayName: "CIE xyY",
			CssName:     "xyy-d50",
			WhitePoint:  "D50",
			Channels:    xyYChannels,

			UseGenericColorFunction: true,

			SnakeName: "xyy_d50",
		},
		{
			Name:        "XyYD65",
			Family:      "XYZ",
			Base:        "XyzD65",
			DisplayName: "CIE xyY",
			CssName:     "xyy-d65",
			WhitePoint:  "D65",
			Channels:    xyYChannels,

			UseGenericColorFunction: true,

			SnakeName: "xyy_d65",
		},
		{
			Name:        "LabD50",
			Aliases:     []string{"Lab"},
			Family:      "Lab",
			Base:        "XyzD50",
			DisplayName: "CIE Lab D50",
			CssName:     "lab",
			WhitePoint:  "D50",
			Channels:    labChannels,
		},
		{
			Name:        "LchD50",
			Aliases:     []string{"Lch"},
			Family:      "Lab",
			Base:        "LabD50",
			DisplayName: "CIE LCh D50",
			CssName:     "lch",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    lchChannels,
		},
		{
			Name:        "LabD65",
			Family:      "Lab",
			Base:        "XyzD65",
			DisplayName: "CIE Lab D65",
			CssName:     "lab-d65",
			WhitePoint:  "D65",
			Channels:    labChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LchD65",
			Family:      "Lab",
			Base:        "LabD65",
			DisplayName: "CIE LCh D65",
			CssName:     "lch-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    lchChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LuvD50",
			Aliases:     []string{"Luv"},
			Family:      "Luv",
			Base:        "XyzD50",
			DisplayName: "CIE Luv D50",
			CssName:     "luv",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    luvChannels,
		},
		{
			Name:        "LchuvD50",
			Aliases:     []string{"Lchuv"},
			Family:      "Luv",
			Base:        "LuvD50",
			DisplayName: "CIE LChuv D50",
			CssName:     "lchuv",
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    lchuvChannels,
		},
		{
			Name:        "LuvD65",
			Family:      "Luv",
			Base:        "XyzD65",
			DisplayName: "CIE Luv D65",
			CssName:     "luv-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    luvChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LchuvD65",
			Family:      "Luv",
			Base:        "LuvD65",
			DisplayName: "CIE LChuv D65",
			CssName:     "lchuv-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    lchuvChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "Oklab",
			Family:      "Oklab",
			Base:        "XyzD65",
			DisplayName: "Oklab",
			CssName:     "oklab",
			WhitePoint:  "D65",
			Channels:    oklabChannels,
		},
		{
			Name:        "Oklch",
			Family:      "Oklab",
			Base:        "Oklab",
			DisplayName: "Oklch",
			CssName:     "oklch",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    oklchChannels,
		},
		{
			Name:        "Hsl",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSL",
			CssName:     "hsl",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hslChannels,
		},
		{
			Name:        "Hsv",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSV",
			CssName:     "hsv",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hsvChannels,
		},
		{
			Name:        "Hwb",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HWB",
			CssName:     "hwb",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hwbChannels,
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
		{Pair: Pair{"Srgb", "LinearSrgb"}, Implemented: true},
		{Pair: Pair{"LinearSrgb", "Srgb"}, Implemented: true},

		{Pair: Pair{"A98", "LinearA98"}, Implemented: true},
		{Pair: Pair{"LinearA98", "A98"}, Implemented: true},

		{Pair: Pair{"DisplayP3", "LinearDisplayP3"}, Implemented: true},
		{Pair: Pair{"LinearDisplayP3", "DisplayP3"}, Implemented: true},

		{Pair: Pair{"ProPhoto", "LinearProPhoto"}, Implemented: true},
		{Pair: Pair{"LinearProPhoto", "ProPhoto"}, Implemented: true},

		{Pair: Pair{"Rec2020", "LinearRec2020"}, Implemented: true},
		{Pair: Pair{"LinearRec2020", "Rec2020"}, Implemented: true},

		{Pair: Pair{"Rec2020OETF", "LinearRec2020"}, Implemented: true},
		{Pair: Pair{"LinearRec2020", "Rec2020OETF"}, Implemented: true},

		// standard converter
		{Pair: Pair{"Srgb", "Hsl"}, Implemented: true},
		{Pair: Pair{"Hsl", "Srgb"}, Implemented: true},

		{Pair: Pair{"Srgb", "Hsv"}, Implemented: true},
		{Pair: Pair{"Hsv", "Srgb"}, Implemented: true},

		{Pair: Pair{"Srgb", "Hwb"}, Implemented: true},
		{Pair: Pair{"Hwb", "Srgb"}, Implemented: true},

		{Pair: Pair{"Hsl", "Hsv"}, Implemented: true},
		{Pair: Pair{"Hsv", "Hsl"}, Implemented: true},

		{Pair: Pair{"Hsl", "Hwb"}, Implemented: true},
		{Pair: Pair{"Hwb", "Hsl"}, Implemented: true},

		{Pair: Pair{"Hsv", "Hwb"}, Implemented: true},
		{Pair: Pair{"Hwb", "Hsv"}, Implemented: true},

		// standard converter
		{Pair: Pair{"LabD50", "XyzD50"}, Implemented: true},
		{Pair: Pair{"XyzD50", "LabD50"}, Implemented: true},

		{Pair: Pair{"LabD65", "XyzD65"}, Implemented: true},
		{Pair: Pair{"XyzD65", "LabD65"}, Implemented: true},

		{Pair: Pair{"LuvD50", "XyzD50"}, Implemented: true},
		{Pair: Pair{"XyzD50", "LuvD50"}, Implemented: true},

		{Pair: Pair{"LuvD65", "XyzD65"}, Implemented: true},
		{Pair: Pair{"XyzD65", "LuvD65"}, Implemented: true},

		{Pair: Pair{"XyzAbsD65", "XyzD65"}, Implemented: true},
		{Pair: Pair{"XyzD65", "XyzAbsD65"}, Implemented: true},

		// generate with Call Ops
		{Pair: Pair{"XyYD50", "XyzD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"XyY", "Xyz"}}}},
		{Pair: Pair{"XyzD50", "XyYD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"Xyz", "XyY"}}}},

		{Pair: Pair{"XyYD65", "XyzD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"XyY", "Xyz"}}}},
		{Pair: Pair{"XyzD65", "XyYD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"Xyz", "XyY"}}}},

		{Pair: Pair{"LabD50", "LchD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lxy", "Lch"}}}},
		{Pair: Pair{"LchD50", "LabD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lch", "Lxy"}}}},

		{Pair: Pair{"LabD65", "LchD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lxy", "Lch"}}}},
		{Pair: Pair{"LchD65", "LabD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lch", "Lxy"}}}},

		{Pair: Pair{"LuvD50", "LchuvD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lxy", "Lch"}}}},
		{Pair: Pair{"LchuvD50", "LuvD50"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lch", "Lxy"}}}},

		{Pair: Pair{"LuvD65", "LchuvD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lxy", "Lch"}}}},
		{Pair: Pair{"LchuvD65", "LuvD65"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lch", "Lxy"}}}},

		{Pair: Pair{"Oklab", "Oklch"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lxy", "Lch"}}}},
		{Pair: Pair{"Oklch", "Oklab"}, Ops: []Op{{Type: OpCall, Func: Pair{"Lch", "Lxy"}}}},

		// generate with Matrix Ops
		// Oklab
		{
			Pair: Pair{"Oklab", "XyzD65"},
			Ops: []Op{
				{Type: OpMatrix, Matrix: &data.OklabLabToLms},
				{Type: OpCube},
				{Type: OpMatrix, Matrix: &data.OklabLmsToXyzD65},
			},
		},
		{
			Pair: Pair{"XyzD65", "Oklab"},
			Ops: []Op{
				{Type: OpMatrix, Matrix: &data.OklabXyzD65ToLms},
				{Type: OpCbrt},
				{Type: OpMatrix, Matrix: &data.OklabLmsToLab},
			},
		},

		// Xyz
		{
			Pair: Pair{"XyzD65", "XyzD50"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToXyzD50}},
		},
		{
			Pair: Pair{"XyzD50", "XyzD65"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD50ToXyzD65}},
		},

		// Xyz* -> Linear*
		{
			Pair: Pair{"XyzD65", "LinearSrgb"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearSrgb}},
		},
		{
			Pair: Pair{"LinearSrgb", "XyzD65"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearSrgbToXyzD65}},
		},
		{
			Pair: Pair{"XyzD65", "LinearDisplayP3"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearDisplayP3}},
		},
		{
			Pair: Pair{"LinearDisplayP3", "XyzD65"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearDisplayP3ToXyzD65}},
		},
		{
			Pair: Pair{"XyzD65", "LinearA98"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearA98}},
		},
		{
			Pair: Pair{"LinearA98", "XyzD65"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearA98ToXyzD65}},
		},
		{
			Pair: Pair{"XyzD50", "LinearProPhoto"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD50ToLinearProPhoto}},
		},
		{
			Pair: Pair{"LinearProPhoto", "XyzD50"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearProPhotoToXyzD50}},
		},
		{
			Pair: Pair{"XyzD65", "LinearRec2020"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.XyzD65ToLinearRec2020}},
		},
		{
			Pair: Pair{"LinearRec2020", "XyzD65"},
			Ops:  []Op{{Type: OpMatrix, Matrix: &data.LinearRec2020ToXyzD65}},
		},
	}
)
