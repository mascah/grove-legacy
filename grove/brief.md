# Grove brief

Direction reconciled 2026-09-23 ([G-260923-fwakw](G-260923-fwakw-reconcile-current-docume.md)).
This is the single current source of product intent. The owner selected the
interactive adoption milestone in [G-260921-tkdwh](G-260921-tkdwh-adopt-the-interactive-ad.md),
the storage and identity direction in [G-260921-gtydy](G-260921-gtydy-keep-identity-and-placem.md),
and the next phase's audience on 2026-09-23 (below). Work records own
acceptance, progress and next actions; the
[record model](../docs/record-model.md) describes the implemented schema.

## Purpose

Grove is a local project workspace where a person shapes intent with agents,
delegates bounded work, and returns to a clear review of the result. It connects
project knowledge, work, and execution evidence so the next human or agent can
continue with **the right context, at the right time**.

The first audience is the owner working across local hobby repositories and
worktrees. The first adoption milestone, a complete interactive
shape → implement → review → integrate loop,
[G-260921-407n6](G-260921-407n6-complete-the-interactive.md), was closed by the owner on 2026-09-22.
Its nullsec acceptance was not exercised in nullsec: the owner accepted this
repository's own use of the loop in its place. Nullsec's cutover to this
Grove ([G-260921-905y3](G-260921-905y3-cut-nullsec-over-to-this.md)) was part of that milestone. Keyborg, selected in G-260921-tkdwh
as the second adoption test, has no work record yet.

The next audience, selected by the owner in a shaping conversation on
2026-09-23, is a small external preview: initial use outside the owner's
projects. "Preview" names that audience only. It is not a selected release
channel, launch commitment, support policy or compatibility promise; release
scope, platforms and licensing remain later choices.

## Selected foundations

- Go CLI, ordinary Markdown with YAML frontmatter, `grove.yaml`, configurable
  record storage, and Git. Core inspection needs no running service.
- One neutral `G-` namespace of date-form IDs issued without coordination;
  revision-checked mutations; ordinary edits preserve identity.
- Humans and agents use the same records. Deterministic validation, retrieval,
  mutation and lifecycle mechanics belong in software; judgment belongs in
  instructions and attributable human or delegated decisions.
- Branch-local editing with a combined project view. Actions bind to an exact
  checkout and revision; browsing should not make source selection the main job.
- Bare `grove` opens the TUI; explicit commands remain noninteractive. The TUI
  should be visually appealing, keyboard-friendly, and use the Charm ecosystem.
- Grove-owned shared workflows with thin Claude and Codex adapters, carried
  inside the binary. They replaced the sibling skills in this repository and
  in nullsec.

## Knowledge, work, and evidence

Keep their lifetimes and owners distinct:

| Information | Owner and purpose |
| --- | --- |
| Brief | Compact project purpose, constraints, selected direction and next investment |
| Terms | Domain vocabulary and relationships; no implementation diary |
| Decisions | Consequential choices, authority, alternatives and reconsideration conditions |
| Questions | Unresolved matters, affected work, evidence, and who can answer |
| Work | Outcome, bounds, acceptance, dependencies and next action |
| Artifacts | Plans, implementation reports and reviews linked to work and relevant revisions |
| Attempts | Particular executions, inputs, owner, progress, result and recovery state |

Retain terms and discoverable artifacts. Plans need not become an
independent ticket lifecycle. Small work may carry its preparation in the work
record; substantial work gets a concise linked plan. Designs, flows, systems,
research, releases, builds and team definitions can begin as general knowledge;
new operational types are introduced when software needs distinct behavior.

**Stable identity, stable placement, flexible content, multiple views.** One
configured Grove root has a flat creation default and recursive discovery;
folders do not determine validity or meaning. General knowledge needs no
predefined semantic category or ticket lifecycle. Known operational records
retain explicit contracts for the facts software acts on. Classification,
title and status changes preserve identity and path. Completed records stay
put; views bound everyday clutter. Per-type folder/prefix settings, automatic
filing and a general schema-extension engine are not selected. This repository
had one deliberate reconciliation into that layout
([G-260921-r491p](G-260921-r491p-reconcile-all-grove-cont.md), mapped by
[G-260921-czt8x](G-260921-czt8x-identity-and-path-migrat.md)); stable placement applies from then on.
Until a first release Grove keeps no backward compatibility: only the current
schema is read, and an old commit is inspected with the CLI it carries.

Retrieve context by activity: shaping starts with the brief and relevant
knowledge; preparation with selected work and affected interfaces; execution
with its mandate and current task; review with acceptance and the candidate;
resume with a checkpoint and changed inputs. Links enable discovery without
preloading everything. A context response is facts, never authorization.

## Workflow and authority

Work lifecycle: **Proposed → Active → Review → Done**, with Abandoned
available only through an explicit human decision for now. Preparation,
implementation and agent review are activities; waiting, failures and process
state are additional facts. A terminal attempt does not automatically enter Review.

For implementation, Done means accepted and integrated into the configured
target ([G-260921-9wkjt](G-260921-9wkjt-hand-implementation-cand.md)); a `done` record without a
candidate predates that meaning and keeps its historical evidence.
Research/design work needs a completion condition suitable to its deliverable.

An assignment or Implement action authorizes bounded execution; merely creating
a proposal or editing status does not. Preparation investigates technical
unknowns and surfaces consequential human decisions before implementation.
Routine technical choices proceed within the mandate. Separate plan-file or
stage-by-stage human approval is not a universal gate.

Independent agent review examines a candidate against acceptance and evidence.
Human review presents changed behavior, decisions, outstanding issues and checks
before optional diffs. Approval is tied to the candidate being integrated;
changed candidates require reconsideration. Feedback that starts another
implementation attempt returns work to Active and preserves prior reviews.
A standing policy in `grove.yaml` may delegate one bounded conflict
resolution, and the approval and integration of a candidate that meets its
written conditions after independent review, which may include a bounded
delegated judgment of the candidate against its record
([G-260928-d8py6](G-260928-d8py6-a-standing-policy-may-de.md), owner,
2026-09-28), each act attributed to the
policy and its evidence ([G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md),
owner, 2026-09-25); everything it does not name awaits human judgment.

Interactive and headless callers share the workflow, with explicit human
availability, authority and resource bounds. Missing human decisions become
durable questions; unchanged waits do not trigger repeated work. Unattended
research publishes proposals on an isolated reviewable branch, with supporting
knowledge identified for selective integration, and cannot authorize its own
implementation or merge.

## Current view and TUI

The normal view is project-wide and independent of the invoking checkout.
Present current work first, collapse identical observations, place demonstrably
superseded states in history, and label unintegrated and uncommitted changes.
Genuine divergence stays visible. Git ancestry supports this view; timestamps,
status rankings, or the newest branch tip do not define authority.

This supersedes G-260919-8jb5s's explicit-versions-first presentation as the default.
Exact source inspection and fresh workspace binding remain available.
Ancestry here means merge bases: a copy is superseded when, since it and
another split, only the other changed the record. The integration target is
optional configuration (`target: main` in `grove.yaml`), chosen by the owner
on 2026-09-22. It labels progress as not yet on the target and never decides
which state is current.

Board columns follow the lifecycle, with bounded recent Done items,
searchable older work, and Abandoned hidden by default. Detail leads with
rendered work Markdown, with artifacts, workflow indicators and a timeline;
source observations are secondary. Review leads with a concise result and
progressively exposes evidence, files and diffs. Design these views together;
visual quality requires the owner's judgment in an actual terminal.

## Adoption and execution

`grove init` bootstraps another repository: valid configuration, a
placeholder brief, the record root and managed entrypoints, preserving existing
instructions. A shaping session develops the brief. Preserve provenance and
repair current references; do not retain parallel editable layouts. Product
and workflow instructions keep their functional homes. Live changes to a
sibling repository need a separately assigned scope, and an explicit cutover
preserves identity, knowledge and evidence with one editable authority for
each fact.

One bounded implementation runs independently of the TUI as a Grove-owned
provider process, Claude Code today and Codex once
[G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) lands, that
continues after the terminal closes, reconnects without duplication and
stops explicitly ([G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)
records the process choice and when to reconsider it;
[G-260928-vdhf0](G-260928-vdhf0-run-attempts-on-codex-as.md) the second
provider, the owner on 2026-09-28). No machine-reboot guarantee is selected. Roles start with
shaping/preparation, implementation and independent review; route model
strength by uncertainty and consequence and retain actual configuration per
attempt. Bound batches and fan-out before adding schedules.

Reports and reviews should support evidence-backed process learning. A bounded
debrief proposes changes to knowledge, guidance or mechanical checks; it does
not silently rewrite acceptance or turn every observation into permanent policy.

## Observed state

At main `768efab` (2026-09-23), the Go CLI inspects, creates, updates and
converts records; shows versions across branches and worktrees and resolves
their checkouts; assembles staged context; bootstraps a repository with
`init` and prints its embedded guides; records approval and feedback and
integrates an approved candidate; and runs, lists and stops headless attempts.
Bare `grove` opens the board on the current view, with record detail, search,
the review actions and attempts. Terms, plans, reviews and pages are records.
This is a dated observation, not a progress log: work records own what has
changed since.

The predecessor is uninstalled.
The archived FastAPI/PostgreSQL application
(the sibling checkout `grove-archive-2026-09-18/`, last commit `be40e46`) is
historical, with no service,
credential, deployment or backlog authority over this project.

## Suggested sequence

The brief states no order among work records. That order is their
`depends_on` edges, which `grove deps` and the board's `g` show
([G-260927-y3pc2](G-260927-y3pc2-order-is-an-edge.md)). Change this section
only to select direction for outcomes not yet shaped as work.

The preview phase is three pieces of work: behavioral evaluations of context
and workflows ([G-260923-p5pt6](G-260923-p5pt6-establish-behavioral-eva.md)),
usable Attempts screens ([G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md)),
and distribution preparation ([G-260923-gsthp](G-260923-gsthp-prepare-grove-for-extern.md)).
Completing them is not by itself evidence that a preview is ready; the owner
judges that.

The brief does not track progress or a current next action. Each work
record's Next owns that. The adoption [roadmap](G-260921-466b5-adoption-roadmap.md) and
[evaluation](G-260921-72chf-first-days-evaluation-an.md) are the closed milestone's
history.
