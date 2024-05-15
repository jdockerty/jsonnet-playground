package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/jdockerty/jsonnet-playground/internal/components"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

// HandleAssets wires up the static asset handling for the server.
func HandleAssets(pattern string, fsHandler http.Handler) http.Handler {
	return http.StripPrefix(pattern, fsHandler)
}

// HandleShare is the rendering of the shared snippet view.
func HandleShare(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shareHash := r.PathValue("shareHash")
		log.Printf("Incoming share view for %+v\n", shareHash)

		if shareHash == "" {
			log.Println("Browsed to share with no hash, rendering root page")
			components.RootPage("").Render(context.Background(), w)
			return
		}
		log.Println("Rendering share page")
		sharePage := components.RootPage(shareHash)
		sharePage.Render(context.Background(), w)
	}
}
