package model

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

type Space struct {
	// Canonical identifier
	Name    string   `json:"name"`
	Aliases []string `json:"aliases,omitempty"`

	Equivalent  string   `json:"equivalent,omitempty"`
	Equivalents []string `json:"-"`

	Family string `json:"family,omitempty"`
	Base   string `json:"base,omitempty"`

	DisplayName string `json:"displayName"`
	CssName     string `json:"cssName"`

	WhitePoint string           `json:"whitePoint"`
	Coordinate CoordinateSystem `json:"coordinate,omitempty"`
	Channels   []Channel        `json:"channels"`

	UseGenericColorFunction bool `json:"useGenericColorFunction"`
	Disable                 bool `json:"disable,omitempty"`

	// for generator
	SnakeName   string `json:"snakeName,omitempty"`
	Description string `json:"description,omitempty"`
}

type Channel struct {
	Name        string `json:"name"`
	Ident       string `json:"ident"`
	Symbol      string `json:"symbol"`
	DisplayName string `json:"displayName"`

	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Circular     bool    `json:"circular,omitempty"`
	Unrestricted bool    `json:"unrestricted,omitempty"`

	Unit      UnitKind `json:"unit,omitempty"`
	Precision int      `json:"precision,omitempty"`
}

type CoordinateSystem uint8

const (
	Cartesian CoordinateSystem = iota
	Polar
)

type UnitKind uint8

const (
	UnitNumber UnitKind = iota
	UnitPercent

	UnitDegree
	UnitRadian
	UnitGradian
	UnitTurn
)

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

func (c Channel) MinDegree() float64 {
	return AngleToDegree(c.Min, c.Unit)
}

func (c Channel) MaxDegree() float64 {
	return AngleToDegree(c.Max, c.Unit)
}

func (u UnitKind) String() string {
	switch u {
	case UnitNumber:
		return "UnitNumber"
	case UnitPercent:
		return "UnitPercent"
	case UnitDegree:
		return "UnitDegree"
	case UnitRadian:
		return "UnitRadian"
	case UnitGradian:
		return "UnitGradian"
	case UnitTurn:
		return "UnitTurn"
	default:
		return "UnitKind(" + strconv.FormatUint(uint64(u), 10) + ")"
	}
}

func (u UnitKind) MarshalJSON() ([]byte, error) {
	var s string
	switch u {
	case UnitPercent:
		s = "%"
	case UnitDegree:
		s = "deg"
	case UnitRadian:
		s = "rad"
	case UnitGradian:
		s = "grad"
	case UnitTurn:
		s = "turn"
	default:
		s = ""
	}

	return json.Marshal(s)
}

func (u *UnitKind) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}

	switch strings.ToLower(rawStr) {
	case "0":
		*u = UnitNumber
	case "1", "%":
		*u = UnitPercent
	case "2", "deg":
		*u = UnitDegree
	case "3", "rad":
		*u = UnitRadian
	case "4", "grad":
		*u = UnitGradian
	case "5", "turn":
		*u = UnitTurn
	default:
		*u = UnitNumber
	}

	return nil
}

func AngleToDegree(value float64, unit UnitKind) float64 {
	switch unit {
	case UnitTurn:
		return value * 360
	case UnitRadian:
		return value * (180 / math.Pi)
	case UnitGradian:
		return value * 0.9
	}

	return value
}

func (s Space) ChannelCount() int {
	return len(s.Channels)
}

func (s Space) ChannelSymbols() []string {
	slice := make([]string, 0, s.ChannelCount())
	for _, c := range s.Channels {
		slice = append(slice, c.Symbol)
	}

	return slice
}

func (s Space) ChannelIdent() []string {
	slice := make([]string, 0, s.ChannelCount())
	for _, c := range s.Channels {
		slice = append(slice, c.Ident)
	}

	return slice
}

func (s Space) HueIndex() int8 {
	for i, c := range s.Channels {
		if c.Circular {
			return int8(i)
		}
	}

	return -1
}
