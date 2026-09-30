---
id: "G-260930-tcc9w"
type: decision
title: "Bound delivery groups and retain cheap completion reads"
status: accepted
created: "2026-09-30T20:07:37Z"
updated: "2026-09-30T20:22:01Z"
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260925-wc2pz", "G-260930-gj9d7", "G-260930-84fnb", "G-260930-npw49", "G-260930-0s29t", "G-260930-yfh91", "G-260928-c5j9d"]
---

## Decision

On 2026-09-30, after reviewing the churn in
[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md), the owner asked
to reconcile the product requirements, design and existing work before
continuing implementation or migrating other projects to schema 4.
In this conversation they selected:

- **Both delivery modes.** A selection defaults to separate deliveries:
  implement, review, approve and deliver each item before starting a
  dependent item from the updated target. An explicit deliver-together
  selection runs its members sequentially on one branch and delivers one
  shared candidate, approved against each member's acceptance.
- **Unattended progression.** One authorized launch continues its selected
  work under the configured policy and aggregate limits. Waiting for a
  prerequisite's delivery does not require a human act when delegated
  approval and integration can complete it. It does not authorize adding
  other work, expanding scope, or spending beyond those limits.
- **One workspace per delivery.** External prerequisites must already be
  delivered. Within a combined delivery, members may build on one another.
  Before delivery, interruption and review fixes resume that workspace.
  After delivery its execution role ends; the next delivery starts fresh
  from the updated target. The owner explicitly accepted: "Repeated
  delivery from an old squashed branch explicitly not supported."
- **Automatic cleanup, optional keep.** Remove a completed Grove-owned
  branch/worktree only when its evidence is retained, no process owns it,
  and it contains no additional committed, uncommitted or untracked work.
  Keeping it permits inspection, not another delivery from it.
- **Verify on delivery, read cheaply, audit on request.** Retain
  [G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md)'s reading
  rule: applicable acceptance in the target's record establishes ordinary
  Done; original candidate/review/approval evidence survives cleanup for
  inspection and optional audit. Missing audit objects in a clone do not
  reopen work. The owner is currently comfortable with a hand-written
  accepted record on the target reading Done until an audit exposes the
  missing correspondence. Supported delivery operations must prevent it.
- **Include the LLM approval judge in the portable milestone.** Asked
  whether it should remain deferred, the owner answered "Include the LLM
  judge". Judgment remains a separately configured, bounded condition of
  the owner's policy, alongside independent review and deterministic checks.

This amends the universal shared-candidate boundary of
[G-260925-wc2pz](G-260925-wc2pz-review-an-explicitly-sel.md), the kept-branch
continuation clause of G-260930-gj9d7, and the judge deferral in
[G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md).
It does not select parallel execution, a general workflow builder or a
new completion schema. The current binary's contracts remain documented
until the replacement work is delivered.

The [brief](brief.md) owns direction, the
[contracts](G-260930-84fnb-portable-workflow-contra.md) own the reconciled
design and its remaining implementation proposals, and the
[experience](G-260930-npw49-portable-workflow-experi.md) owns the journey.
The [milestone](G-260930-60c3d-complete-the-portable-gr.md) owns work
membership, trial evidence and the current adoption hold.

## Alternatives

- Separate deliveries only would lose the owner's useful combined-work
  option. Combined delivery only makes an unfinished later member hold back
  earlier work and grows review scope. Two explicit boundaries share one
  delivery operation; arbitrary mixtures or automatic regrouping are not
  required.
- Supporting work stacked across independently managed unmerged branches,
  and repeated deliveries from kept branches, preserves flexibility at the
  cost demonstrated by the churn review. The owner chose a narrower
  supported execution workflow. Reading unrelated branches still works.
- Re-proving historical delivery on every read repeats the rejected cost.
  Discarding original evidence would lose inspection after cleanup. Cheap
  reading plus retained optional audit evidence preserves both needs.
- Leaving the LLM judge beyond the adoption milestone would postpone testing
  the unattended judgment the owner expects of the complete workflow.

## Reconsideration

Reopen a boundary only on a concrete workflow the owner needs and the cost
of supporting it. Reconsider the trust boundary if manual target metadata
or external integrations commonly make ordinary Done misleading; bring
observed cases and read costs, not a speculative universal verifier.
Reconsider default separate delivery if measured review cost outweighs its
smaller reviews and independent progress. Trial both modes before extending
the supported state space.
