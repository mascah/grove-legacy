---
id: "G-261001-62n4p"
type: work
title: "Isolate workspace eligibility before portable execution"
status: proposed
created: "2026-10-01T02:17:16Z"
updated: "2026-10-01T02:22:08Z"
kind: refactor
size: medium
depends_on: ["G-260930-0s29t", "G-261001-mbjwz"]
relates_to: ["G-261001-fg05z", "G-260928-y2p5h", "G-260930-84fnb", "G-260930-60c3d"]
---

## Outcome

Execution and delivery share a small workspace boundary that answers where
work may run, whether an existing workspace remains eligible, and whether it
may be removed. Provider extraction can use it without inheriting Claude
protocol details or reimplementing branch lifecycle rules.

Owner authorization: [G-261001-fg05z](G-261001-fg05z-simplify-workflow-bounda.md).

## Constraints

Use [the contracts](G-260930-84fnb-portable-workflow-contra.md) and the
[operation table](../docs/commands.md#workspace-operation-table). Preserve
record identity, applicable approval, cheap standing, exact sources,
single-writer ownership, interruption recovery and extra user files.

Observed at a448a6c: `attempt.prepare` selects/prepares the workspace while
the same package runs and decodes Claude; `integrate.keeping` imports that
package to read ownership and Keep; `integrate.accepted` still selects among
branches and `integrate.Run` handles carried candidates. `standing.Retired`
centralizes the retirement rule but depends on local evidence refs.

The plan must inventory each affected rule and caller, separate the current
supported path from compatibility handling for existing assignments, and
identify the obsolete recoveries actually removed. Keep branch inspection
and preserved edits; do not substitute arbitrary branch recency for exact
source selection or add a second authoritative work status. A new package
alone does not meet this outcome.

Imported-branch eligibility is part of this boundary. A fresh clone must
still read Done and start fresh work without audit refs. An imported old
implementation branch may run only when its eligibility is established;
missing information must not be treated as evidence that it is reusable.
The plan compares deriving the answer from transferred Git facts with
requiring explicit evidence transfer/reconciliation. Choose the smallest
mechanism demonstrated by fixtures; do not build a general transport system,
automatically fetch, or add historical checks to ordinary board/list reads.

Scope excludes a second provider, hosted delivery, the LLM judge and
multi-delivery orchestration. Their existing records consume this outcome.
No new schema is presumed necessary.

Depends on G-260930-0s29t for the judged local workspace behavior, and on
G-261001-mbjwz so this refactor uses the settled blocker/follow-up review
contract instead of repeating the broad review churn.

## Acceptance

1. Map every workspace mutation caller to the shared eligibility/ownership
   facts. CLI/TUI actions agree for fresh, interrupted, review-fix, delivered
   and reopened work. Missing/uncertain facts have an explicit wait outcome.
2. Delivery cleanup no longer depends on the provider runner's protocol or
   orchestration. Exercise running/lost owners, Keep, extra commits, dirty,
   untracked/ignored files, replaced paths and retry after target advance.
3. Show which branch-selection or recovery rules were removed, which remain
   necessary for supported operations, and which are bounded compatibility
   handling. Demonstrate replacement/refusal before removing safeguards.
4. A disposable source repository and fresh clone demonstrate Done without
   audit refs, a fresh launch, and an imported kept branch with and without
   evidence. No missing-evidence case silently enables a second delivery;
   transferred unfinished work can continue under a fresh mandate when its
   required facts are present. Inspection preserves all copies.
5. One concise walkthrough shows the ordinary path, interruption, delivery,
   reopening and evidence-unavailable wait. Label fixture versus real-agent
   evidence. The owner can explain which workspace each step uses.
6. Keep the owning command/model documents and operation table accurate.
   Targeted verification covers changed rules and callers; unrelated
   provider features and a new completion migration are excluded.

## Next

Authorized for follow-through in the current conversation. After the two
prerequisites deliver, prepare the boundary and removal inventory in a
separate work branch, implement and demonstrate it with the existing runner.
G-260928-y2p5h consumes the resulting interface; its provider feature work
waits for this outcome.
