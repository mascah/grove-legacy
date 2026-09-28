---
id: "G-260928-zyqn9"
type: plan
title: "Rewritten copies: one equivalence read for the board, integrate, resolve and prerequisites"
status: current
created: "2026-09-28T17:07:39Z"
updated: "2026-09-28T17:08:14Z"
work: ["G-260928-4qv1m"]
---

# Plan for G-260928-4qv1m

Base: `main` at `ac43184`, branch `worktree-G-260928-4qv1m`, one
implementer. No prerequisite, no blocking question.

## Choice: `git log --cherry-mark`, Git's own patch-id engine

Both questions are answered by `git log --cherry-mark`, the engine behind
`git cherry` (patch-ids, whitespace-insensitive diff hashes), in one process
each, through `repo.GitContext`:

- **Is a branch a rewritten copy?**
  `git log --cherry-mark --right-only --format=%m%H TARGET...BRANCH` lists
  each commit of the branch that the target lacks, marked `=` when a
  patch-equivalent commit is on the target's side, `+` when not, and `>` for
  a merge, which has no patch-id. Tried on Git 2.55: after cherry-picking
  both commits of a branch that also holds a merge, the two are `=`, the
  merge `>`. A merge counts as *without a copy*: its resolution may hold
  changes nothing else carries, so it keeps the delete command away (the
  conservative side; `ponytail:` note naming the ceiling).
- **Which base commit is a candidate's copy?**
  `git log --cherry-mark --left-only --no-merges --format=%m%H BASE...C C^!`
  restricts the right side to C alone (`C^!` excludes its parents), so each
  `=` on the base's side is a copy of C and of nothing else. Tried: with a
  two-commit branch rewritten onto a moved target, the query for each
  commit marks exactly its own copy. A merge candidate has no copy (the
  right side is empty after `--no-merges`).

`git cherry` itself skips merges silently and cannot pair a copy with one
commit, and a `git patch-id` pipeline needs two processes per side and the
whole diff in memory; `--cherry-mark` does both in one streaming process.

## Design

### 1. `internal/versions/copies.go`, the shared read and its words

```go
type Copies struct {
	Target, Branch string   // the commits compared
	Commits        int      // the branch's commits the target lacks
	Missing        []string // of those, the ones with no copy on the target
}
func CopiesContext(ctx, root, target, branch string) (*Copies, error)
func (c *Copies) Rewritten() bool   // Commits > 0 && no Missing
func (c *Copies) Text(branch, target, worktree string) string
func CopyOfContext(ctx, root, commit, base string) ([]string, error)
```

`Text` is the one explanation every surface prints: `""` when no commit has
a copy (ordinary divergence says nothing new); when all do, that branch
*B* is a rewritten copy of work already on *target* (for example after a
rebase of *target*), that nothing needs merging, and the commands
`git worktree remove PATH` (when a checkout is on *B*, saying it also
deletes that checkout's ignored files such as build output) then
`git branch -D B`, with `-D` explained (Git sees ancestry, not patches, so
`-d` refuses); when only some do, "K of the N commits … have a copy; these
do not: …", and no delete command.

A repair helper in `internal/deps`, `Rewrite(id, copy string, approved
bool) string`, words the prerequisite repair once:
`grove update ID --set candidate=Y [--set approved=Y] --commit` in the
target's checkout, then a note under the verdict. `approved` is set only
when the record carries one, since the model requires it to equal
`candidate` and an unapproved record must not gain an approval.

### 2. Board (`internal/tui`)

- `Backend.Copies func(ctx, root, target, branch string) (*versions.Copies, error)`,
  wired to `versions.CopiesContext` in `run.go`; nil leaves it out.
- `wantCopies`, beside `wantHistory`: only in the current view, only with
  the detail or versions screen open, only for a group with more than one
  current state where one state is on the target and another is held by a
  committed branch other than the target. It reads one (target commit,
  branch commit) pair at a time through `m.read("copies", …)`, keyed and
  held in `m.copies` like `m.hist` and forgotten on every inspection. The
  board load never asks; any key that starts another read, Esc, `r` and `q`
  cancel it, as with history.
- The sidebar's Sources section and the versions screen's divergence text
  (`currentText`) gain the `Copies.Text` for each such branch, "checking …"
  while reading, or the read's failure. The divergence sentence changes to
  say a rewritten copy also diverges and how it is cleared.
- `a`, `f`, `i` and `m` on a card whose shown record is not in review but
  another current state is: the notice names that state's status and branch
  and points to `v` (`G-… is done on main and in review on branch B; v shows
  both states and how to settle them`).

### 3. `integrate` and `resolve`

- `integrate`: after the approval facts and before the conflict prediction,
  `CopiesContext(target HEAD, branch tip)`; when `Rewritten`, refuse:
  "merge of B into T refused: <Text>; nothing was merged and T is
  unchanged". A partial or failed read changes nothing (today's path).
- `attempt.Resolve`: after `inReview`, the same read against
  `refs/heads/TARGET` and the branch tip, refusing with the same Text and
  "there is nothing to resolve", before any feedback is written.

### 4. Prerequisites (launch, `run --dry-run`, `deps`)

- `attempt`: `containsIn` returns `func(commit) (in bool, copies []string, err error)`,
  the copies read only when the base lacks the commit; `selectionOf` and
  `onBranch` take that type. A done prerequisite whose candidate the base
  lacks and has exactly one copy keeps its wait, extended with the copy and
  `deps.Rewrite`; its Delivery says the base holds a rewritten copy.
- `deps.Deliver` gains `copies func(commit, ref string) ([]string, error)`
  (nil reads none); for a done item whose candidate HEAD lacks, exactly one
  copy in HEAD adds it to the Delivery and a Note with `deps.Rewrite`.
  `deps.Copies(ctx, root)` wraps `versions.CopyOfContext`; the CLI passes
  it, the board's dependency view through `Backend.CopyOf`.

### 5. Documents

`docs/board.md` (divergence and the actions), `docs/commands.md`
(`integrate`, `resolve`, `run`, `deps`), `docs/work-execution.md` (rewriting
the target after an integration leaves the integrated branch diverging, and
what to do). The work guide is shipped: it names no record.

## Steps

1. `versions/copies.go` with a test in a real temporary repository: a
   branch integrated by fast-forward and the target rebased onto a new
   commit (every commit copied, text with both commands, and without a
   worktree only `git branch -D`), one extra uncopied commit (partial, names
   it, no delete command), a merge on the branch (counted as missing), and
   `CopyOfContext` finding exactly the rewritten candidate.
2. `integrate` and `resolve` refusals, tested on the same shape: nothing
   written, `git status` and HEAD unchanged; a partial branch still reaches
   the prediction as today.
3. Prerequisites: `selectionOf` and `Deliver` with fakes for the copies
   function; one real-repository test of `grove run --dry-run` waiting with
   the `update` command and, after running it, reporting delivered.
4. Board: `wantCopies` with the fake backend: not read on the board load,
   read on opening the diverging card, cancelled by a key (as
   `TestHistoryReadYieldsToEveryKey`), the rewritten and partial texts in
   the sidebar and the versions screen, and `a`/`f`/`i`/`m` naming the
   branch.
5. Documents; then `go vet ./...`, `gofmt -l .`, `grove check`, the suite,
   and the terminal script since the board changes. An end-to-end check in
   a disposable repository reproducing the incident (acceptance 1–4), the
   commands printed then run, and the card left with one Done state.

Acceptance 7 (the owner judging the wording in a real terminal) is the
owner's, recorded as pending in the handoff.

## Adjusted after review

Round 1 of the review ([G-260928-0pbpf](G-260928-0pbpf-rewritten-copy-review.md))
changed three mechanics, in `905ece2`:

- The refusal in `integrate` and `resolve`, and the board's comparison,
  apply only to work the target already holds as `done`, the incident's
  shape. A branch landed by a hand cherry-pick while the target still holds
  the record in review is integrated as before, which marks it done.
- The prerequisite repair unsets `approved` instead of setting it to the
  copy: an approval is of one commit and does not carry over
  ([Approval](G-260921-btyck-approval.md)). The verdict stays in the body.
- `resolve` explains a rewritten copy before asking for a checkout on the
  branch.
