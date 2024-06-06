package server_test

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/jdockerty/jsonnet-playground/internal/server"
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
		expected   string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}", expected: "../../testdata/hello-world.json", shouldFail: false},
		{name: "blank", input: "{}", expected: "../../testdata/blank.json", shouldFail: false},
		{name: "kubecfg", input: "local kubecfg = import 'internal:///kubecfg.libsonnet';\n{myVeryNestedObj:: { foo: { bar: { baz: { qux: 'some-val' }}}}, hasValue: kubecfg.objectHasPathAll($.myVeryNestedObj, 'foo.bar.baz.qux')}", expected: "../../testdata/kubecfg.json", shouldFail: false},
		{name: "invalid-jsonnet", input: "{", expected: "Invalid Jsonnet", shouldFail: true},
		{name: "invalid-jsonnet-2", input: "{hello:}", expected: "Invalid Jsonnet", shouldFail: true},
		{name: "file-import-jsonnet", input: "local f = import 'file:///proc/self/environ'; error 'test' + f", expected: "File imports are disabled", shouldFail: true},
	}

	for _, tc := range tests {
		data := url.Values{}
		data.Add("jsonnet-input", tc.input)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/run", nil)
		req.PostForm = data

		srv := server.New(state.New("127.0.0.1", "https://example.com"))
		handler := srv.HandleRun()
		handler.ServeHTTP(rec, req)

		if tc.shouldFail {
			assert.Contains(t, rec.Body.String(), tc.expected)
			return
		}

		f, err := os.Open(tc.expected)
		assert.Nil(t, err, "Unable to open %s for test", err)
		defer f.Close()
		expected, _ := io.ReadAll(f)
		assert.Equal(t, rec.Body.String(), string(expected), "[%s] expected: %s, got: %s", tc.name, expected, rec.Body.String())
	}
}

func TestHandleCreateShare(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}"},
		{name: "blank", input: "{}"},
		{name: "invalid-jsonnet", input: "{"},
		{name: "invalid-jsonnet-2", input: "{hello:}"},
		{name: "kubecfg", input: "local kubecfg = import 'internal:///kubecfg.libsonnet';\n{k8s: kubecfg.isK8sObject({apiVersion: 'v1', kind: 'Pod', spec: {}})}"},
	}

	for _, tc := range tests {
		data := url.Values{}
		data.Add("jsonnet-input", tc.input)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/share", nil)
		req.PostForm = data

		srv := server.New(state.New("127.0.0.1", "https://example.com"))
		handler := srv.HandleCreateShare()
		handler.ServeHTTP(rec, req)

		snippetHash := hex.EncodeToString(sha512.New().Sum([]byte(tc.input)))[:15]
		expected := fmt.Sprintf("Link: https://example.com/share/%s", snippetHash)
		assert.Equal(t, rec.Body.String(), expected, "expected: %s, got: %s", tc.name, expected, rec.Body.String())
	}
}

func TestHandleGetShare(t *testing.T) {
	assert := assert.New(t)
	s := state.New("127.0.0.1", "https://example.com")
	snippet := "{hello: 'world'}"
	snippetHash := hex.EncodeToString(sha512.New().Sum([]byte(snippet)))[:15]

	// Get non-existent snippet
	srv := server.New(s)
	handler := srv.HandleGetShare()
	path := fmt.Sprintf("/api/share/%s", snippetHash)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	handler.ServeHTTP(rec, req)

	assert.Contains(rec.Body.String(), "No share snippet exists")

	// Add snippet to store
	evaluated, _ := vm.EvaluateAnonymousSnippet("", snippet)
	srv.State.Store[snippetHash] = evaluated

	// Get snippet which has been added
	handler = srv.HandleGetShare()
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/share/%s", snippetHash), nil)
	req.SetPathValue("shareHash", snippetHash)
	handler.ServeHTTP(rec, req)

	assert.Equal(evaluated, rec.Body.String())

}

func TestHandleVersions(t *testing.T) {
	assert := assert.New(t)
	s := state.New("127.0.0.1", "https://example.com")

	srv := server.New(s)
	handler := srv.HandleVersions()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/versions", nil)
	handler.ServeHTTP(rec, req)

	assert.Contains(rec.Body.String(), string(server.VersionResponse))
}
