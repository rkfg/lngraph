package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

func findComponents(g graph.Graph, n graph.Node, pubkey string, min int) [][]graph.Node {
	result := [][]graph.Node{}
	compnodes := []graph.Node{}
	visited := map[graph.Node]struct{}{}
	pubkeyFound := false
	w := traverse.DepthFirst{Visit: func(n graph.Node) {
		if nn, ok := n.(Node); ok {
			if nn.Pubkey == pubkey {
				pubkeyFound = true
			}
		}
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
			pubkeyFound = false
			w.Walk(g, c, nil)
			if len(compnodes) > min && !pubkeyFound {
				result = append(result, compnodes)
			}
		}
	}
	return result
}
