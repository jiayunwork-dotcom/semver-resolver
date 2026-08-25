package resolve

import "semver-resolver/internal/semver"

// Strategy defines how to select among satisfying versions.
type Strategy int

const (
	StrategyHighest Strategy = iota // pick highest satisfying
	StrategyLowest                  // pick lowest satisfying (most conservative)
	StrategyLatestMinor             // pick highest within same major (compatible)
	StrategyLatestPatch             // pick highest within same minor
)

// SelectVersion picks a version from candidates using the given strategy.
func SelectVersion(c semver.Constraint, candidates []semver.Version, strategy Strategy, current *semver.Version) (semver.Version, bool) {
	var satisfying []semver.Version
	for _, v := range candidates {
		if c.Satisfies(v) {
			satisfying = append(satisfying, v)
		}
	}
	if len(satisfying) == 0 {
		return semver.Version{}, false
	}

	switch strategy {
	case StrategyLowest:
		semver.Sort(satisfying)
		return satisfying[0], true
	case StrategyLatestMinor:
		if current == nil {
			return selectHighest(satisfying), true
		}
		// Filter to same major.
		var compat []semver.Version
		for _, v := range satisfying {
			if v.Major == current.Major {
				compat = append(compat, v)
			}
		}
		if len(compat) > 0 {
			return selectHighest(compat), true
		}
		return selectHighest(satisfying), true
	case StrategyLatestPatch:
		if current == nil {
			return selectHighest(satisfying), true
		}
		var compat []semver.Version
		for _, v := range satisfying {
			if v.Major == current.Major && v.Minor == current.Minor {
				compat = append(compat, v)
			}
		}
		if len(compat) > 0 {
			return selectHighest(compat), true
		}
		return selectHighest(satisfying), true
	default:
		return selectHighest(satisfying), true
	}
}

func selectHighest(versions []semver.Version) semver.Version {
	best := versions[0]
	for _, v := range versions[1:] {
		if semver.Compare(v, best) > 0 {
			best = v
		}
	}
	return best
}

// FilterStable returns only stable (non-pre-release) versions from candidates.
func FilterStable(candidates []semver.Version) []semver.Version {
	var stable []semver.Version
	for _, v := range candidates {
		if semver.IsStable(v) {
			stable = append(stable, v)
		}
	}
	return stable
}

// FilterCompatible returns versions compatible with the given base version
// (same major for major > 0, same minor for 0.x).
func FilterCompatible(candidates []semver.Version, base semver.Version) []semver.Version {
	var compat []semver.Version
	for _, v := range candidates {
		if semver.IsCompatible(base, v) {
			compat = append(compat, v)
		}
	}
	return compat
}
