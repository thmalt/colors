package codegen

import (
	"slices"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

const chIdentType = "ChIdentFlags"

func genSpacePkgChannelIdent(ctx *Context, w *writer.GoWriter) {
	identMap := make(map[string]string)
	sub := w.SubWriter()
	groups := newGroupSpace()

	alpha := "alpha"
	identMap[alpha] = chIdentName(alpha)

	for i, s := range ctx.BuiltSpaces {
		sub.Reset()
		for j, ident := range s.ChannelIdent() {
			lid := strings.ToLower(ident)
			chId := chIdentName(lid)
			identMap[lid] = chId

			if j > 0 {
				sub.Write(" | ")
			}
			sub.Write(chId)
		}
		sub.Write(" | ", identMap[alpha])
		groups.AppendExtra(string(sub.Bytes()), s.ChannelCount(), i, s, chIdentType+s.ChannelIdentFlagsName())
	}

	flagsGroups := groups.Slice()

	var identList = make([]string, 0, len(identMap))
	identList = append(identList, alpha)

	for k := range identMap {
		if k != alpha {
			identList = append(identList, k)
		}
	}

	slices.SortFunc(identList[1:], func(a, b string) int {
		if n := len(a) - len(b); n != 0 {
			return n
		}
		return strings.Compare(a, b)
	})

	w.Separate()
	w.LineWriteln("const MaxChannelCount = ", ctx.MaxChannelCount)
	w.Separate()
	w.LineWriteln("type ", chIdentType, " uint", uintBits(max(32, len(identList))))

	w.Separate()
	w.BeginGroup("const")
	for i, ident := range identList {
		w.LineWrite(identMap[ident])
		if i == 0 {
			w.Write(" ", chIdentType, " = 1 << iota")
		}
		w.InlineComment(ident)
	}
	w.End()

	w.Separate()
	w.BeginGroup("const")
	for _, group := range flagsGroups {
		sub.Reset()
		for i, ident := range group.Spaces[0].ChannelIdent() {
			if i > 0 {
				sub.Write(", ")
			}
			sub.Write(strings.ToLower(ident))
		}
		sub.Write(", ", alpha)

		w.LineWrite(group.Extra, " = ", group.Key)
		w.InlineComment(sub.Bytes())
	}
	w.End()

	var identLenMap = make(map[int][]string)
	for ident := range identMap {
		identLenMap[len(ident)] = append(identLenMap[len(ident)], ident)
	}
	var identLenList = make([]int, 0, len(identLenMap))
	for k := range identLenMap {
		identLenList = append(identLenList, k)
	}
	slices.Sort(identLenList)

	w.Separate()
	// func (f ChIdentFlags) Has(flags ChIdentFlags) bool
	w.Method("f "+chIdentType, "Has")
	w.FuncParams("flags ", chIdentType)
	w.FuncResults("bool")
	w.FuncBody()
	w.Return("f&flags == flags")
	w.End()

	w.Separate()
	// func (f ChIdentFlags) String() string
	w.Method("f "+chIdentType, "String")
	w.FuncResults("string")
	w.FuncBody()
	w.Switch("f")
	w.Case("0")
	w.Return(`"<nil>"`)
	w.Separate()
	for _, ident := range identList {
		w.Case(identMap[ident])
		w.Return('"', ident, '"')
	}
	w.Separate()
	w.Default()
	w.If("f&(f-1) == 0")
	w.Return(`"<unknown>"`)
	w.End()
	w.Return(`"<multi flags>"`)
	w.End()
	w.End()

	w.Separate()
	// func LookupChIdentFlags(ident string) ChIdentFlags
	w.Func("LookupChIdentFlags")
	w.FuncParams("ident string")
	w.FuncResults(chIdentType)
	w.FuncBody()
	w.Switch("len(ident)")
	for _, l := range identLenList {
		slices.Sort(identLenMap[l])
		if l == 1 {
			w.Case(l)
			w.Switch("ident[0] | ('a' - 'A')")
			for _, v := range identLenMap[l] {
				w.Case("'", v, "'")
				w.Return(identMap[v])
			}
			w.End()
		} else {
			w.Case(l)
			for i, v := range identLenMap[l] {
				if i == 0 {
					w.If("ascii.EqualFold(ident, ", '"', v, '"', ')')
				} else {
					w.ElseIf("ascii.EqualFold(ident, ", '"', v, '"', ')')
				}
				w.Return(identMap[v])
			}
			w.End()
		}
	}
	w.End()
	w.Return("0")
	w.End()

	w.Separate()
	// func (s Space) ChIdentFlags() ChIdentFlags
	w.Method("s Space", "ChIdentFlags")
	w.FuncResults(chIdentType)
	w.FuncBody()
	w.Switch("s")
	for _, group := range flagsGroups {
		sub.Reset()
		for i, s := range group.Spaces {
			if i > 0 {
				if i%4 == 0 {
					sub.Write(",")
					sub.Indent()
				} else {
					sub.Write(", ")
				}
			}
			sub.Write(s.Name)
		}
		w.Case(sub.Bytes())
		w.Return(group.Extra)
	}
	w.Default()
	w.Return("0")
	w.End()
	w.End()

	groups.Reset()
	for i, s := range ctx.BuiltSpaces {
		sub.Reset()
		sub.Switch("flag")
		for j, ident := range s.ChannelIdent() {
			sub.Case(identMap[strings.ToLower(ident)])
			sub.Return(j)
		}
		sub.End()
		groups.Append(string(sub.Bytes()), s.ChannelCount(), i, s)
	}

	w.Separate()
	// func (s Space) ChIndex(flag ChIdentFlags) int
	w.Method("s Space", "ChIndex")
	w.FuncParams("flag ", chIdentType)
	w.FuncResults("int")
	w.FuncBody()
	w.Switch("s")
	for _, group := range groups.Slice() {
		sub.Reset()
		for i, s := range group.Spaces {
			if i > 0 {
				if i%4 == 0 {
					sub.Write(",")
					sub.Indent()
				} else {
					sub.Write(", ")
				}
			}
			sub.Write(s.Name)
		}
		w.Case(sub.Bytes())
		w.Write(group.Key)
	}
	w.End()
	w.Separate()
	w.If("flag == ", identMap[alpha])
	w.If("n := s.ChannelCount(); n > 0")
	w.Return("n")
	w.End()
	w.End()
	w.Separate()
	w.Return("-1")
	w.End()
}

func chIdentName(ident string) string {
	return "ChIdent" + toUpperCaseFirstChar(ident)
}
