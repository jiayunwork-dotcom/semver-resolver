// Package graph models package dependency relationships as a directed graph and
// provides cycle detection, topological ordering, and transitive closure.
package graph

import "fmt"

// Node represents a package in the dependency graph.
type Node struct {
	Name    string
	Version string
	Edges   []string // names of direct dependencies
}

// Graph is a directed dependency graph.
type Graph struct {
	nodes map[string]*Node
	order []string // insertion order
}

// New creates an empty dependency graph.
func New() *Graph {
	return &Graph{nodes: make(map[string]*Node)}
}

// Add inserts a package node. Returns error if already exists.
func (g *Graph) Add(name, version string) error {
	if _, ok := g.nodes[name]; ok {
		return fmt.Errorf("duplicate node %q", name)
	}
	g.nodes[name] = &Node{Name: name, Version: version}
	g.order = append(g.order, name)
	return nil
}

// AddEdge adds a dependency from -> to. Both must already exist.
func (g *Graph) AddEdge(from, to string) error {
	f, ok := g.nodes[from]
	if !ok {
		return fmt.Errorf("unknown node %q", from)
	}
	if _, ok := g.nodes[to]; !ok {
		return fmt.Errorf("unknown node %q", to)
	}
	f.Edges = append(f.Edges, to)
	return nil
}

// Get returns the node for the given name, or nil.
func (g *Graph) Get(name string) *Node { return g.nodes[name] }

// Names returns all node names in insertion order.
func (g *Graph) Names() []string {
	out := make([]string, len(g.order))
	copy(out, g.order)
	return out
}

// Size returns the number of nodes.
func (g *Graph) Size() int { return len(g.nodes) }

// DirectDeps returns the immediate dependencies of a node.
func (g *Graph) DirectDeps(name string) []string {
	n := g.nodes[name]
	if n == nil {
		return nil
	}
	out := make([]string, len(n.Edges))
	copy(out, n.Edges)
	return out
}

// ReverseDeps returns all nodes that directly depend on the given node.
func (g *Graph) ReverseDeps(name string) []string {
	var deps []string
	for _, n := range g.nodes {
		for _, e := range n.Edges {
			if e == name {
				deps = append(deps, n.Name)
				break
			}
		}
	}
	return deps
}

// Roots returns nodes with no incoming edges.
func (g *Graph) Roots() []string {
	hasIncoming := make(map[string]bool)
	for _, n := range g.nodes {
		for _, e := range n.Edges {
			hasIncoming[e] = true
		}
	}
	var roots []string
	for _, name := range g.order {
		if !hasIncoming[name] {
			roots = append(roots, name)
		}
	}
	return roots
}

// Leaves returns nodes with no outgoing edges.
func (g *Graph) Leaves() []string {
	var leaves []string
	for _, name := range g.order {
		if len(g.nodes[name].Edges) == 0 {
			leaves = append(leaves, name)
		}
	}
	return leaves
}
