package create

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

const config = "schema_version: 3\nrecords: grove\n"

func record(id, kind, status string) string {
	return "---\nid: \"" + id + "\"\ntype: " + kind + "\ntitle: T\nstatus: " + status + "\n---\nBody.\n"
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	out, err := repo.Command(context.Background(), dir, full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, path, source string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}

// gitProject returns a committed Git project containing G-001.
func gitProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-001-first.md", record("G-001", "work", "done"))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	return root
}

func load(t *testing.T, root string) *project.Project {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics: %v", ds)
	}
	return p
}

func stateDir(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(git(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), "grove")
}

// day is the date part of every ID issued at now.
var now = time.Date(2026, 9, 25, 16, 0, 0, 0, time.UTC)

const day = "G-260925-"

// draws replaces the random tail with tails, repeated, for one test that
// does not run in parallel, and returns how many were drawn.
func draws(t *testing.T, tails ...string) func() int {
	t.Helper()
	old := tail
	var mu sync.Mutex
	n := 0
	tail = func() string {
		mu.Lock()
		defer mu.Unlock()
		n++
		return tails[(n-1)%len(tails)]
	}
	t.Cleanup(func() { tail = old })
	return func() int { mu.Lock(); defer mu.Unlock(); return n }
}

func TestTailIsFiveCrockfordCharacters(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for range 200 {
		id := day + tail()
		if !project.IDPattern.MatchString(id) {
			t.Fatalf("%s is not a valid ID", id)
		}
		seen[id] = true
	}
	if len(seen) < 199 {
		t.Fatalf("200 draws gave only %d distinct tails", len(seen))
	}
}

// G-195 acceptance 1: a tail already held by a record on a local ref, in
// another worktree or in the loaded project is drawn again.
func TestNewSkipsIDsInRefsAndWorktrees(t *testing.T) {
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	// 00001 exists only in the feature branch's history, so the ref scan
	// alone can find it; 00002 exists only as a live file in the worktree and
	// 00003 only in this checkout. Both others are nested: the scan follows
	// discovery, which is recursive.
	write(t, wt, "grove/deep/committed.md", record(day+"00001", "work", "proposed"))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-q", "-m", "branch record")
	if err := os.Remove(filepath.Join(wt, "grove/deep/committed.md")); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "grove/any/where/live.md", record(day+"00002", "work", "proposed"))
	write(t, root, "grove/here.md", record(day+"00003", "work", "proposed"))
	// Another day's ID with the same tail is no obstacle.
	write(t, root, "grove/yesterday.md", record("G-260924-00004", "work", "proposed"))
	drawn := draws(t, "00001", "00002", "00003", "00004")
	path, err := New(load(t, root), "work", "Next", "next", now)
	if err != nil || path != "grove/"+day+"00004-next.md" || drawn() != 4 {
		t.Fatalf("got %q, %v after %d draws; want %s00004 after 4", path, err, drawn(), day)
	}
}

func TestIssueBoundsItsDraws(t *testing.T) {
	root := gitProject(t)
	write(t, root, "grove/taken.md", record(day+"77777", "work", "proposed"))
	drawn := draws(t, "77777")
	_, err := New(load(t, root), "work", "Never", "never", now)
	if err == nil || !strings.Contains(err.Error(), "every one of 8 random IDs for G-260925 was already in use") || drawn() != tries {
		t.Fatalf("got %v after %d draws; want a refusal after %d", err, drawn(), tries)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 2 {
		t.Fatalf("nothing may be created: %v", entries)
	}
}

// G-195 acceptance 3: new reads and writes no counter or allocator lock; the
// write lock is the only state it creates.
func TestNewCreatesOnlyTheWriteLock(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	if _, err := New(load(t, root), "work", "One", "", now); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(stateDir(t, root))
	if err != nil || len(entries) != 1 || entries[0].Name() != "write.lock" {
		t.Fatalf("state: %v, %v; want only write.lock", entries, err)
	}
}

func TestNewRefusesWhenAWorktreeCannotBeScanned(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	hidden := filepath.Join(wt, "grove/G-030-unreadable.md")
	write(t, wt, "grove/G-030-unreadable.md", record("G-030", "work", "proposed"))
	if err := os.Chmod(hidden, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(hidden, 0o644) })
	if path, err := New(load(t, root), "work", "Blind", "", now); err == nil {
		t.Fatalf("created %s although a worktree record could not be read", path)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("nothing may be created: %v", entries)
	}
}

func TestLockFailureCreatesNoFile(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	root := gitProject(t)
	dir := stateDir(t, root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil { // the write lock cannot be created
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	if _, err := New(load(t, root), "work", "Blocked", "blocked", now); err == nil || !strings.Contains(err.Error(), "nothing created") {
		t.Fatalf("expected a lock diagnostic, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("a failed new must create nothing: %v", entries)
	}
}

func TestNewRequiresGit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-001-first.md", record("G-001", "work", "done"))
	if _, err := New(load(t, root), "work", "T", "", now); err == nil || !strings.Contains(err.Error(), "Git") {
		t.Fatalf("expected a Git requirement error, got %v", err)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 2 {
		t.Fatalf("no state may be created outside Git: %v", entries)
	}
}

func TestNewCreatesValidRecords(t *testing.T) {
	root := gitProject(t)
	draws(t, "7k2qm", "0000z", "zzzzz")
	path, err := New(load(t, root), "work", `Title: with "quotes" & more`, "", now)
	if err != nil || path != "grove/G-260925-7k2qm-title-with-quotes-more.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	p := load(t, root)
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == "G-260925-7k2qm" })
	if i < 0 {
		t.Fatal("G-260925-7k2qm not readable")
	}
	r := p.Records[i]
	if r.Title != `Title: with "quotes" & more` || r.Status != "proposed" || r.Created == nil || !r.Created.Equal(now) || !r.Updated.Equal(now) {
		t.Fatalf("unexpected record %+v", r)
	}
	if !strings.Contains(string(r.Source), "## Outcome") {
		t.Fatalf("missing body skeleton: %s", r.Source)
	}
	if path, err := New(p, "question", "Which?", "which", now); err != nil || path != "grove/G-260925-0000z-which.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	// The ID's date is UTC: late on the 25th west of Greenwich is the 26th.
	late := time.Date(2026, 9, 25, 20, 0, 0, 0, time.FixedZone("UTC-5", -5*3600))
	if path, err := New(p, "decision", "Choose", "", late); err != nil || path != "grove/G-260926-zzzzz-choose.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	load(t, root)
}

func TestNewNeverOverwrites(t *testing.T) {
	root := gitProject(t)
	// A directory at the target name makes O_EXCL creation fail while the
	// reader (which skips directories) still sees a valid project.
	if err := os.MkdirAll(filepath.Join(root, "grove/G-260925-00000-x.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	draws(t, "00000")
	if _, err := New(load(t, root), "work", "X", "x", now); err == nil {
		t.Fatal("expected creation to fail on an existing path")
	}
}

func TestNewRejectsBadSlugAndKind(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	for _, c := range [][2]string{{"work", "Bad Slug"}, {"work", "UPPER"}, {"release", "ok"}} {
		if _, err := New(p, c[0], "T", c[1], now); err == nil {
			t.Fatalf("expected %v to be rejected", c)
		}
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("rejected input must not create files: %v", entries)
	}
}

func TestSlug(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"  Hello,   World!! 1234567890123456789012345678 ":    "hello-world-123456789012",
		"Ünïcode ünd Emoji 🎉":                                 "n-code-nd-emoji",
		"!!!":                                                 "record",
		"trailing-hyphen-at-limit-24-x":                       "trailing-hyphen-at-limit",
		"Identify records by creation date and a random tail": "identify-records-by-crea",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNewRefusesWhenRecordRootChangedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	write(t, root, "other/G-001-first.md", record("G-001", "work", "done"))
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: other\n")
	_, err := New(p, "work", "Moved", "moved", now)
	if err == nil || !strings.Contains(err.Error(), "nothing created: the record root changed") {
		t.Fatalf("expected a configuration-change refusal, got %v", err)
	}
	for _, dir := range []string{"grove", "other"} {
		if entries, _ := os.ReadDir(filepath.Join(root, dir)); len(entries) != 1 {
			t.Fatalf("%s: nothing may be created: %v", dir, entries)
		}
	}
}

// A configuration edit that keeps the parsed meaning is still an observed
// change between the loaded input and publication.
func TestNewRefusesWhenConfigurationBytesChangedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	write(t, root, "grove.yaml", "# concurrent edit\n"+config)
	_, err := New(p, "work", "Late", "late", now)
	if err == nil || !strings.Contains(err.Error(), "nothing created: grove.yaml changed") {
		t.Fatalf("expected a configuration-change refusal, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("nothing may be created: %v", entries)
	}
}

// W-008: a live ID in a checkout whose path Git would quote for display is
// still seen, at the root and at a nested prefix, and the write lock belongs
// under the real common directory.
func TestNewScansOddlyNamedWorktrees(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	for _, prefix := range []string{"", "sub\nproject"} {
		parent, err := filepath.EvalSymlinks(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(parent, "new\nline")
		write(t, filepath.Join(root, prefix), "grove.yaml", config)
		write(t, filepath.Join(root, prefix), "grove/G-001-first.md", record("G-001", "work", "done"))
		git(t, root, "init", "-q", "-b", "main")
		git(t, root, "add", "-A")
		git(t, root, "commit", "-q", "-m", "init")
		odd := filepath.Join(parent, "odd \"quoted\"\n\twt ")
		git(t, root, "worktree", "add", "-q", "-b", "odd", odd)
		write(t, filepath.Join(odd, prefix), "grove/live-only.md", record(day+"00000", "work", "proposed"))
		write(t, odd, "elsewhere/grove/unrelated.md", record(day+"00001", "work", "proposed"))
		// A checkout without the record folder is normal, not a scan error.
		absent := filepath.Join(parent, "absent\nwt")
		git(t, root, "worktree", "add", "-q", "--detach", absent)
		if err := os.RemoveAll(filepath.Join(absent, prefix, "grove")); err != nil {
			t.Fatal(err)
		}

		draws(t, "00000", "00001")
		path, err := New(load(t, filepath.Join(root, prefix)), "work", "Next", "next", now)
		if err != nil || path != "grove/"+day+"00001-next.md" {
			t.Fatalf("prefix %q: got %q, %v; want %s00001 past the live %s00000", prefix, path, err, day, day)
		}
		if _, err := os.Stat(filepath.Join(root, ".git", "grove", "write.lock")); err != nil {
			t.Fatalf("the write lock belongs under the real common directory: %v", err)
		}
		if entries, _ := os.ReadDir(parent); len(entries) != 3 {
			t.Fatalf("nothing may appear beside the checkouts: %q", entries)
		}
	}
}

// Another process can define the term between this caller's load and its
// turn at the write lock; the second record must never reach the disk.
func TestNewRefusesATermDefinedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	stale := load(t, root)
	first, err := New(load(t, root), "term", "Attempt", "", now)
	if err != nil {
		t.Fatal(err)
	}
	_, err = New(stale, "term", "attempt", "", now)
	if err == nil || !strings.Contains(err.Error(), "nothing created: the term attempt is already defined by "+first[len("grove/"):len("grove/G-260925-00000")]) {
		t.Fatalf("expected a refusal under the write lock, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 2 {
		t.Fatalf("nothing may be written for the refused attempt: %v", entries)
	}
	if len(load(t, root).Records) == 0 {
		t.Fatal("the project must stay loadable")
	}
}

// Twelve concurrent new in two worktrees, from a source that draws every
// tail twice: the write lock serializes them and each sees the others' files.
func TestNewConcurrentAcrossWorktrees(t *testing.T) {
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-odd\n\twt ") // W-008: a path Git quotes for display
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	var tails []string
	for i := range 12 {
		tail := fmt.Sprintf("%05d", i)
		tails = append(tails, tail, tail)
	}
	draws(t, tails...)
	projects := []*project.Project{load(t, root), load(t, wt)}
	var wg sync.WaitGroup
	for i := range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := New(projects[i%2], "page", "P", "", now); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	seen := map[string]bool{}
	for _, dir := range []string{root, wt} {
		for _, r := range load(t, dir).Records {
			if r.ID == "G-001" {
				continue // committed baseline, held by both worktrees
			}
			if seen[r.ID] {
				t.Fatalf("%s issued twice", r.ID)
			}
			seen[r.ID] = true
		}
	}
	if len(seen) != 12 {
		t.Fatalf("expected 12 distinct IDs, got %v", seen)
	}
}
