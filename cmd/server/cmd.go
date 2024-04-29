package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-jsonnet"
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

	http.HandleFunc("/api/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Must be POST"))
			return
		}

		b64EncodedSnippet, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		snippet, err := base64.StdEncoding.DecodeString(string(b64EncodedSnippet))
		if err != nil {
			panic(err)
		}

		evaluated, fmtErr := vm.EvaluateAnonymousSnippet("", string(snippet))
		if fmtErr != nil {
			fmt.Println(fmtErr)
			return
		}

		fmt.Println(evaluated)
	})

	log.Printf("Listening on %s\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}
