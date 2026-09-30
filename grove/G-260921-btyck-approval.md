---
id: "G-260921-btyck"
type: term
title: "Approval"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-30T03:18:06Z"
relates_to: ["G-260921-jatts", "G-260921-rz7bn", "G-260921-3qgsf", "G-260930-2qa4a"]
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

## Selected completion contract, 2026-09-29

The owner selected [G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md):
record candidate acceptance and derive Done from its target delivery,
with a redesigned work record that does not remain misleadingly in review.
This changes representation, not the distinction between approval and
integration.

## Schema 4, 2026-09-30

[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) implements it.
`grove approve` sets the work `accepted` and records, beside `approved`,
`approved_by` (`owner`, or `policy sha256:…` for a standing policy's
sweep), the authority of whoever ran it, never read from the verdict's text,
and `approved_context`, the digest of the record's title and body outside
`## Next` and Grove's own appended paragraphs. A changed candidate or a
changed context makes the acceptance no longer apply, and the work reads as
in review again. `grove feedback` withdraws all three. Approval is still
not delivery: accepted work is done only once the target's own copy of the
record holds that acceptance, which only a delivery brings there
([G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md)).
