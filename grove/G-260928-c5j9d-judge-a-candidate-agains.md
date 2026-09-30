---
id: "G-260928-c5j9d"
type: work
title: "Judge a candidate against its record under the policy with a delegated LLM"
status: proposed
created: "2026-09-28T19:28:59Z"
updated: "2026-09-30T20:22:05Z"
kind: feature
size: medium
depends_on: ["G-260928-dtrnw", "G-260928-n4f1q"]
relates_to: ["G-260928-d8py6", "G-260925-5wrn8", "G-260925-wh9ax", "G-260926-a8vyj", "G-260921-btyck", "G-260921-rz7bn", "G-260928-y2p5h", "G-260930-e8jj7", "G-260930-60c3d", "G-260930-tcc9w"]
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

## Portable milestone framing

On 2026-09-30 the owner explicitly brought this judge into
[G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md), superseding its
2026-09-29 deferral, in
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md). The workflow
must demonstrate bounded unattended judgment and progression as part of
the milestone. This remains a separately configured policy condition, not
the implementing or reviewing session accepting its own work.

Depends on [G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md) for
portable roles and independently attributed candidate evidence, and on
[G-260928-dtrnw](G-260928-dtrnw-run-sweep-from-a-finishi.md) for existing
policy triggers. It no longer depends on completing the entire milestone.
The [sequence work](G-260930-yfh91-continue-an-authorized-s.md) consumes its
approve/wait result; this judge never launches arbitrary subsequent work.

The technical proposal below is reconciled against the shared
[contracts](G-260930-84fnb-portable-workflow-contra.md). Recheck provider
capabilities at implementation and choose the judgment's provider and
effective limit explicitly. Do not silently inherit implementation settings.

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

- A `judge:` mapping under `policy.approve` with explicitly selected
  provider, model, effort and an effective per-call resource limit; exact
  spelling follows the portable role contract. Absent, the policy is as
  today. No silent fallback to implementation settings and no claim that a
  time limit is a monetary cap.
- One headless call with structured output that reads the record's outcome,
  constraints and acceptance, its plan, every review of the candidate, the
  record's Evidence and the diff against the target, and returns
  `{verdict: approve|wait, acceptance: [{item, met|not met|needs owner,
  evidence}], reasons}`. `wait` whenever an item needs a human, an item is
  not met, or the candidate exceeds the record's scope.
- Before approval, commit the judgment's input revisions, structured verdict,
  per-member acceptance reasons and evidence references as a validated
  supplemental artifact. The delegated verdict names it (`judgment J:
  approve, 4 of 4 acceptance items met`); a wait records its reasons where
  the board shows them. Raw provider logs can remain with the attempt.
  `sweep` lists `judge` as an act and `--dry-run` says it would run and
  its cap; usage also counts toward the selection's aggregate limits when it owns
  the operation. A new candidate, retry or delivery cannot reset that allowance.
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
3. Judgment input revisions, verdict, per-item reasons and evidence references
   remain in committed project evidence after automatic workspace cleanup and
   transfer. Raw provider logs may remain local. The judgment is inspectable
   and attributed to the actual candidate and policy revision; changed
   candidate or acceptance cannot inherit it.
4. Effective limits are enforced per call and across the owning selection;
   dry run states their kind, scope and known gaps. Usage that the provider
   cannot report is not zero and cannot silently reset the allowance.
5. A bounded real-provider trial compares the judge's verdicts with the
   owner's for an eligible candidate, a needs-human acceptance and an unmet
   criterion, under an assigned usage mandate. It includes both a single
   member and a combined candidate judged per member. These observations
   precede sequence work consuming the judge; they are not deferred until
   the final adoption trial.
6. Record model (`policy.approve.judge`), `docs/commands.md` (Sweep),
   `docs/board.md` and the Policy term say so.

## Next

Needs the declared independent-review and policy-trigger prerequisites.
Implement the configured judge under the reconciled contracts, then present
the bounded comparison with owner judgment before the sequence work consumes
it. It is now a milestone member; no implementation or paid call is assigned
by that inclusion.
