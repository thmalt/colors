package codegen

import (
	"slices"
	"strconv"
)

type VariableScope struct {
	used map[string]struct{}
}

func (s *VariableScope) ReserveUnique(name string) string {
	if s.Reserve(name) {
		return name
	}

	for n := 1; ; n++ {
		out := appendInt(name, n)

		if s.Reserve(out) {
			return out
		}
	}
}

func (s *VariableScope) ReserveUniqueAll(names ...string) []string {
	out := make([]string, len(names))

	for i, name := range names {
		out[i] = s.ReserveUnique(name)
	}

	return out
}

func (s *VariableScope) ReserveUniqueN(name string, count int) []string {
	names := make([]string, count)
	n := 1

	for i := range count {
		for {
			out := appendInt(name, n)
			n++
			if s.Reserve(out) {
				names[i] = out

				break
			}
		}
	}

	return names
}

func (s *VariableScope) ReserveAll(names ...string) {
	for _, name := range names {
		s.Reserve(name)
	}
}

func (s *VariableScope) Reserve(name string) bool {
	if s.used == nil {
		s.used = make(map[string]struct{})
	}

	if s.Contains(name) {
		return false
	}

	s.used[name] = struct{}{}
	return true
}

func (s *VariableScope) CreateUnique(name string) string {
	if !s.Contains(name) {
		return name
	}

	for n := 1; ; n++ {
		out := appendInt(name, n)

		if !s.Contains(out) {
			return out
		}
	}
}

func (s *VariableScope) CreateUniqueN(name string, count int) []string {
	names := make([]string, 0, count)

	if !s.Contains(name) {
		names = append(names, name)
	}

	for n := 1; len(names) < count; n++ {
		out := appendInt(name, n)
		if !s.Contains(out) {
			names = append(names, out)
		}

	}

	return names
}

func (s *VariableScope) CreateIndexedUniqueN(name string, count int) []string {
	names := make([]string, 0, count)

	for n := 1; len(names) < count; n++ {
		out := appendInt(name, n)
		if !s.Contains(out) {
			names = append(names, out)
		}

	}

	return names
}

func (s *VariableScope) Contains(name string) bool {
	_, ok := s.used[name]
	return ok
}

func (s *VariableScope) ContainsAny(names ...string) bool {
	for _, name := range names {
		if s.Contains(name) {
			return true
		}
	}
	return false
}

func (s *VariableScope) ContainsAll(names ...string) bool {
	for _, name := range names {
		if !s.Contains(name) {
			return false
		}
	}
	return true
}

func ContainsAny(a, b []string) bool {
	for _, v := range b {
		if slices.Contains(a, v) {
			return true
		}
	}
	return false
}

func ContainsAll(a, b []string) bool {
	for _, v := range b {
		if !slices.Contains(a, v) {
			return false
		}
	}
	return true
}

func (s *VariableScope) Reset() {
	clear(s.used)
}

// //////
func (s *VariableScope) GetValidVar(prefix string, try ...[]string) []string {
	for _, vars := range try {
		if !s.HasAnyWithPrefix(prefix, vars...) {
			if prefix == "" {
				return vars
			} else {
				out := make([]string, len(vars))
				for i := range out {
					out[i] = prefix + vars[i]
				}
				return out
			}
		}
	}

	return nil
}

func (s *VariableScope) HasAnyWithPrefix(prefix string, names ...string) bool {
	for _, name := range names {
		if s.Contains(prefix + name) {
			return true
		}
	}
	return false
}

func appendInt(s string, i int) string {
	return s + strconv.Itoa(i)
}
