---
id: "G-260925-h8rj5"
type: work
title: "Predict whether a candidate in review merges cleanly into the target"
status: done
created: "2026-09-25T21:39:27Z"
updated: "2026-09-25T23:30:44Z"
kind: feature
relates_to: ["G-260925-g39ga", "G-260925-7c8g9", "G-260925-80w3a", "G-260921-jwk4e", "G-260921-3qgsf", "G-260921-jatts", "G-260920-svpbc", "G-260925-dz10z", "G-260925-5wrn8"]
candidate: "a658921371a89424c0482a22d387fef9c17bc2e1"
approved: "a658921371a89424c0482a22d387fef9c17bc2e1"
---

## Outcome

Before pressing `i`, the owner can see for each candidate in review whether
it merges cleanly into the target at its current tip or conflicts in named
files, and, for several candidates in a stated order, where the first
conflict lands, so an integration order or a resolution can be chosen before
any merge is attempted.

Owner intent, conversation 2026-09-25: four or five items implemented in
parallel end in review together, and their conflicts surface only when
`integrate` refuses one, which starts a manual feedback and relaunch loop
the owner wants to avoid. Prediction is the read-only first layer of that
wish; [G-260925-dz10z](G-260925-dz10z-update-a-conflicting-can.md) acts on it by hand and
[G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md) by policy.

## Constraints

### Observed evidence

At main `47852e3`, [`integrate`](../internal/integrate/integrate.go) merges
the inspected commit with a plain `git merge`, aborts on conflict, and
refuses with the target unchanged and the record left in review. Reproduced
on a disposable repository with the built binary:

```text
merge of worktree-G-260919-4h6pn into main refused: CONFLICT (content): Merge conflict in shared.txt; Automatic merge failed; fix conflicts and then commit the result.; main is unchanged at 762e74b and G-260919-4h6pn stays in review
```

The refusal names the file and no next action, and `integrate` runs no
verification after a merge, so two candidates that each merge cleanly and
fail together are not detected by it.

`git merge-tree --write-tree` (Git 2.38 and later; this machine has 2.55.0)
performs the merge read-only, writes only objects, names the conflicting
files with `--name-only`, and exits 1 on conflict. On the same repository
after G-260919-rt9h9 was integrated it printed `shared.txt` and `CONFLICT (content):
Merge conflict in shared.txt` with no checkout touched. Nothing in Grove
uses it.

The board's Review block ([board](../docs/board.md)) shows whether the
target holds the candidate and lists its changed files against the target.
The G-260925-g39ga candidate on `worktree-G-260925-g39ga` (in review at `83e7f38`) adds `grove
deps` and the board's dependency view, whose delivery text reads `awaiting
review; candidate X not in HEAD, not on main`, computed in
`internal/deps/deps.go` (`Deliver`) through one `merge-base --is-ancestor`
per candidate on demand. The predecessor's close skill
(`../skills/skills/close/SKILL.md`, read as a file) left conflicts to the
person: "If the merge conflicts, resolve or rebase, rerun … then retry."

### Proposed design and scope

- One read-only fact per candidate in review, relative to the target's
  current tip: integrated, merges cleanly, merges cleanly although the
  target moved since the branch's merge base, or conflicts in named files.
  Compute it with `git merge-tree --write-tree` through `repo.Command`, on
  demand for the candidates a view names, never during the board load
  (G-260920-svpbc, G-260920-z8vfp). The fact names the target commit it was computed against
  and is stale, visibly, when the target moves.
- Show it where candidates are already explained: the board's Review block,
  `integrate`'s refusal (which then names the next action: G-260925-dz10z's operation
  or the manual merge), and, once G-260925-g39ga is integrated, the `deps` delivery
  text and the selection preview. Choose the noninteractive contract in
  preparation; extending `versions --json` or `deps --json` is preferred to
  a new command.
- For several candidates and an order the person or `context` states,
  simulate that order: each merge-tree result becomes a throwaway commit
  object (`git commit-tree`, no ref, no checkout) that the next merge reads,
  and the report says where the first conflict lands and in which files.
  Grove neither chooses the order nor calls a clean order safe: it is not
  evidence that the changes are semantically compatible.
- Out of scope: running verification on a predicted merge, which needs a
  checkout and belongs to G-260925-5wrn8; starting anything; changing any record or
  ref; a same-branch shared candidate for several IDs, which
  [G-260925-80w3a](G-260925-80w3a-where-should-review-and.md) governs.

## Acceptance

1. For a candidate in review the owner sees, in the board's Review block and
   noninteractively, that it is integrated, merges cleanly into the target
   at a named commit, or conflicts in named files. The same source gives the
   same answer in both.
2. A refused `integrate` names the conflicting files and the exact next
   action, and the target is unchanged.
3. With several candidates and a stated order, the report names where the
   first conflict lands and in which files, and says that a clean order is
   not evidence of compatibility. Grove never chooses the order.
4. Nothing is checked out, written to a ref, or changed in a record; reads
   stay one process per candidate on demand; a target that moved after the
   fact was computed is shown as such.
5. Exercised on: fast-forward, clean merge after the target moved, textual
   conflict, already integrated, and a target that moved between prediction
   and integration.

## Evidence

Implemented headless on `worktree-G-260925-h8rj5` from main `6b14141`, where G-260925-g39ga
is integrated, so `deps` and the board's preview are in scope. Started from
this record at `sha256:fc7cf59f…` and plan
[G-260925-gzdsd](G-260925-gzdsd-plan-for-merge-predictio.md) at `sha256:828cc538…`
(commit `370df4a`). Code through `829f08c`; the candidate adds this
evidence only.

- **One source.** [`internal/versions/merge.go`](../internal/versions/merge.go):
  `Merge {target, commit, outcome, conflicts}` and `Text`, the one wording
  every surface prints. `PredictContext` resolves the target and the
  commits (abbreviated candidates included) in one `rev-parse`, answers
  `integrated` and `fast-forward` from the merge base, and otherwise runs
  `git merge-tree --write-tree --name-only --no-messages -z` (exit 1 is a
  conflict). Several commits are merged in the order given through
  unreferenced `commit-tree` objects, stopping after the first conflict.
  Everything runs through `repo.GitContext`/`repo.Command`.
- **Acceptance 1.** `ChangesContext` reads every fact against one resolved
  `refs/heads/TARGET` commit and fills `Changes.Merge`; the board's Review
  block prints `Merge.Text` in place of `on/not on TARGET`. `deps` text and
  `--json` (`items[].merge`) print the same through `Deliver`'s predict
  function, as does the board's preview (`Backend.Predict`). `TestPredictMerges`
  checks `ChangesContext` and `PredictContext` agree.
- **Acceptance 2.** `integrate` predicts the approved tip against HEAD
  before merging and refuses a conflict naming the target commit, the files
  and the next action: `git merge TARGET` in the branch's checkout, resolve
  and hand the commit back to review, or a pasteable
  `grove feedback ID '…'`. Checked end to end with the built binary; the
  pasted feedback command exits 0. Where prediction fails (Git before 2.38),
  the old `git merge` abort and refusal applies.
- **Acceptance 3.** A `deps` selection with two or more candidates in
  review not on the target gets `merge_order` and a note naming each step up
  to the first conflict, its files, what was not tried, and that Grove chose
  no order and a clean order is not evidence of compatibility
  (`TestDepsPredictsMergeOrder` on a real repository: reversing the order
  moves the conflict; `TestDeliverMergesInTheSelectionsOrder`; the board
  preview test).
- **Acceptance 4.** Tests compare `for-each-ref`, HEAD and `status` before
  and after every prediction and refused integration; nothing writes a
  record. Predictions run only when a card is open, a preview is read or
  `deps` runs; never in the board load. The fact names its target commit;
  the Review block notes when the board read the target at another commit
  (`TestReviewMergePredictionNamesAMovedTarget`), and `integrate` predicts
  afresh.
- **Acceptance 5.** Fast-forward, clean after the target moved, textual
  conflict, integrated, a fileless split-rename conflict and an order
  conflicting only through its first merge in `TestPredictMerges`; a target
  moved between prediction and integration in `integrate`'s "a file
  conflict"; the `git merge` fallback in "what prediction cannot see".
- **Docs.** `docs/board.md`, `docs/commands.md` (`deps` text and JSON),
  `docs/record-model.md`, `docs/work-execution.md`; the term
  [G-260921-3qgsf](G-260921-3qgsf-integration.md) names the pre-merge refusal.
- **Verification at `829f08c`:** `gofmt -l .` empty, `go vet ./...` clean,
  `go run ./cmd/grove check` OK: 177 records,
  `go test -count=1 -timeout 120s ./...` all ok,
  `python3 internal/tui/testdata/terminal.py BINARY` all ok (at `502b27f`,
  before a change to integrate's error text and docs only).
  `internal/versions` under `-short` alone: 4.3 to 4.6s against about 4.2s
  at the base, under the five-second limit but close.
- **Review:** [G-260925-a05hb](G-260925-a05hb-review-of-candidate-merg.md), three
  rounds by fresh `grove-reviewer` agents; round 3 found nothing
  consequential.

Limits for the owner: a prediction takes three Git processes per candidate
(four per clean step of an order), bounded and on demand, which acceptance
4's "one process per candidate" may or may not mean; the board and `deps`
predict the candidate while `integrate` predicts the tip it merges, which
differ only if a record-only commit after the candidate conflicts; the
Git-before-2.38 fallthrough was verified by a reviewer's shim run, not an
automated test; a target branch name containing `'` would break the
printed feedback command.

## Next

In review. The owner judges the candidate on `worktree-G-260925-h8rj5`:

```sh
grove approve G-260925-h8rj5 "VERDICT"   # in /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-260925-h8rj5
grove integrate G-260925-h8rj5           # then in the main checkout
```

or `grove feedback G-260925-h8rj5 "TEXT"` in the worktree. A demo: `grove deps
G-260925-h8rj5 G-260925-dz10z` or the board's `g`, `p` on candidates in review, and a card's
Review block. G-260925-dz10z and G-260925-5wrn8 act on this fact.

Verdict on candidate a658921, 2026-09-25: approved
