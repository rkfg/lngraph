package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

func findComponents(g graph.Graph, n graph.Node, min int) [][]graph.Node {
	result := [][]graph.Node{}
	compnodes := []graph.Node{}
	visited := map[graph.Node]struct{}{}
	w := traverse.DepthFirst{Visit: func(n graph.Node) {
		compnodes = append(compnodes, n)
		visited[n] = struct{}{}
	}, Traverse: func(e graph.Edge) bool {
		return e.To() != n
	},
	}
	for _, c := range graph.NodesOf(g.From(n.ID())) {
		if _, ok := visited[c]; !ok {
			w.Reset()
			compnodes = []graph.Node{}
			w.Walk(g, c, nil)
			if len(compnodes) > min {
				result = append(result, compnodes)
			}
		}
	}
	return result
}
