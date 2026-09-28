package versions

import (
	"context"
	"fmt"
	"strings"

	"github.com/mascah/grove/internal/repo"
)

// Copies says how much of a branch the target already holds as rewritten
// copies (G-260928-4qv1m): commits with the same patch as the branch's, as
// when the target is rebased, or the branch cherry-picked onto it, after
// the branch was merged. A rewritten branch still diverges from the target by
// ancestry, which the current view keeps; this explains it.
type Copies struct {
	Target, Branch string   // the commits compared
	Commits        int      // the branch's commits the target lacks
	Missing        []string // of those, the ones with no copy on the target
}

// CopiesContext compares branch with target by patch, as git cherry does,
// in one process. A merge has no patch, so it never has a copy.
// ponytail: a merge whose resolution the target does hold still counts as
// missing, which keeps the delete command away; compare its tree if that
// case ever matters.
func CopiesContext(ctx context.Context, root, target, branch string) (*Copies, error) {
	_, commits, err := resolveCommits(ctx, root, target, branch)
	if err != nil {
		return nil, err
	}
	out, err := repo.GitContext(ctx, root, "log", "--cherry-mark", "--right-only", "--no-show-signature", "--format=%m%H", commits[0]+"..."+commits[1])
	if err != nil {
		return nil, err
	}
	c := &Copies{Target: commits[0], Branch: commits[1]}
	for l := range strings.FieldsSeq(out) {
		c.Commits++
		if l[0] != '=' { // + without a copy, > a merge
			c.Missing = append(c.Missing, l[1:])
		}
	}
	return c, nil
}

// Rewritten reports a branch every commit of which the target holds as a
// copy: nothing on it needs merging.
func (c *Copies) Rewritten() bool { return c.Commits > 0 && len(c.Missing) == 0 }

// Text explains the comparison to the owner, naming the branch, the target
// and the branch's checkout ("" when none): how to clear a rewritten branch,
// or which commits are not copies. It is "" when no commit has a copy,
// which is ordinary divergence.
func (c *Copies) Text(branch, target, worktree string) string {
	copied := c.Commits - len(c.Missing)
	if copied == 0 {
		return ""
	}
	if c.Rewritten() {
		clear := "git branch -D " + branch
		if worktree != "" {
			clear = "git worktree remove " + worktree + ", which also deletes that checkout's ignored files such as build output, then " + clear
		}
		return fmt.Sprintf("branch %s is a rewritten copy of work already on %s: each of its %d %s %s lacks has a copy there with the same patch, as after a rebase of %s, so nothing needs merging. To clear it: %s (-D, since Git checks ancestry, not patches, and -d would refuse)",
			branch, target, c.Commits, plural(c.Commits, "commit", "commits"), target, target, clear)
	}
	var missing []string
	for _, m := range c.Missing {
		missing = append(missing, m[:min(len(m), 7)])
	}
	return fmt.Sprintf("%d of the %d commits of branch %s that %s lacks have a copy there with the same patch, as after a rebase of %s; %s not: %s. Merge or judge the branch as usual; it is not a copy to delete",
		copied, c.Commits, branch, target, target, plural(len(missing), "this one has", "these have"), strings.Join(missing, ", "))
}

// CopyOfContext lists base's commits with the same patch as commit, which
// base lacks: its rewritten copies. The right side of the comparison is
// commit alone (commit^! excludes its parents), so a copy of an earlier
// commit on its branch is not one of commit's. A merge has none.
func CopyOfContext(ctx context.Context, root, commit, base string) ([]string, error) {
	_, commits, err := resolveCommits(ctx, root, commit, base)
	if err != nil {
		return nil, err
	}
	out, err := repo.GitContext(ctx, root, "log", "--cherry-mark", "--left-only", "--no-merges", "--no-show-signature", "--format=%m%H", commits[1]+"..."+commits[0], commits[0]+"^!")
	if err != nil {
		return nil, err
	}
	var copies []string
	for l := range strings.FieldsSeq(out) {
		if l[0] == '=' {
			copies = append(copies, l[1:])
		}
	}
	return copies, nil
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
