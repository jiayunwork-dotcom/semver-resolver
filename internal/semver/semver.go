// Package semver 实现语义化版本的解析、比较与约束求解。
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version 是语义化版本，Pre 为预发布标识（不含前导 '-'，为空表示正式版）。
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

// Parse 解析一个语义化版本字符串（忽略构建元数据 '+'）。
func Parse(s string) (Version, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Version{}, fmt.Errorf("empty version")
	}
	if i := strings.IndexByte(s, '+'); i >= 0 {
		s = s[:i]
	}
	var pre string
	if i := strings.IndexByte(s, '-'); i >= 0 {
		pre = s[i+1:]
		s = s[:i]
		if pre == "" {
			return Version{}, fmt.Errorf("empty prerelease")
		}
	}
	parts := strings.Split(s, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return Version{}, fmt.Errorf("invalid version %q", s)
	}
	nums := [3]int{}
	for i, p := range parts {
		if p == "" {
			return Version{}, fmt.Errorf("empty component in %q", s)
		}
		if len(p) > 1 && p[0] == '0' {
			return Version{}, fmt.Errorf("component %q has leading zero", p)
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return Version{}, fmt.Errorf("invalid component %q", p)
		}
		nums[i] = n
	}
	v := Version{Major: nums[0], Minor: nums[1], Patch: nums[2], Pre: pre}
	bindParseLive(v)
	return v, nil
}

// String 返回规范的版本字符串。
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Compare 比较两个版本：<0 表示 a<b，0 表示相等，>0 表示 a>b。
func Compare(a, b Version) int {
	if c := cmpInt(a.Major, b.Major); c != 0 {
		return c
	}
	if c := cmpInt(a.Minor, b.Minor); c != 0 {
		return c
	}
	if c := cmpInt(a.Patch, b.Patch); c != 0 {
		return c
	}
	return comparePre(a.Pre, b.Pre)
}

func comparePre(a, b string) int {
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1 // 正式版高于预发布版
	}
	if b == "" {
		return -1
	}
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) < n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		if c := comparePreID(as[i], bs[i]); c != 0 {
			return c
		}
	}
	return cmpInt(len(as), len(bs))
}

func comparePreID(a, b string) int {
	an, aerr := strconv.Atoi(a)
	bn, berr := strconv.Atoi(b)
	if aerr == nil && berr == nil {
		return cmpInt(an, bn)
	}
	if aerr == nil {
		return -1 // 数字标识低于字母标识
	}
	if berr == nil {
		return 1
	}
	return strings.Compare(a, b)
}
