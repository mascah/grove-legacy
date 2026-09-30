package update

import (
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// TestMigrateSchema3 covers each row of the migration's classification, its
// dry run, the refusals before writing, and the committed result: IDs, paths
// and bodies kept, the way back under refs/grove/schema-3, and a project that
// validates at schema 4.
func TestMigrateSchema3(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	on := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "-b", "feature")
	git(t, root, "commit", "-q", "--allow-empty", "-m", "unmerged")
	off := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	policy := "sha256:" + strings.Repeat("c", 64)
	w := func(id, status, extra, body string) {
		write(t, root, "grove/"+id+".md", "---\nid: \""+id+"\"\ntype: work\ntitle: "+id+"\nstatus: "+status+"\n"+extra+"---\n\n## Outcome\n\nIt works.\n"+body)
	}
	cand := func(c string) string { return "candidate: \"" + c + "\"\n" }
	appr := func(c string) string { return cand(c) + "approved: \"" + c + "\"\n" }
	w("G-260101-00011", "review", cand(on), "")
	w("G-260101-00012", "review", appr(on), "\nVerdict on candidate "+on[:7]+", 2026-01-01: delegated under policy grove.yaml "+policy+": fine\n")
	w("G-260101-00013", "done", appr(on), "\nVerdict on candidate "+on[:7]+", 2026-01-01: Ship it.\n")
	w("G-260101-00014", "done", appr(off), "")
	w("G-260101-00015", "done", cand(on), "")
	w("G-260101-00016", "done", "", "")
	write(t, root, "grove.yaml", "schema_version: 3 # the old one\nrecords: grove\ntarget: main\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "schema 3")
	head := git(t, root, "rev-parse", "HEAD")

	m, err := PlanMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range append(m.Changes, m.Kept...) {
		got = append(got, c.ID+" "+c.From+">"+c.To+" "+c.Why)
	}
	want := []string{
		"G-260101-00012 review>accepted approved " + on[:7] + " by policy " + policy,
		"G-260101-00013 done>accepted approved " + on[:7] + " by owner; main contains it",
		"G-260101-00014 done>done kept as schema 3's claim: main does not contain candidate " + off[:7],
		"G-260101-00015 done>done kept as schema 3's claim: no approval recorded",
		"G-260101-00016 done>done kept as schema 3's claim: no approval recorded",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") || len(m.Problems) != 0 {
		t.Fatalf("plan:\n%s\nproblems %v", strings.Join(got, "\n"), m.Problems)
	}
	if git(t, root, "rev-parse", "HEAD") != head || git(t, root, "status", "--porcelain") != "" {
		t.Fatal("the dry run wrote something")
	}
	write(t, root, "grove/G-260101-00011.md", read(t, root, "grove/G-260101-00011.md")+"\nA local edit.\n")
	if _, _, err := Migrate(root, now); err == nil || !strings.Contains(err.Error(), "uncommitted changes") {
		t.Fatalf("dirty project: %v", err)
	}
	git(t, root, "checkout", "-q", "--", ".")

	m, commit, err := Migrate(root, now)
	if err != nil || commit != git(t, root, "rev-parse", "HEAD") || git(t, root, "rev-parse", "refs/grove/schema-3/main") != head {
		t.Fatalf("migrate: %v %s", err, commit)
	}
	if files := git(t, root, "show", "--format=", "--name-only", "HEAD"); files != "grove.yaml\ngrove/G-260101-00012.md\ngrove/G-260101-00013.md" {
		t.Fatalf("the migration commit must hold what it changed alone: %q", files)
	}
	if config := read(t, root, "grove.yaml"); config != "schema_version: 4 # the old one\nrecords: grove\ntarget: main\n" {
		t.Fatalf("grove.yaml: %q", config)
	}
	for id, by := range map[string]string{"G-260101-00012": "policy " + policy, "G-260101-00013": "owner"} {
		r := record(t, root, id)
		if r.Status != "accepted" || r.ApprovedBy != by || r.ApprovedContext != project.AcceptanceContext(r) || !strings.Contains(string(r.Source), "\n\n## Outcome\n\nIt works.\n") || !strings.HasSuffix(string(r.Source), " became status accepted by "+by+".\n") {
			t.Fatalf("%s:\n%s", id, r.Source)
		}
	}
	if r := record(t, root, "G-260101-00014"); r.Status != "done" || r.Approved != off {
		t.Fatalf("a kept claim changed: %+v", r)
	}
	if _, err := PlanMigration(root); err == nil || !strings.Contains(err.Error(), "unsupported version 4; expected 3") {
		t.Fatalf("a second migration: %v", err)
	}
}
