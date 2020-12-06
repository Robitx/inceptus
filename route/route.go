// Package route and its subpackages provides
// most of what you need for http server
package route

import (
	"os"

	chi "github.com/go-chi/chi"
	docgen "github.com/go-chi/docgen"
	// middleware "github.com/go-chi/chi/middleware"
)

// Router embeds chi router
type Router struct {
	chi.Router
}

// New returns new instance of router
func New() Router {
	return Router{chi.NewRouter()}
}

// GenerateDocs json docs string for the router
func GenerateDocs(r Router) string {
	// docgen need gopath for whatever reason
	if gp := os.Getenv("GOPATH"); gp == "" {
		os.Setenv("GOPATH", "/dummy")
	}
	return docgen.JSONRoutesDoc(r.Router)
}
