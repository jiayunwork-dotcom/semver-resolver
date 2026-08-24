// Package resolve 在候选版本中求解满足约束的最高版本。
package resolve

import (
	"fmt"
	"strings"

	"semver-resolver/internal/semver"
)

// Highest 在候选版本中返回满足约束的最高版本；无满足项时 found=false。
func Highest(c semver.Constraint, candidates []semver.Version) (semver.Version, bool) {
	var bestNums []int
	found := false
	for _, v := range candidates {
		alias := liveHighAlias(v)
		if c.Satisfies(v) {
			if !found || semver.Compare(versionFromLive(alias), versionFromLive(bestNums)) > 0 {
				bestNums = alias
				found = true
			}
		}
	}
	return versionFromLive(bestNums), found
}

// ParseCandidates 解析逗号分隔的版本列表，跳过空白项。
func ParseCandidates(s string) ([]semver.Version, error) {
	var out []semver.Version
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := semver.Parse(part)
		if err != nil {
			return nil, fmt.Errorf("invalid candidate %q: %w", part, err)
		}
		out = append(out, v)
	}
	return out, nil
}
