package topology

import (
	"bytes"
	"io"
	"main/pkg/types"

	"github.com/brynbellomy/go-utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

var LogChannel chan string

type ComputeTopologyRequest struct {
	CurrentHomeNode string   `query:"currentHomeNode"`
	IncludeNodes    []string `query:"includeNodes"`
	CrawlDistance   uint64   `query:"crawlDistance"`
	HighlightNodes  []string `query:"highlightNodes"`
}

func ComputeTopology(state *types.State, req ComputeTopologyRequest) (graph.Graph, error) {
	nodeIDs := make(map[string]int64)
	var highlightedGraphNodes []*Node

	knownRPCs := state.KnownRPCs()
	includeNodes := utils.NewSet[types.RPC]()

	// Gather all nodes to include, factoring in the crawl distance
	for _, url := range req.IncludeNodes {
		rpc, ok := knownRPCs.Get(url)
		if !ok {
			continue
		}
		includeNodes.AddSet(stackBasedCrawl(state, rpc, req.CrawlDistance))
	}

	highlightRPCs := utils.NewSet[string]()
	highlightRPCs.AddAll(req.HighlightNodes...)

	g := simple.NewUndirectedGraph()

	// Render all known RPCs
	for rpc := range includeNodes {
		var node *Node
		if highlightRPCs.Has(rpc.ID) || highlightRPCs.Has(rpc.IP) || highlightRPCs.Has(rpc.URL) || highlightRPCs.Has(rpc.Moniker) {
			node = NewNode(g.NewNode(), rpc, "crimson")
			highlightedGraphNodes = append(highlightedGraphNodes, node)
		} else {
			node = NewNode(g.NewNode(), rpc, "cadetblue")
		}

		nodeIDs[rpc.URL] = node.ID()
		g.AddNode(node)
	}

	for rpc := range includeNodes {
		nodeID1, ok := nodeIDs[rpc.URL]
		if !ok {
			continue
		}

		for _, peer := range state.RPCPeers(rpc.URL) {
			nodeID2, ok := nodeIDs[peer.URL()]
			if !ok {
				continue
			}

			if g.Edge(nodeID1, nodeID2) == nil && g.Edge(nodeID2, nodeID1) == nil {
				g.SetEdge(NewEdge(g.Node(nodeID1), g.Node(nodeID2), "azure4", "1.0"))
			}
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
				edge := g.Edge(npath[e].ID(), npath[e+1].ID()).(*Edge)
				if edge != nil {
					edge.SetColor("crimson")
					edge.SetWidth("3.0")
				}

				// hack: color the reverse path to make sure we don't color the "copied" reversed edge
				edge = g.Edge(npath[e+1].ID(), npath[e].ID()).(*Edge)
				if edge != nil {
					edge.SetColor("crimson")
					edge.SetWidth("3.0")
				}
			}
		}
	}

	return g, nil
}

func stackBasedCrawl(state *types.State, homeNode types.RPC, crawlDistance uint64) utils.Set[types.RPC] {
	visited := utils.NewSet[types.RPC]()
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
			stack = append(stack, stackItem{
				node: types.RPC{
					ID:      string(peer.NodeInfo.DefaultNodeID),
					IP:      peer.RemoteIP,
					URL:     "http://" + peer.RemoteIP + ":26657",
					Moniker: peer.NodeInfo.Moniker,
				},
				depth: item.depth + 1,
			})
		}
	}
	return visited
}

func RenderTopology(topology graph.Graph, w io.Writer) error {
	raw, err := dot.Marshal(topology, "topology", "", "")
	if err != nil {
		return err
	}

	_, err = bytes.NewReader(raw).WriteTo(w)
	return err
}
