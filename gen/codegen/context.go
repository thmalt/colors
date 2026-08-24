package codegen

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/thmalt/colors/gen/codegen/model"
)

type Pkg struct {
	Name string
	Path string
}

func (p Pkg) Join(ident string) string {
	return p.Name + "." + ident
}

type Options struct {
	// EmbedMatrix embeds conversion matrices directly into conversion functions
	// instead of calling generated conversion functions that contain the matrices.
	EmbedMatrix bool

	SeparateAfterComment bool

	// FormatSource formats generated Go source code with gofmt.
	FormatSource bool

	// BuildTags specifies the build constraints to add to generated Go files.
	BuildTags string

	// ForceWrite forces generated files to be written even when their contents have not changed.
	ForceWrite bool
}

type Context struct {
	Module    string
	Path      string
	Directory string

	Spaces      []*model.Space
	WhitePoints []*model.WhitePoint
	BuiltSpaces []*model.Space

	Funcs []ConvertFunc
	Graph Graph

	spaceMap map[string]*model.Space

	MaxChannelCount int

	ConvertPkg Pkg
	RootPkg    Pkg
	SpacePkg   Pkg
	InterpPkg  Pkg
	MixerPkg   Pkg

	// Options
	Opts Options

	impls map[Pair]struct{}

	// for logger
	TotalConversionGenerated int
}

func NewContext(opts Options) *Context {
	return &Context{
		Opts:     opts,
		spaceMap: make(map[string]*model.Space),
	}
}

func (ctx *Context) SetModuleByType(a any) {
	ctx.Module, ctx.Path = ModuleAndPathByType(a)
}

func (ctx *Context) AddSpace(spaces ...model.Space) error {
	for i := range spaces {
		if ctx.hasSpace(spaces[i].Name) {
			return fmt.Errorf("space %q at %d already exists", spaces[i].Name, i)
		}
	}

	for i := range spaces {
		ctx.addSpace(spaces[i])
	}

	return nil
}

func (ctx *Context) addSpace(space model.Space) {
	ctx.spaceMap[space.Name] = &space
	ctx.Spaces = append(ctx.Spaces, &space)
}

func (ctx *Context) hasSpace(name string) bool {
	_, ok := ctx.spaceMap[name]
	return ok
}

func (ctx *Context) AddWhitePoint(whitePoints ...model.WhitePoint) error {
	for i := range whitePoints {
		name := whitePoints[i].Name
		if ctx.WhitePointByName(name) != nil {
			return fmt.Errorf("whitePoint %q at %d already exists", name, i)
		}
	}

	for _, whitePoint := range whitePoints {
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

	for i := range ctx.Funcs {
		fn := &ctx.Funcs[i]
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
	if name == "" {
		return nil
	}

	if s, ok := ctx.spaceMap[name]; ok {
		return s
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
	for i := range ctx.Funcs {
		fn := &ctx.Funcs[i]
		if fn.Pair == pair {
			return fn
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

		if eq := ctx.SpaceByName(s.Equivalent); eq != nil {
			if slices.Contains(eq.Equivalents, s.Name) {
				log.Fatalln("duplicate space equivalent:", s.Name, " -> ", eq.Name)
			}
			eq.Equivalents = append(eq.Equivalents, s.Name)
		}
		out = append(out, s)

		maxChannelCount = max(maxChannelCount, s.ChannelCount())
	}

	ctx.MaxChannelCount = max(MinGeneratedChannelCount, maxChannelCount)
	ctx.BuiltSpaces = out
}
