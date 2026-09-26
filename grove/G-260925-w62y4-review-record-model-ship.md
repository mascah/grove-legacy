---
id: "G-260925-w62y4"
type: review
title: "G-260925-ced1h review: record model shipped as grove guide model"
status: current
created: "2026-09-25T04:36:42Z"
updated: "2026-09-25T04:36:45Z"
work: ["G-260925-ced1h"]
examined: "3b62bff"
---

## Examined

Two independent `grove-reviewer` rounds on `worktree-G-260925-ced1h`, each a fresh
agent, read-only, against [G-260925-ced1h](G-260925-ced1h-give-adopting-projects-t.md)'s
acceptance, [G-260925-02jsj](G-260925-02jsj-how-should-an-adopting-p.md)'s Answer and
the plan [G-260925-6ykpv](G-260925-6ykpv-plan-ship-the-record-mod.md):

- Round 1: `2d5777e..7b841f9`.
- Round 2: `2d5777e..3b62bff`, the fix commit read line by line.

`fa18712`, after round 2, applies round 2's two wording notes to
`docs/record-model.md` only (the unchanged-bytes sentence widened to every
failure before replacement, as G-260919-shnj5 states it; one line rewrapped).

Each reviewer ran `go test -short ./internal/cli`, `go vet ./...`,
`gofmt -l .` and `grove check`, all passing; neither ran the full suite or
a disposable `grove init` project, which the author ran (G-260925-ced1h Evidence).

## Findings

Round 1 (7 findings, none blocking):

1. Should-fix: dropping "G-260919-shnj5 owns update's … failure-reporting contract"
   left `update`'s refusal and applied-failure behaviour unstated; G-260919-5f89v's
   `flock` likewise.
2. Note: "that conversion" had no antecedent.
3. Note: "predecessor skills CLI" names a tool an adopting project does not
   know.
4. Note: the guides still mention Grove's records unlinked (G-260921-tkdwh, G-260921-9wkjt
   in work-execution.md's Lifecycle, G-260920-j2eyp in its Invocation, G-260921-ahbrj in
   work-shaping.md), which the new AGENTS.md rule's reason also covers.
5. Note: README's record-model row not reconciled, as the plan said.
6. Note: `guides.go`'s package comment stale.
7. Note: the test's `G-` allowlist would pass an unlinked provenance
   mention of an allowed number.

Round 2 (no blocking or should-fix): 1, 2, 5 and 6 resolved and checked
against G-260919-shnj5 lines 172-176, `internal/update/update.go`, G-260919-5f89v and
`internal/repo/repo.go`; two new notes: the unchanged-bytes sentence was
narrower than the contract, and one line exceeded the wrap.

## Disposition

- Fixed in `3b62bff`: 1, 2, 5, 6. Fixed in `fa18712`: round 2's two notes.
- Not changed: 3 (true and harmless, outside G-260925-02jsj's listed edits); 7
  (accepted: the link check and review cover it).
- 4 is left to the owner in G-260925-ced1h's Next: out of this plan's scope, and
  widening the rule to "names no record" is the owner's call.
