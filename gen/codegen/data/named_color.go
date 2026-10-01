package data

import (
	"embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type NamedColor struct {
	Name  string   `json:"name"`
	RGBA  [4]uint8 `json:"rgba"`
	lower string   `json:"-"`
}

func (nc *NamedColor) Lower() string {
	if nc.lower == "" {
		nc.lower = strings.ToLower(nc.Name)
	}
	return nc.lower
}

//go:embed *.json
var namedColorFS embed.FS

const namedColorFilePrefix = "named-color"

var NamedColors = loadNamedColors()

var MaxNamedColorLength = 0

func loadNamedColors() []NamedColor {
	entries, err := namedColorFS.ReadDir(".")
	if err != nil {
		return nil
	}

	var (
		result []NamedColor
		temp   []NamedColor
		m      = map[string][4]uint8{}
		maxLen int
	)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), namedColorFilePrefix) {
			continue
		}

		r, err := namedColorFS.Open(entry.Name())
		if err != nil {
			continue
		}

		temp = temp[:0]
		if err := json.NewDecoder(r).Decode(&temp); err != nil {
			continue
		}

		for _, nc := range temp {
			lower := nc.Lower()
			if v, ok := m[lower]; ok {
				if v != nc.RGBA {
					fmt.Printf("named_color: duplicate name %q with different RGBA: existing %v, new %v\n", lower, v, nc.RGBA)
				}
				continue
			}

			// nc.Name = naming.PascalCase(nc.Name)
			m[lower] = nc.RGBA
			result = append(result, nc)

			maxLen = max(maxLen, len(lower))
		}
	}

	MaxNamedColorLength = maxLen

	slices.SortFunc(result, func(a, b NamedColor) int {
		return strings.Compare(a.lower, b.lower)
	})

	return result
}
