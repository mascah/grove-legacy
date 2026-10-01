---
id: "G-261001-mbjwz"
type: work
title: "Separate blocking review findings from follow-ups"
status: accepted
created: "2026-10-01T02:17:15Z"
updated: "2026-10-01T02:51:12Z"
kind: feature
size: small
depends_on: ["G-260930-0s29t"]
relates_to: ["G-261001-fg05z", "G-260928-n4f1q", "G-260930-84fnb", "G-260930-60c3d"]
candidate: "cd8f821026eba0e2e5e752fb14a597c527300495"
approved: "cd8f821026eba0e2e5e752fb14a597c527300495"
approved_by: owner
approved_context: "sha256:b3570bb87328ec81288fc745f68978a18a6d1b2d970e59d80093bf8a8de6f69d"
---

## Outcome

An independent review distinguishes unresolved acceptance/correctness issues
from improvements that can be tracked without blocking the candidate. A
wording improvement alone does not restart a complete review cycle or imply
that the review found the implementation unacceptable.

Owner authorization: [G-261001-fg05z](G-261001-fg05z-simplify-workflow-bounda.md).

## Constraints

Bounded change to the existing review guides and `internal/sweep` policy
consumer; the separately launched review operation remains
[G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md)'s outcome.

Observed at a448a6c: `docs/work-review.md` requires `Open findings: none`
only when no finding of any weight remains, and `internal/sweep.review`
requires that line in every applicable current review. Changing prose alone
would leave the policy on the old contract.

Define blockers by consequence: unmet acceptance, a correctness/safety or
authority defect, or incomplete required review/verification. A documentation
defect that misstates the supported behavior may be blocking regardless of
its severity label. Follow-ups have evidence and a disposition; an
implementer cannot silently downgrade the reviewer's finding. Preferences
are not defects. Keep the complete historical review, including its original
examined commit and closing line.

Selected implementation: an explicit, validated review summary consumed by
the policy, with complete/incomplete state and separate blockers and
follow-ups. Its exact representation and configuration change belong in the
implementation plan. Legacy zero-finding reviews keep their meaning;
nonzero legacy reviews cannot be reinterpreted as only follow-ups. Existing
delegations must not silently acquire broader approval authority.

For repair rounds, verify the fix, its neighboring states and affected
callers. Repeat a full review when the changed scope warrants it. If a
finding changes the lifecycle rule, reconcile that rule and its operation
table before patching callers. Keep the existing cap and report what is
unresolved; exhausting a cap is not acceptance.

Depends on G-260930-0s29t to settle the current candidate and guide baseline
before changing the review contract used to judge later work.

## Acceptance

1. The report and human handoff visibly distinguish blockers, follow-ups and
   incomplete review, bound to the examined candidate and supporting facts.
2. Under explicitly applicable delegation, a complete review with no blockers
   and tracked follow-ups can proceed; a blocker in any applicable review,
   an incomplete/malformed report, or stale evidence still waits. Other
   policy checks remain effective.
3. Legacy reviews and existing policy configurations have explicit tested
   compatibility behavior. Historical findings are not rewritten or
   implicitly downgraded.
4. Regression cases cover follow-ups only, a blocking documentation defect,
   conflicting reviews, incomplete output, stale candidates, malformed
   summary, and legacy nonzero findings. Show the resulting policy preview.
5. Reconcile the owning model, both guides, help and affected policy output.
   Review this bounded contract change independently; keep the report format
   reusable by the later independent-review operation.

## Evidence

Branch `worktree-G-261001-mbjwz`, base `bf86caa`; implementation `4a2f64d`.
Candidate is the evidence commit named in frontmatter, followed only by this record's handoff.
[Plan](G-261001-hcb0h-plan-review-summaries-an.md) committed at `4c9846a`.

1. The review guide owns a strict versioned summary separating completeness, blockers and follow-ups; the work guide preserves the reviewer's judgment.
2. Every applicable report must be complete and blocker-free; follow-ups need explicit `policy.approve.allow_followups: true`. Preview and verdict disclose all applicable follow-ups. Freshness, scope and verification gates remain.
3. Existing configurations default to false; legacy clean reviews still pass and nonzero legacy reviews still wait. Historical reports are unchanged.
4. New parser/configuration and real sweep fixtures cover the required refusal states, conflicting reviews, stale code, overflow and delivery. Actual preview: `review G-260101-00005 has no blockers; G-260101-00005: 2 follow-ups; ... then approve and integrate`.
5. Model, review/work guides, command documentation and help agree; the shared parser is reusable by later review orchestration.

Verification on the implementation tree committed as `4a2f64d`: the new tests first failed against the old contract, then passed. Focused `go test -short . ./internal/project ./internal/sweep`, `go vet ./...`, `gofmt -l .`, `git diff --check` and `go run ./cmd/grove check` passed (283 records); 80 changed-document local links and fences checked.
One `go test -count=1 -timeout 120s ./...` passed every package except the TUI `TestTerminal/attempt_lifecycle` timeout waiting for its expected drawing. Only `go test -count=1 -timeout 120s ./internal/tui` was rerun and passed (11.507s). The full run was not wholly green; the same timing symptom was recorded in the prerequisite.
Independent [review G-261001-5hkbh](G-261001-5hkbh-review-of-blocking-findi.md) examined `4a2f64d`, reran focused tests and reported `Review v1: complete; blockers=0; follow-ups=0`.
Limits: repository policy was not opted into allowing follow-ups. Provider trials, the future review runner and imported-branch eligibility are outside this change; no claim is made for them.

## Next

Ready for owner judgment of the candidate in frontmatter. After approval, the integrator runs these commands from this checkout; keep the source workspace for inspection:

```sh
go run ./cmd/grove approve G-261001-mbjwz "Owner approved the recorded candidate and review-policy behavior."
go run ./cmd/grove --project /Users/mascah/GitHub/mascah/grove integrate G-261001-mbjwz --keep
```

Then start [G-261001-62n4p](G-261001-62n4p-isolate-workspace-eligib.md) in a fresh worktree from main. That work is already authorized; it waits for this prerequisite's delivery. No merge or push of this candidate is implied by the earlier approval of G-260930-0s29t.

Verdict on candidate cd8f821, 2026-10-01: Owner approved candidate cd8f821026eba0e2e5e752fb14a597c527300495 and authorized local integration and continuation with workspace eligibility in this conversation.
