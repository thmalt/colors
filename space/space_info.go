package space

import (
	"strconv"
)

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
	coordinate  CoordinateSystem
	channels    []Channel

	useColorFunction bool
}

type CoordinateSystem uint8

const (
	Cartesian CoordinateSystem = iota
	Polar
)

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

func (c CoordinateSystem) String() string {
	switch c {
	case Cartesian:
		return "Cartesian"
	case Polar:
		return "Polar"
	default:
		return "CoordinateSystem(" + strconv.FormatUint(uint64(c), 10) + ")"
	}
}

func (u UnitKind) String() string {
	switch u {
	case UnitNumber:
		return "UnitNumber"
	case UnitPercent:
		return "UnitPercent"
	case UnitDegree:
		return "UnitDegree"
	default:
		return "UnitKind(" + strconv.FormatUint(uint64(u), 10) + ")"
	}
}
