package server

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/go-jsonnet"
)

var (
	// Do not allow import 'file:///<some_file>' expressions, as this allows
	// snooping throughout the container file system.
	disallowFileImports regexp.Regexp = *regexp.MustCompile(`file:/*`)

	// kubecfg as a library does not show the tagged build version, instead it
	// shows as "(dev build)". For now, this can be updated manually on occasional
	// bumps.
	KubecfgVersion  = "v0.34.3"
	VersionResponse = []byte(fmt.Sprintf("jsonnet: %s kubecfg: %s", jsonnet.Version(), KubecfgVersion))
)

// Health indicates whether the server is running.
func (srv *PlaygroundServer) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.State.Logger.Debug("health")
		_, _ = w.Write([]byte("OK"))
	}
}

// HandleRun receives Jsonnet input via text and evaluates it.
func (srv *PlaygroundServer) HandleRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			srv.State.Logger.Error("incorrect method to run handler", "method", r.Method)
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		srv.State.Logger.Info("run triggered", "jsonnet", incomingJsonnet)
		evaluated, err := srv.State.EvaluateSnippet(incomingJsonnet)
		if err != nil {
			srv.State.Logger.Error("invalid snippet", "jsonnet", incomingJsonnet)
			// TODO: display an error for the bad req rather than using a 200
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		srv.State.Logger.Info("evaluated", "jsonnet", evaluated)
		_, _ = w.Write([]byte(evaluated))
	}
}

// HandleCreateShare is used to create shared snippets.
// This is handled through creating a hash of the input and adding it to the state
// store - this storage mechanism is ephemeral.
//
// At a later date, this will include a persistence layer.
func (srv *PlaygroundServer) HandleCreateShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			srv.State.Logger.Error("incorrect method to create share handler", "method", r.Method)
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		_, err := srv.State.EvaluateSnippet(incomingJsonnet)
		if err != nil {
			srv.State.Logger.Error("invalid share", "jsonnet", incomingJsonnet)
			// TODO: display an error for the bad req rather than using a 200
			_, _ = w.Write([]byte("Share is not available for invalid Jsonnet. Run your snippet to see the result."))
			return
		}

		snippetHash := hex.EncodeToString(srv.State.Hasher.Sum([]byte(incomingJsonnet)))[:15]
		if _, ok := srv.State.Store[snippetHash]; !ok {
			srv.State.Logger.Info("store creation", "hash", snippetHash)
			srv.State.Store[snippetHash] = incomingJsonnet
		} else {
			srv.State.Logger.Info("store update", "hash", snippetHash)
			srv.State.Store[snippetHash] = incomingJsonnet
		}
		shareMsg := fmt.Sprintf("%s/share/%s", srv.State.Config.ShareDomain, snippetHash)
		srv.State.Logger.Debug("created share link", "link", shareMsg)
		_, _ = w.Write([]byte("Link: " + shareMsg))
	}
}

// HandleGetShare attempts to retrieve a shared snippet hash from the internal
// store. If this does not exist, an error is displayed.
func (srv *PlaygroundServer) HandleGetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			srv.State.Logger.Error("incorrect method to get share handler", "method", r.Method)
			http.Error(w, "must be GET", http.StatusBadRequest)
			return
		}
		shareHash := r.PathValue("shareHash")
		srv.State.Logger.Info("attempting to load shared snippet", "hash", shareHash)

		snippet, ok := srv.State.Store[shareHash]
		if !ok {
			srv.State.Logger.Warn("no share snippet exists", "hash", shareHash)
			errMsg := fmt.Errorf("No share snippet exists for %s, it might have expired.\n", shareHash)
			_, _ = w.Write([]byte(errMsg.Error()))
			return
		}
		srv.State.Logger.Info("loaded shared snippet", "hash", shareHash)
		_, _ = w.Write([]byte(snippet))
	}
}

// Format the input Jsonnet according to the standard jsonnetfmt rules.
func (srv *PlaygroundServer) HandleFormat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			srv.State.Logger.Error("incorrect method to format handler", "method", r.Method)
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		srv.State.Logger.Info("attempting to format", "jsonnet", incomingJsonnet)
		formattedJsonnet, err := srv.State.FormatSnippet(incomingJsonnet)
		if err != nil {
			srv.State.Logger.Warn("cannot format invalid jsonnet")
			http.Error(w, "Format is not available for invalid Jsonnet. Run your snippet to see the result.", http.StatusBadRequest)
			return
		}
		srv.State.Logger.Info("formatted", "jsonnet", formattedJsonnet)
		_, _ = w.Write([]byte(formattedJsonnet))
	}
}

// Retrieve the current version of Jsonnet/Kubecfg in use for the running application.
// This is purely informational for the frontend.
func (srv *PlaygroundServer) HandleVersions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			srv.State.Logger.Error("incorrect method to versions handler", "method", r.Method)
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(VersionResponse))
	}
}

// Middleware to stop Jsonnet snippets which contain file:///, typically paired
// with an import, being used and becoming shareable. These are rejected before
// running through the Jsonnet VM and a generic error is displayed.
func DisableFileImports(srv *PlaygroundServer, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			srv.State.Logger.Error("unable to parse form")
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}
		incomingJsonnet := r.FormValue("jsonnet-input")
		if ok := disallowFileImports.Match([]byte(incomingJsonnet)); ok {
			srv.State.Logger.Warn("attempt to import file", "jsonnet", incomingJsonnet)
			_, _ = w.Write([]byte("File imports are disabled."))
			return
		}
		next.ServeHTTP(w, r)
	}
}
