package space

type SpaceInfo interface {
	Name() string
	DisplayName() string
	CssName() string
	WhitePoint() WhitePoint
	ChannelCount() int

	Channel(index int) (Channel, bool)
	Channels() []Channel
}

type spaceInfo struct {
	name        string
	displayName string
	cssName     string
	whitePoint  WhitePoint
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

func NewSpaceInfo(name, displayName, cssName string, whitePoint WhitePoint, channels []Channel) *spaceInfo {
	return &spaceInfo{
		name:        name,
		displayName: displayName,
		cssName:     cssName,
		whitePoint:  whitePoint,
		channels:    channels,
	}
}

func (s *spaceInfo) Name() string           { return s.name }
func (s *spaceInfo) DisplayName() string    { return s.displayName }
func (s *spaceInfo) CssName() string        { return s.cssName }
func (s *spaceInfo) WhitePoint() WhitePoint { return s.whitePoint }
func (s *spaceInfo) ChannelCount() int      { return len(s.channels) }

func (s *spaceInfo) Channel(index int) (Channel, bool) {
	if index < 0 || index >= s.ChannelCount() {
		return Channel{}, false
	}

	return s.channels[index], true
}

func (s *spaceInfo) Channels() []Channel {
	return append([]Channel{}, s.channels...)
}
