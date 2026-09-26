---
id: "G-260925-m9jcr"
type: term
title: "Entrypoint revision"
status: proposed
created: "2026-09-25T23:19:59Z"
updated: "2026-09-25T23:29:12Z"
relates_to: ["G-260925-p2k54", "G-260921-vz0v3", "G-260925-khfe7", "G-260925-358a2"]
---

## Meaning

An integer naming what an entrypoint `grove init` writes (the `grove-work`
and `grove-shape` skills, their Codex policies, the `grove-reviewer` agent
definition) needs of the `grove` it calls: which guides it loads and how.
Each such file states it as `grove entrypoint revision N` and passes it as
`grove guide NAME --entrypoint N`. A binary serves a range of revisions:
`guide` refuses an `--entrypoint` outside it, `init --check` reports a
marked file outside it and exits 1, and `run` refuses to launch with one. A
file with the managed marker and no revision line predates revisions and is
revision 1, legacy.

Not a [revision](G-260921-vz0v3-revision.md): that is one file's `sha256:` content
identity, and two files of one entrypoint revision may differ in every
byte. Not the release version `grove version` prints, nor the record
schema's `schema_version`: a release may keep the entrypoint revision, and
a revision changes only when what entrypoints need of the binary changes,
never for new wording.

## Relationships

Introduced by [G-260925-p2k54](G-260925-p2k54-keep-installed-harness-e.md), whose plan
[G-260925-zx4x0](G-260925-zx4x0-plan-entrypoint-revision.md) records the design; the
command reference's Entrypoint revisions section owns the behaviour.
[G-260925-358a2](G-260925-358a2-give-every-distributed-b.md) keeps release version, schema and
entrypoint compatibility distinct. The files it names are generated text
under the [shipped document](G-260925-khfe7-shipped-document.md) rule.
