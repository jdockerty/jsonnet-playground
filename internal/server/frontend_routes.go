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
		srv.State.Logger.Info("share view loading", "shareHash", shareHash)

		if shareHash == "" {
			srv.State.Logger.Debug("browse to share with no hash, rendering root page")
			_ = components.RootPage("").Render(context.Background(), w)
			return
		}
		sharePage := components.RootPage(shareHash)
		_ = sharePage.Render(context.Background(), w)
	}
}
