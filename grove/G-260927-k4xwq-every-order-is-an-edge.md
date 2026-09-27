---
id: "G-260927-k4xwq"
type: work
title: "Every order shaping states is a depends_on edge"
status: active
created: "2026-09-27T19:57:19Z"
updated: "2026-09-27T19:58:03Z"
kind: fix
size: small
depends_on: ["G-260927-3n027"]
relates_to: ["G-260927-y3pc2", "G-260925-g39ga"]
---

## Outcome

The owner can trust `g` as the whole account of order after a shaping
session: every order the session means is a `depends_on` edge with its
reason, and nothing in the records or the session's return claims an order
that the edges do not show. This is the owner's intent, 2026-09-27, under
decision G-260927-y3pc2: "I can't be baby sitting and triple checking
whether the agent bothered to setup the records correctly."

## Constraints

Observed, 2026-09-27, in two shaping sessions run with the installed `grove`
(`64202ba`):

- **ascah.dev.** The redesign's Next said "G-260927-htwdj first, to see the
  look", and no edge was set (G-260927-3n027 has the detail).
- **keyborg** (uncommitted proposals, transcript `ac1d0fb4`). The session read
  the whole shaping guide and followed its rule. G-260927-2qw40's Next says
  "Proposed first in the brief's sequence because it unblocks every other
  slice; no record depends on it". G-260927-5zhf5's says "Can run in parallel
  with G-260927-2qw40; it touches none of the deleted files". The strip also
  rewrites Bench-contract comments in about thirty files under `app/`,
  including `settings.ts`, `passages.ts`, `lessons.ts` and the engine types,
  which the other three proposals edit. The return said both "Unblocks
  everything else" and "None depends on another, so they can run in
  parallel". It ran `grove check` and never `grove deps`.
- The guide at `64202ba` tells a session to write a preferred sequence "in
  Next" (the paragraph on `depends_on`). Step 6 runs only `grove check`, and
  nothing asks the return's account of order to match the records.
- `grove deps` already prints "Equal layers have no declared order, which is
  not evidence that they can proceed in parallel." Neither session read it.

Needs G-260927-3n027 delivered first (`depends_on`): this work rewrites the
passage that G-260927-3n027 adds to the shaping guide, so the two conflict if
built alongside. The branch is stacked on G-260927-3n027's for that reason.

In scope: the record model's definition of `depends_on` and of member order;
the shaping guide's `depends_on` paragraph, its step 2 reference, the passage
G-260927-3n027 adds, and step 6. Out of scope: code, since `deps`, `run` and
the board already treat every edge alike; the work guide, whose sessions act
on edges rather than write them; the brief; and the owner's other checkouts,
whose records are the owner's to change. The guides ship in the binary, so
their text names no record and no repository.

## Acceptance

- `grove guide model` defines `depends_on` as work that should be delivered
  first, because this work needs its result or building it first or alongside
  would redo or conflict with it. Member order is presentation and orders
  nothing.
- `grove guide shape` says an order among work records is stated only as an
  edge with its reason, never in Next, the brief or the return. It also says
  what earns no edge: importance, grouping, and separable edits to the same
  files.
- Step 6 of `grove guide shape` runs `grove deps` over the work the session
  wrote. It checks every order the records and the return state against that
  output, allows no claim that work can proceed in parallel without evidence,
  and returns the order as `deps` shows it.
- Applied by hand to the keyborg proposals, the guide yields the strip as a
  prerequisite of the other three. Applied to ascah.dev, it yields the spike
  before the redesign.
- The shipped guides still link only within themselves or to `https://` and
  name no record or repository. `go test` passes.

## Next

Active on `worktree-G-260927-k4xwq`, stacked on G-260927-3n027's branch.
