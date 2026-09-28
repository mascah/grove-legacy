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

- **Record root.** Both new facts need it, and `attempt.json` does not hold
  it. It is read from `grove.yaml` at the attempt's base commit
  (`git show BASE:PREFIX/grove.yaml`, `project.RecordRoot`), which outlives
  the worktree. Unreadable, the facts that need it are unknown.
- **Shape, at read time.** `ReadShape(events, root)` scans all of
  `events.jsonl` line by line, bounded per line as `ReadEvents` is (the
  same loop, factored out), and classifies each `tool_use`, subagents'
  included, so its index matches `Metrics`' tool count:
  - process: a Bash command one of whose `;`, `&&`, `||`, `|` or newline
    segments runs `grove` (a program named `grove`, or `go run …/cmd/grove`)
    or `git worktree`, or makes a `worktree-` branch with `git branch`,
    `checkout -b` or `switch -c`; an Edit, Write, MultiEdit or NotebookEdit
    under the record root; a Read under it. A guide is read by `grove guide
    NAME`, so it is a grove command; guide prints are also counted by NAME.
  - the first edit outside the record root: the first Edit, Write, MultiEdit
    or NotebookEdit whose path is not under it, with its tool index and its
    event's timestamp; elapsed is from the launch's `started`.
  - limits, documented: `grove` through a wrapper, `$(…)`, `env` or a
    variable assignment is missed; an edit made through Bash (`sed -i`, a
    script) is not an edit; a path spelled through a symlink of the root
    is outside it.
  Only `grove attempt` reads it (`ShowContext` with events). The board
  keeps its bounded window and does not show it. For a running attempt it
  is labelled "so far", and `≥` where an oversized line was skipped.
- **Files changed, at finish.** `reconcile` adds `changed` to
  `result.json`: `git diff --name-only BASE HEAD` in the worktree, split
  into inside and outside the record root. An attempt finished before this
  prints it as not recorded.
- **Totals.** `attempt.Total(views)` sums attempts, the result events' USD
  and turns, and minutes from start to finish, and counts what it could not
  sum (no result event, still running). `grove attempts ID` prints it after
  the table; the board's attempts screen for one work prints it under the
  header. An attempt of a selection counts in full for each member.
- `attempt --json` carries `shape`, `shape_error` and `result.changed`.

## Steps

1. `internal/attempt`: line loop factored out of `ReadEvents`; `Shape`,
   `ReadShape`, record root, `Changed` at reconcile, `Total`; the `Shape:`
   and `Changed:` facts. Fixture test of each fact and each documented limit.
2. `internal/cli`: total line under `attempts ID`; test.
3. `internal/tui`: total under the attempts screen header for one work; test.
4. `docs/commands.md` "Attempts".
5. Verification as `CLAUDE.md` says, an independent review, handoff.
