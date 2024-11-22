package topology

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"main/pkg/types"
)

type Node struct {
	graph.Node
	moniker, id, url, color string
}

func NewNode(n graph.Node, rpc types.RPC, color string) *Node {
	return &Node{
		Node:    n,
		moniker: rpc.Moniker,
		id:      rpc.ID,
		url:     rpc.URL,
		color:   color,
	}
}

func (n *Node) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "label", Value: n.moniker + "\n" + n.id + "\n(" + n.url + ")"},
		{Key: "style", Value: "filled"},
		{Key: "color", Value: n.color},
	}
}

type Edge struct {
	from, to     graph.Node
	color, width string
}

func NewEdge(from, to graph.Node, color string, width string) *Edge {
	return &Edge{
		from:  from,
		to:    to,
		color: color,
		width: width,
	}
}

func (e *Edge) From() graph.Node {
	return e.from
}

func (e *Edge) To() graph.Node {
	return e.to
}

func (e *Edge) ReversedEdge() graph.Edge {
	return &Edge{from: e.to, to: e.from, color: e.color, width: e.width}
}

func (e *Edge) SetColor(color string) {
	e.color = color
}

func (e *Edge) SetWidth(width string) {
	e.width = width
}

func (e *Edge) Weight() float64 {
	return 1.0
}

func (e *Edge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "color", Value: e.color},
		{Key: "penwidth", Value: e.width},
	}
}
