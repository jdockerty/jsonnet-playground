package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jdockerty/jsonnet-playground/internal/server"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

var (
	host         string
	port         int
	shareAddress string
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind to")
	flag.StringVar(&shareAddress, "share-domain", "http://127.0.0.1", "Address prefix when sharing snippets")
	flag.IntVar(&port, "port", 8080, "Port binding for the server")
	flag.Parse()
}

func main() {
	bindAddress := fmt.Sprintf("%s:%d", host, port)
	state := state.New(shareAddress)
	playground := &server.PlaygroundServer{
		State: state,
	}

	log.Printf("Listening on %s\n", bindAddress)
	log.Fatal(playground.Serve(bindAddress))

}
