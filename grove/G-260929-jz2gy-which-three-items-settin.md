---
id: "G-260929-jz2gy"
type: question
title: "Which three items, settings and cap for the plan-sufficiency experiment?"
status: open
created: "2026-09-29T00:23:42Z"
updated: "2026-09-29T00:24:06Z"
blocks: ["G-260928-pqhyg"]
relates_to: ["G-260924-5b6pz"]
---

## Question

[G-260928-pqhyg](G-260928-pqhyg-measure-whether-a-plan-f.md)'s acceptance 1
runs three items "with the settings and the cap the owner names at
assignment". Its headless assignment (`G-260928-pqhyg --interaction
headless`) named none, so nothing paid can start. The owner answers, with an
explicit value for each:

1. **Items.** Which three work records (here or in keyborg), and which one
   runs its implementation as a subagent in an interactive session.
2. **Plan attempt.** Model and effort for `grove run ID --until plan`.
3. **Implementation attempt.** Model and effort for the attempt that reads
   only the plan and the record.
4. **Review.** Model and effort (the record proposes `high`).
5. **Cap.** Dollars per attempt and in total for the experiment.

A blank or missing item keeps this question open: the recommendation below
is not a decision.

## Evidence

Observed at `main` `9a18f57`, 2026-09-29:

- Proposed work here with no undelivered prerequisite:
  [G-260928-kehya](G-260928-kehya-resume-an-attempt-s-sess.md) (small),
  [G-260928-dtrnw](G-260928-dtrnw-run-sweep-from-a-finishi.md) (medium, a
  worktree `worktree-G-260928-dtrnw` already exists at `9a18f57`, so it may
  be another session's), and
  [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) (large).
  G-260928-369c1, G-260928-c5j9d and G-260928-n4f1q each depend on an
  undelivered record. Keyborg's queue was not inspected.
- [G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) acceptance 5
  proposed Opus 5.5 at `xhigh` for the plan, Opus 5.5 at `medium` for the
  implementation, `grove-reviewer` at `high`, a $30 cap for the pair,
  against a sixteen-attempt mean of about $5 per attempt.
- `grove.yaml` `run:` defaults are `opus`, `high`, $50 per attempt.

## Recommendation

Not a decision. Items: G-260928-kehya here, plus two small or medium
keyborg items the owner picks, since this repository has only one small
unblocked item and y2p5h is large; kehya's implementation as the subagent,
being the smallest. Plan: `opus` at `xhigh`. Implementation: `sonnet` at
`medium`, a genuinely cheaper setting than `opus` at `medium`. Review:
`grove-reviewer` at `high`. Cap: $15 per plan attempt, $15 per
implementation attempt, $90 in total.

Only the owner can answer.

## Next

Open, blocking G-260928-pqhyg. When answered, record the answer here, set
`resolved`, and reassign `/grove-work G-260928-pqhyg`.
