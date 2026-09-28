---
id: "G-260927-cg6rt"
type: work
title: "Print the work guide in stages: a head at start, later steps on demand"
status: proposed
created: "2026-09-27T22:13:35Z"
updated: "2026-09-27T22:28:16Z"
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

## Next

Assign: `/grove-work G-260927-cg6rt`. The attempt facts
[G-260927-dx0yn](G-260927-dx0yn-retain-per-attempt-proce.md) builds (guide
prints and `grove` calls, time to the first commit outside the record root)
are how this change's effect on a session will be read; its owner dropped
the process share and first edit it first proposed. This record does not
need them to be built.
