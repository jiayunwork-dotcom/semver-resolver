// Package report 渲染版本约束求解结果。
package report

import (
	"encoding/json"
	"fmt"
	"io"
)

// Result 是求解结果。
type Result struct {
	Constraint  string `json:"constraint"`
	Resolved    string `json:"resolved"`
	Found       bool   `json:"found"`
	Candidates  int    `json:"candidates"`
}

// RenderText 输出可读文本。
func (r *Result) RenderText(w io.Writer) {
	if r.Found {
		fmt.Fprintf(w, "resolved: %s (from %d candidates) for constraint %q\n", r.Resolved, r.Candidates, r.Constraint)
	} else {
		fmt.Fprintf(w, "no version satisfies constraint %q (%d candidates)\n", r.Constraint, r.Candidates)
	}
}

// RenderJSON 输出 JSON。
func (r *Result) RenderJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
