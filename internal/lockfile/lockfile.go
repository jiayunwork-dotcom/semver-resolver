// Package lockfile handles parsing, generating, and comparing dependency
// lock files. A lock file records exact resolved versions for reproducible builds.
package lockfile

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"semver-resolver/internal/semver"
)

// Entry represents one locked dependency.
type Entry struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Constraint string `json:"constraint,omitempty"`
	Integrity  string `json:"integrity,omitempty"`
}

// LockFile is the in-memory representation of a dependency lock file.
type LockFile struct {
	FormatVersion int     `json:"format_version"`
	Entries       []Entry `json:"entries"`
}

// Parse reads a lock file from JSON.
func Parse(r io.Reader) (*LockFile, error) {
	var lf LockFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&lf); err != nil {
		return nil, fmt.Errorf("parse lockfile: %w", err)
	}
	return &lf, nil
}

// Write serializes the lock file to JSON.
func (lf *LockFile) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(lf)
}

// Sort sorts entries by name for deterministic output.
func (lf *LockFile) Sort() {
	sort.Slice(lf.Entries, func(i, j int) bool {
		return lf.Entries[i].Name < lf.Entries[j].Name
	})
}

// Get returns the entry for a given package name, or nil.
func (lf *LockFile) Get(name string) *Entry {
	for i := range lf.Entries {
		if lf.Entries[i].Name == name {
			return &lf.Entries[i]
		}
	}
	return nil
}

// Set adds or updates an entry for a given package.
func (lf *LockFile) Set(name, version, constraint string) {
	for i := range lf.Entries {
		if lf.Entries[i].Name == name {
			lf.Entries[i].Version = version
			lf.Entries[i].Constraint = constraint
			return
		}
	}
	lf.Entries = append(lf.Entries, Entry{Name: name, Version: version, Constraint: constraint})
}

// Remove deletes an entry by name.
func (lf *LockFile) Remove(name string) bool {
	for i := range lf.Entries {
		if lf.Entries[i].Name == name {
			lf.Entries = append(lf.Entries[:i], lf.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// Versions returns a map of package names to their locked versions.
func (lf *LockFile) Versions() map[string]semver.Version {
	m := make(map[string]semver.Version, len(lf.Entries))
	for _, e := range lf.Entries {
		v, err := semver.Parse(e.Version)
		if err == nil {
			m[e.Name] = v
		}
	}
	return m
}

// FromResolution creates a lock file from a resolution map.
func FromResolution(resolved map[string]semver.Version) *LockFile {
	lf := &LockFile{FormatVersion: 1}
	for name, v := range resolved {
		lf.Entries = append(lf.Entries, Entry{Name: name, Version: v.String()})
	}
	lf.Sort()
	return lf
}

// FormatText renders the lockfile as human-readable text.
func (lf *LockFile) FormatText() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# lock file (format v%d)\n", lf.FormatVersion)
	for _, e := range lf.Entries {
		fmt.Fprintf(&b, "%s@%s", e.Name, e.Version)
		if e.Constraint != "" {
			fmt.Fprintf(&b, " (%s)", e.Constraint)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
