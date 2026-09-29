---
id: "G-260929-2nfwv"
type: question
title: "Which item's implementation runs as a subagent, and from which session?"
status: resolved
created: "2026-09-29T00:37:24Z"
updated: "2026-09-29T00:53:23Z"
blocks: ["G-260928-pqhyg"]
relates_to: ["G-260929-jz2gy", "G-260924-5b6pz"]
---

## Question

[G-260928-pqhyg](G-260928-pqhyg-measure-whether-a-plan-f.md)'s method runs
one of its three items' implementation "as a subagent told to read only the
plan". The answer to
[G-260929-jz2gy](G-260929-jz2gy-which-three-items-settin.md) names the items
(G-260928-kehya here, keyborg G-260928-ej7j0 and G-260928-g133n), the plan
setting (Opus 5.5 at `xhigh`), the implementation setting (Sonnet 5.5 at
`high`), the reviewer (Opus 5.5 at `high`) and the budgets ($15 per plan
attempt, $15 per implementation attempt, $90 in total), but not which item
is the subagent. The work guide (step 5) treats an omitted item as a missing
decision, never the recommendation, so no implementation can start until the
route of each is known. The owner answers, with an explicit value for each:

1. **Subagent item.** Which of the three runs its implementation as a
   subagent; the other two run as `grove run ID --model claude-sonnet-5-5
   --effort high --budget 15` on the worktree the plan attempt left.
2. **Parent session.** The record says the subagent run needs an
   interactive session, since the work guide does not tell an attempt to
   delegate. Either (a) an interactive session the owner starts, or (b) a
   headless `/grove-work G-260928-pqhyg` session acting as the parent: it
   dispatches a `general-purpose` subagent on `sonnet` told to read only the
   plan and the record, then reviews and hands off itself. Under (b) the
   subagent's effort cannot be set (the Agent tool takes a model, not an
   effort), so it runs at the provider's default, not `high`.
3. **Reviewer routing.** `grove-reviewer` is `model: inherit`, so under a
   Sonnet implementation it runs on Sonnet. The attempt process strips
   `CLAUDE*` variables, and the definition is managed by `init` (changing it
   is a product change acceptance 3 excludes). The one route that changes
   nothing tracked is an ignored `.claude/settings.local.json` in each
   implementation worktree setting `CLAUDE_CODE_SUBAGENT_MODEL=claude-opus-5-5`
   and `CLAUDE_CODE_SUBAGENT_MODEL_FORCE=1` (Claude Code documents FORCE as
   overriding the definition's `model` and any per-call model). It routes
   every subagent the implementer starts to Opus, not only the reviewer.
   Accept that, or name another route.

## Evidence

Observed 2026-09-29: Claude Code 2.1.284; `internal/attempt/attempt.go`
builds the child environment without `CLAUDE*` except `CLAUDE_CONFIG_DIR`;
`.claude/settings.local.json` is in the owner's global Git ignore; the
subagents page of the Claude Code documentation gives the order per-call
`model`, definition `model` (`inherit` meaning the parent's),
`CLAUDE_CODE_SUBAGENT_MODEL`, parent model, with `_FORCE=1` putting the
variable first.

## Recommendation

Not a decision. 1: G-260928-kehya, the smallest and in this repository. 2:
(b), so the experiment needs no interactive session, and record the effort
gap as a limit. 3: accept, and report the Opus spend by model from each
attempt's `Cost by model:`, which separates the reviewer's (and any other
subagent's) spend from the implementer's.

## Next

Open, blocking G-260928-pqhyg's implementation attempts; its plan attempts
ran. When answered, record the answer here, set `resolved`, and reassign
`/grove-work G-260928-pqhyg --interaction headless`.

## Answer
1. kehya
2. b
3. accept for now. I want to consider the value of the reviewer agent definition instead of just injecting the prompt/context like we do with shaping or implementation after this is done
