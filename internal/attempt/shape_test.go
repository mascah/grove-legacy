package attempt

import (
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// G-260927-dx0yn: each fact the shape derives, and each limit docs/commands.md
// states, from a recorded stream (testdata/shape-events.jsonl): the
// provider ran in /w/p, whose record root is /w/p/grove.
func TestReadShape(t *testing.T) {
	t.Parallel()
	s, err := ReadShape(filepath.Join("testdata", "shape-events.jsonl"), "/w/p", "/w/p/grove")
	if err != nil {
		t.Fatal(err)
	}
	// Process: go run ./cmd/grove guide, cd && grove context, the Read and
	// the relative Write and MultiEdit under the root, git worktree, a
	// checkout -b of a worktree- branch, grove by path, grove guide without
	// a name, and a heredoc line starting with grove (a known false
	// positive). Not process: a Read outside the root, git checkout main,
	// sed -i, env grove, $(grove …), FOO=1 grove, and an Edit of the root
	// spelled through /private.
	if s.Tools != 19 || s.Process != 10 || s.Skipped != 0 {
		t.Fatalf("%+v", s)
	}
	if !maps.Equal(s.Guides, map[string]int{"work": 1, "model": 1}) {
		t.Fatalf("guides %v", s.Guides)
	}
	// The first edit is a subagent's Edit, the 12th call: sed -i earlier is
	// not an edit, the Write under the root is process.
	want := FirstEdit{Tool: 12, At: time.Date(2026, 9, 22, 18, 34, 10, 0, time.UTC), Path: "/w/p/internal/y.go"}
	if s.FirstEdit == nil || *s.FirstEdit != want {
		t.Fatalf("first edit %+v", s.FirstEdit)
	}
	v := &View{Status: Finished, Launch: Launch{Worktree: "/w", Started: now}, Shape: &s}
	if got := shapeText(v); got != "19 tool calls, 10 process (53%); guides printed: model 1, work 1; first edit outside the record root: tool 12, 4m10s after the start, p/internal/y.go" {
		t.Fatal(got)
	}

	// A line over MaxLine is skipped and counted, so the counts are lower
	// bounds; a running attempt's are so far; no edit and no timestamp are
	// said, never filled.
	data, err := os.ReadFile(filepath.Join("testdata", "shape-events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "events.jsonl")
	huge := `{"type":"assistant","x":"` + strings.Repeat("x", MaxLine) + "\"}\n"
	if err := os.WriteFile(path, append(data, huge...), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err = ReadShape(path, "/w/p", "/w/p/grove")
	if err != nil || s.Skipped != 1 || s.Tools != 19 {
		t.Fatalf("%+v %v", s, err)
	}
	s.FirstEdit.At = time.Time{}
	v = &View{Status: Running, Launch: Launch{Worktree: "/elsewhere", Started: now}, Shape: &s}
	if got := shapeText(v); got != "so far, ≥19 tool calls, ≥10 process (53%); guides printed: model 1, work 1; first edit outside the record root: tool 12, at an unknown time, /w/p/internal/y.go; 1 oversized lines not read: the counts are lower bounds, and the first edit is the first in the lines read" {
		t.Fatal(got)
	}
	v.Shape = &Shape{}
	if got := shapeText(v); got != "so far, 0 tool calls, 0 process (0%); guides printed: none; no edit outside the record root" {
		t.Fatal(got)
	}
	if s, err := ReadShape(filepath.Join(t.TempDir(), "missing.jsonl"), "/w/p", "/w/p/grove"); err != nil || s.Tools != 0 {
		t.Fatalf("%+v %v", s, err)
	}
}

// The record root comes from grove.yaml at the base, and the files changed
// base to HEAD split at it; a base Git cannot read is unknown, not empty.
func TestChangedAndShapeFacts(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	write(t, root, "grove/G-260101-00003-plan.md", "plan")
	write(t, root, "internal/a b.go", "package a")
	write(t, root, "grovey.txt", "not under grove/")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "work")
	c := changedFiles(root, "", base, git(t, root, "rev-parse", "HEAD"))
	if c.Error != "" || c.RecordRoot != "grove" || strings.Join(c.Records, ",") != "grove/G-260101-00003-plan.md" || strings.Join(c.Other, ",") != "grovey.txt,internal/a b.go" {
		t.Fatalf("%+v", c)
	}
	unknown := changedFiles(root, "", strings.Repeat("0", 40), base)
	if unknown.Error == "" || len(unknown.Other) != 0 {
		t.Fatalf("%+v", unknown)
	}
	visible := func(s string) string { return s }
	facts := func(r *Result, v View) string {
		v.Status, v.Result = Finished, r
		return strings.Join(Facts(&v, visible), "\n")
	}
	if f := facts(&Result{Changed: unknown}, View{}); !strings.Contains(f, "\nChanged: unknown: ") || strings.Contains(f, "Shape:") {
		t.Fatal(f)
	}
	if f := facts(&Result{}, View{ShapeError: "no grove.yaml"}); !strings.Contains(f, "\nChanged: unknown: not recorded when it finished\n") || !strings.Contains(f, "\nShape: unknown: no grove.yaml\n") {
		t.Fatal(f)
	}
	many := &Changed{RecordRoot: "grove", Records: []string{"grove/a.md"}, Other: strings.Split("a b c d e f g h i j k l", " ")}
	if f := facts(&Result{Changed: many}, View{}); !strings.Contains(f, "\nChanged: 12 outside grove/ (a, b, c, d, e, f, g, h, i, j and 2 more), 1 under it (grove/a.md)\n") {
		t.Fatal(f)
	}
	if f := facts(&Result{Changed: &Changed{RecordRoot: "p/grove"}}, View{}); !strings.Contains(f, "\nChanged: 0 outside p/grove/, 0 under it\n") {
		t.Fatal(f)
	}
}

func TestSum(t *testing.T) {
	t.Parallel()
	done := func(minutes, turns, lastTurns int, usd float64) View {
		r := &Result{Finished: now.Add(time.Duration(minutes) * time.Minute), Events: Events{Turns: turns}}
		if lastTurns >= 0 {
			r.Events.Result = &Final{CostUSD: usd, Turns: lastTurns}
		}
		return View{Launch: Launch{Started: now}, Result: r}
	}
	if got := Sum(nil).String(); got != "Total: 0 attempts, $0.00, 0 turns, 0m" {
		t.Fatal(got)
	}
	// Turns are summed over a run's result events; an attempt that recorded
	// only the last one's makes the sum a lower bound.
	views := []View{done(30, 109, 14, 9.75), done(10, 0, 32, 1.62), done(5, 0, -1, 0), {Launch: Launch{Started: now}}}
	if got := Sum(views).String(); got != "Total: 4 attempts, $11.37, ≥141 turns, 45m; 1 unfinished in none of these, 1 without a result event not in the cost or turns" {
		t.Fatal(got)
	}
	if got := Sum(views[:1]).String(); got != "Total: 1 attempt, $9.75, 109 turns, 30m" {
		t.Fatal(got)
	}
}
