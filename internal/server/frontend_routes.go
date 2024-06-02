package server

import (
	"context"
	"log"
	"net/http"

	"github.com/jdockerty/jsonnet-playground/internal/components"
)

// HandleAssets wires up the static asset handling for the server.
func HandleAssets(pattern string, fsHandler http.Handler) http.Handler {
	return http.StripPrefix(pattern, fsHandler)
}

// HandleShare is the rendering of the shared snippet view.
func (srv *PlaygroundServer) HandleShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shareHash := r.PathValue("shareHash")
		log.Printf("Incoming share view for %+v\n", shareHash)

		if shareHash == "" {
			log.Println("Browsed to share with no hash, rendering root page")
			_ = components.RootPage("").Render(context.Background(), w)
			return
		}
		log.Println("Rendering share page")
		sharePage := components.RootPage(shareHash)
		_ = sharePage.Render(context.Background(), w)
	}
}