---
id: "G-260930-62nmj"
type: work
title: "Design the portable workflow journey and contracts"
status: accepted
created: "2026-09-30T01:05:45Z"
updated: "2026-09-30T20:22:13Z"
kind: investigation
size: large
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260928-vdhf0", "G-260923-tnn5e", "G-260921-3qgsf", "G-260921-jatts", "G-260921-btyck", "G-260921-sth8q", "G-260930-npw49", "G-260930-84fnb", "G-260930-y6fyy", "G-260930-2qa4a", "G-260930-tcc9w"]
candidate: "f376a6dab7a5999560b25adaacb93d919c36e7a8"
approved: "f376a6dab7a5999560b25adaacb93d919c36e7a8"
approved_by: owner
approved_context: "sha256:75ee6b958c76e5a6240aac5cb3eaed17f2d2a75764ef7b03b3413650f6bc3262"
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

## Evidence and owner acceptance

On 2026-09-29 (owner's local date), the owner reviewed the written designs
at commit f376a6dab7a5999560b25adaacb93d919c36e7a8 and explicitly accepted them:
"I've reviewed the designs and I accept them."

The accepted artifact commit is already on the configured target main.
This closes the design investigation; it does not claim any proposed
runtime behavior has been implemented or any agent trial has passed.
Candidate and approval name that exact reviewed commit. The original
completion-model answer remains attributed in accepted decision
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md), and question
[G-260930-y6fyy](G-260930-y6fyy-should-completion-be-der.md) remains resolved.

| Acceptance | Delivered evidence |
| --- | --- |
| 1: journey and owner judgment | [Experience sketches](G-260930-npw49-portable-workflow-experi.md), including waits/recovery; owner's acceptance above |
| 2: shared contracts | [Contracts design](G-260930-84fnb-portable-workflow-contra.md), assignment through recorded acceptance and verified delivery |
| 3: provider capabilities | Contracts design's capability matrix and failure semantics |
| 4: retention and migration | Contracts design's fresh-clone, retained-ref, legacy completion and schema migration requirements |
| 5: implementation strategy | Contracts design's alternatives and refactoring ownership map |
| 6: reconciliation and approval | Brief and dependent proposals reconciled at f376a6d; explicit local-delivery prerequisite for portable execution; owner's acceptance above |

Mechanical evidence at f376a6d: grove check passed for 269 records, 114 local
links resolved, dependency edges matched the stated order, and diff checks
passed. These checks establish document consistency, not product usability.
No independent agent review or provider trial is claimed for this shaping
work. Closure validation also passed: grove check (269 records), 74 local
links across the nine acceptance updates, balanced code fences and diff
checks. The final dependency preview verifies the design prerequisite is
delivered before reporting local delivery ready.

## Next

The design at f376a6d remains an accepted historical artifact. On 2026-09-30
the owner reconsidered its delivery and execution boundaries in
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md). Review the
reconciled [contracts](G-260930-84fnb-portable-workflow-contra.md) and
[experience](G-260930-npw49-portable-workflow-experi.md) before assigning
new implementation; this accepted record does not pretend the old candidate
approved later writing. Its original evidence is preserved above.

The milestone owns the changed membership, dependency graph and adopter hold.
Local delivery G-260929-gm3m4 is already delivered; new bounded workspace and
sequence work address the selected changes without reopening that migration.
No implementation, paid trial or adopter migration is assigned here.

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
