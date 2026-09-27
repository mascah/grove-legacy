---
id: "G-260927-ngkbz"
type: work
title: "Frame the README for someone adopting Grove in their project"
status: proposed
created: "2026-09-27T17:10:58Z"
updated: "2026-09-27T17:11:30Z"
relates_to: ["G-260923-gsthp", "G-260921-5gz9a"]
size: small
---

## Outcome

Someone setting Grove up in their own project can learn, from the README and
from what `init` prints, how to install it, what `grove.yaml` takes, and how
work moves from shaping to done, without this checkout. The owner first:
adopting Grove in ascah.dev on 2026-09-27, they ran `grove init`, found
nothing saying what `grove.yaml` values to set, and fell back to copying this
repository's. Other adopters have no such fallback.

Owner intent, shaping conversation 2026-09-27: the README is written for
someone working on Grove itself and treats adoption as the side case ("Use
Grove in another repository"); give it a full pass with the adopter as its
reader, as one work item that also covers `init`'s pointer, and show
`grove.yaml` in the README minimally. The owner accepted that this serves the
brief's preview audience even though
[G-260923-gsthp](G-260923-gsthp-prepare-grove-for-extern.md) narrowed
distribution to the bare minimum.

## Constraints

Observed at `main` `67d35da`, 2026-09-27:

- [README.md](../README.md) opens with "Run it" (`go run ./cmd/grove`,
  "From this repository", `show G-260919-rt9h9`), then the command table and
  the board, and only then "Use Grove in another repository". Its intro
  describes mechanism (random-tail IDs, Git ancestry) rather than what a
  person does with Grove. "Where each subject lives" mixes adopter subjects
  with rows only this repository needs (dogfooding evidence records, evals,
  legacy ID mapping), and [CLAUDE.md](../CLAUDE.md) already keeps its own
  list of owners for agents working here.
- The `grove.yaml` keys are documented only in
  [the record model](../docs/record-model.md#configuration-and-discovery)
  (`grove guide model`, "Configuration and discovery", within a 592-line
  document): required `schema_version` and `records`; optional `brief`,
  `target`, `run:` and `policy:`. `init` writes only `schema_version: 3`,
  `records: grove` and `brief: grove/brief.md`, which is what ascah.dev's
  `grove.yaml` holds.
- `init`'s closing note on stderr (`internal/cli/init.go:97`) says to check,
  commit, and how the entrypoints invoke `grove`; it names no configuration
  reference. `internal/cli/init_test.go:92` asserts only its opening words.
  [Init](../docs/commands.md#init) describes the note.
- CLAUDE.md says the README owns "what Grove is and where each subject is
  owned"; this work changes that contract, so the line is reconciled with it.
- Nothing reads README.md programmatically (repository search, 2026-09-27).

Proposed design (binds nobody), agreed in shape by the owner:

1. What Grove is: short, in terms of outcomes.
2. Install: `go install`, `grove version`.
3. Set up a project: `init`, what it writes, the minimal `grove.yaml` with
   `target: main`, one line each on `target`, `run:` and `policy:` (leave it
   out until you want automatic acts), a pointer to `grove guide model` for
   the rules, commit.
4. The loop: brief, `/grove-shape`, `/grove-work` or `grove run`, judge on the
   board (`a`, `f`, `i`), done.
5. Commands: the table, each subject's `grove guide …` command beside its
   docs link where a shipped guide covers it.
6. Developing Grove, last: `go run` from a checkout, CLAUDE.md, and one line
   that `grove/` here is a real project tracked with Grove.

Out of scope: making the configuration easier to find inside
`grove guide model` (a section selector or a separate config guide). The
owner named that as awkward and chose not to solve it now. Also out: the
content of `docs/`, a license, a release pipeline or site.

## Acceptance

- The README reads in the order above, and no contributor-only instruction
  (`go run ./cmd/grove`, a record ID of this repository) appears before the
  section on developing Grove.
- The `grove.yaml` the README shows passes `grove check` in a disposable
  project after `init`, and each optional key it names matches the record
  model.
- The README keeps no row or link that only this repository's development
  needs outside its development section; every link resolves.
- `init`'s closing note names `grove guide model` as where `grove.yaml`'s
  keys are described, a test asserts it, and [Init](../docs/commands.md#init)
  agrees.
- CLAUDE.md's description of what the README owns matches the new README.
- CLAUDE.md's final checks pass (`go vet ./...`, `gofmt -l .`,
  `go run ./cmd/grove check`, `go test -count=1 -timeout 120s ./...`).
- The owner, reading the README as someone who has never seen Grove, judges
  that it gets them from install to a first shaped record.

## Next

Owner: assign with `/grove-work G-260927-ngkbz`.
