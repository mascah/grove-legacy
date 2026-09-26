package renumber

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/repo"
)

func git(t *testing.T, dir string, env []string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	cmd := repo.Command(context.Background(), dir, full...)
	cmd.Env = append(cmd.Env, env...)
	out, err := cmd.CombinedOutput()
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

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

const (
	first = "---\nid: \"G-001\"\ntype: work\ntitle: \"First, with a title longer than the slug cap\"\nstatus: active\nrelates_to: [\"G-002\"]\ncreated: \"2026-09-20T10:00:00Z\"\nupdated: \"2026-09-21T11:00:00Z\"\n---\n\n" +
		"See [second](G-002-converted.md) and branch `worktree-G-001-G-002`; not G-1000, AG-001 or G-0012.\n"
	second = "---\nid: \"G-002\"\ntype: work\ntitle: Converted\nstatus: proposed\nformerly: \"docs/old.md\"\n---\n\nBack to [first](./G-001-first.md#outcome).\n"
	plan   = "---\nid: \"G-003\"\ntype: plan\ntitle: Plan for G-001\nstatus: current\nwork: [\"G-001\"]\ncreated: \"2026-09-20T12:00:00Z\"\nupdated: \"2026-09-20T12:00:00Z\"\n---\n\nFor G-001, reviewed in [a review](G-260921-abcde-review-of-g-001.md).\n"
	review = "---\nid: \"G-260921-abcde\"\ntype: review\ntitle: \"Review of G-001's work\"\nstatus: current\nwork: [\"G-001\"]\n---\n"
)

func TestRenumber(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, nil, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\ntarget: main\n")
	write(t, root, "docs/old.md", "# Old\n")
	git(t, root, nil, "add", "-A")
	git(t, root, []string{"GIT_AUTHOR_DATE=2026-09-14T09:30:00Z"}, "commit", "-q", "-m", "old document")
	git(t, root, nil, "rm", "-q", "docs/old.md")
	write(t, root, "grove/G-001-first.md", first)
	write(t, root, "grove/G-002-converted.md", second)
	write(t, root, "grove/G-003-plan.md", plan)
	write(t, root, "grove/G-260921-abcde-review-of-g-001.md", review)
	git(t, root, nil, "add", "-A")
	git(t, root, nil, "commit", "-q", "-m", "records")
	attempts := filepath.Join(git(t, root, nil, "rev-parse", "--path-format=absolute", "--git-common-dir"), "grove", "attempts")
	write(t, attempts, "G-001.20260920T100000Z/attempt.json", `{"attempt":"G-001.20260920T100000Z","work":"G-001","record_path":"grove/G-001-first.md"}`)
	write(t, attempts, "G-001.20260920T100000Z/result.json", `{"exit_code":0}`)
	write(t, attempts, "G-001.20260920T100000Z/events.jsonl", `{"text":"G-001 in grove/G-001-first.md\nG-002\tG-003"}`+"\n")

	git(t, root, nil, "branch", "stale")
	if _, err := Run(root); err == nil || !strings.Contains(err.Error(), "hold records: stale") {
		t.Fatalf("a second branch with records: %v", err)
	}
	if status := git(t, root, nil, "status", "--porcelain"); status != "" {
		t.Fatalf("a refusal wrote:\n%s", status)
	}
	git(t, root, nil, "branch", "-D", "stale")

	renames, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(renames) != 4 {
		t.Fatalf("map = %+v", renames)
	}
	one, two, three, four := renames[0], renames[1], renames[2], renames[3]
	for i, want := range []struct{ from, fromPath, day, slug string }{
		{"G-001", "grove/G-001-first.md", "G-260920-", "-first-with-a-title-longe.md"},
		{"G-002", "grove/G-002-converted.md", "G-260914-", "-converted.md"},                                   // formerly's first commit
		{"G-003", "grove/G-003-plan.md", "G-260920-", "-plan-for.md"},                                         // a slug leaves out the IDs a title cites
		{"G-260921-abcde", "grove/G-260921-abcde-review-of-g-001.md", "G-260921-abcde", "-review-of-work.md"}, // a date-form record named after a legacy ID
	} {
		m := renames[i]
		if m.From != want.from || m.FromPath != want.fromPath || !strings.HasPrefix(m.ID, want.day) || m.Path != "grove/"+m.ID+want.slug {
			t.Fatalf("map line %d = %+v", i, m)
		}
		if _, err := os.Stat(filepath.Join(root, want.fromPath)); !os.IsNotExist(err) {
			t.Fatalf("%s still exists", want.fromPath)
		}
	}
	gotFirst := read(t, filepath.Join(root, one.Path))
	wantFirst := strings.NewReplacer("G-002-converted.md", filepath.Base(two.Path), "worktree-G-001-G-002", "worktree-"+one.ID+"-"+two.ID,
		`"G-001"`, `"`+one.ID+`"`, `"G-002"`, `"`+two.ID+`"`).Replace(first)
	if gotFirst != wantFirst {
		t.Fatalf("first:\n%s\nwant:\n%s", gotFirst, wantFirst)
	}
	if got := read(t, filepath.Join(root, two.Path)); got != "---\nid: \""+two.ID+"\"\ntype: work\ntitle: Converted\nstatus: proposed\nformerly: \"docs/old.md\"\ncreated: \"2026-09-14T09:30:00Z\"\n---\n\nBack to [first](./"+filepath.Base(one.Path)+"#outcome).\n" {
		t.Fatalf("second:\n%s", got)
	}
	if got := read(t, filepath.Join(root, three.Path)); !strings.Contains(got, `work: ["`+one.ID+`"]`) || !strings.Contains(got, "For "+one.ID+", reviewed in [a review]("+filepath.Base(four.Path)+").") {
		t.Fatalf("plan:\n%s", got)
	}

	views, err := attempt.List(root, one.ID)
	if err != nil || len(views) != 1 || filepath.Base(views[0].Dir) != one.ID+".20260920T100000Z" || views[0].Launch.Work != one.ID {
		t.Fatalf("attempts = %+v, %v", views, err)
	}
	if got := read(t, filepath.Join(views[0].Dir, "events.jsonl")); got != `{"text":"`+one.ID+` in `+one.Path+`\n`+two.ID+`\t`+three.ID+`"}`+"\n" {
		t.Fatalf("events = %s", got)
	}

	if _, err := Run(root); err == nil || !strings.Contains(err.Error(), "no record has a legacy ID") {
		t.Fatalf("rerun: %v", err)
	}
}
