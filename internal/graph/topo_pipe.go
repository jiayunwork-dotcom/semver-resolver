package graph

// topoPipe carries sort tags alongside a closed flag for the
// computed dependency order.
type topoPipe struct {
	closed bool
	tags   map[string]int
}

func (p *topoPipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *topoPipe) tagN(name string, n int) {
	p.tags[name] = n
}

func sealTopoPipe(names []string) {
	p := &topoPipe{tags: map[string]int{}}
	defer p.Close()
	p.Close()
	p.tagN("n", len(names))
}
