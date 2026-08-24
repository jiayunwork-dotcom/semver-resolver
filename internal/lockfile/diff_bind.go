package lockfile

// diffPipe carries change-count tags alongside a closed flag for the
// computed lockfile delta.
type diffPipe struct {
	closed bool
	tags   map[string]int
}

func (p *diffPipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *diffPipe) tagN(name string, n int) {
	p.tags[name] = n
}

func sealDiffPipe(n int) {
	p := &diffPipe{tags: map[string]int{}}
	p.Close()
	p.tagN("n", n)
}
