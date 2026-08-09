package space

func (s Space) IsValid() bool {
	return s > InvalidSpace && s < SpaceCount
}

// ChannelCount returns the number of channels in the color space.
func (s Space) ChannelCount() int {
	if uint(s) >= uint(len(spaceChannelCounts)) {
		return 0
	}
	return int(spaceChannelCounts[s])
}

func (s Space) HueIndex() int {
	if uint(s) >= uint(len(spaceHueIndexes)) {
		return -1
	}
	return int(spaceHueIndexes[s])
}

func (s Space) CoordinateSystem() CoordinateSystem {
	if uint(s) >= uint(len(spaceCoordinateSystems)) {
		return 0
	}
	return spaceCoordinateSystems[s]
}

func (s Space) Info() SpaceInfo {
	if uint(s) >= uint(len(spaceInfos)) {
		return nil
	}
	return spaceInfos[s]
}
