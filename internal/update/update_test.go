package update

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

const (
	work     = "---\nid: \"G-260101-00001\"\ntype: work\ntitle: First\nstatus: proposed\nrelates_to: [\"G-260101-00003\"]\ncreated: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-19T12:00:00Z\"\n---\n\n## Outcome\n\nBody --- stays.\n"
	second   = "---\nid: \"G-260101-00002\"\ntype: work\ntitle: Second\nstatus: proposed\n---\nBody.\n"
	question = "---\nid: \"G-260101-00003\"\ntype: question\ntitle: Which?\nstatus: open\n---\nBody.\n"
	decision = "---\nid: \"G-260101-00004\"\ntype: decision\ntitle: Choose\nstatus: proposed\n---\nBody.\n"
)

var now = time.Date(2026, 9, 19, 18, 30, 0, 500, time.UTC)

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

func read(t *testing.T, root, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// gitProject returns a committed Git project with G-260101-00001, G-260101-00002, G-260101-00003, G-260101-00004.
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
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n")
	write(t, root, "grove/work/G-260101-00001-first.md", work)
	write(t, root, "grove/work/G-260101-00002-second.md", second)
	write(t, root, "grove/questions/G-260101-00003-which.md", question)
	write(t, root, "grove/decisions/G-260101-00004-choose.md", decision)
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	return root
}

func revision(t *testing.T, root, path string) string {
	t.Helper()
	return project.Revision([]byte(read(t, root, path)))
}

func record(t *testing.T, root, id string) *project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("project invalid: %v", ds)
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == id })
	if i < 0 {
		t.Fatalf("%s missing", id)
	}
	return p.Records[i]
}

func apply(t *testing.T, root, id string, sets []Field, unsets ...string) Result {
	t.Helper()
	r := record(t, root, id)
	res, err := Apply(root, Request{ID: id, Expect: project.Revision(r.Source), Set: sets, Unset: unsets}, now, nil)
	if err != nil {
		t.Fatalf("update %s: %v", id, err)
	}
	if res.ID != id || res.Path != r.Path || res.Revision != revision(t, root, r.Path) {
		t.Fatalf("result %+v does not describe the file", res)
	}
	return res
}

func tempFiles(t *testing.T, root string) []string {
	t.Helper()
	var found []string
	filepath.WalkDir(filepath.Join(root, "grove"), func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && filepath.Ext(path) != ".md" {
			found = append(found, path)
		}
		return nil
	})
	return found
}

func TestUpdateEveryFieldOnItsTypes(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	res := apply(t, root, "G-260101-00001", []Field{
		{"title", ` "Quoted": ünïcode — 版本 #1 `}, {"status", "active"}, {"kind", "tooling"}, {"priority", "1"},
		{"size", "large"}, {"members", `["G-260101-00002"]`}, {"depends_on", `["G-260101-00002"]`}, {"relates_to", `["G-260101-00003", "G-260101-00004"]`},
	})
	if !res.Changed {
		t.Fatal("expected a change")
	}
	r := record(t, root, "G-260101-00001")
	if r.Title != ` "Quoted": ünïcode — 版本 #1 ` || r.Status != "active" || r.Kind != "tooling" || *r.Priority != 1 || r.Size != "large" ||
		!slices.Equal(r.Members, []string{"G-260101-00002"}) || !slices.Equal(r.DependsOn, []string{"G-260101-00002"}) || !slices.Equal(r.RelatesTo, []string{"G-260101-00003", "G-260101-00004"}) {
		t.Fatalf("fields not applied: %+v", r)
	}
	if r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || r.Updated.Format(time.RFC3339) != "2026-09-19T18:30:00Z" {
		t.Fatalf("dates: created %v updated %v", r.Created, r.Updated)
	}
	source := string(r.Source)
	if !strings.HasSuffix(source, "---\n\n## Outcome\n\nBody --- stays.\n") || !strings.HasPrefix(source, "---\nid: \"G-260101-00001\"\ntype: work\ntitle: \" \\\"Quoted\\\": ünïcode — 版本 #1 \"\nstatus: active\nrelates_to: [\"G-260101-00003\", \"G-260101-00004\"]\ncreated: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-19T18:30:00Z\"\nkind: tooling\npriority: 1\nsize: large\nmembers: [\"G-260101-00002\"]\ndepends_on: [\"G-260101-00002\"]\n---") {
		t.Fatalf("unexpected source:\n%s", source)
	}
	apply(t, root, "G-260101-00003", []Field{{"blocks", `["G-260101-00001"]`}, {"status", "resolved"}, {"title", "Resolved?"}, {"relates_to", "[]"}})
	q := record(t, root, "G-260101-00003")
	if !slices.Equal(q.Blocks, []string{"G-260101-00001"}) || q.Status != "resolved" || q.RelatesTo == nil || len(q.RelatesTo) != 0 || q.Created != nil || q.Updated == nil {
		t.Fatalf("question: %+v", q)
	}
	apply(t, root, "G-260101-00004", []Field{{"status", "accepted"}, {"relates_to", `["G-260101-00001"]`}})
	if d := record(t, root, "G-260101-00004"); d.Status != "accepted" || d.Created != nil {
		t.Fatalf("decision: %+v", d)
	}
	// Optional removal and explicit empty lists.
	apply(t, root, "G-260101-00001", []Field{{"members", "[]"}}, "kind", "priority", "size", "depends_on")
	r = record(t, root, "G-260101-00001")
	if r.Kind != "" || r.Priority != nil || r.Size != "" || r.DependsOn != nil || r.Members == nil || len(r.Members) != 0 {
		t.Fatalf("removal: %+v", r)
	}
	if strings.Contains(string(r.Source), "kind") || !strings.Contains(string(r.Source), "members: []\n") {
		t.Fatalf("source after removal:\n%s", r.Source)
	}
	if files := tempFiles(t, root); len(files) != 0 {
		t.Fatalf("temporary files left behind: %v", files)
	}
}

func TestUpdateLifecycleAndReopening(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	// Untouched optional fields, including the pointer-valued priority, must
	// compare equal between the original and the candidate.
	apply(t, root, "G-260101-00001", []Field{{"priority", "2"}, {"size", "small"}, {"members", `["G-260101-00002"]`}})
	head := git(t, root, "rev-parse", "HEAD")
	for _, status := range []string{"active", "review", "proposed", "abandoned", "review"} {
		sets := []Field{{"status", status}}
		if status == "review" { // review needs the candidate, and every status keeps it
			sets = append(sets, Field{"candidate", head})
		}
		apply(t, root, "G-260101-00001", sets)
		r := record(t, root, "G-260101-00001")
		if r.Status != status || *r.Priority != 2 || r.Size != "small" || r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || !strings.HasSuffix(string(r.Source), "\n---\n\n## Outcome\n\nBody --- stays.\n") {
			t.Fatalf("%s: %+v", status, r)
		}
		if status != "active" && (r.Candidate != head || !strings.Contains(string(r.Source), "candidate: \""+head+"\"")) {
			t.Fatalf("%s: candidate must stay and stay quoted: %+v", status, r)
		}
	}
	for _, status := range []string{"resolved", "open"} {
		apply(t, root, "G-260101-00003", []Field{{"status", status}})
	}
	for _, status := range []string{"accepted", "superseded", "rejected", "proposed"} {
		apply(t, root, "G-260101-00004", []Field{{"status", status}})
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove/work")); len(entries) != 2 {
		t.Fatalf("status changes must not rename or add files: %v", entries)
	}
}

func TestUpdateNoOpPreservesBytesAndRefusesStale(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	path := filepath.Join(root, "grove/work/G-260101-00001-first.md")
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	past := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	os.Chtimes(path, past, past)
	before, _ := os.Stat(path)
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	for _, req := range []Request{
		{Set: []Field{{"title", "First"}, {"status", "proposed"}, {"relates_to", `["G-260101-00003"]`}}},
		{Unset: []string{"kind", "priority", "members"}},
		{Set: []Field{{"relates_to", `["G-260101-00003"]`}}, Unset: []string{"size"}},
	} {
		req.ID, req.Expect = "G-260101-00001", expect
		// The clock is behind the record's dates; a no-op needs no clock.
		res, err := Apply(root, req, past, nil)
		if err != nil || res.Changed || res.Revision != expect {
			t.Fatalf("%+v: %+v %v", req, res, err)
		}
	}
	after, _ := os.Stat(path)
	if read(t, root, "grove/work/G-260101-00001-first.md") != work || !after.ModTime().Equal(before.ModTime()) || after.Mode() != before.Mode() {
		t.Fatal("no-op changed bytes, mtime, or permissions")
	}
	_, err := Apply(root, Request{ID: "G-260101-00001", Expect: "sha256:" + strings.Repeat("0", 64), Set: []Field{{"status", "proposed"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "changed since the expected revision") || !strings.Contains(err.Error(), expect) {
		t.Fatalf("a stale no-op must be refused and report the current revision: %v", err)
	}
	// Absent versus explicitly empty are different states.
	if res := apply(t, root, "G-260101-00001", []Field{{"members", "[]"}}); !res.Changed {
		t.Fatal("adding an empty list to an absent field is a change")
	}
	if res := apply(t, root, "G-260101-00001", []Field{{"members", "[]"}}); res.Changed {
		t.Fatal("an explicit empty list already present is a no-op")
	}
	if res := apply(t, root, "G-260101-00001", nil, "members"); !res.Changed {
		t.Fatal("removing a present empty list is a change")
	}
	if res := apply(t, root, "G-260101-00001", []Field{{"relates_to", `["G-260101-00003"]`}, {"members", "[]"}}, "kind"); !res.Changed {
		t.Fatal("one effective change among no-ops still applies")
	}
}

func TestUpdateRejectsInvalidRequestsWithoutWriting(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	for _, tc := range []struct {
		name string
		id   string
		set  []Field
		un   []string
		want string
	}{
		{"unknown field", "G-260101-00001", []Field{{"foo", "x"}}, nil, "foo"},
		{"id", "G-260101-00001", []Field{{"id", "G-260101-00009"}}, nil, "id cannot"},
		{"created", "G-260101-00001", []Field{{"created", "\"2026-01-01T00:00:00Z\""}}, nil, "created cannot"},
		{"updated", "G-260101-00001", nil, []string{"updated"}, "updated cannot"},
		{"unset id", "G-260101-00001", nil, []string{"id"}, "id cannot"},
		{"unset required", "G-260101-00001", nil, []string{"title"}, "required"},
		{"wrong type field", "G-260101-00001", []Field{{"blocks", "[]"}}, nil, "blocks"},
		{"empty title", "G-260101-00001", []Field{{"title", "  "}}, nil, "title"},
		{"status value", "G-260101-00001", []Field{{"status", "resolved"}}, nil, "status"},
		{"status injection", "G-260101-00001", []Field{{"status", "active\nkind: fix"}}, nil, "invalid"},
		{"kind value", "G-260101-00001", []Field{{"kind", "magic"}}, nil, "kind"},
		{"size value", "G-260101-00001", []Field{{"size", "huge"}}, nil, "size"},
		{"priority letters", "G-260101-00001", []Field{{"priority", "high"}}, nil, "priority"},
		{"priority negative", "G-260101-00001", []Field{{"priority", "-1"}}, nil, "priority"},
		{"priority range", "G-260101-00001", []Field{{"priority", "6"}}, nil, "priority"},
		{"list null", "G-260101-00001", []Field{{"members", "null"}}, nil, "members"},
		{"list not json", "G-260101-00001", []Field{{"members", "G-260101-00002"}}, nil, "members"},
		{"list numbers", "G-260101-00001", []Field{{"members", "[1]"}}, nil, "members"},
		{"list duplicate", "G-260101-00001", []Field{{"members", `["G-260101-00002", "G-260101-00002"]`}}, nil, "duplicate"},
		{"self link", "G-260101-00001", []Field{{"depends_on", `["G-260101-00001"]`}}, nil, "self"},
		{"unresolved", "G-260101-00001", []Field{{"depends_on", `["G-260101-00404"]`}}, nil, "unresolved"},
		{"non-work target", "G-260101-00001", []Field{{"members", `["G-260101-00003"]`}}, nil, "must be work"},
		{"missing record", "G-260101-00404", []Field{{"status", "active"}}, nil, "not found"},
		{"invalid utf8", "G-260101-00001", []Field{{"title", "bad\xff"}}, nil, "UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Apply(root, Request{ID: tc.id, Expect: expect, Set: tc.set, Unset: tc.un}, now, nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wanted %q, got %v", tc.want, err)
			}
			if read(t, root, "grove/work/G-260101-00001-first.md") != work {
				t.Fatal("a rejected request changed the file")
			}
		})
	}
	if files := tempFiles(t, root); len(files) != 0 {
		t.Fatalf("temporary files left behind: %v", files)
	}
}

func TestUpdateRejectsCyclesAndInvalidProjects(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	apply(t, root, "G-260101-00001", []Field{{"depends_on", `["G-260101-00002"]`}, {"members", `["G-260101-00002"]`}})
	expect := revision(t, root, "grove/work/G-260101-00002-second.md")
	for _, field := range []string{"depends_on", "members"} {
		_, err := Apply(root, Request{ID: "G-260101-00002", Expect: expect, Set: []Field{{field, `["G-260101-00001"]`}}}, now, nil)
		if err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Fatalf("%s: wanted a cycle refusal, got %v", field, err)
		}
	}
	// A group may depend on its own members; that is not a cycle.
	apply(t, root, "G-260101-00002", []Field{{"relates_to", `["G-260101-00001"]`}})
	write(t, root, "grove/work/broken.md", "---\nid: \"G-260101-00005\"\ntype: work\ntitle: Broken\nstatus: imaginary\n---\n")
	_, err := Apply(root, Request{ID: "G-260101-00001", Expect: revision(t, root, "grove/work/G-260101-00001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "not valid") || !strings.Contains(err.Error(), "broken.md") {
		t.Fatalf("an invalid project must be refused: %v", err)
	}
	if strings.Contains(read(t, root, "grove/work/G-260101-00001-first.md"), "status: active") {
		t.Fatal("refused update was written")
	}
}

func TestUpdateClockContract(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	early := time.Date(2026, 9, 19, 11, 59, 59, 0, time.UTC)
	_, err := Apply(root, Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", "active"}}}, early, nil)
	if err == nil || !strings.Contains(err.Error(), "clock") || read(t, root, "grove/work/G-260101-00001-first.md") != work {
		t.Fatalf("a clock behind the record must be refused without writing: %v", err)
	}
	// Equal to created is acceptable; two changes within one second share a stamp.
	same := time.Date(2026, 9, 19, 12, 0, 0, 999_999_999, time.UTC)
	first, err := Apply(root, Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", "active"}}}, same, nil)
	if err != nil {
		t.Fatal(err)
	}
	secondRes, err := Apply(root, Request{ID: "G-260101-00001", Expect: first.Revision, Set: []Field{{"status", "abandoned"}}}, same, nil)
	if err != nil || secondRes.Revision == first.Revision {
		t.Fatalf("same-second change: %+v %v", secondRes, err)
	}
	r := record(t, root, "G-260101-00001")
	if r.Updated.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" {
		t.Fatalf("stamp must truncate to seconds and preserve created: %+v", r)
	}
	// A record without created stays without one.
	apply(t, root, "G-260101-00002", []Field{{"status", "active"}})
	if r := record(t, root, "G-260101-00002"); r.Created != nil || r.Updated == nil || !strings.HasSuffix(string(r.Source), "status: active\nupdated: \"2026-09-19T18:30:00Z\"\n---\nBody.\n") {
		t.Fatalf("missing created must stay absent: %s", r.Source)
	}
}

func TestUpdateStaleAfterBodyOnlyEdit(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	write(t, root, "grove/work/G-260101-00001-first.md", work+"An appended paragraph with unchanged timestamps.\n")
	_, err := Apply(root, Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "changed since") {
		t.Fatalf("body edits must invalidate the revision: %v", err)
	}
}

func TestUpdateDetectsChangesDuringPreparation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		change func(t *testing.T, root string)
	}{
		{"neighbor body", func(t *testing.T, root string) {
			write(t, root, "grove/work/G-260101-00002-second.md", second+"more\n")
		}},
		{"inventory", func(t *testing.T, root string) {
			write(t, root, "grove/work/G-260101-00005-new.md", strings.Replace(second, "G-260101-00002", "G-260101-00005", 1))
		}},
		{"configuration", func(t *testing.T, root string) {
			os.Rename(filepath.Join(root, "grove"), filepath.Join(root, "records"))
			write(t, root, "grove.yaml", "schema_version: 4\nrecords: records\n")
		}},
		{"configuration comment only", func(t *testing.T, root string) {
			write(t, root, "grove.yaml", "# concurrent edit\nschema_version: 4\nrecords: grove\n")
		}},
		{"target permissions", func(t *testing.T, root string) {
			os.Chmod(filepath.Join(root, "grove/work/G-260101-00001-first.md"), 0o600)
		}},
		{"target replaced", func(t *testing.T, root string) {
			os.Remove(filepath.Join(root, "grove/work/G-260101-00001-first.md"))
			write(t, root, "grove/work/G-260101-00001-first.md", work)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := gitProject(t)
			expect := revision(t, root, "grove/work/G-260101-00001-first.md")
			fault := func(step string) error {
				if step == "compare" {
					tc.change(t, root)
				}
				return nil
			}
			_, err := Apply(root, Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", "active"}}}, now, fault)
			if err == nil || errors.As(err, new(*Failure)) {
				t.Fatalf("expected a refusal before publication, got %v", err)
			}
			for _, dir := range []string{"grove", "records"} {
				if data, err := os.ReadFile(filepath.Join(root, dir, "work/G-260101-00001-first.md")); err == nil && string(data) != work {
					t.Fatal("refused update was written")
				}
			}
			if files := tempFiles(t, root); len(files) != 0 {
				t.Fatalf("temporary files left behind: %v", files)
			}
		})
	}
}

func TestUpdateInjectedFailures(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	path := filepath.Join(root, "grove/work/G-260101-00001-first.md")
	if err := os.Chmod(path, 0o640); err != nil {
		t.Fatal(err)
	}
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	req := Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", "active"}}}
	for _, step := range []string{"write", "sync", "close", "compare", "rename"} {
		_, err := Apply(root, req, now, func(s string) error {
			if s == step {
				return errors.New("injected " + step)
			}
			return nil
		})
		if err == nil || !strings.Contains(err.Error(), "injected "+step) || !strings.Contains(err.Error(), "unchanged") || errors.As(err, new(*Failure)) {
			t.Fatalf("%s: %v", step, err)
		}
		if read(t, root, "grove/work/G-260101-00001-first.md") != work {
			t.Fatalf("%s: original bytes lost", step)
		}
		if files := tempFiles(t, root); len(files) != 0 {
			t.Fatalf("%s: temporary files left behind: %v", step, files)
		}
	}
	for _, step := range []string{"dirsync", "validate"} {
		root := gitProject(t)
		if err := os.Chmod(filepath.Join(root, "grove/work/G-260101-00001-first.md"), 0o640); err != nil {
			t.Fatal(err)
		}
		_, err := Apply(root, req, now, func(s string) error {
			if s == step {
				return errors.New("injected " + step)
			}
			return nil
		})
		var failure *Failure
		if !errors.As(err, &failure) || !strings.Contains(err.Error(), "applied") || failure.Path != "grove/work/G-260101-00001-first.md" {
			t.Fatalf("%s: expected an applied-state failure, got %v", step, err)
		}
		got := read(t, root, "grove/work/G-260101-00001-first.md")
		if !strings.Contains(got, "status: active") || failure.Revision != project.Revision([]byte(got)) {
			t.Fatalf("%s: applied state must describe the published bytes", step)
		}
		if info, _ := os.Stat(filepath.Join(root, "grove/work/G-260101-00001-first.md")); info.Mode().Perm() != 0o640 {
			t.Fatalf("%s: permissions not preserved: %v", step, info.Mode())
		}
		if step == "dirsync" && !strings.Contains(err.Error(), "durability") {
			t.Fatalf("directory sync failure must report uncertain durability: %v", err)
		}
	}
}

func TestUpdateSameRevisionRace(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-260101-00001-first.md")
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, status := range []string{"active", "abandoned"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := Apply(root, Request{ID: "G-260101-00001", Expect: expect, Set: []Field{{"status", status}}}, now, nil)
			if err == nil && !res.Changed {
				err = errors.New("unexpected no-op")
			}
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	var failures int
	for err := range results {
		if err != nil {
			failures++
			if !strings.Contains(err.Error(), "changed since") {
				t.Fatal(err)
			}
		}
	}
	if failures != 1 {
		t.Fatalf("exactly one writer must win, got %d refusals", failures)
	}
	if r := record(t, root, "G-260101-00001"); r.Status != "active" && r.Status != "abandoned" {
		t.Fatalf("winning change lost: %+v", r)
	}
}

func TestUpdateReciprocalDependenciesCannotFormACycle(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	pairs := [][2]string{{"G-260101-00001", "G-260101-00002"}, {"G-260101-00002", "G-260101-00001"}}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i, pair := range pairs {
		wg.Add(1)
		expect := revision(t, root, "grove/work/"+pair[0]+"-"+map[string]string{"G-260101-00001": "first", "G-260101-00002": "second"}[pair[0]]+".md")
		go func() {
			defer wg.Done()
			_, errs[i] = Apply(root, Request{ID: pair[0], Expect: expect, Set: []Field{{"depends_on", `["` + pair[1] + `"]`}}}, now, nil)
		}()
	}
	wg.Wait()
	if (errs[0] == nil) == (errs[1] == nil) {
		t.Fatalf("exactly one must succeed: %v / %v", errs[0], errs[1])
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		t.Fatalf("project invalid after concurrent updates: %v", ds)
	}
}

func TestNewAndUpdateShareTheWriteLockAcrossWorktrees(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	common, _, err := repo.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan string, 2)
	go func() {
		_, err := Apply(wt, Request{ID: "G-260101-00001", Expect: revision(t, wt, "grove/work/G-260101-00001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
		if err != nil {
			t.Error(err)
		}
		done <- "update"
	}()
	go func() {
		p, ds := project.Load(root, root)
		if len(ds) != 0 {
			t.Error(ds)
		}
		if _, err := create.New(p, "work", "Third", "third", now); err != nil {
			t.Error(err)
		}
		done <- "new"
	}()
	select {
	case who := <-done:
		t.Fatalf("%s completed while the write lock was held", who)
	case <-time.After(300 * time.Millisecond):
	}
	unlock()
	for range 2 {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("commands did not proceed after the lock was released")
		}
	}
	if r := record(t, wt, "G-260101-00001"); r.Status != "active" {
		t.Fatal("worktree update lost")
	}
	if p, ds := project.Load(root, root); len(ds) != 0 || len(p.Records) != 5 {
		t.Fatalf("main checkout: %v", ds)
	}
	if _, err := os.Stat(filepath.Join(common, "grove", "write.lock")); err != nil {
		t.Fatal("write lock file must remain")
	}
}

func TestUpdateRequiresGit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n")
	write(t, root, "grove/work/G-260101-00001-first.md", work)
	write(t, root, "grove/questions/G-260101-00003-which.md", question)
	_, err := Apply(root, Request{ID: "G-260101-00001", Expect: revision(t, root, "grove/work/G-260101-00001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "Git") || read(t, root, "grove/work/G-260101-00001-first.md") != work {
		t.Fatalf("expected a Git requirement without writes: %v", err)
	}
}

func TestUnchangedGuardCatchesEditorDrift(t *testing.T) {
	t.Parallel()
	before, _ := project.ParseRecord("w.md", []byte(work))
	after, _ := project.ParseRecord("w.md", []byte(strings.Replace(work, "title: First", "title: Other", 1)))
	if err := unchanged(before, after, []change{set("status", "active")}); err == nil || !strings.Contains(err.Error(), "title") {
		t.Fatalf("an untouched field that differs must be refused: %v", err)
	}
	if err := unchanged(before, after, []change{set("title", `"Other"`)}); err != nil {
		t.Fatal(err)
	}
	priority := 2
	before.Priority, after.Priority = &priority, new(int)
	*after.Priority = 2
	after.Title = before.Title
	if err := unchanged(before, after, nil); err != nil {
		t.Fatalf("equal priorities behind different pointers must compare equal: %v", err)
	}
}

// G-260919-z9w13: the review's update reproducers, at the level a user reaches them.
func TestUpdatePreservesAcceptedForms(t *testing.T) {
	t.Parallel()
	const stamped = "updated: \"2026-09-19T18:30:00Z\""
	t.Run("comment before a later-line value", func(t *testing.T) {
		t.Parallel()
		root := gitProject(t)
		write(t, root, "grove/work/G-260101-00001-first.md", strings.Replace(work, "title: First", "title: # retain\n  First", 1))
		apply(t, root, "G-260101-00001", []Field{{"title", "New"}})
		want := strings.NewReplacer("title: First", "title: # retain\n  \"New\"", "updated: \"2026-09-19T12:00:00Z\"", stamped).Replace(work)
		if got := read(t, root, "grove/work/G-260101-00001-first.md"); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})
	t.Run("final two flow entries", func(t *testing.T) {
		t.Parallel()
		for _, unsets := range [][]string{{"kind", "size"}, {"size", "kind"}} {
			root := gitProject(t)
			write(t, root, "grove/work/G-260101-00001-first.md", "---\n{id: G-260101-00001, type: work, title: T, status: proposed, kind: fix, size: small}\n---\nBody\n")
			apply(t, root, "G-260101-00001", nil, unsets...)
			want := "---\n{id: G-260101-00001, type: work, title: T, status: proposed, " + stamped + "}\n---\nBody\n"
			if got := read(t, root, "grove/work/G-260101-00001-first.md"); got != want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, want)
			}
		}
	})
	t.Run("explicit keys", func(t *testing.T) {
		t.Parallel()
		root := gitProject(t)
		source := strings.NewReplacer("status: proposed", "? status\n: proposed", "relates_to: [\"G-260101-00003\"]", "? relates_to\n: [\"G-260101-00003\"] # why").Replace(work)
		write(t, root, "grove/work/G-260101-00001-first.md", source)
		apply(t, root, "G-260101-00001", []Field{{"status", "active"}}, "relates_to")
		want := strings.NewReplacer("status: proposed", "? status\n: active", "relates_to: [\"G-260101-00003\"]\n", "", "updated: \"2026-09-19T12:00:00Z\"", stamped).Replace(work)
		if got := read(t, root, "grove/work/G-260101-00001-first.md"); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})
}

// G-260919-7qv4x: a main checkout named "new\nline" once put the write lock in a
// sibling ".../new/grove". Coordination state belongs under the real common
// directory, from the main and a linked checkout alike.
func TestCoordinationStateStaysUnderTheCommonDirectory(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt := filepath.Join(parent, "new\nline"), filepath.Join(parent, "linked\twt ")
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n")
	write(t, root, "grove/work/G-260101-00001-first.md", work)
	write(t, root, "grove/questions/G-260101-00003-which.md", question)
	git(t, root, "init", "-q", "-b", "main")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)

	apply(t, root, "G-260101-00001", []Field{{"status", "active"}})
	apply(t, wt, "G-260101-00001", []Field{{"status", "abandoned"}})
	for _, dir := range []string{root, wt} {
		p, ds := project.Load(dir, dir)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		if _, err := create.New(p, "work", "Odd", "odd", now); err != nil {
			t.Fatal(err)
		}
	}
	var state []string
	filepath.WalkDir(parent, func(path string, entry os.DirEntry, err error) error {
		if err == nil && slices.Contains([]string{"write.lock", "lock", "neutral-ids"}, entry.Name()) {
			state = append(state, path)
		}
		return nil
	})
	slices.Sort(state)
	common := filepath.Join(root, ".git", "grove")
	if want := []string{filepath.Join(common, "write.lock")}; !slices.Equal(state, want) {
		t.Fatalf("coordination state: %q\nwant %q", state, want)
	}
	if entries, _ := os.ReadDir(parent); len(entries) != 2 {
		t.Fatalf("nothing may appear beside the checkouts: %q", entries)
	}
}

// TestUpdateNeverWritesDone covers schema 4's completion contract
// (G-260930-2qa4a): Done is derived, so no update writes done or changes a
// done record's candidate, while a schema 3 done record migration kept stays
// editable and can be reopened.
func TestUpdateNeverWritesDone(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	head := git(t, root, "rev-parse", "HEAD")
	refuse := func(id, want string, sets ...Field) {
		t.Helper()
		r := record(t, root, id)
		_, err := Apply(root, Request{ID: id, Expect: project.Revision(r.Source), Set: sets}, now, nil)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("%s %v: got %v, want %q", id, sets, err, want)
		}
		if after := record(t, root, id); !bytes.Equal(after.Source, r.Source) {
			t.Fatalf("%s was written despite the refusal", id)
		}
	}
	// A candidate value of letters only must still be written as a quoted string.
	apply(t, root, "G-260101-00002", []Field{{"candidate", "abcdefa"}})
	if src := read(t, root, "grove/work/G-260101-00002-second.md"); !strings.Contains(src, "candidate: \"abcdefa\"\n") {
		t.Fatalf("candidate must be quoted:\n%s", src)
	}
	refuse("G-260101-00001", "candidate: required while status is review: set candidate=COMMIT", Field{"status", "review"})
	refuse("G-260101-00001", "done is derived in schema 4", Field{"status", "done"}, Field{"candidate", head})
	write(t, root, "grove/work/G-260101-00005-done.md", "---\nid: \"G-260101-00005\"\ntype: work\ntitle: Done\nstatus: done\ncandidate: \""+head+"\"\n---\nBody.\n")
	refuse("G-260101-00005", "done is derived in schema 4", Field{"candidate", "abcdef0"})
	apply(t, root, "G-260101-00005", []Field{{"title", "Renamed"}})
	apply(t, root, "G-260101-00005", []Field{{"status", "active"}})
	refuse("G-260101-00005", "done is derived in schema 4", Field{"status", "done"})
}

// TestUpdateOptionalExpectAndCommit covers G-260922-q3cr9: an omitted Expect applies to
// the file as it is while a stale one is still refused, and Commit commits the
// record's file alone, nothing for a no-op, and reports a failed commit as an
// applied update.
func TestUpdateOptionalExpectAndCommit(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1]) // the product's commit uses the repository's own identity
	}
	head := func() string { return git(t, root, "rev-parse", "HEAD") }
	base := head()
	// A direct edit after the last commit: the omitted form applies to it, a stale revision is refused.
	write(t, root, "grove/work/G-260101-00001-first.md", strings.Replace(work, "Body --- stays.", "Edited body.", 1))
	if _, err := Apply(root, Request{ID: "G-260101-00001", Expect: project.Revision([]byte(work)), Set: []Field{{"status", "active"}}}, now, nil); err == nil || !strings.Contains(err.Error(), "changed since the expected revision") {
		t.Fatalf("stale expect: %v", err)
	}
	res, err := Apply(root, Request{ID: "G-260101-00001", Set: []Field{{"status", "active"}}}, now, nil)
	if err != nil || !res.Changed || res.Commit != "" || head() != base {
		t.Fatalf("omitted expect: %+v %v", res, err)
	}
	if src := read(t, root, "grove/work/G-260101-00001-first.md"); !strings.Contains(src, "status: active\n") || !strings.Contains(src, "Edited body.") {
		t.Fatalf("update not applied over the direct edit:\n%s", src)
	}
	// Other changes in the tree, staged and unstaged, are left alone by --commit.
	write(t, root, "grove/work/G-260101-00002-second.md", strings.Replace(second, "Second", "Second edited", 1))
	write(t, root, "notes.txt", "staged\n")
	git(t, root, "add", "notes.txt")
	res, err = Apply(root, Request{ID: "G-260101-00001", Set: []Field{{"status", "review"}, {"candidate", base}}, Commit: true}, now, nil)
	if err != nil || !res.Changed || res.Commit != head() || res.Commit == base {
		t.Fatalf("commit: %+v %v (HEAD %s)", res, err, head())
	}
	if files := git(t, root, "show", "--stat", "--format=", "--name-only", "HEAD"); files != "grove/work/G-260101-00001-first.md" {
		t.Fatalf("the commit must hold the record alone: %q", files)
	}
	if subject := git(t, root, "log", "-1", "--format=%s"); subject != "docs(G-260101-00001): set status=review candidate="+base {
		t.Fatalf("message: %q", subject)
	}
	if status := git(t, root, "status", "--porcelain"); status != "M grove/work/G-260101-00002-second.md\nA  notes.txt" { // git() trims the leading space
		t.Fatalf("other paths must stay as they were:\n%s", status)
	}
	// Acceptance 1: a status from a clean tree in one step, with several fields and an unset in the message.
	git(t, root, "commit", "-qam", "the rest")
	clean := head()
	res, err = Apply(root, Request{ID: "G-260101-00001", Set: []Field{{"status", "abandoned"}}, Unset: []string{"relates_to"}, Commit: true}, now, nil)
	if err != nil || res.Commit != head() || res.Commit == clean {
		t.Fatalf("abandoned --commit: %+v %v", res, err)
	}
	if subject := git(t, root, "log", "-1", "--format=%s"); subject != "docs(G-260101-00001): set status=abandoned unset relates_to" {
		t.Fatalf("message: %q", subject)
	}
	if status := git(t, root, "status", "--porcelain"); status != "" {
		t.Fatalf("tree must be clean after the commit:\n%s", status)
	}
	// A no-op commits nothing.
	done := head()
	res, err = Apply(root, Request{ID: "G-260101-00001", Set: []Field{{"status", "abandoned"}}, Commit: true}, now, nil)
	if err != nil || res.Changed || res.Commit != "" || head() != done {
		t.Fatalf("no-op with commit: %+v %v", res, err)
	}
	// A commit that fails after publication: the file holds the update, the
	// error carries the applied revision, and nothing was committed.
	hooks := filepath.Join(root, "hooks")
	write(t, root, "hooks/pre-commit", "#!/bin/sh\necho refused by hook >&2\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "config", "core.hooksPath", hooks)
	res, err = Apply(root, Request{ID: "G-260101-00001", Set: []Field{{"title", "Renamed"}}, Commit: true}, now, nil)
	var failure *Failure
	if !errors.As(err, &failure) || failure.Revision != revision(t, root, "grove/work/G-260101-00001-first.md") || !strings.Contains(err.Error(), "refused by hook") || !strings.Contains(err.Error(), "the file is staged but nothing was committed") || !strings.Contains(err.Error(), "the update was applied to grove/work/G-260101-00001-first.md") {
		t.Fatalf("failed commit: %+v %v", res, err)
	}
	if head() != done || record(t, root, "G-260101-00001").Title != "Renamed" {
		t.Fatal("the file must hold the update and HEAD must not move")
	}
}

// TestUpdateApprovedBindsToTheCandidate covers the acceptance fields: approved
// is written quoted and must name the candidate, all three hold only while
// accepted, and an acceptance context written must be the record's own, so
// neither a changed candidate nor a reopened record keeps an acceptance.
func TestUpdateApprovedBindsToTheCandidate(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	head := git(t, root, "rev-parse", "HEAD")
	refuse := func(want string, sets []Field, unsets ...string) {
		t.Helper()
		r := record(t, root, "G-260101-00001")
		_, err := Apply(root, Request{ID: "G-260101-00001", Expect: project.Revision(r.Source), Set: sets, Unset: unsets}, now, nil)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("%v %v: got %v, want %q", sets, unsets, err, want)
		}
	}
	apply(t, root, "G-260101-00001", []Field{{"candidate", head}}) // allowed on every status; approval is not
	refuse("approved: approval holds only while status is accepted or done", []Field{{"approved", head}})
	apply(t, root, "G-260101-00001", []Field{{"status", "review"}})
	refuse("approved: approval holds only while status is accepted or done, not review", []Field{{"approved", head}})
	context := project.AcceptanceContext(record(t, root, "G-260101-00001"))
	accept := []Field{{"status", "accepted"}, {"approved", head}, {"approved_by", "owner"}, {"approved_context", context}}
	refuse("approved_context: required while status is accepted", accept[:3])
	refuse("approved: approval is of one commit and must name the candidate", []Field{accept[0], {"approved", "abcdef0"}, accept[2], accept[3]})
	refuse("approved_context must be this record's acceptance context, "+context, []Field{accept[0], accept[1], accept[2], {"approved_context", "sha256:" + strings.Repeat("0", 64)}})
	apply(t, root, "G-260101-00001", accept)
	if r := record(t, root, "G-260101-00001"); r.Approved != head || !strings.Contains(string(r.Source), "approved: \""+head+"\"\napproved_by: owner\napproved_context: \""+context+"\"\n") {
		t.Fatalf("the acceptance must be written quoted: %s", r.Source)
	}
	refuse("approved: approval is of one commit and must name the candidate", []Field{{"candidate", "abcdef0"}})
	refuse("approved: approval holds only while status is accepted or done", []Field{{"status", "active"}})
	refuse("approved_by: belongs to an acceptance", []Field{{"status", "active"}}, "approved")
	apply(t, root, "G-260101-00001", []Field{{"status", "active"}}, "approved", "approved_by", "approved_context")
	if r := record(t, root, "G-260101-00001"); r.Approved != "" || r.Candidate != head || r.Status != "active" {
		t.Fatalf("reopening with the acceptance unset: %+v", r)
	}
}

// TestUpdateRefusalsNameTheirRule is each refusal as a command's author reads
// it: the accepted values or fields are listed, and a lifecycle refusal says
// its rule and the correction, which then succeeds.
func TestUpdateRefusalsNameTheirRule(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	write(t, root, "grove/work/G-260101-00005-done.md", "---\nid: \"G-260101-00005\"\ntype: work\ntitle: Done\nstatus: done\ncandidate: \"abcdef0\"\napproved: \"abcdef0\"\n---\nBody.\n")
	for _, tc := range []struct {
		id   string
		set  Field
		want string
	}{
		{"G-260101-00001", Field{"status", "bogus"}, "status: expected proposed, active, review, accepted, abandoned or done for work"},
		{"G-260101-00001", Field{"size", "huge"}, "size: expected small, medium or large"},
		{"G-260101-00001", Field{"kind", "chore"}, "kind: expected feature, fix, refactor, investigation, tooling or release"},
		{"G-260101-00001", Field{"zzz", "1"}, "zzz is not a field that update accepts on work records; it accepts type, title, status, relates_to, kind, size, priority, members, depends_on, candidate, approved, approved_by or approved_context"},
		{"G-260101-00003", Field{"work", "[]"}, "work is not a field that update accepts on question records; it accepts type, title, status, relates_to or blocks"},
		{"G-260101-00004", Field{"status", "open"}, "status: expected proposed, accepted, rejected or superseded for decision"},
		{"G-260101-00005", Field{"status", "active"}, "approved: approval holds only while status is accepted or done, not active: unset approved, or set status accepted"},
	} {
		if _, err := Apply(root, Request{ID: tc.id, Set: []Field{tc.set}}, now, nil); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s %s=%s: wanted %q, got %v", tc.id, tc.set.Name, tc.set.Value, tc.want, err)
		}
	}
	if _, err := Apply(root, Request{ID: "G-260101-00005", Set: []Field{{"status", "active"}}, Unset: []string{"approved"}}, now, nil); err != nil {
		t.Fatalf("the correction the refusal names was refused: %v", err)
	}
}
