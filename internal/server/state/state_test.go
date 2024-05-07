package state_test

import (
	"testing"

	"github.com/jdockerty/jsonnet-playground/internal/server/state"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateSnippet(t *testing.T) {
	s := state.New("")

	eval, _ := s.EvaluateSnippet("{}")
	assert.Equal(t, eval, "{ }\n")

	eval, _ = s.EvaluateSnippet("{hello: 'world'}")
	expected := `{
   "hello": "world"
}
`
	assert.Equal(t, eval, expected)
}
