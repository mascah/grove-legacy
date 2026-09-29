---
id: "G-260929-1w9x4"
type: review
title: "Does a plan from one session suffice for a cheaper implementer?"
status: current
created: "2026-09-29T03:33:37Z"
updated: "2026-09-29T03:34:08Z"
work: ["G-260928-pqhyg"]
---

## Examined

The plan-sufficiency experiment of
[G-260928-pqhyg](G-260928-pqhyg-measure-whether-a-plan-f.md), run
2026-09-29 with grove `9a18f57` and Claude Code 2.1.284, on three items with
the settings G-260929-jz2gy, G-260929-2nfwv and G-260929-z5nec answered:
plan on Opus 5.5 at `xhigh` as `grove run ID --until plan`, then an
implementation on Sonnet 5.5 that read only the plan and the record.
Examined: each item's attempt facts (`grove attempt`), its branch, its
record's Evidence and its review. This is the parent session's report, not an
independent review; the owner's verdict on each candidate is pending.

| Item | Plan (Opus, xhigh) | Implementation (Sonnet) | Review | Candidate |
| --- | --- | --- | --- | --- |
| G-260928-kehya (here, small) | $1.73, 4 min 41 s, 34 tool calls, 142 lines | subagent of an interactive Opus session, effort unset: 42 tool calls, 4 min 58 s, 118k tokens; fix round 8 tool calls, 67 s; cost not separable; plus $0.27 real-provider check | `grove-reviewer` on Opus from the parent, 1 round, 34 tool calls, 5 min: 3 low, fixed | `63ea818`, review |
| keyborg G-260928-ej7j0 | $2.98, 9 min 23 s, 42 tool calls, 212 lines | `grove run`, `high`: $5.02, 80 turns, 110 tool calls, 26 min 21 s | `grove-reviewer` on Sonnet inside the attempt, 2 rounds: 3 findings, one a real bug ("91 of 90 days"), fixed; `Open findings: none` | `62e2911`, review |
| keyborg G-260928-g133n | $1.75, 5 min 4 s, 31 tool calls, 140 lines | `grove run`, `high`: $2.45, 54 turns, 83 tool calls, 26 min 59 s, first code commit at 10 min 46 s | `grove-reviewer` on Sonnet, 2 rounds: 3 low (keyboard focus, a missing test), fixed, the last not re-reviewed | `3d41e96`, review |

Plan plus implementation: ej7j0 $8.00, g133n $4.20, kehya $2.00 plus the
subagent's unmeasured share. The sixteen-attempt mean in
[G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) is about $5
per whole attempt on Opus; these are other items, so the figures place the
split, they do not measure a saving.

## Findings

1. **A plan sufficed for all three.** Each Sonnet implementation reached a
   reviewed candidate without a question or a stop, and no review finding
   traced back to what the plan said. The one real bug (ej7j0's day count)
   was the implementer's and the review caught it.
2. **What the plans lacked was detail, not direction.** kehya's implementer
   reported about ten silences, each settled as a routine choice the
   reviewer then accepted: which field to compare; where a check sits
   relative to the data it needs; every place a rule is enforced (a refusal
   the plan put in one place lived in three); the format of each
   user-visible line; fixture prerequisites for a test; how to derive a
   baseline the plan asked for; and the state of a related record it was
   told not to read (G-260928-y2p5h). The keyborg items made their own
   additions (one tab stop instead of 25, meter widths, a merge of `main`)
   and recorded them. A sufficient plan names every touch point and each
   visible string, and carries the facts the implementer will not look up.
3. **The subagent route is fast and keeps the parent trim, and is thin on
   facts.** The parent saw one report, not the implementation, but the
   subagent's cost cannot be separated from the parent's, its effort cannot
   be set, "read only the plan" is an instruction rather than a bound, and
   under auto mode a headless parent was denied the write into the sibling
   worktree (G-260929-z5nec), so it needs an interactive parent.
4. **The separate process has every fact and misses reviewer routing.**
   `grove run` bounds cost, records turns, tool calls and duration, and
   outlives the terminal, but the reviewer inherits the implementer's model:
   the owner's Opus reviewer could only be had through an untracked
   settings file, which auto mode refused to write, so both keyborg reviews
   ran on Sonnet.
5. **Limits.** Three items; reviewers on different models; the Sonnet
   attempts took about 26 minutes each against 5 for the subagent, on
   larger items; no owner verdict yet.

## Disposition

Recommendation, for the owner, not a decision: shape the separate process
first. `grove run` already takes a model and effort per phase through
`--until plan`; what it lacks is a reviewer route of its own (a reviewer
model on `run`, or the review injected as a prompt rather than an agent
definition, which the owner raised in G-260929-2nfwv). Leave the subagent
route unshaped until a subagent's cost can be attributed. Whichever route,
the work guide's plan step should require the contents in finding 2.
Judge the three candidates before acting on this.

Open findings: none
