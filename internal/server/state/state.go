package state

import (
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/google/go-jsonnet"
	"github.com/kubecfg/kubecfg/pkg/kubecfg"
)

// New creates a new default State
func New(shareAddress string) *State {
	vm, _ := kubecfg.JsonnetVM(kubecfg.WithAlpha(true))
	return &State{
		Store:  make(map[string]string),
		Vm:     vm,
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

func (s *State) EvaluateSnippet(snippet string) (string, error) {
	evaluated, fmtErr := s.Vm.EvaluateAnonymousSnippet("", snippet)
	if fmtErr != nil {
		// TODO: display an error for the bad req rather than using a 200
		return "", fmt.Errorf("Invalid Jsonnet: %w", fmtErr)
	}
	return evaluated, nil
}

// Config contains server configuration
type Config struct {
	ShareDomain string
}
