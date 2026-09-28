---
id: "G-260928-dbgbw"
type: decision
title: "R on a member of an ended selection relaunches the selection"
status: accepted
created: "2026-09-28T20:54:54Z"
updated: "2026-09-28T20:55:10Z"
relates_to: ["G-260928-63124", "G-260925-wc2pz", "G-260925-7c8g9"]
---

## Decision

On 2026-09-28 the owner decided, in the session implementing
[G-260928-63124](G-260928-63124-show-each-member-s-state.md): `R` on a
member of an ended selection whose current state still stands on the
selection's branch relaunches that selection, the same IDs on the same
branch and worktree, as `grove run` with those IDs would. It is refused
while a started member still waits on an open question, naming the
question to answer first.

This follows [G-260925-wc2pz](G-260925-wc2pz-review-an-explicitly-sel.md):
a started member left incomplete holds the whole branch out of review, so
handing one member off alone would carry another member's unfinished code
with it. The relaunch does not redo members whose checkpoint the branch
confirms. On Enter the board checks every member's record in its checkout
against what it read, as `--expect` does for one work.

## Alternatives

- **Launch the member alone** (the earlier single-work `R`): simple, but
  on the selection's branch it would hand that member off with the other
  members' unfinished code, which G-260925-wc2pz rules out.
- **Refuse `R` and send the owner to `grove run ID...`**: safe, but the
  board then cannot continue the work it shows as waiting on it.

## Reconsideration

Reconsider when a selection can hand off part of its members, or when
parallel implementation is selected.
