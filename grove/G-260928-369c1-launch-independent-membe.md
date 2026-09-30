---
id: "G-260928-369c1"
type: work
title: "Launch independent members of a selection as parallel attempts"
status: proposed
created: "2026-09-28T19:28:59Z"
updated: "2026-09-30T01:16:10Z"
kind: feature
size: medium
depends_on: ["G-260928-63124", "G-260930-60c3d"]
relates_to: ["G-260925-7c8g9", "G-260925-wc2pz", "G-260925-g39ga", "G-260925-h8rj5", "G-260925-5wrn8", "G-260928-c5j9d", "G-260930-e8jj7", "G-260930-60c3d"]
---

## Outcome

The owner can launch a selection whose members are independent as one
attempt per member, each in its own worktree with its own candidate, so one
member's question or failure holds none of the others, and the candidates
are integrated one by one under the policy or by the owner. A selection
launched without that option stays one attempt on one branch, the chain, as
today.

Owner intent, shaping conversation 2026-09-28: "chain default, parallel
optional"; and, from keyborg, "we don't have a good handle on understanding
what work can happen in parallel versus sequential, and how grove can manage
that."

## Portable milestone framing

On 2026-09-29 the owner selected [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md). This remains proposed
work outside [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).
Parallel launch increases throughput once a person can reliably understand
and continue one complete change. Its default chain and optional independent
attempts remain the intended behavior.

Depends on [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md): the owner selected demonstration of the
complete usable loop before this investment. It also consumes the resulting
provider, work ownership and delivery contracts;
building against the current interfaces would risk redoing that integration.
The existing prerequisite G-260928-63124 supplies per-member state and remains in depends_on.

The technical design below is a dated proposal. At assignment, reconcile it
with the delivered portable contracts and capability-specific limits;
a summed time allowance must not be presented as a monetary cap.

## Constraints

Observed 2026-09-28 at main `6fbb888`:

- `run ID...` is one process, one worktree, one budget, members implemented
  in order, handed off together on one shared candidate
  ([G-260925-7c8g9](G-260925-7c8g9-execute-an-explicitly-se.md); decision
  [G-260925-wc2pz](G-260925-wc2pz-review-an-explicitly-sel.md), whose
  reconsideration names "when parallel implementation is selected"). In
  keyborg one member's question held two complete members out of review
  ([G-260928-63124](G-260928-63124-show-each-member-s-state.md)).
- `deps` orders by `depends_on` only; equal layers are "not evidence of
  parallel" (`docs/commands.md`). Independence evidence Grove holds without
  a model: no edge among the members, no prerequisite that waits, the
  plans' named files where plans exist, the overlap of the code paths each
  record names (the board's `described by` tiers), and the merge prediction
  of the resulting candidates
  ([G-260925-h8rj5](G-260925-h8rj5-predict-whether-a-candid.md)), which
  exists only after they are built.
- A running attempt refuses another launch that shares a member; launches
  of disjoint members are allowed today, one command each. The board's `g`
  preview shows a selection's order and outside prerequisites
  ([G-260925-g39ga](G-260925-g39ga-see-work-dependencies-an.md)). The brief:
  "Bound batches and fan-out before adding schedules."

Proposed design, labelled proposed:

- `run ID... --parallel` and the same choice on the board's dependency
  preview: refused when any member needs another in the selection or waits;
  otherwise one ordinary attempt per member, `worktree-ID` each, each with
  its own budget from `run:`, listed and tagged as any attempt. The preview
  and `--dry-run` show the overlap evidence (paths the records and plans
  name in common) and the summed budget; the owner decides on that
  evidence, and Grove refuses nothing on overlap alone.
- After the candidates hand off, the Review column shows their predicted
  merge order (`deps` `merge_order`) and the policy integrates them in ID
  order, resolving conflicts as
  [G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md) does; the owner
  does the same by hand.

Out of scope: an LLM judgment of independence, which a later record may add
on the mechanism of [G-260928-c5j9d](G-260928-c5j9d-judge-a-candidate-agains.md),
since the owner chose deterministic evidence and the chain by default;
parallel implementation inside one attempt; any change to the chain's
shared-candidate rule.

## Acceptance

1. `run A B --parallel --dry-run` prints two launches with their worktrees,
   budgets and the overlap evidence; refused with the reason for an edge
   between A and B or a member that waits.
2. Launching starts two attempts; each ends in its own candidate; a question
   in one leaves the other's handoff untouched (fake provider).
3. The board's preview offers the parallel launch with the same evidence
   and the summed budget.
4. `docs/commands.md` (Attempts, Dependencies), `docs/board.md` and the work
   guide's invocation row say so; the answer to G-260925-wc2pz's
   reconsideration is noted in this record, not a change to that decision.
5. One real trial on two independent items, here or in keyborg, under a
   budget the owner names at assignment.

## Next

Needs [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md) and the existing prerequisite named above.
At assignment, reconcile the dated mechanism with the delivered contracts,
retain the original bounded-trial acceptance, and record the owner's
execution mandate. This record is not a milestone member or assigned work.
