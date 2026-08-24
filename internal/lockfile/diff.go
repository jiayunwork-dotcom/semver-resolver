package lockfile

import (
	"fmt"
	"strings"
)

// ChangeType classifies a difference between two lock files.
type ChangeType int

const (
	Added   ChangeType = iota // package not in old, present in new
	Removed                   // package in old, not in new
	Updated                   // version changed
)

// Change describes a single difference between two lock files.
type Change struct {
	Type       ChangeType
	Name       string
	OldVersion string // empty for Added
	NewVersion string // empty for Removed
}

// Diff computes the differences between two lock files.
func Diff(old, new *LockFile) []Change {
	var changes []Change

	oldMap := make(map[string]string, len(old.Entries))
	for _, e := range old.Entries {
		oldMap[e.Name] = e.Version
	}
	newMap := make(map[string]string, len(new.Entries))
	for _, e := range new.Entries {
		newMap[e.Name] = e.Version
	}

	// Check for removed and updated.
	for name, oldVer := range oldMap {
		newVer, ok := newMap[name]
		if !ok {
			changes = append(changes, Change{Type: Removed, Name: name, OldVersion: oldVer})
		} else if oldVer != newVer {
			changes = append(changes, Change{Type: Updated, Name: name, OldVersion: oldVer, NewVersion: newVer})
		}
	}
	// Check for added.
	for name, newVer := range newMap {
		if _, ok := oldMap[name]; !ok {
			changes = append(changes, Change{Type: Added, Name: name, NewVersion: newVer})
		}
	}
	return changes
}

// FormatDiff renders changes as a human-readable string.
func FormatDiff(changes []Change) string {
	if len(changes) == 0 {
		return "no changes\n"
	}
	var b strings.Builder
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(&b, "+ %s@%s\n", c.Name, c.NewVersion)
		case Removed:
			fmt.Fprintf(&b, "- %s@%s\n", c.Name, c.OldVersion)
		case Updated:
			fmt.Fprintf(&b, "~ %s: %s -> %s\n", c.Name, c.OldVersion, c.NewVersion)
		}
	}
	return b.String()
}

// HasBreakingChanges reports whether any change is a major version bump
// (potentially breaking).
func HasBreakingChanges(changes []Change) bool {
	for _, c := range changes {
		if c.Type != Updated {
			continue
		}
		// Compare major versions.
		oldMajor := majorFromString(c.OldVersion)
		newMajor := majorFromString(c.NewVersion)
		if newMajor > oldMajor && oldMajor > 0 {
			return true
		}
	}
	return false
}

func majorFromString(version string) int {
	for i, c := range version {
		if c == '.' {
			n := 0
			for _, d := range version[:i] {
				if d >= '0' && d <= '9' {
					n = n*10 + int(d-'0')
				}
			}
			return n
		}
		_ = i
	}
	return 0
}

// CountByType counts changes by type.
func CountByType(changes []Change) (added, removed, updated int) {
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Updated:
			updated++
		}
	}
	return
}
