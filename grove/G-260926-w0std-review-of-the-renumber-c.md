---
id: "G-260926-w0std"
type: review
title: "Review of the renumber command and the rename"
status: current
created: "2026-09-26T16:59:12Z"
updated: "2026-09-26T16:59:27Z"
work: ["G-260926-vkv48"]
examined: "24a43fae8cb179669fb6270f493f701ea227b41f"
---

## Examined

Work [G-260926-vkv48](G-260926-vkv48-rename-legacy-records-to.md), plan
[G-260926-yzd8f](G-260926-yzd8f-plan-for-renaming-legacy.md), on branch
`worktree-G-260926-vkv48` from `a9f8fce`. Three rounds by fresh
`grove-reviewer` agents, read-only, with the map of the run
(`/tmp/vkv48-map.keep.jsonl`) and the pre-run attempt store copy
(`/tmp/vkv48-attempts-backup`):

1. `a9f8fce..ac69921`: the command, its test and docs, the run, the
   reference repair.
2. `ac69921..502d4bf`: the round 1 fixes, spot-checking the whole candidate.
3. `502d4bf..24a43fa`: the round 2 fix.

## Findings

Round 1 (4 open):

1. Medium. `Rewrite` treated a JSON escape before an ID (`\nG-153`) as a
   letter joined to it, so about 2,090 legacy IDs stayed in the attempts'
   `events.jsonl`. Fixed in `3bdcf25`: a letter or digit after a backslash
   no longer joins the ID. The regression fixture fails without the fix. The
   live store was repaired by a one-off pass with the same rule and map,
   skipping the running attempt (2,089 tokens). Round 2 applied the fixed
   rule to the backup and got the live store byte for byte.
2. Medium. Nullsec's own numbers were left rewritten in two places
   (G-260925-dzxm6 lines 230-233, G-260921-905y3 line 329). Restored to the
   base text in `3bdcf25`; round 2 found no other nullsec context rewritten.
3. Low. The record model still said legacy IDs are never renumbered.
   Reworded in `3bdcf25`.
4. Low. Nothing said to copy the attempt store, which is not in Git, before
   running. Added to the record model, `--help` and the incomplete-run error
   in `3bdcf25`; a CLI test of the map line, rerun and usage error came with
   it (`502d4bf`).

Round 2 (1 open): Low. The map page G-260921-czt8x overstated what the attempt store
keeps: the running attempt of this work was skipped and quotes the legacy
IDs, and a scratch filename `G-107tail.jsonl` remains. Reworded in
`24a43fa`.

Round 3: the sentence verified against the store; nothing else changed.

Knowledge (round 1): no new domain concept, no conflict with a settled term;
the slug rule that leaves cited IDs out, and the renaming of G-260926-afe5w,
stay within the decision's choice 1 and are recorded in the plan, for the
owner's judgment under acceptance 8.

## Disposition

Every finding fixed. Open findings: none
