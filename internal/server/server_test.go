package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestResolveEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := resolveRequest{
		Constraint: "^1.2.0",
		Versions:   []string{"1.0.0", "1.2.3", "1.3.0", "2.0.0"},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/resolve", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp resolveResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.Found {
		t.Error("expected found=true")
	}
	if resp.Resolved != "1.3.0" {
		t.Errorf("expected 1.3.0, got %s", resp.Resolved)
	}
}

func TestResolveEndpoint_NotFound(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := resolveRequest{
		Constraint: "^3.0.0",
		Versions:   []string{"1.0.0", "2.0.0"},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/resolve", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp resolveResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Found {
		t.Error("expected found=false")
	}
}

func TestCheckEndpoint_Satisfies(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := checkRequest{Constraint: ">=1.0.0", Version: "1.5.0"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/check", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp checkResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.Satisfies {
		t.Error("expected satisfies=true")
	}
}

func TestCheckEndpoint_DoesNotSatisfy(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := checkRequest{Constraint: ">=2.0.0", Version: "1.5.0"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/check", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp checkResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Satisfies {
		t.Error("expected satisfies=false")
	}
}

func TestParseEndpoint_Valid(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := parseRequest{Version: "2.3.4-beta.1"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/parse", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp parseResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.Valid {
		t.Errorf("expected valid=true, got error: %s", resp.Error)
	}
	if resp.Major != 2 || resp.Minor != 3 || resp.Patch != 4 {
		t.Errorf("expected 2.3.4, got %d.%d.%d", resp.Major, resp.Minor, resp.Patch)
	}
}

func TestParseEndpoint_Invalid(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := parseRequest{Version: "not-a-version"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/parse", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp parseResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Valid {
		t.Error("expected valid=false")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/resolve", "/api/check", "/api/parse"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
