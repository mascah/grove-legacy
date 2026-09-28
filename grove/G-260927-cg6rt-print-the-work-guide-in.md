---
id: "G-260927-cg6rt"
type: work
title: "Print the work guide in stages: a head at start, later steps on demand"
status: active
created: "2026-09-27T22:13:35Z"
updated: "2026-09-28T00:44:29Z"
size: medium
relates_to: ["G-260925-khfe7", "G-260925-m9jcr", "G-260925-p2k54", "G-260923-p5pt6", "G-260927-dx0yn"]
---

## Outcome

An agent carrying assigned work reads the part of the work guide its current
step needs, so the guide stops occupying a fixed 34 KB of every turn's
context, while nothing an agent must follow moves out of the guide.

Owner intent, conversation 2026-09-27: "make the rest on demand", agreeing to
a head printed at start and later steps fetched when reached. The split and
the mechanism below are proposed.

## Constraints

Observed 2026-09-27 at main `fef102f`:

- `grove guide work` prints the whole of
  [`docs/work-execution.md`](../docs/work-execution.md), 34,472 bytes
  (`internal/cli/cli.go:229-240` reads `grove.GuideFiles["work"]` entire;
  `--entrypoint` only checks that the revision is served). Section sizes in
  bytes: Lifecycle 852, Inputs 1,491, Authority 1,088, Read in stages 2,333,
  steps 1 to 3 together 5,175; then Prepare 1,624, Implement 6,010, Review
  1,812, Checkpoint 787, When a human decision is missing 2,370, Hand off
  3,860, Judging and integrating 4,159, Invocation 1,666. The intro through
  step 3 is 12,184 bytes, 35% of the file.
- In the 48 interactive `/grove-work` sessions recorded for this repository
  (Claude Code transcripts under `~/.claude/projects`, read 2026-09-27) the
  median session ran 62 turns; the work guide was printed in 27 main threads
  and the record model in 16. Whatever is printed stays in context for every
  later turn, so the guide's cost is its length times the turns after it,
  not one read.
- The guide's own "Read in stages" table says "Starting: This guide" in
  full, and the entrypoint revision 2 skills say "Run `grove guide work
  --entrypoint 2` and follow the guide it prints".
- The [Entrypoint revision](G-260925-m9jcr-entrypoint-revision.md) term: a
  revision changes only when what entrypoints need of the binary changes,
  never for wording. The [Shipped document](G-260925-khfe7-shipped-document.md)
  term: the guide links only within itself or to `https://`, names other
  shipped documents by command and names no record or path; `grove
  version`'s guides digest changes with any edit, and the evals record it.
- The behavioral evals ([G-260923-p5pt6](G-260923-p5pt6-establish-behavioral-eva.md))
  exercise headless shaping only. No case runs the work guide, so its
  regression check is a real attempt through `grove run`, read from its
  events.

**Proposed design.** One file stays the one editable owner. `grove guide
work [--entrypoint N]` prints the head: the intro, Lifecycle, Inputs,
Authority, Read in stages, steps 1 to 3, and a closing table naming each
remaining part, the command that prints it and the step that needs it.
`grove guide work --part NAME` prints one part alone: `prepare` (step 4),
`implement` (5), `review` (6), `checkpoint` (7 together with "When a human
decision is missing"), `handoff` (8), `judge` ("Judging and integrating a
candidate") and `invocation`; `--part all` prints the whole file as today.
The binary splits the file at its `##` headings by a fixed table, and a test
asserts that the head and the parts, in order, are the file byte for byte.
The head's Read in stages table gains a row per part. The shaping and
review guides are out of scope: at 15,940 and 3,058 bytes each is read once.
Whether the entrypoint revision stays 2 is the implementer's reading of
G-260925-p2k54: the skills' command is unchanged and the head tells the
agent how to fetch the rest, so nothing new seems needed of the binary;
record the reasoning either way.

## Acceptance

1. `grove guide work --entrypoint 2` prints a head of at most 13 KB that
   ends with the table of remaining parts; each `--part NAME` prints that
   part alone; an unknown part is a usage error (exit 2) naming the parts;
   `--part all` prints the whole file. A test asserts that the head and the
   parts, in order, equal `docs/work-execution.md`.
2. No instruction an agent must follow exists only in the head's table: each
   step's text is in exactly one part, and the head says when each part is
   read.
3. The shipped-document checks pass unchanged, and the entrypoint revision
   decision is recorded in this record with its reason.
4. One real `grove run` attempt on this repository after the change shows,
   in its events, the head printed once and each later part printed at most
   once, at the step that needs it. Evidence lists those `guide` commands
   with the turn each ran at, and the attempt's cost, turns and duration
   (G-260927-ngkbz's attempt, with the guide printed whole: 46 turns, $2.99,
   10 minutes).
5. `docs/commands.md` "Version and guide", the command list in `README.md`
   and this repository's `CLAUDE.md` line on `grove guide` say what `guide`
   prints now; the evals README's wording on the `grove guide shape`
   retrieval fact is untouched.

## Evidence

On `worktree-G-260927-cg6rt`, from base `main` `38511c1` and record revision
`sha256:72d8fe53923a`. No plan was written: the proposed design above served
as one. Implemented in `e50b955`, with the review fixes in `75675e7`. The
candidate is the commit that records this Evidence.

- **Acceptance 1.** `docs/work-execution.md` gains "The rest of this guide"
  after step 3: `grove guide work` stops there, each later part is printed
  once when reached, and a table names each part, the step it holds and when
  to print it. `grove.WorkParts` (`guides.go`) is the fixed table: each part
  starts at its `##` heading and runs to the next; `grove.WorkGuide` cuts the
  file. `grove guide work [--entrypoint N]` prints the head, 12,995 bytes;
  `--part NAME` one part; `--part all` the file. An unknown part is a usage
  error, exit 2, listing the parts; `--part` on another guide or command is
  exit 2 too. `TestGuideWorkParts` (`internal/cli/init_test.go`) asserts
  the head is at most 13,000 bytes and ends with the table naming every
  part, each part starts at its heading, the head and the parts in order
  and `--part all` equal the embedded file byte for byte, and the errors.
- **Acceptance 2.** Every `##` section is in exactly one part (the
  byte-for-byte test); "When a human decision is missing" travels with step
  7 in `checkpoint`. The table only says when to print each part, and the
  Read in stages table gains one row, "Reaching a later step", pointing at
  it rather than a row per part, so the part list has one owner.
- **Acceptance 3.** The shipped-document checks pass with their assertions
  unchanged; their invocation for `work` is now `guide work --part all`
  (`wholeGuide`), since bare `guide work` would give them only the head.
  **Entrypoint revision stays 2.** G-260925-p2k54 and the Entrypoint
  revision term change it only when an entrypoint needs something of the
  binary an older one lacks, or a newer one stops serving what it asks. The
  entrypoints still ask only for `grove guide work --entrypoint 2` and follow
  what it prints. An older binary prints the whole guide, which mentions no
  `--part`; a newer one prints a head whose `--part` the same `grove` serves.
  So nothing new is needed of the binary, and the skills' wording ("follow
  the guide it prints") stays true.
- **Acceptance 5.** `docs/commands.md` "Version and guide", README.md's
  workflow paragraph and command table, and CLAUDE.md's line on the two
  workflows say what `guide work` prints now. `docs/work-review.md` names
  step 6 by `grove guide work --part review`. `evals/README.md` is
  untouched.
- **Head size.** Intro through step 3 was 12,184 bytes; with the table it
  was 13,282, so the intro's list of what the workflow covers became "from
  an assignment to its handoff", and the table's cells were shortened. No
  instruction was dropped.
- **Final checks.** At `75675e7`: `go vet ./...` passed; `gofmt -l .`
  printed nothing; `go run ./cmd/grove check` printed `OK: 216 records`;
  `go test -count=1 -timeout 120s ./...` passed every package. A script
  resolved every relative link and anchor in README.md, CLAUDE.md,
  docs/commands.md, docs/work-execution.md and docs/work-review.md. There is
  no TUI change, so the terminal checks were not run.
- **Review.** [G-260928-f0ayt](G-260928-f0ayt-review-of-g-260927-cg6rt.md):
  two `grove-reviewer` rounds. Round 1 (`e50b955`) found two minor issues,
  a `--part` error hiding an unknown guide name and a sentence in
  docs/commands.md that lost its referent, both fixed in `75675e7`. Round 2
  ended `Open findings: none`.
- **Limits.** Acceptance 4 is not met: this headless session did not
  launch a `grove run` attempt, which is a paid session it may not start.
  It needs the change where an attempt's session runs its `grove`: after
  integration, or an attempt in a worktree holding this branch.

## Next

In review at the candidate named in the frontmatter. The owner judges
acceptance 1 to 3 and 5 from the Evidence and decides acceptance 4, a real
attempt, before or after integration. The attempt to compare against is
G-260927-ngkbz's: 46 turns, $2.99, 10 minutes, with the guide printed whole.

1. Approve in this checkout:
   `grove approve G-260927-cg6rt "VERDICT"`.
2. Integrate in the target's checkout:
   `grove integrate G-260927-cg6rt --cleanup`.
3. For acceptance 4, after integration: `grove run WORK_ID` on a small
   proposed work item, then read its events for each `guide` command and
   the turn it ran at, and record them here with the attempt's cost, turns
   and duration.
