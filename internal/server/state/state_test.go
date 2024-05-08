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
	expected := `{
   "hello": "world"
}
`
	assert.Equal(t, expected, eval)
}

func TestEvaluateKubecfg(t *testing.T) {
	s := state.New("")

	f, _ := os.ReadFile("../../../testdata/kubecfg.jsonnet")

	expected := `{
   "hasValue": true
}
`
	eval, _ := s.EvaluateSnippet(string(f))
	assert.Equal(t, expected, eval)
}
