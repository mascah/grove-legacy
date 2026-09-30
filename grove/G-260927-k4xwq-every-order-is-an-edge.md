---
id: "G-260927-k4xwq"
type: work
title: "Every order shaping states is a depends_on edge"
status: accepted
created: "2026-09-27T19:57:19Z"
updated: "2026-09-30T13:18:01Z"
kind: fix
size: small
depends_on: ["G-260927-3n027"]
relates_to: ["G-260927-y3pc2", "G-260925-g39ga"]
candidate: "80c93fc1ba36a02dcaaa497f895c28eae2c646e4"
approved: "80c93fc1ba36a02dcaaa497f895c28eae2c646e4"
approved_by: owner
approved_context: "sha256:da038e0c5c04b91af080ffe15d04daa459b402bd90d35c3f69d1b63a0df5c06c"
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
  rewrites Bench-contract comments across `app/`. Among them are
  `app/content/passages.ts`, which the corpus proposal G-260927-wqbaj
  replaces, and the header of `app/lib/difficulty/identify.ts`, which
  describes the thresholds the mastery proposal G-260927-5zhf5 replaces. The
  focus-text proposal G-260927-yjz5k builds on `words.ts`, which the strip
  does not edit. The return said both "Unblocks
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
their text names no record and no repository. This repository's own brief
still states a proposed order in its "Suggested sequence"; decision
G-260927-y3pc2 leaves that section to the owner.

## Acceptance

- `grove guide model` defines `depends_on` as work that should be delivered
  first, because this work needs its result or building it first or alongside
  would redo or conflict with it. Member order is presentation and orders
  nothing.
- `grove guide shape` says an order among work records is made only by an
  edge with its reason. Prose in Next, the brief or the return may repeat an
  order `grove deps` shows, never one it lacks, and claims parallel work only
  with evidence named in a record. It also says what earns no edge:
  importance, grouping, and separable edits to the same files.
- Step 6 of `grove guide shape` runs `grove deps` over the work the session
  wrote. It checks every order the records and the return state against that
  output, allows no claim that work can proceed in parallel without evidence,
  and returns the order as `deps` shows it.
- Applied by hand to the keyborg proposals, the guide yields the strip as a
  prerequisite of the corpus and the mastery proposals, and no edge for the
  focus text, whose files the strip does not edit. Applied to ascah.dev, it
  yields the spike before the redesign.
- The shipped guides still link only within themselves or to `https://` and
  name no record or repository. `go test` passes.

## Evidence

On `worktree-G-260927-k4xwq`, from base `6c1b05a`, the tip of
G-260927-3n027's branch, and record revision `sha256:9ad75cb72cc7` at `5293e4c`.
Implemented in `76fa6dd`. The review fixes are in `4fe7155` and the
decision disclosure in `91d5922`. The candidate is the commit that records
this Evidence.

- **Model.** `grove guide model` defines `depends_on` as "Prerequisite work,
  which should be delivered before this work: this work needs its result,
  or building it first or alongside would redo that work or conflict with
  it". The relationship-fields list agrees. Member order "is presentation;
  only dependencies order work."
- **Shape rule.** Step 5 says only such an edge makes an order, with its
  reason in the dependent's body. Prose in Next, the direction document or
  the return may repeat an order `grove deps` shows, never one it lacks.
  Parallel work needs named evidence. Importance, grouping and separable
  edits to the same files earn no edge. Step 2 and G-260927-3n027's passage
  now use "should".
- **Step 6.** It runs `grove deps` with the written work's IDs and checks
  each order against `NEEDS` and each parallel claim against its evidence.
  The return carries the order as that output shows it.
- **Hand application.** keyborg, uncommitted there: `grove deps` shows the
  strip G-260927-2qw40 in `NEEDS` for G-260927-wqbaj and G-260927-5zhf5,
  each with its file reason in the body, and no edge for G-260927-yjz5k.
  ascah.dev at `f4ce9d4`: the spike precedes the redesign.
- **Shipped-document rule.** No link other than anchors and `https://`, and
  no record or repository name added.
- **Final checks.** At `91d5922`: `go vet ./...` passed, `gofmt -l .`
  printed nothing, `go run ./cmd/grove check` printed `OK: 207 records`, and
  `go test -count=1 -timeout 120s ./...` passed every package. The same
  checks passed at `76fa6dd`. There is no code change.
- **Review.** [G-260927-4m4j6](G-260927-4m4j6-review-of-g-260927-k4xwq.md):
  two rounds. Three findings from round 1 were resolved or left to the owner.
  Round 2's disclosure finding was addressed in `91d5922` and not
  re-reviewed.
- **Limits.**
  - The owner has not read the decision's reworded paragraph. Approving this
    candidate confirms it.
  - No behavioral evaluation was run, so the claim that sessions now
    set these edges is unmeasured. A case in `evals/` would measure it, at
    runs × budget dollars under a mandate.
  - This repository's brief still states a proposed order.
  - Sessions keep the old guide until `just install` runs after the merge.

## Next

In review. G-260927-3n027 integrates first, since this branch holds it.
Then:

```sh
# in .claude/worktrees/worktree-G-260927-k4xwq
go run ./cmd/grove approve G-260927-k4xwq "VERDICT"
# in the main checkout
go run ./cmd/grove integrate G-260927-k4xwq --cleanup
```

Or `go run ./cmd/grove feedback G-260927-k4xwq "TEXT"` in the worktree.

Verdict on candidate 80c93fc, 2026-09-27: approved

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
