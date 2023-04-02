package api

import (
	"net/http"

	"github.com/rajatjindal/preview-compare/internal/proxy"
)

func (s *Server) ProxyFirst(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r, true)
}

func (s *Server) ProxySecond(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r, false)
}
