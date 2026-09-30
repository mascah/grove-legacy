---
id: "G-260924-ecs9m"
type: work
title: "Default an attempt's budget, mode, model and effort from grove.yaml and launch from one line"
status: accepted
created: "2026-09-24T23:52:07Z"
updated: "2026-09-30T13:18:01Z"
relates_to: ["G-260921-h46pb", "G-260921-7trd7", "G-260924-5b6pz", "G-260924-59f5k", "G-260925-04ccr"]
candidate: "b18804e8a9194d596355893698e2d688a3c19130"
approved: "b18804e8a9194d596355893698e2d688a3c19130"
approved_by: owner
approved_context: "sha256:100d09f2b70ffc0795499c96d39c065e9273b12cc916f8c3ae5ba8df26cb121d"
---

## Outcome

The owner launches an attempt from the board with one keypress after `R`,
or with `grove run ID` and no flags, when `grove.yaml` holds the project's
launch defaults; a flag or a typed override still changes any value for one
launch; and every attempt records the values it ran with, as it does now.

Owner intent, shaping conversation 2026-09-24: after
[G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md) there are too many
questions when launching an attempt, and the defaults need a home such as
`grove.yaml`. The owner read the investigation below and chose this design
over the same defaults with five prefilled prompts, over provider-side
settings alone, and over only skipping the optional prompts. The key
spelling and the launch line are proposed design where labelled.

## Constraints

Observed at main `132709e`, 2026-09-24:

- **Five prompts on the board, two required flags.** `R` asks budget,
  permission mode, bound, model and effort in sequence, one line each
  (`launchKey` in [attempts.go](../internal/tui/attempts.go), the prompt
  texts in [review.go](../internal/tui/review.go)); Enter alone on the last
  three asks for nothing. `run` requires `--budget` and `--permission-mode`
  and refuses without them (`Start` in
  [attempt.go](../internal/attempt/attempt.go)); `--until`, `--model` and
  `--effort` are optional. G-260924-5b6pz added the three optional prompts; before it
  there were two.
- **What was typed.** The nineteen attempts under `.git/grove/attempts/`,
  from every `attempt.json`: budget `50` on seventeen, `30` once (G-260924-59f5k's
  bounded run under its experiment cap) and `50` once; permission mode
  `auto` on all nineteen; model empty on seventeen and `opus` on the two
  G-260924-59f5k attempts; effort set once (`high`); bound set once (`plan`). Two
  prompts retype constants and three are Enter.
- **Why they are required today.** [G-260921-h46pb](G-260921-h46pb-run-one-bounded-implemen.md) made
  both flags required "since a default for either would make the choice in
  practice" ([G-260923-0t43m](G-260923-0t43m-durable-attempt-plan.md), Budgets and
  bounds); [G-260921-7trd7](G-260921-7trd7-launch-and-inspect-manag.md) kept them typed with no prefill
  for that reason; G-260924-5b6pz put "role profiles or model defaults in
  `grove.yaml`" out of scope, "as G-260921-h46pb chose". This record reverses that
  one sentence of G-260921-h46pb, on the owner's request: the choice has been made
  in practice by repetition, and a default in a committed file is the choice
  made once, attributably, in Git history, where a retyped value is not.
- **`grove.yaml` today.** `parseConfig` in
  [project.go](../internal/project/project.go) accepts exactly
  `schema_version`, `records`, `brief` and `target` and reports any other key
  as unknown; values are yaml.v3 nodes read as scalars, with no helper for a
  nested mapping; `check` validates the file; `init` writes the first three
  keys (`defaultConfig` in [init.go](../internal/cli/init.go)). The versions
  loader parses every branch's committed `grove.yaml` through the same
  function, but a launch loads only the launching checkout's project, so
  branches need not agree on defaults. The
  [record model](../docs/record-model.md#configuration-and-discovery) owns
  the key list.
- **The provider.** Claude Code 2.1.282 (`claude --help`): `--effort` takes
  low, medium, high, xhigh or max; `--permission-mode` takes acceptEdits,
  auto, bypassPermissions, manual, dontAsk or plan; `--max-budget-usd` works
  only with `--print`. Its documentation, read on 2026-09-24 through the
  claude-code-guide agent: settings.json can set `model` and `effortLevel`;
  nothing but the flag sets a budget; `permissions.defaultMode` is honoured
  from user or managed settings only, never from a project's settings file.
  The owner's user settings already set `defaultMode: auto`. A provider-side
  default reaches an attempt's facts only as the init event's actual model
  and mode, never as a requested effort.
- **Code shape.** `internal/cli` imports `internal/tui`, and both import
  `internal/attempt`, so a flag table shared by `run` and the board belongs
  in the attempt package; today the table is in
  [cli.go](../internal/cli/cli.go), and `Start` validates the budget before
  it loads the project. The board's resolver carries the live checkout's
  `target` (`versions.Result`) and can carry its launch defaults the same
  way. `TestLaunchFromTheDetail`, `TestRefusals`, `TestAttemptCommandsUsage`
  and three scenarios in
  [terminal.py](../internal/tui/testdata/terminal.py) drive the prompts and
  refusals.

In scope, proposed design:

1. **`run:` in `grove.yaml`.** An optional mapping with optional `budget`,
   `permission_mode`, `model` and `effort`, named as `run`'s flags:

   ```yaml
   run:
     budget: 50
     permission_mode: auto
   ```

   `check` refuses a budget `ValidBudget` refuses, a value with whitespace,
   and any other key under `run`. No default for the bound: a bound is one
   launch's mandate (G-260924-5b6pz), and the brief selects no universal plan gate.
   `init` does not write `run:`, so a fresh project keeps G-260921-h46pb's rule until
   its owner adds a default. The alternative spelling, flat top-level keys,
   needs no nested-mapping helper at the cost of a `model:` beside
   `records:`; the implementer chooses and the record model says which.
2. **`run` and `Start` default from it.** `Start` loads the project first,
   fills each empty value from `run:`, then validates as now; the refusal
   for a missing budget or mode stays and names both ways to supply one.
   `attempt.json` records the resolved values, and nothing distinguishes a
   default from a flag. `run`'s usage and [commands.md](../docs/commands.md)
   say the two flags are required unless `grove.yaml` defaults them.
3. **One launch line on the board.** `R` replaces the five prompts with one
   line: the resolved launch (budget, mode, bound, model, effort) and where
   it runs, where Enter launches; text typed in `run`'s own flag spelling,
   such as `--until plan --effort xhigh`, overrides for that launch and is
   parsed by the table `run` uses, moved into the attempt package. Without a
   default for the budget or the mode, the same line says what must be
   typed and refuses to launch until it is. Esc cancels. The plan-ready
   outcome's "read it, then R" is unchanged.
4. **This repository's defaults.** Add `run:` with `budget: 50` and
   `permission_mode: auto`, the values the nineteen attempts used, to this
   repository's `grove.yaml`, as the first use.
5. Reconcile the [record model](../docs/record-model.md) (configuration),
   [commands.md](../docs/commands.md) (run and init) and
   [board.md](../docs/board.md) (R), each saying its part once.

Out of scope: per-phase or role profiles, such as one effort for `--until
plan` and another for the implementation, which G-260924-5b6pz defers until
[G-260924-59f5k](G-260924-59f5k-run-the-eval-pair-on-cod.md) reports; a default
bound; remembering the last launch; user-level Grove configuration; writing
defaults in `init`; and any change to what an attempt records.

## Acceptance

1. `grove.yaml` with `run: {budget: 50, permission_mode: auto}` passes
   `check`; `run: {budget: x}`, `run: {effort: "a b"}` and `run: {until:
   plan}` are refused by `check` with the key named; a file without `run:`
   means what it means today.
2. With those defaults, `grove run ID` with no flags starts an attempt whose
   command carries `--max-budget-usd 50 --permission-mode auto` and whose
   `attempt.json` shows them; `grove run ID --budget 2` carries 2 and
   `auto`; without `run:` and without the flags, `run` is refused as today
   by a message that also names `run:` in `grove.yaml`. Shown with the fake
   provider, and once with the real one in a disposable clone reached by an
   absolute `--project` path.
3. `R` on the board shows one line naming the resolved launch and where it
   runs; Enter starts the attempt with those values; typing `--until plan
   --effort xhigh` on that line starts one with the bound and effort and the
   defaults for the rest; a line that leaves the budget or the mode
   unsupplied is refused without launching; a word the table does not know
   is refused with the message `run` gives. `TestLaunchFromTheDetail` and
   the three pseudo-terminal scenarios cover it, and the five prompts they
   expect today are gone.
4. The usage text, the record model, `commands.md` and `board.md` describe
   the defaults and the launch line once each, every link resolves, and no
   document still states G-260921-h46pb's sentence as current.
5. `gofmt`, `go vet`, `check`, the full uncached suite and `terminal.py`
   pass, with no package over five seconds.

## Evidence

Implemented on `worktree-G-260924-ecs9m` from base `main` `670ca9c`, starting from
this record at `sha256:40169cdb…` with no separate plan: the design above
stood in for one, and the nested `run:` spelling needed only a small
`runField` beside a factored `addFields` in `parseMapping`. The code is
`d9ca2ef`, this repository's defaults `ad52710`, and the guide's row
`b917fc2`.

Per acceptance:

1. `grove.yaml` takes an optional `run:` mapping, read by `runField` in
   [project.go](../internal/project/project.go). `ValidBudget` moved there
   from the attempt package, and problems are named `run.KEY`. With the
   binary at `ad52710` in a disposable clone, `check` passed `run: {budget:
   50, permission_mode: auto}` and reported
   `grove.yaml:5: run.budget: expected a positive decimal dollar amount`,
   `run.effort: expected one word, without whitespace` and `run.until:
   unknown key; run: takes budget, permission_mode, model and effort`.
   Without `run:`, parsing is unchanged. The test is `TestRunDefaults`.
2. `attempt.Start` loads the project first, fills empty values with
   `Defaulted`, then validates, and `attempt.json` records the resolved
   values. The CLI still refuses a launch with neither flags nor defaults
   as a usage error, exit 2, with `ErrUnsupplied`, which names `run:` in
   `grove.yaml`. Tests: `TestDefaults` (fake provider: no flags gives
   `--max-budget-usd 50 --permission-mode auto`, and `--budget 2` gives 2
   and `auto`, in the command and in `attempt.json`), `TestRefusals` and
   `TestAttemptCommandsUsage`. With the real provider (Claude Code
   2.1.282), `grove --project /tmp/g140-clone run G-260925-9ptz2` with no flags in a
   disposable clone started
   `claude -p /grove-work G-260925-9ptz2 --interaction headless … --max-budget-usd 50
   --permission-mode auto --permission-prompts none`. Its `attempt.json`
   held `budget_usd 50` and `permission_mode auto`, and `grove stop` ended
   it after $0.12. The clone was then deleted.
3. `R` opens one line, for example `Launch G-260924-ecs9m ▏ · $50, mode auto, to the
   handoff, model default, effort default, on branch worktree-G-260924-ecs9m · Enter
   launches; …`. The typed flags come first, since truncation drops the
   end. `resolved` parses what is typed with `attempt.Flag`, the table
   `run` uses, now in the attempt package, over this checkout's defaults
   (`versions.Source.Run`), and the row resolves as you type. Enter launches
   exactly the values shown. A line that leaves the budget or the mode
   unsupplied, anything `run` refuses (`unknown option --frob`,
   `--until must be plan`, …) and `--branch` or `--worktree` all launch
   nothing. Tests: `TestLaunchFromTheDetail`, `TestLaunchPlace`, the
   feedback continuation, and `terminal.py`'s `attempt_lifecycle`. That
   scenario launches three times: typed flags without `run:`, Enter on the
   defaults, and `--effort xhigh` over them, and checks that each
   `attempt.json` recorded budget 1 and mode `auto` all three times, with
   effort `xhigh` only on the third. The five prompts are gone.
4. `run`'s usage text, the [record model](../docs/record-model.md#configuration-and-discovery)
   (configuration), [commands.md](../docs/commands.md) (Attempts, Init),
   [board.md](../docs/board.md) (Attempts) and the work guide's invocation
   row each describe their part. A search for "no default", "neither has a
   default" and "typed each time" finds no statement of G-260921-h46pb's sentence as
   current. Every new link resolves.
5. At `ad52710`: `gofmt -l .` empty, `go vet ./...` ok, `go run ./cmd/grove
   check` OK (141 records), `go test -count=1 -timeout 120s ./...` all ok,
   and `python3 internal/tui/testdata/terminal.py` 11 of 11 ok. `go test
   -short -count=1 ./...` was rerun at `b917fc2`, all ok. Under `-short`,
   `cli` (3.4 s at the tip, 3.1 s at base) and `versions` (5.4 s, 5.5 s)
   match base in isolation. `tui`, `versions` and `attempt` exceed five
   seconds in full uncached runs at base too; `attempt` adds the non-short
   `TestDefaults`. The base checkout's full `attempt` run once failed
   `TestOwnerLost` under concurrent load, which is a base flake and not
   this change.

Decisions taken:

- **Nested `run:`, keys named as the flags** (`permission_mode` in YAML).
  A budget YAML types as a number is accepted as its literal text.
- **The board fills the resolved values into the request itself**, so a
  `grove.yaml` edited after the board was read cannot change what the line
  showed. `Start` defaults again only for callers that leave values empty.
- **`--branch` and `--worktree` are refused on the board line**, because
  the board chooses where an attempt runs from what it shows.
- **The CLI's missing-flags refusal stays a usage error (exit 2)**, checked
  after the project loads. `run` in an invalid project now reports the
  project's diagnostics first (exit 1).
- **The author's reading of [G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md)**,
  which was accepted after this record was shaped: a `run:` default
  committed by the owner in the project's `grove.yaml` is the owner's
  explicit answer for what an attempt spends, made once and attributable in
  Git. It is not a provider default like `~/.codex/config.toml`, and not a
  blank an agent reads as approval. Every launch is still the owner's
  keypress or command, and the board shows the values before Enter. The
  owner has not confirmed this reading.

Review: [G-260925-9ptz2](G-260925-9ptz2-launch-defaults-review-2.md), an independent
`grove-reviewer` of `ad52710`, found acceptance met and no correctness
defect. Its six findings and their dispositions are there. After it, only
the guide's row (`b917fc2`) and these records changed.

Limits: the defaults come from whichever checkout launches, so `grove run`
inside a work branch's worktree uses that branch's `run:`, which an attempt
could have edited, and the CLI names the values only once the attempt has
started. A binary older than this change refuses a `grove.yaml` with `run:`.

## Next

In review: candidate on `worktree-G-260924-ecs9m`, base `670ca9c`. For the owner to
judge, most important first:

1. Whether the G-260925-04ccr reading above holds. If it does not, the defaults
   should not spend without a typed confirmation.
2. Whether one keypress after `R` is what you wanted: `go build -o
   /tmp/grove-g140 ./cmd/grove` in this worktree, then `/tmp/grove-g140`
   and `R` on a proposed work item, then Esc.
3. Whether the launching checkout, rather than the target, should own the
   defaults.

Then, in this worktree:
`grove approve G-260924-ecs9m "VERDICT"`, and in `main`'s checkout:
`grove integrate G-260924-ecs9m --cleanup`. After the merge, rebuild the installed
`~/.local/bin/grove`: the one at `670ca9c` refuses the merged `grove.yaml`.

Verdict on candidate b18804e, 2026-09-25: approve

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
