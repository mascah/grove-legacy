---
id: "G-260921-gtydy"
type: decision
title: "Keep identity and placement stable while knowledge evolves"
status: accepted
created: "2026-09-21T15:40:59Z"
updated: "2026-09-26T16:11:44Z"
relates_to: ["G-260919-6mpmw", "G-260919-4h6pn", "G-260921-e8bva", "G-260921-w9x25", "G-260921-ebsby", "G-260921-r491p", "G-260921-407n6", "G-260926-2da4n", "G-260926-yvjy6"]
formerly: "D-006"
---

## Decision and authority

On 2026-09-21, after G-260921-w9x25 was integrated, the owner reopened both mandatory
type folders and the closed set of knowledge types before assigning G-260921-r491p.
They selected folder placement as a convention, leaned toward a flat default
and type-independent identity, and questioned the value of moving completed
work. They then endorsed "stable identity, stable placement, flexible content,
and multiple ways to browse it" and requested a durable decision and realigned
work so execution could continue.

In the same session, the owner clarified that multiple existing layouts were
already causing grief and explicitly selected "Reconcile both IDs and file
locations" for a one-time repository migration. This replaces the earlier
recommendation to leave existing records in their old folders with typed IDs.

The [brief](brief.md) owns the selected direction:

- One configured root, flat creation by default, recursive discovery. Folders
  do not determine record type, lifecycle or validity. Nested records remain
  supported, but this repo adopts one flat layout through the migration below.
  No directory-routing configuration in the initial change.
- One neutral sequential ID namespace, independent of type. After the explicit
  migration, ordinary edits preserve identity and references. Title,
  classification and status changes do not rename or move a file. Keep a short
  descriptive filename slug.
- General knowledge pages need a supported minimal representation without
  choosing from a predefined taxonomy or acquiring a ticket lifecycle. Known
  operational types retain explicit validation for facts software acts on;
  flexibility does not make prose into acceptance or merge authority.
- Completed records stay where they are. Bound ordinary views and retain
  explicit lookup; no automatic filing, completed-work folder or cleanup command
  is selected. No per-type prefix settings or general schema-extension engine.

This revises G-260919-6mpmw/G-260921-e8bva's required type folders and G-260919-4h6pn/G-260921-e8bva's type-bound
identity for future creation. It preserves sequential allocation across local
worktrees, revision-checked writes, G-260921-w9x25's useful types and artifact links,
the one brief, and G-260921-9wkjt's ownership of review/approval semantics. G-260921-w9x25 remains
completed evidence of schema 2, not work to reopen or undo.

## Delivery boundary

[G-260921-ebsby](G-260921-ebsby-decouple-record-identity.md) owns compatible schema/CLI support.
[G-260921-r491p](G-260921-r491p-reconcile-all-grove-cont.md) follows with a complete one-time
reconciliation of this repo's Grove content: all existing records and legacy
plans/reviews receive neutral IDs and flat filenames beneath the record root;
the brief moves there as the separately configured `brief.md`. Preserve content,
relationships and provenance, repair current references, and retain one durable
old-ID/path to new-ID/path mapping. Remove superseded files and type folders.

This is an explicitly authorized migration scope, not automatic behavior on
ordinary edits or completion, and not an assignment to execute G-260921-r491p now. Do
not rewrite historical commits or fabricate new evidence for old reviews.
Old-schema reading remains necessary for historical branches; it does not
require multiple live conventions after the migration. Repository entrypoints
and product/workflow documentation retain their functional homes. Sibling
repositories still need their own explicitly assigned migration.

The selected direction is not the current file format. Until support ships,
use [the implemented record model](../docs/record-model.md), including its
typed creation commands. G-260921-ebsby's preparation settles exact schema version,
minimal fields, generic-page authoring and allocation compatibility; suggested
new ID spelling is `G-NNN`. These mechanical details do not require another
layout-design session unless they cannot meet the selected bounds.

## Evidence and alternatives

At main `76da081`, `internal/project/project.go` requires a type folder and
`internal/project/metadata.go` requires matching type and ID prefix. A disposable
clone passed validation with W-019 beneath `work/knowledge/`, but failed with
the identical record at the root or an ordinary synthesis page there. Current
filenames are already conventions, and ID relationships survive path changes;
ordinary Markdown links still require repair when paths change.

- Required type folders aid scanning, but turn organization into schema.
- A flat directory alone leaves the closed knowledge taxonomy untouched.
- Configurable per-type folders/prefixes add migration and historical-config
  questions before demonstrated need. Read existing records independently of
  creation defaults if configuration is added later.
- Completion-driven moves reduce file-tree clutter at the cost of link repairs
  and branch churn; bounded views solve the immediate board need.

Research inspected on 2026-09-21:
[Karpathy's llm-wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
leaves layout domain-specific and emphasizes maintained synthesis and links.
[Pocock's domain modeling](https://github.com/mattpocock/skills/blob/main/skills/engineering/domain-modeling/SKILL.md)
separates vocabulary and decisions.
[Backlog.md's CLI](https://github.com/MrLesk/Backlog.md/blob/main/CLI-INSTRUCTIONS.md)
supports document moves and a separate completed-task cleanup;
its [document update](https://github.com/MrLesk/Backlog.md/blob/main/src/core/backlog.ts)
preserves identity across type/path changes. These are precedents, not adopted
workflow requirements. None establishes a measured scalability limit for Grove.

## Reconsideration

Revisit organization or creation settings when real browsing/adoption evidence
shows a need. Add operational semantics when software needs them, not merely
because a new topic deserves a page. Future changes must preserve existing
identity and discovery and explicitly handle Markdown links. A configurable
root is not permission for silent root relocation or bulk rewriting.

[G-260926-2da4n](G-260926-2da4n-identify-records-by-crea.md) revised the "neutral
sequential ID namespace" on 2026-09-25 to one neutral namespace without a
sequence; stable identity, placement and short slugs stand.

[G-260926-yvjy6](G-260926-yvjy6-retire-legacy-ids-by-ren.md) authorized a second one-time migration on 2026-09-26, legacy IDs to
the date form; stable identity and placement stand for ordinary edits.
