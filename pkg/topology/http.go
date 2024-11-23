package topology

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/brynbellomy/go-utils"
	"github.com/gorilla/mux"

	tmhttp "main/pkg/http"
	"main/pkg/topology/embed"
	"main/pkg/types"
)

func WithHTTPTopologyAPI(state *types.State, highlightNodes []string) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/topology", utils.UnrestrictedCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ComputeTopologyRequest
		err := utils.UnmarshalHTTPRequest(&req, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
			return
		}

		topoGraph, err := ComputeTopology(state, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not compute topology: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/vnd.graphviz")
		if err := RenderTopology(topoGraph, w); err != nil {
			http.Error(w, fmt.Sprintf("Could not render topology: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	})))
}

func WithHTTPPeersAPI(state *types.State) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/peers", utils.UnrestrictedCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peers := state.KnownRPCs().Values()
		_ = json.NewEncoder(w).Encode(peers)
	})))
}

func WithFrontendStaticAssets() tmhttp.Option {
	return tmhttp.WithRouterOption(func(r *mux.Router) {
		assets, err := fs.Sub(embed.Frontend, "frontend/dist")
		if err != nil {
			panic(err)
		}
		r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.FS(assets)).ServeHTTP(w, r)
		})
	})
}
