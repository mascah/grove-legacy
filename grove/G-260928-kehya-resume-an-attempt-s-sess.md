---
id: "G-260928-kehya"
type: work
title: "Resume an attempt's session after its question is answered"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-09-28T19:34:02Z"
kind: feature
size: small
relates_to: ["G-260924-wp2pe", "G-260924-5b6pz", "G-260923-tnn5e", "G-260928-y2p5h"]
---

## Outcome

## Constraints

## Acceptance

## Next

## Outcome

When the owner relaunches work whose question was answered, or whose plan is
ready, they may continue the previous attempt's session instead of starting
fresh, and can compare what each costs and produces.

Owner intent, shaping conversation 2026-09-28: "Resuming sessions after
answering questions rather than starting fresh each time? Or is it better
to re-launch with a fresh context?" Neither is known; this measures it.

## Constraints

Observed 2026-09-28 at main `6fbb888`:

- `attempt.json` records `session_id`. `claude -p --resume SESSION`
  continues a print session and `--fork-session` copies it (Claude Code
  2.1.283 `--help`); Codex has `exec resume`. The events reader already sums
  result events "one per query of a resumed session" (`docs/commands.md`,
  Attempts).
- Every relaunch today is a fresh process on the reused worktree. Re-entry
  cost $1 to $2.50 and 2 to 4 minutes in the three resumed attempts
  [G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) measured.
  A resumed session carries stale reads, which the work guide's step 7
  ("a changed revision means the input changed, so reread it") covers.
- `question answered: R again` and `plan ready: read it, then R` are the
  board's relaunch points
  ([G-260924-wp2pe](G-260924-wp2pe-answer-a-blocking-questi.md)).

Proposed design, labelled proposed: `run ID --resume` resumes the latest
finished attempt's session of that work on its branch with the same
assignment, refused when the worktree is gone, the session is unknown or
the provider cannot resume; the new attempt records `resumed_from`; off by
default; the board's launch line accepts it. Then run it on the next three
question answers or plan continuations here and report cost, duration and
outcome beside fresh relaunches in this record's Evidence. A default is the
owner's later choice.

## Acceptance

1. `--resume` composes the provider's resume of the recorded session, and
   the attempt facts show it; refusals as listed; fake provider.
2. Three real resumes compared with fresh relaunches, reported.
3. `docs/commands.md` (Attempts) and `docs/board.md` (Attempts) say so.

## Next

Assign: `/grove-work G-260928-kehya`. This touches the command composition
in `internal/attempt/attempt.go`, which
[G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) reshapes; no
order is declared between them, and the owner should not launch them
alongside each other without reading both.
