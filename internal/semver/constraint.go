package semver

import (
	"fmt"
	"strings"
	"unicode"
)

type cmp struct {
	op string // = > >= < <=
	v  Version
}

// Constraint 是一组「或」关系（alts）的约束，每个 alt 内部是「且」关系（ands）。
type Constraint struct {
	alts [][]cmp
}

// ParseConstraint 解析版本约束字符串。
func ParseConstraint(s string) (Constraint, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "*" {
		return Constraint{alts: [][]cmp{{}}}, nil
	}
	var c Constraint
	for _, altStr := range strings.Split(s, "||") {
		altStr = strings.TrimSpace(altStr)
		if altStr == "" {
			return Constraint{}, fmt.Errorf("empty alternative in constraint")
		}
		toks := splitTokens(altStr)
		var ands []cmp
		for _, tok := range toks {
			cs, err := expandToken(tok)
			if err != nil {
				return Constraint{}, err
			}
			ands = append(ands, cs...)
		}
		c.alts = append(c.alts, ands)
	}
	return c, nil
}

func splitTokens(s string) []string {
	f := func(r rune) bool { return r == ',' || unicode.IsSpace(r) }
	return strings.FieldsFunc(s, f)
}

func expandToken(tok string) ([]cmp, error) {
	if tok == "" || tok == "*" {
		return nil, nil
	}
	if strings.HasPrefix(tok, "^") {
		v, err := Parse(tok[1:])
		if err != nil {
			return nil, err
		}
		return caret(v), nil
	}
	if strings.HasPrefix(tok, "~") {
		v, err := Parse(tok[1:])
		if err != nil {
			return nil, err
		}
		return tilde(v), nil
	}
	if op, val, ok := splitOp(tok); ok {
		v, err := Parse(val)
		if err != nil {
			return nil, err
		}
		return []cmp{{op: op, v: v}}, nil
	}
	v, err := Parse(tok)
	if err != nil {
		return nil, err
	}
	return []cmp{{op: "=", v: v}}, nil
}

func splitOp(tok string) (string, string, bool) {
	for _, op := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(tok, op) {
			return op, tok[len(op):], true
		}
	}
	return "", "", false
}

func caret(v Version) []cmp {
	var upper Version
	switch {
	case v.Major > 0:
		upper = Version{Major: v.Major + 1}
	case v.Minor > 0:
		upper = Version{Minor: v.Minor + 1}
	default:
		upper = Version{Patch: v.Patch + 1}
	}
	return []cmp{{op: ">=", v: v}, {op: "<", v: upper}}
}

func tilde(v Version) []cmp {
	var upper Version
	if v.Minor == 0 && v.Patch == 0 {
		upper = Version{Major: v.Major + 1}
	} else {
		upper = Version{Major: v.Major, Minor: v.Minor + 1}
	}
	return []cmp{{op: ">=", v: v}, {op: "<", v: upper}}
}

func satisfyCmp(op string, cv, v Version) bool {
	switch op {
	case "=":
		return Compare(v, cv) == 0
	case ">":
		return Compare(v, cv) > 0
	case ">=":
		return Compare(v, cv) >= 0
	case "<":
		return Compare(v, cv) < 0
	case "<=":
		return Compare(v, cv) <= 0
	}
	return false
}

// Satisfies 判断版本 v 是否满足约束。
func (c Constraint) Satisfies(v Version) bool {
	for _, alt := range c.alts {
		ok := true
		for _, cm := range alt {
			if !satisfyCmp(cm.op, cm.v, v) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
