package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/google/go-jsonnet"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jdockerty/jsonnet-playground/internal/components"
)

var (
	host         string
	port         int
	cacheSize    int
	shareAddress string

	// In-memory store for single running instance of the application.
	// TODO: multiple replicas will require a separate persistence layer.
	cache *expirable.LRU[string, string]
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind to")
	flag.StringVar(&shareAddress, "share-domain", "http://127.0.0.1", "Address prefix when sharing snippets")
	flag.IntVar(&port, "port", 8080, "Port binding for the server")
	flag.IntVar(&cacheSize, "cache-size", 1000, "Expirable LRU cache size")
	flag.Parse()
}

func main() {
	bindAddress := fmt.Sprintf("%s:%d", host, port)
	vm := jsonnet.MakeVM()

	cache = expirable.NewLRU[string, string](cacheSize, nil, time.Minute*2)
	hasher := sha512.New()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Endpoints
	//
	// GET /api/share/<id>. Retrieve shared snippet hash, display in UI
	// POST /api/run <encoded-data>. Load snippet and eval with Jsonnet VM
	// POST /api/share <encoded-data>. Share code snippet, returns hash

	rootPage := components.RootPage()
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.Handle("/", templ.Handler(rootPage))
	http.HandleFunc("/share/{shareHash}", func(w http.ResponseWriter, r *http.Request) {
		shareHash := r.PathValue("shareHash")

		if shareHash == "" {
			log.Println("Browsed to share with no hash, rendering root page")
			rootPage.Render(context.Background(), w)
			return
		}
		rootPage.Render(context.TODO(), w)
	})

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
			errMsg := fmt.Errorf("Invalid Jsonnet: %w", fmtErr)
			// TODO: display an error for the bad req rather than using a 200
			w.Write([]byte(errMsg.Error()))
			return
		}

		log.Printf("Snippet:\n%s\n", evaluated)
		w.Write([]byte(evaluated))
	})
	http.HandleFunc("/api/share/{shareHash}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "must be GET", 400)
			return
		}

		shareHash := r.PathValue("shareHash")

		_, ok := cache.Get(shareHash)
		if !ok {
			errMsg := fmt.Errorf("No share snippet exists for %s\n", shareHash)
			w.Write([]byte(errMsg.Error()))
			return
		}
		log.Printf("Loading shared snippet for %s\n", shareHash)
		rootPage.Render(context.TODO(), w)
	})

	http.HandleFunc("/api/share", func(w http.ResponseWriter, r *http.Request) {
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
			errMsg := fmt.Errorf("Invalid Jsonnet: %w", fmtErr)
			// TODO: display an error for the bad req rather than using a 200
			w.Write([]byte(errMsg.Error()))
			return
		}

		snippetHash := hex.EncodeToString(hasher.Sum([]byte(evaluated)))[:15]
		if _, ok := cache.Get(snippetHash); !ok {
			log.Printf("%s added to cache", snippetHash)
			cache.Add(snippetHash, evaluated)
		} else {
			log.Printf("cache hit for %s\n", snippetHash)
		}
		shareMsg := fmt.Sprintf("Link: %s/share/%s\n", shareAddress, snippetHash)
		w.Write([]byte(shareMsg))
	})

	log.Printf("Listening on %s\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}
