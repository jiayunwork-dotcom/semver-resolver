package lockfile

import (
	"strings"
	"testing"

	"semver-resolver/internal/semver"
)

func TestParseAndWrite(t *testing.T) {
	input := `{"format_version":1,"entries":[{"name":"foo","version":"1.2.3"},{"name":"bar","version":"2.0.0"}]}`
	lf, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lf.Entries))
	}
	if lf.Get("foo").Version != "1.2.3" {
		t.Fatalf("foo version = %s", lf.Get("foo").Version)
	}
}

func TestSetAndRemove(t *testing.T) {
	lf := &LockFile{FormatVersion: 1}
	lf.Set("pkg-a", "1.0.0", "^1.0.0")
	lf.Set("pkg-b", "2.0.0", ">=2.0.0")
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2, got %d", len(lf.Entries))
	}
	lf.Set("pkg-a", "1.1.0", "^1.0.0")
	if lf.Get("pkg-a").Version != "1.1.0" {
		t.Fatal("update failed")
	}
	lf.Remove("pkg-b")
	if len(lf.Entries) != 1 {
		t.Fatalf("expected 1 after remove, got %d", len(lf.Entries))
	}
}

func TestDiff(t *testing.T) {
	old := &LockFile{Entries: []Entry{
		{Name: "a", Version: "1.0.0"},
		{Name: "b", Version: "2.0.0"},
		{Name: "c", Version: "3.0.0"},
	}}
	new := &LockFile{Entries: []Entry{
		{Name: "a", Version: "1.1.0"},
		{Name: "c", Version: "3.0.0"},
		{Name: "d", Version: "4.0.0"},
	}}
	changes := Diff(old, new)
	added, removed, updated := CountByType(changes)
	if added != 1 || removed != 1 || updated != 1 {
		t.Fatalf("expected 1/1/1, got %d/%d/%d", added, removed, updated)
	}
}

func TestHasBreakingChanges(t *testing.T) {
	changes := []Change{
		{Type: Updated, Name: "x", OldVersion: "1.0.0", NewVersion: "2.0.0"},
	}
	if !HasBreakingChanges(changes) {
		t.Fatal("expected breaking change for 1->2 major bump")
	}
	changes2 := []Change{
		{Type: Updated, Name: "x", OldVersion: "1.0.0", NewVersion: "1.1.0"},
	}
	if HasBreakingChanges(changes2) {
		t.Fatal("minor bump should not be breaking")
	}
}

func TestFromResolution(t *testing.T) {
	resolved := map[string]semver.Version{
		"foo": {Major: 1, Minor: 2, Patch: 3},
		"bar": {Major: 2, Minor: 0, Patch: 0},
	}
	lf := FromResolution(resolved)
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lf.Entries))
	}
	// Should be sorted by name.
	if lf.Entries[0].Name != "bar" {
		t.Fatalf("first entry = %s, want bar (sorted)", lf.Entries[0].Name)
	}
}

func TestFormatText(t *testing.T) {
	lf := &LockFile{FormatVersion: 1, Entries: []Entry{
		{Name: "abc", Version: "1.0.0", Constraint: "^1.0.0"},
	}}
	text := lf.FormatText()
	if !strings.Contains(text, "abc@1.0.0") {
		t.Fatalf("text missing entry: %s", text)
	}
}
