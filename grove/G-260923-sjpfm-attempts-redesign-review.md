---
id: "G-260923-sjpfm"
type: review
title: "G-260923-895zb Attempts redesign review"
status: current
created: "2026-09-23T20:34:58Z"
updated: "2026-09-23T20:35:13Z"
work: ["G-260923-895zb"]
examined: "611f7ec"
---

## Examined

An independent reviewer subagent (read-only, headless attempt of
2026-09-23) reviewed [G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) in two rounds on
branch `worktree-G-260923-895zb`. Round 1 covered `0315df9`, the implementation,
against its parent. Round 2 covered `611f7ec`, the fixes, against `0315df9`.
Both rounds checked the code against G-260923-895zb's constraints,
[G-260923-hvnqh](G-260923-hvnqh-which-attempts-list-and.md)'s answer and the section
"Revised after G-260923-hvnqh's answer" in
[G-260923-rz01m](G-260923-rz01m-attempts-layouts-observe.md). The reviewer ran the
short tests and vet, probed widths 40, 41, 60, 99, 100 and 101 with `d` on
and off, and ran `ReadActivity` from the commit on all seven real
`events.jsonl` logs, comparing with counts from Python.

## Findings

Round 1 found no problem with terminal safety or widths, bounded reads, the
grouping rules apart from item 2, or the docs. It found:

1. Should-fix: Subagents was counted too high. Claude repeats `task_started`
   each time a background subagent is resumed, so G-260923-p5pt6 showed ≥3 against
   1 spawned.
2. Should-fix: a work in review whose current candidate differed from the
   attempt's still read "judge candidate", with the attempt's commit.
3. Should-fix: Turns took the largest `num_turns` when it should sum them.
   Each result event of a resumed session counts only its own query, so
   G-260923-fwakw showed ≥47 against a true 64.
4. Should-fix: Model read "not reported yet" when the window was cut, even
   though `result.json` held the model.
5. Nit: a failed read showed zeros as figures.
6. Nit: "spend known at the end" appeared for attempts that had ended
   without a result event.
7. Nit: every settled attempt got ✓.
8. Nit: the context window was the largest of the run's models.

Round 2 found all eight fixed and no regression. On every real log, Turns
and Subagents now equal Claude's own totals: G-260921-7trd7 122/1, G-260923-fwakw 74/2 and
64/3, G-260923-p5pt6 63/1, G-260923-q7tm6 39/1, G-260923-895zb 32/0. It raised two optional nits:

9. Nit: a log small enough to fit the window, but with more than 200
   activity rows, still showed `≥` on every count.
10. Nit: the State line read "is review" instead of "is in review".

## Disposition

- Fixed in `611f7ec`, each with a regression test: 1 (distinct `task_id`),
  2 (W-108 in `TestAttemptStandings`), 3 (summed results and `Ended`, shown
  as `≈N` until a result ends the run), 4, 5 and 6
  (`TestAttemptScreenHonesty`), 7 and 8.
- Fixed in `feb6b42` after round 2, without a third review, both small and
  covered by tests:
  - 9: a separate `Activity.Dropped` flag. `TestReadActivityBounded` now
    asserts it is set without `Cut`.
  - 10: the wording, with its test string updated.
- No finding is open.
