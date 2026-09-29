---
id: "G-260928-kehya"
type: work
title: "Resume an attempt's session after its question is answered"
status: active
created: "2026-09-28T19:29:00Z"
updated: "2026-09-29T03:05:53Z"
kind: feature
size: small
relates_to: ["G-260924-wp2pe", "G-260924-5b6pz", "G-260923-tnn5e", "G-260928-y2p5h"]
---

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

## Evidence

Tested at commit `318fe8e` (`cad3572` code and tests, `318fe8e` docs), then
re-run after review fixes (usage rewrap, escaped `Resume:` line, baseline row) at
`0d114ab`: `go vet`, `gofmt -l .`, `grove check` (250 records), `go test -short` of
attempt, cli and tui, and `go test -count=1 -timeout 120s ./internal/attempt
./internal/cli` all pass. Branch
`worktree-G-260928-kehya`, Claude Code 2.1.284.

Acceptance 1 (fake provider) and 3:

- `go vet ./...` clean; `gofmt -l .` empty; `go run ./cmd/grove check` OK, 250
  records.
- `go test -count=1 -timeout 120s ./...` passes in every package (attempt
  20.7s, tui 18.8s, cli 9.7s; the attempt and tui totals include process-spawning
  tests that skip under `-short`).
- `python3 internal/tui/testdata/terminal.py BINARY` on a build of that
  commit: every case ok, including `attempt_lifecycle`.
- New tests: `TestResume` (fork composition after `--until plan`, `resumed_from`
  in `Requested` and `Facts`; refusals with nothing written for no worktree, a
  different selection, a source without an init event, a removed worktree),
  `TestResumeFlag`, and cases in the CLI usage table and the board's launch and
  resolve lines.
- `docs/commands.md` and `docs/board.md` (Attempts) and `grove --help` describe
  `--resume`; the record links resolve (`grove check`).

Step 5, real provider, disposable fixture (a fresh `grove init` repository with
one work record, reached by an absolute `--project` path, built from
`318fe8e`, removed afterward; not a clone of this repository, to keep the
plan run small):

- `--until plan` attempt `G-260929-2jvac.20260929T030906Z`: success, 4 turns,
  $0.1245, 1 permission denial, init present, session `d5d6a1c0-...`.
- `run --resume` without the bound: attempt `G-260929-2jvac.20260929T030925Z`,
  command carried `--session-id d94ccd8d-... --resume d5d6a1c0-... --fork-session`,
  `grove attempt` showed `resuming G-260929-2jvac.20260929T030906Z` in
  `Requested:` and `resumed_from` is in `--json`. Success, 1 turn, $0.1472,
  0 denials, init present. Its events (9 lines) carry only the new session id
  `d94ccd8d-...`, so cost and turns are that process's alone. The
  `--dry-run` preview showed `Resume: ...` and a digest.
- Total spend about $0.27 against the $3 cap.

Acceptance 2 is pending. Baseline the three resumes will be read against: fresh
relaunches after a question or a bounded plan in `grove attempts` here (cost,
turns, duration, record status after):

| Work | Follows | Fresh attempt | Cost | Turns | Duration | Outcome |
| --- | --- | --- | --- | --- | --- | --- |
| G-260924-59f5k | plan | 20260924T232247Z | $6.01 | 53 | 2644s | stopped (exit 0), error_during_execution, active; not a natural end |
| G-260925-7c8g9 | plan | 20260925T232126Z | $16.15 | 139 | 2519s | success, review |
| G-260925-ced1h | plan | 20260925T042810Z | $3.18 | 56 | 572s | success, review |
| G-260927-n4wvk | question | 20260928T011916Z | $4.00 | 51 | 1345s | success, review |
| G-260928-dtrnw | question | 20260929T024521Z | $2.77 | 38 | 661s | success, review |
| G-260928-pqhyg | question | 20260929T003437Z | $1.49 | 31 | 782s | success, active |
| G-260928-pqhyg | question | 20260929T005419Z | $0.77 | 17 | 191s | success, active |

## Next

The candidate awaits review. Acceptance 2 is not measured here: it is measured
on the next three question answers or plan continuations after this lands
(plan step 6), read against the baseline in Evidence. Nothing else waits.

This touches the command composition in `internal/attempt/attempt.go`, which
[G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) reshapes; whichever
lands second adds the refusal for a provider without resume.
