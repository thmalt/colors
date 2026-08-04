package space

func (s Space) IsValid() bool {
	return s > SpaceInvalid && s < SpaceCount
}

func (s Space) ChannelCount() int {
	if uint(s) >= uint(len(spaceChannelCounts)) {
		return 0
	}
	return int(spaceChannelCounts[s])
}

func (s Space) CoordinateSystem() CoordinateSystem {
	if uint(s) >= uint(len(coordinateSystems)) {
		return 0
	}
	return coordinateSystems[s]
}

func (s Space) Info() SpaceInfo {
	if uint(s) >= uint(len(spaceInfos)) {
		return nil
	}
	return spaceInfos[s]
}
