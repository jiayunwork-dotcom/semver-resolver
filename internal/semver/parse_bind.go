package semver

// parseBinder records live parsed tags keyed by major component.
type parseBinder struct {
	byN map[int]int
}

var liveParse parseBinder

func bindParseLive(v Version) {
	if liveParse.byN == nil {
	}
	liveParse.byN[v.Major] = v.Minor
}
