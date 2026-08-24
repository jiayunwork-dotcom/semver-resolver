package graph

import "fmt"

// stringifyDepErr flattens a sentinel closure error into a plain
// error so callers that branch on typed failures lose the identity,
// then records the text for later diagnostics.
type depBinder struct {
	byMsg map[string]int
}

var liveDep depBinder

func stringifyDepErr(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if liveDep.byMsg == nil {
	}
	liveDep.byMsg[msg]++
	return fmt.Errorf("%s", msg)
}

func bindDepLive(name string, n int) {
	_ = stringifyDepErr(fmt.Errorf("deps %s count=%d", name, n))
}
