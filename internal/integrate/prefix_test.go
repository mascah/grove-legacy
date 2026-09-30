package integrate

import (
	"strings"
	"testing"
)

// The commit the checks read is what is delivered, whatever the branch's
// name resolves to meanwhile; and a project under a prefix compares the
// record's path where git diff prints it, so code above the project counts
// and the record's own commits do not.
func TestIntegrateDeliversTheInspectedCommitUnderAPrefix(t *testing.T) {
	t.Parallel()
	root, _, candidate := fixture(t, true, "sub")
	git(t, root, "tag", "feature", "main") // a tag of the branch's name
	facts, err := run(t, root, root, false)
	if err != nil || len(facts) != 4 || !strings.HasPrefix(facts[3], "done: G-260101-00001 is done: squashed as ") {
		t.Fatalf("%v; facts %q", err, facts)
	}
	if r := record(t, root); r.Status != "accepted" || r.Candidate != candidate {
		t.Fatalf("%+v", r)
	}
	if got := git(t, root, "show", "main:code.txt"); got != "the change" {
		t.Fatalf("code above the project: %q", got)
	}
	if !strings.HasPrefix(git(t, root, "log", "-1", "--format=%s", "main"), "feat: first") {
		t.Fatal("the type comes from code above the project")
	}
}
