package named

import (
	"slices"
	"strings"

	"github.com/thmalt/colors"
)

type NamedColor struct {
	Name  string
	Color colors.Color
}

var all []NamedColor

func init() {
	all = make([]NamedColor, 0, len(lookup))

	for name, color := range lookup {
		all = append(all, NamedColor{Name: name, Color: color})
	}

	slices.SortFunc(all, func(a, b NamedColor) int {
		return strings.Compare(a.Name, b.Name)
	})
}

// All returns a sorted copy of all named colors.
func All() []NamedColor {
	return slices.Clone(all)
}
