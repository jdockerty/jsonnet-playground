package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/jdockerty/jsonnet-playground/internal/components"
	"github.com/jdockerty/jsonnet-playground/internal/server/routes"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

// The playground server
type PlaygroundServer struct {
	State *state.State
}

// Load the available routes for the server
func (srv *PlaygroundServer) Routes() error {

	path, ok := os.LookupEnv("KO_DATA_PATH")
	if !ok || path == "" {
		return fmt.Errorf("KO_DATA_PATH is not set")
	}

	// Frontend routes
	rootPage := components.RootPage()
	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/", routes.HandleAssets("/assets/", fs))
	http.Handle("/", templ.Handler(rootPage))
	http.HandleFunc("/share/{shareHash}", routes.HandleShare(srv.State))

	// Backend/API routes
	http.HandleFunc("/api/health", routes.Health())
	http.HandleFunc("/api/run", routes.HandleRun(srv.State))
	http.HandleFunc("/api/share", routes.HandleCreateShare(srv.State))
	http.HandleFunc("/api/share/{shareHash}", routes.HandleGetShare(srv.State))
	return nil
}

// Serve will listen on the provided address, running the server.
func (srv *PlaygroundServer) Serve(address string) error {
	err := srv.Routes()
	if err != nil {
		return fmt.Errorf("unable to serve: %w", err)
	}
	return http.ListenAndServe(address, nil)
}
