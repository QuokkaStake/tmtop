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

func WithHTTPTopologyAPI(state *types.State) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/topology", utils.UnrestrictedCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ComputeTopologyRequest
		err := utils.UnmarshalHTTPRequest(&req, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if req.Format == "dot" {
			graph, err := ComputeTopology(state, req)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not compute topology: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			topoGraph, err := ComputeTopologyDOT(graph, req)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not convert topology to dot: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/vnd.graphviz")
			if err := RenderTopologyDOT(topoGraph, w); err != nil {
				http.Error(w, fmt.Sprintf("could not marshal dot: %s", err.Error()), http.StatusInternalServerError)
				return
			}

		} else {
			topoGraph, err := ComputeTopology(state, req)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not compute topology: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(topoGraph); err != nil {
				http.Error(w, fmt.Sprintf("could not marshal json: %s", err.Error()), http.StatusInternalServerError)
				return
			}
		}
	})))
}

func WithHTTPPeersAPI(state *types.State) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/peers", utils.UnrestrictedCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peers := state.KnownRPCs().Values()
		_ = json.NewEncoder(w).Encode(peers)
	})))
}

func WithHTTPDebugAPI(state *types.State) tmhttp.Option {
	return tmhttp.WithRoute("GET", "/debug", utils.UnrestrictedCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(state.ChainValidators)
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
