---
id: "G-260927-y3pc2"
type: decision
title: "depends_on records every order worth stating"
status: accepted
created: "2026-09-27T19:57:20Z"
updated: "2026-09-27T19:58:02Z"
relates_to: ["G-260925-g39ga", "G-260927-k4xwq", "G-260927-3n027"]
---

## Decision

`depends_on` names the work that should be delivered before this work: this
work needs its result, or building it first or alongside would redo that work
or conflict with it, as a sweeping removal or rename does for work that edits
the same code. An order among work records is stated only as such an edge,
with its reason in the dependent work's body. It is never a sentence in Next,
the brief or a session's return, because the board's `g` and `grove deps`
cannot read one. What is weaker than that is no order at all: an importance is
`priority`, a grouping is `members`, and work that only touches the same files
in separable places needs no edge.

Decided by the owner on 2026-09-27. Two shaping sessions in a row wrote a real
order only in prose, which `g` then showed as unrelated work. The owner said:
"This is a problem we need to have stop happening, I can't be baby sitting and
triple checking whether the agent bothered to setup the records correctly." Of
the three options put to them, they chose this one.

It replaces the stricter reading, in which `depends_on` held only work that
"must be delivered before this work can proceed" and a preferred sequence
went in Next. That reading came from G-260925-g39ga, whose acceptance item 7
kept "preferred sequence identified separately".

## Alternatives

- **Keep strict prerequisites and state no preferred order at all.** Rejected:
  the order is real. In keyborg, stripping the game layer rewrites comments in
  about thirty files that the other three proposals edit. Dropping that order
  would hide it, not remove it.
- **Keep `depends_on` strict and add a soft field that `g` shows but `run`
  does not wait on.** Rejected for now: it adds a schema field, code in
  `deps`, the board and `check`, and one more thing for an agent to classify.
  That classification is the step that failed twice.

## Reconsideration

Reopen if the wider edges serialize work that could have run alongside often
enough that `run`'s waits cost more than the conflicts they prevent. Also
reopen if a shaping session still leaves an order in prose after the guide
says this; then the fix is enforcement, not wording.
