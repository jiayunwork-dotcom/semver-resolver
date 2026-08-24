package graph

// Stats holds summary statistics about a dependency graph.
type Stats struct {
	NodeCount      int
	EdgeCount      int
	RootCount      int
	LeafCount      int
	MaxDepth       int
	AvgFanOut      float64
	MaxFanOut      int
	IsolatedCount  int // nodes with no edges in either direction
}

// ComputeStats computes summary statistics for the graph.
func (g *Graph) ComputeStats() Stats {
	nodeCount := len(g.nodes)
	edgeCount := 0
	maxFanOut := 0
	for _, n := range g.nodes {
		edgeCount += len(n.Edges)
		if len(n.Edges) > maxFanOut {
			maxFanOut = len(n.Edges)
		}
	}

	avgFanOut := 0.0
	if nodeCount > 0 {
		avgFanOut = float64(edgeCount) / float64(nodeCount)
	}

	hasIncoming := make(map[string]bool)
	for _, n := range g.nodes {
		for _, e := range n.Edges {
			hasIncoming[e] = true
		}
	}

	roots := 0
	leaves := 0
	isolated := 0
	for _, name := range g.order {
		n := g.nodes[name]
		isRoot := !hasIncoming[name]
		isLeaf := len(n.Edges) == 0
		if isRoot {
			roots++
		}
		if isLeaf {
			leaves++
		}
		if isRoot && isLeaf {
			isolated++
		}
	}

	return Stats{
		NodeCount:     nodeCount,
		EdgeCount:     edgeCount,
		RootCount:     roots,
		LeafCount:     leaves,
		MaxDepth:      g.MaxDepth(),
		AvgFanOut:     avgFanOut,
		MaxFanOut:     maxFanOut,
		IsolatedCount: isolated,
	}
}

// Density returns the edge density of the graph (edges / max possible edges).
func (g *Graph) Density() float64 {
	n := len(g.nodes)
	if n <= 1 {
		return 0
	}
	maxEdges := n * (n - 1)
	edgeCount := 0
	for _, node := range g.nodes {
		edgeCount += len(node.Edges)
	}
	return float64(edgeCount) / float64(maxEdges)
}

// FanOut returns the number of direct dependencies for a node.
func (g *Graph) FanOut(name string) int {
	n := g.nodes[name]
	if n == nil {
		return 0
	}
	return len(n.Edges)
}

// FanIn returns the number of direct dependents for a node.
func (g *Graph) FanIn(name string) int {
	count := 0
	for _, n := range g.nodes {
		for _, e := range n.Edges {
			if e == name {
				count++
				break
			}
		}
	}
	return count
}
