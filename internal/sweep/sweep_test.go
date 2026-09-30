package sweep

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
)

// TestMain lets the test binary be the owner an attempt starts, sweeping after
// it as cmd/grove does, and makes every attempt's provider the fake handedOff
// describes unless a test sets its own.
func TestMain(m *testing.M) {
	if dir := os.Getenv(attempt.OwnerEnv); dir != "" {
		os.Exit(Own(dir))
	}
	dir, err := os.MkdirTemp("", "grove-sweep-fake-")
	if err != nil {
		panic(err)
	}
	fake := filepath.Join(dir, "claude")
	// The worktree is <dir>/repo/.claude/worktrees/<branch>, and <dir> holds
	// what to hand off.
	script := "#!/bin/sh\n[ \"$1\" = --version ] && { echo 'fake 0.1'; exit 0; }\nset -e\nD=$(cd ../../../.. && pwd)\n" +
		"cp -R \"$D/files/.\" .\ngit add -A\ngit commit -q --allow-empty -m 'feat: the change'\n" +
		"sed \"s/EXAMINED/$(git rev-parse HEAD)/\" \"$D/review.md\" > grove/G-260101-00005-review.md\ngit add -A\ngit commit -qm 'docs: review'\n" +
		"f=grove/G-260101-00001-first.md\nawk -v c=\"$(git rev-parse HEAD)\" '/^candidate:/ {next} /^status:/ {print \"status: review\"; print \"candidate: \\\"\" c \"\\\"\"; next} {print}' $f > $f.new\nmv $f.new $f\ngit commit -qam 'docs: hand off'\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		panic(err)
	}
	os.Setenv(attempt.ClaudeEnv, fake)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

const (
	work   = "---\nid: \"G-260101-00001\"\ntype: work\ntitle: First\nstatus: %s\n---\n\n## Outcome\n\nBody.\n"
	review = "---\nid: \"G-260101-00005\"\ntype: review\ntitle: Review of G-260101-00001\nstatus: current\nwork: [\"G-260101-00001\"]\nexamined: \"%s\"\n---\n\nFindings: none.\n\n%s\n"
	policy = "schema_version: 4\nrecords: grove\ntarget: main\npolicy:\n  budget: 5\n  resolve:\n    budget: 1\n  approve:\n    verify: [%s]\n    max_lines: 30\n    never: [grove.yaml, secret/**]\n  integrate: true\nrun:\n  permission_mode: acceptEdits\n"
)

var now = time.Date(2026, 9, 25, 18, 30, 0, 0, time.UTC)

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

// fixture is main with config, and branch worktree-G-260101-00001 in its checkout
// where G-260101-00001 is in review: its candidate writes files over code.txt's
// "base", and a review record, added after it, examined it and closes with
// closing. It returns main's checkout and the branch's.
func fixture(t *testing.T, config string, files map[string]string, closing string) (root, wt string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt = filepath.Join(dir, "repo"), filepath.Join(dir, "repo", ".claude", "worktrees", "worktree-G-260101-00001")
	git(t, dir, "init", "-q", "-b", "main", root)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	write(t, root, ".gitignore", ".claude/worktrees/\n")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, root, attempt.SkillPath, "---\nname: grove-work\n---\n")
	write(t, root, "code.txt", "base\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-260101-00001", wt)
	for name, content := range files {
		write(t, wt, name, content)
	}
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: the change")
	examined := git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "grove/G-260101-00005-review.md", strings.Replace(strings.Replace(review, "%s", examined, 1), "%s", closing, 1))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "docs: review")
	candidate := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: candidate}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	return root, wt
}

func record(t *testing.T, dir string) *project.Record {
	t.Helper()
	p, ds := project.Load(dir, dir)
	if len(ds) != 0 {
		t.Fatalf("%v", ds)
	}
	for _, r := range p.Records {
		if r.ID == "G-260101-00001" {
			return r
		}
	}
	t.Fatal("no G-260101-00001")
	return nil
}

func sweep(t *testing.T, root string) (*Sweep, []string) {
	t.Helper()
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	var facts []string
	s.Run(now, func(f string) { facts = append(facts, f) })
	return s, facts
}

func TestSweepIntegratesACandidateInsideThePolicy(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t, strings.Replace(policy, "%s", "grep -q change code.txt, test -f grove/G-260101-00005-review.md", 1), map[string]string{"code.txt": "the change\n"}, ClosingLine)
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Items) != 1 || s.Items[0].Act != Integrate || !strings.Contains(s.Items[0].Why, "review G-260101-00005 has no open finding; 14 changed lines") {
		t.Fatalf("plan: %+v", s.Items)
	}
	if !strings.HasPrefix(s.Attribution, "policy grove.yaml sha256:") {
		t.Fatalf("attribution %q", s.Attribution)
	}
	var facts []string
	s.Run(now, func(f string) { facts = append(facts, f) })
	joined := strings.Join(facts, "\n")
	for _, want := range []string{"G-260101-00001: verified: the merge of", "G-260101-00001: approved under policy grove.yaml sha256:", "G-260101-00001: delivery: squash commit ", "G-260101-00001: done: G-260101-00001 is done: squashed as "} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in:\n%s", want, joined)
		}
	}
	r := record(t, root)
	if r.Status != "accepted" || !update.Delegated(r) || r.ApprovedBy != "policy "+strings.TrimPrefix(s.Attribution, "policy grove.yaml ") {
		t.Fatalf("record on main:\n%s", r.Source)
	}
	for _, want := range []string{
		": delegated under " + s.Attribution + ": review G-260101-00005 examined ",
		"verification passed (grep -q change code.txt; test -f grove/G-260101-00005-review.md); no Grove attempt is recorded as producing it",
	} {
		if !strings.Contains(string(r.Source), want) {
			t.Fatalf("missing %q in:\n%s", want, r.Source)
		}
	}
	if !strings.Contains(git(t, root, "log", "-1", "--format=%B", "main"), "Integrated under "+s.Attribution+".") {
		t.Fatal("the delivery must name the policy")
	}
	if out := git(t, root, "worktree", "list", "--porcelain"); strings.Count(out, "worktree ") != 2 {
		t.Fatalf("the verification worktree remains:\n%s", out)
	}
	// Integrated, it is no longer a candidate in review anywhere but on its
	// branch, which the target now holds.
	if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "G-260101-00001: skipped: ") {
		t.Fatalf("second sweep: %q", facts)
	}
}

func TestSweepLeavesAFailedVerificationUnchanged(t *testing.T) {
	t.Parallel()
	root, wt := fixture(t, strings.Replace(policy, "%s", "'true', 'echo broken; exit 3'", 1), map[string]string{"code.txt": "the change\n"}, ClosingLine)
	main, branch := git(t, root, "rev-parse", "main"), git(t, root, "rev-parse", "worktree-G-260101-00001")
	_, facts := sweep(t, root)
	if len(facts) != 1 || !strings.Contains(facts[0], "G-260101-00001: waits: verification of the merge with main at ") || !strings.Contains(facts[0], "echo broken; exit 3: exit status 3; last output: broken; nothing was approved") {
		t.Fatalf("facts %q", facts)
	}
	if git(t, root, "rev-parse", "main") != main || git(t, root, "rev-parse", "worktree-G-260101-00001") != branch || record(t, wt).Approved != "" {
		t.Fatal("a failed verification changed something")
	}
}

func TestSweepWaitsOutsideThePolicy(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		files   map[string]string
		closing string
		want    string
	}{
		"a never path":    {map[string]string{"secret/key": "x\n"}, ClosingLine, "changes secret/key, which never matches (secret/**)"},
		"the policy":      {map[string]string{"grove.yaml": strings.Replace(policy, "%s", "'true'", 1) + "# more\n"}, ClosingLine, "changes grove.yaml, which never matches (grove.yaml)"},
		"too many lines":  {map[string]string{"code.txt": strings.Repeat("line\n", 30)}, ClosingLine, "lines, over the policy's 30"},
		"an open finding": {map[string]string{"code.txt": "the change\n"}, "Open findings: 1", `review G-260101-00005 of candidate `},
		"a binary":        {map[string]string{"blob": "\x00\x01"}, ClosingLine, "changes the binary file blob"},
		"no closing line": {map[string]string{"code.txt": "the change\n"}, "", `does not end with "Open findings: none"`},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, _ := fixture(t, strings.Replace(policy, "%s", "'true'", 1), c.files, c.closing)
			main := git(t, root, "rev-parse", "main")
			_, facts := sweep(t, root)
			if len(facts) != 1 || !strings.HasPrefix(facts[0], "G-260101-00001: waits: ") || !strings.Contains(facts[0], c.want) {
				t.Fatalf("facts %q, want %q", facts, c.want)
			}
			if git(t, root, "rev-parse", "main") != main {
				t.Fatal("main moved")
			}
		})
	}
	t.Run("approved by the owner", func(t *testing.T) {
		t.Parallel()
		root, wt := fixture(t, strings.Replace(policy, "%s", "'true'", 1), map[string]string{"code.txt": "x\n"}, ClosingLine)
		if _, err := update.Approve(wt, "G-260101-00001", "Mine.", update.Owner, now); err != nil {
			t.Fatal(err)
		}
		if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "waits: accepted by the owner; delivering it is the owner's: grove integrate G-260101-00001") {
			t.Fatalf("facts %q", facts)
		}
	})
	t.Run("no approve section", func(t *testing.T) {
		t.Parallel()
		root, _ := fixture(t, "schema_version: 4\nrecords: grove\ntarget: main\npolicy:\n  budget: 1\n  resolve: {}\n", map[string]string{"code.txt": "x\n"}, ClosingLine)
		if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "the policy does not approve") {
			t.Fatalf("facts %q", facts)
		}
	})
}

func TestSweepRefusesWithoutAPolicy(t *testing.T) {
	t.Parallel()
	root, wt := fixture(t, "schema_version: 4\nrecords: grove\ntarget: main\n", map[string]string{"code.txt": "x\n"}, ClosingLine)
	if _, err := Plan(root); err == nil || !strings.Contains(err.Error(), "grove.yaml has no policy: nothing is automatic") {
		t.Fatalf("got %v", err)
	}
	if _, err := Plan(wt); err == nil || !strings.Contains(err.Error(), "sweep runs in the checkout of the target main") {
		t.Fatalf("from the branch: %v", err)
	}
}

func TestSweepResolvesAConflictOncePerTargetCommit(t *testing.T) {
	if testing.Short() {
		t.Skip("waits on a fake provider's attempt")
	}
	// The candidate's own grove.yaml asks for another permission mode, which
	// a delegated attempt never takes.
	config := strings.Replace(policy, "%s", "'true'", 1)
	root, wt := fixture(t, config, map[string]string{"code.txt": "branch\n", "grove.yaml": strings.Replace(config, "acceptEdits", "bypassPermissions\n  model: branch-model\n  effort: max", 1)}, ClosingLine)
	write(t, root, "code.txt", "main\n")
	git(t, root, "commit", "-qam", "main moves")
	fake := filepath.Join(t.TempDir(), "claude")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n[ \"$1\" = --version ] && { echo 'fake 0.1'; exit 0; }\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(attempt.ClaudeEnv, fake)

	s, facts := sweep(t, root)
	if len(s.Items) != 1 || s.Items[0].Act != Resolve || !strings.Contains(s.Items[0].Why, "in code.txt; one resolution attempt, budget 1 USD") {
		t.Fatalf("plan %+v", s.Items)
	}
	if last := facts[len(facts)-1]; !strings.Contains(last, "resolution attempt ") || !strings.Contains(last, "started under "+s.Attribution+", budget 1 USD") {
		t.Fatalf("facts %q", facts)
	}
	r := record(t, wt)
	if r.Status != "active" || !strings.Contains(string(r.Source), ": delegated under "+s.Attribution+", budget 1 USD: conflicts with main at ") {
		t.Fatalf("record:\n%s", r.Source)
	}
	deadline := time.Now().Add(20 * time.Second)
	for {
		views, err := attempt.List(root, "G-260101-00001")
		if err != nil {
			t.Fatal(err)
		}
		if len(views) == 1 {
			if l := views[0].Launch; l.PermissionMode != "acceptEdits" || l.Model != "" || l.Effort != "" {
				t.Fatalf("the attempt took mode %q, model %q, effort %q from the candidate", l.PermissionMode, l.Model, l.Effort)
			}
		}
		if len(views) == 1 && views[0].Status == attempt.Finished {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt did not finish: %+v", views)
		}
		time.Sleep(50 * time.Millisecond)
	}
	// The attempt handed the same candidate back without resolving: the
	// conflict with the same main commit waits for the owner.
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00001", Set: []update.Field{{Name: "status", Value: "review"}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "waits: ") || !strings.Contains(facts[0], "again after a resolution of main at ") {
		t.Fatalf("second sweep %q", facts)
	}
}

func TestSweepKeepsResolutionsInsideTheBudget(t *testing.T) {
	t.Parallel()
	config := strings.Replace(strings.Replace(policy, "%s", "'true'", 1), "  resolve:\n    budget: 1\n", "  resolve:\n    budget: 6\n", 1)
	root, _ := fixture(t, config, map[string]string{"code.txt": "branch\n"}, ClosingLine)
	write(t, root, "code.txt", "main\n")
	git(t, root, "commit", "-qam", "main moves")
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Items) != 1 || s.Items[0].Act != Wait || !strings.Contains(s.Items[0].Why, "a 6 USD attempt would pass the policy's 5 USD for one sweep") {
		t.Fatalf("plan %+v", s.Items)
	}
	// Without a permission mode beside the policy, resolve would refuse it,
	// so the plan says it waits.
	write(t, root, "grove.yaml", strings.Replace(config, "run:\n  permission_mode: acceptEdits\n", "", 1))
	git(t, root, "commit", "-qam", "no mode")
	if s, err = Plan(root); err != nil || len(s.Items) != 1 || s.Items[0].Act != Wait || !strings.Contains(s.Items[0].Why, "no permission mode") {
		t.Fatalf("plan %+v, %v", s.Items, err)
	}
}

// Every current review covering the candidate must close with no open
// finding: an older clean one does not outvote a newer one that found
// something.
func TestSweepHeedsEveryReviewOfTheCandidate(t *testing.T) {
	t.Parallel()
	r := &project.Record{ID: "G-260101-00001", Candidate: "abcdef1234"}
	rev := func(id, closing string) *project.Record {
		return &project.Record{ID: id, Type: "review", Status: "current", Work: []string{"G-260101-00001"}, Examined: "abcdef1234", Source: []byte("Findings.\n\n" + closing + "\n")}
	}
	s := &Sweep{}
	if got, why := s.review(context.Background(), r, []*project.Record{rev("G-260101-00005", ClosingLine), rev("G-260101-00007", "Open findings: 1")}); got != nil || !strings.Contains(why, "review G-260101-00007 of candidate abcdef1 does not end with") {
		t.Fatalf("got %v, %q", got, why)
	}
	if got, _ := s.review(context.Background(), r, []*project.Record{rev("G-260101-00005", ClosingLine), rev("G-260101-00007", ClosingLine)}); got == nil || got.ID != "G-260101-00007" {
		t.Fatalf("got %v, want the newest", got)
	}
}

// Integrating the first of two candidates moves the target; the second is
// predicted and verified again against the moved target before its
// approval, and integrated too.
func TestSweepIntegratesSeveralCandidatesInOneSweep(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t, strings.Replace(policy, "%s", "'true'", 1), map[string]string{"code.txt": "the change\n"}, ClosingLine)
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-260101-00002")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-260101-00002", wt)
	write(t, wt, "grove/G-260101-00002-second.md", strings.Replace(strings.Replace(work, "G-260101-00001", "G-260101-00002", 1), "%s", "active", 1))
	write(t, wt, "other.txt", "more\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: another")
	examined := git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "grove/G-260101-00006-review.md", strings.NewReplacer("G-260101-00005", "G-260101-00006", "G-260101-00001", "G-260101-00002").Replace(strings.Replace(strings.Replace(review, "%s", examined, 1), "%s", ClosingLine, 1)))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "docs: review")
	if _, err := update.Apply(wt, update.Request{ID: "G-260101-00002", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: git(t, wt, "rev-parse", "HEAD")}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	_, facts := sweep(t, root)
	joined := strings.Join(facts, "\n")
	if strings.Count(joined, ": done: ") != 2 || strings.Contains(joined, "waits") {
		t.Fatalf("facts:\n%s", joined)
	}
}

// handedOff launches one attempt of G-260101-00001, active on main with
// config, through TestMain's fake provider, which hands off a candidate
// writing files with a review that closes clean, and waits until the
// attempt's owner has exited, its sweep included. setup runs in main's
// checkout before the launch. It returns main's checkout, the branch's, and
// the attempt.
func handedOff(t *testing.T, config string, files map[string]string, setup func(root string)) (root, wt string, v *attempt.View) {
	t.Helper()
	if testing.Short() {
		t.Skip("waits on a fake provider's attempt and its owner's sweep")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root = filepath.Join(dir, "repo")
	git(t, dir, "init", "-q", "-b", "main", root)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	write(t, root, ".gitignore", ".claude/worktrees/\n")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-260101-00001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, root, attempt.SkillPath, "---\nname: grove-work\n---\n")
	write(t, root, "code.txt", "base\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	for name, content := range files {
		write(t, dir, filepath.Join("files", name), content)
	}
	write(t, dir, "review.md", strings.Replace(strings.Replace(review, "%s", "EXAMINED", 1), "%s", ClosingLine, 1))
	if setup != nil {
		setup(root)
	}
	l, err := attempt.Start(attempt.Request{Root: root, IDs: []string{"G-260101-00001"}, BudgetUSD: "1", PermissionMode: "acceptEdits"}, now, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	return root, l.Worktree, ended(t, root, l.Attempt)
}

// ended waits until the attempt's owner has exited and reads it.
func ended(t *testing.T, root, id string) *attempt.View {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		v, err := attempt.Show(root, id)
		if err != nil {
			t.Fatal(err)
		}
		// The owner has released the attempt's lock when it sweeps: its exit
		// is the end, and what it wrote is read after it.
		if v.Status == attempt.Finished && syscall.Kill(v.Launch.Owner, 0) != nil {
			if v, err = attempt.Show(root, id); err != nil {
				t.Fatal(err)
			}
			return v
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %s did not end: %+v", id, v)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestAnAttemptsOwnerIntegratesItsCandidateInsideThePolicy(t *testing.T) {
	t.Parallel()
	// The owner's environment must not reach a verify command, or a grove
	// or test binary there would take itself for the attempt's owner.
	config := strings.Replace(policy, "%s", "grep -q change code.txt, 'test -z \"$GROVE_ATTEMPT_OWNER\"'", 1)
	root, _, v := handedOff(t, config, map[string]string{"code.txt": "the change\n"}, nil)
	joined := strings.Join(v.Sweep, "\n")
	for _, want := range []string{"sweep of G-260101-00001 in " + root + " under policy grove.yaml sha256:", "G-260101-00001: verified: ", "G-260101-00001: approved under policy", "G-260101-00001: done: "} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in sweep.log:\n%s", want, joined)
		}
	}
	if r := record(t, root); r.Status != "done" || !update.Delegated(r) || !strings.Contains(string(r.Source), "attempt "+v.Launch.Attempt+" produced it") {
		t.Fatalf("record on main:\n%s", r.Source)
	}
	if facts := strings.Join(attempt.Facts(v, func(s string) string { return s }), "\n"); !strings.Contains(facts, "Sweep: ") {
		t.Fatalf("facts:\n%s", facts)
	}
}

func TestAnAttemptsOwnerLeavesWhatThePolicyDoesNotName(t *testing.T) {
	t.Parallel()
	config := strings.Replace(policy, "%s", "'true'", 1)
	t.Run("out of policy", func(t *testing.T) {
		t.Parallel()
		root, wt, v := handedOff(t, config, map[string]string{"secret/key": "x\n"}, nil)
		if len(v.Sweep) != 2 || !strings.Contains(v.Sweep[1], "G-260101-00001: waits: changes secret/key, which never matches") {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
		if r := record(t, wt); r.Status != "review" || r.Approved != "" || record(t, root).Status != "active" {
			t.Fatalf("something was written:\n%s", r.Source)
		}
	})
	t.Run("a conflict", func(t *testing.T) {
		t.Parallel()
		root, wt, v := handedOff(t, config, map[string]string{"code.txt": "branch\n"}, func(root string) {
			git(t, root, "branch", "worktree-G-260101-00001")
			write(t, root, "code.txt", "main\n")
			git(t, root, "commit", "-qam", "main moves")
		})
		last := v.Sweep[len(v.Sweep)-1]
		if !strings.Contains(last, "G-260101-00001: resolution attempt ") {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
		// The resolution attempt hands the same conflict back; its own
		// owner's sweep waits instead of trying again.
		views, err := attempt.List(root, "G-260101-00001")
		if err != nil || len(views) != 2 {
			t.Fatalf("attempts %v, %v", views, err)
		}
		again := ended(t, root, views[0].Launch.Attempt)
		if again.Launch.Attempt == v.Launch.Attempt {
			again = ended(t, root, views[1].Launch.Attempt)
		}
		if len(again.Sweep) != 2 || !strings.Contains(again.Sweep[1], "waits: ") || !strings.Contains(again.Sweep[1], "again after a resolution of main at ") {
			t.Fatalf("the resolution's sweep.log %q", again.Sweep)
		}
		if r := record(t, wt); r.Status != "review" || r.Approved != "" {
			t.Fatalf("record:\n%s", r.Source)
		}
	})
	t.Run("no policy", func(t *testing.T) {
		t.Parallel()
		root, wt, v := handedOff(t, "schema_version: 4\nrecords: grove\ntarget: main\n", map[string]string{"code.txt": "x\n"}, nil)
		if v.Sweep != nil || record(t, wt).Status != "review" || record(t, root).Status != "active" {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
		if _, err := os.Stat(filepath.Join(v.Dir, attempt.SweepLog)); !os.IsNotExist(err) {
			t.Fatalf("sweep.log written: %v", err)
		}
	})
	t.Run("a dirty target", func(t *testing.T) {
		t.Parallel()
		root, wt, v := handedOff(t, config, map[string]string{"code.txt": "the change\n"}, func(root string) {
			write(t, root, "code.txt", "edited\n")
		})
		if len(v.Sweep) != 1 || !strings.Contains(v.Sweep[0], "not swept: the checkout of main at "+root+" has uncommitted changes; nothing was written") {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
		if r := record(t, wt); r.Status != "review" || r.Approved != "" {
			t.Fatalf("record:\n%s", r.Source)
		}
	})
	t.Run("while another sweep runs", func(t *testing.T) {
		t.Parallel()
		// The owner waits, saying so, until the sweep holding the lock ends.
		root, _, v := handedOff(t, config, map[string]string{"code.txt": "the change\n"}, func(root string) {
			unlock, err := lock(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			go func() {
				defer unlock()
				for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
					views, _ := attempt.List(root, "G-260101-00001")
					if len(views) == 1 && views[0].Status == attempt.Finished && len(views[0].Sweep) != 0 {
						return
					}
				}
			}()
		})
		if len(v.Sweep) < 2 || !strings.Contains(v.Sweep[0], "waiting for another sweep of this repository to end") || !strings.Contains(v.Sweep[len(v.Sweep)-1], "G-260101-00001: done: ") {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
		if r := record(t, root); r.Status != "done" {
			t.Fatalf("record on main:\n%s", r.Source)
		}
	})
	t.Run("no target checkout", func(t *testing.T) {
		t.Parallel()
		_, _, v := handedOff(t, config, map[string]string{"code.txt": "the change\n"}, func(root string) {
			git(t, root, "checkout", "-q", "-b", "elsewhere")
		})
		if len(v.Sweep) != 1 || !strings.Contains(v.Sweep[0], "not swept: no checkout is on the target main") {
			t.Fatalf("sweep.log %q", v.Sweep)
		}
	})
}

// One sweep at a time: a second refuses, having done nothing.
func TestSweepRefusesWhileAnotherRuns(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t, strings.Replace(policy, "%s", "'true'", 1), map[string]string{"code.txt": "x\n"}, ClosingLine)
	unlock, err := lock(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Run(now, func(f string) { t.Fatalf("acted: %s", f) }); err == nil || !strings.Contains(err.Error(), "another sweep of this repository is running") {
		t.Fatalf("got %v", err)
	}
}
