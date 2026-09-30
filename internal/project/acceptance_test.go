package project

import (
	"strings"
	"testing"
)

func TestAcceptanceContext(t *testing.T) {
	t.Parallel()
	base := "---\nid: G-260101-00001\ntype: work\ntitle: T\nstatus: review\ncandidate: \"abcdef0\"\n---\n\n## Outcome\n\nIt works.\n\n## Acceptance\n\n1. Tested.\n\n## Next\n\nHand off.\n"
	context := func(source string) string {
		r, ds := ParseRecord("grove/a.md", []byte(source))
		if len(ds) != 0 {
			t.Fatalf("%s: %v", source, ds)
		}
		return AcceptanceContext(r)
	}
	want := context(base)
	for name, source := range map[string]string{
		"acceptance fields":        strings.Replace(base, "status: review\n", "status: accepted\napproved: \"abcdef0\"\napproved_by: owner\napproved_context: \""+want+"\"\nupdated: \"2026-01-01T00:00:00Z\"\n", 1),
		"next":                     strings.Replace(base, "Hand off.", "Delivered; nothing left.\n\nMore progress.", 1),
		"appended verdict":         base + "\nVerdict on candidate abcdef0, 2026-01-01: ship it\n",
		"appended feedback":        base + "\nFeedback on candidate abcdef0, 2026-01-01: redo\n\nReopened with G-260101-00002's feedback.\n",
		"crlf":                     strings.ReplaceAll(base, "\n", "\r\n"),
		"bom":                      "\ufeff" + base,
		"trailing spaces and gaps": strings.Replace(base, "It works.\n", "It works.  \n\n\n", 1),
	} {
		if got := context(source); got != want {
			t.Errorf("%s changed the context", name)
		}
	}
	for name, source := range map[string]string{
		"outcome":              strings.Replace(base, "It works.", "It works fast.", 1),
		"acceptance":           strings.Replace(base, "1. Tested.", "1. Tested.\n2. Documented.", 1),
		"title":                strings.Replace(base, "title: T", "title: U", 1),
		"a section after next": base + "\n## Scope\n\nNarrower.\n",
		"a paragraph that only mentions a verdict": strings.Replace(base, "It works.", "It works. Verdict on candidate abcdef0 pending.", 1),
	} {
		if got := context(source); got == want {
			t.Errorf("%s left the context unchanged", name)
		}
	}
}

func TestSchema3IsReadOnlyForMigration(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	put(t, root, "grove/G-260101-00002.md", typed("G-260101-00002", "work", "review", "candidate: \"abcdef0\"\napproved: \"abcdef0\"\n"))
	if _, ds := Load(root, root); !strings.Contains(diagnostics(ds), "grove migrate previews") {
		t.Fatalf("schema 3 loaded: %s", diagnostics(ds))
	}
	p, ds := LoadSchema3(root)
	if len(ds) != 0 || p.Records[0].Approved != "abcdef0" {
		t.Fatalf("schema 3 approval in review refused: %s", diagnostics(ds))
	}
	put(t, root, "grove/G-260101-00002.md", typed("G-260101-00002", "work", "accepted", "candidate: \"abcdef0\"\n"))
	if _, ds := LoadSchema3(root); !strings.Contains(diagnostics(ds), "status: expected proposed, active, review, done or abandoned") {
		t.Fatalf("schema 4 status read as schema 3: %s", diagnostics(ds))
	}
	put(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n")
	if _, ds := LoadSchema3(root); !strings.Contains(diagnostics(ds), "unsupported version 4; expected 3") {
		t.Fatalf("schema 4 read as schema 3: %s", diagnostics(ds))
	}
}
