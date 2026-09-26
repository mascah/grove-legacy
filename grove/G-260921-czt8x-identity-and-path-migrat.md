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

## Legacy IDs to date-form IDs

On 2026-09-26 `grove renumber` gave every record with a three-digit legacy
ID a date-form ID dated by its `created`, and a filename from that ID and its
title, and rewrote every reference beneath `grove/` and in the attempt store;
references elsewhere in the repository were repaired from its output
([G-260926-vkv48](G-260926-vkv48-rename-legacy-records-to.md), under
[G-260926-yvjy6](G-260926-yvjy6-retire-legacy-ids-by-ren.md)). Commit
subjects, branch names and anything else in Git history keep the legacy IDs;
this table resolves them. IDs that followed an escape in an attempt's JSON
text (`\nG-153`) were missed by that run and rewritten in a second pass
with the same map once the command was fixed; only fixture IDs drawn after
terminal escapes remain there. The table above already shows each typed ID's
date-form ID. Nullsec's own numbers, cited in its cutover records, are not
this repository's and were left as they were.

| Legacy ID | Legacy path | ID | Path |
| --- | --- | --- | --- |
| G-001 | `grove/G-001-starter-defaults.md` | G-260919-6mpmw | `grove/G-260919-6mpmw-adopt-the-starter-record.md` |
| G-002 | `grove/G-002-branch-versions.md` | G-260919-8jb5s | `grove/G-260919-8jb5s-how-should-the-board-pre.md` |
| G-003 | `grove/G-003-inspect-records.md` | G-260919-rt9h9 | `grove/G-260919-rt9h9-inspect-grove-project-re.md` |
| G-004 | `grove/G-004-sequential-ids.md` | G-260919-4h6pn | `grove/G-260919-4h6pn-use-shared-sequential-id.md` |
| G-005 | `grove/G-005-inspection-plan.md` | G-260919-10zeb | `grove/G-260919-10zeb-inspection-implementatio.md` |
| G-006 | `grove/G-006-allocator-mechanism.md` | G-260919-5f89v | `grove/G-260919-5f89v-allocate-ids-with-flock.md` |
| G-007 | `grove/G-007-create-records.md` | G-260919-92n2y | `grove/G-260919-92n2y-create-records-with-shar.md` |
| G-008 | `grove/G-008-create-plan.md` | G-260919-dmhpz | `grove/G-260919-dmhpz-record-creation-implemen.md` |
| G-009 | `grove/G-009-update-records.md` | G-260919-shnj5 | `grove/G-260919-shnj5-update-record-status-and.md` |
| G-010 | `grove/G-010-record-versions.md` | G-260919-zb0s8 | `grove/G-260919-zb0s8-inspect-record-versions.md` |
| G-011 | `grove/G-011-record-workspace.md` | G-260919-n9t4p | `grove/G-260919-n9t4p-locate-the-workspace-for.md` |
| G-012 | `grove/G-012-update-plan.md` | G-260919-50rcx | `grove/G-260919-50rcx-record-update-implementa.md` |
| G-013 | `grove/G-013-coordination-plan.md` | G-260919-qprdw | `grove/G-260919-qprdw-and-coordination-plan.md` |
| G-014 | `grove/G-014-workspace-provenance.md` | G-260919-8bbvy | `grove/G-260919-8bbvy-bind-workspace-routing-t.md` |
| G-015 | `grove/G-015-preserve-updates.md` | G-260919-z9w13 | `grove/G-260919-z9w13-preserve-accepted-frontm.md` |
| G-016 | `grove/G-016-git-paths.md` | G-260919-7qv4x | `grove/G-260919-7qv4x-preserve-git-paths-throu.md` |
| G-017 | `grove/G-017-terminal-picker.md` | G-260919-k7b8j | `grove/G-260919-k7b8j-browse-a-terminal-kanban.md` |
| G-018 | `grove/G-018-workspace-provenance-plan.md` | G-260919-72c95 | `grove/G-260919-72c95-workspace-provenance-imp.md` |
| G-019 | `grove/G-019-preserve-updates-plan.md` | G-260919-w7ccc | `grove/G-260919-w7ccc-source-preserving-update.md` |
| G-020 | `grove/G-020-git-paths-plan.md` | G-260919-jyfpb | `grove/G-260919-jyfpb-git-path-identity-implem.md` |
| G-021 | `grove/G-021-terminal-picker-plan.md` | G-260919-wxbsx | `grove/G-260919-wxbsx-terminal-kanban-and-vers.md` |
| G-022 | `grove/G-022-integrated-cli-review.md` | G-260919-zrk8t | `grove/G-260919-zrk8t-integrated-cli-review-20.md` |
| G-023 | `grove/G-023-work-handoffs.md` | G-260919-nddsf | `grove/G-260919-nddsf-prepare-reusable-work-in.md` |
| G-024 | `grove/G-024-predecessor-work-review.md` | G-260919-ph0w1 | `grove/G-260919-ph0w1-predecessor-work-review.md` |
| G-025 | `grove/G-025-shaping-entrypoint.md` | G-260919-04z88 | `grove/G-260919-04z88-shape-project-work-throu.md` |
| G-026 | `grove/G-026-shaping-and-runner-evidence-review.md` | G-260919-rr3ae | `grove/G-260919-rr3ae-shaping-and-headless-run.md` |
| G-027 | `grove/G-027-agent-handoffs-plan.md` | G-260919-p1dgm | `grove/G-260919-p1dgm-agent-handoffs-implement.md` |
| G-028 | `grove/G-028-repairs-review.md` | G-260919-syk45 | `grove/G-260919-syk45-cli-repairs-to-evidence.md` |
| G-029 | `grove/G-029-board-review.md` | G-260919-zq52f | `grove/G-260919-zq52f-terminal-board-evidence.md` |
| G-030 | `grove/G-030-card-lineage.md` | G-260920-svpbc | `grove/G-260920-svpbc-show-a-work-item-s-linea.md` |
| G-031 | `grove/G-031-load-scaling.md` | G-260920-z8vfp | `grove/G-260920-z8vfp-keep-a-full-board-load-f.md` |
| G-032 | `grove/G-032-dogfood-review.md` | G-260920-j2eyp | `grove/G-260920-j2eyp-dogfooding-evidence.md` |
| G-033 | `grove/G-033-card-lineage-plan.md` | G-260920-d1qjs | `grove/G-260920-d1qjs-card-lineage-plan.md` |
| G-034 | `grove/G-034-card-lineage-review.md` | G-260920-0sakc | `grove/G-260920-0sakc-card-lineage-evidence-20.md` |
| G-035 | `grove/G-035-interactive-adoption.md` | G-260921-tkdwh | `grove/G-260921-tkdwh-adopt-the-interactive-ad.md` |
| G-036 | `grove/G-036-interactive-adoption.md` | G-260921-407n6 | `grove/G-260921-407n6-complete-the-interactive.md` |
| G-037 | `grove/G-037-knowledge-artifacts.md` | G-260921-w9x25 | `grove/G-260921-w9x25-represent-domain-terms-a.md` |
| G-038 | `grove/G-038-review-lifecycle.md` | G-260921-9wkjt | `grove/G-260921-9wkjt-hand-implementation-cand.md` |
| G-039 | `grove/G-039-interactive-loop.md` | G-260921-9t178 | `grove/G-260921-9t178-prove-the-complete-inter.md` |
| G-040 | `grove/G-040-portable-bootstrap.md` | G-260921-5gz9a | `grove/G-260921-5gz9a-bootstrap-projects-with.md` |
| G-041 | `grove/G-041-nullsec-pilot.md` | G-260921-905y3 | `grove/G-260921-905y3-cut-nullsec-over-to-this.md` |
| G-042 | `grove/G-042-current-view.md` | G-260921-ms6ev | `grove/G-260921-ms6ev-derive-a-project-wide-cu.md` |
| G-043 | `grove/G-043-board-detail.md` | G-260921-k0mwk | `grove/G-260921-k0mwk-make-the-board-and-item.md` |
| G-044 | `grove/G-044-review-integration.md` | G-260921-jwk4e | `grove/G-260921-jwk4e-review-candidates-and-in.md` |
| G-045 | `grove/G-045-durable-attempt.md` | G-260921-h46pb | `grove/G-260921-h46pb-run-one-bounded-implemen.md` |
| G-046 | `grove/G-046-managed-runs.md` | G-260921-7trd7 | `grove/G-260921-7trd7-launch-and-inspect-manag.md` |
| G-047 | `grove/G-047-adoption-roadmap-plan.md` | G-260921-466b5 | `grove/G-260921-466b5-adoption-roadmap.md` |
| G-048 | `grove/G-048-direction-evaluation-review.md` | G-260921-72chf | `grove/G-260921-72chf-first-days-evaluation-an.md` |
| G-049 | `grove/G-049-shaping-entrypoint-plan.md` | G-260921-x53yt | `grove/G-260921-x53yt-plan-shaping-guide-and-g.md` |
| G-050 | `grove/G-050-shaping-review.md` | G-260921-ahbrj | `grove/G-260921-ahbrj-shaping-entrypoint-evide.md` |
| G-051 | `grove/G-051-typed-knowledge-records.md` | G-260921-e8bva | `grove/G-260921-e8bva-represent-terms-plans-an.md` |
| G-052 | `grove/G-052-migrate-knowledge.md` | G-260921-r491p | `grove/G-260921-r491p-reconcile-all-grove-cont.md` |
| G-053 | `grove/G-053-knowledge-artifacts-plan.md` | G-260921-q09km | `grove/G-260921-q09km-plan-terms-plans-reviews.md` |
| G-054 | `grove/G-054-work.md` | G-260921-vr8a8 | `grove/G-260921-vr8a8-work.md` |
| G-055 | `grove/G-055-preparation.md` | G-260921-dqdde | `grove/G-260921-dqdde-preparation.md` |
| G-056 | `grove/G-056-attempt.md` | G-260921-sth8q | `grove/G-260921-sth8q-attempt.md` |
| G-057 | `grove/G-057-candidate.md` | G-260921-jatts | `grove/G-260921-jatts-candidate.md` |
| G-058 | `grove/G-058-review.md` | G-260921-rz7bn | `grove/G-260921-rz7bn-review.md` |
| G-059 | `grove/G-059-approval.md` | G-260921-btyck | `grove/G-260921-btyck-approval.md` |
| G-060 | `grove/G-060-integration.md` | G-260921-3qgsf | `grove/G-260921-3qgsf-integration.md` |
| G-061 | `grove/G-061-source.md` | G-260921-phehv | `grove/G-260921-phehv-source.md` |
| G-062 | `grove/G-062-revision.md` | G-260921-vz0v3 | `grove/G-260921-vz0v3-revision.md` |
| G-063 | `grove/G-063-knowledge-records-review.md` | G-260921-q6e5n | `grove/G-260921-q6e5n-independent-review-of-kn.md` |
| G-064 | `grove/G-064-stable-knowledge.md` | G-260921-gtydy | `grove/G-260921-gtydy-keep-identity-and-placem.md` |
| G-065 | `grove/G-065-flexible-records.md` | G-260921-ebsby | `grove/G-260921-ebsby-decouple-record-identity.md` |
| G-066 | `grove/G-066-flexible-records-plan.md` | G-260921-6n3da | `grove/G-260921-6n3da-flexible-records-schema.md` |
| G-067 | `grove/G-067-flexible-records-review.md` | G-260921-kfd06 | `grove/G-260921-kfd06-independent-review-of-fl.md` |
| G-068 | `grove/G-068-reconciliation-plan.md` | G-260921-awvk8 | `grove/G-260921-awvk8-reconciliation-mapping-c.md` |
| G-069 | `grove/G-069-migration-map.md` | G-260921-czt8x | `grove/G-260921-czt8x-identity-and-path-migrat.md` |
| G-070 | `grove/G-070-reconciliation-review.md` | G-260921-xs7wz | `grove/G-260921-xs7wz-independent-review-of-th.md` |
| G-071 | `grove/G-071-spawn-fewer-git-processes-per-in.md` | G-260922-9cbh6 | `grove/G-260922-9cbh6-spawn-fewer-git-processe.md` |
| G-072 | `grove/G-072-g-071-process-count-review.md` | G-260922-3cgay | `grove/G-260922-3cgay-process-count-review.md` |
| G-073 | `grove/G-073-review-lifecycle-plan.md` | G-260922-4fr84 | `grove/G-260922-4fr84-review-lifecycle-plan.md` |
| G-074 | `grove/G-074-review-lifecycle-review.md` | G-260922-fvqpv | `grove/G-260922-fvqpv-review-lifecycle-review.md` |
| G-075 | `grove/G-075-trial-plan-for-the-interactive-l.md` | G-260922-10dj4 | `grove/G-260922-10dj4-trial-plan-for-the-inter.md` |
| G-076 | `grove/G-076-filter-grove-list-by-status.md` | G-260922-w53bz | `grove/G-260922-w53bz-filter-grove-list-by-sta.md` |
| G-077 | `grove/G-077-g-076-independent-review-of-the.md` | G-260922-22hyw | `grove/G-260922-22hyw-independent-review-of-th.md` |
| G-078 | `grove/G-078-g-039-trial-evidence-for-the-int.md` | G-260922-08wxx | `grove/G-260922-08wxx-trial-evidence-for-the-i.md` |
| G-079 | `grove/G-079-update-a-record-by-hand-without.md` | G-260922-q3cr9 | `grove/G-260922-q3cr9-update-a-record-by-hand.md` |
| G-080 | `grove/G-080-portable-bootstrap-plan.md` | G-260922-6d6jg | `grove/G-260922-6d6jg-portable-bootstrap-plan.md` |
| G-081 | `grove/G-081-github-ci.md` | G-260922-jtsed | `grove/G-260922-jtsed-run-secure-ci-and-depend.md` |
| G-082 | `grove/G-082-portable-bootstrap-review.md` | G-260922-fqf3b | `grove/G-260922-fqf3b-portable-bootstrap-revie.md` |
| G-083 | `grove/G-083-g-079-review.md` | G-260922-0em47 | `grove/G-260922-0em47-review-of-optional-expec.md` |
| G-084 | `grove/G-084-review-of-g-081-ci-dependabot-an.md` | G-260922-q88hk | `grove/G-260922-q88hk-review-of-ci-dependabot.md` |
| G-089 | `grove/G-089-ignore-ambient-git-environment-w.md` | G-260922-g6e7p | `grove/G-260922-g6e7p-ignore-ambient-git-envir.md` |
| G-090 | `grove/G-090-review-of-g-089-ambient-git-envi.md` | G-260922-fze8c | `grove/G-260922-fze8c-review-of-ambient-git-en.md` |
| G-091 | `grove/G-091-nullsec-cutover-plan.md` | G-260922-r1dhw | `grove/G-260922-r1dhw-nullsec-cutover-mapping.md` |
| G-092 | `grove/G-092-nullsec-cutover-review.md` | G-260922-9d399 | `grove/G-260922-9d399-nullsec-cutover-review.md` |
| G-093 | `grove/G-093-current-view-plan.md` | G-260922-9tcff | `grove/G-260922-9tcff-project-wide-current-vie.md` |
| G-094 | `grove/G-094-current-view-review.md` | G-260922-ayftm | `grove/G-260922-ayftm-review-of-current-view.md` |
| G-095 | `grove/G-095-integration-target-review.md` | G-260922-cwns8 | `grove/G-260922-cwns8-review-of-integration-ta.md` |
| G-096 | `grove/G-096-g-043-board-and-detail-design-vi.md` | G-260922-1w0hn | `grove/G-260922-1w0hn-board-and-detail-design.md` |
| G-097 | `grove/G-097-review-of-g-043-board-detail-and.md` | G-260922-3fn26 | `grove/G-260922-3fn26-review-of-board-detail-a.md` |
| G-098 | `grove/G-098-g-044-review-integration-plan.md` | G-260923-twv25 | `grove/G-260923-twv25-review-actions-and-local.md` |
| G-099 | `grove/G-099-g-044-review-integration-review.md` | G-260923-r1w6p | `grove/G-260923-r1w6p-review-integration-revie.md` |
| G-100 | `grove/G-100-g-045-durable-attempt-plan.md` | G-260923-0t43m | `grove/G-260923-0t43m-durable-attempt-plan.md` |
| G-101 | `grove/G-101-attempt-mechanism.md` | G-260923-tnn5e | `grove/G-260923-tnn5e-run-attempts-as-a-grove.md` |
| G-102 | `grove/G-102-g-045-durable-attempt-review.md` | G-260923-ccda0 | `grove/G-260923-ccda0-durable-attempt-review.md` |
| G-103 | `grove/G-103-g-046-managed-runs-plan.md` | G-260923-stkc6 | `grove/G-260923-stkc6-managed-runs-plan.md` |
| G-104 | `grove/G-104-g-046-managed-runs-review.md` | G-260923-a8kzm | `grove/G-260923-a8kzm-managed-runs-review.md` |
| G-105 | `grove/G-105-restore-green-ci-linux-build-go.md` | G-260923-q7tm6 | `grove/G-260923-q7tm6-restore-green-ci-linux-b.md` |
| G-106 | `grove/G-106-g-105-green-ci-review.md` | G-260923-cd2ce | `grove/G-260923-cd2ce-green-ci-review.md` |
| G-107 | `grove/G-107-current-documentation.md` | G-260923-fwakw | `grove/G-260923-fwakw-reconcile-current-docume.md` |
| G-108 | `grove/G-108-workflow-evals.md` | G-260923-p5pt6 | `grove/G-260923-p5pt6-establish-behavioral-eva.md` |
| G-109 | `grove/G-109-attempts-usability.md` | G-260923-895zb | `grove/G-260923-895zb-make-attempts-easy-to-sc.md` |
| G-110 | `grove/G-110-external-preview.md` | G-260923-gsthp | `grove/G-260923-gsthp-prepare-grove-for-extern.md` |
| G-111 | `grove/G-111-g-107-docs-plan.md` | G-260923-2zgsr | `grove/G-260923-2zgsr-documentation-inventory.md` |
| G-112 | `grove/G-112-g-107-review.md` | G-260923-f9yah | `grove/G-260923-f9yah-documentation-reconcilia.md` |
| G-113 | `grove/G-113-g-107-router-review.md` | G-260923-d8xkp | `grove/G-260923-d8xkp-router-restructure-revie.md` |
| G-114 | `grove/G-114-capture-and-reuse-terms-question.md` | G-260923-h9c30 | `grove/G-260923-h9c30-capture-and-reuse-terms.md` |
| G-115 | `grove/G-115-g-108-eval-skeleton-plan.md` | G-260923-v9wby | `grove/G-260923-v9wby-eval-skeleton-plan.md` |
| G-116 | `grove/G-116-g-109-attempts-layouts-observed.md` | G-260923-rz01m | `grove/G-260923-rz01m-attempts-layouts-observe.md` |
| G-117 | `grove/G-117-which-attempts-list-and-detail-l.md` | G-260923-hvnqh | `grove/G-260923-hvnqh-which-attempts-list-and.md` |
| G-118 | `grove/G-118-what-mandate-should-the-g-108-pa.md` | G-260923-659zw | `grove/G-260923-659zw-what-mandate-should-the.md` |
| G-119 | `grove/G-119-g-108-eval-skeleton-review.md` | G-260923-50gkk | `grove/G-260923-50gkk-eval-skeleton-review.md` |
| G-120 | `grove/G-120-g-109-attempts-redesign-review.md` | G-260923-sjpfm | `grove/G-260923-sjpfm-attempts-redesign-review.md` |
| G-121 | `grove/G-121-how-do-the-g-108-eval-runs-log-i.md` | G-260924-y99bx | `grove/G-260924-y99bx-how-do-the-eval-runs-log.md` |
| G-122 | `grove/G-122-g-108-baseline-runs-the-missing.md` | G-260924-frzeg | `grove/G-260924-frzeg-baseline-runs-the-missin.md` |
| G-123 | `grove/G-123-make-board-navigation-and-cards.md` | G-260924-nqkkh | `grove/G-260924-nqkkh-make-board-navigation-an.md` |
| G-124 | `grove/G-124-keep-the-board-fresh-without-pre.md` | G-260924-zxvqf | `grove/G-260924-zxvqf-keep-the-board-fresh-wit.md` |
| G-125 | `grove/G-125-answer-a-blocking-question-from.md` | G-260924-wp2pe | `grove/G-260924-wp2pe-answer-a-blocking-questi.md` |
| G-126 | `grove/G-126-g-108-handoff-review-the-login-c.md` | G-260924-ycwx8 | `grove/G-260924-ycwx8-handoff-review-the-login.md` |
| G-127 | `grove/G-127-review-of-g-123-board-navigation.md` | G-260924-tbx21 | `grove/G-260924-tbx21-review-of-board-navigati.md` |
| G-128 | `grove/G-128-g-114-independent-review-of-the.md` | G-260924-w0h5g | `grove/G-260924-w0h5g-independent-review-of-th.md` |
| G-129 | `grove/G-129-headless-attempts-run-long-comma.md` | G-260924-3bapc | `grove/G-260924-3bapc-headless-attempts-run-lo.md` |
| G-130 | `grove/G-130-review-of-g-124-board-re-reads-o.md` | G-260924-380hs | `grove/G-260924-380hs-review-of-board-re-reads.md` |
| G-131 | `grove/G-131-plan-for-g-125-answer-a-blocking.md` | G-260924-cr3e4 | `grove/G-260924-cr3e4-plan-for-answer-a-blocki.md` |
| G-132 | `grove/G-132-g-129-independent-review-of-the.md` | G-260924-95495 | `grove/G-260924-95495-independent-review-of-th.md` |
| G-133 | `grove/G-133-g-125-review-answering-a-questio.md` | G-260924-szryt | `grove/G-260924-szryt-review-answering-a-quest.md` |
| G-134 | `grove/G-134-bound-an-attempt-at-its-plan-and.md` | G-260924-5b6pz | `grove/G-260924-5b6pz-bound-an-attempt-at-its.md` |
| G-135 | `grove/G-135-run-the-g-108-eval-pair-on-codex.md` | G-260924-59f5k | `grove/G-260924-59f5k-run-the-eval-pair-on-cod.md` |
| G-136 | `grove/G-136-g-134-plan-plan-bound-per-phase.md` | G-260924-204rt | `grove/G-260924-204rt-plan-plan-bound-per-phas.md` |
| G-137 | `grove/G-137-g-134-review-plan-bound-per-phas.md` | G-260924-n0747 | `grove/G-260924-n0747-review-plan-bound-per-ph.md` |
| G-138 | `grove/G-138-g-135-codex-eval-row-plan.md` | G-260924-w07wn | `grove/G-260924-w07wn-codex-eval-row-plan.md` |
| G-139 | `grove/G-139-what-mandate-and-login-should-th.md` | G-260924-7x7p7 | `grove/G-260924-7x7p7-what-mandate-and-login-s.md` |
| G-140 | `grove/G-140-default-an-attempt-s-budget-mode.md` | G-260924-ecs9m | `grove/G-260924-ecs9m-default-an-attempt-s-bud.md` |
| G-141 | `grove/G-141-never-run-gpt-6-astra-unless-the.md` | G-260925-04ccr | `grove/G-260925-04ccr-never-run-gpt-6-astra-un.md` |
| G-142 | `grove/G-142-keep-a-blank-mandate-answer-from.md` | G-260925-beby3 | `grove/G-260925-beby3-keep-a-blank-mandate-ans.md` |
| G-143 | `grove/G-143-g-135-codex-eval-row-pattern-pla.md` | G-260925-42j50 | `grove/G-260925-42j50-codex-eval-row-pattern-p.md` |
| G-144 | `grove/G-144-give-adopting-projects-the-recor.md` | G-260925-ced1h | `grove/G-260925-ced1h-give-adopting-projects-t.md` |
| G-145 | `grove/G-145-g-135-plan-cap-guard-and-report.md` | G-260925-6kap4 | `grove/G-260925-6kap4-plan-cap-guard-and-repor.md` |
| G-146 | `grove/G-146-how-should-an-adopting-project-r.md` | G-260925-02jsj | `grove/G-260925-02jsj-how-should-an-adopting-p.md` |
| G-147 | `grove/G-147-g-140-launch-defaults-review-202.md` | G-260925-9ptz2 | `grove/G-260925-9ptz2-launch-defaults-review-2.md` |
| G-148 | `grove/G-148-g-144-plan.md` | G-260925-6ykpv | `grove/G-260925-6ykpv-plan-ship-the-record-mod.md` |
| G-149 | `grove/G-149-g-144-review.md` | G-260925-w62y4 | `grove/G-260925-w62y4-review-record-model-ship.md` |
| G-150 | `grove/G-150-launch-attempts-only-where-the-w.md` | G-260925-3pj9a | `grove/G-260925-3pj9a-launch-attempts-only-whe.md` |
| G-151 | `grove/G-151-strip-grove-repository-pointers.md` | G-260925-ej1xh | `grove/G-260925-ej1xh-strip-grove-repository-p.md` |
| G-152 | `grove/G-152-shipped-document.md` | G-260925-khfe7 | `grove/G-260925-khfe7-shipped-document.md` |
| G-153 | `grove/G-153-search-and-code-links.md` | G-260925-dzxm6 | `grove/G-260925-dzxm6-search-record-bodies-and.md` |
| G-154 | `grove/G-154-listed-constraint-eval.md` | G-260925-pbx81 | `grove/G-260925-pbx81-evaluate-whether-agents.md` |
| G-155 | `grove/G-155-g-154-listed-and-code-constraint.md` | G-260925-a4kn8 | `grove/G-260925-a4kn8-listed-and-code-constrai.md` |
| G-156 | `grove/G-156-g-151-review-shipped-guides-and.md` | G-260925-00h58 | `grove/G-260925-00h58-review-shipped-guides-an.md` |
| G-157 | `grove/G-157-g-150-skill-guard-review.md` | G-260925-1n5x7 | `grove/G-260925-1n5x7-skill-guard-review.md` |
| G-158 | `grove/G-158-what-mandate-should-the-g-154-wi.md` | G-260925-gymkr | `grove/G-260925-gymkr-what-mandate-should-the.md` |
| G-159 | `grove/G-159-g-154-runner-cases-review.md` | G-260925-yygmz | `grove/G-260925-yygmz-runner-cases-review.md` |
| G-160 | `grove/G-160-g-154-without-row-both-constrain.md` | G-260925-khwkq | `grove/G-260925-khwkq-without-row-both-constra.md` |
| G-161 | `grove/G-161-dependency-view.md` | G-260925-g39ga | `grove/G-260925-g39ga-see-work-dependencies-an.md` |
| G-162 | `grove/G-162-bounded-work-selection.md` | G-260925-7c8g9 | `grove/G-260925-7c8g9-execute-an-explicitly-se.md` |
| G-163 | `grove/G-163-selected-work-review-boundary.md` | G-260925-80w3a | `grove/G-260925-80w3a-where-should-review-and.md` |
| G-164 | `grove/G-164-g-153-board-body-search-and-revi.md` | G-260925-nf4hz | `grove/G-260925-nf4hz-board-body-search-and-re.md` |
| G-165 | `grove/G-165-g-161-dependency-view-plan.md` | G-260925-e5qhz | `grove/G-260925-e5qhz-dependency-view-layouts.md` |
| G-166 | `grove/G-166-g-161-dependency-layout.md` | G-260925-t2nb3 | `grove/G-260925-t2nb3-which-dependency-view-la.md` |
| G-167 | `grove/G-167-g-153-board-search-and-review-li.md` | G-260925-yjds8 | `grove/G-260925-yjds8-board-search-and-review.md` |
| G-168 | `grove/G-168-g-161-deps-review.md` | G-260925-be4e3 | `grove/G-260925-be4e3-deps-first-review-gate-o.md` |
| G-169 | `grove/G-169-harness-upgrade-compatibility.md` | G-260925-p2k54 | `grove/G-260925-p2k54-keep-installed-harness-e.md` |
| G-170 | `grove/G-170-release-identity.md` | G-260925-358a2 | `grove/G-260925-358a2-give-every-distributed-b.md` |
| G-173 | `grove/G-173-what-should-g-154-s-with-row-bec.md` | G-260925-9bjrx | `grove/G-260925-9bjrx-what-should-with-row-bec.md` |
| G-174 | `grove/G-174-plan-for-g-170-shared-build-iden.md` | G-260925-sv063 | `grove/G-260925-sv063-plan-for-shared-build-id.md` |
| G-175 | `grove/G-175-g-161-deps-second-review-gate-on.md` | G-260925-h2c5a | `grove/G-260925-h2c5a-deps-second-review-gate.md` |
| G-176 | `grove/G-176-review-of-g-170-release-identity.md` | G-260925-cm2r8 | `grove/G-260925-cm2r8-review-of-release-identi.md` |
| G-177 | `grove/G-177-merge-prediction.md` | G-260925-h8rj5 | `grove/G-260925-h8rj5-predict-whether-a-candid.md` |
| G-178 | `grove/G-178-candidate-target-update.md` | G-260925-dz10z | `grove/G-260925-dz10z-update-a-conflicting-can.md` |
| G-179 | `grove/G-179-standing-policy-question.md` | G-260925-w33j7 | `grove/G-260925-w33j7-what-may-a-standing-owne.md` |
| G-180 | `grove/G-180-policy-driven-integration.md` | G-260925-5wrn8 | `grove/G-260925-5wrn8-resolve-approve-and-inte.md` |
| G-181 | `grove/G-181-g-154-final-review-runner-fixes.md` | G-260925-j4j17 | `grove/G-260925-j4j17-final-review-runner-fixe.md` |
| G-182 | `grove/G-182-standing-policy-delegation.md` | G-260925-wh9ax | `grove/G-260925-wh9ax-delegate-conflict-resolu.md` |
| G-183 | `grove/G-183-plan-for-g-177-merge-prediction.md` | G-260925-gzdsd | `grove/G-260925-gzdsd-plan-for-merge-predictio.md` |
| G-184 | `grove/G-184-g-169-plan-entrypoint-revisions.md` | G-260925-zx4x0 | `grove/G-260925-zx4x0-plan-entrypoint-revision.md` |
| G-185 | `grove/G-185-g-162-selected-work-plan.md` | G-260925-t70h8 | `grove/G-260925-t70h8-selected-work-one-attemp.md` |
| G-186 | `grove/G-186-entrypoint-revision.md` | G-260925-m9jcr | `grove/G-260925-m9jcr-entrypoint-revision.md` |
| G-187 | `grove/G-187-review-of-g-169-entrypoint-revis.md` | G-260925-jd94s | `grove/G-260925-jd94s-review-of-entrypoint-rev.md` |
| G-188 | `grove/G-188-selected-work-shared-candidate.md` | G-260925-wc2pz | `grove/G-260925-wc2pz-review-an-explicitly-sel.md` |
| G-189 | `grove/G-189-review-of-g-177-candidate-merge.md` | G-260925-a05hb | `grove/G-260925-a05hb-review-of-candidate-merg.md` |
| G-190 | `grove/G-190-g-162-review-1.md` | G-260925-2v889 | `grove/G-260925-2v889-final-review-of-in-three.md` |
| G-191 | `grove/G-191-plan-for-g-178-candidate-resolution.md` | G-260926-bp82x | `grove/G-260926-bp82x-plan-for-resolve-a-confl.md` |
| G-192 | `grove/G-192-resolution.md` | G-260926-kfcpp | `grove/G-260926-kfcpp-resolution.md` |
| G-193 | `grove/G-193-review-of-g-178-candidate-resolution.md` | G-260926-wmyet | `grove/G-260926-wmyet-review-of-candidate-reso.md` |
| G-194 | `grove/G-194-identify-records-by-creation-dat.md` | G-260926-2da4n | `grove/G-260926-2da4n-identify-records-by-crea.md` |
| G-195 | `grove/G-195-coordination-free-record-ids.md` | G-260926-pgj43 | `grove/G-260926-pgj43-coordination-free-record.md` |
| G-196 | `grove/G-196-plan-for-g-180-policy-integration.md` | G-260926-vpvhf | `grove/G-260926-vpvhf-plan-for-policy-driven-i.md` |
| G-197 | `grove/G-197-policy.md` | G-260926-a8vyj | `grove/G-260926-a8vyj-policy.md` |
| G-198 | `grove/G-198-review-of-g-180-policy-integration.md` | G-260926-mg6g8 | `grove/G-260926-mg6g8-review-of-policy-driven.md` |
| G-199 | `grove/G-199-coordination-free-record-ids.md` | G-260926-bz4n5 | `grove/G-260926-bz4n5-coordination-free-record.md` |

One record already in the date form was named after a legacy ID and took the
same kind of name: G-260926-afe5w, `grove/G-260926-afe5w-review-of-g-195-coordina.md`, is now `grove/G-260926-afe5w-review-of-coordination-f.md`.
