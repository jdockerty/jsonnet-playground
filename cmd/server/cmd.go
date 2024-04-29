package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    log.Printf("Listening on %s\n", bindAddress)
    log.Fatal(http.ListenAndServe(bindAddress, nil))
}
