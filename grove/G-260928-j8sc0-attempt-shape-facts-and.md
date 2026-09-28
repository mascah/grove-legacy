---
id: "G-260928-j8sc0"
type: plan
title: "Attempt shape facts and per-work totals"
status: current
created: "2026-09-28T00:46:00Z"
updated: "2026-09-28T00:46:21Z"
work: ["G-260927-dx0yn"]
---

## Design

Implements [G-260927-dx0yn](G-260927-dx0yn-retain-per-attempt-proce.md)'s proposed
design, from `main` at `38511c1`.

- **Record root.** Files changed needs it, and `attempt.json` does not hold
  it. It is read from `grove.yaml` at the attempt's base commit
  (`git show BASE:PREFIX/grove.yaml`, `project.RecordRoot`). Unreadable,
  files changed is unknown.
- **Shape, at read time.** `ReadShape(events)` scans all of
  `events.jsonl` line by line, bounded per line as `ReadEvents` is (the
  same loop, factored out), and counts each `tool_use`, subagents'
  included, as `Metrics` does; the Bash calls one of whose `;`, `&&`, `||`,
  `|` or newline segments runs `grove` (a program named `grove`, or `go run
  …/cmd/grove`); and `grove guide NAME` runs by NAME. Limits, documented:
  `grove` through a wrapper, `$(…)`, `env` or a variable assignment is
  missed; a heredoc line starting with `grove` counts. Only `grove attempt`
  reads it (`ShowContext` with events). The board keeps its bounded window
  and does not show it. For a running attempt it is labelled "so far", and
  `≥` where an oversized line was skipped.
- **Files changed, at finish.** `reconcile` adds `changed` to
  `result.json`: `git diff --name-only BASE HEAD` in the worktree, split
  into inside and outside the record root, which `grove.yaml` at the base
  names, and the committer time of the first commit, along first parents,
  that touched a file outside it. An attempt finished before this prints it
  as not recorded.
- **Totals.** `attempt.Total(views)` sums attempts, the result events' USD
  and turns, and minutes from start to finish, and counts what it could not
  sum (no result event, still running). `grove attempts ID` prints it after
  the table; the board's attempts screen for one work prints it under the
  header. An attempt of a selection counts in full for each member.
- `attempt --json` carries `shape`, `shape_error` and `result.changed`.

Revised after the owner's feedback on `8cf2e4c` (in the work record's
Next): the process share and the first edit outside the record root were
dropped, and the time to the first commit outside it added.

## Steps

1. `internal/attempt`: line loop factored out of `ReadEvents`; `Shape`,
   `ReadShape`, record root, `Changed` at reconcile, `Total`; the `Shape:`
   and `Changed:` facts. Fixture test of each fact and each documented limit.
2. `internal/cli`: total line under `attempts ID`; test.
3. `internal/tui`: total under the attempts screen header for one work; test.
4. `docs/commands.md` "Attempts".
5. Verification as `CLAUDE.md` says, an independent review, handoff.

All five steps done at `e096c33`, and redone for the feedback in `1779bdc`
and `655fc05`; the evidence is in
[G-260927-dx0yn](G-260927-dx0yn-retain-per-attempt-proce.md). Built as designed,
with two additions: `result.json`'s events gain `turns`, `num_turns` summed
over a run's result events, since each counts one query of a resumed session
while the cost is cumulative, and the total's turns show `≥` where an older
attempt holds only its last result's; and an unreadable HEAD at finish is
recorded as `Changed`'s error.
