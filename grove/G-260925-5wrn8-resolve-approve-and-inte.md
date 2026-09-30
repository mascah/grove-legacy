---
id: "G-260925-5wrn8"
type: work
title: "Resolve, approve and integrate candidates under an explicit owner policy"
status: accepted
created: "2026-09-25T21:39:28Z"
updated: "2026-09-30T13:18:01Z"
kind: feature
depends_on: ["G-260925-h8rj5", "G-260925-dz10z"]
relates_to: ["G-260921-jwk4e", "G-260921-rz7bn", "G-260921-btyck", "G-260921-3qgsf", "G-260923-tnn5e", "G-260924-5b6pz", "G-260924-ecs9m", "G-260925-beby3", "G-260925-7c8g9", "G-260925-80w3a", "G-260925-w33j7", "G-260925-wh9ax", "G-260926-a8vyj"]
candidate: "2251569c061a951ad7ba5d09da8c5d242a678e8d"
approved: "2251569c061a951ad7ba5d09da8c5d242a678e8d"
approved_by: owner
approved_context: "sha256:96002868ffaf331dc3294d68aeafc1e48caf19aa8a14599906a0735e6ef0f255"
---

## Outcome

Where the owner has written a standing policy, Grove starts a resolution
attempt for a candidate that conflicts with the target, and approves and
integrates a candidate that meets the policy's conditions after its
independent review, with each act attributed to the policy and its
evidence, so a set of independently implemented items can go from review to
the target while the owner is away, and the owner judges only what the
policy leaves to them.

Decision [G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md), the owner on
2026-09-25 answering [G-260925-w33j7](G-260925-w33j7-what-may-a-standing-owne.md): option 4,
automatic resolution plus delegated approval and integration under a
narrow written policy the owner extends. The initial policy below is
proposed at the owner's request and binds nobody until the owner writes it
into `grove.yaml`.

## Constraints

### Observed evidence

At main `47852e3`, `approved` is a commit that must equal `candidate`, and
the verdict is a body paragraph (record model, work lifecycle); nothing in
the schema distinguishes the owner's verdict from a delegated one, and
`check` accepts any equal value, which [G-260921-jwk4e](G-260921-jwk4e-review-candidates-and-in.md)
notes as the hand-set escape that skips the verdict and tip checks. The
independent reviewer (`.claude/agents/grove-reviewer.md`,
[G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md)) returns findings and
limits and is "never approval"; the [work guide](../docs/work-execution.md)
dispatches it per gate and caps fix rounds at three. `integrate` merges,
then writes done, and runs no verification; approval needs the branch's
clean checkout and integration the target's. Each attempt records its
model, effort, reviewer definition hash and cost
([G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md),
[G-260924-ecs9m](G-260924-ecs9m-default-an-attempt-s-bud.md)), and nothing resident runs between
attempts ([G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)). G-260925-w33j7's evidence lists the
authority and spend constraints.

### Proposed design and scope

- A `policy:` mapping in `grove.yaml` beside `run:`, absent by default.
  Proposed initial policy for this repository, narrow on purpose:

  ```yaml
  policy:                 # standing delegation (G-260925-wh9ax); absent means nothing automatic
    budget: 30            # dollars, aggregate for every automatic act in one sweep
    resolve:
      budget: 10          # one attempt per candidate per target movement; run: defaults otherwise
    approve:
      verify:             # must pass on the merged result, in a temporary worktree
        - go test -count=1 -timeout 120s ./...
        - go vet ./...
        - go run ./cmd/grove check
      max_lines: 300      # added plus removed, the record's own file excluded
      never:              # a change to any of these always waits for the owner
        - grove.yaml
        - .github/**
        - .claude/**
        - .agents/**
        - go.mod
        - go.sum
        - lefthook.yml
    integrate: true       # merge and write done after a delegated approval
  ```

  `never` is the trust boundary: files that change what the automation
  itself does or what CI runs. `max_lines` bounds the blast radius. `verify`
  is this repository's final evidence from its agent instructions. The
  budgets are guesses to adjust. Likely extensions, deliberately absent
  from the initial policy: a `paths:` allowlist, a `kinds:` list (for
  example only `fix` and `tooling`), a `sizes:` list, and a higher
  `max_lines`.
- Conditions the policy cannot switch off, so that its narrowness holds by
  construction: the record is in review with its candidate unchanged since
  the review; a review record examined that candidate and reports no open
  finding and no knowledge finding; no open question blocks the work; the
  merge is clean (a conflict goes to resolution, never to approval); the
  merged result passed `verify`. How "no open finding" is read from a
  review record is chosen in preparation, since the schema does not
  structure findings: a conventional closing line the reviewer definition
  writes, or a review field.
- Attribution: a delegated approval writes `approved` and a verdict
  paragraph that keeps the `Verdict on candidate X, DATE:` prefix
  `integrate` already quotes, with text beginning `delegated under policy`
  and naming the `grove.yaml` revision, the review record, the attempt and
  the verification result, so it is distinguishable from the owner's own
  verdict wherever the record is shown; the merge and `done` are the
  ordinary `integrate` acts. Introduce a field only if prose attribution
  proves insufficient.
- Verify before the target moves: make the merge in a temporary worktree,
  run the policy's commands there, and only then merge on the target. A
  failure leaves the target unchanged and the candidate in review with the
  reason, exactly as a refused `integrate` does.
- Trigger: the owner process of a finishing attempt, or an explicit sweep
  command the owner or an external scheduler runs. Grove starts no resident
  service (G-260923-tnn5e); an unchanged wait does not retry.
- Everything the policy does not name waits for the owner as today.
  Reconcile the brief's lifecycle paragraph and the [approval](G-260921-btyck-approval.md)
  term if the answer changes their meaning; the existing reviewer
  definition suffices initially, with its findings consumed by the policy.
- Out of scope: choosing the policy's conditions for the owner; a second
  reviewer definition; executing several selected items
  ([G-260925-7c8g9](G-260925-7c8g9-execute-an-explicitly-se.md)); a shared candidate for a chain
  ([G-260925-80w3a](G-260925-80w3a-where-should-review-and.md)).

## Acceptance

1. With no policy, nothing changes. With one, every automatic act is visible
   in advance, as what would happen to each candidate and why, and
   attributable afterwards, with the policy, evidence, attempt and cost in
   the record and on the board.
2. A candidate outside the policy waits for the owner. One inside it is
   integrated only after the merged result passed the policy's verification;
   a failure leaves the target and the record unchanged with the reason.
3. Automatic resolution runs at most once per target movement within the
   aggregate budget; a repeated conflict waits with the reason.
4. A delegated verdict is distinguishable from the owner's in the record and
   on the board, and reversible: the record names the merge to revert.
5. Exercised with fake providers on clean, conflicting, failing-verification
   and out-of-policy candidates, then a separately bounded real-provider
   trial on a disposable project before use in this repository.

## Evidence

Implemented headless on `worktree-G-260925-5wrn8` from main `fd7744e` (G-260925-h8rj5 and
G-260925-dz10z integrated). It started from this record at `sha256:805b71a6…` and
plan [G-260926-vpvhf](G-260926-vpvhf-plan-for-policy-driven-i.md) at
`sha256:ccb8d406…` (commit `42e07ed`). Code runs through `8eb741f`; the
candidate adds only this evidence and the review record.

**What was built.**

- `policy:` in `grove.yaml` (`internal/project/policy.go`), keys as the
  proposed policy above. `resolve` needs `budget`; `approve` needs a
  non-empty `verify`; `integrate: true` needs `approve`; `never` takes
  `path.Match` patterns or `DIR/**`. Every problem is a `check` diagnostic
  named `policy.KEY`. At the owner's request the proposed policy above was
  then written into this repository's `grove.yaml` (`2251569`), before the
  real-provider trial.
- `grove sweep [--dry-run]` (`internal/sweep`), in the target's checkout
  with `grove.yaml` committed, so every act names the revision it ran
  under (`policy grove.yaml sha256:…`). Refused without a policy. For each
  candidate in review it plans skip, wait (with the reason), resolve,
  approve or integrate; `--dry-run` prints that and writes nothing.
- Resolution is G-260925-dz10z's `attempt.Resolve` with a new `Request.Policy`: the
  feedback begins `delegated under policy grove.yaml sha256:…, budget N
  USD:`, and the attempt takes its budget, permission mode, model and
  effort from the target's `run:`, never the candidate branch's.
- Approval first predicts the merge again, then merges the branch's tip
  into that target commit in a temporary worktree outside every checkout
  and runs each `verify` command there (`sh -c`, Git's location variables
  removed). Only then `update.Approve` writes the verdict `delegated under
  policy grove.yaml sha256:…: review ID examined X with no open finding;
  merged with TARGET at T, verification passed (COMMANDS); attempt A
  produced it for N USD` (or that no Grove attempt is recorded as producing
  it), and `integrate.Run` with the new `Expect` (refused unmerged if the
  target moved from T) and `Policy` (the done update appends `Integrated
  under policy … as merge M on TARGET (was B); to reverse it: git revert
  -m 1 M`, or the range of a fast-forward).
- How the policy reads a review, chosen in preparation: the review guide
  now ends every report with `Open findings: none` or `Open findings: N`,
  and the work guide asks the author to copy that line into the review
  record. The policy needs at least one `current` review of the work
  covering the candidate (it examined it, or an earlier commit from which
  only records changed), and every such review must close `Open findings:
  none`. No schema field.
- `update.Delegated` and the board's Review block: `approved under
  policy` for a delegated verdict.
- Knowledge: term [G-260926-a8vyj](G-260926-a8vyj-policy.md) (proposed), "Policy", and a
  delegated act.

**Acceptance 1.** With no policy, `sweep` is refused and nothing else
changed: every other command's behaviour is untouched (full suite).
`sweep --dry-run` shows each candidate's act and why. Afterwards the
record holds the policy revision, the review, the verification and the
attempt with its cost (approval), the policy and budget (resolution; the
attempt and its cost are in `grove attempts`), and the merge and its
revert (integration). The board shows `approved under policy` and the
body's verdict.

**Acceptance 2.** `TestSweepWaitsOutsideThePolicy`: a `never` path,
`grove.yaml` itself, too many lines, a binary file, an open finding, no
closing line, an owner's approval, no `approve` section each wait with the
reason and main unchanged. `TestSweepLeavesAFailedVerificationUnchanged`:
main, the branch and the record unchanged, with the command's last output.
`TestSweepIntegratesACandidateInsideThePolicy` and
`TestSweepIntegratesSeveralCandidatesInOneSweep`: integrated only after
verification, the second against the target the first moved.

**Acceptance 3.** `TestSweepResolvesAConflictOncePerTargetCommit`: one
attempt, attributed, with the target's defaults; a repeated conflict with
the same target commit then waits. `TestSweepKeepsResolutionsInsideTheBudget`:
an attempt the aggregate budget does not cover waits, as does one without a
permission mode.

**Acceptance 4.** The verdict begins `delegated under policy`, which
`update.Delegated`, the board and `integrate`'s approval line show; the done
update names the merge and the revert command
(`TestIntegrateUnderAPolicy`).

**Acceptance 5.** Fake provider and real repositories for clean,
conflicting, failing-verification and out-of-policy candidates and no
policy, above. An end-to-end run of the built binary on a disposable
repository printed the dry run, then verified, approved, fast-forwarded
main and wrote done with the revert line. The real-provider trial on a
disposable project was not run: the record bounds it separately, and no
spend for it was authorized here.

**Docs.** `docs/record-model.md` (`policy:`, delegated approval, Review
status), `docs/work-review.md` (the closing line), `docs/work-execution.md`
(the review record's closing line; the sweep among the dispositions),
`docs/commands.md` (Sweep), `docs/board.md`, `README.md`, `grove --help`.

**Verification at `8eb741f`:** `gofmt -l .` empty; `go vet ./...` clean;
`go run ./cmd/grove check` OK, 191 records; `go test -count=1 -timeout
120s ./...` all ok; `python3 internal/tui/testdata/terminal.py BINARY` all
12 ok. `internal/sweep -short` about 1.1s alone (the resolve test, which
waits on a fake provider, skips under `-short`).

**Review:** [G-260926-mg6g8](G-260926-mg6g8-review-of-policy-driven.md), three
rounds by fresh `grove-reviewer` agents. Four consequential findings (a
stale target for a sweep's later candidates, an outvoted review, and the
candidate's `run:` choosing a delegated attempt's mode, model and effort)
were fixed with regressions; round 3 found nothing consequential.

**Limits for the owner:**

- The real-provider trial of acceptance 5 is open.
- `verify` commands have no timeout; a hanging one hangs the sweep.
- `never` guards approval, not resolution: a resolution attempt runs in
  the candidate's checkout, whatever it changed under `.claude/` or
  `.agents/`.
- `max_lines` counts every file but the record's own, so the plan and
  review records a candidate adds count toward it.
- An integration refused after a delegated approval (the target moved in
  between, or another refusal) leaves the record approved under the
  policy; the next sweep says it waits for the owner's `grove integrate`.
- Trigger: `grove sweep` only, run by the owner or a scheduler; the owner
  process of a finishing attempt does not start one.
- A candidate shared by a group, or in review on several branches, always
  waits.

## Next

In review. The owner judges the candidate on `worktree-G-260925-5wrn8`:

```sh
grove approve G-260925-5wrn8 "VERDICT"   # in /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-260925-5wrn8
grove integrate G-260925-5wrn8           # then in the main checkout
```

Or give feedback with `grove feedback G-260925-5wrn8 "TEXT"` in the worktree. To
demo: add the proposed `policy:` above to a disposable project's
`grove.yaml`, commit it, and run `grove sweep --dry-run`, then `grove
sweep`, in its target checkout. This repository's `grove.yaml` carries the
proposed policy from `2251569`, so the first `grove sweep` on main after
integration acts under it; the real-provider trial is still open.

Verdict on candidate 2251569, 2026-09-26: approved

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
