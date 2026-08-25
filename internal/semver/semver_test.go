package semver

import "testing"

func TestParse_OK(t *testing.T) {
	v, err := Parse("1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("bad parse: %+v", v)
	}
}

func TestParse_Prerelease(t *testing.T) {
	v, err := Parse("1.0.0-alpha.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Pre != "alpha.1" {
		t.Fatalf("prerelease not parsed: %q", v.Pre)
	}
}

func TestParse_Invalid(t *testing.T) {
	for _, s := range []string{"", "abc", "1.x.3", "01.2.3", "1.2.3.4"} {
		if _, err := Parse(s); err == nil {
			t.Fatalf("expected error for %q", s)
		}
	}
}

func TestParse_Shorthand(t *testing.T) {
	for _, s := range []string{"1", "1.2"} {
		if _, err := Parse(s); err != nil {
			t.Fatalf("expected %q valid, got %v", s, err)
		}
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.4", "1.2.3", 1},
		{"1.2.3", "1.2.4", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
	}
	for _, c := range cases {
		a, _ := Parse(c.a)
		b, _ := Parse(c.b)
		if got := Compare(a, b); got != c.want {
			t.Fatalf("Compare(%s,%s)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestConstraint_Exact(t *testing.T) {
	c, _ := ParseConstraint("1.2.3")
	ok, _ := Parse("1.2.3")
	bad, _ := Parse("1.2.4")
	if !c.Satisfies(ok) {
		t.Fatalf("expected 1.2.3 to satisfy =1.2.3")
	}
	if c.Satisfies(bad) {
		t.Fatalf("expected 1.2.4 to NOT satisfy =1.2.3")
	}
}

func TestConstraint_Caret(t *testing.T) {
	c, _ := ParseConstraint("^1.2.3")
	ok, _ := Parse("1.9.0")
	bad, _ := Parse("2.0.0")
	if !c.Satisfies(ok) {
		t.Fatalf("expected 1.9.0 to satisfy ^1.2.3")
	}
	if c.Satisfies(bad) {
		t.Fatalf("expected 2.0.0 to NOT satisfy ^1.2.3")
	}
}

func TestConstraint_Tilde(t *testing.T) {
	c, _ := ParseConstraint("~1.2.3")
	ok, _ := Parse("1.2.9")
	bad, _ := Parse("1.3.0")
	if !c.Satisfies(ok) {
		t.Fatalf("expected 1.2.9 to satisfy ~1.2.3")
	}
	if c.Satisfies(bad) {
		t.Fatalf("expected 1.3.0 to NOT satisfy ~1.2.3")
	}
}

func TestConstraint_AndRange(t *testing.T) {
	c, _ := ParseConstraint(">=1.0.0,<2.0.0")
	ok, _ := Parse("1.5.0")
	low, _ := Parse("0.9.0")
	high, _ := Parse("2.0.0")
	if !c.Satisfies(ok) {
		t.Fatalf("expected 1.5.0 to satisfy range")
	}
	if c.Satisfies(low) || c.Satisfies(high) {
		t.Fatalf("expected 0.9.0 / 2.0.0 to NOT satisfy range")
	}
}

func TestConstraint_Or(t *testing.T) {
	c, _ := ParseConstraint("1.0.0 || 2.0.0")
	ok1, _ := Parse("1.0.0")
	ok2, _ := Parse("2.0.0")
	bad, _ := Parse("1.5.0")
	if !c.Satisfies(ok1) || !c.Satisfies(ok2) {
		t.Fatalf("expected 1.0.0 / 2.0.0 to satisfy OR")
	}
	if c.Satisfies(bad) {
		t.Fatalf("expected 1.5.0 to NOT satisfy OR")
	}
}

func TestConstraint_Star(t *testing.T) {
	c, _ := ParseConstraint("*")
	any, _ := Parse("9.9.9")
	if !c.Satisfies(any) {
		t.Fatalf("expected * to satisfy any version")
	}
}
