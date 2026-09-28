package versions

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// TestCopiesAfterARebaseOfTheTarget reproduces G-260928-4qv1m's incident: a
// branch fast-forwarded into main, kept, and main rebased onto a commit it
// lacked. Every branch commit has a copy; one more commit does not; a merge
// never does; and each candidate finds exactly its own copy.
func TestCopiesAfterARebaseOfTheTarget(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	ctx := context.Background()
	base := git(t, root, "rev-parse", "HEAD")
	wt := addWorktree(t, root, "work", "main", "-b", "work")
	write(t, wt, "a.txt", "a\n")
	first := commit(t, wt, "a")
	write(t, wt, "b.txt", "b\n")
	second := commit(t, wt, "b")
	git(t, root, "merge", "-q", "--ff-only", "work")
	git(t, root, "branch", "upstream", base)
	up := addWorktree(t, root, "up", "upstream")
	write(t, up, "u.txt", "u\n")
	commit(t, up, "upstream")
	git(t, root, "rebase", "-q", "upstream")

	c, err := CopiesContext(ctx, root, "main", "work")
	if err != nil || c.Commits != 2 || len(c.Missing) != 0 || !c.Rewritten() || c.Branch != second {
		t.Fatalf("rewritten: %+v %v", c, err)
	}
	text := c.Text("work", "main", wt)
	for _, want := range []string{"branch work is a rewritten copy of work already on main", "each of its 2 commits main lacks", "nothing needs merging",
		"git worktree remove " + wt + ", which also deletes", "then git branch -D work (-D, since"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text lacks %q: %s", want, text)
		}
	}
	if text := c.Text("work", "main", ""); strings.Contains(text, "worktree") || !strings.Contains(text, "To clear it: git branch -D work") {
		t.Fatalf("without a checkout: %s", text)
	}
	for _, commit := range []string{first, second} {
		copies, err := CopyOfContext(ctx, root, commit[:7], "main")
		if err != nil || len(copies) != 1 || copies[0] == commit {
			t.Fatalf("copy of %s: %v %v", commit, copies, err)
		}
		if got := git(t, root, "log", "-1", "--format=%s", copies[0]); got != git(t, root, "log", "-1", "--format=%s", commit) {
			t.Fatalf("copy of %s is %s", commit, got)
		}
	}
	if copies, err := CopyOfContext(ctx, root, base, "upstream"); err != nil || copies != nil {
		t.Fatalf("a commit the base holds has no copy: %v %v", copies, err)
	}

	// One commit more, which main lacks: the partial answer names it.
	write(t, wt, "c.txt", "c\n")
	third := commit(t, wt, "c")
	c, err = CopiesContext(ctx, root, "main", "work")
	if err != nil || c.Commits != 3 || !reflect.DeepEqual(c.Missing, []string{third}) || c.Rewritten() {
		t.Fatalf("partial: %+v %v", c, err)
	}
	if text := c.Text("work", "main", wt); !strings.Contains(text, "2 of the 3 commits of branch work that main lacks have a copy") ||
		!strings.Contains(text, "this one has not: "+third[:7]) || strings.Contains(text, "git branch") || strings.Contains(text, "remove") {
		t.Fatalf("partial text: %s", text)
	}

	// A merge has no patch, so it has no copy.
	git(t, wt, "reset", "-q", "--hard", second)
	git(t, wt, "merge", "-q", "--no-ff", "--no-edit", "upstream")
	c, err = CopiesContext(ctx, root, "main", "work")
	if err != nil || c.Commits != 3 || len(c.Missing) != 1 || c.Missing[0] != git(t, wt, "rev-parse", "HEAD") {
		t.Fatalf("merge: %+v %v", c, err)
	}

	// Ordinary divergence has nothing to say.
	c, err = CopiesContext(ctx, root, "upstream", "work")
	if err != nil || c.Text("work", "upstream", "") != "" {
		t.Fatalf("no copies: %+v %v", c, err)
	}
}
