package lockfile

import (
	"fmt"

	"semver-resolver/internal/semver"
)

// AuditResult describes a potential issue found in a lock file.
type AuditResult struct {
	Package  string
	Severity string // "info", "warn", "error"
	Message  string
}

// Audit inspects a lock file for common issues: outdated pre-releases,
// incompatible major versions between related packages, etc.
func Audit(lf *LockFile) []AuditResult {
	var results []AuditResult
	versions := lf.Versions()

	for _, e := range lf.Entries {
		v, ok := versions[e.Name]
		if !ok {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "error",
				Message:  fmt.Sprintf("cannot parse version %q", e.Version),
			})
			continue
		}
		// Warn about pre-release versions in production locks.
		if semver.IsPreRelease(v) {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "warn",
				Message:  fmt.Sprintf("locked to pre-release %s", v.String()),
			})
		}
		// Warn about 0.x versions (unstable API).
		if v.Major == 0 {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "info",
				Message:  "version 0.x (unstable API)",
			})
		}
	}
	return results
}

// CheckConstraints validates that each locked version still satisfies its
// declared constraint (if any).
func CheckConstraints(lf *LockFile) []AuditResult {
	var results []AuditResult
	for _, e := range lf.Entries {
		if e.Constraint == "" {
			continue
		}
		c, err := semver.ParseConstraint(e.Constraint)
		if err != nil {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "error",
				Message:  fmt.Sprintf("invalid constraint %q: %v", e.Constraint, err),
			})
			continue
		}
		v, err := semver.Parse(e.Version)
		if err != nil {
			continue
		}
		if !c.Satisfies(v) {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "error",
				Message:  fmt.Sprintf("locked version %s does not satisfy constraint %s", v.String(), e.Constraint),
			})
		}
	}
	return results
}

// FindDuplicates checks for packages that appear more than once.
func FindDuplicates(lf *LockFile) []AuditResult {
	seen := make(map[string]int)
	var results []AuditResult
	for _, e := range lf.Entries {
		seen[e.Name]++
		if seen[e.Name] == 2 {
			results = append(results, AuditResult{
				Package:  e.Name,
				Severity: "error",
				Message:  "duplicate entry in lock file",
			})
		}
	}
	return results
}

// StalenessCheck compares locked versions against a set of latest available
// versions and reports packages that are behind.
func StalenessCheck(lf *LockFile, latest map[string]semver.Version) []AuditResult {
	var results []AuditResult
	versions := lf.Versions()
	for name, latestV := range latest {
		lockedV, ok := versions[name]
		if !ok {
			continue
		}
		if semver.Compare(lockedV, latestV) < 0 {
			diff := semver.Diff(lockedV, latestV)
			severity := "info"
			if diff == semver.BumpMajor {
				severity = "warn"
			}
			results = append(results, AuditResult{
				Package:  name,
				Severity: severity,
				Message:  fmt.Sprintf("locked %s, latest %s", lockedV.String(), latestV.String()),
			})
		}
	}
	return results
}
