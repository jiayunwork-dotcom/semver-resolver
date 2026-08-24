package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestResult_TextFound(t *testing.T) {
	r := &Result{Constraint: "^1.2.3", Resolved: "1.9.0", Found: true, Candidates: 4}
	var buf bytes.Buffer
	r.RenderText(&buf)
	if !strings.Contains(buf.String(), "resolved: 1.9.0") {
		t.Fatalf("unexpected: %q", buf.String())
	}
}

func TestResult_TextNone(t *testing.T) {
	r := &Result{Constraint: ">=3.0.0", Found: false, Candidates: 2}
	var buf bytes.Buffer
	r.RenderText(&buf)
	if !strings.Contains(buf.String(), "no version satisfies") {
		t.Fatalf("unexpected: %q", buf.String())
	}
}

func TestResult_JSON(t *testing.T) {
	r := &Result{Constraint: "*", Resolved: "1.0.0", Found: true, Candidates: 1}
	var buf bytes.Buffer
	if err := r.RenderJSON(&buf); err != nil {
		t.Fatalf("json error: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if out["found"] != true || out["resolved"] != "1.0.0" {
		t.Fatalf("unexpected json: %+v", out)
	}
}
