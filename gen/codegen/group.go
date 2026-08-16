package codegen

import (
	"sort"

	"github.com/thmalt/colors/gen/codegen/model"
)

type groupSpaceValue struct {
	Key     string
	Count   int
	Indexes []int
	Spaces  []*model.Space
}

type groupSpace struct {
	M map[string]groupSpaceValue
}

func newGroupSpace() groupSpace {
	return groupSpace{
		M: make(map[string]groupSpaceValue),
	}
}

func (g groupSpace) Append(key string, count int, index int, space *model.Space) {
	v := g.M[key]

	v.Key = key
	v.Count = count

	v.Spaces = append(v.Spaces, space)
	v.Indexes = append(v.Indexes, index)

	g.M[key] = v
}

func (g groupSpace) Slice() []groupSpaceValue {
	groups := make([]groupSpaceValue, 0, len(g.M))
	for _, v := range g.M {
		groups = append(groups, v)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Indexes[0] < groups[j].Indexes[0]
	})

	return groups
}

func (g groupSpace) SortedSlice() []groupSpaceValue {
	groups := make([]groupSpaceValue, 0, len(g.M))
	for _, v := range g.M {
		groups = append(groups, v)
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Count != groups[j].Count {
			return groups[i].Count < groups[j].Count
		}
		if len(groups[i].Spaces) != len(groups[j].Spaces) {
			return len(groups[i].Spaces) < len(groups[j].Spaces)
		}
		return groups[i].Indexes[0] < groups[j].Indexes[0]
	})

	return groups
}
