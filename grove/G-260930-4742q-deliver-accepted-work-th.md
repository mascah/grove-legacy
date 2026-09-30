---
id: "G-260930-4742q"
type: work
title: "Deliver accepted work through a protected hosted pull request"
status: proposed
created: "2026-09-30T01:05:45Z"
updated: "2026-09-30T01:16:06Z"
kind: feature
size: large
depends_on: ["G-260930-gwnb1"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260929-gm3m4", "G-260921-3qgsf"]
---

## Outcome

A developer using a protected hosted target can deliver Grove work through
the repository's PR workflow and see accurate completion and waiting reasons
in Grove. Required host checks and approvals retain their authority, and an
external merge can be observed without requiring Grove to perform it.

Owner direction: [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md); milestone: [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).
GitHub is the proposed first implementation because this project already
uses it; the design must name the supported host and exercised PR policy.
This does not promise every host or acquire git-host responsibilities.

## Scope and constraints

Use the delivery/candidate contracts established by the local proof.
Bind the submitted head, Grove candidate, reviewed evidence, approval and
delivered result. A host PR status or a trailer alone cannot establish that
the approved code was delivered. Host-native requirements can impose
additional waits; Grove approval does not bypass them.

Proposed operations are publishing/linking the selected candidate, presenting
checks and approval waits, requesting integration when authorized, and
observing the final destination. Posting, pushing and merging require an
explicit assignment or applicable standing policy. Inspection has no such
side effect. Reconcile work completion without an extra unapproved direct
write to a protected target.

Local-only repositories remain supported without credentials, a remote,
a network service or hosted CI. A host outage leaves delivery unknown or
waiting with a reason rather than creating a local success claim.

## Dependencies

Depends on [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md): it consumes the exercised candidate, checkpoint
and transformed-delivery contracts. Designing against the earlier
ancestry-only completion mechanism would require redoing the adapter.

## Acceptance

1. In an authorized hosted test repository, a bounded candidate follows a
   protected PR path through checks, approval and the supported merge method.
   Grove reports the same standing through CLI and TUI.
2. A failing check, missing host approval, moved target, changed PR head,
   closed-unmerged PR and unavailable host each produce an accurate wait or
   invalidation. Old Grove approval cannot authorize a different candidate.
3. A merge performed outside Grove is reconciled idempotently with evidence
   connecting the accepted candidate to the resulting destination. A squash
   is covered, including retained evidence and cleanup behavior.
4. Completion tracking requires no bypass of branch protection and no
   mandatory metadata commit after the protected merge.
5. Existing local-only operation works without host setup or network access.
   Documentation names the exact host capabilities and limitations tested.

## Next

Needs [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md). At assignment, confirm the proposed host against the
accepted design and obtain a specific repository/action mandate for the real
trial; fixtures cover mechanics without external publication.
