package attempt

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestWorkspaceStartsFromTheTarget: a new workspace starts from the target,
// not the launching checkout's HEAD (G-260930-tcc9w), and takes from HEAD
// only the selected records the target lacks with the records they need:
// what they relate to and the plan for them, never code or unrelated
// records. A selected record the target holds otherwise is refused.
func TestWorkspaceStartsFromTheTarget(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	main := git(t, root, "rev-parse", "main")
	git(t, root, "checkout", "-q", "-b", "shape")
	member(t, root, "G-260101-00003", "proposed")
	write(t, root, "grove/G-260101-00003.md", strings.Replace(readFile(t, root, "grove/G-260101-00003.md"), "---\n\n", "relates_to: [\"G-260101-00001\", \"G-260101-00004\"]\n---\n\n", 1))
	write(t, root, "grove/G-260101-00004-page.md", "---\nid: \"G-260101-00004\"\ntype: page\ntitle: Design\n---\n\nIt.\n")
	write(t, root, "grove/G-260101-00007-plan.md", "---\nid: \"G-260101-00007\"\ntype: plan\ntitle: Plan\nstatus: current\nwork: [\"G-260101-00003\"]\n---\n\nSteps.\n")
	write(t, root, "grove/G-260101-00008-other.md", "---\nid: \"G-260101-00008\"\ntype: page\ntitle: Unrelated\n---\n\nNo.\n")
	write(t, root, "shaping.txt", "code on the shaping branch\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "shape")
	shape := git(t, root, "rev-parse", "HEAD")

	l := preview(t, root, "G-260101-00003")
	var admitted []string
	for _, m := range l.Selection.Admitted {
		admitted = append(admitted, m.ID)
	}
	if l.Base != main || l.Selection.Source != shape || strings.Join(admitted, " ") != "G-260101-00003 G-260101-00007 G-260101-00004" {
		t.Fatalf("base %s source %s admitted %v", l.Base, l.Selection.Source, admitted)
	}
	if text := strings.Join(Explain(l, func(v string) string { return v }), "\n"); !strings.Contains(text, "Admit: G-260101-00007 at grove/G-260101-00007-plan.md, record sha256:") {
		t.Fatalf("the preview names no admission:\n%s", text)
	}

	// Other work further along than proposed, or an uncommitted admitted
	// file, is refused: its code or its bytes would not come along.
	write(t, root, "grove/G-260101-00004-page.md", "---\nid: \"G-260101-00004\"\ntype: page\ntitle: Design, edited\n---\n\nIt.\n")
	if _, err := Preview(Request{Root: root, IDs: []string{"G-260101-00003"}, BudgetUSD: "3", PermissionMode: "auto"}); err == nil || !strings.Contains(err.Error(), "grove/G-260101-00004-page.md has uncommitted changes") {
		t.Fatal(err)
	}
	git(t, root, "checkout", "--", "grove")
	member(t, root, "G-260101-00005", "active")
	write(t, root, "grove/G-260101-00003.md", strings.Replace(readFile(t, root, "grove/G-260101-00003.md"), "G-260101-00004\"]", "G-260101-00004\", \"G-260101-00005\"]", 1))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "relate active work")
	if _, err := Preview(Request{Root: root, IDs: []string{"G-260101-00003"}, BudgetUSD: "3", PermissionMode: "auto"}); err == nil || !strings.Contains(err.Error(), "the selection names G-260101-00005, which the target lacks and which is active at HEAD") {
		t.Fatal(err)
	}

	// The target's copy of a selected record differs: a changed input.
	write(t, root, "grove/G-260101-00001-first.md", readFile(t, root, "grove/G-260101-00001-first.md")+"\nRefined.\n")
	git(t, root, "commit", "-qam", "refine")
	if _, err := Preview(Request{Root: root, IDs: []string{"G-260101-00001"}, BudgetUSD: "3", PermissionMode: "auto"}); err == nil || !strings.Contains(err.Error(), "grove/G-260101-00001-first.md differs between HEAD") {
		t.Fatal(err)
	}
}

// TestWorkspaceAdmissionCommit: the launch copies the admitted records onto
// a branch from the target and commits them with their provenance; the
// shaping branch's code stays behind.
func TestWorkspaceAdmissionCommit(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	main := git(t, root, "rev-parse", "main")
	git(t, root, "checkout", "-q", "-b", "shape")
	member(t, root, "G-260101-00003", "proposed")
	write(t, root, "shaping.txt", "code on the shaping branch\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "shape")
	shape := git(t, root, "rev-parse", "HEAD")
	fake(t, initLine+"\n"+resultLine("success", false))
	var facts []string
	l, err := Start(Request{Root: root, IDs: []string{"G-260101-00003"}, BudgetUSD: "1", PermissionMode: "auto", Keep: true}, now, func(f string) { facts = append(facts, f) })
	if err != nil {
		t.Fatal(err)
	}
	await(t, root, l.Attempt, Finished)
	branch := l.Branch
	if git(t, root, "rev-parse", branch+"^") != main || !l.Keep || !strings.Contains(strings.Join(facts, "\n"), "admitted: G-260101-00003 from "+shape[:7]) {
		t.Fatalf("%+v %q", l, facts)
	}
	if msg := git(t, root, "log", "-1", "--format=%B", branch); !strings.HasPrefix(msg, "chore: admit G-260101-00003 from ") || !strings.Contains(msg, "Grove-Admitted-From: "+shape) {
		t.Fatalf("message:\n%s", msg)
	}
	if files := git(t, root, "diff", "--name-only", main, branch); files != "grove/G-260101-00003.md" {
		t.Fatalf("admitted more than the record: %q", files)
	}
}

// TestWorkspaceRetiredIdentity: a branch once delivered is retired, so the
// default name moves on to a fresh one from the target, and naming it is
// refused; a submission retained whose target never advanced retires
// nothing.
func TestWorkspaceRetiredIdentity(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	git(t, root, "branch", "worktree-G-260101-00001")
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-260101-00001")
	git(t, root, "worktree", "add", "-q", wt, "worktree-G-260101-00001")
	write(t, wt, "code.txt", "delivered\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: it")
	tip := git(t, wt, "rev-parse", "HEAD")
	git(t, root, "update-ref", "refs/grove/submitted/"+tip, tip)
	// Not yet delivered: the branch still runs, where it is.
	if l := preview(t, root, "G-260101-00001"); l.Branch != "worktree-G-260101-00001" || l.Base != tip {
		t.Fatalf("%+v", l)
	}
	d := git(t, root, "commit-tree", tip+"^{tree}", "-p", "main", "-m", "feat: it\n\nGrove-Work: G-260101-00001\nGrove-Submitted: "+tip)
	git(t, root, "merge", "-q", "--ff-only", d)
	l := preview(t, root, "G-260101-00001")
	if l.Branch != "worktree-G-260101-00001-2" || l.Base != d || l.WorktreeReused {
		t.Fatalf("%+v", l)
	}
	// The board names the default branch and its checkout: the same.
	if b, err := Preview(Request{Root: root, IDs: []string{"G-260101-00001"}, BudgetUSD: "3", PermissionMode: "auto", Branch: "worktree-G-260101-00001", Worktree: wt}); err != nil || b.Branch != l.Branch || !strings.Contains(strings.Join(b.Selection.Notes, "\n"), "worktree-G-260101-00001 holds the delivery "+d[:7]+" of G-260101-00001 and is kept for inspection, never run again; this launch uses worktree-G-260101-00001-2") {
		t.Fatalf("%v %+v", err, b)
	}
	k, err := Preview(Request{Root: root, IDs: []string{"G-260101-00001"}, BudgetUSD: "3", PermissionMode: "auto", Keep: true})
	if err != nil || !k.Keep || k.Selection.Digest == l.Selection.Digest || !strings.Contains(strings.Join(Explain(k, func(v string) string { return v }), "\n"), "Keep: the workspace stays after its delivery") {
		t.Fatalf("%v %+v", err, k)
	}
	// Another name holding that delivery is refused, whatever work would
	// come next on it: a kept branch never becomes another delivery.
	git(t, root, "branch", "other", tip)
	member(t, root, "G-260101-00003", "proposed")
	git(t, root, "add", "grove")
	git(t, root, "commit", "-qm", "another")
	_, err = Preview(Request{Root: root, IDs: []string{"G-260101-00003"}, BudgetUSD: "3", PermissionMode: "auto", Branch: "other"})
	if err == nil || !strings.Contains(err.Error(), "other holds the delivery "+d[:7]+" of G-260101-00001") {
		t.Fatal(err)
	}
}
