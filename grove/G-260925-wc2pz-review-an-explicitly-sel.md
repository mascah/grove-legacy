---
id: "G-260925-wc2pz"
type: decision
title: "Review an explicitly selected set of work together, on one shared candidate"
status: superseded
created: "2026-09-25T23:22:15Z"
updated: "2026-09-30T20:22:03Z"
relates_to: ["G-260925-7c8g9", "G-260925-80w3a", "G-260925-t70h8", "G-260921-jatts", "G-260921-btyck", "G-260921-3qgsf", "G-260923-tnn5e", "G-260930-tcc9w"]
---

## Superseded boundary, 2026-09-30

The owner selected both modes, with separate delivery as default and together
explicit, in [G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md).
The original decision below explains current assignments and the combined
mode's provenance; it no longer requires every future selection to share a
candidate. Existing assignments keep their recorded boundary until completed.
New behavior is proposed until its implementing work lands.

## Decision

On 2026-09-25 the owner answered
[G-260925-80w3a](G-260925-80w3a-where-should-review-and.md) with its option 1: when
the owner explicitly selects several work items, they are implemented in
order on one branch, and a later item may build on an earlier item's
unmerged changes. Human review and integration stop the selection once,
at its end, not between items. The owner approved the consequences in plan
[G-260925-t70h8](G-260925-t70h8-selected-work-one-attemp.md) by launching its continuation:

- One Grove-owned attempt runs the whole selection sequentially under one
  aggregate budget; Grove never retries it.
- The complete members enter review together with one shared
  [candidate](G-260921-jatts-candidate.md): the records on a branch whose
  `candidate` is the same commit are one group, and nothing else defines
  it.
- [Approval](G-260921-btyck-approval.md) stays per member, so each keeps its own
  acceptance and verdict. [Integration](G-260921-3qgsf-integration.md) is per group:
  merging the commit merges all of it, so it waits until every member is
  approved.
- Feedback on any member reopens the group, dropping every member's
  approval, since the next candidate replaces the shared one.
- A started member that is incomplete holds the whole branch out of
  review; members that never started keep their wait and do not.

A dependency edge still means only that the prerequisite is needed: it
does not itself authorize shared execution, which the owner's explicit
selection does. The meanings of work, attempt, candidate, approval and
integration do not change; a group of one behaves exactly as before.

## Alternatives

- **Review and merge between dependent items** (G-260925-80w3a option 2): smaller
  reviews, but an unattended chain stops at its first boundary.
- **A boundary chosen per assignment** (option 3): flexible, but doubles
  the execution, recovery and verification contracts before either is
  proven.

## Reconsideration

Reconsider when combined reviews prove too large to judge, when owners
need intermediate judgment inside a selection, or when parallel
implementation is selected.
