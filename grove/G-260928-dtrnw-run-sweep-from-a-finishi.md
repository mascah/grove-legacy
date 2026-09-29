---
id: "G-260928-dtrnw"
type: work
title: "Run sweep from a finishing attempt and from the board"
status: active
created: "2026-09-28T19:28:59Z"
updated: "2026-09-29T00:25:25Z"
kind: feature
size: medium
relates_to: ["G-260925-5wrn8", "G-260925-wh9ax", "G-260926-a8vyj", "G-260923-tnn5e"]
---

## Outcome

A candidate that meets the owner's standing policy is integrated without
anyone running `grove sweep`: the attempt that produced it triggers the
sweep when it ends in review, and the board shows and runs one, so Review
holds only what the policy leaves to the owner.

Owner intent, shaping conversation 2026-09-28: "The workflow should
automatically try and integrate it. I'm expecting a continuous stream of
work to be shaped and launched over time, and the review phase is one that
causes a lot of wasted time sitting there for a human in some cases."

## Constraints

Observed 2026-09-28 at main `6fbb888`:

- `grove sweep` (`internal/sweep`) runs in the target's checkout under
  `policy:` and refuses without one.
  [G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md) named the
  finishing attempt as a possible trigger in its design and then listed as a
  limit: "Trigger: `grove sweep` only, run by the owner or a scheduler; the
  owner process of a finishing attempt does not start one." Grove starts no
  resident service ([G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)).
- The owner process outlives the terminal, knows the launching project and
  target, and writes `result.json` at the end (`docs/commands.md`, Attempts).
  The board re-reads on focus, on `r`, and when an attempt ends or a branch
  tip moves (`docs/board.md`); nothing in `internal/tui` references sweep.
- A sweep merges in the target's checkout and runs `verify` in a temporary
  worktree; an unchanged wait does not retry; a resolution it starts is an
  attempt whose end is itself a trigger, bounded by once per target commit
  and the policy's budget (`docs/commands.md`, Sweep).
- keyborg's `grove.yaml` has no `policy:`, so nothing there is automatic
  until its owner writes one. This repository's policy allows 300 changed
  lines and never `.claude/**`.

Proposed design, labelled proposed:

- The owner process, after writing `result.json` for an attempt whose record
  ended in review with a candidate, runs the equivalent of `grove sweep`
  for that candidate alone, in the target's checkout, under the target's
  committed policy, and writes what it did or why it waited into the
  attempt's files. Refused, and recorded as such, when the target's checkout
  is missing, dirty or holds a running sweep; never for a project without a
  policy. A resolution attempt's end triggers the same.
- The board runs a dry run on each re-read and shows the pending act on
  each Review card; `S` runs the sweep. Whether the board also runs it
  unprompted is a routine choice for preparation, defaulting to no.
- The Review block and the attempt screen name the sweep's act and reason.

Out of scope: what the policy checks
([G-260928-c5j9d](G-260928-c5j9d-judge-a-candidate-agains.md) extends it), a
resident process, a scheduler, a policy for keyborg (its owner writes it
there).

## Acceptance

1. With a policy and the fake provider, an attempt that hands off a clean,
   in-policy candidate is integrated by its own owner process; the record
   shows the delegated verdict and the merge; the attempt's files say the
   sweep ran and what it did.
2. An out-of-policy or conflicting candidate waits or resolves as `sweep`
   would, and the attempt's files say why; a later re-read does not retry
   an unchanged wait.
3. Without a policy nothing changes; a dirty or absent target checkout is
   reported and nothing is written.
4. The board shows the pending act per Review card and runs a sweep on `S`.
5. `docs/commands.md` (Attempts, Sweep) and `docs/board.md` say so; one
   real-provider trial in a disposable project under a budget the owner
   names at assignment.

## Evidence

Branch `worktree-G-260928-dtrnw`, base `main` at `9a18f57`, plan
[G-260929-9sszy](G-260929-9sszy-sweep-from-a-finishing-a.md); code at
`0006993`. Not a candidate yet: see Next.

- The owner process runs `sweep.Own`: `attempt.Own`, then `sweep.After`,
  which sweeps only the members the result holds in review with a
  candidate, in the target's checkout, under its committed policy, and
  appends each fact to the attempt's `sweep.log` (`View.Sweep`; `grove
  attempt` prints `Sweep:` lines; the attempt screen shows them). It unsets
  `GROVE_ATTEMPT_OWNER` before any `verify` command.
- One sweep per repository at a time (`.git/grove/sweep.lock`): `grove
  sweep` and `S` refuse while it is held, an owner waits for it.
- The board plans once per re-read, after the predictions and cancellable
  (`sweep.PlanContext`), tags each Review card `sweep: ACT`, and adds a
  `Sweep:` row to the Review block; `S` asks y/n, then sweeps.

Acceptance, all in `internal/sweep/sweep_test.go` (TestMain's fake provider
through the real owner) and `internal/tui`:

1. Met: `TestAnAttemptsOwnerIntegratesItsCandidateInsideThePolicy`: done on
   main, delegated verdict naming the attempt, `sweep.log` says verified,
   approved and done.
2. Met: subtests "out of policy" (waits, never path named) and "a
   conflict" (a resolution attempt starts; its own end sweeps and waits:
   "again after a resolution of main at …").
3. Met: "no policy" (no `sweep.log`, nothing written), "a dirty target",
   "no target checkout"; also "while another sweep runs".
4. Met: `TestReviewCardsShowTheSweepAndSRunsIt`, and the attempt screen's
   `Sweep` rows in `TestAttemptScreenHonesty`.
5. Docs met: `docs/commands.md` (Attempts; Sweep, "After an attempt"),
   `docs/board.md`, and README, `docs/record-model.md` and
   `docs/work-execution.md` where they said who runs a sweep. The
   real-provider trial is not run: no budget was named
   ([G-260929-s0f25](G-260929-s0f25-what-budget-may-the-real.md)).

Verification at `0006993`: `go vet ./...` and `gofmt -l .` clean; `go run
./cmd/grove check` OK, 250 records; `go test -count=1 -timeout 120s ./...`
all ok; `python3 internal/tui/testdata/terminal.py` all six ok. `go test
-count=1 ./internal/sweep` alone takes about 3 s.

Review: [G-260929-zznyj](G-260929-zznyj-review-of-g-260928-dtrnw.md), three
rounds, the last at `0006993`: one minor finding open (a comment in
`internal/sweep/sweep.go` overclaims ctx), the rest fixed.

Also removed a duplicate set of empty headings at the top of this body.

## Next

Checkpoint 2026-09-29, headless attempt: active, stopped at the review
cap with one open finding, and waiting on question
[G-260929-s0f25](G-260929-s0f25-what-budget-may-the-real.md) (the trial's
budget, model and effort).

1. Owner: answer G-260929-s0f25.
2. Then resume on this branch: `/grove-work G-260928-dtrnw`. Fix the open
   finding (use `repo.GitContext` for the branch check in
   `sweep.PlanContext`, or narrow its comment), review it at a new gate,
   run the trial under the answered budget, and hand off with the candidate.

G-260928-c5j9d names this record in its `depends_on`.
