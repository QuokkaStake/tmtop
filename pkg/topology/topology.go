package topology

import (
	"bytes"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"io"
	"main/pkg/types"
)

func ComputeTopology(state *types.State, highlightNodes []string) (graph.Graph, error) {
	nodeIDs := make(map[string]int64)
	var highlightedNodes []*Node
	g := simple.NewUndirectedGraph()

	for _, rpc := range state.KnownRPCs() {
		var node *Node
		if slices.ContainsFunc(highlightNodes, func(n string) bool {
			return n == rpc.Moniker || n == rpc.IP
		}) {
			node = NewNode(g.NewNode(), rpc, "crimson")
			highlightedNodes = append(highlightedNodes, node)
		} else {
			node = NewNode(g.NewNode(), rpc, "cadetblue")
		}

		nodeIDs[rpc.URL] = node.ID()
		g.AddNode(node)
	}

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

			g.SetEdge(NewEdge(g.Node(rpcID), g.Node(peerID), "azure4"))
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
