package topology

import (
	"bytes"
	"io"

	butils "github.com/brynbellomy/go-utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"

	"main/pkg/types"
)

var LogChannel chan string

type ComputeTopologyRequest struct {
	CurrentHomeNode string   `query:"currentHomeNode"`
	IncludeNodes    []string `query:"includeNodes"`
	CrawlDistance   uint64   `query:"crawlDistance"`
	MinBytesSec     uint64   `query:"minBytesSec"`
	HighlightNodes  []string `query:"highlightNodes"`
	Format          string   `query:"format"`
}

func ComputeTopology(state *types.State, req ComputeTopologyRequest) (Graph, error) {
	knownRPCs := state.KnownRPCs()
	includeNodes := butils.NewSet[types.RPC]()
	includeIDs := butils.NewSet[string]()

	// Gather all nodes to include, factoring in the crawl distance
	for _, url := range req.IncludeNodes {
		rpc, ok := knownRPCs.Get(url)
		if !ok {
			continue
		}
		includeNodes.AddSet(stackBasedCrawl(state, rpc, req.CrawlDistance, req.MinBytesSec))
	}

	var g Graph

	// Add nodes
	for rpc := range includeNodes {
		g.Nodes = append(g.Nodes, rpc)
		includeIDs.Add(rpc.ID)
	}

	// Add edges
	for rpc := range includeNodes {
		for _, peer := range state.RPCPeers(rpc.URL) {
			if !includeIDs.Has(string(peer.NodeInfo.DefaultNodeID)) {
				continue
			}

			if peer.IsOutbound {
				g.AddConn(rpc.ID, string(peer.NodeInfo.DefaultNodeID), peer.ConnectionStatus)
			} else {
				g.AddConn(string(peer.NodeInfo.DefaultNodeID), rpc.ID, peer.ConnectionStatus)
			}
		}
	}

	return g, nil
}

func stackBasedCrawl(state *types.State, homeNode types.RPC, crawlDistance uint64, minBytesSec uint64) butils.Set[types.RPC] {
	visited := butils.NewSet[types.RPC]()
	visited.Add(homeNode)

	type stackItem struct {
		node  types.RPC
		depth uint64
	}
	stack := []stackItem{{node: homeNode, depth: 0}}

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visited.Add(item.node)

		if item.depth >= crawlDistance {
			continue
		}

		for _, peer := range state.RPCPeers(item.node.URL) {
			if uint64(peer.ConnectionStatus.SendMonitor.AvgRate)+uint64(peer.ConnectionStatus.RecvMonitor.AvgRate) < minBytesSec {
				continue
			}
			stack = append(stack, stackItem{
				node:  types.NewRPCFromPeer(peer),
				depth: item.depth + 1,
			})
		}
	}
	return visited
}

func ComputeTopologyDOT(topology Graph, req ComputeTopologyRequest) (graph.Graph, error) {
	nodeIDs := make(map[string]int64)
	var highlightedGraphNodes []*DOTPeerNode

	highlightRPCs := butils.NewSet[string]()
	highlightRPCs.AddAll(req.HighlightNodes...)

	g := simple.NewUndirectedGraph()

	// Render all known RPCs
	for _, rpc := range topology.Nodes {
		var node *DOTPeerNode
		if highlightRPCs.Has(rpc.ID) || highlightRPCs.Has(rpc.IP) || highlightRPCs.Has(rpc.URL) || highlightRPCs.Has(rpc.Moniker) {
			node = NewDOTPeerNode(g.NewNode(), rpc, "crimson")
			highlightedGraphNodes = append(highlightedGraphNodes, node)
		} else {
			node = NewDOTPeerNode(g.NewNode(), rpc, "cadetblue")
		}

		nodeIDs[rpc.URL] = node.Node.ID()
		g.AddNode(node)
	}

	// Add edges
	for _, conn := range topology.Conns {
		nodeID1, ok := nodeIDs[conn.From]
		if !ok {
			continue
		}

		nodeID2, ok := nodeIDs[conn.To]
		if !ok {
			continue
		}

		if g.Edge(nodeID1, nodeID2) == nil && g.Edge(nodeID2, nodeID1) == nil {
			g.SetEdge(NewDOTEdge(g.Node(nodeID1), g.Node(nodeID2), "azure4", "1.0"))
		}
	}

	for i, n := range highlightedGraphNodes {
		paths := path.DijkstraFrom(n, g)
		for j := i + 1; j < len(highlightedGraphNodes); j++ {
			npath, _ := paths.To(highlightedGraphNodes[j].ID())
			if npath == nil {
				continue
			}

			for e := 0; e < len(npath)-1; e++ {
				edge := g.Edge(npath[e].ID(), npath[e+1].ID()).(*DOTEdge)
				if edge != nil {
					edge.SetColor("crimson")
					edge.SetWidth("3.0")
				}

				// hack: color the reverse path to make sure we don't color the "copied" reversed edge
				edge = g.Edge(npath[e+1].ID(), npath[e].ID()).(*DOTEdge)
				if edge != nil {
					edge.SetColor("crimson")
					edge.SetWidth("3.0")
				}
			}
		}
	}

	return g, nil
}

func RenderTopologyDOT(topology graph.Graph, w io.Writer) error {
	raw, err := dot.Marshal(topology, "topology", "", "")
	if err != nil {
		return err
	}

	_, err = bytes.NewReader(raw).WriteTo(w)
	return err
}
