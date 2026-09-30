package project

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func put(t *testing.T, root, path, content string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	put(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n")
	if err := os.Mkdir(filepath.Join(root, "grove"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func record(id, kind, extra string) string {
	status := "proposed"
	if kind == "question" {
		status = "open"
	}
	return "---\nid: " + id + "\ntype: " + kind + "\ntitle: Example\nstatus: " + status + "\n" + extra + "---\nBody with --- inside.\n"
}

func diagnostics(ds []Diagnostic) string {
	var out []string
	for _, d := range ds {
		out = append(out, d.String())
	}
	return strings.Join(out, "\n")
}

func TestLoadRenamedNestedRecordsAndPreserveSource(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	source := record("G-260101-00001", "work", "kind: investigation\npriority: 2\nsize: small\ncreated: \"2026-09-19T14:08:40Z\"\n")
	source = string(rune(0xFEFF)) + strings.ReplaceAll(source, "\n", "\r\n")
	put(t, root, "grove/work/nested/custom.md", source)
	put(t, root, "grove/questions/anything.md", record("G-260101-00002", "question", ""))
	put(t, root, "grove/notes.txt", "not a record")
	p, ds := Load(filepath.Join(root, "grove", "work", "nested"), "")
	if len(ds) != 0 {
		t.Fatal(diagnostics(ds))
	}
	if len(p.Records) != 2 || p.Records[0].ID != "G-260101-00001" {
		t.Fatalf("records: %+v", p.Records)
	}
	if string(p.Records[0].Source) != source || p.Records[0].Path != "grove/work/nested/custom.md" {
		t.Fatal("source bytes or path changed")
	}
	if p.Records[1].Priority != nil || p.Records[1].Size != "" {
		t.Fatal("missing planning values must remain unspecified")
	}
}

func TestDiscoveryAndExplicitProject(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "nested/.git", "gitdir: unused-by-reader\n")
	_, ds := Load(filepath.Join(root, "nested"), "")
	if !strings.Contains(diagnostics(ds), "grove.yaml") {
		t.Fatal("nested Git checkout inherited parent configuration")
	}
	p, ds := Load(filepath.Join(root, "nested"), "..")
	if len(ds) != 0 || p.Root != root {
		t.Fatalf("explicit parent project: %v", diagnostics(ds))
	}
	put(t, root, "nested/grove.yaml", "schema_version: 4\nrecords: records\n")
	if err := os.Mkdir(filepath.Join(root, "nested", "records"), 0755); err != nil {
		t.Fatal(err)
	}
	p, ds = Load(filepath.Join(root, "nested"), "")
	if len(ds) != 0 || p.Root != filepath.Join(root, "nested") {
		t.Fatalf("nearest project: %v", diagnostics(ds))
	}
	_, ds = Load(root, "grove")
	if len(ds) == 0 {
		t.Fatal("explicit project must not search upward")
	}
}

func TestInvalidConfiguration(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, source, field string }{
		{"version missing", "records: grove\n", "schema_version"},
		{"unsupported", "schema_version: 5\nrecords: grove\n", "schema_version"},
		{"schema 1 no longer supported", "schema_version: 1\nrecords: grove\n", "schema_version"},
		{"schema 2 no longer supported", "schema_version: 2\nrecords: grove\n", "schema_version"},
		{"schema 3 needs migration", "schema_version: 3\nrecords: grove\n", "schema_version"},
		{"quoted version", "schema_version: \"4\"\nrecords: grove\n", "schema_version"},
		{"fractional version", "schema_version: 4.0\nrecords: grove\n", "schema_version"},
		{"null root", "schema_version: 4\nrecords: null\n", "records"},
		{"escape", "schema_version: 4\nrecords: ../outside\n", "records"},
		{"hidden escape", "schema_version: 4\nrecords: a/../grove\n", "records"},
		{"project root", "schema_version: 4\nrecords: .\n", "records"},
		{"absolute root", "schema_version: 4\nrecords: /tmp\n", "records"},
		{"unknown", "schema_version: 4\nrecords: grove\nextra: true\n", "extra"},
		{"duplicate", "schema_version: 4\nrecords: grove\nrecords: other\n", "records"},
		{"extra document", "schema_version: 4\nrecords: grove\n---\n{}\n", "document"},
		{"missing root", "schema_version: 4\nrecords: absent\n", "records"},
		{"non mapping", "- records\n", "mapping"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := fixture(t)
			put(t, root, "grove.yaml", tc.source)
			_, ds := Load(root, "")
			msg := diagnostics(ds)
			if !strings.Contains(msg, "grove.yaml") || !strings.Contains(msg, tc.field) {
				t.Fatalf("wanted file and %q diagnostic; got %s", tc.field, msg)
			}
		})
	}
}

func TestStrictRecordMetadata(t *testing.T) {
	t.Parallel()
	base := record("G-260101-00001", "work", "")
	for _, tc := range []struct{ name, source, field string }{
		{"missing title", strings.Replace(base, "title: Example\n", "", 1), "title"},
		{"numeric title", strings.Replace(base, "title: Example", "title: 123", 1), "title"},
		{"empty title", strings.Replace(base, "title: Example", "title: \"  \"", 1), "title"},
		{"null", record("G-260101-00001", "work", "size: null\n"), "size"},
		{"priority float", record("G-260101-00001", "work", "priority: 2.0\n"), "priority"},
		{"priority range", record("G-260101-00001", "work", "priority: 6\n"), "priority"},
		{"size", record("G-260101-00001", "work", "size: spike\n"), "size"},
		{"kind", record("G-260101-00001", "work", "kind: magic\n"), "kind"},
		{"unknown", record("G-260101-00001", "work", "priorty: 2\n"), "priorty"},
		{"duplicate key", record("G-260101-00001", "work", "title: Again\n"), "title"},
		{"duplicate relationship", record("G-260101-00001", "work", "relates_to: [G-260101-00002, G-260101-00002]\n"), "relates_to"},
		{"sequence required", record("G-260101-00001", "work", "depends_on: G-260101-00002\n"), "depends_on"},
		{"nonstring target", record("G-260101-00001", "work", "depends_on: [42]\n"), "depends_on"},
		{"wrong type field", record("G-260101-00001", "work", "blocks: []\n"), "blocks"},
		{"bare timestamp", record("G-260101-00001", "work", "created: 2026-09-19T12:00:00Z\n"), "created"},
		{"offset timestamp", record("G-260101-00001", "work", "created: \"2026-09-19T12:00:00+00:00\"\n"), "created"},
		{"fraction timestamp", record("G-260101-00001", "work", "created: \"2026-09-19T12:00:00.123Z\"\n"), "created"},
		{"invalid date", record("G-260101-00001", "work", "created: \"2026-02-30T12:00:00Z\"\n"), "created"},
		{"backward date", record("G-260101-00001", "work", "created: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-18T12:00:00Z\"\n"), "updated"},
		{"backward zero date", record("G-260101-00001", "work", "created: \"0001-01-01T00:00:00Z\"\nupdated: \"0000-01-01T00:00:00Z\"\n"), "updated"},
		{"short id", record("G-1", "work", ""), "id"},
		{"extra padding", record("G-0001", "work", ""), "id"},
		{"short numeric id", record("G-01", "work", ""), "id"},
		{"long numeric id", record("G-1234", "work", ""), "id"},
		{"retired numeric id", record("G-"+"001", "work", ""), "id"}, // split so a search for the retired form finds no fixture (G-260926-19gzg)
		{"short tail", record("G-260925-7k2q", "work", ""), "id"},
		{"uppercase tail", record("G-260925-7K2QM", "work", ""), "id"},
		{"ambiguous tail letter", record("G-260925-7k2ql", "work", ""), "id"},
		{"short date", record("G-26092-7k2qm", "work", ""), "id"},
		{"unknown prefix", record("X-001", "work", ""), "id"},
		{"lifecycle", strings.Replace(base, "status: proposed", "status: resolved", 1), "status"},
		{"unclosed header", strings.TrimSuffix(base, "---\nBody with --- inside.\n"), "frontmatter"},
		{"no header", "Hello\n", "frontmatter"},
		{"alias", "---\nid: G-260101-00001\ntype: work\ntitle: &label Example\nstatus: *label\n---\n", "alias"},
		{"custom tag", strings.Replace(base, "title: Example", "title: !special Example", 1), "title"},
		{"invalid utf8", base + "\xff", "UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := fixture(t)
			put(t, root, "grove/work/custom.md", tc.source)
			_, ds := Load(root, "")
			msg := diagnostics(ds)
			if !strings.Contains(msg, "custom.md") || !strings.Contains(msg, tc.field) {
				t.Fatalf("wanted file and %q diagnostic; got %s", tc.field, msg)
			}
		})
	}
}

func TestRejectSymlinks(t *testing.T) {
	t.Parallel()
	for _, target := range []string{"grove.yaml", "grove", "grove/work", "grove/work/link.md", "grove/link.txt"} {
		t.Run(target, func(t *testing.T) {
			root := fixture(t)
			outside := t.TempDir()
			put(t, outside, "target", record("G-260101-00001", "work", ""))
			path := filepath.Join(root, target)
			os.Remove(path)
			os.MkdirAll(filepath.Dir(path), 0755)
			if err := os.Symlink(filepath.Join(outside, "target"), path); err != nil {
				t.Fatal(err)
			}
			_, ds := Load(root, "")
			if !strings.Contains(diagnostics(ds), "symlink") {
				t.Fatalf("followed %s: %s", target, diagnostics(ds))
			}
		})
	}
}

func TestExplicitZeroTimeIsNotAnAbsentDate(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/undated.md", record("G-260101-00001", "work", ""))
	put(t, root, "grove/work/dated.md", record("G-260101-00999", "work", "created: \"0001-01-01T00:00:00Z\"\n"))
	p, ds := Load(root, "")
	if len(ds) != 0 || p.Records[0].ID != "G-260101-00999" {
		t.Fatalf("a present timestamp must sort before undated records: %s", diagnostics(ds))
	}
}

func configured(t *testing.T, config string) string {
	t.Helper()
	root := fixture(t)
	put(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\n"+config)
	return root
}

func typed(id, kind, status, extra string) string {
	return "---\nid: " + id + "\ntype: " + kind + "\ntitle: " + id + " title\nstatus: " + status + "\n" + extra + "---\nBody.\n"
}

func TestKnowledgeRecords(t *testing.T) {
	t.Parallel()
	root := configured(t, "")
	put(t, root, "grove/G-260101-00001.md", typed("G-260101-00001", "work", "proposed", ""))
	put(t, root, "grove/terms/G-260101-00002-attempt.md", typed("G-260101-00002", "term", "settled", "relates_to: [\"G-260101-00001\"]\n"))
	put(t, root, "grove/plans/G-260101-00003-shared.md", typed("G-260101-00003", "plan", "current", "work: [\"G-260101-00001\"]\n"))
	put(t, root, "grove/reviews/G-260101-00004-first.md", typed("G-260101-00004", "review", "current", "work: [\"G-260101-00001\"]\nexamined: \"42c077d\"\n"))
	p, ds := Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics:\n%s", diagnostics(ds))
	}
	if len(p.Records) != 4 {
		t.Fatalf("got %d records", len(p.Records))
	}
	for _, r := range p.Records {
		if r.ID == "G-260101-00004" && (r.Examined != "42c077d" || len(r.Work) != 1) {
			t.Fatalf("review fields not read: %+v", r)
		}
	}
}

func TestKnowledgeRecordProblemsNameFileAndField(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, source, want string }{
		{"missing work target", "grove/plans/G-260101-00002.md", typed("G-260101-00002", "plan", "current", "work: [\"G-260101-00009\"]\n"), "grove/plans/G-260101-00002.md: work: unresolved target G-260101-00009"},
		{"work target not work", "grove/plans/G-260101-00002.md", typed("G-260101-00002", "plan", "current", "work: [\"G-260101-00002\"]\n"), "work: self-reference"},
		{"work names a term", "grove/reviews/G-260101-00002.md", typed("G-260101-00002", "review", "current", "work: [\"G-260101-00001\"]\n"), "grove/reviews/G-260101-00002.md: work: target G-260101-00001 must be work"},
		{"unquoted numeric commit", "grove/reviews/G-260101-00002.md", typed("G-260101-00002", "review", "current", "examined: 1234567\n"), "examined: expected a nonempty string"},
		{"not a commit", "grove/reviews/G-260101-00002.md", typed("G-260101-00002", "review", "current", "examined: \"main\"\n"), "examined: expected a quoted Git commit"},
		{"examined on a plan", "grove/plans/G-260101-00002.md", typed("G-260101-00002", "plan", "current", "examined: \"42c077d\"\n"), "examined: not a field of plan records, which take"},
		{"work on a term", "grove/terms/G-260101-00003.md", typed("G-260101-00003", "term", "proposed", "work: [\"G-260101-00001\"]\n"), "work: not a field of term records, which take id, type, title, status, relates_to, created, updated or formerly"},
		{"bad status", "grove/terms/G-260101-00003.md", typed("G-260101-00003", "term", "done", ""), "status: expected proposed or settled for term"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := configured(t, "")
			put(t, root, "grove/G-260101-00001.md", typed("G-260101-00001", "term", "settled", ""))
			put(t, root, tc.path, tc.source)
			_, ds := Load(root, root)
			if got := diagnostics(ds); !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
		})
	}
}

func TestDuplicateTermTitle(t *testing.T) {
	t.Parallel()
	root := configured(t, "")
	put(t, root, "grove/terms/G-260101-00001.md", "---\nid: G-260101-00001\ntype: term\ntitle: Attempt\nstatus: settled\n---\n")
	put(t, root, "grove/terms/G-260101-00002.md", "---\nid: G-260101-00002\ntype: term\ntitle: \" attempt\"\nstatus: proposed\n---\n")
	_, ds := Load(root, root)
	if got := diagnostics(ds); !strings.Contains(got, "grove/terms/G-260101-00002.md: title: term already defined by G-260101-00001") {
		t.Fatalf("got %s", got)
	}
}

// briefFixture is a project with one baseline work record, for tests
// that also need the brief: configuration key.
func briefFixture(t *testing.T, config string) string {
	t.Helper()
	root := configured(t, config)
	put(t, root, "grove/G-260101-00001.md", record("G-260101-00001", "work", ""))
	return root
}

func TestTarget(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, config, target, want string }{
		{"none", "", "", ""},
		{"a branch", "target: main\n", "main", ""},
		{"empty", "target: \"\"\n", "", "target: expected a nonempty string"},
		{"not a string", "target: [main]\n", "", "target: expected a nonempty string"},
		{"a full ref", "target: refs/heads/main\n", "", "target: must name a local branch"},
		{"spaced", "target: \" main\"\n", "", "target: must name a local branch"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, ds := Load(briefFixture(t, tc.config), "")
			if got := diagnostics(ds); tc.want == "" && got != "" || !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
			if tc.want == "" && p.Target != tc.target {
				t.Fatalf("Target = %q, want %q", p.Target, tc.target)
			}
		})
	}
}

func TestRunDefaults(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, config string
		run          RunDefaults
		want         string
	}{
		{"none", "", RunDefaults{}, ""},
		{"budget and mode", "run:\n  budget: 50\n  permission_mode: auto\n", RunDefaults{BudgetUSD: "50", PermissionMode: "auto"}, ""},
		{"all four", "run: {budget: 0.5, permission_mode: acceptEdits, model: opus, effort: xhigh}\n", RunDefaults{"0.5", "acceptEdits", "opus", "xhigh"}, ""},
		{"a word for a budget", "run: {budget: x}\n", RunDefaults{}, "grove.yaml:3: run.budget: expected a positive decimal dollar amount"},
		{"a zero budget", "run: {budget: 0}\n", RunDefaults{}, "run.budget: expected a positive"},
		{"a spaced value", "run: {effort: \"a b\"}\n", RunDefaults{}, "run.effort: expected one word"},
		{"a bound", "run: {until: plan}\n", RunDefaults{}, "run.until: unknown key"},
		{"not a mapping", "run: 50\n", RunDefaults{}, "run: expected a mapping"},
		{"a duplicate", "run:\n  model: a\n  model: b\n", RunDefaults{}, "grove.yaml:4: run.model: duplicate YAML key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, ds := Load(briefFixture(t, tc.config), "")
			if got := diagnostics(ds); tc.want == "" && got != "" || !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
			if tc.want == "" && p.Run != tc.run {
				t.Fatalf("Run = %+v, want %+v", p.Run, tc.run)
			}
		})
	}
}

func TestBrief(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, config, file, want string }{
		{"under the record root", "brief: grove/brief.md\n", "grove/brief.md", ""},
		{"elsewhere in the project", "brief: docs/restart-brief.md\n", "docs/restart-brief.md", ""},
		{"missing", "brief: grove/brief.md\n", "", "grove.yaml: brief: "},
		{"escaping", "brief: ../brief.md\n", "", "brief: must name a project-relative .md file"},
		{"unclean", "brief: ./grove/brief.md\n", "grove/brief.md", "brief: must name a project-relative .md file"},
		{"not Markdown", "brief: grove/brief.txt\n", "grove/brief.txt", "brief: must name a project-relative .md file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := briefFixture(t, tc.config)
			if tc.file != "" {
				put(t, root, tc.file, "# Brief\n")
			}
			p, ds := Load(root, root)
			got := diagnostics(ds)
			if tc.want == "" {
				if got != "" {
					t.Fatalf("unexpected diagnostics: %s", got)
				}
				if source, err := p.ReadBrief(); err != nil || string(source) != "# Brief\n" {
					t.Fatalf("ReadBrief = %q, %v", source, err)
				}
				if len(p.Records) != 1 {
					t.Fatalf("the brief must not be a record; got %d records", len(p.Records))
				}
			} else if !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
		})
	}
	t.Run("spelled with another case than the disk", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: grove/Notes/brief.md\n")
		put(t, root, "grove/notes/brief.md", "# Brief\n")
		if _, err := os.Stat(filepath.Join(root, "grove/Notes/brief.md")); err != nil {
			t.Skip("case-sensitive filesystem")
		}
		if _, ds := Load(root, root); len(ds) != 0 {
			t.Fatalf("the brief must be exempt however it is spelled: %s", diagnostics(ds))
		}
	})
	t.Run("symlinked parent directory", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: docs/brief.md\n")
		put(t, root, "real/brief.md", "# Brief\n")
		if err := os.Symlink("real", filepath.Join(root, "docs")); err != nil {
			t.Fatal(err)
		}
		_, ds := Load(root, root)
		if got := diagnostics(ds); !strings.Contains(got, "grove.yaml: brief: docs is a symlink") {
			t.Fatalf("got %s", got)
		}
	})
	t.Run("symlinked brief", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: docs/brief.md\n")
		put(t, root, "real.md", "# Brief\n")
		if err := os.Mkdir(filepath.Join(root, "docs"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../real.md", filepath.Join(root, "docs", "brief.md")); err != nil {
			t.Fatal(err)
		}
		_, ds := Load(root, root)
		if got := diagnostics(ds); !strings.Contains(got, "brief: symlink files are not supported") {
			t.Fatalf("got %s", got)
		}
	})
}

const page = "---\nid: G-260101-00002\ntype: page\ntitle: Synthesis\nrelates_to: [\"G-260101-00001\"]\n---\nNotes.\n"

// Location carries no meaning: the root, a type-named folder, and an
// arbitrary nested folder all hold any type, under recursive discovery.
func TestDiscoversByIdentityNotLocation(t *testing.T) {
	t.Parallel()
	root := configured(t, "brief: grove/brief.md\n")
	put(t, root, "grove/brief.md", "# Brief, not a record\n")
	put(t, root, "grove/G-260101-00001-work.md", typed("G-260101-00001", "work", "proposed", ""))
	put(t, root, "grove/G-260101-00002-synthesis.md", page)
	put(t, root, "grove/decisions/deep/er/G-260101-00003.md", typed("G-260101-00003", "plan", "current", "work: [\"G-260101-00001\"]\n"))
	put(t, root, "grove/work/G-260101-00004-typed-folder.md", typed("G-260101-00004", "work", "done", "depends_on: [\"G-260101-00001\"]\n"))
	put(t, root, "grove/G-260101-00005-reclassified.md", typed("G-260101-00005", "term", "settled", "formerly: \"docs/old.md\"\n"))
	put(t, root, "grove/notes.txt", "ignored")
	p, ds := Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics:\n%s", diagnostics(ds))
	}
	if len(p.Records) != 5 {
		t.Fatalf("got %d records", len(p.Records))
	}
	for _, r := range p.Records {
		if r.ID == "G-260101-00002" && (r.Type != "page" || r.Status != "" || r.Title != "Synthesis") {
			t.Fatalf("page not read: %+v", r)
		}
	}
}

func TestRecordProblems(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, source, want string }{
		{"page with a lifecycle", "grove/a.md", "---\nid: G-260101-00002\ntype: page\ntitle: T\nstatus: proposed\n---\n", "status: not a field of page records"},
		{"page with a work field", "grove/a.md", "---\nid: G-260101-00002\ntype: page\ntitle: T\npriority: 1\n---\n", "priority: not a field of page records"},
		{"no type is not a page", "grove/a.md", "---\nid: G-260101-00002\ntitle: T\n---\n", "type: required field is missing"},
		{"unknown type", "grove/a.md", typed("G-260101-00002", "note", "open", ""), "type: unknown record type"},
		{"work record missing status", "grove/a.md", "---\nid: G-260101-00002\ntype: work\ntitle: T\n---\n", "status: required field is missing"},
		{"work with a bad status", "grove/a.md", typed("G-260101-00002", "work", "settled", ""), "status: expected proposed, active, review, accepted, abandoned or done for work"},
		{"plain Markdown", "grove/a.md", "# Just prose\n", "frontmatter: expected an opening --- line"},
		{"uncanonical ID", "grove/a.md", typed("G-02", "work", "proposed", ""), "id: expected a canonical ID, e.g. G-260925-7k2qm"},
		{"unknown prefix", "grove/a.md", typed("X-002", "work", "proposed", ""), "id: expected a canonical ID"},
		{"duplicate ID across folders", "grove/x/y/a.md", typed("G-260101-00001", "term", "settled", ""), "id: duplicate G-260101-00001 in grove/G-260101-00001.md, grove/x/y/a.md"},
		{"non-work target", "grove/a.md", typed("G-260101-00002", "plan", "current", "work: [\"G-260101-00009\"]\n"), "work: target G-260101-00009 must be work"},
		{"formerly twice", "grove/a.md", typed("G-260101-00002", "work", "done", "formerly: \"W-007\"\n"), "formerly: W-007 was already converted to G-260101-00001"},
		{"formerly twice by another case", "grove/a.md", typed("G-260101-00002", "work", "done", "formerly: \"w-007\"\n"), "formerly: w-007 was already converted to G-260101-00001"},
		{"formerly still present", "grove/a.md", typed("G-260101-00002", "work", "done", "formerly: \"G-260101-00009\"\n"), "formerly: G-260101-00009 still exists in grove/G-260101-00009.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := configured(t, "")
			put(t, root, "grove/G-260101-00001.md", typed("G-260101-00001", "work", "proposed", "formerly: \"W-007\"\n"))
			put(t, root, "grove/G-260101-00009.md", "---\nid: G-260101-00009\ntype: page\ntitle: P\n---\n")
			put(t, root, tc.path, tc.source)
			_, ds := Load(root, root)
			if got := diagnostics(ds); !strings.Contains(got, tc.path+":") && !strings.Contains(got, tc.path+": ") || !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q naming %s; got:\n%s", tc.want, tc.path, got)
			}
		})
	}
}

func TestReviewLifecycleAndCandidate(t *testing.T) {
	t.Parallel()
	commit := "candidate: \"0123456789abcdef0123456789abcdef01234567\"\n"
	approved := "approved: \"0123456789abcdef0123456789abcdef01234567\"\n"
	digest := "sha256:" + strings.Repeat("ab", 32)
	accepted := "approved_by: owner\napproved_context: \"" + digest + "\"\n"
	for _, tc := range []struct{ name, source, want string }{
		{"review needs a candidate", typed("G-260101-00002", "work", "review", ""), "grove/G-260101-00002.md: candidate: required while status is review"},
		{"review with one", typed("G-260101-00002", "work", "review", commit), ""},
		{"historical done keeps validating", typed("G-260101-00002", "work", "done", ""), ""},
		{"done with one", typed("G-260101-00002", "work", "done", commit), ""},
		{"not a commit", typed("G-260101-00002", "work", "active", "candidate: \"main\"\n"), "candidate: expected a quoted Git commit"},
		{"unquoted", typed("G-260101-00002", "work", "active", "candidate: 1234567\n"), "candidate: expected a nonempty string"},
		{"only on work", typed("G-260101-00002", "plan", "current", commit), "candidate: not a field of plan records"},
		{"review is a work status only", typed("G-260101-00002", "plan", "review", ""), "status: expected current or superseded for plan"},
		{"approved in review", typed("G-260101-00002", "work", "review", commit+approved), "approved: approval holds only while status is accepted or done, not review"},
		{"accepted", typed("G-260101-00002", "work", "accepted", commit+approved+accepted), ""},
		{"accepted needs its authority", typed("G-260101-00002", "work", "accepted", commit+approved+"approved_context: \""+digest+"\"\n"), "approved_by: required while status is accepted"},
		{"accepted needs its context", typed("G-260101-00002", "work", "accepted", commit+approved+"approved_by: owner\n"), "approved_context: required while status is accepted"},
		{"accepted needs its approval", typed("G-260101-00002", "work", "accepted", commit+accepted), "approved: required while status is accepted"},
		{"accepted needs a candidate", typed("G-260101-00002", "work", "accepted", ""), "candidate: required while status is accepted"},
		{"a policy's authority", typed("G-260101-00002", "work", "accepted", commit+approved+"approved_by: \"policy "+digest+"\"\napproved_context: \""+digest+"\"\n"), ""},
		{"an unknown authority", typed("G-260101-00002", "work", "accepted", commit+approved+"approved_by: me\napproved_context: \""+digest+"\"\n"), "approved_by: expected owner, or policy sha256:HEX"},
		{"a malformed context", typed("G-260101-00002", "work", "accepted", commit+approved+"approved_by: owner\napproved_context: \"sha256:12\"\n"), "approved_context: expected sha256:"},
		{"acceptance fields off accepted", typed("G-260101-00002", "work", "review", commit+accepted), "approved_by: belongs to an acceptance and holds only while status is accepted, not review"},
		{"no acceptance fields on schema 3 done", typed("G-260101-00002", "work", "done", commit+approved+accepted), "approved_context: belongs to an acceptance"},
		{"approved in done", typed("G-260101-00002", "work", "done", commit+approved), ""},
		{"approved must name the candidate", typed("G-260101-00002", "work", "review", "candidate: \"abcdef0\"\n"+approved), "approved: approval is of one commit and must name the candidate abcdef0"},
		{"approved needs a candidate", typed("G-260101-00002", "work", "done", approved), "approved: approval is of the candidate, and there is none"},
		{"approved before review", typed("G-260101-00002", "work", "active", commit+approved), "approved: approval holds only while status is accepted or done, not active: unset approved"},
		{"approved not a commit", typed("G-260101-00002", "work", "review", commit+"approved: \"HEAD\"\n"), "approved: expected a quoted Git commit"},
		{"approved only on work", typed("G-260101-00002", "review", "current", approved), "approved: not a field of review records"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := configured(t, "")
			put(t, root, "grove/G-260101-00002.md", tc.source)
			p, ds := Load(root, root)
			if got := diagnostics(ds); tc.want == "" && got != "" || tc.want != "" && !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %q", tc.want, got)
			}
			if tc.want == "" && strings.Contains(tc.source, "candidate:") && p.Records[0].Candidate != "0123456789abcdef0123456789abcdef01234567" {
				t.Fatalf("candidate not read: %+v", p.Records[0])
			}
			if tc.want == "" && strings.Contains(tc.source, "approved:") && p.Records[0].Approved != p.Records[0].Candidate {
				t.Fatalf("approved not read: %+v", p.Records[0])
			}
		})
	}
}

func TestPolicy(t *testing.T) {
	t.Parallel()
	full := "policy:\n  budget: 30\n  resolve:\n    budget: 10\n  approve:\n    verify: [go test ./..., go vet ./...]\n    max_lines: 300\n    never: [grove.yaml, .github/**, \"*.sum\"]\n  integrate: true\n"
	for _, tc := range []struct {
		name, config string
		policy       *Policy
		want         string
	}{
		{"none", "", nil, ""},
		{"full", full, &Policy{BudgetUSD: "30", Resolve: true, ResolveBudgetUSD: "10", Approve: true, Verify: []string{"go test ./...", "go vet ./..."}, MaxLines: 300, Never: []string{"grove.yaml", ".github/**", "*.sum"}, Integrate: true}, ""},
		{"resolve only", "policy: {budget: 5, resolve: {}}\n", &Policy{BudgetUSD: "5", Resolve: true}, ""},
		{"empty", "policy: {}\n", &Policy{}, ""},
		{"resolve without a budget", "policy: {resolve: {}}\n", nil, "policy.budget: required with resolve"},
		{"approve without verify", "policy: {approve: {max_lines: 3}}\n", nil, "policy.approve.verify: required"},
		{"a zero bound", "policy: {approve: {verify: [x], max_lines: 0}}\n", nil, "policy.approve.max_lines: expected a positive"},
		{"a bad pattern", "policy: {approve: {verify: [x], never: [\"a/**/b\"]}}\n", nil, "policy.approve.never: bad pattern a/**/b"},
		{"an escaping pattern", "policy: {approve: {verify: [x], never: [../x]}}\n", nil, "policy.approve.never: bad pattern"},
		{"integrate alone", "policy: {integrate: true}\n", nil, "policy.integrate: needs approve"},
		{"integrate as a word", "policy: {approve: {verify: [x]}, integrate: yes}\n", nil, "policy.integrate: expected true or false"},
		{"an unknown key", "policy: {merge: true}\n", nil, "policy.merge: unknown key"},
		{"an unknown nested key", "policy: {budget: 1, resolve: {model: x}}\n", nil, "policy.resolve.model: unknown key"},
		{"not a mapping", "policy: true\n", nil, "policy: expected a mapping"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, ds := Load(briefFixture(t, tc.config), "")
			if got := diagnostics(ds); tc.want == "" && got != "" || !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
			if tc.want == "" && !reflect.DeepEqual(p.Policy, tc.policy) {
				t.Fatalf("Policy = %+v, want %+v", p.Policy, tc.policy)
			}
		})
	}
	p := &Policy{Never: []string{"grove.yaml", ".github/**", "*.sum"}}
	for file, want := range map[string]string{"grove.yaml": "grove.yaml", ".github/w/x.yml": ".github/**", "go.sum": "*.sum", "a/go.sum": "", ".github": "", "x/grove.yaml": ""} {
		if got, _ := p.Matches(file); got != want {
			t.Errorf("Matches(%q) = %q, want %q", file, got, want)
		}
	}
}

// TestRefusalsNameTheirRule holds every refusal of an enumerated value or an
// unknown key to listing what is accepted, from the same tables validation
// checks, and each lifecycle refusal to stating its rule and the correction,
// so a refused command can be corrected without the record model.
func TestRefusalsNameTheirRule(t *testing.T) {
	t.Parallel()
	const id, commit = "G-260101-00002", "0123456789abcdef0123456789abcdef01234567"
	var names []string
	type refusal struct{ name, source, want string }
	cases := []refusal{
		{"kind", typed(id, "work", "proposed", "kind: chore\n"), "kind: expected " + Choices(Kinds)},
		{"size", typed(id, "work", "proposed", "size: huge\n"), "size: expected " + Choices(Sizes)},
		{"priority", typed(id, "work", "proposed", "priority: 9\n"), "priority: expected 1 (highest) through 5 (lowest)"},
		{"candidate in review", typed(id, "work", "review", ""), "candidate: required while status is review: set candidate=COMMIT"},
		{"approved off review", typed(id, "work", "active", "candidate: \""+commit+"\"\napproved: \""+commit+"\"\n"), "approved: approval holds only while status is accepted or done, not active: unset approved, or set status accepted"},
		{"approved another commit", typed(id, "work", "review", "candidate: \"abcdef0\"\napproved: \""+commit+"\"\n"), "approved: approval is of one commit and must name the candidate abcdef0"},
		{"approved beside a bad status", typed(id, "work", "bogus", "candidate: \""+commit+"\"\napproved: \""+commit+"\"\n"), "status: expected " + Choices(Types[0].Statuses) + " for work"},
		{"approved without candidate", typed(id, "work", "done", "approved: \""+commit+"\"\n"), "approved: approval is of the candidate, and there is none: set candidate"},
	}
	for _, ty := range Types {
		names = append(names, ty.Name)
		status := "status: bogus\n"
		if len(ty.Statuses) == 0 {
			status = ""
		} else {
			cases = append(cases, refusal{ty.Name + " status", "---\nid: " + id + "\ntype: " + ty.Name + "\ntitle: T\n" + status + "---\n", "status: expected " + Choices(ty.Statuses) + " for " + ty.Name})
			status = "status: " + ty.Statuses[0] + "\n"
		}
		cases = append(cases, refusal{ty.Name + " field", "---\nid: " + id + "\ntype: " + ty.Name + "\ntitle: T\n" + status + "zzz: x\n---\n", "zzz: not a field of " + ty.Name + " records, which take " + Choices(ty.Keys())})
	}
	cases = append(cases, refusal{"type", typed(id, "note", "open", ""), "type: unknown record type; expected " + Choices(names)})
	for _, tc := range cases {
		_, ds := ParseRecord("grove/a.md", []byte(tc.source))
		if got := diagnostics(ds); !strings.Contains(got, tc.want) || strings.Count(got, "\n") != 0 {
			t.Errorf("%s: wanted one line holding %q; got %q", tc.name, tc.want, got)
		}
	}
	for _, tc := range []struct{ config, want string }{
		{"extra: 1\n", "extra: unknown configuration key; grove.yaml takes " + Choices(ConfigKeys)},
		{"run:\n  bound: 1\n", "run.bound: unknown key; run: takes " + Choices(RunKeys)},
		{"policy:\n  x: 1\n", "policy.x: unknown key; policy: takes " + Choices(PolicyKeys)},
		{"policy:\n  budget: 1\n  resolve:\n    x: 1\n", "policy.resolve.x: unknown key; policy.resolve takes " + Choices(PolicyResolveKeys)},
		{"policy:\n  approve:\n    verify: [true]\n    x: 1\n", "policy.approve.x: unknown key; policy.approve takes " + Choices(PolicyApproveKeys)},
	} {
		_, _, ds := ParseConfig([]byte("schema_version: 4\nrecords: grove\n" + tc.config))
		if got := diagnostics(ds); !strings.Contains(got, tc.want) {
			t.Errorf("%q: wanted %q; got %q", tc.config, tc.want, got)
		}
	}
}
