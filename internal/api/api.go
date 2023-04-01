package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rajatjindal/preview-compare/internal/preview"
)

const uuidRegex = "[a-fA-F0-9]{8}-?[a-fA-F0-9]{4}-?4[a-fA-F0-9]{3}-?[8|9|aA|bB][a-fA-F0-9]{3}-?[a-fA-F0-9]{12}"

// Server is api server
type Server struct {
	Router *mux.Router
	store  preview.Store
}

// New returns new server
func New() (*Server, error) {
	router := mux.NewRouter()
	server := &Server{
		Router: router,
	}

	server.addRoutes()
	return server, nil
}

func (s *Server) addRoutes() {
	s.Router.NewRoute().Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Name("options")

	s.Router.NewRoute().Methods(http.MethodPost).Path("/api/preview").HandlerFunc(s.CreateNewPreviewRequest).Name("CreateNewPreviewRequest")

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"authorization", "content-type"}),
	)

	s.Router.Use(cors)
}
