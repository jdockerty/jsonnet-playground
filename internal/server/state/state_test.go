package state_test

import (
	"os"
	"testing"

	"github.com/jdockerty/jsonnet-playground/internal/server/state"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateSnippet(t *testing.T) {
	s := state.New("")

	eval, _ := s.EvaluateSnippet("{}")
	assert.Equal(t, "{ }\n", eval)

	eval, _ = s.EvaluateSnippet("{hello: 'world'}")
	expected, _ := os.ReadFile("../../../testdata/hello-world.json")
	assert.Equal(t, string(expected), eval)
}

func TestEvaluateKubecfg(t *testing.T) {
	s := state.New("")

	f, _ := os.ReadFile("../../../testdata/kubecfg.jsonnet")

	expected, _ := os.ReadFile("../../../testdata/kubecfg.json")
	eval, _ := s.EvaluateSnippet(string(f))
	assert.Equal(t, string(expected), eval)
}
