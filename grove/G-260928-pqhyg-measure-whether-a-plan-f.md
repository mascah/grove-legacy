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

## Next

Waiting on [G-260929-jz2gy](G-260929-jz2gy-which-three-items-settin.md),
checkpoint 2026-09-29. A headless assignment `G-260928-pqhyg --interaction
headless` named no items, settings or cap, so nothing ran and nothing was
spent; status stays `proposed`. Branch `worktree-G-260928-pqhyg`, base
`main` `9a18f57`. Completed: context and state inspected, question written.
No commands owned, no plan, no attempt launched. Once the question is
resolved, reassign `/grove-work G-260928-pqhyg` naming its answer; the
bounded plan attempt and the implementation attempt already exist as
`grove run --until plan` and `R`.
