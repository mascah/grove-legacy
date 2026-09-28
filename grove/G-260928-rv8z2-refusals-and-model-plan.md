---
id: "G-260928-rv8z2"
type: plan
title: "Refusals that name their rule, and a contract-only record model"
status: current
created: "2026-09-28T16:22:59Z"
updated: "2026-09-28T16:23:31Z"
work: ["G-260927-fwcrr"]
---

# Plan for G-260927-fwcrr

Base: `main` at `bb39668`, branch `worktree-G-260927-fwcrr`, one implementer.
Reproduced there in a disposable clone: `--set status=bogus` prints
`status: unsupported lifecycle value for work`, `--set size=huge` prints
`size: expected small, medium, or large`, `--set kind=x` prints `kind:
unsupported work kind`, `--set zzz=1` prints `zzz is not a field that update
accepts on work records`.

**Correction to the record's Constraints.** `--set status=active` on a done
record *succeeds* when it carries no `approved`: the CLI enforces no
transition order (the model's "Not enforced by software" says so). The
refusal seen on G-260927-ngkbz came from the `approved` rule alone. So the
rule that message must name is "approval holds only while status is review
or done", not "done stays done". Adding a done-stays-done rule would be new
behaviour the record does not authorize; the handoff raises it for the owner.

## Design

### 1. One table per enumeration, used by the checks and their messages

In `internal/project/metadata.go`:

- `TypeInfo` gains `Fields []string`, the type's own fields beyond the
  envelope; `Envelope` lists the fields every record may carry (`id`,
  `type`, `title`, `status`, `relates_to`, `created`, `updated`,
  `formerly`; `status` dropped for a type without statuses). `ParseRecord`
  builds `allowed` from these instead of appending per type.
- `Kinds` and `Sizes` are exported slices; `ConfigKeys`, `RunKeys` and the
  policy key lists in `project.go`/`policy.go` become slices the checks and
  messages both read.
- A small `list(values)` helper renders `a, b or c`.

`internal/update/update.go` derives the fields `update` accepts from the
same table (`title`, `status` where the type has statuses, `relates_to`,
then `Fields`), so the map it keeps today goes away.

### 2. Messages (rule first, one line each)

- `status: expected proposed, active, review, done or abandoned for work`
  (every type; page: status is unknown, as now).
- `kind: expected feature, fix, refactor, investigation, tooling or release`;
  `size: expected small, medium or large`.
- `type: unknown record type; expected work, question, …` from `Types`.
- Unknown frontmatter key: `KEY: not a field of TYPE records, which take
  id, type, title, …` (unknown type: the old wording).
- `update`: `KEY is not a field that update accepts on TYPE records; it
  accepts title, status, …, and type` .
- Unknown `grove.yaml` key lists the six keys.
- `approved`: `approval holds only while status is review or done; status is
  active, so unset approved in the same update`; `approved must equal
  candidate COMMIT: approval is of one commit, and a changed candidate
  needs its own`.
- `candidate: required while status is review: set candidate=COMMIT, the
  commit offered for judgment, in the same update`.
- The `integrated` refusals already lead with their rule; audit the rest of
  `internal/project` (about forty sites) and change only those that state a
  consequence without the rule or an enumeration without its values.

Existing tests that match old substrings are updated to the new text.

### 3. Record model split

`docs/record-model.md` (shipped as `grove guide model`) becomes the
contract in reference form, at most 12 KB:

- `## On-disk contract` with `### Configuration and discovery` (a key
  table: `schema_version`, `records`, `brief`, `target`, `run.*`,
  `policy.*`, and the approval conditions a policy cannot waive), `###
  Folders and files`, `### Identity and dates` (kept: README, commands.md
  and older records link these anchors).
- `## Types and fields`: one table of type, statuses (first is what `new`
  writes) and own fields; one table of every field, its form, allowed
  values and meaning; relationship rules.
- `## Work lifecycle`: statuses and the enforced transitions (`candidate`,
  `approved`, `done`), groups, and what is not enforced.
- `## Reading and writing records`: what `check` refuses, what `update`
  refuses, and `convert`, with command behaviour `grove --help` already
  gives left to it.
- `## Not records`.

Design reasoning (why neutral IDs, why pages are a type, the planning
distinctions, historical Done, the lock's placement and the clone collision
analysis, the brief's committed-source check) moves to a new, unshipped
`docs/record-design.md`, each paragraph in one of the two files only.
`CLAUDE.md`'s owner list names it. `README.md` (the "Configuration and
discovery" pointer and the command-table anchors), `docs/commands.md`'s
opening pointers and its `run:` link, and `internal/cli/init.go`'s heading
are reconciled. The work and shaping guides' "read the model when the CLI
refuses a change" becomes "when a refusal does not say how to correct it".

### 4. Tests

- `internal/project`: a table-driven test builds a refusal for each
  enumerated field on each type and asserts the message holds every value
  from the code's table; unknown field lists the type's fields; unknown
  configuration key lists `ConfigKeys`.
- `internal/update`: the same for `update`'s unaccepted field, and the three
  reproduced refusals end to end with their new text.
- Root package `model_test.go`: `grove guide model` is at most 12,288 bytes;
  its type table's statuses and fields, its `kind` and `size` values and its
  configuration keys equal `project.Types`, `Kinds`, `Sizes`, `ConfigKeys`,
  `RunKeys` and the policy keys.

## Steps

1. Tables and messages in `internal/project` and `internal/update`, with
   tests; `go test -short ./internal/...`.
2. Split the model; `model_test.go`; reconcile README, commands.md,
   CLAUDE.md, init text, guides.
3. Evidence: the three refusals rerun in a disposable clone; size of
   `guide model`; shipped-document test; full verification per `CLAUDE.md`.
4. Independent review, then handoff.
