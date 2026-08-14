// Copyright (c) 2011 Google Inc. All rights reserved.
//
// This file is a Go port of Skia's dithering implementation.
//
// Original source:
// src/gpu/DitherUtils.cpp
//
// See third_party/LICENSES, section "Skia - DitherUtils".

package dither

const (
	ditherSize = 8
	ditherMask = ditherSize - 1
	ditherLen  = ditherSize * ditherSize
)

var ditherLUT = makeDitherLUT()

func makeDitherLUT() [ditherLen]float64 {
	var data [ditherLen]float64
	for x := range ditherSize {
		for y := range ditherSize {
			m := (y&1)<<5 | (x&1)<<4 |
				(y&2)<<2 | (x&2)<<1 |
				(y&4)>>1 | (x&4)>>2
			data[y*ditherSize+x] = float64(m)/64.0 - 63.0/128.0
		}
	}
	return data
}

// Offset returns the normalized dithering offset at pixel position (x, y).
// The offset is intended to be applied to normalized sRGB color components.
func Offset(x, y int) float64 {
	return ditherLUT[(y&ditherMask)*ditherSize+(x&ditherMask)]
}
