package semver

import "sort"

// Range represents a closed interval [Min, Max] of versions. If either bound
// is the zero Version and Unbounded* is true, that side has no limit.
type Range struct {
	Min          Version
	Max          Version
	UnboundedMin bool
	UnboundedMax bool
}

// Contains reports whether version v falls within the range [Min, Max].
func (r Range) Contains(v Version) bool {
	if !r.UnboundedMin && Compare(v, r.Min) < 0 {
		return false
	}
	if !r.UnboundedMax && Compare(v, r.Max) > 0 {
		return false
	}
	return true
}

// IsEmpty reports whether the range cannot contain any version (Min > Max).
func (r Range) IsEmpty() bool {
	if r.UnboundedMin || r.UnboundedMax {
		return false
	}
	return Compare(r.Min, r.Max) > 0
}

// Intersect returns the intersection of two ranges, or an empty range if they
// do not overlap.
func Intersect(a, b Range) Range {
	var result Range
	// Determine the higher minimum.
	if a.UnboundedMin {
		result.Min = b.Min
		result.UnboundedMin = b.UnboundedMin
	} else if b.UnboundedMin {
		result.Min = a.Min
		result.UnboundedMin = a.UnboundedMin
	} else {
		if Compare(a.Min, b.Min) > 0 {
			result.Min = a.Min
		} else {
			result.Min = b.Min
		}
	}
	// Determine the lower maximum.
	if a.UnboundedMax {
		result.Max = b.Max
		result.UnboundedMax = b.UnboundedMax
	} else if b.UnboundedMax {
		result.Max = a.Max
		result.UnboundedMax = a.UnboundedMax
	} else {
		if Compare(a.Max, b.Max) < 0 {
			result.Max = a.Max
		} else {
			result.Max = b.Max
		}
	}
	return result
}

// Union returns the smallest range that contains both a and b. Note: this is
// only exact if the ranges overlap or are adjacent; otherwise it may include
// versions between them.
func Union(a, b Range) Range {
	var result Range
	if a.UnboundedMin || b.UnboundedMin {
		result.UnboundedMin = true
	} else {
		if Compare(a.Min, b.Min) < 0 {
			result.Min = a.Min
		} else {
			result.Min = b.Min
		}
	}
	if a.UnboundedMax || b.UnboundedMax {
		result.UnboundedMax = true
	} else {
		if Compare(a.Max, b.Max) > 0 {
			result.Max = a.Max
		} else {
			result.Max = b.Max
		}
	}
	return result
}

// Overlaps reports whether two ranges share at least one version.
func Overlaps(a, b Range) bool {
	inter := Intersect(a, b)
	return !inter.IsEmpty()
}

// FilterRange returns all versions from candidates that fall within the range.
func FilterRange(candidates []Version, r Range) []Version {
	var out []Version
	for _, v := range candidates {
		if r.Contains(v) {
			out = append(out, v)
		}
	}
	return out
}

// Sort sorts a slice of versions in ascending order.
func Sort(versions []Version) {
	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) < 0
	})
}

// SortDescending sorts a slice of versions in descending order.
func SortDescending(versions []Version) {
	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) > 0
	})
}

// Unique removes duplicate versions from a sorted slice.
func Unique(versions []Version) []Version {
	if len(versions) <= 1 {
		return versions
	}
	out := []Version{versions[0]}
	for i := 1; i < len(versions); i++ {
		if Compare(versions[i], out[len(out)-1]) != 0 {
			out = append(out, versions[i])
		}
	}
	return out
}

// Latest returns the highest version from a slice, or zero Version if empty.
func Latest(versions []Version) (Version, bool) {
	if len(versions) == 0 {
		return Version{}, false
	}
	best := versions[0]
	for _, v := range versions[1:] {
		if Compare(v, best) > 0 {
			best = v
		}
	}
	return best, true
}

// Oldest returns the lowest version from a slice, or zero Version if empty.
func Oldest(versions []Version) (Version, bool) {
	if len(versions) == 0 {
		return Version{}, false
	}
	best := versions[0]
	for _, v := range versions[1:] {
		if Compare(v, best) < 0 {
			best = v
		}
	}
	return best, true
}
