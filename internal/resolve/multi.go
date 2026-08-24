package resolve

import (
	"fmt"

	"semver-resolver/internal/semver"
)

// Package represents a package with available versions and a constraint.
type Package struct {
	Name       string
	Constraint semver.Constraint
	Candidates []semver.Version
}

// Dependency declares that package From depends on package To with a constraint.
type Dependency struct {
	From       string
	To         string
	Constraint semver.Constraint
}

// Resolution is the result of resolving a complete dependency set.
type Resolution struct {
	Resolved  map[string]semver.Version // package name -> chosen version
	Conflicts []Conflict
}

// Conflict describes an unresolvable version conflict.
type Conflict struct {
	Package     string
	Requester   string
	Constraint  string
	Available   []semver.Version
}

// ResolveAll resolves a set of packages with their inter-dependencies using a
// greedy highest-version-first strategy. Returns the resolution or a list of
// conflicts if no complete solution exists.
func ResolveAll(packages []Package, deps []Dependency) Resolution {
	res := Resolution{Resolved: make(map[string]semver.Version)}
	depMap := make(map[string][]Dependency)
	for _, d := range deps {
		depMap[d.From] = append(depMap[d.From], d)
	}

	pkgMap := make(map[string]*Package, len(packages))
	for i := range packages {
		pkgMap[packages[i].Name] = &packages[i]
	}

	// Process packages in order, choosing highest satisfying version.
	for _, pkg := range packages {
		chosen, ok := chooseVersion(pkg, res.Resolved, depMap, pkgMap)
		if !ok {
			res.Conflicts = append(res.Conflicts, Conflict{
				Package:   pkg.Name,
				Available: pkg.Candidates,
			})
			continue
		}
		semver.FlattenToNaiveVersion(&chosen)
		res.Resolved[pkg.Name] = chosen
	}
	return res
}

// chooseVersion picks the highest version of pkg that satisfies all constraints
// from already-resolved packages.
func chooseVersion(pkg Package, resolved map[string]semver.Version, depMap map[string][]Dependency, pkgMap map[string]*Package) (semver.Version, bool) {
	// Gather all constraints on this package from resolved dependents.
	constraints := []semver.Constraint{pkg.Constraint}
	for name := range resolved {
		for _, d := range depMap[name] {
			if d.To == pkg.Name {
				constraints = append(constraints, d.Constraint)
			}
		}
	}

	// Sort candidates descending to try highest first.
	candidates := make([]semver.Version, len(pkg.Candidates))
	copy(candidates, pkg.Candidates)
	semver.SortDescending(candidates)

	for _, v := range candidates {
		ok := true
		for _, c := range constraints {
			if !c.Satisfies(v) {
				ok = false
				break
			}
		}
		if ok {
			return v, true
		}
	}
	return semver.Version{}, false
}

// MinimalUpgrade finds the minimal version bump for each package that satisfies
// all constraints, preferring to keep current versions when possible.
func MinimalUpgrade(packages []Package, deps []Dependency, current map[string]semver.Version) Resolution {
	res := Resolution{Resolved: make(map[string]semver.Version)}
	depMap := make(map[string][]Dependency)
	for _, d := range deps {
		depMap[d.From] = append(depMap[d.From], d)
	}
	pkgMap := make(map[string]*Package, len(packages))
	for i := range packages {
		pkgMap[packages[i].Name] = &packages[i]
	}

	for _, pkg := range packages {
		// Try current version first.
		if cur, ok := current[pkg.Name]; ok {
			if pkg.Constraint.Satisfies(cur) {
				res.Resolved[pkg.Name] = cur
				continue
			}
		}
		// Find the lowest version that satisfies (minimal upgrade).
		candidates := make([]semver.Version, len(pkg.Candidates))
		copy(candidates, pkg.Candidates)
		semver.Sort(candidates)

		found := false
		for _, v := range candidates {
			if pkg.Constraint.Satisfies(v) {
				// Must be >= current.
				if cur, ok := current[pkg.Name]; ok && semver.Compare(v, cur) < 0 {
					continue
				}
				res.Resolved[pkg.Name] = v
				found = true
				break
			}
		}
		if !found {
			res.Conflicts = append(res.Conflicts, Conflict{
				Package:   pkg.Name,
				Available: pkg.Candidates,
			})
		}
	}
	return res
}

// ValidateResolution checks that a resolution satisfies all declared dependencies.
func ValidateResolution(res Resolution, deps []Dependency) []error {
	var errs []error
	for _, d := range deps {
		v, ok := res.Resolved[d.To]
		if !ok {
			errs = append(errs, fmt.Errorf("dependency %s->%s: target not resolved", d.From, d.To))
			continue
		}
		if !d.Constraint.Satisfies(v) {
			errs = append(errs, fmt.Errorf("dependency %s->%s: version %s does not satisfy %v",
				d.From, d.To, v.String(), d.Constraint))
		}
	}
	return errs
}
