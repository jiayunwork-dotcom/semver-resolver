package graph

import "testing"

func makeTestGraph() *Graph {
	g := New()
	g.Add("app", "1.0.0")
	g.Add("web", "2.1.0")
	g.Add("db", "3.0.0")
	g.Add("log", "1.2.0")
	g.Add("util", "0.5.0")
	g.AddEdge("app", "web")
	g.AddEdge("app", "db")
	g.AddEdge("web", "log")
	g.AddEdge("web", "util")
	g.AddEdge("db", "log")
	return g
}

func TestTopologicalSort(t *testing.T) {
	g := makeTestGraph()
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 5 {
		t.Fatalf("expected 5, got %d", len(order))
	}
	// app must come before web and db.
	pos := make(map[string]int)
	for i, n := range order {
		pos[n] = i
	}
	if pos["app"] >= pos["web"] {
		t.Fatal("app should appear before web")
	}
	if pos["web"] >= pos["log"] {
		t.Fatal("web should appear before log")
	}
}

func TestCycleDetection(t *testing.T) {
	g := New()
	g.Add("A", "1.0.0")
	g.Add("B", "1.0.0")
	g.Add("C", "1.0.0")
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")
	if !g.HasCycle() {
		t.Fatal("expected cycle")
	}
	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Fatal("expected at least one cycle")
	}
}

func TestNoCycle(t *testing.T) {
	g := makeTestGraph()
	if g.HasCycle() {
		t.Fatal("should not have cycle")
	}
}

func TestTransitiveDeps(t *testing.T) {
	g := makeTestGraph()
	deps := g.TransitiveDeps("app")
	// app -> web, db, log, util
	if len(deps) != 4 {
		t.Fatalf("expected 4 transitive deps, got %d: %v", len(deps), deps)
	}
}

func TestTransitiveReverseDeps(t *testing.T) {
	g := makeTestGraph()
	revDeps := g.TransitiveReverseDeps("log")
	// log is depended on by web, db, and transitively by app.
	if len(revDeps) != 3 {
		t.Fatalf("expected 3 reverse deps, got %d: %v", len(revDeps), revDeps)
	}
}

func TestRootsAndLeaves(t *testing.T) {
	g := makeTestGraph()
	roots := g.Roots()
	if len(roots) != 1 || roots[0] != "app" {
		t.Fatalf("roots = %v, want [app]", roots)
	}
	leaves := g.Leaves()
	if len(leaves) != 2 {
		t.Fatalf("leaves = %v, want [log, util]", leaves)
	}
}

func TestInstallOrder(t *testing.T) {
	g := makeTestGraph()
	order, err := g.InstallOrder()
	if err != nil {
		t.Fatal(err)
	}
	pos := make(map[string]int)
	for i, n := range order {
		pos[n] = i
	}
	// log and util should be installed before web.
	if pos["log"] >= pos["web"] {
		t.Fatalf("log should install before web: %v", order)
	}
}

func TestMaxDepth(t *testing.T) {
	g := makeTestGraph()
	d := g.MaxDepth()
	// app -> web -> log is depth 2.
	if d != 2 {
		t.Fatalf("max depth = %d, want 2", d)
	}
}

func TestSharedDeps(t *testing.T) {
	g := makeTestGraph()
	shared := g.SharedDeps("web", "db")
	// Both web and db depend on log.
	if len(shared) != 1 || shared[0] != "log" {
		t.Fatalf("shared = %v, want [log]", shared)
	}
}
