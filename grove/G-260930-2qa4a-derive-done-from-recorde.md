---
id: "G-260930-2qa4a"
type: decision
title: "Derive Done from recorded acceptance and verified delivery"
status: accepted
created: "2026-09-30T03:12:14Z"
updated: "2026-09-30T03:17:58Z"
relates_to: ["G-260930-y6fyy", "G-260930-62nmj", "G-260930-84fnb", "G-260929-gm3m4", "G-260930-4742q", "G-260921-vr8a8", "G-260921-btyck", "G-260921-3qgsf"]
---

## Decision

On 2026-09-29 (owner's local date), the owner answered
[G-260930-y6fyy](G-260930-y6fyy-should-completion-be-der.md):
"Derive Done with a redesigned record contract."
They asked to record the choice and finish preparation up to implementation.

For new work under the redesigned contract, persist the candidate and its
attributable acceptance. Derive Done from applicable acceptance plus verified
delivery to the configured target. Redesign the recorded lifecycle so a
completed item does not permanently declare status=review. A mandatory
post-delivery done commit or follow-up completion PR is not the selected path.

## Reasons and alternatives

A board-only Done label would leave agents reading misleading raw records.
Keeping explicit done writes is viable, including follow-up PRs for protected
targets, but adds another delivery cycle, failure recovery and dependency
wait solely for bookkeeping. The owner found that machinery clunky and
selected a redesigned record contract instead.

The tradeoff is explicit: the record tells a reader what was accepted;
Git/Grove establishes whether it reached the target. An unchanged standalone
file cannot certify a later external event. Agent-readable context and
mechanical action checks must carry that observation with its evidence and
freshness, and missing evidence cannot be represented as success.

## Consequences and boundaries

- The brief owns the selected direction; the
  [contracts design](G-260930-84fnb-portable-workflow-contra.md) develops it.
  CLI, board, context, dependency readiness, policy and cleanup must consume
  the same verified completion result, with exact source and candidate binding.
- Local delivery [G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md)
  owns the record redesign, migration and common completion/delivery checks.
  The hosted path consumes them without a second completion workflow.
- Candidate identity, approval authority and evidence remain distinct from
  provider success, PR labels and merge messages. Changed work cannot inherit
  completion from an earlier acceptance. Evidence must survive cleanup and
  be transportable to a fresh clone without a required service.
- This selects the completion model, not every field spelling, screen or
  interface in the design. No product implementation, paid trial, merge or
  push is assigned. Current schema 3 commands and lifecycle remain in force
  until a deliberate migration delivers the new contract.
- Preserve historical completion claims without inventing missing candidate,
  review or approval evidence. A raw-file-only agent behavior evaluation is
  required; a prompt to ignore contradictory status is not a solution.

## Reconsideration

Reopen if a bounded implementation cannot reconstruct trustworthy delivery
from retained project evidence, or if direct-file and tool-mediated agent
trials show unacceptable confusion despite the redesigned representation.
Bring that evidence to the owner; do not silently restore completion PRs or
weaken what Done promises.
