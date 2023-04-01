package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rajatjindal/preview-compare/internal/preview"
	"github.com/rajatjindal/preview-compare/internal/roles"
)

const uuidRegex = "[a-fA-F0-9]{8}-?[a-fA-F0-9]{4}-?4[a-fA-F0-9]{3}-?[8|9|aA|bB][a-fA-F0-9]{3}-?[a-fA-F0-9]{12}"

// Server is api server
type Server struct {
	Router *mux.Router
	store  *preview.Store
}

// New returns new server
func New() (*Server, error) {
	router := mux.NewRouter()
	store, err := preview.NewStore()
	if err != nil {
		return nil, err
	}

	server := &Server{
		Router: router,
		store:  store,
	}

	switch roles.GetRole() {
	case roles.RolePreviewMain:
		server.addRoutesForPreviewMain()
	case roles.RolePreviewFirst:
		server.addRoutesForPreviewFirst()
	case roles.RolePreviewSecond:
		server.addRoutesForPreviewSecond()
	default:
		return nil, fmt.Errorf("unexpected role: %s", roles.GetRole())
	}

	return server, nil
}

func (s *Server) addRoutesForPreviewMain() {
	s.Router.NewRoute().Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Name("options")

	s.Router.NewRoute().Methods(http.MethodPost).Path("/api/preview").HandlerFunc(s.CreateNewPreviewRequest).Name("CreateNewPreviewRequest")
	s.Router.NewRoute().Methods(http.MethodGet).Path("/api/preview/{id:preq-" + uuidRegex + "}").HandlerFunc(s.ComparePreviewWithId).Name("ComparePreviewWithId")

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"authorization", "content-type"}),
	)

	s.Router.Use(cors)
}

func (s *Server) addRoutesForPreviewFirst() {
	s.Router.NewRoute().Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Name("options")

	s.Router.NewRoute().Path("/").HandlerFunc(s.ProxyFirst).Name("ProxyFirst")
}

func (s *Server) addRoutesForPreviewSecond() {
	s.Router.NewRoute().Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Name("options")

	s.Router.NewRoute().Path("/").HandlerFunc(s.ProxySecond).Name("ProxySecond")
}
