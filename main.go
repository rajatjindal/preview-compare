package main

import (
	"fmt"
	"net/http"
	"os"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/preview-compare/internal/api"
	"github.com/rajatjindal/preview-compare/internal/roles"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("role=%s url=%s\n", roles.GetRole(), r.URL.String()))
		s, err := api.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR: %v\n", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		s.Router.ServeHTTP(w, r)
	})
}

func main() {}
