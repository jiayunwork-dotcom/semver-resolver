package main

import (
	"fmt"
	"io"
	"os"

	"semver-resolver/internal/report"
	"semver-resolver/internal/resolve"
	"semver-resolver/internal/semver"
	"semver-resolver/internal/server"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  semver-resolver resolve --constraint C --versions v1,v2,... [--format text|json]")
	fmt.Fprintln(w, "  semver-resolver check --constraint C --version V [--format text|json]")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runServer(nil)
		return
	}
	if args[0] == "serve" {
		runServer(args[1:])
		return
	}
	if err := run(args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runServer(args []string) {
	addr := ":8080"
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-addr" || args[i] == "--addr" {
			addr = args[i+1]
			break
		}
	}
	cfg := server.Config{Addr: addr}
	fmt.Fprintf(os.Stdout, "semver-resolver server listening on %s\n", server.FormatAddr(addr))
	if err := server.ListenAndServe(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// run 是可测试的入口。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		usage(stderr)
		return fmt.Errorf("missing command")
	}
	switch args[0] {
	case "resolve":
		return cmdResolve(args[1:], stdout, stderr)
	case "check":
		return cmdCheck(args[1:], stdout, stderr)
	default:
		usage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func cmdResolve(args []string, stdout, stderr io.Writer) error {
	var constraint, versions, format string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--constraint":
			if i+1 < len(args) {
				constraint = args[i+1]
				i++
			}
		case "--versions":
			if i+1 < len(args) {
				versions = args[i+1]
				i++
			}
		case "--format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		}
	}
	if constraint == "" || versions == "" {
		usage(stderr)
		return fmt.Errorf("resolve requires --constraint and --versions")
	}
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		return err
	}
	cands, err := resolve.ParseCandidates(versions)
	if err != nil {
		return err
	}
	best, found := resolve.Highest(c, cands)
	res := &report.Result{Constraint: constraint, Candidates: len(cands), Found: found}
	if found {
		res.Resolved = best.String()
	}
	return emit(res, format, stdout)
}

func cmdCheck(args []string, stdout, stderr io.Writer) error {
	var constraint, version, format string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--constraint":
			if i+1 < len(args) {
				constraint = args[i+1]
				i++
			}
		case "--version":
			if i+1 < len(args) {
				version = args[i+1]
				i++
			}
		case "--format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		}
	}
	if constraint == "" || version == "" {
		usage(stderr)
		return fmt.Errorf("check requires --constraint and --version")
	}
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		return err
	}
	v, err := semver.Parse(version)
	if err != nil {
		return err
	}
	res := &report.Result{Constraint: constraint, Candidates: 1, Found: c.Satisfies(v)}
	if res.Found {
		res.Resolved = v.String()
	}
	return emit(res, format, stdout)
}

func emit(r *report.Result, format string, stdout io.Writer) error {
	if format == "json" {
		return r.RenderJSON(stdout)
	}
	r.RenderText(stdout)
	return nil
}
