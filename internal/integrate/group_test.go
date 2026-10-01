package integrate

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/update"
)

// groupFixture is a main checkout and a worktree on branch feature where
// G-260101-00001 and G-260101-00003, selected together, share one candidate: the code commit,
// then one commit putting both in review with it (G-260925-wc2pz).
func groupFixture(t *testing.T) (root, wt, candidate string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt = filepath.Join(dir, "repo"), filepath.Join(dir, "feat")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	second := strings.NewReplacer(`"G-260101-00001"`, `"G-260101-00003"`, "title: First", "title: Second\ndepends_on: [\"G-260101-00001\"]").Replace(work)
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "proposed", 1))
	write(t, root, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "proposed", 1))
	write(t, root, "code.txt", "before\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wt, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, wt, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "active", 1))
	write(t, wt, "code.txt", "both changes\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: implement both")
	candidate = git(t, wt, "rev-parse", "HEAD")
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if _, err := update.Apply(wt, update.Request{ID: id, Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: candidate}}}, now, nil); err != nil {
			t.Fatal(err)
		}
	}
	git(t, wt, "commit", "-qam", "docs: hand both off")
	return root, wt, candidate
}

func records(t *testing.T, root string) map[string]*project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("%v", ds)
	}
	out := map[string]*project.Record{}
	for _, r := range p.Records {
		out[r.ID] = r
	}
	return out
}

func TestGroupIntegratesOnlyWhenEveryMemberIsApproved(t *testing.T) {
	t.Parallel()
	root, wt, candidate := groupFixture(t)
	// Approval stays per member, and the other member's record commits are
	// not a new candidate.
	if _, err := update.Approve(wt, "G-260101-00001", "First is right.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	main := git(t, root, "rev-parse", "HEAD")
	var facts []string
	err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(f string) { facts = append(facts, f) })
	if err == nil || !strings.Contains(err.Error(), "shared by G-260101-00001, G-260101-00003, and delivering it delivers all of them, but G-260101-00003 is awaiting judgment; judge each first (grove approve ID VERDICT in "+wt+")") || facts != nil {
		t.Fatalf("%v %q", err, facts)
	}
	if git(t, root, "rev-parse", "HEAD") != main {
		t.Fatal("a refused group integration moved main")
	}
	if _, err := update.Approve(wt, "G-260101-00003", "Second too.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	// Integrating either member delivers the group: one squash commit
	// naming both, each proved and done, no record commit.
	if err := Run(Request{Root: root, ID: "G-260101-00003", Keep: true}, now, func(f string) { facts = append(facts, f) }); err != nil {
		t.Fatalf("%v %q", err, facts)
	}
	if len(facts) != 7 || !strings.Contains(facts[0], "of G-260101-00001 accepted") || !strings.Contains(facts[1], "of G-260101-00003 accepted") || !strings.HasPrefix(facts[3], "delivery: squash commit ") ||
		!strings.HasPrefix(facts[4], "done: G-260101-00001 is done: delivered to main, proved: squashed as ") || !strings.HasPrefix(facts[5], "done: G-260101-00003 is done: delivered to main, proved: squashed as ") {
		t.Fatalf("%q", facts)
	}
	for id, r := range records(t, root) {
		if r.Status != "accepted" || r.Candidate != candidate {
			t.Fatalf("%s %+v", id, r)
		}
	}
	if n := git(t, root, "rev-list", "--count", main+"..HEAD"); n != "1" {
		t.Fatalf("%s commits delivered the group", n)
	}
	if body := git(t, root, "log", "-1", "--format=%B"); !strings.Contains(body, "Grove-Work: G-260101-00001\nGrove-Work: G-260101-00003\n") || !strings.HasPrefix(body, "feat: first\n") {
		t.Fatalf("message:\n%s", body)
	}
}

func TestGroupFeedbackReopensEveryMember(t *testing.T) {
	t.Parallel()
	root, wt, candidate := groupFixture(t)
	if _, err := update.Approve(wt, "G-260101-00001", "First is right.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	before := git(t, wt, "rev-parse", "HEAD")
	res, err := update.Feedback(wt, "G-260101-00003", "Second misses a case.", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Reopened) != 1 || res.Reopened[0].ID != "G-260101-00001" || res.Reopened[0].Commit == "" {
		t.Fatalf("%+v", res)
	}
	rs := records(t, wt)
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if r := rs[id]; r.Status != "active" || r.Approved != "" || r.Candidate != candidate {
			t.Fatalf("%s %+v", id, r)
		}
	}
	if !strings.Contains(string(rs["G-260101-00003"].Source), "Feedback on candidate "+candidate[:7]+", 2026-09-22: Second misses a case.") ||
		!strings.Contains(string(rs["G-260101-00001"].Source), "Reopened with G-260101-00003's feedback on candidate "+candidate[:7]+", 2026-09-22") {
		t.Fatal("the feedback lines are missing")
	}
	// One commit per record, each touching only its own file.
	commits := strings.Fields(git(t, wt, "rev-list", before+"..HEAD"))
	if len(commits) != 2 {
		t.Fatalf("%v", commits)
	}
	for _, c := range commits {
		if files := strings.Fields(git(t, wt, "diff-tree", "--no-commit-id", "--name-only", "-r", c)); len(files) != 1 || !slices.Contains([]string{"grove/G-260101-00001-first.md", "grove/G-260101-00003-second.md"}, files[0]) {
			t.Fatalf("%s touches %v", c, files)
		}
	}
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(string) {}); err == nil || !strings.Contains(err.Error(), "no branch holds G-260101-00001 accepted") {
		t.Fatal(err)
	}
}

// A member reopened with the group's feedback and continued alone: the
// next candidate carries its sibling's code, which lost its approval, so
// integration refuses with the target unchanged.
func TestGroupRefusesToCarryAReopenedSibling(t *testing.T) {
	t.Parallel()
	root, wt, _ := groupFixture(t)
	if _, err := update.Feedback(wt, "G-260101-00001", "One more case.", now); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "code.txt", "both changes, and the case\n")
	git(t, wt, "commit", "-qam", "feat: the case")
	next := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: next}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(wt, "G-260101-00001", "Now right.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	main := git(t, root, "rev-parse", "HEAD")
	err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(string) {})
	if err == nil || !strings.Contains(err.Error(), "delivering feature would also carry G-260101-00003's candidate") || !strings.Contains(err.Error(), "which is active without an acceptance") {
		t.Fatal(err)
	}
	if git(t, root, "rev-parse", "HEAD") != main {
		t.Fatal("main moved")
	}
}
