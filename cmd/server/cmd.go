package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/go-jsonnet"
	"github.com/jdockerty/jsonnet-playground/internal/components"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind to")
	flag.IntVar(&port, "port", 6000, "Port binding for the server")
	flag.Parse()
}

func main() {
	bindAddress := fmt.Sprintf("%s:%d", host, port)
	vm := jsonnet.MakeVM()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Endpoints
	//
	// GET /api/share/<id>. Retrieve shared snippet hash, display in UI
	// POST /api/run <encoded-data>. Load snippet and eval with Jsonnet VM
	// POST /api/share <encoded-data>. Share code snippet, returns hash

	component := components.Page(0, 0)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.Handle("/", templ.Handler(component))

	http.HandleFunc("/api/run", func(w http.ResponseWriter, r *http.Request) {
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
		evaluated, fmtErr := vm.EvaluateAnonymousSnippet("", incomingJsonnet)
		if fmtErr != nil {
			http.Error(w, fmtErr.Error(), 400)
			return
		}

		log.Printf("Snippet:\n%s\n", evaluated)
		w.Write([]byte(evaluated))
	})

	log.Printf("Listening on %s\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}
