---
id: "G-260927-ngkbz"
type: work
title: "Frame the README for someone adopting Grove in their project"
status: active
created: "2026-09-27T17:10:58Z"
updated: "2026-09-27T17:15:58Z"
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

## Evidence

On `worktree-G-260927-ngkbz`, from base `main` `90e82d4` and record revision
`sha256:57e1efe2d3d9`. No plan was written; the design above served as the
plan. Implemented in `558bce8`, with the review fixes in `cc87c6b`. The
candidate is the commit that records this Evidence.

- **Order.** The README runs:
  - an intro in outcome terms;
  - Install: `go install …@latest` or `@v0.1.0`, then `grove version`;
  - Set up a project: `init`, what it writes, the minimal `grove.yaml` with
    `target: main`, one bullet each on `target`, `run:` and `policy:`,
    `grove guide model`, then `check` and commit, then upgrading;
  - The loop: brief, shape, work, judge, with `a`, `f` and `i`;
  - Commands: the board paragraph and the table, with `grove guide model` or
    `grove guide work` beside the rows those guides cover;
  - Develop Grove.

  `go run` and record IDs appear only in Develop Grove.
- **`grove.yaml`.** I built the binary to a temp path and ran `git init -b
  main` then `init` in a `mktemp -d` project. With the README's yaml block
  written verbatim, `check` printed `OK: 0 records` (exit 0) and `list`
  exited 0. The key descriptions agree with the record model's
  "Configuration and discovery". The `target` bullet names what refuses
  without it (`integrate`, `resolve`, `sweep`).
- **Development-only rows.** "Where each subject lives" is gone, along with
  its dogfooding-evidence records, legacy-ID row and adapter links. Develop
  Grove keeps `grove/` and the brief, `go run`, `just install`, CLAUDE.md
  and `evals/README.md`, which nothing else in the tree links. A script
  resolved every relative link and anchor in README.md, CLAUDE.md and
  docs/commands.md.
- **`init` note.** `internal/cli/init.go` adds this sentence: "grove guide
  model describes grove.yaml's keys under "Configuration and discovery",
  such as target: BRANCH for the branch work merges into".
  `internal/cli/init_test.go` asserts it, and docs/commands.md's Init
  describes it. `grove guide model` prints that heading.
- **CLAUDE.md.** The README line now says the README covers what Grove is
  and how a project adopts it, with developing Grove last.
- **G-260921-czt8x.** Its note on `W-001` examples no longer names the
  README, which dropped that row.
- **Final checks.** At `cc87c6b`:
  - `go vet ./...` passed;
  - `gofmt -l .` printed nothing;
  - `go run ./cmd/grove check` printed `OK: 202 records`;
  - `go test -count=1 -timeout 120s ./...` passed every package.

  The same checks also passed at `558bce8`. There is no TUI change, so the
  terminal checks were not run.
- **Review.** [G-260927-917tb](G-260927-917tb-review-of-g-260927-ngkbz.md)
  ran two `grove-reviewer` rounds. Round 1 found two minor issues: the
  `target` bullet understated what depends on the key, and czt8x was stale.
  Both were fixed in `cc87c6b`. Round 2 ended `Open findings: none`.
- **Limits.**
  - `@latest` resolves to `v0.1.0`, which predates the new `init` note, so
    installed adopters see it only after the next tag. A release is out of
    scope.
  - The last acceptance item, the owner reading the README as a newcomer,
    is not judged yet.

## Next

In review. The owner reads `README.md` on `worktree-G-260927-ngkbz` as
someone who has never seen Grove (the last acceptance item). Then:

```sh
# in .claude/worktrees/worktree-G-260927-ngkbz
go run ./cmd/grove approve G-260927-ngkbz "VERDICT"
# in the main checkout
go run ./cmd/grove integrate G-260927-ngkbz --cleanup
```

Or `go run ./cmd/grove feedback G-260927-ngkbz "TEXT"` in the worktree.
