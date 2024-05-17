package routes

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

var (
	// Do not allow import 'file:///<some_file>' expressions, as this allows
	// snooping throughout the container file system.
	disallowFileImports regexp.Regexp = *regexp.MustCompile(`file:/*`)
)

// Health indicates whether the server is running.
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	}
}

// HandleRun receives Jsonnet input via text and evaluates it.
func HandleRun(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		evaluated, err := state.EvaluateSnippet(incomingJsonnet)
		if err != nil {
			log.Printf("Attempted to run invalid snippet: %s", incomingJsonnet)
			// TODO: display an error for the bad req rather than using a 200
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		log.Printf("Snippet:\n%s", evaluated)
		_, _ = w.Write([]byte(evaluated))
	}
}

// HandleCreateShare is used to create shared snippets.
// This is handled through creating a hash of the input and adding it to the state
// store - this storage mechanism is ephemeral.
//
// At a later date, this will include a persistence layer.
func HandleCreateShare(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		_, err := state.EvaluateSnippet(incomingJsonnet)
		if err != nil {
			log.Println("Attempted share of invalid snippet", incomingJsonnet)
			// TODO: display an error for the bad req rather than using a 200
			_, _ = w.Write([]byte("Share is not available for invalid Jsonnet. Run your snippet to see the result."))
			return
		}

		snippetHash := hex.EncodeToString(state.Hasher.Sum([]byte(incomingJsonnet)))[:15]
		if _, ok := state.Store[snippetHash]; !ok {
			log.Printf("%s added to cache", snippetHash)
			state.Store[snippetHash] = incomingJsonnet
		} else {
			log.Printf("cache hit for %s, updating snippet\n", snippetHash)
			state.Store[snippetHash] = incomingJsonnet
		}
		shareMsg := fmt.Sprintf("Link: %s/share/%s", state.Config.ShareDomain, snippetHash)
		_, _ = w.Write([]byte(shareMsg))
	}
}

// HandleGetShare attempts to retrieve a shared snippet hash from the internal
// store. If this does not exist, an error is displayed.
func HandleGetShare(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "must be GET", http.StatusBadRequest)
			return
		}
		shareHash := r.PathValue("shareHash")
		log.Printf("Call to /api/share/%s\n", shareHash)

		snippet, ok := state.Store[shareHash]
		if !ok {
			errMsg := fmt.Errorf("No share snippet exists for %s, it might have expired.\n", shareHash)
			_, _ = w.Write([]byte(errMsg.Error()))
			return
		}
		log.Printf("Loading shared snippet for %s\n", shareHash)
		_, _ = w.Write([]byte(snippet))
	}
}

// Format the input Jsonnet according to the standard jsonnetfmt rules.
func HandleFormat(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "must be POST", http.StatusBadRequest)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		log.Println("Attempting to format:", incomingJsonnet)
		formattedJsonnet, err := state.FormatSnippet(incomingJsonnet)
		if err != nil {
			log.Println("Unable to format invalid Jsonnet")
			http.Error(w, "Format is not available for invalid Jsonnet. Run your snippet to see the result.", http.StatusBadRequest)
			return
		}
		log.Println("Formatted:", formattedJsonnet)
		_, _ = w.Write([]byte(formattedJsonnet))
	}
}

func DisableFileImports(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}
		incomingJsonnet := r.FormValue("jsonnet-input")
		if ok := disallowFileImports.Match([]byte(incomingJsonnet)); ok {
			log.Println("Attempt to import file", incomingJsonnet)
			w.Write([]byte("File imports are disabled."))
			return
		}
		next.ServeHTTP(w, r)
	}
}
