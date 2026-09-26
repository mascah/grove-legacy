package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

func TestPagesAndConvertThroughTheCLI(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	page := newID(t, root, "page", "Probe synthesis")
	// A page has no status, and is never selectable work.
	if code, out, errOut := run(t, root, "list"); code != 0 || !strings.Contains(out, page+"  page      -") {
		t.Fatalf("list: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "context", page); code != 1 || out != "" || !strings.Contains(errOut, page+" is a page; only work can be selected") {
		t.Fatalf("context: code=%d stdout=%q stderr=%s", code, out, errOut)
	}

	write(t, root, "docs/legacy-note.md", "# Legacy note\n\nBody.\n")
	code, out, errOut := run(t, root, "convert", "docs/legacy-note.md", "--type", "plan", "--title", "Legacy note", "--slug", "note")
	var c struct{ ID string }
	json.Unmarshal([]byte(out), &c)
	want := `{"from":"docs/legacy-note.md","from_path":"docs/legacy-note.md","id":"` + c.ID + `","path":"docs/records/` + c.ID + `-note.md"}` + "\n"
	if code != 0 || out != want || !project.IDPattern.MatchString(c.ID) || len(c.ID) != len("G-260925-7k2qm") {
		t.Fatalf("convert: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if source := showJSON(t, root, c.ID)["source"].(string); !strings.Contains(source, "type: plan\ntitle: \"Legacy note\"\nstatus: current\nformerly: \"docs/legacy-note.md\"\n") || !strings.HasSuffix(source, "# Legacy note\n\nBody.\n") {
		t.Fatalf("converted record:\n%s", source)
	}
	// A rerun is refused and writes nothing.
	if code, out, errOut := run(t, root, "convert", "docs/legacy-note.md", "--type", "plan", "--title", "Legacy note"); code != 1 || out != "" || !strings.Contains(errOut, "already converted to "+c.ID) {
		t.Fatalf("rerun: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	newID(t, root, "page", "Next")
	if code, _, errOut := run(t, root, "check"); code != 0 {
		t.Fatal(errOut)
	}
}

func TestConvertUsageErrors(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	for _, args := range [][]string{
		{"convert"},
		{"convert", "a.md", "b.md", "--type", "plan", "--title", "T"},
	} {
		if code, out, errOut := run(t, root, args...); code != 2 || out != "" || !strings.Contains(errOut, "convert requires exactly one document path") {
			t.Fatalf("%v: code=%d stdout=%q stderr=%s", args, code, out, errOut)
		}
	}
}
