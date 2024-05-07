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
	"github.com/stretchr/testify/assert"
)

func TestHandleRun(t *testing.T) {

	vm := jsonnet.MakeVM()
	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}", shouldFail: false},
		{name: "blank", input: "{}", shouldFail: false},
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

func TestHandleShare(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{name: "hello-world", input: "{hello: 'world'}", shouldFail: false},
		{name: "blank", input: "{}", shouldFail: false},
		{name: "invalid-jsonnet", input: "{", shouldFail: true},
		{name: "invalid-jsonnet-2", input: "{hello:}", shouldFail: true},
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
