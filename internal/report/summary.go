package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"semver-resolver/internal/semver"
)

// MultiResult summarizes a multi-package resolution.
type MultiResult struct {
	Resolved  map[string]string `json:"resolved"`
	Conflicts []ConflictInfo    `json:"conflicts,omitempty"`
	Strategy  string            `json:"strategy"`
}

// ConflictInfo describes one conflict in a resolution.
type ConflictInfo struct {
	Package   string   `json:"package"`
	Requester string   `json:"requester,omitempty"`
	Versions  []string `json:"available_versions"`
}

// RenderMultiText writes a human-readable multi-resolution summary.
func (r *MultiResult) RenderMultiText(w io.Writer) {
	fmt.Fprintf(w, "resolution strategy: %s\n", r.Strategy)
	if len(r.Resolved) == 0 && len(r.Conflicts) == 0 {
		fmt.Fprintln(w, "  (empty)")
		return
	}
	if len(r.Resolved) > 0 {
		fmt.Fprintln(w, "resolved:")
		for name, version := range r.Resolved {
			fmt.Fprintf(w, "  %s@%s\n", name, version)
		}
	}
	if len(r.Conflicts) > 0 {
		fmt.Fprintln(w, "conflicts:")
		for _, c := range r.Conflicts {
			fmt.Fprintf(w, "  %s: no satisfying version from [%s]\n", c.Package, strings.Join(c.Versions, ", "))
		}
	}
}

// RenderMultiJSON writes the multi-resolution result as JSON.
func (r *MultiResult) RenderMultiJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// VersionListResult summarizes a version list query (e.g., all satisfying versions).
type VersionListResult struct {
	Constraint string   `json:"constraint"`
	Versions   []string `json:"versions"`
	Count      int      `json:"count"`
}

// RenderVersionList writes a version list as text.
func (r *VersionListResult) RenderVersionList(w io.Writer) {
	fmt.Fprintf(w, "constraint: %s\n", r.Constraint)
	fmt.Fprintf(w, "matching versions (%d):\n", r.Count)
	for _, v := range r.Versions {
		fmt.Fprintf(w, "  %s\n", v)
	}
}

// FormatVersionSlice converts a Version slice to a string slice.
func FormatVersionSlice(vs []semver.Version) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = v.String()
	}
	return out
}
