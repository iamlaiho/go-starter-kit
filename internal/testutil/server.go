package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestServer starts an httptest.Server for the given handler and registers
// cleanup to close it when the test ends.
func NewTestServer(t *testing.T, h http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}
