---
id: "G-260930-gj9d7"
type: decision
title: "Prove delivery once, at integrate, and read Done from the target's record"
status: accepted
created: "2026-09-30T16:20:33Z"
updated: "2026-09-30T16:27:32Z"
relates_to: ["G-260930-2qa4a", "G-260929-gm3m4", "G-260930-44q35", "G-260921-3qgsf", "G-260930-84fnb"]
---

## Decision

On 2026-09-30 the owner reviewed
[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) at its review
cap and rejected one part of it: a verifier that re-proved every past
delivery from Git on every read (list, show, deps, context, run selection,
the board). In their words:

- "I need grove to NOT be a complete performance hog."
- "We have to be mindful that other people may not have good git branch
  hygiene. They could have tons of local branches for repos that are many
  years old with thousands of commits. We can't assume grove will only be
  used in projects with tiny git history or footprint."
- "I didn't authorize adding more git work per load explicitly."
- "I also want to make sure that were not junking up other users git
  histories."
- On re-delivering from a branch kept after a squash: "dropping the
  continuating case is fine."

They said yes to this direction: keep the clean squash commits and stop
re-proving on every read.

- **Reading.** A work record whose acceptance applies is done when the
  target branch tip's own copy of that record, at the same path, is accepted
  for the same candidate; otherwise it is accepted, awaiting delivery. An
  acceptance reaches the target only through a delivery, so nothing else is
  consulted: no history, trailers or `refs/grove`. A read starts at most one
  Git process beyond what it did before, and none on the board, whatever the
  history, branches or deliveries.
- **Delivery.** `integrate` proves the one squash commit it makes before it
  reports done, scoped to that commit.
- **Audit.** Re-proving past deliveries is `grove check --deliveries`, run
  on purpose.
- **Continuation.** A branch kept after its delivery merges the target
  before its next one, as any branch does; there is no special merge base.
- **What Grove leaves in a repository.** One squash commit per delivery,
  with its trailers, and one local ref per delivery under
  `refs/grove/submitted/`, which only the audit reads. Nothing a teammate
  must fetch for Grove to read correctly.

## Alternatives

Re-proving on every read, bounded to a fixed number of Git processes, was
built and measured: 11 to 14 processes per board load against 6 on main,
and about 61 ms of CPU for `grove list` on this repository against 25 ms,
with no measurement of a large repository and a cost that grew with the
target's history. A cache was not considered: it adds state to keep right.
The owner accepted the trade the cheap reading makes: an acceptance written
on the target by hand reads as done, and only the audit notices.

## Reconsideration

Amends [G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md), whose
choice stands: Done is derived from a recorded acceptance, with no done
commit. Reopen if hand-written acceptances on the target, or deliveries the
audit cannot prove, turn out common enough that reading must notice them;
bring the evidence and the cost to the owner.
