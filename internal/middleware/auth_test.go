package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/iamlaiho/go-starter-kit/internal/middleware"
	"github.com/iamlaiho/go-starter-kit/internal/testutil"
)

var testSecret = []byte("test-secret-key")

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func buildToken(t *testing.T, sub string, ttl time.Duration) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub":  sub,
		"kind": "access",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(ttl).Unix(),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(testSecret)
	if err != nil {
		t.Fatalf("build token: %v", err)
	}
	return tok
}

func TestAuthenticate_MissingHeader(t *testing.T) {
	srv := testutil.NewTestServer(t, middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler)))

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	srv := testutil.NewTestServer(t, middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler)))

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	srv := testutil.NewTestServer(t, middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler)))

	tok := buildToken(t, "user-1", -1*time.Minute)
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

func TestAuthenticate_ValidToken(t *testing.T) {
	srv := testutil.NewTestServer(t, middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler)))

	tok := buildToken(t, "user-1", 15*time.Minute)
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}
