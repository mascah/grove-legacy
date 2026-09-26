---
id: "G-260925-khfe7"
type: term
title: "Shipped document"
status: proposed
created: "2026-09-25T19:06:02Z"
updated: "2026-09-25T19:06:02Z"
---

## Meaning

A document the `grove` binary carries and prints or writes wherever it runs:
the work, shaping and review guides (`grove guide work|shape|review`) and
the record model (`grove guide model`). It is read in projects that hold none of Grove's files and
whose own `G-` IDs are live, so it must read the same, and be true, in every
one of them: it links only within itself or to `https://`, names another
shipped document by the command that prints it, and names no Grove record,
path, invocation or history. This repository's file is the one editable
owner; the binary's copy is what an adopting project reads, and `grove
version` names that copy: its guides digest covers the three guides and the
model, and its content digest covers every shipped document together with
the adapters `init` generates ([G-260925-358a2](G-260925-358a2-give-every-distributed-b.md)).

Not shipped: the brief, the command reference, the board guide, AGENTS.md,
this repository's records, and its own `.claude/skills/`,
`.agents/skills/` and `.claude/agents/` adapters, which read the documents
as files. The adapters `init` writes, the `grove-reviewer` agent definition
among them, are generated text under the same rule, not documents: they
load a shipped guide by command with their entrypoint revision
([G-260925-p2k54](G-260925-p2k54-keep-installed-harness-e.md)).

## Relationships

The rule lives in AGENTS.md's constraints; its reason is
[G-260925-02jsj](G-260925-02jsj-how-should-an-adopting-p.md)'s answer, which
[G-260925-ej1xh](G-260925-ej1xh-strip-grove-repository-p.md) extends from links to
mentions. `TestGuideAndVersionNeedNoProject` enforces it. Delivered
by [G-260921-5gz9a](G-260921-5gz9a-bootstrap-projects-with.md) (the guides, `init`),
[G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) (the reviewer definition,
whose body [G-260925-p2k54](G-260925-p2k54-keep-installed-harness-e.md) moved into the
review guide)
and [G-260925-ced1h](G-260925-ced1h-give-adopting-projects-t.md) (the model).
