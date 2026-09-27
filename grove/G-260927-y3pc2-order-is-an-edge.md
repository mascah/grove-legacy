---
id: "G-260927-y3pc2"
type: decision
title: "depends_on records every order worth stating"
status: accepted
created: "2026-09-27T19:57:20Z"
updated: "2026-09-27T20:04:44Z"
relates_to: ["G-260925-g39ga", "G-260925-wc2pz", "G-260927-k4xwq", "G-260927-3n027"]
---

## Decision

`depends_on` names the work that should be delivered before this work: this
work needs its result, or building it first or alongside would redo that work
or conflict with it, as a sweeping removal or rename does for work that edits
the same code. An order among work records is made only by such an edge,
with its reason in the dependent work's body, because the board's `g` and
`grove deps` read nothing else. Prose in Next, the brief or a session's
return may repeat an order the edges show, never one they lack, and a claim
that work can proceed in parallel needs evidence named in a record. What is
weaker than an edge is no order at all: an importance is `priority`, a
grouping is `members`, and work that only touches the same files in
separable places needs no edge. A brief may still lay out direction for
outcomes not yet shaped as work; once they are work records, the edges carry
their order.

Decided by the owner on 2026-09-27. Two shaping sessions in a row wrote a real
order only in prose, which `g` then showed as unrelated work. The owner said:
"This is a problem we need to have stop happening, I can't be baby sitting and
triple checking whether the agent bothered to setup the records correctly." Of
the three options put to them, they chose this one.

After acceptance, review round 1 of G-260927-k4xwq led the session to reword
the first paragraph, and the owner has not yet read the new text. The owner
accepted "never a sentence in Next, the brief or a session's return". The
session changed that to "may repeat an order the edges show, never one they
lack", added that a parallel claim needs evidence, and added the sentence on
unshaped outcomes in a brief. The owner's approval of G-260927-k4xwq's
candidate confirms this wording.

It replaces the stricter reading, in which `depends_on` held only work that
"must be delivered before this work can proceed" and a preferred sequence
went in Next. That reading came from G-260925-g39ga, whose acceptance item 7
kept "preferred sequence identified separately". G-260925-wc2pz's "A
dependency edge still means only that the prerequisite is needed" now reads
with this wider meaning of needed. Its point, that an edge authorizes no
shared execution, stands.

This repository's own brief, in its "Suggested sequence", still calls the
order of three independent proposals "proposed". That section is the
owner's to change.

## Alternatives

- **Keep strict prerequisites and state no preferred order at all.** Rejected:
  the order is real. In keyborg, stripping the game layer rewrites comments in
  `passages.ts`, which the corpus proposal replaces, and in the header of
  `identify.ts`, whose thresholds the mastery proposal replaces. Dropping
  that order would hide it, not remove it.
- **Keep `depends_on` strict and add a soft field that `g` shows but `run`
  does not wait on.** Rejected for now: it adds a schema field, code in
  `deps`, the board and `check`, and one more thing for an agent to classify.
  That classification is the step that failed twice.

## Reconsideration

Reopen if the wider edges serialize work that could have run alongside often
enough that `run`'s waits cost more than the conflicts they prevent. Also
reopen if a shaping session still leaves an order in prose after the guide
says this; then the fix is enforcement, not wording.
