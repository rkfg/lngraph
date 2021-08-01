package main

import (
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type graphCase struct {
	num      int
	edges    [][]int64
	expected map[int64]struct{}
}

func TestPoints(t *testing.T) {
	cases := []graphCase{
		{
			num: 5,
			edges: [][]int64{
				{0, 1},
				{1, 2},
				{0, 2},
				{0, 3},
				{3, 4},
			},
			expected: map[int64]struct{}{
				0: {},
				3: {},
			},
		},
		{
			num: 4,
			edges: [][]int64{
				{0, 1},
				{1, 2},
				{2, 3},
			},
			expected: map[int64]struct{}{
				1: {},
				2: {},
			},
		},
		{
			num: 7,
			edges: [][]int64{
				{0, 1},
				{0, 2},
				{2, 1},
				{1, 3},
				{1, 4},
				{1, 6},
				{3, 5},
				{4, 5},
			},
			expected: map[int64]struct{}{
				1: {},
			},
		},
	}
	for _, c := range cases {
		g := simple.NewUndirectedGraph()
		for i := 0; i < c.num; i++ {
			n := g.NewNode()
			g.AddNode(n)
		}
		for _, e := range c.edges {
			g.SetEdge(g.NewEdge(g.Node(e[0]), g.Node(e[1])))
		}
		t.Logf("Nodes: %+v Edges: %+v", g.Nodes(), g.Edges())
		pts := findArticulationPoints(g)
		checkPts(pts, c.expected, t)
		for i := 0; i < c.num; i++ {
			n := g.Node(int64(i))
			pts := findArticulationPointsFromRoot(g, n)
			checkPts(pts, c.expected, t)
		}
	}
}

func checkPts(pts []graph.Node, expected map[int64]struct{}, t *testing.T) {
	if len(pts) != len(expected) {
		t.Fatalf("Found %d points %+v, expected %d", len(pts), pts, len(expected))
	}
	for i := range pts {
		if _, ok := expected[pts[i].ID()]; !ok {
			t.Fatalf("Point with id %d isn't expected", pts[i].ID())
		}
	}
}
