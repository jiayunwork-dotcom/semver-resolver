// Package server exposes semver resolution via HTTP/JSON endpoints.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"semver-resolver/internal/resolve"
	"semver-resolver/internal/semver"
)

// Config holds server configuration.
type Config struct {
	Addr string
}

// New creates a configured http.ServeMux with all routes registered.
func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/resolve", handleResolve)
	mux.HandleFunc("/api/check", handleCheck)
	mux.HandleFunc("/api/parse", handleParse)
	return mux
}

// ListenAndServe starts the HTTP server.
func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type resolveRequest struct {
	Constraint string   `json:"constraint"`
	Versions   []string `json:"versions"`
}

type resolveResponse struct {
	Constraint string `json:"constraint"`
	Candidates int    `json:"candidates"`
	Found      bool   `json:"found"`
	Resolved   string `json:"resolved,omitempty"`
}

func handleResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req resolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Constraint == "" {
		httpError(w, http.StatusBadRequest, "constraint is empty")
		return
	}
	if len(req.Versions) == 0 {
		httpError(w, http.StatusBadRequest, "versions array is empty")
		return
	}
	c, err := semver.ParseConstraint(req.Constraint)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse constraint: "+err.Error())
		return
	}
	cands, err := resolve.ParseCandidates(strings.Join(req.Versions, ","))
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse versions: "+err.Error())
		return
	}
	best, found := resolve.Highest(c, cands)
	resp := resolveResponse{
		Constraint: req.Constraint,
		Candidates: len(cands),
		Found:      found,
	}
	if found {
		resp.Resolved = best.String()
	}
	writeJSON(w, http.StatusOK, resp)
}

type checkRequest struct {
	Constraint string `json:"constraint"`
	Version    string `json:"version"`
}

type checkResponse struct {
	Constraint string `json:"constraint"`
	Version    string `json:"version"`
	Satisfies  bool   `json:"satisfies"`
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Constraint == "" || req.Version == "" {
		httpError(w, http.StatusBadRequest, "constraint and version are required")
		return
	}
	c, err := semver.ParseConstraint(req.Constraint)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse constraint: "+err.Error())
		return
	}
	v, err := semver.Parse(req.Version)
	if err != nil {
		httpError(w, http.StatusBadRequest, "parse version: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, checkResponse{
		Constraint: req.Constraint,
		Version:    req.Version,
		Satisfies:  c.Satisfies(v),
	})
}

type parseRequest struct {
	Version string `json:"version"`
}

type parseResponse struct {
	Valid      bool   `json:"valid"`
	Major     int    `json:"major,omitempty"`
	Minor     int    `json:"minor,omitempty"`
	Patch     int    `json:"patch,omitempty"`
	Prerelease string `json:"prerelease,omitempty"`
	Error     string `json:"error,omitempty"`
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req parseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Version == "" {
		httpError(w, http.StatusBadRequest, "version is empty")
		return
	}
	v, err := semver.Parse(req.Version)
	if err != nil {
		writeJSON(w, http.StatusOK, parseResponse{Valid: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, parseResponse{
		Valid:      true,
		Major:     v.Major,
		Minor:     v.Minor,
		Patch:     v.Patch,
		Prerelease: v.Pre,
	})
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// ParsePort extracts port number from addr string.
func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

// FormatAddr produces a display-friendly address.
func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
