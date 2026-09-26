---
id: "G-260921-czt8x"
type: page
title: "Identity and path migration map"
created: "2026-09-21T21:11:16Z"
updated: "2026-09-21T21:11:16Z"
---

W-029 reconciled every record, legacy plan and legacy review into neutral IDs
directly under `grove/` on 2026-09-21. This page is the one durable mapping
from an old ID or path to its current counterpart. To read an old commit, use
the old ID or path there with that commit's own CLI (`go run ./cmd/grove`);
to find what it became, look it up here or search `formerly:`. Numbers follow
document date (a record's `created`, a legacy document's first commit, which
became its `created` on 2026-09-22), ties
broken on the old ID or path; the order conveys no authority or priority.

| Old ID | Old path | ID | Path | Document date |
| --- | --- | --- | --- | --- |
| D-001 | `grove/decisions/D-001-starter-defaults.md` | G-260919-6mpmw | `grove/G-260919-6mpmw-adopt-the-starter-record.md` | 2026-09-19T14:08:40Z |
| Q-001 | `grove/questions/Q-001-branch-versions.md` | G-260919-8jb5s | `grove/G-260919-8jb5s-how-should-the-board-pre.md` | 2026-09-19T14:08:40Z |
| W-001 | `grove/work/W-001-inspect-records.md` | G-260919-rt9h9 | `grove/G-260919-rt9h9-inspect-grove-project-re.md` | 2026-09-19T14:08:40Z |
| D-002 | `grove/decisions/D-002-sequential-ids.md` | G-260919-4h6pn | `grove/G-260919-4h6pn-use-shared-sequential-id.md` | 2026-09-19T14:32:32Z |
| - | `docs/plans/W-001-inspection.md` | G-260919-10zeb | `grove/G-260919-10zeb-inspection-implementatio.md` | 2026-09-19T15:19:41Z |
| D-003 | `grove/decisions/D-003-allocator-mechanism.md` | G-260919-5f89v | `grove/G-260919-5f89v-allocate-ids-with-flock.md` | 2026-09-19T15:25:51Z |
| W-002 | `grove/work/W-002-create-records.md` | G-260919-92n2y | `grove/G-260919-92n2y-create-records-with-shar.md` | 2026-09-19T15:25:51Z |
| - | `docs/plans/W-002-create.md` | G-260919-dmhpz | `grove/G-260919-dmhpz-record-creation-implemen.md` | 2026-09-19T15:31:29Z |
| W-003 | `grove/work/W-003-update-records.md` | G-260919-shnj5 | `grove/G-260919-shnj5-update-record-status-and.md` | 2026-09-19T15:36:19Z |
| W-004 | `grove/work/W-004-record-versions.md` | G-260919-zb0s8 | `grove/G-260919-zb0s8-inspect-record-versions.md` | 2026-09-19T17:49:58Z |
| W-005 | `grove/work/W-005-record-workspace.md` | G-260919-n9t4p | `grove/G-260919-n9t4p-locate-the-workspace-for.md` | 2026-09-19T17:50:01Z |
| - | `docs/plans/W-003-update.md` | G-260919-50rcx | `grove/G-260919-50rcx-record-update-implementa.md` | 2026-09-19T18:42:44Z |
| - | `docs/plans/W-004-W-005-coordination.md` | G-260919-qprdw | `grove/G-260919-qprdw-and-coordination-plan.md` | 2026-09-19T19:28:26Z |
| W-006 | `grove/work/W-006-workspace-provenance.md` | G-260919-8bbvy | `grove/G-260919-8bbvy-bind-workspace-routing-t.md` | 2026-09-19T20:13:44Z |
| W-007 | `grove/work/W-007-preserve-updates.md` | G-260919-z9w13 | `grove/G-260919-z9w13-preserve-accepted-frontm.md` | 2026-09-19T20:13:54Z |
| W-008 | `grove/work/W-008-git-paths.md` | G-260919-7qv4x | `grove/G-260919-7qv4x-preserve-git-paths-throu.md` | 2026-09-19T20:13:57Z |
| W-009 | `grove/work/W-009-terminal-picker.md` | G-260919-k7b8j | `grove/G-260919-k7b8j-browse-a-terminal-kanban.md` | 2026-09-19T20:23:01Z |
| - | `docs/plans/W-006-workspace-provenance.md` | G-260919-72c95 | `grove/G-260919-72c95-workspace-provenance-imp.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-007-preserve-updates.md` | G-260919-w7ccc | `grove/G-260919-w7ccc-source-preserving-update.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-008-git-paths.md` | G-260919-jyfpb | `grove/G-260919-jyfpb-git-path-identity-implem.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-009-terminal-picker.md` | G-260919-wxbsx | `grove/G-260919-wxbsx-terminal-kanban-and-vers.md` | 2026-09-19T20:32:24Z |
| - | `docs/reviews/2026-09-19-integrated-cli.md` | G-260919-zrk8t | `grove/G-260919-zrk8t-integrated-cli-review-20.md` | 2026-09-19T20:32:24Z |
| W-010 | `grove/work/W-010-work-handoffs.md` | G-260919-nddsf | `grove/G-260919-nddsf-prepare-reusable-work-in.md` | 2026-09-19T20:44:28Z |
| - | `docs/reviews/2026-09-19-predecessor-work.md` | G-260919-ph0w1 | `grove/G-260919-ph0w1-predecessor-work-review.md` | 2026-09-19T20:59:57Z |
| W-011 | `grove/work/W-011-shaping-entrypoint.md` | G-260919-04z88 | `grove/G-260919-04z88-shape-project-work-throu.md` | 2026-09-19T21:20:41Z |
| - | `docs/reviews/2026-09-19-shaping-and-runner-evidence.md` | G-260919-rr3ae | `grove/G-260919-rr3ae-shaping-and-headless-run.md` | 2026-09-19T21:24:18Z |
| - | `docs/plans/W-010-W-011-agent-handoffs.md` | G-260919-p1dgm | `grove/G-260919-p1dgm-agent-handoffs-implement.md` | 2026-09-19T21:39:23Z |
| - | `docs/reviews/2026-09-19-repairs-W-006-W-008.md` | G-260919-syk45 | `grove/G-260919-syk45-cli-repairs-to-evidence.md` | 2026-09-19T22:07:59Z |
| - | `docs/reviews/2026-09-19-board-W-009.md` | G-260919-zq52f | `grove/G-260919-zq52f-terminal-board-evidence.md` | 2026-09-19T23:28:01Z |
| W-012 | `grove/work/W-012-card-lineage.md` | G-260920-svpbc | `grove/G-260920-svpbc-show-a-work-item-s-linea.md` | 2026-09-20T04:36:57Z |
| W-013 | `grove/work/W-013-load-scaling.md` | G-260920-z8vfp | `grove/G-260920-z8vfp-keep-a-full-board-load-f.md` | 2026-09-20T04:36:57Z |
| - | `docs/reviews/2026-09-19-W-010-dogfood.md` | G-260920-j2eyp | `grove/G-260920-j2eyp-dogfooding-evidence.md` | 2026-09-20T05:50:31Z |
| - | `docs/plans/W-012-card-lineage.md` | G-260920-d1qjs | `grove/G-260920-d1qjs-card-lineage-plan.md` | 2026-09-20T14:55:55Z |
| - | `docs/reviews/2026-09-20-card-lineage-W-012.md` | G-260920-0sakc | `grove/G-260920-0sakc-card-lineage-evidence-20.md` | 2026-09-20T16:11:14Z |
| D-004 | `grove/decisions/D-004-interactive-adoption.md` | G-260921-tkdwh | `grove/G-260921-tkdwh-adopt-the-interactive-ad.md` | 2026-09-21T00:54:13Z |
| W-018 | `grove/work/W-018-interactive-adoption.md` | G-260921-407n6 | `grove/G-260921-407n6-complete-the-interactive.md` | 2026-09-21T00:54:13Z |
| W-019 | `grove/work/W-019-knowledge-artifacts.md` | G-260921-w9x25 | `grove/G-260921-w9x25-represent-domain-terms-a.md` | 2026-09-21T00:54:14Z |
| W-020 | `grove/work/W-020-review-lifecycle.md` | G-260921-9wkjt | `grove/G-260921-9wkjt-hand-implementation-cand.md` | 2026-09-21T00:54:14Z |
| W-021 | `grove/work/W-021-interactive-loop.md` | G-260921-9t178 | `grove/G-260921-9t178-prove-the-complete-inter.md` | 2026-09-21T00:54:14Z |
| W-022 | `grove/work/W-022-portable-bootstrap.md` | G-260921-5gz9a | `grove/G-260921-5gz9a-bootstrap-projects-with.md` | 2026-09-21T00:54:15Z |
| W-023 | `grove/work/W-023-nullsec-pilot.md` | G-260921-905y3 | `grove/G-260921-905y3-cut-nullsec-over-to-this.md` | 2026-09-21T00:54:15Z |
| W-024 | `grove/work/W-024-current-view.md` | G-260921-ms6ev | `grove/G-260921-ms6ev-derive-a-project-wide-cu.md` | 2026-09-21T00:54:15Z |
| W-025 | `grove/work/W-025-board-detail.md` | G-260921-k0mwk | `grove/G-260921-k0mwk-make-the-board-and-item.md` | 2026-09-21T00:54:15Z |
| W-026 | `grove/work/W-026-review-integration.md` | G-260921-jwk4e | `grove/G-260921-jwk4e-review-candidates-and-in.md` | 2026-09-21T00:54:16Z |
| W-027 | `grove/work/W-027-durable-attempt.md` | G-260921-h46pb | `grove/G-260921-h46pb-run-one-bounded-implemen.md` | 2026-09-21T00:54:16Z |
| W-028 | `grove/work/W-028-managed-runs.md` | G-260921-7trd7 | `grove/G-260921-7trd7-launch-and-inspect-manag.md` | 2026-09-21T00:54:16Z |
| - | `docs/plans/W-018-adoption-roadmap.md` | G-260921-466b5 | `grove/G-260921-466b5-adoption-roadmap.md` | 2026-09-21T01:08:42Z |
| - | `docs/reviews/2026-09-20-direction-evaluation.md` | G-260921-72chf | `grove/G-260921-72chf-first-days-evaluation-an.md` | 2026-09-21T01:08:42Z |
| - | `docs/plans/W-011-shaping-entrypoint.md` | G-260921-x53yt | `grove/G-260921-x53yt-plan-shaping-guide-and-g.md` | 2026-09-21T02:53:09Z |
| - | `docs/reviews/2026-09-20-W-011-shaping.md` | G-260921-ahbrj | `grove/G-260921-ahbrj-shaping-entrypoint-evide.md` | 2026-09-21T03:02:56Z |
| D-005 | `grove/decisions/D-005-typed-knowledge-records.md` | G-260921-e8bva | `grove/G-260921-e8bva-represent-terms-plans-an.md` | 2026-09-21T04:37:58Z |
| W-029 | `grove/work/W-029-migrate-knowledge.md` | G-260921-r491p | `grove/G-260921-r491p-reconcile-all-grove-cont.md` | 2026-09-21T04:37:58Z |
| - | `docs/plans/W-019-knowledge-artifacts.md` | G-260921-q09km | `grove/G-260921-q09km-plan-terms-plans-reviews.md` | 2026-09-21T04:55:41Z |
| T-001 | `grove/terms/T-001-work.md` | G-260921-vr8a8 | `grove/G-260921-vr8a8-work.md` | 2026-09-21T05:01:55Z |
| T-002 | `grove/terms/T-002-preparation.md` | G-260921-dqdde | `grove/G-260921-dqdde-preparation.md` | 2026-09-21T05:01:56Z |
| T-003 | `grove/terms/T-003-attempt.md` | G-260921-sth8q | `grove/G-260921-sth8q-attempt.md` | 2026-09-21T05:01:56Z |
| T-004 | `grove/terms/T-004-candidate.md` | G-260921-jatts | `grove/G-260921-jatts-candidate.md` | 2026-09-21T05:01:56Z |
| T-005 | `grove/terms/T-005-review.md` | G-260921-rz7bn | `grove/G-260921-rz7bn-review.md` | 2026-09-21T05:01:56Z |
| T-006 | `grove/terms/T-006-approval.md` | G-260921-btyck | `grove/G-260921-btyck-approval.md` | 2026-09-21T05:01:56Z |
| T-007 | `grove/terms/T-007-integration.md` | G-260921-3qgsf | `grove/G-260921-3qgsf-integration.md` | 2026-09-21T05:01:57Z |
| T-008 | `grove/terms/T-008-source.md` | G-260921-phehv | `grove/G-260921-phehv-source.md` | 2026-09-21T05:01:57Z |
| T-009 | `grove/terms/T-009-revision.md` | G-260921-vz0v3 | `grove/G-260921-vz0v3-revision.md` | 2026-09-21T05:01:57Z |
| R-001 | `grove/reviews/R-001-w-019-knowledge-records.md` | G-260921-q6e5n | `grove/G-260921-q6e5n-independent-review-of-kn.md` | 2026-09-21T05:59:58Z |
| D-006 | `grove/decisions/D-006-stable-knowledge.md` | G-260921-gtydy | `grove/G-260921-gtydy-keep-identity-and-placem.md` | 2026-09-21T15:40:59Z |
| W-030 | `grove/work/W-030-flexible-records.md` | G-260921-ebsby | `grove/G-260921-ebsby-decouple-record-identity.md` | 2026-09-21T15:40:59Z |
| P-001 | `grove/plans/P-001-w-030-flexible-records.md` | G-260921-6n3da | `grove/G-260921-6n3da-flexible-records-schema.md` | 2026-09-21T16:02:50Z |
| R-002 | `grove/reviews/R-002-w-030-flexible-records.md` | G-260921-kfd06 | `grove/G-260921-kfd06-independent-review-of-fl.md` | 2026-09-21T16:32:01Z |
| P-002 | `grove/plans/P-002-w-029-reconciliation.md` | G-260921-awvk8 | `grove/G-260921-awvk8-reconciliation-mapping-c.md` | 2026-09-21T21:00:03Z |
| - | `docs/restart-brief.md` | - | `grove/brief.md` | the brief, not a record |

## Kept in place

- `docs/prompts/*.txt`: three spent handoff prompts, historical evidence that
  the handoff work links. Their old IDs and paths are literals of their time.
- `README.md`, `AGENTS.md`, the skill adapters, `docs/record-model.md`,
  `docs/work-execution.md` and `docs/work-shaping.md`: functional homes, with
  references updated.

## Historical literals

Old identifiers remain in `formerly` fields, in this table, in Git branch
names such as `worktree-W-012` and in commit messages, which name what existed
then. The typed IDs named on these lines were left as written: they belong to
test fixtures, disposable trial clones or the predecessor project rather than
to this repository's records, or they sit in a verbatim quotation (the
owner's words, a commit subject, the 2026-09-19 renumbering table). Every
other typed ID that had a counterpart was rewritten, command transcripts
included:

- `grove/G-260919-dmhpz-record-creation-implemen.md:71: W-003, W-005`
- `grove/G-260919-dmhpz-record-creation-implemen.md:73: W-004`
- `grove/G-260919-dmhpz-record-creation-implemen.md:93: W-030`
- `grove/G-260919-dmhpz-record-creation-implemen.md:95: W-007`
- `grove/G-260920-0sakc-card-lineage-evidence-20.md:51: W-001`
- `grove/G-260919-zrk8t-integrated-cli-review-20.md:77: W-001`
- `grove/G-260919-zrk8t-integrated-cli-review-20.md:131: W-001`
- `grove/G-260919-zrk8t-integrated-cli-review-20.md:134: W-001`
- `grove/G-260920-svpbc-show-a-work-item-s-linea.md:29: W-006`
- `grove/G-260919-rr3ae-shaping-and-headless-run.md:24: W-003`
- `grove/G-260919-rr3ae-shaping-and-headless-run.md:38: W-003`
- `grove/G-260920-j2eyp-dogfooding-evidence.md:54: W-001, W-002, W-003`
- `grove/G-260919-p1dgm-agent-handoffs-implement.md:329: W-001, W-002, W-003`
- `grove/G-260919-p1dgm-agent-handoffs-implement.md:330: W-001, W-002, W-003`
- `grove/G-260919-p1dgm-agent-handoffs-implement.md:402: W-001`
- `grove/G-260919-p1dgm-agent-handoffs-implement.md:406: W-001`
- `grove/G-260919-p1dgm-agent-handoffs-implement.md:409: W-001`
- `grove/G-260919-5f89v-allocate-ids-with-flock.md:14: D-003`
- `grove/G-260919-zq52f-terminal-board-evidence.md:52: W-001`
- `grove/G-260919-zq52f-terminal-board-evidence.md:53: W-002`
- `grove/G-260919-zq52f-terminal-board-evidence.md:66: W-001`
- `grove/G-260919-4h6pn-use-shared-sequential-id.md:40: W-001`
- `grove/G-260919-4h6pn-use-shared-sequential-id.md:41: Q-001`
- `grove/G-260919-4h6pn-use-shared-sequential-id.md:42: D-001`
- `grove/G-260919-jyfpb-git-path-identity-implem.md:88: W-001`
- `grove/G-260919-10zeb-inspection-implementatio.md:117: W-001`
- `grove/G-260919-ph0w1-predecessor-work-review.md:32: W-002`
- `grove/G-260919-qprdw-and-coordination-plan.md:264: W-001`
- `grove/G-260919-qprdw-and-coordination-plan.md:293: W-001`
- `grove/G-260921-q6e5n-independent-review-of-kn.md:31: T-001`
- `grove/G-260921-awvk8-reconciliation-mapping-c.md:62: W-001`
- `grove/G-260919-72c95-workspace-provenance-imp.md:63: W-001`
- `grove/G-260919-72c95-workspace-provenance-imp.md:106: W-001`
- `grove/G-260919-w7ccc-source-preserving-update.md:56: W-001`
- `grove/G-260919-w7ccc-source-preserving-update.md:85: W-001`
- `grove/G-260919-w7ccc-source-preserving-update.md:86: W-001`
- `grove/G-260919-w7ccc-source-preserving-update.md:119: W-001`
- `grove/G-260919-8jb5s-how-should-the-board-pre.md:43: W-003`
- `grove/G-260919-50rcx-record-update-implementa.md:144: W-002`
- `grove/G-260921-ahbrj-shaping-entrypoint-evide.md`: every `W-029` and `W-030`, which the trial wrote in its own clones before these numbers existed here
- `grove/G-260919-z9w13-preserve-accepted-frontm.md:106: W-003`
- `grove/G-260921-ebsby-decouple-record-identity.md:130-131: W-019, R-001, W-030, P-001, D-006` (the inputs of a `convert` rehearsal, which takes typed IDs)
- `grove/G-260921-gtydy-keep-identity-and-placem.md:80: W-019` (evidence about schema 2's folder rule)

A typed ID shown as an example of the retired spelling (`W-001` in the README
and the record model) is an illustration, not a reference. A typed ID inside a filename or path (`docs/prompts/W-009-implementation.txt`,
a fixture's `grove/work/W-001-first.md`, a hypothetical `work/W-019/plan.md`)
was never rewritten either: it names a file, not a record. A typed ID with no row in the table (`W-014` to `W-017`, `W-090`, `Q-002` and
the like) never named a record here: it is a fixture's or the predecessor's.
Records and evidence dated before this migration also name the folders of
their time (`grove/work/`, `docs/plans/`, `docs/reviews/`) where they describe
what was then true; nothing lives there now.
