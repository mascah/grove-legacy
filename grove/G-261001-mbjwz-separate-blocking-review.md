---
id: "G-261001-mbjwz"
type: work
title: "Separate blocking review findings from follow-ups"
status: proposed
created: "2026-10-01T02:17:15Z"
updated: "2026-10-01T02:22:08Z"
kind: feature
size: small
depends_on: ["G-260930-0s29t"]
relates_to: ["G-261001-fg05z", "G-260928-n4f1q", "G-260930-84fnb", "G-260930-60c3d"]
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

Proposed implementation: an explicit, validated review summary consumed by
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

## Next

Authorized for follow-through in the current conversation. Wait for
G-260930-0s29t's walkthrough judgment and delivery, then prepare the exact
summary and delegation-compatibility design in a separate work branch. Do
not expand this into provider execution or a new review runner.
