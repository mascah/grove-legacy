package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// knowledgeFixture is the Git fixture with a brief added.
func knowledgeFixture(t *testing.T) string {
	t.Helper()
	root := gitFixture(t)
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: docs/records\nbrief: docs/brief.md\n")
	write(t, root, "docs/brief.md", "# Brief\n")
	return root
}

func run(t *testing.T, root string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Run(args, root, &out, &errOut)
	return code, out.String(), errOut.String()
}

// newID runs new with args and returns the ID it issued, read from the
// path it prints: IDs are random, so a test names none in advance.
func newID(t *testing.T, root string, args ...string) string {
	t.Helper()
	code, out, errOut := run(t, root, append([]string{"new"}, args...)...)
	id := filepath.Base(out)[:min(len(filepath.Base(out)), len("G-260925-7k2qm"))]
	if code != 0 || !project.IDPattern.MatchString(id) || !strings.HasSuffix(out, ".md\n") {
		t.Fatalf("new %v: code=%d stdout=%q stderr=%s", args, code, out, errOut)
	}
	return id
}

func TestKnowledgeRecordsThroughNewUpdateAndContext(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	// G-260101-00001 (work) and G-260101-00002 (question) already exist from the fixture.
	term := newID(t, root, "term", "Attempt")
	plan := newID(t, root, "plan", "Shared plan")
	review := newID(t, root, "review", "First review")
	for id, want := range map[string]string{term: "-attempt.md", plan: "-shared-plan.md", review: "-first-review.md"} {
		if _, err := os.Stat(filepath.Join(root, "docs/records", id+want)); err != nil {
			t.Fatal(err)
		}
	}
	// A second record for one term is refused before anything is written, so
	// the project stays loadable.
	if code, out, errOut := run(t, root, "new", "term", " attempt "); code != 1 || out != "" || !strings.Contains(errOut, "the term attempt is already defined by "+term+" in docs/records/"+term+"-attempt.md") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	newID(t, root, "term", "Candidate")
	revision := func(id string) string { return showJSON(t, root, id)["revision"].(string) }
	if code, _, errOut := run(t, root, "update", plan, "--expect", revision(plan), "--set", `work=["G-260101-00001"]`); code != 0 {
		t.Fatal(errOut)
	}
	// An all-letter or all-digit commit must stay a YAML string.
	for _, commit := range []string{"abcdefa", "1234567"} {
		if code, _, errOut := run(t, root, "update", review, "--expect", revision(review), "--set", `work=["G-260101-00001"]`, "--set", "examined="+commit); code != 0 {
			t.Fatalf("examined=%s: %s", commit, errOut)
		}
		if source := showJSON(t, root, review)["source"].(string); !strings.Contains(source, `examined: "`+commit+`"`) {
			t.Fatalf("examined must be quoted:\n%s", source)
		}
	}
	for _, c := range []struct{ id, set, want string }{
		{review, "examined=main", "examined: expected a quoted Git commit"},
		{plan, `work=["G-260101-00404"]`, "work: unresolved target G-260101-00404"},
		{plan, `work=["` + term + `"]`, "work: target " + term + " must be work"},
		{plan, "examined=abcdefa", "examined is not a field that update accepts on plan records"},
		{term, `work=["G-260101-00001"]`, "work is not a field that update accepts on term records"},
	} {
		before := revision(c.id)
		if code, _, errOut := run(t, root, "update", c.id, "--expect", before, "--set", c.set); code != 1 || !strings.Contains(errOut, c.want) || revision(c.id) != before {
			t.Fatalf("%s %s: code=%d stderr=%s", c.id, c.set, code, errOut)
		}
	}
	if code, out, errOut := run(t, root, "check"); code != 0 || out != "OK: 6 records\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	code, out, errOut := run(t, root, "context", "G-260101-00001")
	if code != 0 {
		t.Fatal(errOut)
	}
	for _, want := range []string{plan + "  plan  current  listed  plan for G-260101-00001", review + "  review  current  listed  review of G-260101-00001"} {
		if !strings.Contains(out, want) {
			t.Fatalf("context must list %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Source: docs/records/"+plan) || strings.Contains(out, "Source: docs/brief.md") || strings.Contains(out, term) {
		t.Fatalf("context must not include plans or the brief, or list unrelated terms:\n%s", out)
	}
	// One plan shared by several work items is listed for each, never copied.
	second := newID(t, root, "work", "Second")
	if code, _, errOut := run(t, root, "update", plan, "--expect", revision(plan), "--set", `work=["G-260101-00001", "`+second+`"]`); code != 0 {
		t.Fatal(errOut)
	}
	if code, out, errOut := run(t, root, "context", second, "G-260101-00001"); code != 0 || !strings.Contains(out, plan+"  plan  current  listed  plan for "+second+"; plan for G-260101-00001") {
		t.Fatalf("code=%d stderr=%s\n%s", code, errOut, out)
	}
	if code, out, errOut = run(t, root, "context", "G-260101-00001", "--include", "docs/records/"+plan+"-shared-plan.md"); code != 0 || !strings.Contains(out, plan+"  plan  current  included") {
		t.Fatalf("a caller can include the plan: code=%d stderr=%s\n%s", code, errOut, out)
	}
}

func TestBriefCommand(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	if code, out, errOut := run(t, root, "brief"); code != 0 || out != "# Brief\n" || !strings.Contains(errOut, "File: docs/brief.md\n") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, _ := run(t, root, "brief", "--json"); code != 0 || !strings.Contains(out, `"path":"docs/brief.md"`) || !strings.Contains(out, `"revision":"sha256:`) {
		t.Fatalf("code=%d stdout=%q", code, out)
	}
	if code, _, _ := run(t, root, "brief", "extra"); code != 2 {
		t.Fatalf("brief takes no arguments: %d", code)
	}
	if err := os.Remove(filepath.Join(root, "docs/brief.md")); err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"brief", "check", "list"} {
		if code, out, errOut := run(t, root, command); code != 1 || out != "" || !strings.Contains(errOut, "grove.yaml: brief: ") {
			t.Fatalf("%s with a missing brief: code=%d stdout=%q stderr=%s", command, code, out, errOut)
		}
	}
	plain := gitFixture(t)
	if code, _, errOut := run(t, plain, "brief"); code != 1 || !strings.Contains(errOut, "grove.yaml names no brief") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
}

// A committed tree goes through the versions tree reader, which gives the
// loader only grove.yaml and the record folder.
func TestVersionsReadsCommittedKnowledgeRecords(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: docs/records\nbrief: docs/records/brief.md\n")
	write(t, root, "docs/records/brief.md", "# Brief\n")
	var ids []string
	for _, kind := range []string{"term", "plan", "review"} {
		ids = append(ids, newID(t, root, kind, "A "+kind))
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "commit", "-q", "-m", "knowledge records")
	code, out, errOut := run(t, root, "versions")
	if code != 0 || strings.Contains(errOut, "invalid") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	for _, want := range []string{ids[0] + "  proposed  committed refs/heads/main", ids[1] + "  current   committed refs/heads/main", ids[2] + "  current   committed refs/heads/main"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "brief") {
		t.Fatalf("the brief is not a record:\n%s", out)
	}
	selector := versionSelector(t, root, ids[0], "committed refs/heads/main")
	if code, out, errOut := run(t, root, "workspace", "--source", selector); code != 0 || out != root+"\n" {
		t.Fatalf("a term selector must resolve: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
}
