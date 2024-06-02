package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/jdockerty/jsonnet-playground/internal/components"
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
	http.Handle("/assets/", HandleAssets("/assets/", fs))
	http.Handle("/", templ.Handler(rootPage))
	http.HandleFunc("/share/{shareHash}", srv.HandleShare())

	// Backend/API routes
	http.HandleFunc("/api/health", srv.Health())
	http.HandleFunc("/api/run", DisableFileImports(srv, srv.HandleRun()))
	http.HandleFunc("/api/format", DisableFileImports(srv, srv.HandleFormat()))
	http.HandleFunc("/api/share", DisableFileImports(srv, srv.HandleCreateShare()))
	http.HandleFunc("/api/share/{shareHash}", DisableFileImports(srv, srv.HandleGetShare()))
	http.HandleFunc("/api/versions", srv.HandleVersions())
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
