package mixer

import "github.com/thmalt/colors/interp"

// UnsafeMixer performs color interpolation without validating its inputs.
// The caller is responsible for ensuring that the inputs are valid and compatible.
type UnsafeMixer struct {
	hueIndex      int8
	premultiplied bool
	hue           interp.HueInterpolation
}

// NewUnsafeMixer creates an unsafe mixer with the specified hue channel,
// premultiplied alpha mode, and hue interpolation method.
func NewUnsafeMixer(hueIndex int, premultiplied bool, hue interp.HueInterpolation) UnsafeMixer {
	return UnsafeMixer{
		hueIndex:      int8(hueIndex),
		premultiplied: premultiplied,
		hue:           hue,
	}
}
