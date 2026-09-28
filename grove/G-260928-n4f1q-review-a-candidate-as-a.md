---
id: "G-260928-n4f1q"
type: work
title: "Review a candidate as a separately launched independent process"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-09-28T19:34:02Z"
kind: feature
size: medium
depends_on: ["G-260928-y2p5h"]
relates_to: ["G-260928-917h8", "G-260924-5b6pz", "G-260921-rz7bn", "G-260925-5wrn8", "G-260924-3bapc", "G-260928-c5j9d"]
---

## Outcome

## Constraints

## Acceptance

## Next

## Outcome

An independent review of a candidate can be launched as its own Grove-owned
process, on a provider and model chosen for review, from the CLI, from the
board, or by the implementing session in place of a subagent. It writes the
review record with `examined`, and the attempt facts show which reviews ran
that way and what they cost, so review runs at the strength and on the
harness the owner picks, and an implementing session on any provider can
hand off reviewed work.

Owner intent, shaping conversation 2026-09-28: "codex for the plan, claude
code for the implementation, and codex for the review"; "an orchestrator
model, like fable, and then execution models like opus 5.5, and review
models"; the owner chose the incremental path, where the work guide stays
the one workflow and phases gain launchable boundaries.

## Constraints

Observed 2026-09-28 at main `6fbb888`:

- Step 6 of the work guide dispatches `grove-reviewer` as a fresh subagent
  per gate, Claude only (`.claude/agents/grove-reviewer.md`: `model:
  inherit`, `effort: high`, read-only tools;
  [G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md)). The review
  guide (`grove guide review`) is the brief; the reviewer returns findings
  and a closing line, and the implementing session writes the review
  record. Codex has no equivalent (`docs/commands.md`, Init).
- The implementing session's context peaks after the reviewer, 278k to
  352k tokens at handoff, and the fix rounds and record writing run there
  (G-260924-5b6pz's observations).
- The policy reads `Open findings: none` from a `current` review record
  that examined the candidate (record model).
- `attempt.json` records the reviewer definition's digest; a subagent
  review shows in `Metrics` as a subagent and in the result's per-model
  cost.
- A headless session must own a long command in the foreground with a
  timeout ([G-260924-3bapc](G-260924-3bapc-headless-attempts-run-lo.md)).

Proposed design, labelled proposed:

- `grove review ID [--commit C] [--provider P] [--model M] [--effort E]
  [--budget USD | --max-minutes N]`: one Grove-owned process in the branch's
  checkout, read-only for the tree (a Claude permission mode that disallows
  edits, or `codex -s read-only`), whose prompt is the review guide with
  what step 6 passes: checkout, commit or range, record, plan, allowed
  commands, and prior findings for a re-review. It returns findings and the
  closing line, as structured output where the provider supports it, and
  Grove writes the review record with `work` and `examined` and commits it,
  attributed to the launch. Refused when the branch's checkout is dirty, or
  when an attempt of the work runs, unless that attempt launched it.
- Step 6 gains the alternative: where the checkout's provider has no
  reviewer definition, or the assignment names a review provider, dispatch
  through `grove review` and wait for it in the foreground; the closing line
  comes from the record it wrote.
- `run: review:` defaults (provider, model, effort, cap) beside the launch
  defaults; the attempt records the reviews it launched.

Out of scope: replacing the subagent path where it works; a second reviewer
definition; the delegated judgment of
[G-260928-c5j9d](G-260928-c5j9d-judge-a-candidate-agains.md), which is not a
review.

## Acceptance

1. `grove review ID` on a candidate writes a `current` review record with
   `examined` the commit, ending with the closing line, committed on the
   branch; covered with a fake provider on both providers; refusals as
   listed.
2. A headless implementation attempt on Codex hands off reviewed work
   through it: the review record exists and the policy would read it, fake
   providers end to end.
3. The attempt facts and the board's attempt screen list the reviews an
   attempt launched, with provider, model and cost or cap.
4. Work guide step 6, the review guide, `docs/commands.md` and the record
   model (`run.review`) say so; the [review](G-260921-rz7bn-review.md) term's
   meaning is unchanged.
5. One real review of a real candidate on Codex and one on Claude at a
   chosen model, which the owner compares with a subagent review of the
   same commit, under a budget the owner names at assignment.

## Next

Assign after [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) is
delivered: `/grove-work G-260928-n4f1q`. `depends_on` names it because the
review process launches through the provider seam that record introduces
(command composition, events reader, cap, stop); built first on Claude
alone it would be redone on the seam.
