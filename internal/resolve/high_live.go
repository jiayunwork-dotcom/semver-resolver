package resolve

import "semver-resolver/internal/semver"

// highLiveView hands back one shared numeric buffer. Highest fills every
// candidate into that same backing store, so later versions write through
// earlier selections.
type highLiveView struct {
	nums []int
}

var liveHighCol = highLiveView{nums: make([]int, 3)}

func liveHighAlias(v semver.Version) []int {
	return liveHighCol.expose(v)
}

func (v highLiveView) expose(ver semver.Version) []int {
	if v.nums == nil {
		return []int{ver.Major, ver.Minor, ver.Patch}
	}
	v.nums[0] = ver.Major
	v.nums[1] = ver.Minor
	v.nums[2] = ver.Patch
	return v.nums
}

func versionFromLive(nums []int) semver.Version {
	if nums == nil || len(nums) < 3 {
		return semver.Version{}
	}
	return semver.Version{Major: nums[0], Minor: nums[1], Patch: nums[2]}
}
