---
id: "G-260923-h9c30"
type: work
title: "Capture and reuse terms, questions and decisions across shaping, work and review"
status: accepted
created: "2026-09-23T19:44:56Z"
updated: "2026-09-30T13:18:01Z"
relates_to: ["G-260921-w9x25", "G-260921-e8bva", "G-260921-sth8q", "G-260923-fwakw", "G-260923-p5pt6", "G-260923-659zw", "G-260924-frzeg", "G-260924-wp2pe"]
candidate: "fdfd933e334eef1aa27631e44df7e535b762d244"
approved: "fdfd933e334eef1aa27631e44df7e535b762d244"
approved_by: owner
approved_context: "sha256:742babbccc099ab9715fa7bf2c2eff7d8d7f4cee48623aa33753f06b4d09194e"
---
## Outcome

Shaping, work and review share one knowledge-capture procedure, so that a
term settled in conversation, a consequential choice left open in execution,
and the answer that resolves it each reach the record that owns it and can be
retrieved by the next session instead of contradicted.

Owner intent, 2026-09-23: the owner is concerned that Grove does not reliably
save terms and questions during shaping, work or review, and favours the
pattern of Matt Pocock's grilling and domain-modeling skills, where
definitions are challenged against scenarios and code and recorded when they
settle, and consequential trade-offs become selective architecture decision
records. An external source-based assessment the owner read the same day
found the storage adequate and the discipline missing; this record follows
its diagnosis. Its formats and thresholds are proposed design.

Refined 2026-09-24 in an interactive shaping session after
[G-260923-p5pt6](G-260923-p5pt6-establish-behavioral-eva.md) closed: the sequencing this record waited on
is met, and the questions written since narrow the diagnosis (below).

## Constraints

Observed at main `f81f7e9`, first observed at `a54513f`; both guides are
byte-identical between the two and `grove version` prints guides digest
`3f5487904c61`, the digest G-260923-p5pt6's baseline ran at:

- The [shaping guide](../docs/work-shaping.md) says settled vocabulary belongs
  in a term record and when a question or decision may be written. The
  [work guide](../docs/work-execution.md) names decision records only under
  rejection and terms nowhere; its interactive missing-decision path asks in
  chat and relies on the checkpoint's pending judgments, while only the
  headless path persists a question. Review (step 6) examines acceptance and
  evidence, not knowledge.
- Every term (G-260921-vr8a8 to G-260921-vz0v3) was created in one batch under G-260921-w9x25 on
  2026-09-21 and none since; no work record after G-260923-twv25 links a term, other
  than this one. No decision has been recorded since G-260923-tnn5e.
- Questions are now created when a headless attempt meets the choice.
  [G-260923-hvnqh](G-260923-hvnqh-which-attempts-list-and.md),
  [G-260923-659zw](G-260923-659zw-what-mandate-should-the.md) and
  [G-260924-y99bx](G-260924-y99bx-how-do-the-eval-runs-log.md) were each written by a
  headless work attempt following the work guide's missing-decision path,
  with options, evidence, recommendation and `blocks`, and each was resolved
  by the owner writing the answer into the body. On the eval fixture,
  [G-260924-frzeg](G-260924-frzeg-baseline-runs-the-missin.md) reports 5 of 5 headless
  shaping runs persisting a blocking question and 5 of 5 companion runs
  asking nothing. The gap is after resolution and in interactive sessions,
  not in the question mechanics: none of the three resolved questions links
  a decision or term; G-260923-659zw's answer carries a product preference about the
  headless bound that no decision record holds and G-260924-frzeg reports as still
  open; the choice the G-260923-p5pt6 session took about synced account content in
  the clean configuration lives only in G-260923-p5pt6's Evidence; and the owner's
  interactive choices in G-260923-p5pt6's Constraints sit in the work record, which
  the shaping guide allows for routine ones and which nothing distinguishes
  from consequential ones.
- [G-260921-sth8q](G-260921-sth8q-attempt.md) still says Grove has no attempt record and that
  an attempt is visible only as a branch, a worktree and a checkpoint;
  [attempt.go](../internal/attempt/attempt.go) has written attempt files
  since G-260921-h46pb, and [docs/commands.md](../docs/commands.md#attempts) is where
  attempts are documented. The definition mixed an implementation observation
  into a settled meaning, and that part went stale unnoticed.
- `grove context` lists only what the selected work names in `depends_on`,
  `blocks`, `relates_to` or `members`, or links from its body
  ([context.go](../internal/handoff/context.go)). A term or decision nothing
  links is invisible to the next agent. G-260924-frzeg adds that the eval pair cannot
  see the listing at all, since the fixture's only knowledge is its brief.
- The [record model](../docs/record-model.md) already gives every type it
  needs: terms with meaning and relationships, questions with `blocks`,
  decisions with authority, alternatives and reconsideration conditions, and
  `relates_to` on any record. A resolved question keeps its answer in its
  body or links the durable decision.
  [G-260924-wp2pe](G-260924-wp2pe-answer-a-blocking-questi.md) proposes answering a
  question from the board under an `## Answer` heading; the resolver's step
  this record adds must fit that path, and G-260924-wp2pe must not presume this
  record's wording.
- The owner's open choice about the headless bound, whether a shippable
  proposal with a non-blocking question is preferred to blocking (G-260923-659zw
  item 5, G-260924-frzeg Disposition), edits the same "Missing human choice"
  paragraph this record's shaping change touches. This record does not make
  that choice and its text must fit either answer: it says when a question
  is persisted and what it carries and links, never whether the proposal
  blocks. If the owner makes the choice, it is a separate guide edit paired
  with a change to the eval check `question-blocks-proposal`.
- G-260923-p5pt6 is done at `60c9537`. Its case pair reruns for about $3 and ten
  minutes and is the regression check for a guide edit; its companion case
  is where a guide change that makes agents write more questions or terms
  would show as over-asking. The end-to-end knowledge sequence in G-260923-p5pt6's
  Next is the eval of this procedure and is follow-on work, not this
  record's.

In scope, all as guide text unless stated: before introducing or changing a
concept, read the terms and decisions that touch it and name conflicts and
synonyms; capture a definition in a term record when it settles, as meaning,
relationships, boundaries and misleading alternatives, with implementation
state and progress kept elsewhere; persist a consequential open question
before any wait or handoff in either interaction mode, with the choice,
evidence, recommendation, who can answer and what it blocks; on resolution,
keep the question, and when the answer is consequential by the threshold
below record it as a decision the question links, so an answer does not live
only in a resolved question's body; record a
decision only under a selective threshold (proposed: reversal cost, reasoning
a future reader would otherwise lack, and real alternatives), short, with
authority attributed as the shaping guide already requires; add to review
whether the candidate introduces a concept, contradicts a settled term,
depends on an unanswered choice, or implements a consequential decision
without its rationale, reported as findings for the author to reconcile; and
require work to link the terms and decisions that govern it so `context`
lists them. Repairing G-260921-sth8q is the first exercised instance of the term rule.

Out of scope: `context` supplying terms or decisions automatically, a
generated glossary or index, a new record type or field, any parallel
glossary file or separately numbered decision tree, the headless bound's
block-or-surface choice, and recording G-260923-659zw's preference as a decision
before the owner makes it. Reconsider the retrieval change only if evidence
shows links present and unread rather than absent.

## Acceptance

1. Both guides carry the procedure once, in the step where each activity
   meets it, and the review step names the knowledge check; the adapters are
   unchanged, `grove guide shape` and `grove guide work` print the change,
   and the Evidence records the new guides digest `grove version` prints, so
   a G-260923-p5pt6 rerun is comparable with its baseline at `3f5487904c61`.
2. G-260921-sth8q states the meaning without implementation state, and the observation
   it dropped lives where attempts are documented, with a link.
3. One real assignment after the change is walked through the procedure and
   its outcome reported honestly: which terms, questions and decisions it
   read, created or linked, and which rule it found unclear or unnecessary.
   Absence of a new record is a valid result when nothing settled.
4. `grove check` passes and every new link resolves. No product change beyond
   guide text and record bodies.

## Evidence

Branch `worktree-G-260923-h9c30`, based on main `25525cf` (record revision
`f4b81aa`). The candidate is the commit holding this text, named by the
status change that follows it; the reviewed content is `933094d`: guide edits `d2379f0`, G-260921-sth8q repair
`08a30b6`, review fixes `bf60f18`, `70ecd0c`, `54cb163`, `ebfcb7b`, and the
review record [G-260924-w0h5g](G-260924-w0h5g-independent-review-of-th.md) at
`933094d`. No plan: the In scope list names each rule and the change is
guide text, three documents and two record bodies. Two headless attempts
(`G-260923-h9c30.20260924T042824Z`, `G-260923-h9c30.20260924T043857Z`) wrote the commits
through `ebfcb7b` and G-260924-w0h5g, then each ended its turn on the eval pair
running as a background job, which died with the session; the interactive
session of 2026-09-24 finished from G-260924-w0h5g's state. That failure is proposed
on main as G-260924-3bapc and is not this record's.

Acceptance, checked at `933094d`, which the candidate differs from only by
this record:

1. `git diff main...HEAD` touches both guides, `docs/record-model.md`,
   `docs/commands.md`, G-260921-sth8q, this record and G-260924-w0h5g, and nothing under
   `.claude/skills/` or `.agents/skills/`. The work guide carries the rules
   in step 5 (read linked terms and decisions before a unit; read what
   touches a concept before changing it; capture what settles on the shaping
   guide's threshold), step 6 (the knowledge check, self-check included) and
   the interactive missing-decision bullet; the shaping guide carries them
   in step 2 (read the terms and decisions the topic touches), step 5 (term
   content, question persistence in either mode, answers kept and
   consequential ones recorded as decisions, the decision threshold,
   `relates_to`) and the headless bounds. `grove guide work` and `grove
   guide shape` print them. Guides digest at the candidate:
   `5a224350feae`, against the baseline's `3f5487904c61`.
2. [G-260921-sth8q](G-260921-sth8q-attempt.md) states the meaning only; the dropped
   observation is under [attempts](../docs/commands.md#attempts), and G-260921-sth8q
   links it.
3. Deferred by the owner's decision below to G-260924-wp2pe, which walks the
   resolution rule. This assignment's own walk, reported for what it is:
   before editing the guides it read [G-260921-e8bva](G-260921-e8bva-represent-terms-plans-an.md)
   and every settled term, linked G-260921-e8bva in `relates_to`, and created no term
   or decision, since nothing settled that the records lacked; G-260924-w0h5g ran the
   knowledge check and found no missing term and no contradiction. The
   review's finding 9 shows the one rule that needs a walk: an open part of
   a resolved answer (G-260923-659zw's preference) still has no persisted form until
   someone acting on G-260923-659zw or G-260924-frzeg writes the question.
4. `grove check`: `OK: 123 records` in the worktree; the new anchors
   `work-shaping.md#5-write-the-records` and `commands.md#attempts` resolve.
   `gofmt -l .` empty, `go vet ./...` clean, `go test -count=1 -timeout 120s
   ./...` all packages ok at `933094d` (several over five seconds while the
   eval pair ran beside them, not a change here).

G-260923-p5pt6 pair rerun, the owner's mandate, run 2026-09-24 04:50 UTC from this
worktree at `933094d`, guides digest `5a224350feae`, fixture `db18050`,
output `~/.cache/grove-evals/runs/2026-09-24-G-260923-h9c30-933094d` outside every
checkout, $2.88 for ten runs against the $50 cap:

- Missing-choice, against [G-260924-frzeg](G-260924-frzeg-baseline-runs-the-missin.md)
  finding 1: 5 of 5 runs persisted a question with `blocks` naming the
  proposal, and every check passed in every run, as in the baseline. What
  the guide edit added is visible: all five question bodies now name who
  can answer ("the owner"), which the shaping guide's step 5 already asked
  for and the work guide's headless list now asks for too. No run wrote a
  term or a decision. $0.28 to $0.30 and 44 to 49 seconds per run, against
  $0.27 to $0.30 and 45 to 62 seconds.
- Companion, over-asking: 5 of 5 runs proposed one work record and no
  question, term or decision, every check passed, $0.27 to $0.30 and 42 to
  51 seconds per run. The new step 2 reading (terms and decisions the topic
  touches) cost nothing visible: the fixture has none, and each run's one
  search was the same `grep` over the record bodies the baseline ran.
- Same limits as G-260924-frzeg: Claude Opus 5.5, five runs per case, one fixture
  with no term or decision to find, so the pair shows the edit did not
  make agents ask or write more, not that they read terms when they exist.
  The knowledge-sequence case in G-260923-p5pt6's Next is the eval of that. Two
  half-finished output directories from the abandoned headless attempts
  remain under the same `runs/` directory and are the owner's to remove.

## Next

In review; `candidate` names the commit. The integrator's next action, from a
clean checkout of `worktree-G-260923-h9c30`, then of `main`:

```sh
go run ./cmd/grove approve G-260923-h9c30 "VERDICT"
go run ./cmd/grove integrate G-260923-h9c30 --cleanup
```

Owner decisions, 2026-09-24, in the shaping session that refined this
record:

- The G-260923-p5pt6 pair reruns after the guide edit, as this record's own mandate,
  about $3 and a $50 cap, reported in Evidence with G-260923-p5pt6's Next untouched.
  Done above.
- G-260924-wp2pe walks acceptance 3, since it changes how a question is answered and
  so meets the resolution rule directly. Assign it after this record lands.

The headless bound's block-or-surface choice (G-260923-659zw, G-260924-frzeg) stays open and
is not this record's; G-260923-p5pt6's Next lists the knowledge-sequence eval case as
the follow-on once this lands.

Verdict on candidate fdfd933, 2026-09-24: needs more work, but later. Good enough for now

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
