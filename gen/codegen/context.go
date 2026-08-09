package codegen

import (
	"errors"
	"fmt"
	"log"

	"github.com/thmalt/colors/gen/codegen/model"
)

type Pkg struct {
	Name string
	Path string
}

func (p Pkg) Join(ident string) string {
	return p.Name + "." + ident
}

type Optimization int

const (
	OptimizeSize Optimization = iota
	OptimizeSpeed
)

type Context struct {
	Module    string
	Directory string

	Spaces      []*model.Space
	WhitePoints []*model.WhitePoint

	BuildSpaces []*model.Space
	Funcs       []ConvertFunc
	Graph       Graph

	MaxChannelCount int

	ConvertPkg Pkg
	RootPkg    Pkg
	SpacePkg   Pkg
	InterpPkg  Pkg
	MixerPkg   Pkg

	Path string

	FormatSource         bool
	SeparateAfterComment bool

	Optimization Optimization

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
	maxChannelCount := 0
	var out []*model.Space

	for _, s := range ctx.Spaces {
		if s.Disable {
			continue
		}

		out = append(out, s)

		maxChannelCount = max(maxChannelCount, s.ChannelCount())
	}

	ctx.MaxChannelCount = maxChannelCount
	ctx.BuildSpaces = out
}
