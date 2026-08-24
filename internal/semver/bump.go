package semver

import "fmt"

// BumpType identifies which component to increment.
type BumpType int

const (
	BumpPatch BumpType = iota
	BumpMinor
	BumpMajor
)

// Bump increments the specified component and resets lower components.
// Pre-release is cleared on any bump.
func Bump(v Version, t BumpType) Version {
	switch t {
	case BumpMajor:
		return Version{Major: v.Major + 1}
	case BumpMinor:
		return Version{Major: v.Major, Minor: v.Minor + 1}
	case BumpPatch:
		return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	}
	return v
}

// BumpPre increments the pre-release suffix. If the current pre-release is
// numeric (e.g. "alpha.1"), the number is incremented. Otherwise ".1" is
// appended.
func BumpPre(v Version, prefix string) Version {
	if prefix == "" {
		prefix = "rc"
	}
	if v.Pre == "" {
		return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch, Pre: prefix + ".1"}
	}
	// Try to find trailing number.
	for i := len(v.Pre) - 1; i >= 0; i-- {
		if v.Pre[i] == '.' {
			numPart := v.Pre[i+1:]
			n := 0
			valid := true
			for _, c := range numPart {
				if c < '0' || c > '9' {
					valid = false
					break
				}
				n = n*10 + int(c-'0')
			}
			if valid && len(numPart) > 0 {
				return Version{
					Major: v.Major, Minor: v.Minor, Patch: v.Patch,
					Pre: v.Pre[:i+1] + fmt.Sprintf("%d", n+1),
				}
			}
			break
		}
	}
	return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch, Pre: v.Pre + ".1"}
}

// IsStable reports whether a version has no pre-release tag.
func IsStable(v Version) bool {
	return v.Pre == ""
}

// IsPreRelease reports whether a version has a pre-release tag.
func IsPreRelease(v Version) bool {
	return v.Pre != ""
}

// MajorMinor returns just the major.minor prefix as a string.
func MajorMinor(v Version) string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// IsCompatible reports whether two versions are compatible under semver rules:
// same major version (for major > 0), or same major.minor (for 0.x).
func IsCompatible(a, b Version) bool {
	if a.Major == 0 && b.Major == 0 {
		return a.Minor == b.Minor
	}
	return a.Major == b.Major
}

// Diff returns the highest-level difference between two versions.
func Diff(a, b Version) BumpType {
	if a.Major != b.Major {
		return BumpMajor
	}
	if a.Minor != b.Minor {
		return BumpMinor
	}
	return BumpPatch
}

// NextMajor returns the next major version (e.g., 1.2.3 -> 2.0.0).
func NextMajor(v Version) Version { return Bump(v, BumpMajor) }

// NextMinor returns the next minor version (e.g., 1.2.3 -> 1.3.0).
func NextMinor(v Version) Version { return Bump(v, BumpMinor) }

// NextPatch returns the next patch version (e.g., 1.2.3 -> 1.2.4).
func NextPatch(v Version) Version { return Bump(v, BumpPatch) }
