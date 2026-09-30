---
id: "G-260929-gm3m4"
type: work
title: "Deliver accepted work locally with clean history and retained evidence"
status: proposed
created: "2026-09-29T16:15:25Z"
updated: "2026-09-30T03:18:02Z"
kind: feature
size: large
depends_on: ["G-260930-62nmj"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-4742q", "G-260930-r2k4g", "G-260921-3qgsf", "G-260921-jatts", "G-260921-btyck", "G-260930-84fnb", "G-260930-2qa4a"]
---

## Outcome

An accepted change can be delivered locally with useful Conventional Commits
history and durable evidence connecting it to the reviewed candidate.
A project can use its preferred release tooling without exposing every
agent checkpoint and review repair as a separate product change.

Owner intent, 2026-09-29: "fix the commit problem more than anything so that
things are not noisy no matter what tool people use." The portable adoption
direction [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md) retains this record as the local delivery
component of [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).

## Scope and constraints

Keep a complete local path without a remote, credentials, hosted CI or a
resident service. Squash delivery is in scope. A single work item normally
produces one meaningful product-change commit; a shared candidate is
delivered as its approved group without duplicating the code change.

Preserve exact review/approval attribution and verification of the result
against the actual target. A squash changes commit identity. The design
must establish correspondence between approved candidate, integration
inputs and resulting content, retain the evidence objects, and make
completion, dependency delivery and cleanup use that same fact.

A work ID, a matching patch or unverified trailers alone are insufficient
authority to close work. Changed candidates or targets require appropriate
reconsideration/reverification. Partial failure after integration must be
recoverable without a second delivery or destruction of evidence.

The owner selected [G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md):
record acceptance and derive Done from verified delivery, using a redesigned
record contract. This item owns that schema/migration change and the shared
standing and delivery operations. It must remove misleading permanent review
state, not add a board-only override. No mandatory post-delivery done commit
or follow-up completion PR is part of the selected path.

Update every existing lifecycle consumer needed to preserve coherent local
operation: parser/validation, create/update/approve/feedback, list/show/context,
deps and current views, attempt selection/reporting/resolution, integration,
policy sweep, cleanup, and the CLI/TUI adapters and guides that expose them.
The later experience item owns redesigned onboarding and presentation, not
this correctness work. Preserve current Claude execution while its caller
contract changes; provider extraction follows through the declared edge.

Carry acceptance and discoverable delivery inputs in the prepared submission;
retain candidate and evidence before target advance. Inspection must derive
completion from those facts after a crash or an external merge without a
second delivery. Dependencies also check the execution base contains the
delivered result. Unknown evidence waits without spending an attempt.

Message formatting remains an implementation design detail within useful
Conventional Commit history; deriving a type from the current kind field is
not assumed correct. Do not replace release tools or select version numbers.

Hosted PR delivery belongs to [G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md); this record owns the reusable
local delivery/correspondence behavior that it will consume. No host API,
push or protected-branch bypass is part of local integration.

## Observed evidence

The original 2026-09-29 shaping found 145 of 189 direct main commits since
2026-09-24 touched records only, and 7 of 14 fix commits in the keyborg pilot
were review-round fixes. Those are dated observations, not new measurements.

At main c8fa070ef9ff, internal/update/update.go checks candidate ancestry for
done, internal/deps checks delivery through ancestry, internal/integrate
performs a plain merge followed by separate done commits, and just
clean-merged selects branches by Git's merged relation. The earlier proposal
suggested Grove-Work/Grove-Candidate trailers and retained refs. That is a
candidate design, not sufficient proof of delivery and not a claim that
these are the only affected callers.

The settled [Integration](G-260921-3qgsf-integration.md) concept and the
implemented lifecycle must be reconciled when this work lands. The current
mechanism remains in force until then.

Other threads retained from the original conversation: release membership
was owner-selected, while a release/roadmap feature and release-please setup
remain separate future shaping; Grove does not own version numbers.
A catch-up digest is now in [G-260930-r2k4g](G-260930-r2k4g-make-the-complete-grove.md); hosted delivery is in
[G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md). The owner's correction that headless subagents work (the
observed denial was a sibling-worktree write) remains relevant to the
separate-process versus subagent comparison in the design.

## Dependencies

Depends on [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md): candidate/delivery correspondence, metadata
placement and the visible approval/delivery journey determine which
mechanisms and lifecycle contracts change.

## Acceptance

1. Local integration of one approved change and of a shared approved group
   produces the designed clean product history and traceable delivery
   evidence. Record-only bookkeeping does not reintroduce the noise this
   outcome removes; the owner judges representative main/changelog output.
2. Review, approval and verification are bound to the candidate and actual
   integration inputs. A moved target or altered candidate cannot use
   stale evidence; forged or mismatched delivery claims are rejected.
3. Recorded acceptance is truthful before and after delivery; the redesigned
   work contract no longer declares permanent review or requires a stored
   done transition. CLI list/show/context and JSON, board, dependency/base
   readiness, attempt guards, policy integration and cleanup consume the same
   verified standing. Raw source remains distinguishable from observed facts.
   Unknown or stale evidence cannot authorize duplicate work or false Done.
4. Refusal and recovery tests cover dirty targets, conflicts, shared groups,
   target movement, changed acceptance/candidate, external delivery, failure
   just after target advance and safe retry. Done is reconstructed without a
   completion commit. Reopening invalidates old completion for current work.
5. Candidate, review, approval and correspondence evidence survive cleanup
   and ordinary garbage collection. An ordinary fresh-clone transport path
   is documented and exercised; missing evidence and shallow history report
   unknown rather than success. Evidence refs do not masquerade as current
   work branches or conceal real post-delivery edits. Neither a ref name nor
   a trailer is proof.
6. A dry-run and explicit versioned migration preserve IDs, paths, original
   provenance and recoverable Git state. Legacy done claims retain their
   historical meaning and dependency semantics without invented evidence;
   ambiguous or mixed-schema branches have a clear reconciliation path.
   The current schema is not silently reinterpreted. Delivered behavior is
   reconciled in the model, repository policy, guides, entrypoint revisions
   where required, Work/Approval/Integration terms and cleanup commands.
7. A bounded local trial demonstrates useful history, inspectable evidence,
   existing Claude workflow continuity and recovery without hosted services.
   Fixtures exercise direct-file and tool-mediated entry paths; the complete
   local proof additionally evaluates real agents against both paths.

## Next

Ready for assignment: $grove-work G-260929-gm3m4. The owner accepted
[G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md)'s design at f376a6d,
already on main. The completion-model question is settled in
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md); do not reopen it
as an implementation preference. The
[contracts design](G-260930-84fnb-portable-workflow-contra.md) owns its
semantics and migration requirements.

At assignment, write the implementation plan around these concrete units:

- Inventory raw work-status consumers and define the versioned persisted
  acceptance and shared standing types. Include CLI/JSON compatibility and
  the migration dry-run before changing live project records.
- Build deterministic acceptance/delivery fixtures and the shared verifier,
  including current-record selection, grouped candidates, context changes,
  source freshness, target/base reachability and conservative unknowns.
- Connect existing operations, Claude runner guards, context and views to
  that result; retain ownership, environment and exact-source safety tests.
- Implement prepared local squash delivery and evidence transport/retention,
  with atomic target checks, full-result verification and crash recovery.
- Exercise the local journey and migration, reconcile shipped documentation
  and repository policy, and present representative history for judgment.

The plan must specify descriptor serialization without a hash cycle, retained
ref transport, schema revision and command migration. Those are technical
implementation choices bounded by the selected contract, not reasons to
create another generic architecture project. No implementation is assigned
by these preparation notes.
