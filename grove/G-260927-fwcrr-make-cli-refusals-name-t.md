---
id: "G-260927-fwcrr"
type: work
title: "Make CLI refusals name the rule they enforce, and shrink the record model to its contract"
status: review
created: "2026-09-27T22:13:35Z"
updated: "2026-09-28T16:57:28Z"
size: medium
relates_to: ["G-260919-shnj5", "G-260921-ebsby", "G-260925-khfe7"]
candidate: "f2d9894"
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

## Evidence

Branch `worktree-G-260927-fwcrr`, base main `bb39668`; implementation
`5f9c204` (refusals), `e420055` (model split), `c0fb447` (review fixes).
Started from this record at `sha256:c18793e5…` and plan
[G-260928-rv8z2](G-260928-rv8z2-refusals-and-model-plan.md) at
`sha256:5cbddbd4…` (commit `40844ab`).

**Correction to Constraints.** The CLI enforces no transition order: `--set
status=active` on a done record succeeds unless it carries `approved`. The
third refusal on G-260927-ngkbz came from the approval rule alone, so its
message names that rule. A "done stays done" rule would be new behaviour
this record does not select; if the owner wants one, it is a new record.

1. **Refusals.** The record types carry their own fields, and `Kinds`,
   `Sizes`, `ConfigKeys`, `RunKeys` and the policy keys are tables that
   validation, its messages and `update` read alike
   (`internal/project/metadata.go`, `project.go`, `policy.go`,
   `internal/update/update.go`). `TestRefusalsNameTheirRule`
   (`internal/project`) builds each expected message from those tables for
   every type's status and fields, kind, size, type, configuration, `run:`
   and `policy:` keys, and each lifecycle rule; `TestUpdateRefusalsNameTheirRule`
   (`internal/update`) checks the text end to end and that the correction a
   lifecycle refusal names succeeds. The three refusals, rerun with a binary
   built at `c0fb447` in a disposable clone:

   ```text
   $ grove update G-260919-rt9h9 --set status=bogus
   …:5: status: expected proposed, active, review, done or abandoned for work
   $ grove update G-260919-rt9h9 --set size=huge
   …:8: size: expected small, medium or large
   $ grove update G-260927-n4wvk --set status=active     # done, approved
   …:11: approved: approval holds only while status is review or done, not active: unset approved, or set status review or done
   $ grove update G-260927-n4wvk --set status=active --unset approved
   {"changed":true,…}
   $ grove update G-260919-rt9h9 --set zzz=1
   grove: zzz is not a field that update accepts on work records; it accepts type, title, status, relates_to, kind, size, priority, members, depends_on, candidate or approved
   ```

   `list --status bogus`, `new note` and `convert --type note` list the
   statuses and types the same way.
2. **Model.** `grove guide model` is 11,484 bytes (was 33,844): the
   configuration keys, files and identity under "On-disk contract" (its
   three linked headings kept), one table of types, statuses and own fields,
   one of every field's form, values and meaning, the enforced lifecycle and
   what is not enforced, and what `check`, `update` and `convert` refuse.
   `model_test.go` fails if it exceeds 12 KB or if its statuses, fields,
   envelope, kinds, sizes or configuration keys differ from the code's
   tables (checked by deleting a status: it failed).
3. **Design.** The reasoning moved to the unshipped
   [`docs/record-design.md`](../docs/record-design.md), named in
   `CLAUDE.md`'s owner list; command behaviour the old model held beyond
   `grove --help` moved to `docs/commands.md` ("Records", "Judging and
   integrating"). A sentence-level comparison finds nothing shared between
   the model, the design document and the command reference.
4. Owner judgment: the three refusals above.
5. `README.md`'s command table points at the new command sections;
   `README.md` and `init`'s text name "Configuration and discovery", which
   the model keeps; every `record-model.md#…` anchor in the repository,
   records included, resolves. The work and shaping guides now send an agent
   to the model when "a refusal leaves the fix unclear". The
   [Shipped document](G-260925-khfe7-shipped-document.md) term lists the
   record design as not shipped.

Verification at `c0fb447`: `go vet ./...` clean, `gofmt -l .` empty,
`go run ./cmd/grove check` `OK: 220 records`, `go test -count=1 -timeout
120s ./...` all packages ok. No TUI change, so the terminal script was not
run; no concurrency change, so no `-race`.

Review: [G-260928-a372k](G-260928-a372k-fwcrr-review.md), two rounds by
fresh `grove-reviewer` agents; round 1 found six, all fixed in `c0fb447`;
round 2 `Open findings: none`.

Limits: `model_test.go` checks the policy subkeys by substring and does not
tie priority's range to code; `grove --help`'s text was not changed.

**Merge of main at `2633d8d`** (the delegated feedback, 2026-09-28, that
candidate `16af7ee` conflicted with main in `docs/record-model.md`): merge
`ecca7cf` joins the previous candidate's branch tip `7e10a56` and main
`2633d8d`, no rebase, so `16af7ee` and review G-260928-a372k's `c0fb447`
stay ancestors. The one conflicting file, `docs/record-model.md`, was
settled by keeping this branch's contract-only model and carrying each of
main's three sentence changes to where this branch had moved that sentence;
nothing of main's was dropped:

- `size: small` selects the compact handoff: the model's `size` row, and
  the planning reasoning in `docs/record-design.md`.
- A review record before Review is unenforced only "where the work guide's
  handoff calls for one": the model's "Not enforced" paragraph.
- `show --json` adds `approved_by` while `approved` is set: the `show ID`
  bullet in `docs/commands.md` (main's own policy paragraph there merged
  cleanly).

Every other file merged without conflict and is main's. The model is 11,604
bytes; every `record-model.md#` anchor resolves. Verification at `ecca7cf`:
`go vet ./...` clean, `gofmt -l .` empty, `go run ./cmd/grove check` `OK:
228 records`, `go test -count=1 -timeout 120s ./...` all packages ok.
Review: [G-260928-wp37y](G-260928-wp37y-fwcrr-merge-review.md), scoped to
the resolution, two rounds by fresh `grove-reviewer` agents; round 1 found
the missing `approved_by` clause, fixed in the merge; round 2 `Open
findings: none`. Main has since moved to `fc5a095` (one record commit);
`git merge-tree` of `ecca7cf` with it is clean.

## Next

In review: the merge of main `2633d8d` (`ecca7cf`) with its evidence. Judge the
three refusals in Evidence (acceptance 4) and read `grove guide model`
against acceptance 2. Then, in this checkout:
`grove approve G-260927-fwcrr "VERDICT"`, and in the target's checkout:
`grove integrate G-260927-fwcrr --cleanup`.

Feedback on candidate 16af7ee, 2026-09-28: delegated under policy grove.yaml sha256:182036ce84beda7043a09222a7e22419798a46d27e4511f5e747a5032a5d7dcb, budget 10 USD: conflicts with main at 2633d8d in docs/record-model.md. Resolve only that (grove resolve): in this branch, git merge 2633d8dd4cfb7591a3089649ed6d7b3762bfc997, that commit of main even if main has moved since, never a rebase; resolve those files keeping both sides' intent; rerun the repository's verification; and hand off the merge as the new candidate, with the previous candidate 16af7ee, the merged commit and the resolved files in Evidence. Change nothing else. If a resolution needs a choice this record does not settle, stop with a checkpoint naming it. Done at `ecca7cf`; see Evidence.
