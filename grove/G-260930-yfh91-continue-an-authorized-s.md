---
id: "G-260930-yfh91"
type: work
title: "Continue an authorized selection across separate deliveries"
status: proposed
created: "2026-09-30T20:07:53Z"
updated: "2026-09-30T20:22:02Z"
kind: feature
size: large
depends_on: ["G-260930-0s29t", "G-260928-y2p5h", "G-260928-n4f1q", "G-260928-c5j9d"]
relates_to: ["G-260930-tcc9w", "G-260930-60c3d", "G-260930-84fnb", "G-260930-npw49"]
---

## Outcome

One authorized selection can run to its delivery boundary without the owner
manually relaunching each item. Separate delivery is the default; deliver
together is explicit. Policy and the configured LLM judge can approve and
deliver eligible work so its dependents continue from the updated target.

Owner decision:
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md), 2026-09-30.

## Scope and constraints

Implement the finite selection and continuation contract in the
[contracts design](G-260930-84fnb-portable-workflow-contra.md). Preview binds
the selected IDs, mode, dependency order, source revisions, authority,
cleanup preference and aggregate resource limits. Default separate mode
uses one item per delivery; together mode uses one shared candidate judged
per member. Do not add arbitrary partitions, automatic regrouping, parallel
fan-out, scheduling, or dependencies outside the authorized selection.

The on-demand owner coordinates execution, independent review, applicable
policy approval and delivery. It advances after observing delivery, including
delivery completed externally or just before a crash. It never interprets
provider success as approval, silently grants merge authority, or launches
a second owner. Native resume is optional; new delivery workspaces use
durable context. Share operations with CLI/TUI, not another work-status store.

Persist enough selection progress and remaining authority to reconcile a
restart. Work records and the target determine delivered facts; local
operation state only coordinates the launch. Missing resource accounting or
authority after recovery requires an explicit wait rather than a fresh
budget. Stop and aggregate exhaustion stop further launch. A member blocked
on a human choice holds its dependents; other ready selected members may
proceed sequentially within the same mandate. Inspection never wakes work.

No automatic splitting of a partially implemented combined group and no
unbounded retry or repair loop. Existing bounded review/resolve policies
remain explicit. A clean independent review is evidence for judgment, not
approval itself. The LLM judge is configured as a policy condition and
cannot waive deterministic or host requirements.

## Acceptance

1. Preview and execution expose both modes with separate as default. A
   three-item dependency chain delivers A before B begins, then B before C,
   using fresh target-based workspaces without another launch from the owner.
2. Together mode runs the same selected members in order on one workspace,
   reviews the combined candidate, checks each acceptance, and delivers
   once. Every selected member must finish: an incomplete or never-started
   blocked member holds the group's delivery while ready members may execute.
   No partial group is delivered; changing membership requires a reconciled
   owner mandate. Older assignments retain their recorded semantics.
3. Policy-approved work advances automatically; a judge wait or missing
   approval authority names the reason and does not spin. Human acceptance
   or answering a question permits an explicit continuation under the
   existing mandate. No outside item is added.
4. Duplicate requests, owner loss, stop, exhausted limits, changed inputs,
   target movement and interruption immediately after delivery are covered.
   Recovery does not redo delivered work, duplicate an owner or reset the
   aggregate budget. Cleanup failure cannot turn successful delivery into
   an implementation retry.
5. A small local two-item walkthrough, using the existing primary UI/CLI,
   demonstrates mode selection, delegated approval, an actual delivery
   boundary, fresh context and a human wait. Record interventions and costs
   under an owner-assigned real-provider mandate. Do not defer this evidence
   to the hosted or final presentation work.
6. Guides, CLI/TUI preview, selection/attempt reporting and compatible old
   assignment handling describe the same boundaries. Existing assignments
   that named a shared boundary keep it; missing historical mode is not
   reinterpreted as separate. The combined cross-harness and fresh-clone
   evaluation remains G-260930-gwnb1's outcome.

## Dependencies

- G-260930-0s29t supplies workspace admission, retirement and cleanup.
- G-260928-y2p5h supplies portable execution and effective limits.
- G-260928-n4f1q supplies the independent review operation.
- G-260928-c5j9d supplies configured policy judgment and its waits.

## Next

Proposed and unassigned. Needs the declared prerequisites. Plan around the
two mode examples and recovery boundaries, then demonstrate the local
sequence before the complete local proof consumes it. No universal workflow
engine, new lifecycle schema or implicit spending mandate is selected.
