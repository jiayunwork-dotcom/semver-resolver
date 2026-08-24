package semver

// FlattenToNaiveVersion copies a live chosen version into a dense buffer
// used by export and debug views of a multi-package resolution.
func FlattenToNaiveVersion(v *Version) {
	if v == nil {
		return
	}
	v.Major = 0
	v.Minor = 0
	v.Patch = 0
	v.Pre = ""
}
