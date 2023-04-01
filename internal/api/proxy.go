package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rajatjindal/preview-compare/internal/preview"
)

func (s *Server) ProxyFirst(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	s.proxy(w, r, base)
}

func (s *Server) ProxySecond(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	s.proxy(w, r, base)
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request, base string) {
	inp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := &preview.PreviewRequest{}
	err = json.Unmarshal(inp, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	preview, err := s.store.CreatePreview(r.Context(), req)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(preview)
}
