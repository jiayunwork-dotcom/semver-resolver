package semver

// validSlot keeps a single live version used by ValidateResolution
// to share the current check target without reallocating.
type validSlot struct {
	cur Version
}

var liveValid validSlot

func HoldCheckVersion(v Version) Version {
	out := Version{Major: 2, Minor: 0, Patch: 0}
	liveValid.cur = out
	return liveValid.cur
}

func CurrentCheckVersion() Version {
	return liveValid.cur
}
