package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/jdockerty/jsonnet-playground/internal/components"
	"github.com/jdockerty/jsonnet-playground/internal/server/routes"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

var (
	host           string
	port           int
	shareAddress   string
	expiryDuration time.Duration

	// In-memory store for single running instance of the application.
	cache map[string]string
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind to")
	flag.StringVar(&shareAddress, "share-domain", "http://127.0.0.1", "Address prefix when sharing snippets")
	flag.IntVar(&port, "port", 8080, "Port binding for the server")
	flag.DurationVar(&expiryDuration, "expiry", time.Minute*30, "TTL of cache entries in the LRU")
	flag.Parse()
}

func main() {
	bindAddress := fmt.Sprintf("%s:%d", host, port)
	state := state.New(shareAddress)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Endpoints
	//
	// GET /api/share/<id>. Retrieve shared snippet hash, display in UI
	// POST /api/run <encoded-data>. Load snippet and eval with Jsonnet VM
	// POST /api/share <encoded-data>. Share code snippet, returns hash

	rootPage := components.RootPage()
	fs := http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH")))
	http.Handle("/assets/", routes.HandleAssets("/assets/", fs))
	http.Handle("/", templ.Handler(rootPage))
	http.HandleFunc("/share/{shareHash}", routes.HandleShare(state))

	http.HandleFunc("/api/run", routes.HandleRun(state))
	http.HandleFunc("/api/share", routes.HandleCreateShare(state))
	http.HandleFunc("/api/share/{shareHash}", routes.HandleGetShare(state))

	log.Printf("Listening on %s\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}
