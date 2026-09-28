---
id: "G-260927-dx0yn"
type: work
title: "Retain per-attempt process facts and show cost per work record"
status: active
created: "2026-09-27T22:13:36Z"
updated: "2026-09-28T01:17:49Z"
size: medium
relates_to: ["G-260921-sth8q", "G-260923-tnn5e", "G-260923-895zb", "G-260921-h46pb", "G-260923-p5pt6"]
candidate: "8cf2e4c5e19f4bcac1141e13c1879e5e5d8950de"
---

## Outcome

Each attempt Grove ran tells the owner what it spent its turns on and what
it changed, and each work record shows what its attempts have cost
together, so the effect of a workflow change can be read from the attempts
Grove already keeps instead of from transcripts by hand.

Owner intent, conversation 2026-09-27: "we should definitely instrument
attempts further as you suggested": process share, time to first code
edit, files changed outside the record root, and cost per record. The
facts and their placement are proposed below.

## Constraints

Observed 2026-09-27 at main `fef102f`:

- `grove attempt ATTEMPT` prints the launch (`attempt.json`: model, effort,
  budget, permission mode, command, worktree, base, entrypoint digests),
  the result event (turns, USD, denials), cost by model, event counts, the
  worktree after (HEAD, clean or dirty) and the record on the branch.
  `result.json` holds `dirty`, `events`, `exit_code`, `finished`, `head`,
  `record` and `stopped`. `Metrics` (`internal/attempt/activity.go:54-66`)
  already derives turns, input and output tokens, context, subagents,
  compactions, tool calls and tool errors from a bounded window of
  `events.jsonl` for the board's attempts screen.
- Not retained or derived: which tool calls were process (a `grove`
  command, worktree or branch setup, an edit under the record root, a guide
  print) and which were the task; the time and tool index of the first edit
  outside the record root; the files changed between base and HEAD, inside
  and outside the record root; guide prints by name.
- `grove attempts` lists a COST column per attempt and no total per work.
  The 43 finished attempts here total $250.93, mean $5.84, range $0.64 to
  $16.15 (each attempt's result event, read 2026-09-27); several records
  had two to four attempts.
- Measured by hand from this repository's 48 interactive work sessions:
  about a quarter of tool calls were process by that classification, and in
  the 19 sessions that edited a file outside `grove/` the first such edit
  came after a median 25 tool calls and 4 minutes. Interactive sessions are
  outside Grove's reach; attempts are not.
- `docs/commands.md` "Attempts" documents what an attempt retains, and the
  [Attempt](G-260921-sth8q-attempt.md) term points there. Reads of a running
  attempt are bounded (`ReadActivity`), a rule
  [G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) kept.

**Proposed design.** Derive at read time from `events.jsonl`, as `Metrics`
does: process tool calls and their share (Bash whose program is `grove`,
`git worktree` or a branch created under the worktree prefix, Edit or
Write under the record root, a Read of a record or guide), guide prints by
name, and the elapsed time and tool index of the first edit outside the
record root. At finish, `result.json` gains the files changed between base
and HEAD, inside and outside the record root, since the worktree may be
removed later. `grove attempt` prints these under one `Shape:` line; `grove
attempts ID` adds a total line (attempts, USD, turns, minutes); the board's
attempts pane shows the same total on the record; `--json` carries all of
it. For a running attempt the window's facts are lower bounds, labelled as
`Metrics` labels its own. A `grove` run through a wrapper or `$(…)` is
missed, as the evals runner already states of its own reader.

## Acceptance

1. `grove attempt` shows for a finished attempt: process tool calls and
   share, guide prints by name, the first edit outside the record root
   (elapsed, tool index, or none), and the files changed inside and outside
   the record root; an unknown fact is labelled unknown, never filled.
2. `grove attempts ID` and the board show the total cost and count over the
   work's attempts.
3. A test on a recorded `events.jsonl` fixture asserts each derived fact and
   the classification's documented limits; bounded reads of a running
   attempt are preserved.
4. `docs/commands.md` "Attempts" states the new facts and how they are
   derived.
5. Owner judgment in the terminal on one real attempt.

## Evidence

Headless session of 2026-09-28 on branch `worktree-G-260927-dx0yn`, from
`main` at `38511c1`. It started from this record at `sha256:ac01fd5a…` and
wrote plan [G-260928-j8sc0](G-260928-j8sc0-attempt-shape-facts-and.md), at
`sha256:06e6087b…` when implementation began. Commits: `6959b93` (plan),
`133afcb` (active), `6fbc8ff` (implementation), `e096c33` (review fixes).

What changed:

- **`Shape:`** (`internal/attempt/shape.go`, `ReadShape`). It is derived
  when `grove attempt` reads the attempt, from all of `events.jsonl`,
  bounded per line by the loop `ReadEvents` now shares (`eachLine`). The
  record root comes from `grove.yaml` at the attempt's base commit. The
  line shows the tool calls, the process calls and their share, guide
  prints by name, and the first edit outside the record root: its tool
  index, the time from the start, and the path. A running attempt's facts
  are "so far"; a skipped line makes the counts `≥`; an unknown record
  root reads `Shape: unknown: …`. The board does not read it, and keeps its
  bounded window.
- **`Changed:`** is recorded in `result.json` at finish: `git diff
  --name-only BASE HEAD`, split under and outside the record root, with up
  to ten names each. An older attempt reads `Changed: unknown: not recorded
  when it finished`, and a Git failure says why.
- **Totals.** `attempt.Sum`. `grove attempts ID` ends with, for example,
  `Total: 4 attempts, $13.79, ≥124 turns, 80m`, and the board's attempts
  list for one work opens with the same line. `events.turns` sums
  `num_turns` over a run's result events. Each counts one query, and the
  cost is cumulative: this repository's `G-260923-895zb.20260923T200725Z`
  has 95 and 14 turns at $9.75 in both. The `≥` marks attempts that only
  hold their last result's turns. An unfinished attempt, or one without a
  result event, is said beside the total, not summed.
- `attempt --json` carries `shape`, `shape_error`, `result.changed` and
  `result.events.turns`.

Against the acceptance:

1. Real output, `grove attempt G-260927-n4wvk.20260928T004357Z`:
   `Shape: 36 tool calls, 10 process (28%); guides printed: work 1; first
   edit outside the record root: tool 23, 3m59s after the start,
   docs/commands.md` and `Changed: unknown: not recorded when it finished`.
   The second is right, because the attempt predates this. Every existing
   attempt finished before `Changed` was recorded, so only tests show it
   filled: `TestChangedAndShapeFacts`, and `TestRunToResult` through a
   fake provider.
2. `grove attempts G-260923-fwakw` ends `Total: 2 attempts, $11.79, ≥22
   turns, 26m`. `TestAttemptScreensReconnectAndStop` asserts the board's
   row for one work, and asserts that the list of all attempts has none.
   `TestAttemptCommandsUsage` asserts the CLI line.
3. `TestReadShape` runs on the recorded fixture
   `internal/attempt/testdata/shape-events.jsonl`: 19 calls, 10 of them
   process. It checks guide prints, the first edit (a subagent's) and each
   documented limit: `env grove`, `$(grove …)`, `FOO=1 grove`, `sed -i`,
   the `/private` spelling and the heredoc false positive. It also checks
   an oversized line, a running attempt, and no edit. `TestSum` covers the
   totals and their labels. `TestReadEventsBounded` and the board's
   `TestAttemptActivityIsBounded` pass unchanged.
4. `docs/commands.md` "Attempts" states the facts, how they are derived,
   and the known misses. `docs/board.md` names the board's total.
5. This is for the owner, in a terminal.

Verification at `e096c33`, all passing: `go vet ./...`, `gofmt -l .`,
`go run ./cmd/grove check` (OK: 217 records), `go test -count=1 -timeout
120s ./...`, and `go test -short ./internal/attempt` (2.8 s).
`python3 internal/tui/testdata/terminal.py` passed at `6fbc8ff` (12
scenarios). `e096c33` changes no TUI code.

Review: [G-260928-4ae1r](G-260928-4ae1r-review-of-g-260927-dx0yn.md), by an
independent `grove-reviewer` subagent in two rounds. Round 1 on `6fbc8ff`
found four low findings, all fixed in `e096c33`. Round 2 closed with `Open
findings: none`.

Limits:

- The classification misses some things, as `docs/commands.md` states.
  Edits made through Bash and guides read as files are the large ones:
  earlier attempts here read `docs/work-execution.md` with `cat` and wrote
  records with scripts. So these shares are lower than the hand count's
  quarter, for example 16% on `G-260923-fwakw.20260923T182316Z`.
- `Result event:` still prints the last result event's turns, as before.
- `attempts` without an ID prints no total, and `attempts` has no `--json`.

## Next

In review: the candidate (the `candidate` field) is on branch
`worktree-G-260927-dx0yn`, from `main` at `38511c1`. The owner judges one
real attempt in a terminal (acceptance 5). Run `go run ./cmd/grove attempt
ATTEMPT` on any attempt listed by `go run ./cmd/grove attempts`. After
merging, attempts launched from that point also show `Changed:`. Also run
`go run ./cmd/grove attempts G-260925-pbx81`, and on the board `A` on a
work's detail. Then, in this worktree:

```sh
go run ./cmd/grove approve G-260927-dx0yn "VERDICT"   # or: go run ./cmd/grove feedback G-260927-dx0yn "TEXT"
```

and in the `main` checkout:

```sh
go run ./cmd/grove integrate G-260927-dx0yn
```

Feedback on candidate 8cf2e4c, 2026-09-28: Process share and first edit are noise: across 46 attempts 91% of tool calls are Bash, ~800 file writes go through Bash vs 110 Edit/Write, record reads 337 via cat vs 31 via Read; first edit reads 'none' in 17/46 and is late or missing in 34/46. Drop both; keep Changed, totals, grove/guide counts. Optionally add time to first commit touching files outside the record root, from git at finish.
