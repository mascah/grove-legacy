---
id: "G-260921-btyck"
type: term
title: "Approval"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:33Z"
relates_to: ["G-260921-jatts", "G-260921-rz7bn", "G-260921-3qgsf"]
formerly: "T-006"
---

## Meaning

A person with authority, or someone they delegated to, accepting one specific
[candidate](G-260921-jatts-candidate.md) as meeting the work's outcome. It is bound to
that candidate: if the candidate changes, the approval does not carry over and
must be reconsidered.

Not a [review](G-260921-rz7bn-review.md), which informs it; not passing checks, which
cannot supply judgment; and not [integration](G-260921-3qgsf-integration.md), which
follows it. Since G-260921-jwk4e, `grove approve` records it: the work record's
`approved` field names the candidate, the verdict is appended to the body,
and `check` rejects an approval of any other commit, so a changed candidate
needs its own. `grove feedback` withdraws it and returns the work to active.

## Relationships

Follows review, precedes integration.
