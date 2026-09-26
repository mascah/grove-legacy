---
id: "G-260925-beby3"
type: work
title: "Keep a blank mandate answer from approving spend"
status: active
created: "2026-09-25T03:47:15Z"
updated: "2026-09-26T20:46:28Z"
relates_to: ["G-260924-59f5k", "G-260924-7x7p7", "G-260925-04ccr", "G-260925-gymkr", "G-260925-02jsj", "G-260924-ecs9m"]
---

## Outcome

A work session never spends on an answer the owner did not give: a question
about what a paid step spends (model, effort, runs, budget or cap) takes
only explicit values, and a blank or missing item keeps it open and
blocking. Owner intent, recorded as
[G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md) on 2026-09-24.

## Constraints

Observed at `worktree-G-260924-59f5k` `48358b7`:

- [G-260924-7x7p7](G-260924-7x7p7-what-mandate-and-login-s.md), written by a
  headless `/grove-work G-260924-59f5k --until plan` session, said "any item left
  blank takes the recommendation"; the owner's answer covered the login
  only, and the continuation ran 9½ paid Codex runs on the recommended
  `gpt-6-astra` at `high`, 12% to 95% of the ChatGPT Plus five-hour window.
- G-260924-59f5k excluded changing the guides, so this is separate work, and
  its handoff records that the owner rejected putting the rule in
  `AGENTS.md`.

Observed at main `2bd50b9`, shaping 2026-09-26:

- [docs/work-shaping.md](../docs/work-shaping.md) says "Your recommendation
  is not a decision. Silence is not agreement."
  [docs/work-execution.md](../docs/work-execution.md) has no equivalent: step
  5 says a resolved question's answer is the constraint, and "When a human
  decision is missing" forbids a headless session picking a default for a
  product choice, but nothing says a blank item in an answer is not an
  answer, at the point where a session writes a question or reads a resolved
  one.
- The rule held once without the guide:
  [G-260925-gymkr](G-260925-gymkr-what-mandate-should-the.md) (headless,
  2026-09-25) refused to let a blank take the recommendation because
  G-260925-pbx81's body named G-260925-04ccr. It failed once without it:
  [G-260925-02jsj](G-260925-02jsj-how-should-an-adopting-p.md), written two
  hours after G-260925-04ccr for a non-spending choice, repeated "any item
  left blank takes the recommendation". Nothing was spent and the owner
  answered every item, but the phrase recurs unless the guide forbids it,
  and an adopting project never sees G-260925-04ccr.
- Software already bounds the rest: `evals/run.py` requires an explicit
  `--model` and `--max-plan-percent` (G-260924-59f5k), and `grove run`
  spends only the values the owner set in `grove.yaml` under `run:` or
  passed on the line (G-260924-ecs9m). This work is only about questions a
  session writes and reads.

Proposed design: two sentences in `docs/work-execution.md`, written for
every consequential item rather than spend alone, since a blank on a product
choice is the default the headless list already forbids: where a session
writes a question, it never offers that a blank item takes the
recommendation; where step 5 reads a resolved question, an item the step
depends on that has no answer is a missing human decision for that item,
and what a paid step spends is the named example. No entrypoint revision
changes, since the entrypoints need nothing new of the binary.

Out of scope: `AGENTS.md`, the shaping guide, the record model, and any
software check on a question's body.

## Acceptance

1. `docs/work-execution.md` says, in "When a human decision is missing"
   where a question is written and in step 5 where a resolved question is
   read, that a blank item is not an answer, naming what a paid step spends
   (model, effort, runs, budget, cap) among the items that must be
   explicit, consistent with the shaping guide's "Silence is not agreement".
2. `go run ./cmd/grove guide work` prints the change and
   `go run ./cmd/grove check` passes.

## Next

Captured 2026-09-24 from the G-260924-59f5k incident. Shaped 2026-09-26:
recommendation is to assign it as written, a documentation-only change
with link and consistency checks; the alternative is to abandon it and
rely on every paid work record naming G-260925-04ccr, as G-260925-pbx81 did,
which leaves unlinked sessions and adopting projects uncovered. Owner
decides: `/grove-work G-260925-beby3`, or `update` to `abandoned`.
