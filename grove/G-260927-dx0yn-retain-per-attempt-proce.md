---
id: "G-260927-dx0yn"
type: work
title: "Retain per-attempt process facts and show cost per work record"
status: active
created: "2026-09-27T22:13:36Z"
updated: "2026-09-28T00:46:25Z"
size: medium
relates_to: ["G-260921-sth8q", "G-260923-tnn5e", "G-260923-895zb", "G-260921-h46pb", "G-260923-p5pt6"]
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

## Acceptance

1. `grove attempt` shows for a finished attempt: process tool calls and
   share, guide prints by name, the first edit outside the record root
   (elapsed, tool index, or none), and the files changed inside and outside
   the record root; an unknown fact is labelled unknown, never filled.
2. `grove attempts ID` and the board show the total cost and count over the
   work's attempts.
3. A test on a recorded `events.jsonl` fixture asserts each derived fact and
   the classification's documented limits; bounded reads of a running
   attempt are preserved.
4. `docs/commands.md` "Attempts" states the new facts and how they are
   derived.
5. Owner judgment in the terminal on one real attempt.

## Next

Assign: `/grove-work G-260927-dx0yn`. The evals' work row, listed in
G-260923-p5pt6's Next, is not this record.
