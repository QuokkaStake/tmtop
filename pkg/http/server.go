package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func NewServer(addr string, opts ...Option) *Server {
	router := mux.NewRouter()
	s := &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
		},
		router: router,
	}

	WithOptions(opts...)(s)

	return s
}

type Option func(*Server)

func WithOptions(opts ...Option) Option {
	return func(s *Server) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

func WithRouterOption(opt func(*mux.Router)) Option {
	return func(s *Server) {
		opt(s.router)
	}
}

func WithServerOption(opt func(*http.Server)) Option {
	return func(s *Server) {
		opt(s.server)
	}
}

func WithRoute(method, path string, handler http.Handler) Option {
	return WithRouterOption(func(r *mux.Router) {
		r.Methods(method).Path(path).Handler(handler)
	})
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}
