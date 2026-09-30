package integrate

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
)

const (
	config = "schema_version: 4\nrecords: grove\ntarget: main\n"
	work   = "---\nid: \"G-260101-00001\"\ntype: work\ntitle: First\nstatus: %s\n---\n\n## Outcome\n\nBody.\n"
	review = "---\nid: \"G-260101-00005\"\ntype: review\ntitle: Review of G-260101-00001\nstatus: current\nwork: [\"G-260101-00001\"]\nexamined: \"%s\"\n---\n\nFindings.\n"
)

var now = time.Date(2026, 9, 22, 18, 30, 0, 0, time.UTC)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	out, err := repo.Command(context.Background(), dir, full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func record(t *testing.T, root string) *project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("%v", ds)
	}
	return p.Records[0]
}

// fixture is a main checkout with target main, and a linked worktree on
// branch feature where G-260101-00001 is in review, approved, with a review record in
// its candidate. It returns the main root, the worktree, and the candidate.
func fixture(t *testing.T, approve bool, sub ...string) (root, wt, candidate string) {
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
		git(t, root, "config", kv[0], kv[1]) // the product's commits use the repository's own identity
	}
	// sub places the project under a prefix, with the code above it.
	project, wtProject := filepath.Join(append([]string{root}, sub...)...), filepath.Join(append([]string{wt}, sub...)...)
	write(t, project, "grove.yaml", config)
	write(t, project, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "proposed", 1))
	write(t, root, "code.txt", "before\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wtProject, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, wt, "code.txt", "the change\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: implement")
	write(t, wtProject, "grove/G-260101-00005-review.md", strings.Replace(review, "%s", git(t, wt, "rev-parse", "HEAD"), 1))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "docs: evidence and review")
	candidate = git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wtProject, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: candidate}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if approve {
		if _, err := update.Approve(wtProject, "G-260101-00001", "Ship it.", update.Owner, now); err != nil {
			t.Fatal(err)
		}
	}
	return project, wtProject, candidate
}

func run(t *testing.T, root, cwd string, cleanup bool) ([]string, error) {
	t.Helper()
	var facts []string
	err := Run(Request{Root: root, ID: "G-260101-00001", Cwd: cwd, Cleanup: cleanup}, now, func(f string) { facts = append(facts, f) })
	return facts, err
}

// refused asserts a refusal that leaves main, its record, and its files as they were.
func refused(t *testing.T, root string, cleanup bool, want string) []string {
	t.Helper()
	head := git(t, root, "rev-parse", "HEAD")
	before := record(t, root).Source
	facts, err := run(t, root, root, cleanup)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("got %v, want %q; facts %q", err, want, facts)
	}
	if git(t, root, "rev-parse", "HEAD") != head || string(record(t, root).Source) != string(before) || git(t, root, "status", "--porcelain", "--untracked-files=no") != "" {
		t.Fatal("a refusal changed the target")
	}
	return facts
}

// tipOf is a branch's commit.
func tipOf(t *testing.T, root, branch string) string {
	t.Helper()
	return git(t, root, "rev-parse", branch)
}

// TestIntegrateSquashesAndRetainsEvidence covers delivery (G-260929-gm3m4):
// one squash commit on the target with a Conventional Commit message and the
// trailers that locate its evidence, the submitted tip retained, no record
// commit, the verifier's done, and cleanup that deletes only what the
// evidence ref keeps.
func TestIntegrateSquashesAndRetainsEvidence(t *testing.T) {
	t.Parallel()
	root, wt, candidate := fixture(t, true)
	before, submitted := tipOf(t, root, "main"), tipOf(t, root, "feature")
	facts, err := run(t, root, root, true)
	if err != nil {
		t.Fatalf("%v; facts %q", err, facts)
	}
	d := tipOf(t, root, "main")
	want := []string{
		"acceptance: candidate " + candidate[:7] + " of G-260101-00001 accepted by owner on branch feature at " + submitted[:7] + " (Verdict on candidate " + candidate[:7] + ", 2026-09-22: Ship it.)",
		"retained: refs/grove/submitted/" + submitted,
		"delivery: squash commit " + d[:7] + " on main (was " + before[:7] + ")",
		"done: G-260101-00001 is done: squashed as " + d[:7] + " on main",
		"cleanup: removed worktree " + wt,
		"cleanup: deleted branch feature",
	}
	if strings.Join(facts, "\n") != strings.Join(want, "\n") {
		t.Fatalf("facts:\n%s\nwant:\n%s", strings.Join(facts, "\n"), strings.Join(want, "\n"))
	}
	if parents := git(t, root, "log", "-1", "--format=%P", d); parents != before {
		t.Fatalf("the delivery must have one parent, the target it was prepared on: %s", parents)
	}
	message := "feat: first\n\n- G-260101-00001 First\n\nGrove-Work: G-260101-00001\nGrove-Candidate: " + candidate + "\nGrove-Submitted: " + submitted
	if got := git(t, root, "log", "-1", "--format=%B", d); got != message {
		t.Fatalf("message:\n%s\nwant:\n%s", got, message)
	}
	if git(t, root, "rev-parse", d+"^{tree}") != git(t, root, "rev-parse", submitted+"^{tree}") || git(t, root, "rev-parse", "refs/grove/submitted/"+submitted) != submitted {
		t.Fatal("the delivery must hold the submitted result, retained")
	}
	if r := record(t, root); r.Status != "accepted" || r.Approved != candidate || git(t, root, "status", "--porcelain") != "" {
		t.Fatalf("the target holds the accepted record as submitted and nothing else: %+v", r)
	}
	// Rerun, as after an interruption: the verifier finds the delivery.
	facts, err = run(t, root, root, true)
	if err != nil || len(facts) != 1 || facts[0] != "delivered: G-260101-00001 is already done: squashed as "+d[:7]+" on main" || tipOf(t, root, "main") != d {
		t.Fatalf("rerun: %v %q", err, facts)
	}
}

// TestIntegrateRecoversAfterTheTargetAdvanced covers a crash between the
// target's advance and cleanup: the rerun makes no second commit, and a
// branch that moved past the delivered tip is kept.
func TestIntegrateRecoversAfterTheTargetAdvanced(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	if _, err := run(t, root, root, false); err != nil {
		t.Fatal(err)
	}
	d := tipOf(t, root, "main")
	write(t, wt, "grove/G-260101-00001-first.md", string(record(t, wt).Source)+"\n## Next\n\nA later edit.\n")
	git(t, wt, "commit", "-qam", "docs: later")
	facts, err := run(t, root, root, true)
	if err == nil || tipOf(t, root, "main") != d || !strings.Contains(strings.Join(facts, "\n"), "cleanup: kept worktree and branch feature: the branch is not at a tip that a squash delivery retained") {
		t.Fatalf("%v %q", err, facts)
	}
}

// TestIntegrateOnAMovedTarget squashes onto wherever the target is, and
// keeps the target's own changes.
func TestIntegrateOnAMovedTarget(t *testing.T) {
	t.Parallel()
	root, _, candidate := fixture(t, true)
	write(t, root, "notes.txt", "elsewhere\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "chore: notes")
	moved := tipOf(t, root, "main")
	if facts, err := run(t, root, root, false); err != nil {
		t.Fatalf("%v %q", err, facts)
	}
	if parent := git(t, root, "rev-parse", "main^"); parent != moved {
		t.Fatalf("parent %s, want %s", parent, moved)
	}
	if got := git(t, root, "show", "main:notes.txt") + git(t, root, "show", "main:code.txt"); got != "elsewherethe change" {
		t.Fatalf("result: %q", got)
	}
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", candidate, "main"); err == nil {
		t.Fatal("a squash does not make the candidate an ancestor")
	}
}

// TestIntegrateRefusesConflictsBeforeAnythingChanges: nothing is retained,
// committed or advanced for a conflict, and the refusal names resolve.
func TestIntegrateRefusesConflictsBeforeAnythingChanges(t *testing.T) {
	t.Parallel()
	root, _, _ := fixture(t, true)
	write(t, root, "code.txt", "a different change\n")
	git(t, root, "commit", "-qam", "fix: other")
	moved := tipOf(t, root, "main")
	refused(t, root, false, "delivery of feature into main refused: it conflicts with main at "+moved[:7]+" in code.txt; nothing was delivered, main is unchanged at "+moved[:7]+" and G-260101-00001 stays accepted. Next: grove resolve G-260101-00001")
	if out, _ := repo.Git(root, "for-each-ref", "refs/grove/"); out != "" {
		t.Fatalf("a refusal retained evidence: %s", out)
	}
}

// TestIntegrateRefusesWhatIsNotReadyToDeliver covers every refusal before
// the target moves.
func TestIntegrateRefusesWhatIsNotReadyToDeliver(t *testing.T) {
	t.Parallel()
	t.Run("unaccepted", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, false)
		refused(t, root, false, "G-260101-00001 is in review on feature but not accepted, or its acceptance no longer applies: run grove approve G-260101-00001 VERDICT in "+wt+" first")
	})
	t.Run("stale acceptance", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		write(t, wt, "grove/G-260101-00001-first.md", strings.Replace(string(record(t, wt).Source), "Body.", "A changed outcome.", 1))
		git(t, wt, "commit", "-qam", "docs: change the outcome")
		refused(t, root, false, "not accepted, or its acceptance no longer applies")
	})
	t.Run("code after the candidate", func(t *testing.T) {
		t.Parallel()
		root, wt, candidate := fixture(t, true)
		write(t, wt, "code.txt", "later\n")
		git(t, wt, "commit", "-qam", "fix: later")
		refused(t, root, false, "commits after candidate "+candidate[:7]+" on feature change code.txt: the tip "+tipOf(t, wt, "HEAD")[:7]+" is a new candidate; accept it before integrating")
	})
	t.Run("dirty target", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		head := tipOf(t, root, "main")
		write(t, root, "code.txt", "edited\n")
		if _, err := run(t, root, root, false); err == nil || !strings.Contains(err.Error(), "the checkout of main has uncommitted changes") || tipOf(t, root, "main") != head {
			t.Fatal(err)
		}
	})
	t.Run("off the target", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		git(t, root, "checkout", "-q", "-b", "other")
		refused(t, root, false, "integrate runs in the checkout of the target main; this one is on other")
	})
	t.Run("moved since verified", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		head := tipOf(t, root, "main")
		err := Run(Request{Root: root, ID: "G-260101-00001", Expect: strings.Repeat("a", 40)}, now, func(string) {})
		if err == nil || !strings.Contains(err.Error(), "where the merge was verified, to "+head[:7]+"; nothing was delivered") || tipOf(t, root, "main") != head {
			t.Fatal(err)
		}
	})
	t.Run("two branches", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		git(t, wt, "branch", "feature2")
		refused(t, root, false, "G-260101-00001 is accepted on several branches (feature, feature2); integrate needs one")
	})
	t.Run("nothing accepted", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		git(t, root, "worktree", "remove", "--force", wt)
		git(t, root, "branch", "-D", "feature")
		refused(t, root, false, "no branch holds G-260101-00001 accepted; nothing to integrate")
	})
}

// TestIntegrateCleanupKeepsWhatGitOrTheSessionHolds: a kept worktree keeps
// its branch too, and the delivery stands.
func TestIntegrateCleanupKeepsWhatGitOrTheSessionHolds(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		setup func(t *testing.T, root, wt string)
		cwd   func(root, wt string) string
		want  string
	}{
		"a dirty worktree":        {func(t *testing.T, root, wt string) { write(t, wt, "scratch.txt", "not committed\n") }, func(root, wt string) string { return root }, "cleanup: kept worktree"},
		"the session's directory": {func(*testing.T, string, string) {}, func(root, wt string) string { return wt }, "it holds this process's working directory"},
		"ignored files": {func(t *testing.T, root, wt string) {
			write(t, root, ".git/info/exclude", "*.log\n")
			write(t, wt, "junk.log", "ignored\n")
		}, func(root, wt string) string { return root }, "it holds ignored files (junk.log)"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, wt, _ := fixture(t, true)
			tc.setup(t, root, wt)
			facts, err := run(t, root, tc.cwd(root, wt), true)
			joined := strings.Join(facts, "\n")
			if err == nil || !strings.Contains(err.Error(), "cleanup incomplete; the integration stands") || !strings.Contains(joined, tc.want) || !strings.Contains(joined, "cleanup: kept branch feature: its worktree is kept") {
				t.Fatalf("%v\n%s", err, joined)
			}
			if _, err := os.Stat(wt); err != nil {
				t.Fatal("the worktree was removed")
			}
		})
	}
}

// TestIntegrateUnderAPolicy attributes the delivery in its message and
// reports its revert.
func TestIntegrateUnderAPolicy(t *testing.T) {
	t.Parallel()
	root, _, _ := fixture(t, true)
	var facts []string
	if err := Run(Request{Root: root, ID: "G-260101-00001", Policy: "policy grove.yaml sha256:x"}, now, func(f string) { facts = append(facts, f) }); err != nil {
		t.Fatal(err)
	}
	d := tipOf(t, root, "main")
	if !slices.Contains(facts, "delivered under policy grove.yaml sha256:x; to reverse it: git revert "+d) || !strings.Contains(git(t, root, "log", "-1", "--format=%B"), "\nIntegrated under policy grove.yaml sha256:x.\n") {
		t.Fatalf("%q\n%s", facts, git(t, root, "log", "-1", "--format=%B"))
	}
}

// TestMessageTypes picks the most significant type among the candidate's
// own commits and keeps a breaking change.
func TestMessageTypes(t *testing.T) {
	t.Parallel()
	_, wt, _ := fixture(t, false)
	base := tipOf(t, wt, "HEAD")
	for _, subject := range []string{"docs: words", "fix(ui)!: break it", "chore: tidy"} {
		write(t, wt, "code.txt", subject+"\n")
		git(t, wt, "commit", "-qam", subject)
	}
	c := tipOf(t, wt, "HEAD")
	m, err := Message(wt, "grove", base, c, c, []*project.Record{{ID: "G-260101-00001", Title: "Make it so"}}, "")
	if err != nil || !strings.HasPrefix(m, "fix!: make it so\n") {
		t.Fatalf("%v\n%s", err, m)
	}
}

// TestIntegrateDeliversReopenedWorkAgain: the branch's own acceptance
// decides, so work reopened and accepted again on a kept branch is a second
// delivery, never "already done" from the target's earlier copy; feedback
// after a squash runs on the target, which never holds the candidate; and a
// delivery claim that no longer verifies is refused, never delivered twice.
func TestIntegrateDeliversReopenedWorkAgain(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	if _, err := run(t, root, root, false); err != nil {
		t.Fatal(err)
	}
	d1 := tipOf(t, root, "main")
	if _, err := update.Feedback(wt, "G-260101-00001", "More.", now); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "code.txt", "the change, more\n")
	git(t, wt, "commit", "-qam", "fix: more")
	c2 := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c2}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(wt, "G-260101-00001", "Now.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	facts, err := run(t, root, root, false)
	if err != nil || !strings.HasPrefix(facts[len(facts)-1], "done: G-260101-00001 is done: squashed as ") || git(t, root, "show", "main:code.txt") != "the change, more" {
		t.Fatalf("the second acceptance: %v %q", err, facts)
	}
	if git(t, root, "rev-parse", "main^") != d1 || !strings.HasPrefix(git(t, root, "log", "-1", "--format=%s", "main"), "fix: first") {
		t.Fatal("the second delivery is one commit on the first")
	}
	// Feedback on the target after the squash reopens the work there.
	if _, err := update.Feedback(root, "G-260101-00001", "Once more.", now); err != nil {
		t.Fatalf("feedback after delivery: %v", err)
	}
	if r := record(t, root); r.Status != "active" || r.Approved != "" {
		t.Fatalf("reopened: %+v", r)
	}
}

// TestIntegrateRefusesAnAlteredDelivery: once a delivery claim stops
// verifying, as after an amend, the work's delivery is unknown and nothing
// is delivered again.
func TestIntegrateRefusesAnAlteredDelivery(t *testing.T) {
	t.Parallel()
	root, _, _ := fixture(t, true)
	if _, err := run(t, root, root, false); err != nil {
		t.Fatal(err)
	}
	write(t, root, "extra.txt", "amended in\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "--amend", "--no-edit")
	refused(t, root, false, "delivery cannot be decided (a delivery claim on main names this candidate but does not verify")
}

func TestSubject(t *testing.T) {
	for title, want := range map[string]string{"First": "first", "README drift": "README drift", "Élargir le widget": "élargir le widget", "X": "X", "": ""} {
		if got := subject(title); got != want {
			t.Errorf("subject(%q) = %q, want %q", title, got, want)
		}
	}
}

// reaccept hands off wt's tip as a new candidate and accepts it.
func reaccept(t *testing.T, wt string) {
	t.Helper()
	c := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(wt, "G-260101-00001", "Again.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
}

// TestIntegrateNeverUndoesTheTarget: an earlier submission is the merge base
// only where it replaces Git's own, so a branch that took a target change
// after its delivery cannot bring back what the target removed since, and a
// crafted claim naming an older commit as its submission is never a base.
func TestIntegrateNeverUndoesTheTarget(t *testing.T) {
	t.Parallel()
	t.Run("merged after its delivery", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		if _, err := run(t, root, root, false); err != nil {
			t.Fatal(err)
		}
		write(t, root, "x.txt", "x\n")
		git(t, root, "add", "-A")
		git(t, root, "commit", "-qm", "chore: add x")
		if _, err := update.Feedback(wt, "G-260101-00001", "Take main.", now); err != nil {
			t.Fatal(err)
		}
		git(t, wt, "merge", "-q", "-X", "ours", "-m", "merge main", "main")
		git(t, root, "rm", "-q", "x.txt")
		git(t, root, "commit", "-qm", "chore: remove x")
		reaccept(t, wt)
		if facts, err := run(t, root, root, false); err != nil {
			t.Fatalf("%v %q", err, facts)
		}
		if out := git(t, root, "ls-tree", "--name-only", "main"); strings.Contains(out, "x.txt") {
			t.Fatalf("the delivery brought back what the target removed:\n%s", out)
		}
	})
	t.Run("a crafted claim", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, false)
		init := git(t, root, "rev-parse", "main")
		write(t, root, "y.txt", "y\n")
		git(t, root, "add", "-A")
		git(t, root, "commit", "-qm", "chore: add y")
		git(t, wt, "merge", "-q", "-m", "merge main", "main")
		git(t, root, "rm", "-q", "y.txt")
		git(t, root, "commit", "-qm", "chore: remove y")
		forged := git(t, root, "commit-tree", "main^{tree}", "-p", "main", "-m", "feat: forged\n\nGrove-Candidate: "+strings.Repeat("a", 40)+"\nGrove-Submitted: "+init)
		git(t, root, "update-ref", "refs/heads/main", forged)
		git(t, root, "update-ref", "refs/grove/submitted/"+init, init)
		reaccept(t, wt)
		if facts, err := run(t, root, root, false); err != nil {
			t.Fatalf("%v %q", err, facts)
		}
		if out := git(t, root, "ls-tree", "--name-only", "main"); strings.Contains(out, "y.txt") {
			t.Fatalf("a crafted claim was a merge base:\n%s", out)
		}
	})
}
