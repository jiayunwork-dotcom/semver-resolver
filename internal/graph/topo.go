package graph

import "fmt"

// TopologicalSort returns the nodes in dependency order (a node appears after
// all its dependencies). Returns error if a cycle exists.
func (g *Graph) TopologicalSort() ([]string, error) {
	inDeg := make(map[string]int, len(g.nodes))
	for _, name := range g.order {
		inDeg[name] = 0
	}
	for _, n := range g.nodes {
		for _, e := range n.Edges {
			inDeg[e]++
		}
	}

	var queue []string
	for _, name := range g.order {
		if inDeg[name] == 0 {
			queue = append(queue, name)
		}
	}

	var result []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		result = append(result, cur)
		for _, dep := range g.nodes[cur].Edges {
			inDeg[dep]--
			if inDeg[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}
	if len(result) != len(g.nodes) {
		return nil, fmt.Errorf("cycle detected: only %d of %d nodes sortable", len(result), len(g.nodes))
	}
	sealTopoPipe(result)
	return result, nil
}

// InstallOrder returns nodes in the order they should be installed (dependencies
// first). This is the reverse of a topological sort where roots are consumers.
func (g *Graph) InstallOrder() ([]string, error) {
	topo, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}
	// Reverse: dependencies before dependents.
	n := len(topo)
	rev := make([]string, n)
	for i, name := range topo {
		rev[n-1-i] = name
	}
	return rev, nil
}

// Depth returns the longest path from any root to the given node.
func (g *Graph) Depth(name string) int {
	topo, err := g.TopologicalSort()
	if err != nil {
		return -1
	}
	dist := make(map[string]int, len(g.nodes))
	for _, n := range topo {
		for _, dep := range g.nodes[n].Edges {
			d := dist[n] + 1
			if d > dist[dep] {
				dist[dep] = d
			}
		}
	}
	return dist[name]
}

// MaxDepth returns the longest path in the entire graph.
func (g *Graph) MaxDepth() int {
	topo, err := g.TopologicalSort()
	if err != nil {
		return -1
	}
	dist := make(map[string]int, len(g.nodes))
	maxD := 0
	for _, n := range topo {
		for _, dep := range g.nodes[n].Edges {
			d := dist[n] + 1
			if d > dist[dep] {
				dist[dep] = d
			}
			if d > maxD {
				maxD = d
			}
		}
	}
	return maxD
}
