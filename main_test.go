package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoCommand(t *testing.T) {
	if err := run([]string{}, nil, &bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Fatalf("expected error for missing command")
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	if err := run([]string{"bogus"}, nil, &bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Fatalf("expected error for unknown command")
	}
}

func TestRun_Resolve_PicksHighest(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{"resolve", "--constraint", "^1.2.3", "--versions", "1.2.0,1.2.3,1.9.0,2.0.0"}, nil, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v (stderr=%s)", err, errOut.String())
	}
	if !strings.Contains(out.String(), "resolved: 1.9.0") {
		t.Fatalf("expected 1.9.0 resolved, got %q", out.String())
	}
}

func TestRun_Resolve_InvalidConstraint(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{"resolve", "--constraint", "!!!", "--versions", "1.0.0"}, nil, &out, &errOut)
	if err == nil {
		t.Fatalf("expected error for invalid constraint")
	}
}

func TestRun_Check_Satisfies(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{"check", "--constraint", ">=1.0.0,<2.0.0", "--version", "1.5.0"}, nil, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "resolved: 1.5.0") {
		t.Fatalf("expected satisfied, got %q", out.String())
	}
}

func TestRun_Check_NotSatisfies(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{"check", "--constraint", ">=2.0.0", "--version", "1.5.0"}, nil, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "no version satisfies") {
		t.Fatalf("expected not satisfied, got %q", out.String())
	}
}
