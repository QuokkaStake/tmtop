package topology

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"main/pkg/types"
)

type Node struct {
	graph.Node
	label, color string
}

func NewNode(n graph.Node, rpc types.RPC, color string) *Node {
	return &Node{
		Node:  n,
		label: rpc.Moniker + "\n(" + rpc.URL + ")",
		color: color,
	}
}

func (n *Node) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "label", Value: n.label},
		{Key: "style", Value: "filled"},
		{Key: "color", Value: n.color},
	}
}

type Edge struct {
	from, to graph.Node
	color    string
}

func NewEdge(from, to graph.Node, color string) *Edge {
	return &Edge{
		from:  from,
		to:    to,
		color: color,
	}
}

func (e *Edge) From() graph.Node {
	return e.from
}

func (e *Edge) To() graph.Node {
	return e.to
}

func (e *Edge) ReversedEdge() graph.Edge {
	return &Edge{from: e.to, to: e.from, color: e.color}
}

func (e *Edge) SetColor(color string) {
	e.color = color
}

func (e *Edge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "color", Value: e.color},
	}
}
