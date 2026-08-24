package semver

// ConstraintIntersect computes the intersection of two constraints by
// combining their alt clauses. A version must satisfy both constraints.
func ConstraintIntersect(a, b Constraint) Constraint {
	if len(a.alts) == 0 {
		return b
	}
	if len(b.alts) == 0 {
		return a
	}
	// Intersection: every alt from a must be ANDed with every alt from b.
	var newAlts [][]cmp
	for _, aAlt := range a.alts {
		for _, bAlt := range b.alts {
			combined := make([]cmp, 0, len(aAlt)+len(bAlt))
			combined = append(combined, aAlt...)
			combined = append(combined, bAlt...)
			newAlts = append(newAlts, combined)
		}
	}
	return Constraint{alts: newAlts}
}

// ConstraintSatisfiesAny reports whether any version in candidates satisfies
// the constraint.
func ConstraintSatisfiesAny(c Constraint, candidates []Version) bool {
	for _, v := range candidates {
		if c.Satisfies(v) {
			return true
		}
	}
	return false
}

// ConstraintSatisfiesAll reports whether all candidates satisfy the constraint.
func ConstraintSatisfiesAll(c Constraint, candidates []Version) bool {
	for _, v := range candidates {
		if !c.Satisfies(v) {
			return false
		}
	}
	return true
}

// FilterSatisfying returns all versions from candidates that satisfy the
// constraint.
func FilterSatisfying(c Constraint, candidates []Version) []Version {
	var out []Version
	for _, v := range candidates {
		if c.Satisfies(v) {
			out = append(out, v)
		}
	}
	return out
}

// ConstraintIsEmpty returns true if no version can satisfy the constraint
// (determined by testing against a broad set of versions). This is an
// approximation — it tests a representative set.
func ConstraintIsEmpty(c Constraint) bool {
	test := []Version{
		{Major: 0, Minor: 0, Patch: 0},
		{Major: 0, Minor: 0, Patch: 1},
		{Major: 0, Minor: 1, Patch: 0},
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 1, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
		{Major: 5, Minor: 0, Patch: 0},
		{Major: 10, Minor: 0, Patch: 0},
		{Major: 100, Minor: 0, Patch: 0},
	}
	return !ConstraintSatisfiesAny(c, test)
}

// ConstraintString returns a simplified string representation of a constraint.
func (c Constraint) String() string {
	if len(c.alts) == 0 {
		return "*"
	}
	var parts []string
	for _, alt := range c.alts {
		var altParts []string
		for _, cm := range alt {
			altParts = append(altParts, cm.op+cm.v.String())
		}
		if len(altParts) == 0 {
			parts = append(parts, "*")
		} else {
			joined := ""
			for i, p := range altParts {
				if i > 0 {
					joined += " "
				}
				joined += p
			}
			parts = append(parts, joined)
		}
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " || "
		}
		result += p
	}
	return result
}
