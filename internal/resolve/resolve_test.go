package resolve

import (
	"testing"

	"semver-resolver/internal/semver"
)

func TestHighest(t *testing.T) {
	c, _ := semver.ParseConstraint("^1.2.0")
	cands := []semver.Version{
		{Major: 1, Minor: 1, Patch: 0},
		{Major: 1, Minor: 2, Patch: 0},
		{Major: 1, Minor: 3, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
	}
	v, ok := Highest(c, cands)
	if !ok || v.Minor != 3 {
		t.Fatalf("expected 1.3.0, got %v (found=%v)", v, ok)
	}
}

func TestParseCandidates(t *testing.T) {
	vs, err := ParseCandidates("1.0.0, 2.0.0, 3.0.0-beta.1")
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) != 3 {
		t.Fatalf("expected 3 candidates, got %d", len(vs))
	}
}

func TestResolveAll(t *testing.T) {
	c1, _ := semver.ParseConstraint("^1.0.0")
	c2, _ := semver.ParseConstraint("^2.0.0")
	cDep, _ := semver.ParseConstraint(">=1.1.0")

	packages := []Package{
		{Name: "app", Constraint: c1, Candidates: []semver.Version{
			{Major: 1, Minor: 0, Patch: 0}, {Major: 1, Minor: 1, Patch: 0}, {Major: 1, Minor: 2, Patch: 0},
		}},
		{Name: "lib", Constraint: c2, Candidates: []semver.Version{
			{Major: 2, Minor: 0, Patch: 0}, {Major: 2, Minor: 1, Patch: 0},
		}},
	}
	deps := []Dependency{
		{From: "app", To: "lib", Constraint: cDep},
	}
	res := ResolveAll(packages, deps)
	if len(res.Conflicts) > 0 {
		t.Fatalf("unexpected conflicts: %+v", res.Conflicts)
	}
	if res.Resolved["app"].String() != "1.2.0" {
		t.Fatalf("app = %s, want 1.2.0", res.Resolved["app"])
	}
}

func TestMinimalUpgrade(t *testing.T) {
	c1, _ := semver.ParseConstraint(">=1.1.0")
	packages := []Package{
		{Name: "foo", Constraint: c1, Candidates: []semver.Version{
			{Major: 1, Minor: 0, Patch: 0}, {Major: 1, Minor: 1, Patch: 0},
			{Major: 1, Minor: 2, Patch: 0}, {Major: 1, Minor: 3, Patch: 0},
		}},
	}
	current := map[string]semver.Version{"foo": {Major: 1, Minor: 1, Patch: 0}}
	res := MinimalUpgrade(packages, nil, current)
	if res.Resolved["foo"].String() != "1.1.0" {
		t.Fatalf("expected 1.1.0 (keep current), got %s", res.Resolved["foo"])
	}
}

func TestBacktrackSolver(t *testing.T) {
	c1, _ := semver.ParseConstraint("^1.0.0")
	c2, _ := semver.ParseConstraint(">=1.2.0")
	cDep, _ := semver.ParseConstraint("^1.0.0")

	packages := []Package{
		{Name: "A", Constraint: c1, Candidates: []semver.Version{
			{Major: 1, Minor: 0, Patch: 0}, {Major: 1, Minor: 1, Patch: 0}, {Major: 1, Minor: 2, Patch: 0},
		}},
		{Name: "B", Constraint: c2, Candidates: []semver.Version{
			{Major: 1, Minor: 2, Patch: 0}, {Major: 1, Minor: 3, Patch: 0},
		}},
	}
	deps := []Dependency{
		{From: "A", To: "B", Constraint: cDep},
	}
	solver := NewBacktrackSolver(packages, deps)
	res := solver.Solve()
	if res == nil {
		t.Fatal("expected a solution")
	}
	if res.Resolved["A"].String() != "1.2.0" {
		t.Fatalf("A = %s, want 1.2.0", res.Resolved["A"])
	}
	if res.Resolved["B"].Minor < 2 {
		t.Fatalf("B = %s, want >= 1.2.0", res.Resolved["B"])
	}
}

func TestSelectVersionStrategy(t *testing.T) {
	c, _ := semver.ParseConstraint(">=1.0.0")
	cands := []semver.Version{
		{Major: 1, Minor: 0, Patch: 0}, {Major: 1, Minor: 1, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
	}
	// Lowest strategy.
	v, ok := SelectVersion(c, cands, StrategyLowest, nil)
	if !ok || v.String() != "1.0.0" {
		t.Fatalf("lowest = %v, want 1.0.0", v)
	}
	// Highest strategy.
	v, ok = SelectVersion(c, cands, StrategyHighest, nil)
	if !ok || v.String() != "2.0.0" {
		t.Fatalf("highest = %v, want 2.0.0", v)
	}
}

func TestFilterStable(t *testing.T) {
	cands := []semver.Version{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 1, Patch: 0, Pre: "beta.1"},
		{Major: 2, Minor: 0, Patch: 0},
	}
	stable := FilterStable(cands)
	if len(stable) != 2 {
		t.Fatalf("expected 2 stable, got %d", len(stable))
	}
}

func TestValidateResolution(t *testing.T) {
	c, _ := semver.ParseConstraint(">=2.0.0")
	res := Resolution{Resolved: map[string]semver.Version{
		"A": {Major: 1, Minor: 0, Patch: 0},
		"B": {Major: 1, Minor: 5, Patch: 0},
	}}
	deps := []Dependency{{From: "A", To: "B", Constraint: c}}
	errs := ValidateResolution(res, deps)
	if len(errs) == 0 {
		t.Fatal("expected validation error: B=1.5.0 doesn't satisfy >=2.0.0")
	}
}
