package api

import (
	"net/http"

	"github.com/rajatjindal/preview-compare/internal/proxy"
)

func (s *Server) Proxy(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r)
}
