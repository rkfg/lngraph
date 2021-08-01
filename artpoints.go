package main

import (
	"rkfg.me/lngraph/utils"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

type ArticulationPoints struct {
	traverse.DepthFirst
	g            graph.Graph
	visited      utils.Int64s
	stack        utils.EdgeStack
	timer        uint64
	tin          map[graph.Node]uint64
	low          map[graph.Node]uint64
	result       map[graph.Node]struct{}
	rootChildren uint64
	root         graph.Node
	Visit        func(n graph.Node, parent graph.Node)
	OnSubtree    func(n graph.Edge)
}

type SubtreeTraversedEdge struct {
	graph.Edge
}

func (a *ArticulationPoints) visit(v graph.Node, parent graph.Node) {
	a.tin[v] = a.timer
	a.low[v] = a.timer
	if parent == a.root {
		a.rootChildren++
	}
	a.timer++
}

func (a *ArticulationPoints) traverse(e graph.Edge) bool {
	if a.visited.Has(e.To().ID()) {
		if a.tin[e.To()] < a.low[e.From()] {
			a.low[e.From()] = a.tin[e.To()]
		}
		return false
	}
	return true
}

func (a *ArticulationPoints) onSubtree(e graph.Edge) {
	if a.low[e.To()] < a.low[e.From()] {
		a.low[e.From()] = a.low[e.To()]
	}
	if a.low[e.To()] >= a.tin[e.From()] && e.From() != a.root {
		a.result[e.From()] = struct{}{}
	}
}

func (a *ArticulationPoints) findArticulationPointsFromRoot(n graph.Node) {
	a.rootChildren = 0
	a.root = n
	a.Walk(a.g, n)
	if a.rootChildren > 1 {
		a.result[n] = struct{}{}
	}
}

func newArticulationPoints(g graph.Graph) *ArticulationPoints {
	w := ArticulationPoints{
		g:      g,
		tin:    map[graph.Node]uint64{},
		low:    map[graph.Node]uint64{},
		result: map[graph.Node]struct{}{},
	}
	w.Visit = w.visit
	w.Traverse = w.traverse
	w.OnSubtree = w.onSubtree
	return &w
}

func (a *ArticulationPoints) resultSlice() []graph.Node {
	result := []graph.Node{}
	for k := range a.result {
		result = append(result, k)
	}
	return result
}

func findArticulationPoints(g graph.Graph) []graph.Node {
	nodes := g.Nodes()
	w := newArticulationPoints(g)
	for nodes.Next() {
		n := nodes.Node()
		if !w.visited.Has(n.ID()) {
			w.findArticulationPointsFromRoot(n)
		}
	}
	return w.resultSlice()
}

func findArticulationPointsFromRoot(g graph.Graph, n graph.Node) []graph.Node {
	w := newArticulationPoints(g)
	w.findArticulationPointsFromRoot(n)
	return w.resultSlice()
}

func (d *ArticulationPoints) Walk(g graph.Graph, from graph.Node) graph.Node {
	if d.visited == nil {
		d.visited = make(utils.Int64s)
	}
	d.stack.Push(simple.Edge{F: simple.Node(-1), T: from})
	for d.stack.Len() != 0 {
		u := d.stack.Pop()
		uid := u.To().ID()
		if st, ok := u.(SubtreeTraversedEdge); ok {
			if d.OnSubtree != nil {
				d.OnSubtree(st)
			}
			continue
		}
		if d.visited.Has(uid) {
			continue
		}
		d.visited.Add(uid)
		if d.Visit != nil {
			d.Visit(u.To(), u.From())
		}
		to := g.From(uid)
		if u.From().ID() != -1 {
			d.stack.Push(SubtreeTraversedEdge{u}) // subtree processed
		}
		for to.Next() {
			v := to.Node()
			if u.From() == v {
				// skip going back to parent
				continue
			}
			vid := v.ID()
			e := g.Edge(uid, vid)
			if d.Traverse != nil && !d.Traverse(e) {
				continue
			}
			d.stack.Push(e)
		}
	}

	return nil
}
