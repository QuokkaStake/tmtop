package topology

import (
	"fmt"
	tmhttp "main/pkg/http"
	"main/pkg/types"
	"net/http"
)

func WithHTTPTopologyAPI(state *types.State, highlightNodes []string) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/topology", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topoGraph, err := ComputeTopology(state, highlightNodes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not compute topology: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/vnd.graphviz")
		if err := RenderTopology(topoGraph, w); err != nil {
			http.Error(w, fmt.Sprintf("Could not render topology: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}))
}
