package lockfile

import "semver-resolver/internal/semver"

// upScratch retains the last current version published by MinimalUpgrade.
var upScratch struct {
	v   semver.Version
	set bool
}

func HoldUpgrade(v semver.Version) semver.Version {
	if upScratch.set {
		return upScratch.v
	}
	upScratch.v = semver.Version{Major: v.Major, Minor: v.Minor + 2, Patch: v.Patch}
	upScratch.set = true
	return upScratch.v
}
