---
id: "G-260928-4qv1m"
type: work
title: "Say when a diverging state is a rewritten copy already on the target, and how to clear it"
status: done
created: "2026-09-28T16:36:30Z"
updated: "2026-09-28T17:49:50Z"
kind: feature
size: medium
relates_to: ["G-260921-ms6ev", "G-260921-jwk4e", "G-260925-h8rj5", "G-260925-dz10z", "G-260920-svpbc"]
candidate: "b107ba0e0893f9a81f126714971b0a3c7b97254f"
approved: "b107ba0e0893f9a81f126714971b0a3c7b97254f"
---

## Outcome

When a work record diverges only because its branch was rewritten onto the
target (a rebase, squash or cherry-pick of the target after the branch
merged), the owner opening its card, or running `integrate` or `resolve` on
it, learns what happened and the exact commands that clear it, instead of a
card in Review whose detail says Done and whose `a` and `i` refuse.

Owner intent, conversation 2026-09-28, after the ascah.dev incident below:
make "this situation easier to understand what happened and how to resolve
going forward".

## Constraints

**Observed: the ascah.dev incident, 2026-09-28** (installed grove
`2ed99f4`; ascah.dev read from its own checkout).

- `grove integrate` fast-forwarded ascah.dev `main` to
  `worktree-G-260927-5gh2k` at `c136c64` on 2026-09-27 and wrote `done`
  there. The branch and its worktree stayed. The worktree holds ignored
  build output (`dist/`, `node_modules/`, `.astro/`), which `cleanup`
  keeps by design (`internal/integrate/integrate.go:301`).
- On 2026-09-28 at 09:07, `main` was rebased onto `origin/main`, which held
  one commit from May that local `main` lacked. The reflog shows
  `rebase (start): checkout origin/main`. All 11 local commits since then
  got new hashes, and the new `main` was pushed. `git cherry main
  worktree-G-260927-5gh2k` marked all four branch commits `-`, meaning
  each has a patch-equivalent copy on `main`.
- The merge base of the two is `f4ce9d4`, where the record was `proposed`.
  Both sides changed it since, so `grove versions G-260927-5gh2k` showed
  two current states: `done` on `main` and `review`, approved, on the
  branch. That is correct under the current view's rule
  ([G-260921-ms6ev](G-260921-ms6ev-derive-a-project-wide-cu.md)).
- The board filed the card in Review, the earliest status among the states
  (`currentCards`, `internal/tui/model.go:1031`). Opening it showed `main`'s
  `done`, because `shown()` (`internal/tui/model.go:381`) takes the first
  current version. So `a`, `f` and `i` said "is not in review: nothing to
  approve, give feedback on, or integrate" (`internal/tui/review.go:528`).
  That notice does not mention the state in review, or where it is.
- Run from the branch's state, `integrate` would find the approval and
  predict a conflict on the record file (`review` against `done`;
  `git merge-tree` reproduced it). Its refusal points to `grove resolve`,
  which would start an agent attempt to merge `main` into a branch whose
  work `main` already holds.
- The version detail's divergence text (`currentText`,
  `internal/tui/view.go:708`) says the card stays put "until one side takes
  the other's change, by a merge or an edit". Neither applies here; the fix
  is deleting the branch.
- `main`'s `done` record still named `candidate: d799309cf137`, the
  pre-rebase hash. After the branch was deleted, launching
  G-260927-kkgke (which `depends_on` it) refused: "G-260927-kkgke needs
  G-260927-5gh2k, whose candidate d799309 the base lacks; resolve that
  before another attempt" (`internal/attempt/selection.go:115`). Neither
  the refusal nor `grove deps` says why the base lacks it or how to repair
  it. The record model's existing rule applies: "A squash or rebase that
  lands a different commit is a manual merge that names that commit as the
  candidate" (`docs/record-model.md`). It was repaired by hand in ascah.dev
  `5120f11`: `43b74ac` has the same `git patch-id --stable` as `d799309`,
  so `grove update G-260927-5gh2k --set candidate=43b74acd29b5 --set
  approved=43b74acd29b5 --commit` with a note under the verdict. After
  that, `grove run G-260927-kkgke --dry-run` reported "delivered" and
  "can start".

**Owner decision, 2026-09-28 (this conversation).** Keep the divergence:
the current view stays Git ancestry, as the brief and G-260921-ms6ev
select. A rewritten copy is still a current state. Grove explains it; it
does not mark it older.

**Repository constraints that bind the design** (this repository's
`CLAUDE.md`):

- Read a record's Git history only while its card is open, as a read any
  key may cancel, never during the board load
  ([G-260920-svpbc](G-260920-svpbc-show-a-work-item-s-linea.md)). A
  patch-equivalence read is the same kind of cost, so it follows the same
  rule.
- Every Git process goes through `repo.Command`
  ([G-260922-g6e7p](G-260922-g6e7p-ignore-ambient-git-envir.md)).

**Proposed design** (it binds nobody; preparation settles the mechanics):

- A function in `internal/versions` that, given a target commit and a
  branch commit, reports how many of the branch's non-merge commits missing
  from the target have a patch-equivalent copy on it. `git cherry` or
  `git patch-id` both give this.
- Board detail, only in the current view and only for a diverging card
  where one state is on the target and another is a committed branch. On
  open, a cancellable read like history. When every branch commit has a
  copy, the detail says that branch *B*'s state is a rewritten copy of work
  already on *target* (for example after a rebase of *target*), that
  nothing needs merging, and gives the commands to clear it:
  `git worktree remove PATH` when a checkout is on *B*, then
  `git branch -D B`, with `-D` explained. When only some commits have a
  copy, it says how many and names the others, and gives no delete
  command.
- On a diverging card, `a`, `f`, `i` and `m` refuse by naming the state that
  is in review and its branch, and point to `v`, rather than saying the
  record is not in review.
- `integrate` and `resolve`, before predicting a conflict, refuse a branch
  whose every commit already has a copy on the target, with the same
  explanation and commands. Nothing is written.
- Where a done prerequisite's candidate is missing from the base (launch,
  `grove run --dry-run`, `grove deps`), look for the candidate's
  patch-equivalent copy in the base. When exactly one exists, say the base
  holds it as *Y*, a rewritten copy, and give the repair:
  `grove update ID --set candidate=Y --set approved=Y --commit` in the
  target's checkout, with a note under the verdict. It still waits: a copy
  counts as delivered only once the record names it, so the record stays
  the fact.
- Out of scope: marking such a copy older or placing the card by it; a
  board key that deletes branches or worktrees (Grove prints the commands,
  the owner runs them); rewriting a done record's candidate automatically;
  a `grove versions` note, which would cost a read per diverging branch
  on every listing.

## Acceptance

1. In a disposable repository reproducing the incident (integrate by
   fast-forward, keep the branch, rebase the target onto a new commit):
   - the board detail for the record names the branch as a rewritten copy
     already on the target and prints the removal commands;
   - `a`, `f`, `i` and `m` name the review state's branch;
   - `grove integrate ID` and `grove resolve ID` refuse with the same
     explanation and write nothing, and `git status` is unchanged.
2. After running the printed commands, the card has one state, Done, with
   no `⑂` tag.
3. In the same repository, after the branch is deleted, launching work
   that `depends_on` the record waits with a message naming the rewritten
   copy and the exact `update` command. After that command, the launch
   reports the prerequisite delivered.
4. A branch with one commit that has no copy on the target gets the partial
   message and no delete command, and `integrate` behaves as it does today.
5. The equivalence read never runs during the board load, and a key press
   cancels it. A test proves both, as the history tests do.
6. `docs/board.md` (the divergence text and the actions),
   `docs/commands.md` (`integrate`, `resolve`, `run`, `deps`) and the
   version detail's divergence text describe the rewritten case and how to clear it.
   `docs/work-execution.md` says that rewriting the target after an
   integration leaves the integrated branch diverging, and what to do.
7. The owner judges the detail wording in a real terminal against the
   ascah.dev shape.

## Evidence

Branch `worktree-G-260928-4qv1m`, base main `ac43184`. Started from this
record at `sha256:1e00c3cc…` and plan
[G-260928-zyqn9](G-260928-zyqn9-rewritten-copy-plan.md) at
`sha256:1673f192…` (commit `9e85b76`); implementation `fc1e6a9`
(integrate and resolve), `8e5eea5` (board and prerequisites), `ff18761`
(documents), `905ece2` (review fixes, recorded in the plan's "Adjusted
after review").

**How it works.** `versions.CopiesContext` compares a branch with the
target by patch in one `git log --cherry-mark --right-only TARGET...BRANCH`;
a merge never has a copy. `versions.CopyOfContext` finds a commit's own
copies with `git log --cherry-mark --left-only --no-merges BASE...C C^!`.
`Copies.Text` is the one explanation every surface prints. Everything goes
through `repo.GitContext`. The explanation applies only to work the target
already holds as `done`, so a branch cherry-picked by hand while its record
is still in review on the target is integrated as before.

1. **Met.** Reproduced with a built binary in a disposable repository
   (fast-forward integrate, branch and worktree kept, main rebased onto a
   new commit; `versions` shows `done` on main, `review` on the branch):
   - `grove integrate ID` refused with "merge of BRANCH into main refused:
     branch BRANCH is a rewritten copy of work already on main: each of its
     4 commits main lacks has a copy there with the same patch, as after a
     rebase of main, so nothing needs merging. To clear it: git worktree
     remove PATH, which also deletes that checkout's ignored files such as
     build output, then git branch -D BRANCH (-D, since Git checks ancestry,
     not patches, and -d would refuse); nothing was merged, main is
     unchanged at …".
   - `grove resolve ID` refused with the same text and "; there is nothing
     to resolve".
   - HEAD, refs and `git status` were unchanged in both checkouts.
   - The board, driven at `905ece2` with the live backend, shows that text
     under the detail's Sources and in the versions' divergence text; `a`,
     `f`, `i` and `m` say "ID is done on branch main and in review on branch
     BRANCH: its states diverge, and v shows both and how to settle them".
   - Tests: `TestCopiesAfterARebaseOfTheTarget`,
     `TestIntegrateRefusesARewrittenCopy`, `TestResolveRefusesARewrittenCopy`,
     `TestRewrittenCopyExplainedOnOpen`.
2. **Met** (at `8e5eea5`, and `versions` again at `905ece2` in review
   round 2). After running the printed `git worktree remove` (it deleted
   the ignored build output) and `git branch -D`, `versions` lists only
   `done` and the board shows one Done card, no `⑂`.
3. **Met.** At `905ece2` (review round 2; my own run at `8e5eea5` agreed
   but printed the earlier `--set approved`), `grove run DEPENDENT
   --dry-run` refused: "needs ID, whose candidate OLD the base lacks; the
   base holds COPY, a rewritten copy with the same patch, as after a
   rebase: in the target's checkout, grove update ID --set candidate=COPY
   --unset approved --commit names it as delivered; then note the rewrite
   under the record's verdict", and `grove deps DEPENDENT` gave the same as
   a note. After running that update, `check` was OK and the dry-run said
   "delivered: candidate COPY is in the base … can start". Tests:
   `TestSelectionNamesARewrittenPrerequisite` (update applied, then
   delivered) and `TestDeliverExplainsEachStatus`.
4. **Met, from tests with real Git.** One extra commit gives "2 of the 3
   commits of branch work that main lacks have a copy there …; this one has
   not: …", with no delete command (`TestCopiesAfterARebaseOfTheTarget`,
   and on the board `TestRewrittenCopyExplainedOnOpen`). `integrate` then
   reaches today's conflict refusal naming `grove resolve`
   (`TestIntegrateRefusesARewrittenCopy`, second half).
5. **Met.** `wantCopies` runs only with the detail or versions screen open,
   through `m.read` like history. `TestRewrittenCopyExplainedOnOpen` asserts
   no read on the board load; `TestRewrittenCopyReadYieldsToEveryKey`
   asserts that Esc, `r` and `q` cancel it and ignore its reply. The
   dependency preview's copy lookup runs only inside the explicit `p`
   preview read.
6. **Met.** `docs/board.md` has the divergence sentence, a new "Rewritten
   copies" section and the key table. `docs/commands.md` covers
   `integrate`, `resolve`'s refusals, `run`'s waits and `deps`. The
   versions screen's divergence text (`currentText`) names deletion as a
   way it ends. `docs/work-execution.md` has "A target rewritten after
   integration", naming only its example ID.
7. **Pending**: the owner's judgment of the wording in a real terminal.

**Decision taken, for the owner.** This record's proposed design, and the
hand repair in ascah.dev, set `approved` to the copy. The settled
[Approval](G-260921-btyck-approval.md) term says an approval does not carry
over to another commit, so the printed repair is `--unset approved` where
the record is approved; the verdict stays in the body. If an approval
should carry over to a patch-equivalent copy, that amends the term and
changes one line in `deps.Rewrite`.

Verification at `905ece2`: `go vet ./...` clean, `gofmt -l .` empty,
`go run ./cmd/grove check` `OK: 230 records`, `go test -count=1 -timeout
120s ./...` all packages ok, and `python3 internal/tui/testdata/terminal.py
BINARY` all 12 scenarios ok. No `-race`: the change adds a read on the
existing one-at-a-time read path and no concurrency.

Review: [G-260928-0pbpf](G-260928-0pbpf-rewritten-copy-review.md), two
rounds by fresh `grove-reviewer` agents. Round 1 found three, all fixed in
`905ece2`; round 2 `Open findings: none`.

Limits:
- A merge commit on the branch never counts as a copy. A branch that went
  through `grove resolve`, was integrated and then had its target rebased
  gets the partial message, and `integrate` still points to `resolve`
  (`ponytail:` note in `internal/versions/copies.go`).
- The repair is worded for the target's checkout, but `run` and `deps`
  look for the copy in their base or HEAD. From a base other than the
  target, `update` refuses a copy the target lacks. Not tested.
- `internal/versions` takes about 5 s under `-short`, with or without its
  new test; that is the existing baseline.

## Next

Judge the candidate this record names, and acceptance 7 in a terminal
against the ascah.dev shape. In this checkout:
`grove approve G-260928-4qv1m "VERDICT"`; then, in main's checkout,
`grove integrate G-260928-4qv1m`.

Verdict on candidate b107ba0, 2026-09-28: approved
