---
id: "G-260927-dx0yn"
type: work
title: "Retain per-attempt process facts and show cost per work record"
status: review
created: "2026-09-27T22:13:36Z"
updated: "2026-09-28T16:37:15Z"
size: medium
relates_to: ["G-260921-sth8q", "G-260923-tnn5e", "G-260923-895zb", "G-260921-h46pb", "G-260923-p5pt6"]
candidate: "f1f75b2"
approved: "f1f75b2"
---

## Outcome

Each attempt Grove ran tells the owner what it spent its turns on and what
it changed, and each work record shows what its attempts have cost
together, so the effect of a workflow change can be read from the attempts
Grove already keeps instead of from transcripts by hand.

Owner intent, conversation 2026-09-27: "we should definitely instrument
attempts further as you suggested": process share, time to first code
edit, files changed outside the record root, and cost per record. The
facts and their placement are proposed below.

## Constraints

Observed 2026-09-27 at main `fef102f`:

- `grove attempt ATTEMPT` prints the launch (`attempt.json`: model, effort,
  budget, permission mode, command, worktree, base, entrypoint digests),
  the result event (turns, USD, denials), cost by model, event counts, the
  worktree after (HEAD, clean or dirty) and the record on the branch.
  `result.json` holds `dirty`, `events`, `exit_code`, `finished`, `head`,
  `record` and `stopped`. `Metrics` (`internal/attempt/activity.go:54-66`)
  already derives turns, input and output tokens, context, subagents,
  compactions, tool calls and tool errors from a bounded window of
  `events.jsonl` for the board's attempts screen.
- Not retained or derived: which tool calls were process (a `grove`
  command, worktree or branch setup, an edit under the record root, a guide
  print) and which were the task; the time and tool index of the first edit
  outside the record root; the files changed between base and HEAD, inside
  and outside the record root; guide prints by name.
- `grove attempts` lists a COST column per attempt and no total per work.
  The 43 finished attempts here total $250.93, mean $5.84, range $0.64 to
  $16.15 (each attempt's result event, read 2026-09-27); several records
  had two to four attempts.
- Measured by hand from this repository's 48 interactive work sessions:
  about a quarter of tool calls were process by that classification, and in
  the 19 sessions that edited a file outside `grove/` the first such edit
  came after a median 25 tool calls and 4 minutes. Interactive sessions are
  outside Grove's reach; attempts are not.
- `docs/commands.md` "Attempts" documents what an attempt retains, and the
  [Attempt](G-260921-sth8q-attempt.md) term points there. Reads of a running
  attempt are bounded (`ReadActivity`), a rule
  [G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) kept.

**Proposed design.** Derive at read time from `events.jsonl`, as `Metrics`
does: process tool calls and their share (Bash whose program is `grove`,
`git worktree` or a branch created under the worktree prefix, Edit or
Write under the record root, a Read of a record or guide), guide prints by
name, and the elapsed time and tool index of the first edit outside the
record root. At finish, `result.json` gains the files changed between base
and HEAD, inside and outside the record root, since the worktree may be
removed later. `grove attempt` prints these under one `Shape:` line; `grove
attempts ID` adds a total line (attempts, USD, turns, minutes); the board's
attempts pane shows the same total on the record; `--json` carries all of
it. For a running attempt the window's facts are lower bounds, labelled as
`Metrics` labels its own. A `grove` run through a wrapper or `$(…)` is
missed, as the evals runner already states of its own reader.

The owner's feedback on `8cf2e4c` (Next) dropped the process share and the
first edit as noise, since most edits and record reads go through Bash, and
added the time to the first commit outside the record root, from Git at
finish.

## Acceptance

1. `grove attempt` shows for a finished attempt: tool calls, calls running
   `grove`, guide prints by name, the files changed inside and outside the
   record root, and the time to the first commit outside it; an unknown
   fact is labelled unknown, never filled. (Revised by the owner's feedback
   on `8cf2e4c`, in Next: it dropped the process share and first edit.)
2. `grove attempts ID` and the board show the total cost and count over the
   work's attempts.
3. A test on a recorded `events.jsonl` fixture asserts each derived fact and
   the classification's documented limits; bounded reads of a running
   attempt are preserved.
4. `docs/commands.md` "Attempts" states the new facts and how they are
   derived.
5. Owner judgment in the terminal on one real attempt.

## Evidence

The first headless session (2026-09-28) branched `worktree-G-260927-dx0yn`
from `main` at `38511c1`. It started from this record at `sha256:ac01fd5a…`
and wrote plan [G-260928-j8sc0](G-260928-j8sc0-attempt-shape-facts-and.md).
Its commits were `6959b93` (plan), `133afcb` (active), `6fbc8ff`
(implementation) and `e096c33` (review fixes). The candidate was `8cf2e4c`,
reviewed in [G-260928-4ae1r](G-260928-4ae1r-review-of-g-260927-dx0yn.md).

The owner's feedback on `8cf2e4c` (in Next) set the second attempt, a
headless session of 2026-09-28 on the same branch. It started from this
record at `sha256:3afc3841…` and plan `sha256:e68bf1d1…`. Its commits:

- `1779bdc`: drop the process share and the first edit; add the first
  commit outside the record root.
- `655fc05`: review fixes.
- `992017f`: `G-260927-cg6rt`'s Next.
- `01c856b`: the plan revised.

Acceptance 1 and the proposed design now say what the feedback kept.

What changed:

- **`Shape:`** (`internal/attempt/shape.go`, `ReadShape(events)`). It is
  derived when `grove attempt` reads the attempt, from all of
  `events.jsonl`, bounded per line by the loop `ReadEvents` shares
  (`eachLine`). It counts the tool calls, the Bash calls with a segment
  running `grove`, and guide prints by name. A running attempt's facts are
  "so far", and a skipped line makes the counts `≥`. The board does not
  read it and keeps its bounded window. The process share, the first edit,
  and the record-root classification of reads and edits are gone.
- **`Changed:`** is recorded in `result.json` at finish. It is `git diff
  --name-only BASE HEAD`, split under and outside the record root that
  `grove.yaml` at the base names, with up to ten names each. It also holds
  the committer time of the first commit that touched a file outside the
  root, following first parents, so a resolution's merge counts and the
  target's older commits do not. That time prints as "first commit outside
  it D after the start", or "before the start" when the clock says so. An
  older attempt reads `Changed: unknown: not recorded when it finished`,
  and a Git failure says why.
- **Totals.** `attempt.Sum`. `grove attempts ID` ends with the total line,
  and the board's attempts list for one work opens with it. `events.turns`
  sums `num_turns` over a run's result events. `≥` marks attempts that hold
  only their last result's turns. An unfinished attempt, or one without a
  result event, is said beside the total, not summed.
- `attempt --json` carries `shape` (`tool_calls`, `grove_calls`,
  `guide_prints`, `skipped_lines`), `shape_error`, `result.changed`
  (`first_other_commit` included) and `result.events.turns`.

Against the acceptance:

1. Real output of `grove attempt G-260927-n4wvk.20260928T004357Z`:
   - `Shape: 36 tool calls, 8 running grove; guides printed: work 1`
   - `Changed: unknown: not recorded when it finished`. This is right: the
     attempt predates the fact.

   No existing attempt has `Changed` recorded, so only tests show it filled
   in. `TestChangedAndShapeFacts` covers a commit only under the root and
   then two outside it, a resolution's merge of older target commits (it
   fails without `--first-parent`), a time before the start, and a base Git
   cannot read. `TestRunToResult` covers it through a fake provider.
2. `grove attempts G-260923-fwakw` ends `Total: 2 attempts, $11.79, ≥22
   turns, 26m`. `TestAttemptScreensReconnectAndStop` asserts the board's
   total for one work, and that the list of all attempts has none.
   `TestAttemptCommandsUsage` asserts the CLI line.
3. `TestReadShape` runs on the recorded fixture
   `internal/attempt/testdata/shape-events.jsonl`: 19 calls, 5 of them
   running `grove`, and guide prints `work 1, model 1`. It checks each
   documented limit: `env grove`, `$(grove …)` and `FOO=1 grove` are
   missed, and the heredoc line is a false positive. It also checks an
   oversized line, a running attempt and an empty shape. `TestSum` covers
   the totals. `TestReadEventsBounded` and the board's
   `TestAttemptActivityIsBounded` pass unchanged.
4. `docs/commands.md` "Attempts" states the facts, how they are derived,
   and the known misses. `docs/board.md` names the board's total.
5. This is for the owner, in a terminal.

Verification at `01c856b`, all passing:

- `go vet ./...`
- `gofmt -l .` (empty)
- `go run ./cmd/grove check` (OK: 218 records)
- `go test -count=1 -timeout 120s ./...`
- `go test -short ./internal/attempt` (2.8 to 4.2 s)

At `1779bdc`, `TestOwnerLost` failed once. It found no init event in an
attempt whose provider the test kills right after it starts. This work does
not touch that path. It passed three isolated reruns and the package
rerun. `terminal.py` was not rerun: nothing since `6fbc8ff`, where it
passed, changed TUI code.

Review: [G-260928-av630](G-260928-av630-review-of-g-260927-dx0yn.md),
independent `grove-reviewer` subagents in three rounds. Round 1 found one
medium and two low findings, round 2 one low finding. All were fixed.
Round 3 closed with `Open findings: none`.

Limits:

- `Shape` misses what `docs/commands.md` names: `grove` through a wrapper
  or `$(…)`, and a guide read as a file.
- `Changed`'s file lists come from the whole diff, so a resolution attempt
  counts the merged target's files as outside the record root. The first
  commit follows first parents; the lists do not. This was already so at
  `8cf2e4c`.
- `Result event:` still prints the last result event's turns.
- `attempts` without an ID prints no total, and `attempts` has no `--json`.

**Resolution of the conflict with `main` at `3ec0f68`.** The owner approved
candidate `24fc77a`, then fed back that it conflicts with `main` at
`3ec0f68` (Next). A third headless session of 2026-09-28, on the same
branch from record `sha256:0450c168…`, merged that commit, not a later
`main` and not a rebase, as `7ac3118`. So `24fc77a` and every review's
`examined` stay ancestors.

- One file conflicted: `grove/G-260927-cg6rt-print-the-work-guide-in.md`,
  in its Next. `main` had replaced it with the review handoff, three steps
  and its verdict line; this branch (`992017f`) had reworded the paragraph
  naming the facts this work builds. Both are kept: `main`'s Next whole, the
  verdict line still last, and this branch's paragraph after the steps. Only
  its `Assign: /grove-work G-260927-cg6rt.` lead is dropped, since that
  record is `done` on `main`. No side was taken whole.
- `docs/commands.md` and `internal/cli/cli.go` merged cleanly in disjoint
  regions (`main`'s `guide --part`, this branch's `attempts` total). Nothing
  else changed.
- Verification at `7ac3118`, all passing: `go vet ./...`; `gofmt -l .`
  (empty); `go run ./cmd/grove check` (OK: 220 records); `go test -count=1
  -timeout 120s ./...`.
- Review: [G-260928-a8tfe](G-260928-a8tfe-review-of-g-260927-dx0yn.md), one
  independent `grove-reviewer` round on `7ac3118`, `Open findings: none`.
  Its observation: `Shape` counts `grove guide work --part NAME` as a `work`
  print, the approved behaviour, which `docs/commands.md` states; a count per
  part would be its own work.

## Next

In review: the `candidate` field names the candidate, merge `7ac3118` and
the record commit after it, on branch `worktree-G-260927-dx0yn` from `main`
at `38511c1`, with `main` at `3ec0f68` merged. The previous candidate was
`24fc77a`. The owner judges the resolution (`git show --remerge-diff
7ac3118`) and one real attempt in a terminal (acceptance 5):

- Run `go run ./cmd/grove attempt ATTEMPT` on any attempt that `go run
  ./cmd/grove attempts` lists. Attempts launched after the merge also show
  the filled `Changed:` line.
- Run `go run ./cmd/grove attempts G-260925-pbx81`.
- On the board, press `A` on a work's detail.

Then, in this worktree:

```sh
go run ./cmd/grove approve G-260927-dx0yn "VERDICT"   # or: go run ./cmd/grove feedback G-260927-dx0yn "TEXT"
```

and in the `main` checkout:

```sh
go run ./cmd/grove integrate G-260927-dx0yn
```

Feedback on candidate 8cf2e4c, 2026-09-28: Process share and first edit are noise: across 46 attempts 91% of tool calls are Bash, ~800 file writes go through Bash vs 110 Edit/Write, record reads 337 via cat vs 31 via Read; first edit reads 'none' in 17/46 and is late or missing in 34/46. Drop both; keep Changed, totals, grove/guide counts. Optionally add time to first commit touching files outside the record root, from git at finish. Addressed in `1779bdc`, `655fc05` and `01c856b`, including the
optional first commit (see Evidence).

Verdict on candidate 24fc77a, 2026-09-28: approved

Feedback on candidate 24fc77a, 2026-09-28: conflicts with main at 3ec0f68 in grove/G-260927-cg6rt-print-the-work-guide-in.md. Resolve only that (grove resolve): in this branch, git merge 3ec0f688c3cd3e6acf3fa42183dbe1d10e43b6e1, that commit of main even if main has moved since, never a rebase; resolve those files keeping both sides' intent; rerun the repository's verification; and hand off the merge as the new candidate, with the previous candidate 24fc77a, the merged commit and the resolved files in Evidence. Change nothing else. If a resolution needs a choice this record does not settle, stop with a checkpoint naming it.

Verdict on candidate f1f75b2, 2026-09-28: approved
