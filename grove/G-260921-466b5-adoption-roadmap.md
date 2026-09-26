---
id: "G-260921-466b5"
type: plan
title: "G-260921-407n6 adoption roadmap"
status: current
formerly: "docs/plans/W-018-adoption-roadmap.md"
work: ["G-260921-407n6"]
created: "2026-09-21T01:08:42Z"
updated: "2026-09-21T21:11:16Z"
---

# G-260921-407n6 adoption roadmap

Selected 2026-09-20; storage/identity sequence revised 2026-09-21. The
[brief](brief.md) owns direction and
[G-260921-tkdwh](G-260921-tkdwh-adopt-the-interactive-ad.md) and
[G-260921-gtydy](G-260921-gtydy-keep-identity-and-placem.md) record authority.
[G-260921-407n6](G-260921-407n6-complete-the-interactive.md) owns the first milestone's
acceptance. This is its coordination plan, not a batch assignment or a detailed
implementation plan. Records below own their scope, acceptance and next action.

## How to continue without this conversation

1. Start with the chosen work record, repository instructions and the shared
   execution guide through `/grove-work W-ID` (Claude) or `$grove-work W-ID`
   (Codex). In this repository every command is `go run ./cmd/grove …`.
2. Inspect its current context, branch/worktree, prerequisites and linked current
   plan. Prepare a concise linked plan record through the supported CLI when
   needed. Legacy plans linked here remain until G-260921-r491p migrates them; do not
   create new files under `docs/plans/`.
3. Read only the relevant research section when its uncertainty matters. This
   roadmap and historical reviews are not mandatory context for every task.
4. Complete the selected assignment and its acceptance, preserve evidence and
   next action, and obtain its explicit integration disposition. Do not expand
   into all roadmap entries automatically.

Work remains Proposed until individually assigned and started. Approval of this
roadmap does not assign sibling migrations, background execution or automatic
merges; those actions belong to their explicitly selected work and mandate. Existing assignment authority still applies within a unit;
routine technical decisions do not require repeated permission.

## First adoption milestone

Preferred order is top to bottom. `depends_on` in the records expresses actual
prerequisites; order here also reflects product priorities.

| Work | Deliverable | Completion evidence |
| --- | --- | --- |
| [G-260919-04z88](G-260919-04z88-shape-project-work-throu.md) | Native interactive shaping and shared authoring guidance | Real requirements session and fresh harness discovery; no fabricated knowledge records |
| [G-260921-w9x25](G-260921-w9x25-represent-domain-terms-a.md) | Domain terms and discoverable linked plans/reports/reviews | Compatible schema and staged retrieval, preserved identity/links |
| [G-260921-ebsby](G-260921-ebsby-decouple-record-identity.md) | Stable identity/placement with flexible knowledge and flat creation | General pages, neutral new IDs, recursive discovery and explicit old-schema/allocator compatibility |
| [G-260921-r491p](G-260921-r491p-reconcile-all-grove-cont.md) | Complete Grove-content reconciliation into one flat layout and neutral IDs | All records and legacy artifacts mapped, provenance retained, references repaired and no parallel old layout |
| [G-260921-9wkjt](G-260921-9wkjt-hand-implementation-cand.md) | Review lifecycle and durable candidate handoff | Revision-bound evidence, explicit completion migration, manual review/integration path |
| [G-260921-9t178](G-260921-9t178-prove-the-complete-inter.md) | Complete loop on real Grove work | Fresh-session continuation, independent review, owner verdict and integration evidence |
| [G-260921-5gz9a](G-260921-5gz9a-bootstrap-projects-with.md) | Minimal setup and versioned portable workflows | Disposable project adoption and observed Claude/Codex entrypoints |
| [G-260921-905y3](G-260921-905y3-cut-nullsec-over-to-this.md) | Nullsec cut over to this Grove, predecessor uninstalled | Rehearsal, live migration, fresh-session discovery and the owner's judgment of the converted tree; the first real change and the verdict stay with G-260921-407n6 (owner, 2026-09-22) |

G-260921-ebsby extends G-260921-w9x25 without reopening its completion. G-260921-r491p needs that support;
G-260921-9wkjt also needs G-260921-ebsby's contract but does not technically depend on moving this
repo's content. The preferred sequence performs the full reconciliation first.
G-260921-9t178 checks the connected shaping/knowledge/review result before
G-260921-5gz9a carries it into another repository. G-260921-905y3 cuts nullsec over, and the
first real nullsec change completes G-260921-407n6; member status
alone cannot establish the milestone's human acceptance.

For the initial loop, plans stay proportional and human review can read CLI/files.
Design board/detail/review together early, but do not make the full TUI redesign
a prerequisite for the real hobby-project change.

## Following investments

| Work | Deliverable | Boundary |
| --- | --- | --- |
| [G-260921-ms6ev](G-260921-ms6ev-derive-a-project-wide-cu.md) | Project-wide current projection | Resolve concrete ambiguous histories; preserve explicit sources and batched reads |
| [G-260921-k0mwk](G-260921-k0mwk-make-the-board-and-item.md) | Polished board, detail, artifacts, timeline and list/search | Owner-reviewed visuals and connected terminal checks |
| [G-260921-jwk4e](G-260921-jwk4e-review-candidates-and-in.md) | Progressive review plus local approval/integration actions | Candidate-bound approval, conflict/refusal and safe cleanup |
| [G-260921-h46pb](G-260921-h46pb-run-one-bounded-implemen.md) | One bounded durable implementation attempt | Fake-process failure probes, then a bounded actual-provider trial |
| [G-260921-7trd7](G-260921-7trd7-launch-and-inspect-manag.md) | TUI launch, runs overview, reconnect, stop and feedback continuation | UI observes the independent owner and reuses the same review contract |

G-260921-ms6ev depends on the knowledge/lifecycle contracts, not on the act of migrating
nullsec. G-260921-h46pb needs portable assignment and review contracts, not the TUI merge
screen; its later place is investment order. These distinctions allow deliberate
reordering without manufacturing dependencies or silently widening an assignment.

## Preparation questions, at their point of use

These are design tasks, not unanswered human preferences that block assignment.
Investigate routine choices; create a question through the CLI only when an
actual consequential human decision remains, linked to the work it blocks.

| Boundary | Questions to resolve | Owner |
| --- | --- | --- |
| Knowledge foundation | Minimal general-page envelope and authoring; neutral allocator compatibility; explicit schema migration; keep operational validation | G-260921-ebsby |
| Repository reconciliation | Full record/document inventory, old-ID/path mapping, references/evidence, rehearsal/recovery and old-branch reintegration | G-260921-r491p |
| Review lifecycle | Candidate/input identity; disposition and integration receipts; historical Done migration; research/design completion | G-260921-9wkjt |
| Portability | Workflow packaging, managed adapter updates, executable coexistence and minimal initialization | G-260921-5gz9a |
| Live adoption | Current nullsec scope, retained/archive mapping, collision handling, rollback and explicit cutover | G-260921-905y3 |
| Current projection | Target branch, record-level ancestry, dirty overlays, deletion/reverts, divergent-card placement | G-260921-ms6ev |
| Visual experience | Concrete layouts, narrow terminals, bounded Done defaults, timeline navigation, compatible Charm modules | G-260921-k0mwk |
| Local integration | Supported merge strategy, moved target, stale approval, status publication and cleanup failures | G-260921-jwk4e |
| Runtime | Native Claude background versus owned process; role/model controls, budgets, event protocol and owner-loss recovery | G-260921-h46pb |
| Run experience | Separate run screen versus detail section, activity summarization, feedback relaunch | G-260921-7trd7 |

## Research and source map

The [direction evaluation](G-260921-72chf-first-days-evaluation-an.md) holds
the repository observations and external references behind this sequence.
Read the relevant subsection at preparation, not every predecessor skill.

[G-260921-gtydy](G-260921-gtydy-keep-identity-and-placem.md) records the later layout
investigation, owner decision and research links. Its selected bounds replace
mandatory folder/type/prefix coupling; the current record model stays in force
until G-260921-ebsby ships.

- Shaping, terms and context: evaluation's **Knowledge and context**;
  prior [shaping/runner research](G-260919-rr3ae-shaping-and-headless-run.md).
- Preparation, review and roles: evaluation's **Workflow and roles**;
  prior [predecessor work review](G-260919-ph0w1-predecessor-work-review.md).
- Projection and TUI: evaluation's **Current state and presentation**, G-260919-8jb5s's
  historical answer, and the current versions/history implementation.
- Process ownership: evaluation's **Future runtime boundary**, then current
  provider documentation and installed capabilities. Old adapter code is
  evidence, not a supported implementation to copy unexamined.

## Later candidates, not implementation-ready assignments

- **Learning/debrief:** mine linked plans, reports and reviews for repeated
  failures or missing context; retain supporting examples and propose bounded
  changes to knowledge, guidance or checks. Collect useful evidence now; measure
  whether a change improves subsequent work before making universal policy.
- **Multiple work items:** explicit bounded selection, dependency/ownership
  analysis, sequential execution by default, per-item acceptance and shared
  integration evidence. Define budgets, maximum fan-out, stop conditions and
  partial completion before enabling unattended batches.
- **Research/dream loop and scheduling:** one research mandate, no automatic
  expansion into implementation; separate reviewable branch, evidence-backed
  proposals, durable human questions, deduplication and unchanged-wait handling.
  Scheduling is another caller of the proven run contract.
- **Dependency graph:** visualize actual dependency edges and distinguish them
  from membership or preferred order. Use real graph size/readability evidence
  before choosing a terminal rendering approach.
- **Additional operational semantics:** designs, flows, systems, research,
  releases, builds and team material can start as general knowledge. Add a
  built-in contract only when software needs distinct behavior; capturing a
  page does not require expanding the schema or recreating Bench's taxonomy.
- **Keyborg adoption:** a second test of packaging and migration with its pinned
  Bench contract; shape it from the nullsec results and current project needs.
- **Autonomous judging, other providers and recovery after reboot:** separate
  policy/capability decisions. Interactive Claude/Codex support does not establish
  managed-provider parity or unattended acceptance authority.

## Completion and reconsideration

Use G-260921-9t178 and G-260921-905y3 evidence to adjust this order and the owning records.
Changes to product direction belong in the brief with an attributable decision.
Current progress belongs in work records, not another status table here. When
the first milestone completes, record the owner's verdict in G-260921-407n6 and select
the next bounded investment from the observed friction.
