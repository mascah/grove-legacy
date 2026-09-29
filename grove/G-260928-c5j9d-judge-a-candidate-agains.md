---
id: "G-260928-c5j9d"
type: work
title: "Judge a candidate against its record under the policy with a delegated LLM"
status: proposed
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T19:34:01Z"
kind: feature
size: medium
depends_on: ["G-260928-dtrnw"]
relates_to: ["G-260928-d8py6", "G-260925-5wrn8", "G-260925-wh9ax", "G-260926-a8vyj", "G-260921-btyck", "G-260921-rz7bn", "G-260928-y2p5h"]
---

## Outcome

Under the standing policy, a candidate that passed verification and whose
independent review closed clean is approved and integrated when a delegated,
bounded judgment finds that it did what its record asked and needs no
owner's eye; otherwise it waits with the judgment's reasons on the record,
so the owner reads reasons instead of diffs.

Decision [G-260928-d8py6](G-260928-d8py6-a-standing-policy-may-de.md), the
owner on 2026-09-28: "An LLM decides ideally." The design below is proposed
and the owner writes the policy that enables it.

## Constraints

Observed 2026-09-28 at main `6fbb888`:

- The policy's approval conditions (record model, `policy:`): a `current`
  review examined the candidate and closed `Open findings: none`; no open
  question; nothing after the candidate but its record; no `never` path, no
  binary, within `max_lines`; a clean merge; `verify` passing on the merged
  result. None reads the record's acceptance.
  [G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md) named `paths`,
  `kinds` and `sizes` as likely extensions.
- Both providers return structured output headless: Claude Code 2.1.283
  `--json-schema`, Codex 0.157.0 `--output-schema` (`--help`, 2026-09-28).
  The eval runner (`evals/run.py`) already composes bounded headless calls
  with a clean configuration directory.
- Acceptance items that name the owner's judgment are common: keyborg's
  G-260928-vb6rt ("Owner judgment: not yet given"),
  [G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) here ("The
  owner judges the revised screens in a terminal"). No deterministic check
  can pass them.
- A delegated verdict begins `delegated under policy grove.yaml REVISION`
  and names its evidence (record model;
  [G-260926-a8vyj](G-260926-a8vyj-policy.md)).

Proposed design, labelled proposed:

- A `judge:` mapping under `policy.approve` with `model`, `effort` and
  `budget` per judgment, all required to enable it; absent, the policy is
  as today. Its provider is the target's `run:` provider once
  [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) lands, else
  Claude Code.
- One headless call with structured output that reads the record's outcome,
  constraints and acceptance, its plan, every review of the candidate, the
  record's Evidence and the diff against the target, and returns
  `{verdict: approve|wait, acceptance: [{item, met|not met|needs owner,
  evidence}], reasons}`. `wait` whenever an item needs a human, an item is
  not met, or the candidate exceeds the record's scope.
- Its input and output are kept with the attempt's files and quoted in the
  delegated verdict (`judgment J: approve, 4 of 4 acceptance items met`); a
  wait appends the reasons to the record as a paragraph the board shows.
  `sweep` lists `judge` as an act and `--dry-run` says it would run and
  its cap; its spend counts in the sweep's aggregate budget.
- Kept deterministic and unchanged: `never`, `verify`, the review's closing
  line, the clean merge.

Out of scope: a `feedback` verdict that reopens the work and spends an
attempt, which the owner may add once the wait reasons prove reliable; the
judgment as a reviewer, since the review stays the reviewer's; learning
from the owner's verdicts.

## Acceptance

1. With `judge:` absent, `sweep` is unchanged. With it and a fake judge:
   approve leads to the delegated approval and integration whose verdict
   names the judgment and its result; wait leaves the record in review with
   the reasons appended and visible on the board.
2. A judgment that fails, times out, exceeds its budget or returns malformed
   output is a wait with that reason; it never approves.
3. The judgment's input and output are retained, printable by `grove attempt`
   or a listing, and attributed to the policy revision.
4. Spend is bounded per call and counted in the aggregate budget; dry run
   states it.
5. A real-provider trial on this repository's next few candidates compares
   the judgment's verdicts with the owner's, reported in a review record,
   under a budget the owner names at assignment.
6. Record model (`policy.approve.judge`), `docs/commands.md` (Sweep),
   `docs/board.md` and the Policy term say so.

## Next

Assign after [G-260928-dtrnw](G-260928-dtrnw-run-sweep-from-a-finishi.md)
is delivered: `/grove-work G-260928-c5j9d`. `depends_on` names it because
both change the sweep's per-candidate plan and the policy's keys in
`internal/sweep` and `internal/project`; the trigger is the smaller change
and this one extends the same loop, so building them alongside would
conflict.
