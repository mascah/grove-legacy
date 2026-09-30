---
id: "G-260930-y6fyy"
type: question
title: "Should completion be derived from verified delivery?"
status: resolved
created: "2026-09-30T01:31:26Z"
updated: "2026-09-30T03:17:59Z"
relates_to: ["G-260930-62nmj", "G-260930-84fnb", "G-260929-gm3m4", "G-260930-2qa4a"]
blocks: ["G-260930-62nmj"]
---

## Question

Should Grove derive a work item's effective Done state from verified
delivery on the target, without requiring a separate committed done update?

Asked during the portable design continuation on 2026-09-29 (local date).
The owner could decide; it blocked finalizing the lifecycle/delivery contract
in [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md).

## Evidence and alternatives considered

At main 9b5c5a55f355, internal/update requires candidate ancestry to write
done, and internal/integrate writes a separate done commit after merging.
Context, show and direct-file reads expose the stored status. Keeping
status=review while only the board says Done would mislead those consumers.

Two alternatives were examined: explicitly close work with a post-delivery
metadata commit (through a follow-up PR where the target requires it), or
redesign the record contract and derive completion. The owner proposed the
two-step mechanism to test its viability, then raised its operational cost.
It is feasible but was not selected. Its historical assessment is not a
requirement on the implementation.

## Answer

On 2026-09-29 (owner's local date), the owner explicitly selected:
"Derive Done with a redesigned record contract."

The consequential choice is recorded in accepted decision
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md).
Persist candidate acceptance; derive Done from applicable acceptance and
verified delivery. Remove the misleading permanent review representation,
rather than overlaying a contradictory UI label. No mandatory follow-up
completion commit or PR is required by the selected model.

The owner also asked to finish preparation up to implementation. This
resolves the completion-model question; it does not claim acceptance of
all proposed interface or screen details. The
[contracts](G-260930-84fnb-portable-workflow-contra.md) and
[experience](G-260930-npw49-portable-workflow-experi.md) are reconciled around
this answer. The owner subsequently reviewed and accepted both designs at
f376a6d; that acceptance is recorded in
[G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md). No product code
has changed.
