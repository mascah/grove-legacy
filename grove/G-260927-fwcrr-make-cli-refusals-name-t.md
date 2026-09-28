---
id: "G-260927-fwcrr"
type: work
title: "Make CLI refusals name the rule they enforce, and shrink the record model to its contract"
status: active
created: "2026-09-27T22:13:35Z"
updated: "2026-09-28T16:23:41Z"
size: medium
relates_to: ["G-260919-shnj5", "G-260921-ebsby", "G-260925-khfe7"]
---

## Outcome

An agent whose `update` is refused can correct the command from the refusal
alone, and an agent that must know a field's meaning or allowed values reads
a record model that states the contract and nothing else.

Owner intent, conversation 2026-09-27: "Refusals should carry the rule. The
record model is likely too lengthy right now." The message changes and the
split below are proposed.

## Constraints

Observed 2026-09-27 at main `fef102f`:

- `update` validates the changed record and wraps the diagnostics as "the
  update would leave ID invalid:" with one line per problem
  (`internal/update/update.go:132`). Three refusals tried on
  G-260927-ngkbz: `--set status=bogus` printed `status: unsupported
  lifecycle value for work` and named no value
  (`internal/project/metadata.go:276`, where the type's `Statuses` slice is
  in scope); `--set size=huge` printed `size: expected small, medium, or
  large`; `--set status=active` on a done record printed only `approved:
  applies only while status is review or done; unset it when reopening the
  work`, a consequence, not the rule that a done record stays done. `X is
  not a field that update accepts on TYPE records`
  (`internal/update/update.go:243`) lists no accepted field though `allowed`
  is in scope. Forty message sites in `internal/project` carry a rule.
- [`docs/record-model.md`](../docs/record-model.md), what `grove guide model`
  prints, is 593 lines and 33,844 bytes: Identity and placement apart from
  classification 5,140; Knowledge records and the brief 2,197; Minimum
  representation 918; Work planning metadata 1,953; On-disk contract
  15,467; Lifecycle and validation boundary 7,357; Not records 498. The
  first three sections and much of the rest are design reasoning; the
  contract an agent needs (configuration keys, types, fields and their
  values, statuses and transitions, what `check` and `update` refuse) is
  spread through them.
- The work guide sends an agent to the model "when a field's meaning or
  allowed values matter or the CLI refuses a change"; in 16 of the 48
  interactive work sessions recorded for this repository the agent printed
  it, about 8.5k tokens each time that then stay in context.
- `internal/cli/init.go:99` names the model's "Configuration and discovery"
  heading in the text `init` prints; `CLAUDE.md` names `docs/record-model.md`
  as the owner of configuration, schema, validation and lifecycle; the
  [Shipped document](G-260925-khfe7-shipped-document.md) rules apply to
  whatever `guide model` prints.

**Proposed design.** Two changes in one record, because the second decides
how much the first must say. (1) Every refusal that enforces an enumerated
value lists the values (`status: expected proposed, active, review, done or
abandoned for work`), every refusal of a field lists the fields `update`
accepts on that type, and every lifecycle refusal names the rule before its
consequence (`done stays done: a done record's status cannot change; …`).
Messages stay one line each. (2) The record model is split: `grove guide
model` prints the contract, at most 12 KB, in reference form (a table per
type of fields, values and meaning; the statuses and transitions per type;
the configuration keys; what `check` refuses), and the design reasoning
moves to a repository document that is not shipped (proposed
`docs/record-design.md`), linked from `CLAUDE.md`'s owner list, so no fact
has two homes. Both documents are reconciled with `README.md`,
`docs/commands.md`, the `init` text and this repository's `CLAUDE.md`.

## Acceptance

1. A table-driven test asserts that every refusal of an enumerated field
   lists that field's values, that an unaccepted field lists the accepted
   ones, and that each lifecycle refusal names its rule; the three refusals
   above print the new form, shown in Evidence.
2. `grove guide model` is at most 12 KB and states every field, value,
   status, transition, configuration key and refusal the CLI enforces
   exactly once; a test checks the enumerations it prints against the
   code's tables so they cannot drift.
3. The design reasoning has one editable home in this repository, reachable
   from `CLAUDE.md`'s owner list, with no paragraph duplicated between the
   two documents.
4. Owner judgment: for each of the three refusals above, the corrected
   command is evident from the message without printing the model.
5. Shipped-document checks and the Go suite pass; `README.md`,
   `docs/commands.md` and the `init` text name the right headings.

## Next

Assign: `/grove-work G-260927-fwcrr`.
