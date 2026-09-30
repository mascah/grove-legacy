---
id: "G-260929-gm3m4"
type: work
title: "Deliver accepted work locally with clean history and retained evidence"
status: accepted
created: "2026-09-29T16:15:25Z"
updated: "2026-09-30T20:22:12Z"
kind: feature
size: large
depends_on: ["G-260930-62nmj"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-4742q", "G-260930-r2k4g", "G-260921-3qgsf", "G-260921-jatts", "G-260921-btyck", "G-260930-84fnb", "G-260930-2qa4a", "G-260930-gj9d7", "G-260930-tcc9w"]
candidate: "b59f3de066bb41c5e97367c0f9b83b5e2a7e9f1d"
approved: "b59f3de066bb41c5e97367c0f9b83b5e2a7e9f1d"
approved_by: owner
approved_context: "sha256:f0006ce02ab18dbf2966ad9891bad89bef1b457bb038dd234735fc0b8039abab"
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
inputs and resulting content, retain the evidence objects, prove that
correspondence when a delivery is made and on audit, and make completion,
dependency delivery and cleanup read the same fact: the target's own copy
of the accepted record.

A work ID, a matching patch or unverified trailers alone are insufficient
authority to close work. Changed candidates or targets require appropriate
reconsideration/reverification. Partial failure after integration must be
recoverable without a second delivery or destruction of evidence.

The owner selected [G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md):
record acceptance and derive Done from its delivery, using a redesigned
record contract. They amended it with
[G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md): a delivery is
proved once, when it is made, and again only by an audit a person runs;
reading costs what it did before, whatever the history, branches or
deliveries, and leaves nothing in a repository beyond one squash commit and
one local ref per delivery. This item owns that schema/migration change and the shared
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
retain candidate and evidence before target advance. Reading derives
completion from the target's record after a crash or an external merge
without a second delivery. Dependencies also check that the execution
base's copy of the record holds the target's acceptance. An unreadable
target waits without spending an attempt.

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
   stale evidence; forged or mismatched delivery claims are rejected where
   a delivery is proved: at delivery and by audit.
3. Recorded acceptance is truthful before and after delivery; the redesigned
   work contract no longer declares permanent review or requires a stored
   done transition. CLI list/show/context and JSON, board, dependency/base
   readiness, attempt guards, policy integration and cleanup consume the same
   standing, read from the target's own copy of the record; a delivery is
   proved at delivery and by audit, never on a read. Raw source remains
   distinguishable from observed facts. An unreadable target or a stale
   acceptance cannot authorize duplicate work or Done; an acceptance
   written on the target by hand reads as done until an audit reports it.
4. Refusal and recovery tests cover dirty targets, conflicts, shared groups,
   target movement, changed acceptance/candidate, external delivery, failure
   just after target advance and safe retry. Done is reconstructed without a
   completion commit. Reopening invalidates old completion for current work.
5. Candidate, review, approval and correspondence evidence survive cleanup
   and ordinary garbage collection. An ordinary fresh-clone transport path
   is documented and exercised; missing evidence and shallow history make
   the audit report a delivery it cannot prove here, never success, and
   reading needs no evidence. Evidence refs do not masquerade as current
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

## Evidence

Implemented on branch `worktree-G-260929-gm3m4` from base `556f362` (main),
following [the plan](G-260930-3hcv4-plan-for-g-260929-gm3m4.md). Commits:
`ce12b1a` schema 4 contract, `6159f4b` `grove migrate`, `04edb08` this
repository's records migrated, `4c2ee0c` standing, `7466c1e` squash
integrate, `a08bdae` consumers, `8548c07`/`48bb0c1`/`af6c8df` documents and
terms, `e17ddd4`, `0e5e8cd`, `826b97f`, `da5e5ed`, `9c9fed2` review fixes,
`6b7ab0b` cleanup retains the submitted tip, then the owner's redesign
([G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md)): `c71b0c2`
code, `bfe2c80` and the commit holding this Evidence, documents and records.

1. One squash commit per delivery, `TYPE: title` from the candidate's own
   commits, members listed, trailers `Grove-Work`, `Grove-Candidate`,
   `Grove-Submitted`; a group is one commit; no record commit
   (`TestIntegrateSquashesAndRetainsEvidence`,
   `TestGroupIntegratesOnlyWhenEveryMemberIsApproved`). /tmp trial: main
   read `feat: add greeting` over `chore: unrelated work on main`, with a
   `fixup:` and a `test:` commit folded in. A branch kept after its delivery
   merges the target first, and its next delivery's type comes from all its
   commits since the fork (`TestIntegrateDeliversReopenedWorkAgain`). The
   owner's judgment of representative history is still owed.
2. Acceptance binds candidate and context (`approved_by`,
   `approved_context`); a moved target and a changed candidate or context
   are refused. Forged, altered or code-carrying delivery claims are not
   proved by integrate's proof of its own commit or by
   `grove check --deliveries` (`TestAuditRejectsForgedDeliveries`,
   `TestAuditSurvivesUnrelatedClaims`,
   `TestIntegrateLeavesAnAlteredDeliveryToTheAudit`).
3. Status `accepted`; done read from the target tip's own copy of the
   record, by list (STANDING), show (stderr, JSON `standing`), deps,
   context, run selection, resolve, feedback, sweep and the board
   (`TestStandingReadsTheTargetsRecord`, `TestJudgeStartsNothing`). A
   reading starts one Git process, a `cat-file --batch`, and no log,
   rev-list, merge-tree, merge-base or for-each-ref, whatever the
   deliveries and versions (`TestStandingGitWorkIsConstant`, GIT_TRACE).
   Traced once each on this repository: `grove list` starts 0 Git
   processes on main and 1 here; `grove versions`, the board's inspection,
   8 on both, and the board adds none. One timed `list` each: about 20 ms
   CPU on main, 40 ms here, at the timer's resolution. An acceptance
   written on the target by hand reads done and only the audit reports it
   (`TestAuditProvesDeliveries`), the trade the owner accepted.
4. Dirty target, conflicts, shared groups, target movement, changed
   acceptance, crash after advance and retry, reopen then redeliver after
   merging the target, a candidate merged alone then delivered, and
   feedback after delivery are tested in `internal/integrate` and
   `internal/standing`; done is reconstructed from the target's record.
5. `refs/grove/submitted/S` survives cleanup and `gc --prune=now`; cleanup
   writes it before deleting a branch at S
   (`TestIntegrateCleanupRetainsTheSubmittedTip`). A `--no-local` clone
   reads done without it, and its audit says it cannot prove the delivery
   there until `git fetch origin 'refs/grove/*:refs/grove/*'`; a shallow
   clone cannot audit (`TestAuditTransport`).
6. `grove migrate` dry run and `--commit` with `refs/grove/schema-3/BRANCH`;
   on a clone of main: 50 accepted, 29 kept as schema 3's claim, 190
   unchanged, check OK. A schema 3 branch is refused with the path. Model,
   commands, board, work guide, record design, README, CLAUDE.md, the
   brief, `just clean-merged`, the contract, and the Work, Approval,
   Integration, Candidate and Standing terms reconciled with the redesign.
   The entrypoint revision is unchanged: the entrypoints need nothing new
   of the binary.
7. Fixtures cover both entry paths; the real-agent evaluation is the
   owner's.

Verification after the redesign: `go vet ./...` and `gofmt -l .` clean;
`go run ./cmd/grove check` OK: 273 records; one
`go test -count=1 -timeout 120s ./...` run failed two size budgets (the
record model at 12594 bytes, the work guide head at 13030), which the last
commit trims; the two packages then passed on rerun (`.` and
`internal/cli`), every other package passed in the full run;
`python3 internal/tui/testdata/terminal.py` all ok.

Review: [G-260930-c0g5p](G-260930-c0g5p-review-of-g-260929-gm3m4.md): three
rounds, the owner's review, a fix round and a re-review at `9c9fed2`, Open
findings: 3, of which the redesign removes the cause of two and `6b7ab0b`
fixes the third. The redesign's review: `2cb7bdd`, 5 findings, fixed in
`f1fa1b7` and `aaf2d5a`; `aaf2d5a`, 2 findings on carried work, fixed in
`a4453d1`; `281a8f1`, 2 findings on carried work delivered for another
candidate, fixed in `e144aa6` by one rule: the target's own copy of carried
work decides, and the latest of several acceptances is delivered
(`TestIntegrateDeliversCarriedWorkFirst`,
`TestIntegrateCarriedWorkDeliveredElsewhere`,
`TestIntegrateCarriedWorkRevisedElsewhere`,
`TestIntegrateCarriedWorkOnOneBranch`,
`TestIntegrateDeliversTheLatestAcceptance`). At `e144aa6`: `go vet ./...`,
`gofmt -l .` clean, `grove check` OK, and one
`go test -count=1 -timeout 120s ./...` passed every package.

## Next

Accepted by the owner at b59f3de and present on main. Its implementation and
review evidence above remain historical facts, including the last fix's
independent-review limit. Do not rerun or silently expand this completed item.

On 2026-09-30 the owner paused adopter migrations and selected narrower
delivery boundaries in
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md).
[G-260930-0s29t](G-260930-0s29t-use-one-fresh-workspace.md) owns workspace
retirement/cleanup and removal of repeated squashed-branch delivery;
[G-260930-yfh91](G-260930-yfh91-continue-an-authorized-s.md) owns the two
selection modes and automatic progression. Cheap standing and retained
optional audit evidence remain selected. The
[milestone](G-260930-60c3d-complete-the-portable-gr.md) owns the migration
hold and remaining journey evidence. The earlier suggestion to drop all
retained refs is not selected: the owner chose evidence to survive cleanup.

Verdict on candidate b59f3de, 2026-09-30: approved
