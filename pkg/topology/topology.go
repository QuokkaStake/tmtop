package topology

import (
	"bytes"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"io"
	"main/pkg/types"
)

func ComputeTopology(state *types.State, highlightNodes []string) (graph.Graph, error) {
	nodeIDs := make(map[string]int64)
	var highlightedNodes []*Node
	g := simple.NewUndirectedGraph()

	// Render all known RPCs
	for _, rpc := range state.KnownRPCs() {
		var node *Node
		// Highlight nodes matched by ID, IP, or moniker
		if slices.ContainsFunc(highlightNodes, func(n string) bool {
			return n == rpc.ID || n == rpc.IP || n == rpc.Moniker
		}) {
			node = NewNode(g.NewNode(), rpc, "crimson")
			highlightedNodes = append(highlightedNodes, node)
		} else {
			node = NewNode(g.NewNode(), rpc, "cadetblue")
		}

		nodeIDs[rpc.URL] = node.ID()
		g.AddNode(node)
	}

	// Render all edges between nodes
	for _, rpc := range state.KnownRPCs() {
		rpcID, ok := nodeIDs[rpc.URL]
		if !ok {
			continue
		}

		for _, peer := range state.RPCPeers(rpc.URL) {
			peerID, ok := nodeIDs[peer.URL()]
			if !ok {
				continue
			}

			g.SetEdge(NewEdge(g.Node(rpcID), g.Node(peerID), "azure4", "1.0"))
		}
	}

	// Highlight the shortest paths between highlighted nodes
	for i, n := range highlightedNodes {
		paths := path.DijkstraFrom(n, g)
		for j := i + 1; j < len(highlightedNodes); j++ {
			npath, _ := paths.To(highlightedNodes[j].ID())
			if npath == nil {
				continue
			}

			for e := 0; e < len(npath)-1; e++ {
				edge := g.Edge(npath[e].ID(), npath[e+1].ID()).(*Edge)
				edge.SetColor("crimson")
				edge.SetWidth("3.0")
			}
		}
	}

	return g, nil
}

func RenderTopology(topology graph.Graph, w io.Writer) error {
	raw, err := dot.Marshal(topology, "topology", "", "")
	if err != nil {
		return err
	}

	_, err = bytes.NewReader(raw).WriteTo(w)
	return err
}
