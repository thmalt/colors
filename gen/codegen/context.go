package codegen

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/thmalt/colors/gen/codegen/model"
)

type Pkg struct {
	Name string
	Path string
}

func (p Pkg) Join(ident string) string {
	return p.Name + "." + ident
}

type Context struct {
	Module    string
	Directory string

	Spaces      []*model.Space
	WhitePoints []*model.WhitePoint

	BuildSpaces []*model.Space
	Funcs       []ConvertFunc
	Graph       Graph

	ConvertPkg Pkg
	RootPkg    Pkg
	SpacePkg   Pkg
	InterpPkg  Pkg

	Path string

	SplitFile bool

	FormatSource bool

	impls map[Pair]struct{}
}

func (ctx *Context) SetModuleByType(a any) {
	ctx.Module, ctx.Path = ModuleAndPathByType(a)
}

func (ctx *Context) AddSpaces(spaces []model.Space) error {
	for _, space := range spaces {
		err := ctx.AddSpace(&space)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) AddSpace(spaces ...*model.Space) error {
	if len(spaces) == 0 {
		return nil
	}

	for i, space := range spaces {
		if space == nil {
			return fmt.Errorf("space at index %d is nil", i)
		}

		if ctx.SpaceByName(space.Name) != nil {
			return fmt.Errorf("space %q at %d already exists", space.Name, i)
		}

		ctx.Spaces = append(ctx.Spaces, space)
	}

	return nil
}

func (ctx *Context) AddWhitePoint(whitePoints ...model.WhitePoint) error {
	if len(whitePoints) == 0 {
		return nil
	}

	for i, whitePoint := range whitePoints {
		if ctx.WhitePointByName(whitePoint.Name) != nil {
			return fmt.Errorf("whitePoint %q at %d already exists", whitePoint.Name, i)
		}

		ctx.WhitePoints = append(ctx.WhitePoints, &whitePoint)
	}

	return nil
}

func (ctx *Context) AddConvertFunc(funcs ...ConvertFunc) error {
	if len(funcs) == 0 {
		return nil
	}

	for _, fn := range funcs {
		if fn.Pair.From == fn.Pair.To {
			return errors.New(`conversion: "from" and "to" must be different`)
		}

		if ctx.ConvertFuncByPair(fn.Pair) != nil {
			return fmt.Errorf("conversion: Pair{%q, %q} already exists", fn.Pair.From, fn.Pair.To)
		}

		ctx.Funcs = append(ctx.Funcs, fn)
	}

	return nil
}

func (ctx *Context) Build() error {
	ctx.buildSpaces()

	ctx.impls = make(map[Pair]struct{})

	for _, fn := range ctx.Funcs {
		if fn.Implemented {
			from, to := ctx.ResolveSpacePair(fn.Pair)
			if from == nil || to == nil {
				if from == nil {
					log.Printf("space of %s not found\n", fn.Pair.From)
				}

				if to == nil {
					log.Printf("space of %s not found\n", fn.Pair.To)
				}

				continue
			}
			ctx.impls[Pair{from.Name, to.Name}] = struct{}{}
		}
	}

	return ctx.Graph.Build(ctx, ctx.Funcs)
}

func (ctx *Context) SpaceByName(name string) *model.Space {
	for _, space := range ctx.Spaces {
		if space != nil && space.Name == name {
			return space
		}
	}

	return nil
}

func (ctx *Context) WhitePointByName(name string) *model.WhitePoint {
	for _, whitePoint := range ctx.WhitePoints {
		if whitePoint != nil && whitePoint.Name == name {
			return whitePoint
		}
	}

	return nil
}

func (ctx *Context) ConvertFuncByPair(pair Pair) *ConvertFunc {
	for _, fn := range ctx.Funcs {
		if fn.Pair == pair {
			return &fn
		}
	}

	return nil
}

func (ctx *Context) ResolveSpacePairName(from, to string) (*model.Space, *model.Space) {
	return ctx.SpaceByName(from), ctx.SpaceByName(to)
}

func (ctx *Context) ResolveSpacePair(pair Pair) (from, to *model.Space) {
	return ctx.ResolveSpacePairName(pair.From, pair.To)
}

func (ctx *Context) buildSpaces() {
	families := make(map[string][]*model.Space)
	for _, s := range ctx.Spaces {
		if s.Disable {
			continue
		}

		families[s.Family] = append(families[s.Family], s)
	}

	var out []*model.Space

	for _, family := range FamilyOrder {
		group := families[family]
		slices.SortStableFunc(group, func(a, b *model.Space) int {
			oa := rgbOrder(a.Name)
			ob := rgbOrder(b.Name)
			if oa != ob {
				if oa < ob {
					return -1
				} else {
					return 1
				}
			}

			aLinear := strings.HasPrefix(a.Name, "Linear")
			bLinear := strings.HasPrefix(b.Name, "Linear")

			if aLinear != bLinear {
				return -1
			}

			return strings.Compare(a.Name, b.Name)
		})

		if len(group) == 0 {
			continue
		}

		byName := make(map[string]*model.Space, len(group))
		for _, s := range group {
			byName[s.Name] = s
		}

		visited := make(map[string]bool, len(group))

		var visit func(s *model.Space)
		visit = func(s *model.Space) {
			if visited[s.Name] {
				return
			}
			visited[s.Name] = true

			if s.Base != "" {
				if base := byName[s.Base]; base != nil {
					visit(base)
				}
			}

			out = append(out, s)
		}

		for _, s := range group {
			visit(s)
		}

		delete(families, family)
	}

	if len(families) > 0 {
		var names []string
		for name := range families {
			names = append(names, name)
		}
		slices.Sort(names)

		for _, family := range names {
			out = append(out, families[family]...)
		}
	}

	ctx.BuildSpaces = out
}

func rgbOrder(name string) int {
	name = strings.TrimPrefix(name, "Linear")

	switch name {
	case "Srgb":
		return 0
	case "DisplayP3":
		return 1
	case "A98":
		return 2
	case "ProPhoto":
		return 3
	case "Rec2020":
		return 4
	case "Hsl":
		return 100
	case "Hsv":
		return 101
	case "Hwb":
		return 102
	default:
		return 1000
	}
}
