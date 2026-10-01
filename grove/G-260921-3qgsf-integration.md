---
id: "G-260921-3qgsf"
type: term
title: "Integration"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-30T20:22:12Z"
relates_to: ["G-260921-vr8a8", "G-260921-jatts", "G-260921-btyck", "G-260930-e8jj7", "G-260929-gm3m4", "G-260930-4742q", "G-260930-2qa4a", "G-260930-tcc9w"]
formerly: "T-007"
---

## Meaning

An accepted [candidate](G-260921-jatts-candidate.md)'s result reaching the
project's configured target. It follows
[approval](G-260921-btyck-approval.md) and establishes delivery for
[work](G-260921-vr8a8-work.md). Approval, passing checks and provider success
are separate facts; none alone is integration.

An ordinary merge can retain the candidate's identity. A squash transforms
it, so the delivery operation establishes correspondence with the accepted
result and retains the original candidate and judgment evidence. Trailers
locate that evidence; they are not proof by themselves.

## Relationship to standing

Schema 4 records acceptance and derives completion.
[Standing](G-260930-44q35-standing.md) reads the target's accepted copy of
the record cheaply, under the trust boundary selected in
[G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md). Delivery is
proved when made and optionally audited later. Missing historical evidence
limits the audit, not ordinary Done. Manual target metadata can make a false
receipt; the owner currently accepts that limit and requires supported
operations to prevent it.

Reverting a delivery including its record restores the earlier standing.
A later code-only change that leaves that receipt does not erase historical
delivery; corrective work is explicit. Legacy schema 3 claims retain their
historical meaning without fabricated evidence.

## Selected delivery boundary

[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md) selects one
workspace for a delivery, containing one item by default or an explicitly
combined group. Members of a group are judged against their own acceptance;
external prerequisites are delivered before execution. After delivery the
workspace is retired and normally cleaned up safely, with optional retention
for inspection. A subsequent delivery uses a fresh target-based workspace;
repeated delivery from an old squashed branch is unsupported in this selected
workflow.

The [contracts design](G-260930-84fnb-portable-workflow-contra.md) owns the
proposed mechanics. [Command documentation](../docs/commands.md) describes
implemented behavior: G-260929-gm3m4 delivered schema 4 and local squash;
G-260930-0s29t fresh target-based workspaces, retirement and default
cleanup; G-260930-4742q owns the hosted path. Historical ancestry-only mechanisms and kept-branch
recoveries remain in Git history, not a second definition here.
