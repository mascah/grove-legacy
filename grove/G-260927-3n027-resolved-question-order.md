---
id: "G-260927-3n027"
type: work
title: "Shaping keeps the order a resolved question was holding"
status: done
created: "2026-09-27T19:15:26Z"
updated: "2026-09-27T20:11:27Z"
kind: fix
size: small
relates_to: ["G-260925-g39ga"]
candidate: "58a7855687c5ec319b0e715ab12ddf701b5834a1"
approved: "58a7855687c5ec319b0e715ab12ddf701b5834a1"
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
  is not this work. It also made the redesign depend on the content refresh
  G-260927-5gh2k, an order no question held: the session had written "Independent
  of the redesign, so it can ship first" although the redesign restyles what
  the refresh adds. That is a prerequisite missed while shaping, which the
  guide's `depends_on` paragraph already covers; this work does not address it.

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

## Evidence

On `worktree-G-260927-3n027`, from base `main` `64202ba` and record revision
`sha256:73a27b541f62`. No plan: one passage of documentation. Implemented in
`c36c5e0`, with the review fix in `8f9dfd0`. The candidate is the commit that
records this Evidence.

- **Text.** `grove guide shape` prints, at the end of the paragraph on
  answered questions: "The resolved question's `blocks` no longer holds
  anything back. When the answer was also meant to wait for other work, such
  as a spike, and that work must still be delivered first, add it to the
  blocked work's `depends_on` and give the reason in that work's body, so the
  order does not survive only in prose." Applied to the ascah.dev case, it
  yields the edge that session missed: the redesign depends on the spike.
- **Same test.** "Must still be delivered first" is the `depends_on`
  paragraph's test; a preferred sequence still fails it. `grove guide model`,
  docs/commands.md and docs/board.md needed no change: `blocks` still never
  orders.
- **Shipped-document rule.** The passage adds no link, record ID or
  repository name.
- **Final checks.** At `8f9dfd0`: `go vet ./...` passed, `gofmt -l .`
  printed nothing, `go run ./cmd/grove check` printed `OK: 204 records`, and
  `go test -count=1 -timeout 120s ./...` passed every package. The same
  checks passed at `c36c5e0`. No code or TUI change.
- **Review.** [G-260927-e6p0r](G-260927-e6p0r-review-of-g-260927-3n027.md):
  two rounds by one `grove-reviewer` agent, resumed for round 2 rather than
  fresh. Round 1 found the example clause hard to parse, fixed in `8f9dfd0`;
  round 2 ended `Open findings: none`.
- **Limits.** A question the owner resolves from the board (`y`) involves no
  guide, so the lifted order is caught only when a session next reconciles
  the answer. The installed `grove` carries the old guide until `just
  install` after the merge.

## Next

In review. Then:

```sh
# in .claude/worktrees/worktree-G-260927-3n027
go run ./cmd/grove approve G-260927-3n027 "VERDICT"
# in the main checkout
go run ./cmd/grove integrate G-260927-3n027 --cleanup
```

Or `go run ./cmd/grove feedback G-260927-3n027 "TEXT"` in the worktree.

Verdict on candidate 58a7855, 2026-09-27: approved
