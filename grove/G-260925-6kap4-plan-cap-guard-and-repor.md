---
id: "G-260925-6kap4"
type: review
title: "G-260924-59f5k plan-cap guard and report: independent review"
status: current
created: "2026-09-25T04:01:37Z"
updated: "2026-09-25T04:01:49Z"
work: ["G-260924-59f5k"]
examined: "4513c69"
---

## Examined

Three rounds by fresh `grove-reviewer` agents, 2026-09-25, of the changes
made after the owner stopped the G-260924-59f5k attempt, range `1e1a345..4513c69`
on `worktree-G-260924-59f5k`: the plan-cap guard in `evals/run.py` and its README
section, [G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md),
[G-260925-beby3](G-260925-beby3-keep-a-blank-mandate-ans.md), the report
[G-260925-42j50](G-260925-42j50-codex-eval-row-pattern-p.md) and
[G-260925-ced1h](G-260925-ced1h-give-adopting-projects-t.md). Each ran `python3
evals/run.py selftest` and `grove check` and read the retained runs under
`~/.cache/grove-evals` without writing; none ran a paid command. Round 1
examined `48358b7`, round 2 `1e1a345..643d29d`, round 3 `4513c69`, the
`examined` commit.

## Findings

Round 1 (`48358b7`), no blocking defect: the first run was gated on a
possibly stale reading with no per-run allowance; a stale prior reading
could be counted as the first run's use; a reading without `resets_at`
read as a reset window; use across a reset could go negative; the stop
message was inexact; README and G-260925-04ccr said 8 to 12 points where the
rollouts show 6 to 13. Fixed in `f37be4e`: a 13-point floor per run,
same-window subtraction clamped at 0, a missing `resets_at` read as
current, the README saying the cap is not a ceiling. The guardian rollout
picked by mtime was checked safe on all nine pairs.

Round 2 (`1e1a345..643d29d`), no blocking defect: a new home with no
reading skipped the floor on its first run; `resets_at` jitters by a second
within a window, so exact equality misjudged it; the selftest covered only
the happy path and the stop; G-260925-42j50 understated the reads outside the clone
and named one fixture commit for all batches; G-260925-ced1h said the binary embeds
only the guides. Fixed in `4513c69`, with selftest checks for the new-home
floor, a reset window and a missing `resets_at`. G-260925-42j50's numbers, the
pattern per case and the Claude pair were verified against the evidence.

Round 3 (`4513c69`): every earlier finding resolved; nothing blocks.

## Disposition

Open, low, left at the review cap: no selftest for the stop after a run
that leaves no reading or for the clamp at 0; a run whose readings lack
`resets_at` is measured from its own first reading (conservative for the
gate, the floor bounds it, no retained reading lacks it); G-260925-42j50 calls all
of `~/.cache/grove-evals` the eval's output where one search was of the
Codex home; one README line is longer than its neighbours.
