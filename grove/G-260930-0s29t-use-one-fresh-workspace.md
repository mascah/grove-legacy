---
id: "G-260930-0s29t"
type: work
title: "Use one fresh workspace per delivery with safe automatic cleanup"
status: proposed
created: "2026-09-30T20:07:50Z"
updated: "2026-09-30T20:22:02Z"
kind: feature
size: medium
depends_on: ["G-260929-gm3m4"]
relates_to: ["G-260930-tcc9w", "G-260930-60c3d", "G-260930-84fnb", "G-260930-npw49"]
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

## Next

Proposed and unassigned. Review the reconciled contracts, then plan this
bounded change against the current implementation. The adopter migration
hold belongs to the milestone; this record does not release it.
