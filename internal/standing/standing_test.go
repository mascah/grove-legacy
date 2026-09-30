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

// accepted is id's source accepted for candidate.
func accepted(id, candidate string) string {
	r, _ := project.ParseRecord("grove/"+id+".md", []byte(work(id, "review", "candidate: \""+candidate+"\"\n")))
	return work(id, "accepted", "candidate: \""+candidate+"\"\napproved: \""+candidate+"\"\napproved_by: owner\napproved_context: \""+project.AcceptanceContext(r)+"\"\n")
}

// accept writes id accepted for candidate on the current branch and commits it.
func accept(t *testing.T, root, id, candidate string) {
	t.Helper()
	write(t, root, "grove/"+id+".md", accepted(id, candidate))
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

func inspect(t *testing.T, root string) map[string]*Standing {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	return Inspect(context.Background(), root, p.Target, p.Records)
}

func audit(t *testing.T, root string) map[string]*Proof {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	proofs, err := Audit(context.Background(), root, p.Target, p.RecordDir, p.Records)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]*Proof{}
	for _, p := range proofs {
		out[p.ID] = p
	}
	return out
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

func mustRecords(t *testing.T, root string) []*project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	return p.Records
}

// TestStandingReadsTheTargetsRecord covers the completion table: accepted
// before delivery, a squash delivery read from the target and from the kept
// branch, cleanup and gc, a stale acceptance, a Next edit, feedback after
// delivery, acceptance again on the kept branch, an ordinary merge, and
// schema 3's claim.
func TestStandingReadsTheTargetsRecord(t *testing.T) {
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
	if s.State != Accepted || s.Tip != git(t, root, "rev-parse", "main") || s.Text() != "accepted, awaiting delivery to main" {
		t.Fatalf("accepted before delivery: %+v", s)
	}
	sub := git(t, root, "rev-parse", "work")
	git(t, root, "checkout", "-q", "main")
	squash(t, root, sub, "Grove-Work: G-260101-00001\nGrove-Candidate: "+c+"\nGrove-Submitted: "+sub)
	git(t, root, "reset", "-q", "--hard", "main")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done || s.Text() != "done: delivered to main" {
		t.Fatalf("squash delivery: %+v", s)
	}
	// The kept branch's own copy reads the same.
	git(t, root, "checkout", "-q", "work")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("from the kept branch: %+v", s)
	}
	git(t, root, "checkout", "-q", "main")
	// The work branch is gone, nothing retains it, and gc runs: reading
	// needs none of it.
	git(t, root, "branch", "-D", "work")
	git(t, root, "reflog", "expire", "--expire=now", "--all")
	git(t, root, "gc", "-q", "--prune=now")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("after cleanup and gc: %+v", s)
	}
	// An edit of the requirements after delivery withdraws the acceptance.
	src := readFile(t, root, "grove/G-260101-00001.md")
	write(t, root, "grove/G-260101-00001.md", strings.Replace(src, "It works.", "It works fast.", 1))
	if s = inspect(t, root)["G-260101-00001"]; s.State != Review || !strings.Contains(s.Text(), "no longer applies") {
		t.Fatalf("stale acceptance: %+v", s)
	}
	// A Next edit does not.
	write(t, root, "grove/G-260101-00001.md", src+"\n## Next\n\nDelivered.\n")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("a Next edit withdrew the acceptance: %+v", s)
	}
	// Feedback after delivery reopens it for current work.
	write(t, root, "grove/G-260101-00001.md", work("G-260101-00001", "active", "candidate: \""+c+"\"\n"))
	if s = inspect(t, root)["G-260101-00001"]; s.State != Active {
		t.Fatalf("reopened: %+v", s)
	}
	git(t, root, "checkout", "-q", "--", ".")
	// Accepted again for another candidate on a branch: a new delivery.
	git(t, root, "checkout", "-q", "-b", "again")
	write(t, root, "code.txt", "newer\n")
	git(t, root, "commit", "-qam", "fix: code")
	c2 := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c2)
	if s = inspect(t, root)["G-260101-00001"]; s.State != Accepted {
		t.Fatalf("accepted again: %+v", s)
	}
	// An ordinary merge brings the acceptance with it.
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "-q", "--no-ff", "-m", "merge", "again")
	if s = inspect(t, root)["G-260101-00001"]; s.State != Done {
		t.Fatalf("ordinary merge: %+v", s)
	}
}

// TestStandingUnknowns: without a target, or with one that is not here,
// accepted work is unknown, never done or awaiting delivery.
func TestStandingUnknowns(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	accept(t, root, "G-260101-00001", git(t, root, "rev-parse", "HEAD"))
	st := Inspect(context.Background(), root, "", mustRecords(t, root))
	if s := st["G-260101-00001"]; s.State != Unknown || !strings.Contains(s.Why, "no target") {
		t.Fatalf("no target: %+v", s)
	}
	st = Inspect(context.Background(), root, "gone", mustRecords(t, root))
	if s := st["G-260101-00001"]; s.State != Unknown || s.Text() != "accepted; delivery unknown: the target branch gone cannot be read here" {
		t.Fatalf("no target branch: %+v", s)
	}
}

// processes counts the Git processes one Each starts over versions copies
// of every record, as a board passes each checkout's copy, and requires
// every accepted record done and none of them a walk of history.
func processes(t *testing.T, root string, versions int) int {
	t.Helper()
	var records []*project.Record
	for range versions {
		records = append(records, mustRecords(t, root)...)
	}
	trace := filepath.Join(t.TempDir(), "trace")
	t.Setenv("GIT_TRACE", trace)
	res := Each(context.Background(), root, "main", records)
	t.Setenv("GIT_TRACE", "0")
	for r, s := range res {
		if r.Status == "accepted" && s.State != Done {
			t.Fatalf("%s: %+v", r.ID, s)
		}
	}
	log, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	for _, walk := range []string{" log ", " rev-list ", " merge-tree ", " merge-base ", " for-each-ref "} {
		if strings.Contains(string(log), walk) {
			t.Errorf("a reading ran git%s:\n%s", walk, log)
		}
	}
	return strings.Count(string(log), "trace: built-in: git ")
}

// TestStandingGitWorkIsConstant: one reading starts one Git process whatever
// the number of deliveries, accepted records and versions of each record
// passed in, and walks no history. Not parallel: it sets GIT_TRACE.
func TestStandingGitWorkIsConstant(t *testing.T) {
	root := fixture(t)
	delivered(t, root, "G-260101-00010")
	if one := processes(t, root, 1); one != 1 {
		t.Errorf("one delivery started %d Git processes", one)
	}
	delivered(t, root, "G-260101-00011")
	delivered(t, root, "G-260101-00012")
	if many := processes(t, root, 3); many != 1 {
		t.Errorf("three deliveries in three versions started %d Git processes", many)
	}
}

// TestJudgeStartsNothing: the board's reading, from records it holds.
func TestJudgeStartsNothing(t *testing.T) {
	t.Parallel()
	here, _ := project.ParseRecord("grove/G-260101-00001.md", []byte(accepted("G-260101-00001", strings.Repeat("c", 40))))
	target, _ := project.ParseRecord("grove/G-260101-00001.md", []byte(accepted("G-260101-00001", strings.Repeat("c", 40))))
	other, _ := project.ParseRecord("grove/G-260101-00001.md", []byte(accepted("G-260101-00001", strings.Repeat("d", 40))))
	for name, tc := range map[string]struct {
		copies map[string]*project.Record
		want   string
	}{
		"held":              {map[string]*project.Record{here.Path: target}, Done},
		"another":           {map[string]*project.Record{here.Path: other}, Accepted},
		"absent":            {map[string]*project.Record{}, Accepted},
		"target unreadable": {nil, Unknown},
	} {
		if s := Judge("main", "t", tc.copies, []*project.Record{here})[here]; s.State != tc.want {
			t.Errorf("%s: %+v", name, s)
		}
	}
}

func readFile(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestAuditProvesDeliveries: a squash delivery and an ordinary merge are
// proved; an acceptance written on the target by hand reads done, which is
// the trade a cheap reading makes, and only the audit says it is not proved.
func TestAuditProvesDeliveries(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	delivered(t, root, "G-260101-00010")
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	accept(t, root, "G-260101-00001", c)
	git(t, root, "checkout", "-q", "-b", "side")
	write(t, root, "side.txt", "side\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: side")
	side := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "-q", "--no-ff", "-m", "merge", "work")
	accept(t, root, "G-260101-00003", side)
	if s := inspect(t, root)["G-260101-00003"]; s.State != Done {
		t.Fatalf("a hand-written acceptance: %+v", s)
	}
	proofs := audit(t, root)
	if p := proofs["G-260101-00010"]; !p.Proved || p.Submitted == "" || !strings.HasPrefix(p.Text(), "proved: squashed as ") {
		t.Errorf("squash: %+v", p)
	}
	if p := proofs["G-260101-00001"]; !p.Proved || p.Delivered != c {
		t.Errorf("ordinary merge: %+v", p)
	}
	if p := proofs["G-260101-00003"]; p.Proved || p.Unauditable || !strings.Contains(p.Why, "no delivery on the target names candidate "+side[:7]) {
		t.Errorf("a hand-written acceptance: %+v", p)
	}
	accept(t, root, "G-260101-00003", strings.Repeat("a", 40))
	if p := audit(t, root)["G-260101-00003"]; p.Proved || !p.Unauditable || !strings.Contains(p.Why, "is not in this repository") {
		t.Errorf("a candidate that is not here: %+v", p)
	}
}

// TestAuditRejectsForgedDeliveries covers claims the trailers make and Git
// does not support: a tree that is not the merge, code after the candidate,
// a record not accepted there, and a submitted tip that is not here.
func TestAuditRejectsForgedDeliveries(t *testing.T) {
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
	record := git(t, root, "show", good+":grove/G-260101-00001.md") + "\n"
	for name, tc := range map[string]struct{ tree, sub, want string }{
		"a tree that is not the merge": {git(t, root, "rev-parse", "main^{tree}"), good, "is not what merging"},
		"code after the candidate":     {strings.Fields(git(t, root, "merge-tree", "--write-tree", base, sneaky))[0], sneaky, "changes sneaky.txt after the candidate"},
		"not accepted there":           {git(t, root, "rev-parse", c+"^{tree}"), c, "is not accepted for candidate"},
	} {
		d := git(t, root, "commit-tree", tc.tree, "-p", base, "-m", "feat: forged\n\n"+trailers(tc.sub))
		git(t, root, "update-ref", "refs/heads/main", d)
		git(t, root, "reset", "-q", "--hard", "main")
		// The target's record claims the acceptance, as a copied file would.
		write(t, root, "grove/G-260101-00001.md", record)
		git(t, root, "commit", "-q", "--allow-empty", "-am", "claim")
		if p := audit(t, root)["G-260101-00001"]; p.Proved || p.Unauditable || !strings.Contains(p.Why, tc.want) {
			t.Errorf("%s: %+v", name, p)
		}
		git(t, root, "update-ref", "refs/heads/main", base)
		git(t, root, "reset", "-q", "--hard", "main")
	}
	missing := strings.Repeat("d", 40)
	d := git(t, root, "commit-tree", git(t, root, "rev-parse", "main^{tree}"), "-p", base, "-m", "feat: lost\n\n"+trailers(missing))
	git(t, root, "update-ref", "refs/heads/main", d)
	git(t, root, "reset", "-q", "--hard", "main")
	write(t, root, "grove/G-260101-00001.md", record)
	git(t, root, "commit", "-q", "--allow-empty", "-am", "claim")
	if p := audit(t, root)["G-260101-00001"]; p.Proved || !p.Unauditable || !strings.Contains(p.Why, "fetch refs/grove/*") {
		t.Errorf("missing submitted tip: %+v", p)
	}
}

// TestAuditSurvivesUnrelatedClaims: a claim naming as its submitted tip a
// commit unrelated to its parent, which git merge-tree refuses by default,
// fails its own record's audit and leaves every other's intact.
func TestAuditSurvivesUnrelatedClaims(t *testing.T) {
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
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "claims")
	proofs := audit(t, root)
	if p := proofs["G-260101-00001"]; p.Proved || !strings.Contains(p.Why, "does not verify: delivery "+d[:7]+" names submitted tip "+orphan[:7]+", which does not contain the candidate") {
		t.Errorf("a claim on a tip without the candidate: %+v", p)
	}
	if p := proofs["G-260101-00003"]; p.Proved || !strings.Contains(p.Why, "does not verify: delivery "+d3[:7]+" is not what merging") {
		t.Errorf("a claim on a history of its own: %+v", p)
	}
	if p := proofs["G-260101-00010"]; !p.Proved {
		t.Errorf("another record: %+v", p)
	}
}

// TestAuditTransport covers a group whose member is later reopened, which
// leaves its sibling proved, and what a clone can audit: without the
// retained evidence its reading is unchanged and the audit says it cannot
// prove the delivery here, until it fetches refs/grove/*; a shallow clone
// cannot audit.
func TestAuditTransport(t *testing.T) {
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
	write(t, root, "grove/G-260101-00001.md", work("G-260101-00001", "review", "candidate: \""+d+"\"\n"))
	git(t, root, "commit", "-qam", "reopen first")
	if s := inspect(t, root)["G-260101-00003"]; s.State != Done {
		t.Fatalf("the sibling of a reopened member: %+v", s)
	}
	if p := audit(t, root)["G-260101-00003"]; !p.Proved || p.Delivered != d {
		t.Fatalf("the sibling's audit: %+v", p)
	}
	clone := filepath.Join(t.TempDir(), "clone")
	git(t, root, "clone", "-q", "--no-local", root, clone)
	if s := inspect(t, clone)["G-260101-00003"]; s.State != Done {
		t.Fatalf("a clone without the evidence reads: %+v", s)
	}
	if p := audit(t, clone)["G-260101-00003"]; p.Proved || !p.Unauditable || !strings.Contains(p.Why, "refs/grove") {
		t.Fatalf("a clone without the evidence audits: %+v", p)
	}
	git(t, clone, "fetch", "-q", "origin", "refs/grove/*:refs/grove/*")
	if p := audit(t, clone)["G-260101-00003"]; !p.Proved || p.Delivered != d {
		t.Fatalf("after fetching refs/grove/*: %+v", p)
	}
	// A shallow clone lacks the candidate, so it cannot audit, never fails.
	shallow := filepath.Join(t.TempDir(), "shallow")
	git(t, root, "clone", "-q", "--depth", "1", "file://"+root, shallow)
	if p := audit(t, shallow)["G-260101-00003"]; p.Proved || !p.Unauditable {
		t.Fatalf("shallow: %+v", p)
	}
}

// TestAuditReadsTheTargetsRecords: a record added after the candidate, as a
// sibling's or a review's, is allowed because the target holds it, even
// audited from a checkout that lacks it.
func TestAuditReadsTheTargetsRecords(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	git(t, root, "checkout", "-q", "-b", "work")
	write(t, root, "code.txt", "new\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: code")
	c := git(t, root, "rev-parse", "HEAD")
	write(t, root, "grove/G-260101-00009.md", work("G-260101-00009", "proposed", ""))
	accept(t, root, "G-260101-00001", c)
	sub := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	squash(t, root, sub, "Grove-Work: G-260101-00001\nGrove-Candidate: "+c+"\nGrove-Submitted: "+sub)
	git(t, root, "reset", "-q", "--hard", "main")
	if err := os.Remove(filepath.Join(root, "grove/G-260101-00009.md")); err != nil {
		t.Fatal(err)
	}
	if p := audit(t, root)["G-260101-00001"]; p == nil || !p.Proved {
		t.Fatalf("%+v", p)
	}
}
