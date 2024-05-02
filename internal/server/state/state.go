package state

import (
	"crypto/sha512"
	"hash"

	"github.com/google/go-jsonnet"
)

// New creates a new default State
func New(shareAddress string) *State {
	return &State{
		Store:  make(map[string]string),
		Vm:     jsonnet.MakeVM(),
		Hasher: sha512.New(),
		Config: &Config{
			ShareDomain: shareAddress,
		},
	}
}

// State contains the shared state of the running server across all routes.
type State struct {
	Store  map[string]string
	Vm     *jsonnet.VM
	Hasher hash.Hash
	Config *Config
}

// Config contains server configuration
type Config struct {
	ShareDomain string
}
