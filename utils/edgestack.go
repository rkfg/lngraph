package utils

import "gonum.org/v1/gonum/graph"

// EdgeStack implements a LIFO stack of graph.Node.
type EdgeStack []graph.Edge

// Len returns the number of graph.Nodes on the stack.
func (s *EdgeStack) Len() int { return len(*s) }

// Pop returns the last graph.Node on the stack and removes it
// from the stack.
func (s *EdgeStack) Pop() graph.Edge {
	v := *s
	v, n := v[:len(v)-1], v[len(v)-1]
	*s = v
	return n
}

// Push adds the node n to the stack at the last position.
func (s *EdgeStack) Push(n graph.Edge) { *s = append(*s, n) }
