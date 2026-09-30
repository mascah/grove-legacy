package standing

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	out, err := repo.Command(context.Background(), dir, full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, name, source string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}

const body = "\n## Outcome\n\nIt works.\n"

func work(id, status, extra string) string {
	return "---\nid: \"" + id + "\"\ntype: work\ntitle: " + id + "\nstatus: " + status + "\n" + extra + "---\n" + body
}

// accept writes id accepted for candidate on the current branch and commits it.
func accept(t *testing.T, root, id, candidate string) {
	t.Helper()
	r, _ := project.ParseRecord("grove/"+id+".md", []byte(work(id, "review", "candidate: \""+candidate+"\"\n")))
	write(t, root, "grove/"+id+".md", work(id, "accepted", "candidate: \""+candidate+"\"\napproved: \""+candidate+"\"\napproved_by: owner\napproved_context: \""+project.AcceptanceContext(r)+"\"\n"))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "accept "+id)
}

// squash delivers submitted onto main as delivery does, with message's trailers.
func squash(t *testing.T, root, submitted, trailers string) string {
	t.Helper()
	base := git(t, root, "rev-parse", "main")
	tree := strings.Fields(git(t, root, "merge-tree", "--write-tree", base, submitted))[0]
	d := git(t, root, "commit-tree", tree, "-p", base, "-m", "feat: deliver\n\n"+trailers)
	git(t, root, "update-ref", "refs/heads/main", d, base)
	return d
}

// delivered accepts id on a branch of its own, squash-delivers it onto main
// and retains the submitted tip, as integrate does.
func delivered(t *testing.T, root, id string) {
	t.Helper()
	git(t, root, "checkout", "-q", "-b", id, "main")
	write(t, root, id+".txt", id+"\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: "+id)
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, id, c)
	sub := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	squash(t, root, sub, "Grove-Work: "+id+"\nGrove-Candidate: "+c+"\nGrove-Submitted: "+sub)
	git(t, root, "reset", "-q", "--hard", "main")
	git(t, root, "update-ref", Ref(sub), sub)
	git(t, root, "branch", "-D", id)
}

// processes counts the Git processes one Each starts over versions copies
// of every record, as a board passes each checkout's copy, and requires
// every accepted record done.
func processes(t *testing.T, root string, versions int) int {
	t.Helper()
	var records []*project.Record
	for range versions {
		records = append(records, mustRecords(t, root)...)
	}
	trace := filepath.Join(t.TempDir(), "trace")
	t.Setenv("GIT_TRACE", trace)
	res, err := Each(context.Background(), root, "main", records)
	t.Setenv("GIT_TRACE", "0")
	if err != nil {
		t.Fatal(err)
	}
	for r, s := range res {
		if r.Status == "accepted" && s.State != Done {
			t.Fatalf("%s: %+v", r.ID, s)
		}
	}
	log, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(log), "trace: built-in: git ")
}

func inspect(t *testing.T, root string) map[string]*Standing {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	res, err := Inspect(context.Background(), root, p.Target, p.Records)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func fixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\ntarget: main\n")
	write(t, root, "grove/G-260101-00001.md", work("G-260101-00001", "active", ""))
	write(t, root, "grove/G-260101-00002.md", work("G-260101-00002", "done", ""))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	return root
}

// TestStandingDerivesDoneFromVerifiedDelivery covers the completion table:
// accepted before delivery, ordinary and squash delivery, a crash after the
// target advanced (nothing but Git to read), feedback after delivery, a
// stale acceptance, and schema 3's claim.
func TestStandingDerivesDoneFromVerifiedDelivery(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	if s := inspect(t, root)["G-260101-00002"]; s.State != Done || !s.Legacy || s.Text() != "done (schema 3 claim)" {
		t.Fatalf("legacy done: %+v", s)
	}
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	s := inspect(t, root)["G-260101-00001"]
	if s.State != Accepted || s.Tip == "" || s.Text() != "accepted, awaiting delivery to main" {
		t.Fatalf("accepted before delivery: %+v", s)
	}
	// The target moves on without it: still awaiting delivery.
	git(t, root, "checkout", "-q", "main")
	write(t, root, "other.txt", "other\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "chore: other")
	sub := git(t, root, "rev-parse", "work")
	d := squash(t, root, sub, "Grove-Work: G-260101-00001\nGrove-Candidate: "+c+"\nGrove-Submitted: "+sub)
	git(t, root, "reset", "-q", "--hard", "main")
	s = inspect(t, root)["G-260101-00001"]
	if s.State != Done || s.Delivered != d || s.Submitted != sub || s.Text() != "done: squashed as "+d[:7]+" on main" {
		t.Fatalf("squash delivery: %+v", s)
	}
	// The work branch is gone and gc runs: the retained ref keeps the evidence.
	git(t, root, "update-ref", Ref(sub), sub)
	git(t, root, "branch", "-D", "work")
	git(t, root, "reflog", "expire", "--expire=now", "--all")
	git(t, root, "gc", "-q", "--prune=now")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("after cleanup and gc: %+v", s)
	}
	// An edit of the requirements after delivery withdraws the acceptance.
	src, _ := os.ReadFile(filepath.Join(root, "grove/G-260101-00001.md"))
	write(t, root, "grove/G-260101-00001.md", strings.Replace(string(src), "It works.", "It works fast.", 1))
	if s = inspect(t, root)["G-260101-00001"]; s.State != Review || !strings.Contains(s.Text(), "no longer applies") {
		t.Fatalf("stale acceptance: %+v", s)
	}
	// A Next edit does not.
	write(t, root, "grove/G-260101-00001.md", string(src)+"\n## Next\n\nDelivered.\n")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("a Next edit withdrew the acceptance: %+v", s)
	}
	// Feedback after delivery reopens it for current work.
	write(t, root, "grove/G-260101-00001.md", work("G-260101-00001", "active", "candidate: \""+c+"\"\n"))
	if s = inspect(t, root)["G-260101-00001"]; s.State != Active {
		t.Fatalf("reopened: %+v", s)
	}
}

// TestStandingOrdinaryMergeAndUnknowns covers ancestry delivery, a missing
// candidate, no target, and a shallow clone that cannot assert absence.
func TestStandingOrdinaryMergeAndUnknowns(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	git(t, root, "checkout", "-q", "main")
	// A shallow clone of main alone cannot see a delivery that is not there.
	clone := filepath.Join(t.TempDir(), "clone")
	git(t, root, "clone", "-q", "--depth", "2", "--branch", "work", "file://"+root, clone)
	git(t, clone, "fetch", "-q", "--depth", "1", "origin", "main:main")
	p, _ := project.Load(clone, clone)
	res, err := Inspect(context.Background(), clone, "main", p.Records)
	if err != nil {
		t.Fatal(err)
	}
	if s := res["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "shallow") {
		t.Fatalf("shallow: %+v", s)
	}
	git(t, root, "merge", "-q", "--no-ff", "-m", "merge", "work")
	if s := inspect(t, root)["G-260101-00001"]; s.State != Done || s.Delivered != c {
		t.Fatalf("ordinary merge: %+v", s)
	}
	res, _ = Inspect(context.Background(), root, "", mustRecords(t, root))
	if s := res["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "no target") {
		t.Fatalf("no target: %+v", s)
	}
	// A candidate this repository lacks is unknown, not absent.
	write(t, root, "grove/G-260101-00001.md", strings.ReplaceAll(readFile(t, root, "grove/G-260101-00001.md"), c, strings.Repeat("e", 40)))
	if s := inspect(t, root)["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "is not in this repository") {
		t.Fatalf("missing candidate: %+v", s)
	}
}

func mustRecords(t *testing.T, root string) []*project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	return p.Records
}

func readFile(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestStandingRejectsForgedDeliveries covers claims the trailers make and Git
// does not support: a tree that is not the merge, code after the candidate,
// a record not accepted there, and a submitted tip that is not here. Each is
// unknown, never done.
func TestStandingRejectsForgedDeliveries(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	good := git(t, root, "rev-parse", "HEAD")
	write(t, root, "sneaky.txt", "unreviewed\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "chore: sneak")
	sneaky := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	trailers := func(sub string) string {
		return "Grove-Work: G-260101-00001\nGrove-Candidate: " + c + "\nGrove-Submitted: " + sub
	}
	base := git(t, root, "rev-parse", "main")
	for name, tc := range map[string]struct{ tree, sub, want string }{
		"a tree that is not the merge": {git(t, root, "rev-parse", "main^{tree}"), good, "is not what merging"},
		"code after the candidate":     {strings.Fields(git(t, root, "merge-tree", "--write-tree", base, sneaky))[0], sneaky, "changes sneaky.txt after the candidate"},
		"not accepted there":           {git(t, root, "rev-parse", c+"^{tree}"), c, "is not accepted for candidate"},
	} {
		d := git(t, root, "commit-tree", tc.tree, "-p", base, "-m", "feat: forged\n\n"+trailers(tc.sub))
		git(t, root, "update-ref", "refs/heads/main", d)
		git(t, root, "reset", "-q", "--hard", "main")
		// The live record claims the acceptance, as a copied file would.
		write(t, root, "grove/G-260101-00001.md", git(t, root, "show", good+":grove/G-260101-00001.md")+"\n")
		// Not done, and not proof of none either: unknown, so nothing
		// delivers it again before someone reconciles the claim.
		if s := inspect(t, root)["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, tc.want) {
			t.Errorf("%s: %+v", name, s)
		}
		git(t, root, "update-ref", "refs/heads/main", base)
		git(t, root, "reset", "-q", "--hard", "main")
	}
	missing := strings.Repeat("d", 40)
	d := git(t, root, "commit-tree", git(t, root, "rev-parse", "main^{tree}"), "-p", base, "-m", "feat: lost\n\n"+trailers(missing))
	git(t, root, "update-ref", "refs/heads/main", d)
	write(t, root, "grove/G-260101-00001.md", git(t, root, "show", good+":grove/G-260101-00001.md")+"\n")
	if s := inspect(t, root)["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "fetch refs/grove/*") {
		t.Errorf("missing submitted tip: %+v", s)
	}
}

// TestStandingSurvivesUnrelatedClaims: a claim naming as its submitted tip a
// commit unrelated to its parent, which git merge-tree refuses by default,
// leaves its record unknown and every other record's reading intact, both
// where that tip lacks the candidate and where it holds one of its own.
func TestStandingSurvivesUnrelatedClaims(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	delivered(t, root, "G-260101-00010")
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	good := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "--orphan", "alone", "main")
	write(t, root, "alone.txt", "alone\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: alone")
	c3 := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00003", c3)
	s3 := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	orphan := git(t, root, "commit-tree", "main^{tree}", "-m", "orphan")
	d := git(t, root, "commit-tree", "main^{tree}", "-p", "main", "-m", "feat: forged\n\nGrove-Candidate: "+c+"\nGrove-Submitted: "+orphan)
	d3 := git(t, root, "commit-tree", "main^{tree}", "-p", d, "-m", "feat: forged\n\nGrove-Candidate: "+c3+"\nGrove-Submitted: "+s3)
	git(t, root, "update-ref", "refs/heads/main", d3)
	git(t, root, "reset", "-q", "--hard", "main")
	write(t, root, "grove/G-260101-00001.md", git(t, root, "show", good+":grove/G-260101-00001.md")+"\n")
	write(t, root, "grove/G-260101-00003.md", git(t, root, "show", s3+":grove/G-260101-00003.md")+"\n")
	st := inspect(t, root)
	if s := st["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "does not verify: delivery "+d[:7]+" names submitted tip "+orphan[:7]+", which does not contain the candidate") {
		t.Errorf("a claim on a tip without the candidate: %+v", s)
	}
	if s := st["G-260101-00003"]; s.State != Unknown || !strings.Contains(s.Why, "does not verify: delivery "+d3[:7]+" is not what merging") {
		t.Errorf("a claim on a history of its own: %+v", s)
	}
	if s := st["G-260101-00010"]; s.State != Done {
		t.Errorf("another record: %+v", s)
	}
}

// TestStandingCrissCrossIsGitsOwnMerge: with two best merge bases, as when
// an unrelated history is merged into both the target and a branch kept
// after its delivery, the earlier submission is never the base, so a
// delivery that brings back what the target removed since does not verify.
// The unrelated root is dated first, so that git merge-base without --all
// names the fork point, the base that would bring it back.
func TestStandingCrissCrossIsGitsOwnMerge(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "one\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: one")
	c1 := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c1)
	s1 := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	squash(t, root, s1, "Grove-Candidate: "+c1+"\nGrove-Submitted: "+s1)
	git(t, root, "reset", "-q", "--hard", "main")
	git(t, root, "update-ref", Ref(s1), s1)
	git(t, root, "checkout", "-q", "--orphan", "other")
	git(t, root, "rm", "-rqf", ".")
	write(t, root, "r.txt", "r\n")
	git(t, root, "add", "-A")
	old := repo.Command(context.Background(), root, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "commit", "-qm", "chore: other root")
	old.Env = append(old.Env, "GIT_AUTHOR_DATE=2000-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2000-01-01T00:00:00Z")
	if out, err := old.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	for _, b := range []string{"main", "work"} {
		git(t, root, "checkout", "-q", b)
		git(t, root, "merge", "-q", "--allow-unrelated-histories", "-m", "merge other", "other")
	}
	git(t, root, "checkout", "-q", "main")
	git(t, root, "rm", "-q", "r.txt")
	git(t, root, "commit", "-qm", "chore: remove r")
	tip := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "work")
	write(t, root, "code.txt", "two\n")
	git(t, root, "commit", "-qam", "feat: two")
	c2 := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c2)
	s2 := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	if bases := strings.Fields(git(t, root, "merge-base", "--all", tip, s2)); len(bases) != 2 {
		t.Fatalf("fixture: merge bases %v", bases)
	}
	if base, err := Base(context.Background(), root, tip, s2); err != nil || base != "" {
		t.Fatalf("a criss-cross continued from %q: %v", base, err)
	}
	// Merged from the earlier submission, r.txt comes back.
	tree := strings.Fields(git(t, root, "merge-tree", "--write-tree", "--merge-base="+s1, tip, s2))[0]
	d := git(t, root, "commit-tree", tree, "-p", tip, "-m", "feat: deliver\n\nGrove-Candidate: "+c2+"\nGrove-Submitted: "+s2)
	git(t, root, "update-ref", "refs/heads/main", d)
	git(t, root, "reset", "-q", "--hard", "main")
	if s := inspect(t, root)["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "is not what merging") {
		t.Fatalf("a delivery bringing back r.txt: %+v", s)
	}
}

// TestStandingGitWorkIsConstant: one reading starts the same Git processes
// whatever the number of deliveries, open acceptances and versions of each
// record passed in. Not parallel: it sets GIT_TRACE.
func TestStandingGitWorkIsConstant(t *testing.T) {
	root := fixture(t)
	delivered(t, root, "G-260101-00010")
	one := processes(t, root, 1)
	delivered(t, root, "G-260101-00011")
	delivered(t, root, "G-260101-00012")
	if many := processes(t, root, 3); many != one {
		t.Errorf("one delivery started %d Git processes, three deliveries in three versions %d", one, many)
	}
}

// TestStandingGroupAndTransport covers a group whose member is later
// reopened and given another candidate, which leaves its sibling done, and
// the transport: a clone without the retained evidence reads unknown until
// it fetches refs/grove/*.
func TestStandingGroupAndTransport(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	write(t, root, "grove/G-260101-00003.md", work("G-260101-00003", "active", ""))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "second")
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "both\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: both")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	accept(t, root, "G-260101-00003", c)
	sub := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	d := squash(t, root, sub, "Grove-Work: G-260101-00001\nGrove-Work: G-260101-00003\nGrove-Candidate: "+c+"\nGrove-Submitted: "+sub)
	git(t, root, "reset", "-q", "--hard", "main")
	git(t, root, "update-ref", Ref(sub), sub)
	git(t, root, "branch", "-D", "work")
	// The first member is reopened with another candidate.
	write(t, root, "grove/G-260101-00001.md", work("G-260101-00001", "review", "candidate: \""+d+"\"\n"))
	git(t, root, "commit", "-qam", "reopen first")
	st := inspect(t, root)
	if s := st["G-260101-00003"]; s.State != Done || s.Delivered != d {
		t.Fatalf("the sibling of a reopened member: %+v", s)
	}
	clone := filepath.Join(t.TempDir(), "clone")
	git(t, root, "clone", "-q", "--no-local", root, clone)
	if s := inspect(t, clone)["G-260101-00003"]; s.State != Unknown || !strings.Contains(s.Why, "refs/grove") {
		t.Fatalf("a clone without the evidence: %+v", s)
	}
	git(t, clone, "fetch", "-q", "origin", "refs/grove/*:refs/grove/*")
	if s := inspect(t, clone)["G-260101-00003"]; s.State != Done || s.Delivered != d {
		t.Fatalf("after fetching refs/grove/*: %+v", s)
	}
}
