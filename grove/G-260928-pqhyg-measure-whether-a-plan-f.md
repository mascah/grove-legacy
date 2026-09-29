---
id: "G-260928-pqhyg"
type: work
title: "Measure whether a plan from one session suffices for another implementer"
status: active
created: "2026-09-28T19:29:00Z"
updated: "2026-09-29T00:37:15Z"
kind: investigation
size: small
relates_to: ["G-260924-5b6pz", "G-260924-59f5k", "G-260927-dx0yn", "G-260921-dqdde"]
---

## Outcome

Evidence, on real work, of whether a plan produced by a bounded preparation
attempt is enough for a separate implementer, at a cheaper model or effort,
as a subagent, or on another provider, to implement to acceptance without
the preparer's context, reported against the sixteen-attempt figures in
[G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md), so the owner can
decide whether to push implementation into subagents or separate processes
and what the work guide must say for that.

Owner intent, shaping conversation 2026-09-28: the owner chose incremental
orchestration and wants "to explore pushing more of the execution
phases/steps down into sub agents to maintain a trim parent context" with
different models per phase; on 2026-09-24 the owner expected preparation and
review to run at higher effort than implementation.

## Constraints

Observed 2026-09-28:

- G-260924-5b6pz's acceptance 5 designed this experiment (an Opus plan at
  `xhigh`, an Opus implementation at `medium`, the reviewer at `high`, a $30
  cap) and named [G-260924-59f5k](G-260924-59f5k-run-the-eval-pair-on-cod.md)
  as its target; that attempt was stopped after the gpt-6-astra incident and
  its Evidence reports no routing result. The experiment never ran. The
  coupling it named stands: routing implementation to a cheaper setting
  pays only if the plan suffices.
- Today's plans are written by the session that implements them. The plan
  bound, `--model`, `--effort` and the attempt facts exist; `attempts ID`
  totals cost; `Shape:` and `Changed:` give tool calls and the time to the
  first commit outside the record root
  ([G-260927-dx0yn](G-260927-dx0yn-retain-per-attempt-proce.md)).
- The work guide's step 5 does not say to implement through subagents. A
  subagent's only inputs are what the parent passes: the plan and the
  record.
- [G-260921-dqdde](G-260921-dqdde-preparation.md): a technical plan does not
  wait for human sign-off unless it needs a product choice.

Proposed method, labelled proposed: three real work items, here or in
keyborg, each run as a bounded preparation attempt at the owner's chosen
model and effort, then an implementation attempt at a cheaper setting that
reads only the plan and the record, with the review at `high`; for one of
them, the implementation as a subagent told to read only the plan, in an
interactive session since the guide does not say so. Report per item: plan
and implementation cost, duration, review findings and rounds, what the
plan lacked, and the owner's verdict. A review record holds the report and
ends with a recommendation for the guide. No product change.

## Acceptance

1. Three items run as above, with the settings and the cap the owner names
   at assignment, each recorded in its own record's Evidence.
2. A review record reports the comparison and answers: does a plan suffice,
   what must it contain, and which of subagent or separate process is worth
   shaping.
3. Nothing in the binary or the guides changes in this work.

## Evidence

Settings from [G-260929-jz2gy](G-260929-jz2gy-which-three-items-settin.md)'s
answer: items G-260928-kehya (here), keyborg G-260928-ej7j0 and
G-260928-g133n; plan on Opus 5.5 at `xhigh`, implementation on Sonnet 5.5 at
`high`, `grove-reviewer` on Opus 5.5 at `high`; $15 per plan attempt, $15
per implementation attempt, $90 in total. Which item's implementation runs
as a subagent was not answered: [G-260929-2nfwv](G-260929-2nfwv-which-items-implementation-subagent.md).

Plan attempts, 2026-09-29, each `grove run ID --until plan --model
claude-opus-5-5 --effort xhigh --budget 15`, launched together, Claude Code
2.1.284, grove at `9a18f57`; each exited 0, left its record `proposed`,
its worktree clean, and changed only its record and one new plan record:

| Item | Attempt | Cost | Duration | Turns | Tool calls | Plan |
| --- | --- | --- | --- | --- | --- | --- |
| G-260928-kehya | `G-260928-kehya.20260929T003720Z` | $1.73 | 4 min 41 s | 35 | 34 | G-260929-pjqxp, 142 lines, `a6b7cdf` |
| keyborg G-260928-ej7j0 | `G-260928-ej7j0.20260929T003721Z` | $2.98 | 9 min 23 s | 43 | 42 | G-260929-4fkdm, 212 lines, `b7ed1ae` |
| keyborg G-260928-g133n | `G-260928-g133n.20260929T003721Z` | $1.75 | 5 min 4 s | 32 | 31 | G-260929-ehvzq, 140 lines, `0bd0108` |

Plans total $6.46 of the $90 cap. No implementation attempt has run.

## Next

Waiting on [G-260929-2nfwv](G-260929-2nfwv-which-items-implementation-subagent.md),
checkpoint 2026-09-29: which item's implementation runs as a subagent, from
which session, and how the reviewer reaches Opus under a Sonnet
implementer. Branch `worktree-G-260928-pqhyg`, base `main` `9a18f57`.
Completed: the three plan attempts (Evidence); their worktrees
`worktree-G-260928-kehya` here and `worktree-G-260928-ej7j0`,
`worktree-G-260928-g133n` in keyborg hold the plans, statuses `proposed`.
No command owned. Once answered, reassign `/grove-work G-260928-pqhyg
--interaction headless`: add the reviewer routing the answer names to each
implementation worktree, run the two separate-process implementations as
`grove run ID --model claude-sonnet-5-5 --effort high --budget 15` on those
worktrees, the subagent one as answered, then add each item's facts to its
own record's Evidence and write the review record acceptance 2 asks for.
