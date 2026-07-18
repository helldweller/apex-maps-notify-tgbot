package apexapi

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// roundTripFunc lets a test act as the HTTP transport without a real network.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func newStubClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn, Timeout: 5 * time.Second}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

const sampleResponse = `{
  "battle_royale": {
    "current": {"start": 1000, "end": 2000, "map": "World's Edge", "code": "we", "asset": "https://example.com/we.png", "ignored": "x"},
    "next":    {"start": 2000, "end": 3000, "map": "Olympus", "code": "ol", "asset": "https://example.com/ol.png"}
  },
  "ranked": {
    "current": {"start": 1100, "end": 2100, "map": "Storm Point", "code": "sp", "asset": ""},
    "next":    {"start": 2100, "end": 3100, "map": "Broken Moon", "code": "bm", "asset": ""}
  },
  "ltm": {
    "current": {"start": 1200, "end": 2200, "map": "Kings Canyon", "code": "kc", "asset": ""},
    "next":    {"start": 2200, "end": 3200, "map": "Fragment", "code": "fr", "asset": ""}
  }
}`

func TestUpdateParsesResponse(t *testing.T) {
	var gotURL *url.URL
	SetClient(newStubClient(func(req *http.Request) (*http.Response, error) {
		gotURL = req.URL
		return jsonResponse(http.StatusOK, sampleResponse), nil
	}))
	t.Cleanup(func() { SetClient(nil) })

	var m Modes
	if err := m.Update("secret-key"); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	// Endpoint and required query params.
	if gotURL.Host != "api.mozambiquehe.re" || gotURL.Path != "/maprotation" {
		t.Errorf("unexpected endpoint: %s", gotURL)
	}
	if got := gotURL.Query().Get("auth"); got != "secret-key" {
		t.Errorf("auth query = %q, want %q", got, "secret-key")
	}
	if got := gotURL.Query().Get("version"); got != "2" {
		t.Errorf("version query = %q, want %q", got, "2")
	}

	// Battle royale (Pub) mapping.
	if m.Pub.Current.Map != "World's Edge" || m.Pub.Current.Start != 1000 || m.Pub.Current.End != 2000 {
		t.Errorf("Pub.Current mismatch: %+v", m.Pub.Current)
	}
	if m.Pub.Current.Asset != "https://example.com/we.png" {
		t.Errorf("Pub.Current.Asset = %q", m.Pub.Current.Asset)
	}
	if m.Pub.Next.Map != "Olympus" {
		t.Errorf("Pub.Next.Map = %q", m.Pub.Next.Map)
	}

	// Ranked and Ltm mapping.
	if m.Ranked.Current.Map != "Storm Point" || m.Ranked.Next.Map != "Broken Moon" {
		t.Errorf("Ranked mismatch: %+v", m.Ranked)
	}
	if m.Ltm.Current.Map != "Kings Canyon" || m.Ltm.Next.Map != "Fragment" {
		t.Errorf("Ltm mismatch: %+v", m.Ltm)
	}
}

func TestUpdateInvalidJSON(t *testing.T) {
	SetClient(newStubClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "not json"), nil
	}))
	t.Cleanup(func() { SetClient(nil) })

	var m Modes
	err := m.Update("k")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	// The error is expected to include the raw body for debuggability.
	if !strings.Contains(err.Error(), "not json") {
		t.Errorf("error should contain raw body, got: %v", err)
	}
}

func TestUpdateTransportError(t *testing.T) {
	SetClient(newStubClient(func(req *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	}))
	t.Cleanup(func() { SetClient(nil) })

	var m Modes
	if err := m.Update("k"); err == nil {
		t.Fatal("expected transport error, got nil")
	}
}

func TestSetClientNilResetsToDefault(t *testing.T) {
	SetClient(nil)
	if client == nil {
		t.Fatal("SetClient(nil) must leave a usable client")
	}
	if client.Timeout != 30*time.Second {
		t.Errorf("default client timeout = %v, want 30s", client.Timeout)
	}
}
