package space

func (s Space) IsValid() bool {
	return s < SpaceCount
}

func (s Space) Info() *SpaceInfo {
	if !s.IsValid() {
		return nil
	}
	return &spaceInfos[s]
}
