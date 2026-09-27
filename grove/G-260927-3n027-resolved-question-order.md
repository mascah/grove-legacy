---
id: "G-260927-3n027"
type: work
title: "Shaping keeps the order a resolved question was holding"
status: active
created: "2026-09-27T19:15:26Z"
updated: "2026-09-27T19:16:07Z"
kind: fix
size: small
relates_to: ["G-260925-g39ga"]
---

## Outcome

When a session records the answer to a question whose `blocks` also held its
work behind other work, that order survives the answer as `depends_on`, so
`grove deps` and the board's `g` still show it once the gate lifts. The owner
first: after a shaping session in ascah.dev on 2026-09-27 they opened `g`,
saw no work depend on any other, and asked whether it could all truly be done
in any order. It could not.

## Constraints

Observed in ascah.dev at `c1afc99`, 2026-09-27, from that session's
transcript and records:

- Its question G-260927-pgdn2 (the diorama pipeline) blocked the redesign
  G-260927-kkgke, and its Next said the owner would answer "after seeing the
  spike" G-260927-htwdj. The gate was also standing for spike-before-redesign.
- The owner answered before the spike. The session resolved the question,
  recorded the accepted decision, and rewrote the redesign's Next to
  "Unblocked. G-260927-htwdj first, to see the look". No `depends_on` was
  set anywhere, and `grove deps` listed the four work items as four groups at
  layer 0; selected together, `run` would have taken them in the order given.
- The shaping guide at `64202ba` says to set `depends_on` only for work that
  must be delivered first and that a preferred sequence is not a
  prerequisite ([work-shaping](../docs/work-shaping.md), the paragraph on
  `depends_on`); its paragraph on answered questions says nothing about the
  lifted `blocks`. The session read "the spike first" as a sequence.
- ascah.dev's records were repaired separately in that checkout; that repair
  is not this work.

In scope: the shaping guide's paragraph on answered questions. Out of scope:
code (`blocks` still never orders, and `deps` and `g` are unchanged), the
work guide, whose sessions meet an answer on work already being executed,
and the record model. The guide ships in the
binary, so its text names no record and no repository.

## Acceptance

- `grove guide shape`, where it covers answered questions, tells a session
  that a resolved question's `blocks` no longer gates, and to set
  `depends_on`, with the reason in the dependent work's body, when the gate
  also held work behind other work that must still be delivered first.
- The added text uses the same test as the `depends_on` paragraph (must be
  delivered first), so it does not turn a preferred sequence into an edge.
- The guide still links only within itself or to `https://` and names no
  record or repository; `go test` passes.

## Next

Active on `worktree-G-260927-3n027`.
