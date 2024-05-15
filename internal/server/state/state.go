package state

import (
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/formatter"
	"github.com/kubecfg/kubecfg/pkg/kubecfg"
)

// Providing a name for error message purposes, this has no use expect to provide
// more presentable error messages as it shows the error being in 'play.jsonnet',
// as opposed to no file name.
const PlaygroundFile = "play.jsonnet"

// New creates a new default State
func New(shareAddress string) *State {
	vm, _ := kubecfg.JsonnetVM()
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
	evaluated, fmtErr := s.Vm.EvaluateAnonymousSnippet(PlaygroundFile, snippet)
	if fmtErr != nil {
		// TODO: display an error for the bad req rather than using a 200
		return "", fmt.Errorf("Invalid Jsonnet: %w", fmtErr)
	}
	return evaluated, nil
}

func (s *State) FormatSnippet(snippet string) (string, error) {
	_, err := s.EvaluateSnippet(snippet)
	if err != nil {
		return "", err
	}

	opts := formatter.DefaultOptions()
	output, err := formatter.Format(PlaygroundFile, snippet, opts)
	if err != nil {
		return "", err
	}
	return output, nil
}

// Config contains server configuration
type Config struct {
	ShareDomain string
}
