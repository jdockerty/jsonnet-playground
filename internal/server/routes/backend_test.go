package routes_test

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/jdockerty/jsonnet-playground/internal/server/routes"
	"github.com/jdockerty/jsonnet-playground/internal/server/state"
	"github.com/kubecfg/kubecfg/pkg/kubecfg"
	"github.com/stretchr/testify/assert"
)

var (
	vm *jsonnet.VM
)

func init() {
	vm, _ = kubecfg.JsonnetVM()
}

func TestHandleRun(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}", shouldFail: false},
		{name: "blank", input: "{}", shouldFail: false},
		{name: "kubecfg", input: "local kubecfg = import 'internal:///kubecfg.libsonnet';\n{k8s: kubecfg.isK8sObject({apiVersion: 'v1', kind: 'Pod', spec: {}})}", shouldFail: false},
		{name: "invalid-jsonnet", input: "{", shouldFail: true},
		{name: "invalid-jsonnet-2", input: "{hello:}", shouldFail: true},
	}

	for _, tc := range tests {
		data := url.Values{}
		data.Add("jsonnet-input", tc.input)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/run", nil)
		req.PostForm = data

		handler := routes.HandleRun(state.New("https://example.com"))
		handler.ServeHTTP(rec, req)

		if tc.shouldFail {
			assert.Contains(t, rec.Body.String(), "Invalid Jsonnet")
			return
		}
		expected, _ := vm.EvaluateAnonymousSnippet("", tc.input)
		assert.Equal(t, rec.Body.String(), expected, "[%s] expected: %s, got: %s", tc.name, expected, rec.Body.String())
	}
}

func TestHandleCreateShare(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}", shouldFail: false},
		{name: "blank", input: "{}", shouldFail: false},
		{name: "invalid-jsonnet", input: "{", shouldFail: true},
		{name: "invalid-jsonnet-2", input: "{hello:}", shouldFail: true},
		{name: "kubecfg", input: "local kubecfg = import 'internal:///kubecfg.libsonnet';\n{k8s: kubecfg.isK8sObject({apiVersion: 'v1', kind: 'Pod', spec: {}})}", shouldFail: false},
	}

	for _, tc := range tests {
		data := url.Values{}
		data.Add("jsonnet-input", tc.input)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/share", nil)
		req.PostForm = data

		handler := routes.HandleCreateShare(state.New("https://example.com"))
		handler.ServeHTTP(rec, req)

		if tc.shouldFail {
			assert.Contains(t, rec.Body.String(), "Share is not available for invalid Jsonnet")
			return
		}
		snippetHash := hex.EncodeToString(sha512.New().Sum([]byte(tc.input)))[:15]
		expected := fmt.Sprintf("Link: https://example.com/share/%s", snippetHash)
		assert.Equal(t, rec.Body.String(), expected, "expected: %s, got: %s", tc.name, expected, rec.Body.String())
	}
}

func TestHandleGetShare(t *testing.T) {
	assert := assert.New(t)
	s := state.New("https://example.com")
	snippet := "{hello: 'world'}"
	snippetHash := hex.EncodeToString(sha512.New().Sum([]byte(snippet)))[:15]

	// Get non-existent snippet
	handler := routes.HandleGetShare(s)
	path := fmt.Sprintf("/api/share/%s", snippetHash)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	handler.ServeHTTP(rec, req)

	assert.Contains(rec.Body.String(), "No share snippet exists")

	// Add snippet to store
	evaluated, _ := vm.EvaluateAnonymousSnippet("", snippet)
	s.Store[snippetHash] = evaluated

	// Get snippet which has been added
	handler = routes.HandleGetShare(s)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/share/%s", snippetHash), nil)
	req.SetPathValue("shareHash", snippetHash)
	handler.ServeHTTP(rec, req)

	assert.Equal(evaluated, rec.Body.String())

}
