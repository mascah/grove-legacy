---
id: "G-260930-0s29t"
type: work
title: "Use one fresh workspace per delivery with safe automatic cleanup"
status: accepted
created: "2026-09-30T20:07:50Z"
updated: "2026-10-01T02:29:45Z"
kind: feature
size: medium
depends_on: ["G-260929-gm3m4"]
relates_to: ["G-260930-tcc9w", "G-260930-60c3d", "G-260930-84fnb", "G-260930-npw49"]
candidate: "a8cc52da94fe98c38dbfe764c1cd59edf4ff4dd8"
approved: "a8cc52da94fe98c38dbfe764c1cd59edf4ff4dd8"
approved_by: owner
approved_context: "sha256:0e13f723d25b6c2009072c345c3ac4a4861a4c43bea7eda82d3183f6cfb1059e"
---

## Outcome

Each local delivery has one execution workspace, rooted in the delivered
target. Completed work is cleaned up safely by default, with an explicit
keep option. Kept squashed branches cannot become another delivery.

Owner decision:
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md), 2026-09-30.
This is the bounded follow-up to
[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md), not a repeat of
its schema migration or delivery verifier.

## Scope and constraints

Apply the workspace and admission rules in the
[contracts design](G-260930-84fnb-portable-workflow-contra.md). A delivery
may contain one item or an explicitly combined selection. Its external
prerequisites must be delivered into its execution base. Refuse inherited
unfinished implementation from another delivery; preserve ordinary branch
inspection. Before delivery, resume and fix the same workspace.

Prepare selected committed proposal records on the fresh target-based
workspace without inheriting a shaping branch's implementation. Preserve
IDs, paths, exact source revisions and provenance. The design specifies
bounded admission of records and required knowledge, conflict handling and
validation; no arbitrary branch copying or automatic authority expansion.

After delivery, retire that workspace from execution. Cleanup is on by
default, with a keep option, and runs only after the delivery is known,
evidence retained and its owner has exited. Preserve changed tips, dirty or
untracked files, replaced paths and live owners. A refused cleanup does not
undo Done or prevent the next fresh workspace when otherwise safe. A
reopened item uses a fresh workspace with an identity that cannot collide
with a kept one; inspection of old evidence stays available.

Keep schema 4 acceptance, cheap standing and optional audit. No history walk
or audit enters board/list/dependency reads. Do not add special continuation
merge bases or recoveries for repeated delivery from squashed branches.
Remove obsolete handling only after supported-path and refusal evidence
shows what replaces it. Multi-delivery orchestration belongs to
[G-260930-yfh91](G-260930-yfh91-continue-an-authorized-s.md).

## Acceptance

1. A single item and a combined group start from the delivered target, retain
   exact assignment inputs, and deliver through the existing squash
   operation. A committed proposal existing only on a shaping branch is
   admitted without its unrelated files or code; ambiguous or changed
   inputs wait before execution.
2. Undelivered external prerequisites and an already delivered workspace
   are refused before implementation or another delivery. A kept branch's
   new commits cannot silently become a second squash. Pre-delivery fixes
   and interrupted attempts still resume normally.
3. Successful delivery automatically cleans up an unchanged, unowned
   workspace after retaining evidence. Keep mode, dirty/untracked files,
   extra commits, live ownership and replacement paths preserve the work.
   Retry after target advance or a cleanup failure cannot redeliver.
4. Reopening delivered work starts fresh from the target, preserves old
   judgment and uses a distinct workspace identity. Old kept work remains
   inspectable without becoming the active execution source by accident.
5. A disposable local walkthrough exercises one delivery, a combined group,
   keep, cleanup, interruption and refusal of a reused squashed branch.
   Present the actual commands and history for owner judgment before the
   provider work depends on this boundary. Record real-agent observations
   only under an assigned usage mandate; fixtures alone are labelled so.
6. Commands, work guide, model/design, repository branch policy, entrypoint
   revision where needed and CLI/TUI controls agree. Changed defaults have
   explicit compatibility handling; no silent reinterpretation of an
   existing assignment. No new schema is presumed necessary.

## Dependencies

Depends on G-260929-gm3m4 for the delivered schema, standing, squash and
retention operations. It narrows their supported workflow before provider
extraction touches the same workspace/selection code.

## Evidence

Branch `worktree-G-260930-0s29t` from base `e17c646` (main), following
[the plan](G-260930-h29p8-plan-for-g-260930-0s29t.md) at `8730119`, from
this record at `sha256:dbb29e50`. A headless attempt made `8730119` plan,
`1c15afc` code, `9195b45` documents, `f7a2bb9` tests and the review fixes
`31fa7d2` and `195e235`, and stopped at the review cap. The owner then
allowed further rounds: `5bc4fbd`, `14e54f8` and `0123ce7`. Behaviour:

- A branch holding a delivery is retired: it contains a retained submission
  the target lacks, whose squash is on the target's first parents.
  `standing.Retired` is the one check, whichever work the delivery was of,
  and every verb that would act asks it; no read walks history.
- `run`: a new workspace starts from the target's tip (HEAD without a
  target); from a shaping branch it admits, as one
  `chore: admit IDS from COMMIT` commit with `Grove-Admitted-From`, only the
  selected records the target lacks, the records their relationships name
  and the plans and questions naming them. A differing selected record,
  other work past proposed, or an uncommitted admitted file is refused. A
  retired default name moves on to `-2`, `-3` with a preview note; a
  retired `--branch` is refused. `--keep` is recorded in `attempt.json`.
- `approve` and `feedback` refuse in a retired checkout, and so `resolve`
  does before its feedback is written. `feedback` on the target, where
  delivered work reopens, says to start a fresh workspace from it.
- `integrate`: never delivers from a retired branch, in review there or
  accepted; cleans up by default, kept by `--keep`, the newest attempt's
  `run --keep` (unless `--cleanup` or the board's `y`), a running or
  orphaned attempt, or a replaced path, besides Git's own refusals. A kept
  workspace is a fact, exit 0; a rerun retries only the cleanup. The
  sweep's delegated delivery cleans up too.

Decisions taken in review, beyond the plan: retirement is not scoped to the
work a branch delivered, since a kept branch could otherwise deliver other
work's squash; `approve` is refused on a retired branch, since an acceptance
there can never be delivered, while a hand edit with `update` stays
possible and `integrate` is the backstop; `--cleanup` overrides `run --keep`.
[commands.md](../docs/commands.md#retired-branches) owns the rule.

Against the acceptance:

1. Target base, exact inputs (admitted records' revisions and source in
   `attempt.json` and the digest), squash for one item and a group, and
   admission without code: `TestWorkspaceStartsFromTheTarget`,
   `TestWorkspaceAdmissionCommit`; walkthrough sections 1, 2 and 5.
2. Undelivered outside prerequisites still wait at the (target) base; a
   delivered workspace is refused by every verb, for its own work and for
   other work on it (`TestWorkspaceRetiredIdentity`,
   `TestNothingIsJudgedOnARetiredBranch`, `TestResolveRefusals`,
   `TestIntegrateRefusesAKeptDeliveredBranch`; walkthrough sections 3 and
   6); an interrupted delivery retries, and an interrupted attempt resumes
   its workspace (section 4).
3. `TestIntegrateSquashesAndRetainsEvidence`,
   `TestIntegrateCleanupKeepsWhatGitOrTheSessionHolds`,
   `TestIntegrateCleanupKeepsWhatAnAttemptHolds`,
   `TestIntegrateRecoversAfterTheTargetAdvanced`; walkthrough section 4.
4. Reopen on the target, fresh `-2` from the target, old judgment and the
   kept branch left inspectable (walkthrough section 3; board `R` path in
   `TestWorkspaceRetiredIdentity`).
5. The [owner walkthrough](#owner-walkthrough) below explains the six
   scenarios, actual commands and observed history. Fixtures only; no real
   agent was used. The owner's judgment remains separate from the fixture
   results and from approval of the candidate.
6. `--help`, README, docs/commands.md, docs/board.md, the work guide's step
   3 and judging part, docs/record-design.md, CLAUDE.md's branch names and
   retained-refs rule, and the Integration term agree; `--cleanup` stays
   accepted and now overrides `run --keep`; existing branches continue
   where they are. No schema change; the entrypoints need nothing new, so
   no revision.

Verification: `go vet ./...` and `gofmt -l .` clean and
`go run ./cmd/grove check` OK, 279 records, at `0123ce7`. One
`go test -count=1 -timeout 120s ./...` passed every package, the terminal
checks included, at `14e54f8`. At `0123ce7`, which changes one sentence of
the work guide, the same run at load average 12 failed `TestOwnerLost` and
the terminal check's `attempt_lifecycle` on timing; `internal/attempt` and
`internal/tui` each passed rerun alone. `internal/integrate` under `-short`
measures 3.7 to 4.0 s alone.

Review: [G-260930-3yyg8](G-260930-3yyg8-review-of-g-260930-0s29t.md), five
rounds: 8, 3, 2, 3, then 1 open finding at `14e54f8`, every one fixed. The
last, a sentence of the work guide, is fixed in `0123ce7` and checked by its
author, not independently.

Limits:

- The retired check finds deliveries through local `refs/grove/submitted/`
  refs. A clone that fetched a kept branch without them cannot tell, until
  it fetches `refs/grove/*`.
- Once delivered work is reopened on the target, `integrate ID --cleanup`
  refuses its kept branch as retired; only Git removes it.
- The sweep learns a branch is retired when its delegated `approve` is
  refused, after it verified the merge.
- The board was read, not driven by hand, beyond the terminal checks.

### Documentation follow-up, 2026-09-30

At the owner's request, the recorded fixture walkthrough was rewritten as
six scenarios and checked against its original script/log. The command
reference now owns a workspace operation table. The contract, milestone and
provider/review proposals reflect the authorized sequence through
G-261001-mbjwz and G-261001-62n4p. No runtime code changed in this follow-up.
The owner accepted the walkthrough behavior as recorded below.

Self-check of these documentation changes: `go vet ./...`, `gofmt -l .`,
`git diff --check` and `go run ./cmd/grove check` passed (282 records);
124 local links and code fences in nine changed documents checked. One
`go test -count=1 -timeout 120s ./...` passed every package, including the
terminal checks. This is a documentation self-check and fresh regression
run, not another independent review of the implementation. The historical
review and its examined commit remain unchanged.

## Owner walkthrough

This asks the owner to judge **whether the demonstrated local behavior is
what they want**. It does not ask them to inspect every Git object or rerun
the commands. Agreement with the behavior is distinct from accepting the
candidate, its known limits, or merging it.

Evidence: the original disposable walkthrough used the CLI built from
`14e54f8` and a `GROVE_CLAUDE` shell fixture that commits one change per work
ID and hands off. It did not exercise a real coding agent or independent
agent review. The original script and log were read again during the
2026-09-30 owner follow-up; they are `walk.sh` and `walk.log` under:

```text
/private/tmp/claude-501/-Users-mascah-GitHub-mascah-grove--claude-worktrees-worktree-G-260930-0s29t/a93df3ab-c812-4bef-96c3-1fce4ecb5663/scratchpad
```

These temporary files are local evidence, not a portable artifact. The
commands and results summarized here remain in the record. A, B, C, etc.
stand for the fixture's generated work IDs. Approval runs in the work's
checkout; integration runs in the target checkout. All fixture commands
used explicit absolute `--project` paths.

### 1. One change, one delivery

1. `grove run A` creates a workspace from the current `main`. The fixture
   makes the change and hands it into review with a candidate commit.
2. `grove approve A "Fixture: ok."` records acceptance in that workspace.
3. `grove integrate A` delivers the result to `main` as one squash commit.
4. Grove reports Done, retains the original evidence and removes the
   completed branch and workspace.

Observed: the implementation and bookkeeping commits on the work branch
became one `feat: greet in one delivery` commit on `main`; no work branch
remained. Judge whether that history and default cleanup match the intent.

### 2. Two items deliberately delivered together

1. `grove run B C` uses one workspace and one final candidate for both.
2. Approve B and C separately against that shared candidate.
3. `grove integrate C` delivers the whole accepted group once.

Observed: both read Done, one commit `feat: first of a group` reaches
`main`, and the shared workspace is removed. Both members' acceptance
travels with the delivery. This demonstrates the current combined mode;
separate deliveries from one multi-item launch remain later work.

### 3. Keep a workspace, then reopen work

1. `grove run D --keep`, approve D and integrate D. Grove leaves the
   completed workspace available for inspection.
2. `feedback D` in that old workspace is refused. The fixture deliberately
   used `update` and ordinary Git edits to put more work there anyway;
   `approve`, an explicitly named `run`, and `integrate` still refused it.
3. `grove feedback D "Fixture: reopened on main."` in the target checkout
   reopens the work properly.
4. `grove run D` creates `worktree-D-2` from the updated target. Approve and
   deliver the new candidate from that workspace.

Observed: the old kept workspace and its edits survive; the second delivery
uses a fresh workspace. Reopening adds an explicit metadata commit on main.
Judge whether Keep means inspection, and whether this fresh-start behavior
is the right way to request more work after a delivery.

### 4. Interruption and cleanup that preserves extra files

1. Start E, `grove stop ATTEMPT`, then `grove run E` again before delivery.
   Grove reuses E's unfinished workspace.
2. After handoff and approval, add an untracked `scratch.txt` there.
3. `grove integrate E` delivers successfully but keeps the workspace because
   removing it would lose the scratch file.
4. Remove the scratch file, then repeat `grove integrate E`.

Observed: the second integration says already Done and only cleans up;
there is no duplicate delivery. Judge the distinction: unfinished work
resumes in place; cleanup trouble does not turn a successful delivery into
failed work.

### 5. Admit a proposal without importing unrelated implementation

1. On a shaping branch, create and commit F's proposed record and a separate
   `shaping.txt` code file.
2. `grove run F --dry-run` shows which records it will admit.
3. `grove run F` starts a workspace from `main`, copies the proposal record
   with its source provenance, and runs the fixture implementation.
4. Inspect the diff, approve F and integrate it.

Observed: the workspace contains F's record and implementation, not
`shaping.txt`. Judge whether this preserves shaping context without treating
the shaping branch as an implementation base.

### 6. Refuse an old stacked branch without destroying its contents

1. Implement and approve P. Manually create H's branch on P's branch and
   run H there, exercising compatibility with an existing branch.
2. Deliver P to `main`.
3. Try `feedback`, `approve`, `resolve` (after a target conflict), and
   `integrate` for H on its old branch.
4. Preview a normal new `grove run H`.

Observed: each mutation attempt is refused naming P's delivery, H's branch
tip stays unchanged, and no further attempt starts there. The fresh preview
chooses `worktree-H-2`. This is evidence of refusal and preservation, not a
supported recommendation to stack independently managed work branches.

### Resulting history and remaining judgment

The original run ended with this main history, newest first:

```text
feat: main touches the same file
feat: delivered under a stack
feat: shaped on a branch
feat: interrupted attempt
feat: kept after delivery
docs(D): set status=active unset approved approved_by approved_context note
feat: kept after delivery
feat: first of a group
feat: greet in one delivery
docs: propose work
chore: init grove
```

The first entry was the deliberately conflicting target edit in scenario 6.
The `docs(D)` entry was the deliberate reopening on main in scenario 3.
`grove check --deliveries` proved all seven deliveries. Only the kept D
workspace and refused H workspace remained.

Owner walkthrough judgment: accepted on 2026-09-30 (owner's local date).
After reading the six scenarios, the owner answered "Accept the walkthrough
behavior" to the question explicitly limited to those local outcomes. This
satisfies the behavior judgment in acceptance 5; it is not approval of the
candidate, permission to merge, or acceptance of the missing-evidence limit.
The owner also authorized
[the simplification follow-up](G-261001-fg05z-simplify-workflow-bounda.md).

## Next

In review. The owner accepted the local walkthrough behavior. The candidate
includes the operation table and authorized follow-up records. Present that
exact candidate and its missing-evidence limitation for approval and
delivery. The original implementation candidate remains in history; the new
documents are included explicitly in the refreshed candidate.

The ordered follow-through is recorded in
[G-261001-mbjwz](G-261001-mbjwz-separate-blocking-review.md) and
[G-261001-62n4p](G-261001-62n4p-isolate-workspace-eligib.md), then the existing
provider work. Their dependency edges own that order. This conversation
has not authorized merging or paid provider trials. The milestone retains
the adopter migration hold.

Verdict on candidate a8cc52d, 2026-10-01: Owner approved candidate a8cc52da94fe98c38dbfe764c1cd59edf4ff4dd8 and authorized local integration with --keep in this conversation. The six local walkthrough behaviors are accepted; the documented imported-branch evidence limitation is assigned to G-261001-62n4p before provider work.
