package xrequesttrace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sp-jcberleur/xrequesttrace"
)

func TestNotExists(t *testing.T) {
	cfg := xrequesttrace.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := xrequesttrace.New(ctx, next, cfg, "xrequesttrace-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Request-ID", "")
}

func TestOtelExistsNotMatch(t *testing.T) {
	cfg := xrequesttrace.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := xrequesttrace.New(ctx, next, cfg, "xrequesttrace-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("traceparent", "not-matching-header")

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Request-ID", "")
}

func TestRequestIdAlreadyExists(t *testing.T) {
	cfg := xrequesttrace.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := xrequesttrace.New(ctx, next, cfg, "xrequesttrace-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Request-ID", "123")
	req.Header.Set("traceparent", "00-12345678901234567890123456789012-1234567890123456-01")

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Request-ID", "123")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
