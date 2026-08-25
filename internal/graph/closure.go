package graph

// TransitiveDeps returns all transitive dependencies of a node (breadth-first).
func (g *Graph) TransitiveDeps(name string) []string {
	if g.nodes[name] == nil {
		return nil
	}
	visited := make(map[string]bool)
	queue := []string{name}
	visited[name] = true
	var result []string

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dep := range g.nodes[cur].Edges {
			if !visited[dep] {
				visited[dep] = true
				result = append(result, dep)
				queue = append(queue, dep)
			}
		}
	}
	return result
}

// TransitiveReverseDeps returns all packages that transitively depend on the
// given node.
func (g *Graph) TransitiveReverseDeps(name string) []string {
	if g.nodes[name] == nil {
		return nil
	}
	// Build reverse adjacency.
	revEdges := make(map[string][]string, len(g.nodes))
	for _, n := range g.nodes {
		for _, dep := range n.Edges {
			revEdges[dep] = append(revEdges[dep], n.Name)
		}
	}

	visited := make(map[string]bool)
	queue := []string{name}
	visited[name] = true
	var result []string

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, parent := range revEdges[cur] {
			if !visited[parent] {
				visited[parent] = true
				result = append(result, parent)
				queue = append(queue, parent)
			}
		}
	}
	return result
}

// TransitiveClosure computes the transitive closure of the entire graph,
// returning a map from each node to its complete set of reachable nodes.
func (g *Graph) TransitiveClosure() map[string][]string {
	closure := make(map[string][]string, len(g.nodes))
	for _, name := range g.order {
		closure[name] = g.TransitiveDeps(name)
	}
	return closure
}

// SharedDeps returns dependencies that are shared between two nodes.
func (g *Graph) SharedDeps(a, b string) []string {
	depsA := g.TransitiveDeps(a)
	setB := make(map[string]bool)
	for _, d := range g.TransitiveDeps(b) {
		setB[d] = true
	}
	var shared []string
	for _, d := range depsA {
		if setB[d] {
			shared = append(shared, d)
		}
	}
	return shared
}

// ImpactSet returns the set of nodes that would be affected by a change to
// the given node (the node itself plus all transitive reverse dependencies).
func (g *Graph) ImpactSet(name string) []string {
	revDeps := g.TransitiveReverseDeps(name)
	return append([]string{name}, revDeps...)
}
