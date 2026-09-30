---
id: "G-260928-dtrnw"
type: work
title: "Run sweep from a finishing attempt and from the board"
status: accepted
created: "2026-09-28T19:28:59Z"
updated: "2026-09-30T13:18:01Z"
kind: feature
size: medium
relates_to: ["G-260925-5wrn8", "G-260925-wh9ax", "G-260926-a8vyj", "G-260923-tnn5e"]
candidate: "6b110517b0deb262f4ec42b023ac2ee4783db18d"
approved: "6b110517b0deb262f4ec42b023ac2ee4783db18d"
approved_by: owner
approved_context: "sha256:dd9330a84ff7bc6f98933ff5eac4efb0f13610b934d095c3e6956e34c7d412cf"
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
[G-260929-9sszy](G-260929-9sszy-sweep-from-a-finishing-a.md) (step 5
reconciled with the answer); code at `6c0a1d2`, the candidate the commit
after it that holds this evidence. Started from record `sha256:c9a00063…`
and plan `sha256:425158d0…`; resumed from record `sha256:57d0c624…`.

- The owner process runs `sweep.Own`: `attempt.Own`, then `sweep.After`,
  which sweeps only the members the result holds in review with a
  candidate, in the target's checkout, under its committed policy, and
  appends each fact to the attempt's `sweep.log` (`View.Sweep`; `grove
  attempt` prints `Sweep:` lines; the attempt screen shows them). It unsets
  `GROVE_ATTEMPT_OWNER` before any `verify` command.
- One sweep per repository at a time (`.git/grove/sweep.lock`): `grove
  sweep` and `S` refuse while it is held, an owner waits for it.
- The board plans once per re-read, after the predictions and cancellable
  (`sweep.PlanContext`, every Git read of which now takes ctx), tags each
  Review card `sweep: ACT`, and adds a `Sweep:` row to the Review block;
  `S` asks y/n, then sweeps. It never sweeps unprompted (the record's
  default).

Acceptance, 1 to 4 in `internal/sweep/sweep_test.go` (TestMain's fake
provider through the real owner) and `internal/tui`:

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
5. Met. Docs: `docs/commands.md` (Attempts; Sweep, "After an attempt"),
   `docs/board.md`, and README, `docs/record-model.md` and
   `docs/work-execution.md` where they said who runs a sweep. Trial, under
   [G-260929-s0f25](G-260929-s0f25-what-budget-may-the-real.md)'s answer:
   a disposable project (`/tmp/dtrnw-trial.IXSY/proj`, `grove init`, target
   `main`, policy verify `true` and `test -f hello.txt`, `max_lines: 200`,
   `integrate: true`) with one work record, "Add a hello file"; a binary
   built from `51823e2` launched `grove run G-260929-fedy7` with budget 30
   USD, model `claude-sonnet-5-5`, effort high. Attempt
   `G-260929-fedy7.20260929T024646Z` handed off in 86 s for 0.35 USD, and
   its owner's `Sweep:` lines say: verified the merge with main at
   `e1b0fc2`, approved under policy (the verdict names review `G-260929-d1kw2`
   and the attempt), fast-forwarded main to `09649d9`, done at `ceaecb6`.
   `git log main` there ends with those two commits; `hello.txt` says
   `hello`. `6c0a1d2` changed only how the plan lists attempts after it.

Decisions: the resumed session fixed round 3's finding in code rather than
narrowing the comment (`update.BranchContext`), and gate 2's by reading
attempts from the common directory the plan had already read.

Verification at `6c0a1d2`: `go vet ./...` and `gofmt -l .` clean; `go run
./cmd/grove check` OK, 252 records; `go test -count=1 -timeout 120s ./...`
all ok. `terminal.py` last ran at `0006993` (six ok); nothing after it
touches `internal/tui`.

Review: [G-260929-zznyj](G-260929-zznyj-review-of-g-260928-dtrnw.md),
examined `6c0a1d2`. Gate 1, three rounds to its cap, one minor finding
open; gate 2, two rounds: that finding and one more minor one fixed, the
last round "Open findings: none".

Limits: no sweep test makes the plan wait on a running attempt (untested at
the base too). The trial binary, built with `go build` in this linked
worktree, is stamped `9a18f57`, so the attempt's `Started:` line names
`main`'s commit; its symbols are `51823e2`'s. The diff touches no path the
policy's `never` lists but exceeds `max_lines`, so it waits for the owner.

## Next

Owner: judge the candidate in this checkout, then integrate in the target's:

```sh
grove approve G-260928-dtrnw "VERDICT"
grove --project /Users/mascah/GitHub/mascah/grove integrate G-260928-dtrnw
```

G-260928-c5j9d names this record in its `depends_on`.

Verdict on candidate 6b11051, 2026-09-29: approvd

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
