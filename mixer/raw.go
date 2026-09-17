package mixer

import "github.com/thmalt/colors/interp"

// RawMixer performs color interpolation directly on channel values.
type RawMixer struct {
	premultiplied bool
	hueIndex      int8
	hueInterp     interp.HueInterpolation
}

// NewRawMixer creates a [RawMixer] for color interpolation.
func NewRawMixer(premultiplied bool, hueIndex int, hueInterp interp.HueInterpolation) RawMixer {
	return RawMixer{
		premultiplied: premultiplied,
		hueIndex:      int8(hueIndex),
		hueInterp:     hueInterp,
	}
}
