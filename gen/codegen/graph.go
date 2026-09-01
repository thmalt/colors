package codegen

import (
	"math"
	"slices"

	"github.com/thmalt/colors/gen/codegen/model"
)

type Node struct {
	From *model.Space
	To   *model.Space
	Fn   *ConvertFunc

	Weight int
}

type Path = []*Node

type Graph struct {
	Nodes map[*model.Space][]Node

	cached map[pair]Path
}

type pair struct {
	From *model.Space
	To   *model.Space
}

func (g *Graph) Build(ctx *Context, funcs []ConvertFunc) error {
	g.Nodes = make(map[*model.Space][]Node)
	g.cached = make(map[pair]Path)

	for _, fn := range funcs {
		from, to := ctx.ResolveSpacePair(fn.Pair)
		g.Nodes[from] = append(g.Nodes[from], Node{
			From:   from,
			To:     to,
			Fn:     &fn,
			Weight: fn.Cost,
		})
	}

	return nil
}

func (g *Graph) FindPath(from, to *model.Space) Path {
	if path, ok := g.cached[pair{from, to}]; ok {
		return slices.Clone(path)
	}

	path := g.findPath(from, to)
	g.cached[pair{from, to}] = path

	return slices.Clone(path)
}

func (g *Graph) findPath(from, to *model.Space) Path {
	if from == to {
		return nil
	}

	dist := map[*model.Space]int{
		from: 0,
	}

	prev := map[*model.Space]*Node{}
	visited := map[*model.Space]bool{}

	for {
		var space *model.Space
		best := math.MaxInt

		for s, d := range dist {
			if visited[s] {
				continue
			}
			if d < best {
				best = d
				space = s
			}
		}

		if space == nil {
			return nil
		}

		if space == to {
			break
		}

		visited[space] = true

		for _, node := range g.Nodes[space] {
			if visited[node.To] {
				continue
			}

			weight := best + node.Weight

			if old, ok := dist[node.To]; !ok || weight < old {
				dist[node.To] = weight
				prev[node.To] = &node
			}
		}
	}

	var path Path

	for space := to; space != from; {
		node, ok := prev[space]
		if !ok {
			return nil
		}

		path = append(path, node)
		space = node.From
	}

	slices.Reverse(path)

	return path
}

func (g *Graph) FindAllPath(from, to *model.Space) []Path {
	var (
		out     []Path
		path    Path
		visited = map[*model.Space]bool{}
	)

	var dfs func(*model.Space)
	dfs = func(space *model.Space) {
		if space == to {
			p := make(Path, len(path))
			copy(p, path)
			out = append(out, p)
			return
		}

		visited[space] = true
		defer delete(visited, space)

		for i := range g.Nodes[space] {
			node := &g.Nodes[space][i]

			if visited[node.To] {
				continue
			}

			path = append(path, node)
			dfs(node.To)
			path = path[:len(path)-1]
		}
	}

	dfs(from)
	return out
}
