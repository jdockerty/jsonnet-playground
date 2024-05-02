package routes

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-jsonnet"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

var (
	vm = jsonnet.MakeVM()
)

func HandleRun(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "must be POST", 400)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", 400)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		evaluated, fmtErr := state.Vm.EvaluateAnonymousSnippet("", incomingJsonnet)
		if fmtErr != nil {
			errMsg := fmt.Errorf("Invalid Jsonnet: %w", fmtErr)
			// TODO: display an error for the bad req rather than using a 200
			w.Write([]byte(errMsg.Error()))
			return
		}

		log.Printf("Snippet:\n%s\n", evaluated)
		w.Write([]byte(evaluated))
	}
}

func HandleCreateShare(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "must be POST", 400)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", 400)
			return
		}

		incomingJsonnet := r.FormValue("jsonnet-input")
		_, fmtErr := state.Vm.EvaluateAnonymousSnippet("", incomingJsonnet)
		if fmtErr != nil {
			// TODO: display an error for the bad req rather than using a 200
			w.Write([]byte("Share is available for invalid Jsonnet\nRun your snippet to see the result."))
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
		shareMsg := fmt.Sprintf("Link: %s/share/%s\n", state.Config.ShareDomain, snippetHash)
		w.Write([]byte(shareMsg))
	}
}

func HandleGetShare(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "must be GET", 400)
			return
		}
		shareHash := r.PathValue("shareHash")
		log.Printf("Call to /api/share/%s\n", shareHash)

		snippet, ok := state.Store[shareHash]
		if !ok {
			errMsg := fmt.Errorf("No share snippet exists for %s, it might have expired.\n", shareHash)
			w.Write([]byte(errMsg.Error()))
			return
		}
		log.Printf("Loading shared snippet for %s\n", shareHash)
		w.Write([]byte(snippet))
	}
}
