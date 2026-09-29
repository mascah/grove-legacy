---
id: "G-260929-pjqxp"
type: plan
title: "Plan for G-260928-kehya: resume an attempt's session with run --resume"
status: current
created: "2026-09-29T00:40:25Z"
updated: "2026-09-29T00:40:30Z"
work: ["G-260928-kehya"]
---

## Design

Plan for [G-260928-kehya](G-260928-kehya-resume-an-attempt-s-sess.md),
written headless from main `9a18f57` on `worktree-G-260928-kehya`. It
implements the record's proposed design and settles the choices that design
leaves open. Each choice is routine.

**Flag.** `--resume` is a boolean on `run`, parsed by `attempt.Flag`, so
`grove run` and the board's launch line take it alike. `Request.Resume bool`
gets a case at the top of `Flag`, since the table is string-valued. A second
`--resume` is refused with "may only be supplied once", and `--resume=…`
stays an unknown option. In the CLI it applies only to `run`, through the
check that already covers `--until`, `--branch` and `--worktree`.
`resolve` and the board's conflict line refuse it as they refuse `--until`:
a resolution carries a new mandate and does not continue an earlier
assignment's session.

**What is resumed.** When `req.Resume` is set, `prepare` makes these checks
after `locate` and before anything is written:

- The branch's registered worktree must be the one reused (`reused`).
  Otherwise the launch is refused: "--resume continues the session in the
  worktree it ran in; BRANCH has no worktree at WORKTREE; launch without
  --resume". Claude Code keeps a session under its working directory, so a
  recreated checkout would not find it anyway.
- The source is the newest attempt from `List(root, IDs[0])` that meets four
  conditions. It is `Finished`. Its `Branch`, `Worktree` (`samePath`) and
  `Prefix` equal this launch's. Its selected IDs equal `req.IDs` in order,
  or, for an attempt from before selections, its `Work` alone does. The
  bound may differ, so a plan continuation resumes the `--until plan`
  attempt.
  Interrupted attempts are skipped, because the record says "latest
  finished". `Start` already refuses a running or orphaned one. When no
  attempt qualifies, the launch is refused: "no finished attempt of IDS on
  BRANCH to resume".
- The session counts as known when the source recorded a `session_id` and
  its result's events hold the provider's init event
  (`Result.Events.Init != nil`). A provider that never started never
  created the session. Otherwise the launch is refused: "attempt X never
  started its provider's session; launch without --resume". Grove does not
  read the provider's private session store. A session the provider no
  longer holds, for example after its cleanup period, ends the attempt at
  once with the provider's own "No conversation found" result at $0, and
  `grove attempt` shows it.
- The provider cannot resume: at this base Claude Code is the only provider,
  and it can resume, so this refusal has no code path yet.
  [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) reshapes the
  command composition. Whichever of the two lands second adds the refusal
  for a provider without resume and for a source attempt on another
  provider, and may compose Codex's `exec resume`. If y2p5h is on the base
  when this is implemented, do it here and say so in Evidence.

**Composition.** `Start` appends `--resume SESSION --fork-session` after
`--session-id NEW`. Forking keeps one session per attempt. `session_id`
stays this attempt's own, and the source's transcript stays as it ended, so
a second resume or a fresh comparison starts from the same place. Each
attempt's events and cost are then its own process's alone. On 2026-09-29,
against an unknown session and with no API call, Claude Code 2.1.284
accepted `--resume S --fork-session --session-id N`. It refused
`--session-id` with `--resume` but without `--fork-session`. The unknown
session ended with an `error_during_execution` result costing $0. The
prompt is unchanged: `/grove-work IDS [--until plan] --interaction
headless`. The resumed session runs the guide again from the top, and step
7 of that guide rereads whatever changed. A comparison with a fresh launch
then varies only the session.

**Recorded and shown.** `Launch.ResumedFrom string
json:"resumed_from,omitempty"` holds the source attempt's ID, and the
session appears in `command`. `prepared` carries the source session to
`Start`. `Requested` appends `, resuming ATTEMPT`. That one change shows the
resume in the CLI's launch lines, in `grove attempt`'s `Requested:` and in
the board's `Asked`. `Explain` (`--dry-run`) adds `Resume: ATTEMPT, session
S`. `digest` gains a `resume ATTEMPT` line only when resuming. A newer
finished attempt between the preview and the launch then changes the
digest, and existing digests stay the same. The board's `launchText` names
the resume.

**Docs.**

- `docs/commands.md`, Attempts: `[--resume]` in the synopsis, and one
  paragraph on what it resumes, the fork, the refusals, and that it is off
  by default.
- `internal/cli/cli.go`: the same in the usage synopsis and the `run`
  paragraph.
- `docs/board.md`, Attempts: the launch line takes `--resume`, which is
  meant for `question answered: R again` and `plan ready`. The conflict line
  refuses it.

Out of scope: a default, which is the owner's later choice; a board key or
prompt for resuming; resume for `resolve`; and any change to the guide.

## Steps

1. `internal/attempt/attempt.go`: `Request.Resume`, the `Flag` case,
   `Launch.ResumedFrom`, the source lookup and refusals in `prepare` (one
   helper beside `locate`), and the composition in `Start`.
   `selection.go`: `digest` and `Explain`. `facts.go`: `Requested`.
2. Tests in `internal/attempt/attempt_test.go` with the fake provider:
   - After a finished attempt, `Start` with `Resume` composes `--resume <its
     session> --fork-session` after a new `--session-id` and records
     `resumed_from`. `Requested` and `Facts` show it.
   - A `--until plan` attempt followed by a resume without the bound
     resumes across the change of bound.
   - Refusals, each before anything is written: no finished attempt, a
     source without an init event, a different selection, and a worktree
     removed after the source (`git worktree remove`).
   - `Flag` parses `--resume` once and refuses it twice.
   Use `-short` skips as the neighbouring tests do.
3. `internal/cli`: `--resume` applies only to `run`, the usage text is
   updated, and a parse test goes beside `TestAttemptCommandsUsage`.
   `internal/tui/attempts.go`: `resolved` refuses `--resume` on a conflict
   line, `launchText` names it, and the launch-line test gets a case.
4. The docs as above, with a link and consistency check.
5. Real provider, capped at $3: in a disposable clone of the candidate,
   reached by an absolute `--project` path (CLAUDE.md), on a fixture work
   record there, run a bounded `--until plan` attempt, then `run ID
   --resume` without the bound. Verify:
   - the resumed attempt's init and result;
   - that its events hold only the new process's output;
   - its cost and turns;
   - that `grove attempt` shows `resumed_from`.
   Then remove the clone.
6. Acceptance 2 is not in this candidate. By the record's own terms it is
   measured on the next three question answers or plan continuations here
   after this lands, as acceptance 5 of
   [G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) was.
   Evidence records the baseline those three will be read against: the
   fresh relaunches after a question or a plan-ready in `grove attempts`
   here, with cost, duration and outcome. The handoff says that acceptance
   2 is pending and why.
7. Verification per CLAUDE.md. The launch line changes, so also run
   `python3 internal/tui/testdata/terminal.py BINARY`.
