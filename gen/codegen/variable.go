package codegen

import (
	"fmt"
	"slices"
	"strings"
)

type VarState struct {
	used map[string]struct{}
}

func (s *VarState) ReserveName(name string, prefixes ...string) (variable string, renamed bool) {
	if name == "" {
		panic("variable name is empty")
	}

	if s.used == nil {
		s.used = make(map[string]struct{})
	}

	if !s.Contains(name) {
		s.Reserve(name)
		return name, false
	}

	base := name
	if prefix := strings.TrimSpace(strings.Join(prefixes, "")); prefix != "" {
		base = prefix + strings.ToUpper(name[:1]) + name[1:]
		if !s.Contains(base) {
			s.Reserve(base)
			return base, true
		}
	}

	for i := 2; ; i++ {
		name := fmt.Sprintf("%s%d", base, i)
		if _, ok := s.used[name]; !ok {
			s.used[name] = struct{}{}
			return name, true
		}
	}
}

func (s *VarState) ReserveNames(names []string, prefixes ...string) (newNames []string, hasNew bool) {
	newNames = make([]string, len(names))
	prefix := strings.Join(prefixes, "")

	for i, name := range names {
		newName, renamed := s.ReserveName(name, prefix)
		if renamed {
			hasNew = true
		}
		newNames[i] = newName
	}

	return
}

// x, y, z -> x*, y*, z*
func (s *VarState) ReserveNumberNames(names ...string) []string {
	var out = make([]string, len(names))

	found := false
	for i := 1; !found; i++ {
		for x, name := range names {
			newName := fmt.Sprintf("%s%d", name, i)
			out[x] = newName
			if s.Contains(newName) {
				break
			}

			found = true
		}
	}

	s.Reserve(out...)
	return out
}

// x, x, x -> x*, x*+1, x*+2, ...
func (s *VarState) ReserveNumberAddNames(names ...string) []string {
	var out = make([]string, len(names))

	num := 1
	for i, name := range names {
		for n := num; ; n++ {
			newName := fmt.Sprintf("%s%d", name, n)
			if s.Contains(newName) {
				continue
			}

			out[i] = newName
			num = n + 1
			break
		}
	}

	s.Reserve(out...)
	return out
}

func (s *VarState) Reserve(names ...string) {
	if s.used == nil {
		s.used = make(map[string]struct{})
	}

	for _, name := range names {
		s.used[name] = struct{}{}
	}
}

func (s *VarState) Contains(name string) bool {
	_, ok := s.used[name]
	return ok
}

func (s *VarState) ContainsAny(names ...string) bool {
	for _, name := range names {
		if s.Contains(name) {
			return true
		}
	}
	return false
}

func (s *VarState) ContainsAll(names ...string) bool {
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

func (s *VarState) Reset() {
	clear(s.used)
}

// //////
func (s *VarState) GetValidVar(prefix string, try ...[]string) []string {
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

func (s *VarState) HasAnyWithPrefix(prefix string, names ...string) bool {
	for _, name := range names {
		if s.Contains(prefix + name) {
			return true
		}
	}
	return false
}
