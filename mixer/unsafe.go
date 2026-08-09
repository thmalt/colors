package mixer

import "github.com/thmalt/colors/interp"

type UnsafeMixer struct {
	hueIndex      int8
	premultiplied bool
	hue           interp.HueInterpolation
}

func NewUnsafeMixer(hueIndex int, premultiplied bool, hue interp.HueInterpolation) UnsafeMixer {
	return UnsafeMixer{
		hueIndex:      int8(hueIndex),
		premultiplied: premultiplied,
		hue:           hue,
	}
}
