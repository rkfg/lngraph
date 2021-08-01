package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type Node struct {
	graph.Node `json:"node"`
	Pubkey     string `json:"pub_key"`
	Alias      string `json:"alias"`
}

type strint uint64

func (v *strint) UnmarshalJSON(data []byte) error {
	u, err := strconv.ParseUint(strings.ReplaceAll(string(data), "\"", ""), 10, 64)
	if err != nil {
		return err
	}
	*v = strint(u)
	return nil
}

type Channel struct {
	graph.Edge
	ChannelID strint `json:"channel_id"`
	Node1Pub  string `json:"node1_pub"`
	Node2Pub  string `json:"node2_pub"`
	Capacity  strint
}

type LNGraph struct {
	Nodes []Node
	Edges []Channel
}

func loadGraph(name string) graph.Graph {
	f, err := os.Open(name)
	if err != nil {
		log.Fatalf("Can't open %s: %s", name, err)
	}
	lng := LNGraph{}
	err = json.NewDecoder(f).Decode(&lng)
	if err != nil {
		log.Fatalf("Error parsing graph: %s", err)
	}
	g := simple.NewUndirectedGraph()
	pkid := map[string]Node{}
	for _, n := range lng.Nodes {
		n.Node = g.NewNode()
		g.AddNode(n)
		pkid[n.Pubkey] = n
	}
	totalcap := 0
	for _, c := range lng.Edges {
		c.Edge = g.NewEdge(pkid[c.Node1Pub], pkid[c.Node2Pub])
		g.SetEdge(c)
		totalcap += int(c.Capacity)
	}
	return g
}

type components struct {
	Point      graph.Node     `json:"point"`
	Components [][]graph.Node `json:"components"`
}

func main() {
	g := loadGraph("graph.json")
	points := findArticulationPoints(g)
	result := []components{}
	for i := range points {
		fmt.Fprintf(os.Stderr, "Processing point %d/%d\n", i, len(points))
		comps := findComponents(g, points[i], 3)
		if len(comps) > 1 {
			fmt.Fprintf(os.Stderr, "Added %d components\n", len(comps))
			result = append(result, components{Point: points[i], Components: comps})
		}
	}
	json.NewEncoder(os.Stdout).Encode(result)
	// log.Printf("Nodes: %d, Channels: %d, Capacity: %d, Connected comps: %d", g.Nodes().Len(), g.Edges().Len(), totalcap, len(comps))
}
