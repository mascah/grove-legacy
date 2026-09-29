---
id: "G-260929-zznyj"
type: review
title: "Review of G-260928-dtrnw: sweep from a finishing attempt and the board"
status: current
created: "2026-09-29T00:57:22Z"
updated: "2026-09-29T02:55:38Z"
work: ["G-260928-dtrnw"]
examined: "6c0a1d2d076bad90f28fa889f765e3f0b501e175"
---

## Examined

Branch `worktree-G-260928-dtrnw`, base `main` at `9a18f57`, against
[G-260928-dtrnw](G-260928-dtrnw-run-sweep-from-a-finishi.md) and plan
[G-260929-9sszy](G-260929-9sszy-sweep-from-a-finishing-a.md). Three rounds
through the `grove-reviewer` agent, a fresh one each, read-only:

- Round 1 at `f211c97` (the diff `9a18f57..f211c97`).
- Round 2 at `c5bbec1`.
- Round 3 at `0006993`, the last round the work guide allows.

Each ran `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove check` and
`go test -count=1` on `internal/sweep` and `internal/tui`, plus `-short` on
the rest; acceptance 1 to 4 met by the tests they ran each round.

## Findings

Round 1, five:

1. Consequential: an owner sweep that met a running sweep was refused, and
   no later sweep would pick its candidate up. Fixed in `c5bbec1`: the
   owner waits on the sweep lock, saying so in `sweep.log`, and plans and
   acts under it; `grove sweep` and the board's `S` still refuse. Test:
   "while another sweep runs".
2. Minor: `After` was silent on a `grove.yaml` that does not parse, and
   called a failed `git status` uncommitted changes. Fixed in `c5bbec1`.
3. Minor: the board's plan read could not be cancelled. Fixed in `c5bbec1`
   (`sweep.PlanContext`) and `0006993` (the remaining `repo.Git` reads).
4. Minor: the board does not follow an attempt's sweep, which runs after
   the attempt reads as ended. Documented in `docs/board.md`; no code
   change.
5. Minor: doc wording (lock, README, record model, work guide). Fixed in
   `c5bbec1`.

Round 2, two, both minor: `PlanContext`'s comment overclaimed (three Git
reads without ctx), and a stale test comment. Fixed in `0006993`.

Round 3, one, minor, open: `internal/sweep/sweep.go:78` says "PlanContext
is Plan, whose Git reads end when ctx does", but the branch check through
`update.Branch` (`symbolic-ref`) still ignores ctx. Fix: check the branch
with `repo.GitContext` in `PlanContext`, or narrow the comment.

Knowledge, every round: no new concept needs a term; nothing contradicts
G-260925-wh9ax, G-260923-tnn5e or the Policy term G-260926-a8vyj. The limit
G-260925-5wrn8 listed, that a finishing attempt's owner starts no sweep, is
what this record changes.

Gate 2, after the resumed session, two rounds through a fresh
`grove-reviewer` each, read-only:

- Round 1 at `51823e2` (the fix, `0006993..51823e2`; the combined
  `9a18f57..51823e2`; and the real-provider trial's evidence).
- Round 2 at `6c0a1d2` (`51823e2..6c0a1d2`, the combined tip briefly).

Each ran `go vet ./...`, `gofmt -l .`, `go run ./cmd/grove check`,
`go test -count=1` on `internal/sweep` and `internal/update` (round 2 also
`internal/attempt`) and `go test -short ./...`: all clean.

Gate 2, round 1, one finding, minor: the round-3 finding was fixed in `51823e2`
(`update.BranchContext` in `sweep.PlanContext`), but `(*Sweep).plan` still
listed attempts through `attempt.List`, whose `repo.Locate` ignores ctx.
Fixed in `6c0a1d2`: `attempt.ListDir` on the common directory the plan's
`versions.InspectContext` already read, as the board does.

Gate 2, round 2: none. Every Git read reachable from `PlanContext` now ends with
ctx; the attempts directory is the same one by construction.

Noted, not counted as findings: no sweep test makes the plan wait on a
running or orphaned attempt (untested at the base too); the trial binary,
built with `go build` in this linked worktree, is stamped with `main`'s
`9a18f57` rather than `51823e2`, so the attempt's `Started:` line names
that commit, while its symbols (`update.BranchContext`) are `51823e2`'s.

Trial (acceptance 5), judged in round 1: met on the evidence.

## Disposition

Gate 1 stopped at the review cap with its round-3 finding open; gate 2
fixed it and the one it found after; nothing is open.

Open findings: none
