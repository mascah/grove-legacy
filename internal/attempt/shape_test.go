package attempt

import (
	"context"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mascah/grove/internal/repo"
)

// G-260927-dx0yn: each fact the shape derives, and each limit docs/commands.md
// states, from a recorded stream (testdata/shape-events.jsonl).
func TestReadShape(t *testing.T) {
	t.Parallel()
	s, err := ReadShape(filepath.Join("testdata", "shape-events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	// Running grove: go run ./cmd/grove guide, cd && grove context, grove by
	// path, grove guide without a name, and a heredoc line starting with
	// grove (a known false positive). Not seen: env grove, $(grove …) and
	// FOO=1 grove. The rest are git, edits, reads and sed.
	if s.Tools != 19 || s.Grove != 5 || s.Skipped != 0 {
		t.Fatalf("%+v", s)
	}
	if !maps.Equal(s.Guides, map[string]int{"work": 1, "model": 1}) {
		t.Fatalf("guides %v", s.Guides)
	}
	v := &View{Status: Finished, Shape: &s}
	if got := shapeText(v); got != "19 tool calls, 5 running grove; guides printed: model 1, work 1" {
		t.Fatal(got)
	}

	// A line over MaxLine is skipped and counted, so the counts are lower
	// bounds; a running attempt's are so far.
	data, err := os.ReadFile(filepath.Join("testdata", "shape-events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "events.jsonl")
	huge := `{"type":"assistant","x":"` + strings.Repeat("x", MaxLine) + "\"}\n"
	if err := os.WriteFile(path, append(data, huge...), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err = ReadShape(path)
	if err != nil || s.Skipped != 1 || s.Tools != 19 {
		t.Fatalf("%+v %v", s, err)
	}
	v = &View{Status: Running, Shape: &s}
	if got := shapeText(v); got != "so far, ≥19 tool calls, ≥5 running grove; guides printed: model 1, work 1; 1 oversized lines not read: the counts are lower bounds" {
		t.Fatal(got)
	}
	v.Shape = &Shape{}
	if got := shapeText(v); got != "so far, 0 tool calls, 0 running grove; guides printed: none" {
		t.Fatal(got)
	}
	if s, err := ReadShape(filepath.Join(t.TempDir(), "missing.jsonl")); err != nil || s.Tools != 0 {
		t.Fatalf("%+v %v", s, err)
	}
}

// The record root comes from grove.yaml at the base, and the files changed
// base to HEAD split at it; a base Git cannot read is unknown, not empty.
func TestChangedAndShapeFacts(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	// A plan under the record root, then two commits outside it: the first
	// of those is the first commit outside it.
	gitAt := func(at time.Time, args ...string) {
		cmd := repo.Command(context.Background(), root, append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false"}, args...)...)
		cmd.Env = append(cmd.Env, "GIT_COMMITTER_DATE="+at.Format(time.RFC3339))
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
	}
	commit := func(at time.Time, files ...string) {
		for _, f := range files {
			write(t, root, f, f)
		}
		git(t, root, "add", "-A")
		gitAt(at, "commit", "-qm", "work")
	}
	commit(now.Add(time.Minute), "grove/G-260101-00003-plan.md")
	commit(now.Add(5*time.Minute), "internal/a b.go")
	commit(now.Add(9*time.Minute), "grovey.txt")
	c := changedFiles(root, "", base, git(t, root, "rev-parse", "HEAD"))
	if c.Error != "" || c.RecordRoot != "grove" || strings.Join(c.Records, ",") != "grove/G-260101-00003-plan.md" || strings.Join(c.Other, ",") != "grovey.txt,internal/a b.go" || !c.FirstOther.Equal(now.Add(5*time.Minute)) {
		t.Fatalf("%+v", c)
	}
	if f := strings.Join(Facts(&View{Status: Finished, Launch: Launch{Started: now}, Result: &Result{Changed: c}}, func(s string) string { return s }), "\n"); !strings.Contains(f, "\nChanged: 2 outside grove/ (grovey.txt, internal/a b.go), 1 under it (grove/G-260101-00003-plan.md); first commit outside it 5m0s after the start\n") {
		t.Fatal(f)
	}
	// A resolution attempt merges the target: the target's commits, older
	// than the attempt, are not its first; the merge is. A time before the
	// start is said, not clamped.
	tip := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "target", base)
	commit(now.Add(-time.Hour), "target.go")
	git(t, root, "checkout", "-q", "-")
	gitAt(now.Add(20*time.Minute), "merge", "-q", "--no-ff", "-m", "merge", "target")
	merged := changedFiles(root, "", tip, git(t, root, "rev-parse", "HEAD"))
	if strings.Join(merged.Other, ",") != "target.go" || !merged.FirstOther.Equal(now.Add(20*time.Minute)) {
		t.Fatalf("%+v", merged)
	}
	if f := strings.Join(Facts(&View{Status: Finished, Launch: Launch{Started: now.Add(30 * time.Minute)}, Result: &Result{Changed: merged}}, func(s string) string { return s }), "\n"); !strings.Contains(f, "; first commit outside it 10m0s before the start\n") {
		t.Fatal(f)
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
	if f := facts(&Result{}, View{ShapeError: "open events.jsonl: permission denied"}); !strings.Contains(f, "\nChanged: unknown: not recorded when it finished\n") || !strings.Contains(f, "\nShape: unknown: open events.jsonl: permission denied\n") {
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
