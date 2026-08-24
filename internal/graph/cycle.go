package graph

// DetectCycles returns all cycles found in the graph. Each cycle is a list of
// node names forming the cycle (last element points back to first).
func (g *Graph) DetectCycles() [][]string {
	white := make(map[string]bool, len(g.nodes))
	gray := make(map[string]bool)
	black := make(map[string]bool)
	for name := range g.nodes {
		white[name] = true
	}

	var cycles [][]string
	var path []string

	var dfs func(string)
	dfs = func(name string) {
		delete(white, name)
		gray[name] = true
		path = append(path, name)

		for _, dep := range g.nodes[name].Edges {
			if gray[dep] {
				// Found a cycle: extract it from path.
				start := -1
				for i, p := range path {
					if p == dep {
						start = i
						break
					}
				}
				if start >= 0 {
					cycle := make([]string, len(path)-start)
					copy(cycle, path[start:])
					cycles = append(cycles, cycle)
				}
			} else if white[dep] {
				dfs(dep)
			}
		}

		path = path[:len(path)-1]
		delete(gray, name)
		black[name] = true
	}

	for _, name := range g.order {
		if white[name] {
			dfs(name)
		}
	}
	return cycles
}

// HasCycle reports whether the graph contains any cycle.
func (g *Graph) HasCycle() bool {
	_, err := g.TopologicalSort()
	return err != nil
}

// FindCyclePath returns the first cycle found as a path, or nil if acyclic.
func (g *Graph) FindCyclePath() []string {
	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		return nil
	}
	return cycles[0]
}
