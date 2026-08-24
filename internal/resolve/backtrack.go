package resolve

import (
	"semver-resolver/internal/semver"
)

// BacktrackSolver performs constraint-based resolution using backtracking.
// It attempts to find a consistent assignment of versions to all packages such
// that all inter-package constraints are satisfied simultaneously.
type BacktrackSolver struct {
	Packages   []Package
	Deps       []Dependency
	MaxTries   int // maximum attempts before giving up
	depMap     map[string][]Dependency
	pkgMap     map[string]*Package
	assignment map[string]semver.Version
	tried      int
}

// NewBacktrackSolver creates a solver with default settings.
func NewBacktrackSolver(packages []Package, deps []Dependency) *BacktrackSolver {
	depMap := make(map[string][]Dependency)
	for _, d := range deps {
		depMap[d.From] = append(depMap[d.From], d)
	}
	pkgMap := make(map[string]*Package, len(packages))
	for i := range packages {
		pkgMap[packages[i].Name] = &packages[i]
	}
	return &BacktrackSolver{
		Packages:   packages,
		Deps:       deps,
		MaxTries:   10000,
		depMap:     depMap,
		pkgMap:     pkgMap,
		assignment: make(map[string]semver.Version),
	}
}

// Solve attempts to find a consistent resolution. Returns nil resolution if
// no solution exists within MaxTries.
func (s *BacktrackSolver) Solve() *Resolution {
	s.tried = 0
	if s.backtrack(0) {
		res := &Resolution{Resolved: make(map[string]semver.Version)}
		for k, v := range s.assignment {
			res.Resolved[k] = v
		}
		return res
	}
	return nil
}

func (s *BacktrackSolver) backtrack(idx int) bool {
	if idx >= len(s.Packages) {
		return true // all assigned
	}
	s.tried++
	if s.tried > s.MaxTries {
		return false
	}

	pkg := s.Packages[idx]
	// Try candidates in descending order.
	candidates := make([]semver.Version, len(pkg.Candidates))
	copy(candidates, pkg.Candidates)
	semver.SortDescending(candidates)

	for _, v := range candidates {
		if !pkg.Constraint.Satisfies(v) {
			continue
		}
		if !s.consistent(pkg.Name, v) {
			continue
		}
		s.assignment[pkg.Name] = v
		if s.backtrack(idx + 1) {
			return true
		}
		delete(s.assignment, pkg.Name)
	}
	return false
}

// consistent checks whether assigning version v to pkg is consistent with
// all currently assigned packages.
func (s *BacktrackSolver) consistent(pkgName string, v semver.Version) bool {
	// Check constraints from already-assigned packages that depend on pkgName.
	for assignedPkg := range s.assignment {
		for _, d := range s.depMap[assignedPkg] {
			if d.To == pkgName {
				if !d.Constraint.Satisfies(v) {
					return false
				}
			}
		}
	}
	// Check constraints from pkgName to already-assigned packages.
	for _, d := range s.depMap[pkgName] {
		if assignedV, ok := s.assignment[d.To]; ok {
			if !d.Constraint.Satisfies(assignedV) {
				return false
			}
		}
	}
	return true
}

// Attempts returns the number of tries the solver used.
func (s *BacktrackSolver) Attempts() int { return s.tried }
