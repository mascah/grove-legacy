# Grove

Grove gives you and your agents one loop for a project's work: shape an idea
into proposed work, hand it to an agent, judge the candidate commit it hands
back, and merge it. The work, questions, decisions, plans and reviews along
the way are Markdown files in your own Git repository, read and changed
through a CLI and a terminal board, with no service to run. The board shows
each piece of work in its latest state across every local branch and
worktree, and the workflows the agents follow ship inside the binary.

## Install

Grove needs Git and Go 1.26.8 or later.

```sh
GOBIN="$HOME/.local/bin" go install github.com/mascah/grove/cmd/grove@latest   # or @v0.1.0 to pin
grove version   # grove v0.1.0 (…) guides sha256:… content sha256:…
```

Put that directory on `PATH`. If `grove version` prints a usage error,
another `grove` answered: fix `PATH` before `init`.

## Set up a project

At the top of the project's Git checkout (or with `--project /absolute/path`):

```sh
grove init
```

`init` writes `grove.yaml`, a `grove/` folder for records with a placeholder
brief, the `grove-work` and `grove-shape` skills for Claude Code (`.claude/`)
and Codex (`.agents/`), and Claude Code's `grove-reviewer` agent. It never
touches `AGENTS.md` or `CLAUDE.md`: the skills defer to them, so say there if
`grove` is not the one on `PATH`, or how work branches should be named.
[Init](docs/commands.md#init) has the details, including Codex's `PATH`.

Then name the branch work is delivered to, so `grove.yaml` reads:

```yaml
schema_version: 4
records: grove
brief: grove/brief.md
target: main
```

- `target` is that branch: `integrate` delivers accepted work onto it as
  one squash commit, done is read from what it contains, `resolve` and
  `sweep` need it, and the board marks work not yet on it. A project at
  schema 3 converts with `grove migrate` ([Migrate](docs/commands.md#migrate)).
- `run:` sets `grove run`'s defaults (`budget`, `permission_mode`, `model`,
  `effort`); without it, every launch passes `--budget` and
  `--permission-mode`.
- `policy:` lets `grove sweep`, which every attempt that hands work off
  runs when it ends, resolve, approve and integrate candidates without you.
  Leave it out until you want automatic acts.

`grove guide model` prints the rules for each key, under "Configuration and
discovery", with the record types, fields and lifecycle. Check and commit
before the first attempt, since an attempt's worktree holds only committed
files:

```sh
grove check
git add grove.yaml grove .claude .agents && git commit -m "chore: set up grove"
```

To upgrade, run the same `go install` line, then `grove init --check` and
`grove init` in each project, and commit what changed
([Entrypoint revisions](docs/commands.md#entrypoint-revisions)).

## The loop

Work runs `proposed`, `active`, `review`, `done`.

1. **Brief.** `grove/brief.md` says what the project is for and the direction
   you have chosen. Write it, or let shaping develop it.
2. **Shape.** `/grove-shape TOPIC` in Claude Code (`$grove-shape TOPIC` in
   Codex) works an idea through with you and proposes work records, each with
   an outcome and acceptance. Shaping implements nothing.
3. **Work.** `/grove-work ID` (`$grove-work ID`) carries out one proposed item
   on its own branch and worktree, has it reviewed, and hands a candidate
   commit into `review`. `grove run ID --budget USD --permission-mode MODE`,
   or `R` on the board, does the same as a headless attempt that outlives the
   terminal.
4. **Judge.** Open the work on the board (`grove`) and read the candidate:
   `a` accepts it, `f` gives feedback and returns it to `active`, and `i`
   delivers the accepted candidate onto the target as one squash commit,
   after which it reads as done. The same are `grove approve`, `grove
   feedback` and `grove integrate`.

`grove guide shape`, `grove guide work` and `grove guide review` print the
workflows the skills load ([shaping](docs/work-shaping.md),
[work](docs/work-execution.md), [review](docs/work-review.md)); `guide work`
prints the work guide's head, and `--part NAME` each later step as it is
reached.

## Commands

`grove` with no command opens the board: work by status, in its current
state across every branch and checkout. From it you open a record's detail,
history and versions, search every record, judge and integrate a candidate,
and start and stop attempts ([The board](docs/board.md)).

`grove --help` has each command's usage. Every command takes
`--project DIR`; without it, Grove searches upward for `grove.yaml`.

| Command | What it does | Details |
| --- | --- | --- |
| (none) | Opens the terminal board | [Board](docs/board.md) |
| `list`, `show`, `brief`, `check` | List, print and validate records and the brief, reading only | [Records](docs/commands.md#records), `grove guide model` |
| `new`, `update` | Create a record with a new date-form ID; change its frontmatter | [Records](docs/commands.md#records), `grove guide model` |
| `convert` | Make a record from a Markdown document outside the record root | [Records](docs/commands.md#records), `grove guide model` |
| `approve`, `feedback`, `integrate` | Judge a candidate in review; deliver an accepted one as one squash commit | [Judging and integrating](docs/commands.md#judging-and-integrating), `grove guide model` |
| `resolve` | Give a candidate that conflicts with the target to one attempt that merges the target and resolves it | [Resolving a conflict](docs/commands.md#resolving-a-conflict) |
| `sweep` | Resolve, approve and integrate candidates in review under the owner's standing `policy:`, each act attributed to it | [Sweep](docs/commands.md#sweep) |
| `run`, `attempts`, `attempt`, `stop` | Start one headless agent attempt of a work item that outlives the terminal; list, inspect and stop attempts | [Attempts](docs/commands.md#attempts) |
| `versions` | Show each record's versions across branches and worktrees, current or older | [Versions](docs/commands.md#versions) |
| `workspace` | Print the checkout holding a selected version | [Workspace](docs/commands.md#workspace) |
| `context` | Assemble staged context for selected work | [Context](docs/commands.md#context), `grove guide work` |
| `deps` | Show how unfinished work depends on other work, or preview a selection's order | [Dependencies](docs/commands.md#dependencies) |
| `init` | Set up Grove in a Git checkout | [Init](docs/commands.md#init) |
| `migrate` | Convert a schema 3 project to schema 4, previewed first | [Migrate](docs/commands.md#migrate) |
| `guide`, `version` | Print a workflow or review guide (the work guide's head, or one `--part`) or the record model; name this build | [Version and guide](docs/commands.md#version-and-guide) |

## Develop Grove

This repository tracks its own work with Grove: `grove/` holds its real
records and [brief](grove/brief.md). In a checkout, `go run ./cmd/grove …`
runs that checkout's code (`go run ./cmd/grove` opens the board over this
repository's records), and `just install` replaces `~/.local/bin/grove` with
a build of it. [CLAUDE.md](CLAUDE.md) holds the development and verification
policy, the hooks, CI and tooling, and where each subject is owned for work
here; `AGENTS.md` is a symlink to it, so Codex reads the same file.
[evals/README.md](evals/README.md) measures how well the guides steer an
agent.
