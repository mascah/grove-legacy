---
id: "G-260925-wh9ax"
type: decision
title: "Delegate conflict resolution, approval and integration to a standing owner policy"
status: accepted
created: "2026-09-25T21:55:20Z"
updated: "2026-09-28T19:34:03Z"
relates_to: ["G-260925-w33j7", "G-260925-5wrn8", "G-260925-dz10z", "G-260925-h8rj5", "G-260921-jatts", "G-260921-rz7bn", "G-260921-btyck", "G-260921-3qgsf", "G-260923-tnn5e", "G-260925-beby3", "G-260925-80w3a", "G-260928-d8py6"]
---

## Decision

On 2026-09-25 the owner answered [G-260925-w33j7](G-260925-w33j7-what-may-a-standing-owne.md)
with its option 4. Grove may act on a standing policy written in
`grove.yaml`, with no per-candidate human act, to start one bounded
resolution attempt when a candidate in review conflicts with the target, at
most once per target movement, and to approve and integrate a candidate
that meets the policy's written conditions after its independent review.
Each act is attributed in the record to the policy and its evidence, is
distinguishable from the owner's own verdict, and is reversible by an
ordinary revert. Anything the policy does not name waits for human judgment
as before. The policy's presence is the owner's standing instruction and
its absence means nothing automatic. The initial policy is narrow, proposed
in [G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md) at the owner's request, and
the owner extends it.

The meanings of [candidate](G-260921-jatts-candidate.md), [review](G-260921-rz7bn-review.md),
[approval](G-260921-btyck-approval.md) and [integration](G-260921-3qgsf-integration.md) do
not change: approval was already the owner's or a delegate's acceptance of
one commit, and a review stays evidence, which the policy consumes.

## Alternatives

- **Nothing automatic** (G-260925-w33j7 option 1): unattended progress for
  implementation only; integration stays a morning chore.
- **Automatic resolution only** (option 2): delegates spend, not judgment.
  Not chosen as the initial state, though a policy without an approval
  section yields exactly it.
- **Delegated approval without automatic resolution** (option 3 alone):
  leaves conflicts to the manual feedback and relaunch loop that started
  this exploration.
- **The review agent decides by itself, with no written policy**: rejected.
  A review is evidence, never a verdict (G-260921-rz7bn); a policy the owner can
  read is what makes each act attributable and extensible.

## Reconsideration

Reopen the conditions, not the delegation, when a delegated integration
lands a defect the written policy should have caught. Reopen the
delegation itself if the owner stops unattended use, if the conditions the
owner actually judges by cannot be expressed as deterministic checks plus
review evidence, or if the harness gains native gated merges that make
Grove's own policy redundant. A later decision supersedes this one if
[G-260925-80w3a](G-260925-80w3a-where-should-review-and.md) selects a shared candidate
for a chain, whose review boundary this decision does not cover.
