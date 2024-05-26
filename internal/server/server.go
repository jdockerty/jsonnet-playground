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
	Server http.Server
	State  *state.State
}

func New(state *state.State) *PlaygroundServer {
	return &PlaygroundServer{
		State: state,
		Server: http.Server{
			Addr: state.Config.Address,
		},
	}
}

// Load the available routes for the server
func (srv *PlaygroundServer) Routes() error {

	path, ok := os.LookupEnv("KO_DATA_PATH")
	if !ok || path == "" {
		return fmt.Errorf("KO_DATA_PATH is not set")
	}

	// Frontend routes
	rootPage := components.RootPage("")
	fs := http.FileServer(http.Dir(path))
	http.Handle("/assets/", routes.HandleAssets("/assets/", fs))
	http.Handle("/", templ.Handler(rootPage))
	http.HandleFunc("/share/{shareHash}", routes.HandleShare(srv.State))

	// Backend/API routes
	http.HandleFunc("/api/health", routes.Health(srv.State))
	http.HandleFunc("/api/run", routes.DisableFileImports(srv.State, routes.HandleRun(srv.State)))
	http.HandleFunc("/api/format", routes.DisableFileImports(srv.State, routes.HandleFormat(srv.State)))
	http.HandleFunc("/api/share", routes.DisableFileImports(srv.State, routes.HandleCreateShare(srv.State)))
	http.HandleFunc("/api/share/{shareHash}", routes.DisableFileImports(srv.State, routes.HandleGetShare(srv.State)))
	http.HandleFunc("/api/versions", routes.HandleVersions(srv.State))
	return nil
}

// Serve will listen on the provided address, running the server.
func (srv *PlaygroundServer) Serve() error {
	err := srv.Routes()
	if err != nil {
		return fmt.Errorf("unable to serve: %w", err)
	}
	return srv.Server.ListenAndServe()
}
