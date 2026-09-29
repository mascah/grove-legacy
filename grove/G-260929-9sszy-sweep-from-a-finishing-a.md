---
id: "G-260929-9sszy"
type: plan
title: "Sweep from a finishing attempt's owner and from the board"
status: current
created: "2026-09-29T00:24:52Z"
updated: "2026-09-29T00:25:12Z"
work: ["G-260928-dtrnw"]
---

# Plan for G-260928-dtrnw

Base: `main` at `9a18f57`, branch `worktree-G-260928-dtrnw`, one
implementer. No prerequisite, no blocking question.

## Design

**The owner's sweep.** `sweep.Own(dir)` runs `attempt.Own(dir)` and then
`sweep.After(dir)`; `cmd/grove` and every test binary that is its own owner
call it instead of `attempt.Own` (`attempt` cannot import `sweep`). `Own`
closes the attempt's lock when it returns, so the attempt reads as finished
and the sweep's own "attempt is running" wait does not hold it back. After:

- reads `attempt.json` and `result.json`; does nothing unless a member's
  result state is `review` with a candidate, and the launch names a target;
- unsets `GROVE_ATTEMPT_OWNER` first, so a verify command's `go test` or
  `grove` is never taken for the attempt's owner;
- finds the target's checkout among the worktrees (the project inside it by
  the launch's prefix); none, or tracked changes in it, is reported and
  nothing runs; no policy there means nothing is written at all;
- plans for the handed-off IDs alone (`Plan(root, ids...)`), runs, and
  appends each fact to the attempt's `sweep.log`.

**One sweep at a time.** `Run` takes a non-blocking lock beside the
attempts directory (`<common>/grove/sweep.lock`) and refuses while another
sweep holds it: CLI, board and owner alike.

**Reading it.** `attempt.View.Sweep` holds `sweep.log`'s lines; `Facts`
(`grove attempt`) and the attempt screen print them.

**The board.** `Backend.SweepPlan` (dry run) and `Backend.Sweep` (run), both
in the target's checkout. After a re-read, once no other read is pending and
the Review cards' predictions are held, the board plans once; each Review
card is tagged `sweep: ACT` and the Review block names act and reason. `S`
on the board asks y/n, then runs the sweep as an action whose facts show on
the result screen. The board never sweeps unprompted (the record's default).

## Steps

1. `internal/sweep`: `Plan` ID filter, the lock in `Run`, `Own`/`After`,
   `sweep.log`; tests with the fake provider for acceptance 1–3.
2. `internal/attempt`: `View.Sweep`, facts line; `cmd/grove` and the
   TestMains call `sweep.Own`.
3. `internal/tui` and `internal/cli/tui.go`: backend, planning read, card
   tag, Review row, `S`; tests.
4. `docs/commands.md` (Attempts, Sweep), `docs/board.md`, help text.
5. The real-provider trial needs a budget the owner names at assignment;
   none was named at first, so it waited on question G-260929-s0f25,
   whose answer (30 USD, sonnet 5.5, effort high, run by the resumed
   session) the resumed session followed.
