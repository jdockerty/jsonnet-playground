package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/jdockerty/jsonnet-playground/internal/server"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
)

var (
	host         string
	port         int
	shareAddress string
	logLevel     string
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host address to bind to")
	flag.StringVar(&shareAddress, "share-domain", "http://127.0.0.1", "Address prefix when sharing snippets")
	flag.StringVar(&logLevel, "log-level", "info", "Log verbosity level")
	flag.IntVar(&port, "port", 8080, "Port binding for the server")
	flag.Parse()
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug.Level()
	case "error":
		return slog.LevelError.Level()
	case "warn":
		return slog.LevelWarn.Level()
	default:
		return slog.LevelInfo.Level()
	}
}

func main() {
	bindAddress := fmt.Sprintf("%s:%d", host, port)
	logLevel := parseLogLevel(logLevel)
	log.Println("Log level set to", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	state := state.NewWithLogger(shareAddress, logger)
	playground := &server.PlaygroundServer{
		State: state,
	}

	slog.Info("Listening on", "address", bindAddress)
	log.Fatal(playground.Serve(bindAddress))
}
