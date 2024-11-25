package topology

import (
	"main/pkg/types"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

type Graph struct {
	Nodes []types.RPC `json:"nodes"`
	Conns []Conn      `json:"conns"`

	connsMap map[string]bool
}

func (g *Graph) AddConn(from, to string, connectionStatus types.ConnectionStatus) {
	if g.connsMap == nil {
		g.connsMap = make(map[string]bool)
	}

	_, ok := g.connsMap[from+to]
	if !ok {
		g.Conns = append(g.Conns, NewConn(from, to, connectionStatus))
		g.connsMap[from+to] = true
	}
}

type Conn struct {
	From             string                 `json:"from"`
	To               string                 `json:"to"`
	ConnectionStatus types.ConnectionStatus `json:"connectionStatus"`
}

func NewConn(from, to string, connectionStatus types.ConnectionStatus) Conn {
	return Conn{
		From:             from,
		To:               to,
		ConnectionStatus: connectionStatus,
	}
}

type DOTPeerNode struct {
	graph.Node
	types.RPC
	Color string
}

func NewDOTPeerNode(n graph.Node, rpc types.RPC, color string) *DOTPeerNode {
	return &DOTPeerNode{
		Node:  n,
		RPC:   rpc,
		Color: color,
	}
}

func (n *DOTPeerNode) ID() int64 {
	return n.Node.ID()
}

func (n *DOTPeerNode) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "label", Value: n.Moniker + "\n" + n.RPC.ID + "\n(" + n.URL + ")"},
		{Key: "style", Value: "filled"},
		{Key: "color", Value: n.Color},
	}
}

type DOTEdge struct {
	from, to     graph.Node
	color, width string
}

func NewDOTEdge(from, to graph.Node, color string, width string) *DOTEdge {
	return &DOTEdge{
		from:  from,
		to:    to,
		color: color,
		width: width,
	}
}

func (e *DOTEdge) From() graph.Node {
	return e.from
}

func (e *DOTEdge) To() graph.Node {
	return e.to
}

func (e *DOTEdge) ReversedEdge() graph.Edge {
	return &DOTEdge{from: e.to, to: e.from, color: e.color, width: e.width}
}

func (e *DOTEdge) SetColor(color string) {
	e.color = color
}

func (e *DOTEdge) SetWidth(width string) {
	e.width = width
}

func (e *DOTEdge) Weight() float64 {
	return 1.0
}

func (e *DOTEdge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "color", Value: e.color},
		{Key: "penwidth", Value: e.width},
	}
}
