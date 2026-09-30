---
id: "G-260930-62nmj"
type: work
title: "Design the portable workflow journey and contracts"
status: proposed
created: "2026-09-30T01:05:45Z"
updated: "2026-09-30T01:16:05Z"
kind: investigation
size: large
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260928-vdhf0", "G-260923-tnn5e", "G-260921-3qgsf", "G-260921-jatts", "G-260921-btyck", "G-260921-sth8q"]
---

## Outcome

The portable milestone has a reviewable interaction design and architecture
contract that developers can implement without inferring the product from
the existing commands. The owner can follow one complete journey, including
recovery and waiting, and judge whether it matches the intended workflow.

Owner-selected direction: [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md). This work elaborates the
acceptance in [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md); it does not reopen that direction by
default or create another editable product brief.

## Scope and constraints

Use realistic examples and terminal sketches for adoption, shaping,
delegation, returning, answering, continuation, review and delivery.
Expose outcome, evidence and available actions before branch/session detail.
Make configuration and authority understandable without teaching the record
schema. Both interactive and unattended use must fit.

Proposed responsibilities to test in the design:

- Grove owns work identity, authority, checkpoints, evidence and delivery
  standing, with deterministic transitions where software can check them.
- Harness integrations own launch/stop, capabilities, provider events,
  permissions/limits and native-session handling.
- Delivery integrations perform or observe integration and establish
  correspondence with the approved candidate.
- CLI and TUI consume the same facts and operations.

These are hypotheses for boundaries, not mandated Go interfaces. Inspect
which workflow guarantees are currently instructions and which are enforced
by code; identify which must become mechanical for reliable continuation.
Avoid both a forced fresh process per activity and an assumption that every
activity shares one private session.

Observed at main c8fa070ef9ff: internal/attempt mixes Claude command/session
handling with ownership and worktrees; internal/tui already has a backend
boundary; internal/update's completion check requires candidate ancestry.
The record/Git/approval packages and tests supply reusable evidence. The
current record model and guides remain implemented contracts during design.

## Acceptance

1. Linked design artifacts, represented through Grove's supported records,
   show one coherent journey and the failure/wait states from the milestone,
   with the primary screen/action at each moment and advanced details
   reachable. The owner reviews the actual sketches.
2. Explicit examples define assignment input, checkpoint/continuation,
   candidate/review/approval evidence, and local/hosted delivery
   correspondence. Each fact has one authoritative owner; process events
   and UI labels cannot silently create competing lifecycle state.
3. Provider capabilities distinguish supported, unsupported and unknown,
   including cost reporting, enforceable limits, permissions, cancellation,
   interactive launch and native resume. Missing capabilities fail clearly
   or follow a documented supported path.
4. The design explains retention and migration of existing work, reviews,
   candidate objects, configuration and native sessions, including what a
   fresh clone can continue without local logs.
5. Compare adaptation, targeted replacement and wider replacement using
   the complete local trial's requirements. Name code to retain, refactor,
   replace or investigate, and behavioral tests to retain or reconsider.
   No blanket rewrite or blanket preservation follows from this record.
6. Reconcile dependent proposals with the chosen contracts, recording any
   owner decisions and new dependency edges. The owner accepts the
   interaction and architectural design before dependent implementation.
   This design does not authorize paid trials or changes in adopter repos.

## Next

Ready for a design assignment: $grove-work G-260930-62nmj.
Inspect the milestone and its member records, then produce concrete journey
sketches and contract examples for the owner's review. Preserve unresolved
consequential choices as questions; investigate routine technical uncertainty.
