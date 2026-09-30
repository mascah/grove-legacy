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
	"github.com/mascah/grove/internal/standing"
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
// commit, the one commit proved and each member read done, and cleanup that
// deletes only what the evidence ref keeps.
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
		"done: G-260101-00001 is done: delivered to main, proved: squashed as " + d[:7] + " from submitted tip " + submitted[:7],
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
	// Rerun, as after an interruption: the target's record says delivered.
	facts, err = run(t, root, root, true)
	if err != nil || len(facts) != 1 || facts[0] != "delivered: G-260101-00001 is already done: delivered to main" || tipOf(t, root, "main") != d {
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
	if err == nil || tipOf(t, root, "main") != d || !strings.Contains(strings.Join(facts, "\n"), "cleanup: kept worktree and branch feature: the target does not contain its tip") {
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
		// Two branches that each moved on: neither holds only the other's.
		git(t, wt, "branch", "feature2")
		git(t, root, "update-ref", "refs/heads/feature2", git(t, root, "commit-tree", "feature2^{tree}", "-p", "feature2", "-m", "chore: elsewhere"))
		write(t, wt, "grove/G-260101-00001-first.md", string(record(t, wt).Source)+"\n## Next\n\nNoted.\n")
		git(t, wt, "commit", "-qam", "docs: next")
		refused(t, root, false, "G-260101-00001 is accepted on several branches (feature, feature2); integrate needs one")
	})
	t.Run("arrived some other way", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		git(t, root, "checkout", "feature", "--", ":/")
		git(t, root, "commit", "-qm", "chore: copy feature")
		head := tipOf(t, root, "main")
		// The acceptance came with the copy, so it reads delivered: only the
		// audit could say how it got there.
		facts, err := run(t, root, root, false)
		if err != nil || len(facts) != 1 || facts[0] != "delivered: G-260101-00001 is already done: delivered to main" || tipOf(t, root, "main") != head {
			t.Fatalf("%v %q", err, facts)
		}
		if out, _ := repo.Git(root, "for-each-ref", "refs/grove/"); out != "" {
			t.Fatalf("evidence was retained for no delivery: %s", out)
		}
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

// TestIntegrateCleanupAfterAnOrdinaryMerge: work merged whole has its
// acceptance on the target and nothing retained under refs/grove, so its
// branch goes because the target contains the tip; a merge of the candidate
// alone leaves the acceptance off the target, so the work is delivered as
// usual, which brings the acceptance over and retains the branch's tip.
func TestIntegrateCleanupAfterAnOrdinaryMerge(t *testing.T) {
	t.Parallel()
	t.Run("tip contained", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		git(t, root, "merge", "-q", "--no-ff", "-m", "merge feature", "feature")
		facts, err := run(t, root, root, true)
		if err != nil || !slices.Contains(facts, "cleanup: removed worktree "+wt) || !slices.Contains(facts, "cleanup: deleted branch feature") {
			t.Fatalf("%v %q", err, facts)
		}
	})
	t.Run("candidate alone", func(t *testing.T) {
		t.Parallel()
		root, wt, candidate := fixture(t, true)
		tip := tipOf(t, root, "feature")
		git(t, root, "merge", "-q", "--no-ff", "-m", "merge the candidate", candidate)
		facts, err := run(t, root, root, true)
		if err != nil || !slices.Contains(facts, "retained: refs/grove/submitted/"+tip) || !slices.Contains(facts, "cleanup: removed worktree "+wt) || !slices.Contains(facts, "cleanup: deleted branch feature") {
			t.Fatalf("%v %q", err, facts)
		}
		if r := record(t, root); r.Status != "accepted" || r.Approved != candidate {
			t.Fatalf("the acceptance did not arrive: %+v", r)
		}
	})
}

// TestIntegrateCleanupRetainsTheSubmittedTip: a rerun whose submitted tip
// nothing retains, as in a clone that fetched the branch without refs/grove,
// retains it before it deletes the branch that held it, so the audit can
// still prove the delivery.
func TestIntegrateCleanupRetainsTheSubmittedTip(t *testing.T) {
	t.Parallel()
	root, _, _ := fixture(t, true)
	submitted := tipOf(t, root, "feature")
	if _, err := run(t, root, root, false); err != nil {
		t.Fatal(err)
	}
	git(t, root, "update-ref", "-d", "refs/grove/submitted/"+submitted)
	facts, err := run(t, root, root, true)
	if err != nil || !slices.Contains(facts, "cleanup: deleted branch feature") {
		t.Fatalf("%v %q", err, facts)
	}
	if held, _ := repo.Git(root, "rev-parse", "-q", "--verify", "refs/grove/submitted/"+submitted); strings.TrimSpace(held) != submitted {
		t.Fatalf("the branch was deleted with nothing retaining %s", submitted)
	}
}

// TestIntegrateDeliversCarriedWorkFirst: work accepted on a branch based on
// another accepted branch would carry that work in its squash with no
// delivery naming it, so it is refused until the earlier work is delivered,
// which comes from its own branch, the one the later branch contains.
func TestIntegrateDeliversCarriedWorkFirst(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	later := filepath.Join(filepath.Dir(wt), "later")
	git(t, root, "worktree", "add", "-q", "-b", "later", later, "feature")
	second := strings.NewReplacer(`"G-260101-00001"`, `"G-260101-00003"`, "title: First", "title: Second").Replace(work)
	write(t, later, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "active", 1))
	write(t, later, "later.txt", "later\n")
	git(t, later, "add", "-A")
	git(t, later, "commit", "-qm", "feat: second")
	c3 := git(t, later, "rev-parse", "HEAD")
	if _, err := update.Apply(later, update.Request{ID: "G-260101-00003", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c3}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(later, "G-260101-00003", "Yes.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	// A checkpoint on the earlier branch after later was based on it: the
	// branches diverge, and the earlier one still changes only its record.
	write(t, wt, "grove/G-260101-00001-first.md", string(record(t, wt).Source)+"\n## Next\n\nNoted.\n")
	git(t, wt, "commit", "-qam", "docs: checkpoint")
	var facts []string
	say := func(f string) { facts = append(facts, f) }
	head := tipOf(t, root, "main")
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, say); err == nil || !strings.Contains(err.Error(), "would also carry G-260101-00001's candidate") || !strings.Contains(err.Error(), "integrate G-260101-00001 first, from feature") || tipOf(t, root, "main") != head {
		t.Fatalf("%v %q", err, facts)
	}
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, say); err != nil || !slices.Contains(facts, "retained: refs/grove/submitted/"+tipOf(t, root, "feature")) {
		t.Fatalf("the earlier work, from its own branch: %v %q", err, facts)
	}
	// The squash and the checkpoint change the earlier record apart, so the
	// later branch merges the target first and is judged again, as any
	// branch behind a squash is.
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, say); err == nil || !strings.Contains(err.Error(), "conflicts with main") {
		t.Fatalf("the later work, behind: %v", err)
	}
	git(t, later, "merge", "-q", "-X", "theirs", "-m", "merge main", "main")
	if _, err := update.Feedback(later, "G-260101-00003", "Take main first.", now); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Apply(later, update.Request{ID: "G-260101-00003", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: tipOf(t, later, "HEAD")}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(later, "G-260101-00003", "Still yes.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, say); err != nil {
		t.Fatalf("the later work: %v %q", err, facts)
	}
	p, _ := project.Load(root, root)
	proofs, err := standing.Audit(context.Background(), root, "main", p.RecordDir, p.Records)
	if err != nil || len(proofs) != 2 || !proofs[0].Proved || !proofs[1].Proved {
		t.Fatalf("%v %+v %+v", err, proofs[0], proofs[len(proofs)-1])
	}
}

// TestIntegrateCarriedWorkDeliveredElsewhere: later work based on earlier
// work while it was in review carries it unaccepted; once the earlier work
// is delivered from its own branch, the refusal says to merge the target.
func TestIntegrateCarriedWorkDeliveredElsewhere(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, false)
	later := filepath.Join(filepath.Dir(wt), "later")
	git(t, root, "worktree", "add", "-q", "-b", "later", later, "feature")
	second := strings.NewReplacer(`"G-260101-00001"`, `"G-260101-00003"`, "title: First", "title: Second").Replace(work)
	write(t, later, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "active", 1))
	write(t, later, "later.txt", "later\n")
	git(t, later, "add", "-A")
	git(t, later, "commit", "-qm", "feat: second")
	c3 := git(t, later, "rev-parse", "HEAD")
	if _, err := update.Apply(later, update.Request{ID: "G-260101-00003", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c3}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	for dir, id := range map[string]string{wt: "G-260101-00001", later: "G-260101-00003"} {
		if _, err := update.Approve(dir, id, "Yes.", update.Owner, now); err != nil {
			t.Fatal(err)
		}
	}
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(string) {}); err != nil {
		t.Fatal(err)
	}
	head := tipOf(t, root, "main")
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, func(string) {}); err == nil || !strings.Contains(err.Error(), "but main already holds G-260101-00001 done, from candidate ") || !strings.Contains(err.Error(), "merge main into later") || tipOf(t, root, "main") != head {
		t.Fatal(err)
	}
}

// TestIntegrateCarriedWorkRevisedElsewhere: earlier work reopened, revised
// and delivered from its own branch after later work was based on it. The
// target's copy decides: the later branch is behind it and merges it, never
// delivering the superseded candidate, and of the earlier work's two
// acceptances, the one containing the other is the one delivered.
func TestIntegrateCarriedWorkRevisedElsewhere(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	later := filepath.Join(filepath.Dir(wt), "later")
	git(t, root, "worktree", "add", "-q", "-b", "later", later, "feature")
	second := strings.NewReplacer(`"G-260101-00001"`, `"G-260101-00003"`, "title: First", "title: Second").Replace(work)
	write(t, later, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "active", 1))
	write(t, later, "later.txt", "later\n")
	git(t, later, "add", "-A")
	git(t, later, "commit", "-qm", "feat: second")
	handoff := func(dir, id string) {
		t.Helper()
		if _, err := update.Apply(dir, update.Request{ID: id, Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: tipOf(t, dir, "HEAD")}}, Commit: true}, now, nil); err != nil {
			t.Fatal(err)
		}
		if _, err := update.Approve(dir, id, "Yes.", update.Owner, now); err != nil {
			t.Fatal(err)
		}
	}
	handoff(later, "G-260101-00003")
	if _, err := update.Feedback(wt, "G-260101-00001", "One more case.", now); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "code.txt", "the change, revised\n")
	git(t, wt, "commit", "-qam", "fix: the case")
	handoff(wt, "G-260101-00001")
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(string) {}); err != nil {
		t.Fatal(err)
	}
	head := tipOf(t, root, "main")
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, func(string) {}); err == nil || !strings.Contains(err.Error(), "but main already holds G-260101-00001 done") || !strings.Contains(err.Error(), "merge main into later") {
		t.Fatalf("the later work: %v", err)
	}
	var facts []string
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(f string) { facts = append(facts, f) }); err != nil || len(facts) != 1 || !strings.HasPrefix(facts[0], "delivered: G-260101-00001 is already done") {
		t.Fatalf("the earlier work again: %v %q", err, facts)
	}
	if tipOf(t, root, "main") != head {
		t.Fatal("main moved")
	}
	// The recovery the refusal names delivers the later work alone.
	git(t, later, "merge", "-q", "-X", "theirs", "-m", "merge main", "main")
	if _, err := update.Feedback(later, "G-260101-00003", "Take main first.", now); err != nil {
		t.Fatal(err)
	}
	handoff(later, "G-260101-00003")
	if err := Run(Request{Root: root, ID: "G-260101-00003"}, now, func(string) {}); err != nil {
		t.Fatal(err)
	}
	if got := git(t, root, "show", "main:code.txt"); got != "the change, revised" {
		t.Fatalf("code.txt: %q", got)
	}
	p, _ := project.Load(root, root)
	if proofs, err := standing.Audit(context.Background(), root, "main", p.RecordDir, p.Records); err != nil || len(proofs) != 2 || !proofs[0].Proved || !proofs[1].Proved {
		t.Fatalf("%v %+v", err, proofs)
	}
}

// TestIntegrateDeliversTheLatestAcceptance: work accepted on its branch,
// then reopened, revised and accepted again on a branch based on it, is
// delivered from the later acceptance, never the one it superseded.
func TestIntegrateDeliversTheLatestAcceptance(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	later := filepath.Join(filepath.Dir(wt), "later")
	git(t, root, "worktree", "add", "-q", "-b", "later", later, "feature")
	if _, err := update.Feedback(later, "G-260101-00001", "One more case.", now); err != nil {
		t.Fatal(err)
	}
	write(t, later, "code.txt", "the change, revised\n")
	git(t, later, "commit", "-qam", "fix: the case")
	if _, err := update.Apply(later, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: tipOf(t, later, "HEAD")}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(later, "G-260101-00001", "Now right.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	var facts []string
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(f string) { facts = append(facts, f) }); err != nil || !slices.Contains(facts, "retained: refs/grove/submitted/"+tipOf(t, root, "later")) {
		t.Fatalf("%v %q", err, facts)
	}
	if got := git(t, root, "show", "main:code.txt"); got != "the change, revised" {
		t.Fatalf("code.txt: %q", got)
	}
}

// TestIntegrateCarriedWorkOnOneBranch: two works accepted for different
// candidates on one branch cannot be delivered apart, so each refusal says
// to hand both off as one candidate, and that recovery delivers both.
func TestIntegrateCarriedWorkOnOneBranch(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	second := strings.NewReplacer(`"G-260101-00001"`, `"G-260101-00003"`, "title: First", "title: Second").Replace(work)
	write(t, wt, "grove/G-260101-00003-second.md", strings.Replace(second, "%s", "active", 1))
	write(t, wt, "later.txt", "later\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: second")
	c3 := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00003", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c3}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := update.Approve(wt, "G-260101-00003", "Yes.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	head := tipOf(t, root, "main")
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if err := Run(Request{Root: root, ID: id}, now, func(string) {}); err == nil || !strings.Contains(err.Error(), "hand them off as one candidate") || strings.Contains(err.Error(), "first") || tipOf(t, root, "main") != head {
			t.Fatalf("%s: %v", id, err)
		}
	}
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if _, err := update.Feedback(wt, id, "Deliver them together.", now); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if _, err := update.Apply(wt, update.Request{ID: id, Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c3}}, Commit: true}, now, nil); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"G-260101-00001", "G-260101-00003"} {
		if _, err := update.Approve(wt, id, "Together.", update.Owner, now); err != nil {
			t.Fatal(err)
		}
	}
	if err := Run(Request{Root: root, ID: "G-260101-00001"}, now, func(string) {}); err != nil {
		t.Fatal(err)
	}
	p, _ := project.Load(root, root)
	if proofs, err := standing.Audit(context.Background(), root, "main", p.RecordDir, p.Records); err != nil || len(proofs) != 2 || !proofs[0].Proved || !proofs[1].Proved {
		t.Fatalf("%v %+v", err, proofs)
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
// decides, so work reopened on a kept branch is refused while in review and
// is a second delivery once accepted again, never "already done" from the
// target's earlier copy. The kept branch conflicts with its own earlier
// delivery until it merges the target, as any branch does; feedback after a
// squash runs on the target, which never holds the candidate.
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
	// What was delivered is the branch's own, so its side wins.
	git(t, wt, "merge", "-q", "-X", "ours", "-m", "merge main", "main")
	c2 := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: c2}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	refused(t, root, false, "G-260101-00001 is in review on feature but not accepted")
	if _, err := update.Approve(wt, "G-260101-00001", "Now.", update.Owner, now); err != nil {
		t.Fatal(err)
	}
	facts, err := run(t, root, root, false)
	if err != nil || !strings.HasPrefix(facts[len(facts)-1], "done: G-260101-00001 is done: delivered to main, proved: squashed as ") || git(t, root, "show", "main:code.txt") != "the change, more" {
		t.Fatalf("the second acceptance: %v %q", err, facts)
	}
	if git(t, root, "rev-parse", "main^") != d1 || !strings.HasPrefix(git(t, root, "log", "-1", "--format=%s", "main"), "feat: first") {
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

// TestIntegrateLeavesAnAlteredDeliveryToTheAudit: a delivery amended after
// the fact still reads done, since the target holds the acceptance, and is
// never delivered again; only the audit, on request, says it is not proved.
func TestIntegrateLeavesAnAlteredDeliveryToTheAudit(t *testing.T) {
	t.Parallel()
	root, _, _ := fixture(t, true)
	if _, err := run(t, root, root, false); err != nil {
		t.Fatal(err)
	}
	write(t, root, "extra.txt", "amended in\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "--amend", "--no-edit")
	head := tipOf(t, root, "main")
	if facts, err := run(t, root, root, false); err != nil || len(facts) != 1 || tipOf(t, root, "main") != head {
		t.Fatalf("%v %q", err, facts)
	}
	p, _ := project.Load(root, root)
	proofs, err := standing.Audit(context.Background(), root, "main", p.RecordDir, p.Records)
	if err != nil || len(proofs) != 1 || proofs[0].Proved || !strings.Contains(proofs[0].Why, "is not what merging") {
		t.Fatalf("%v %+v", err, proofs)
	}
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
