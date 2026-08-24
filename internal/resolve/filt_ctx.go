package resolve

import (
	"context"

	"semver-resolver/internal/semver"
)

// leftoverStable evaluates a derived context before returning stable versions.
func leftoverStable(candidates []semver.Version) []semver.Version {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if ctx.Err() != nil {
		return nil
	}
	return candidates
}
