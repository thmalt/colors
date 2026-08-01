package space

type SpaceInfo struct {
	name        string
	displayName string
	cssName     string
	whitePoint  [3]float64
	channels    []Channel

	useColorFunction bool
}

type Channel struct {
	Name        string
	Symbol      string
	DisplayName string

	Min, Max     float64
	Circular     bool
	Unrestricted bool

	Unit      UnitKind
	Precision int
}

type UnitKind uint8

const (
	UnitNumber UnitKind = iota
	UnitPercent
	UnitDegree
)

func NewSpaceInfo(name, displayName, cssName string, whitePoint [3]float64, channels []Channel) *SpaceInfo {
	return &SpaceInfo{
		name:        name,
		displayName: displayName,
		cssName:     cssName,
		whitePoint:  whitePoint,
		channels:    channels,
	}
}

func (s *SpaceInfo) Name() string           { return s.name }
func (s *SpaceInfo) DisplayName() string    { return s.displayName }
func (s *SpaceInfo) CssName() string        { return s.cssName }
func (s *SpaceInfo) WhitePoint() [3]float64 { return s.whitePoint }
func (s *SpaceInfo) ChannelCount() int      { return len(s.channels) }

func (s *SpaceInfo) Channel(index int) (Channel, bool) {
	if index < 0 || index >= s.ChannelCount() {
		return Channel{}, false
	}

	return s.channels[index], true
}

func (s *SpaceInfo) Channels() []Channel {
	return append([]Channel{}, s.channels...)
}
