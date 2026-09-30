---
id: "G-260929-04svs"
type: work
title: "Enter and resume interactive work on the selected harness"
status: proposed
created: "2026-09-29T03:04:23Z"
updated: "2026-09-30T01:16:08Z"
kind: feature
size: medium
relates_to: ["G-260928-y50a4", "G-260923-tnn5e", "G-260928-kehya", "G-260928-y2p5h", "G-260919-k7b8j", "G-260921-sth8q", "G-260930-e8jj7", "G-260930-60c3d", "G-260930-r2k4g"]
depends_on: ["G-260930-gwnb1"]
---

## Outcome

From the work they are inspecting, a person can enter a supported
interactive harness with the correct project context, help or continue the
work, and return to the same Grove view. Native resume is offered where
supported; a fresh session can use the durable handoff when it is not.

The original owner request on 2026-09-28 was to open Claude from a record
with that record in scope, type an instruction, return to Grove on exit,
and resume later. The owner selected portable adoption on 2026-09-29 in
[G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md); this proposal now serves both milestone harnesses.

## Scope and constraints

Use the project's selected supported harness and expose capability
differences. Do not maintain a second Claude-only provider selector or
launch path. Bind context and any native session to the actual work and
checkout, rechecking freshness before launch. Keep local session identifiers
under the Git common directory; durable work must not require them.

An interactive session has a present person and does not acquire the
authority of a bounded unattended assignment merely because it was opened
from Grove. Opening a record does not silently assign implementation, send
an expensive initial turn, or grant integration permission. Explain any
provider behavior that requires submitting context as a turn.

Do not attach a second writer to an actively owned worktree. Use the
accepted design's stop/wait/isolation path, and preserve unfinished work.
The board suspends and restores the terminal, then refreshes the same work
and explains the outcome. Exact key bindings and prompt seeding are design
choices rather than this record's acceptance.

## Observed evidence

Observed 2026-09-28 at main `9a18f57`:

- `e` ([G-260928-y50a4](G-260928-y50a4-edit-any-record-in-your.md)) is the
  pattern: `tea.ExecProcess` of the editor on the board's own stdin and
  stdout (`liveAnswer`, `internal/tui/answer.go`), in the one checkout on
  the branch of the shown version (`checkoutOf`, `internal/tui/review.go`;
  `projectDir` is its root and the version's `Path` the record's path in
  it), refused on a commit's copy, a diff, an older state, a deleted record,
  and no or two checkouts on the branch; after the editor, `edited` re-reads
  the board. The detail key switch (`internal/tui/detail.go`) leaves `c`
  unbound; `c` is bound only in the dependencies screen.
- Claude Code 2.1.284 (`claude --help`; the CLI reference at
  https://code.claude.com/docs/en/cli-reference) has no flag, setting, hook
  or environment variable that pre-fills the interactive input without
  submitting it: the positional prompt is sent as the first turn.
  `--append-system-prompt` adds to the system prompt without a turn;
  `--name` titles the session in the prompt box and the `/resume` picker;
  `--session-id UUID` sets the id; `--resume UUID` resumes it.
- Resume, probed on this machine with that version: `claude -p --session-id
  UUID --name NAME "…"`, then interactive `claude --resume UUID` in the
  same directory, opened straight into that conversation. Interactive
  `claude --resume NAME` did not: it opened the picker with NAME as a search
  term and said "No sessions match", even for a session created there with
  that `--name`; print-mode `--resume NAME` resolves a title only in the
  directory the session was created in. An unknown UUID prints "No
  conversation found with session ID: …". A reliable resume therefore needs
  the id Grove chose, kept by Grove, per checkout.
- Attempts already do that for headless runs: `uuid()` and `--session-id`
  (`internal/attempt/attempt.go`, around 551 and 589), `session_id` in
  `attempt.json`, state under the Git common directory (`Dir`, around 305:
  `.git/grove/attempts/`, shared by every worktree, never committed), and
  `environ(os.Environ(), OwnerEnv, "CLAUDE*", "!CLAUDE_CONFIG_DIR", …)`
  (around 861) so the provider is not a child of a launching Claude session
  ([G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)). The probe's
  resume, run from inside a Claude session, showed "Transcript saving is
  off — inherited CLAUDE_CODE_CHILD_SESSION marker": that nesting.
- `GROVE_CLAUDE` names the provider executable (`provider()`, around 733);
  `internal/tui/testdata/terminal.py` runs the board with a fake provider
  and a fake editor that checks it has a terminal (`edit_record`, around
  783).
- Bare `grove` opens the board and explicit subcommands stay noninteractive
  ([G-260919-k7b8j](G-260919-k7b8j-browse-a-terminal-kanban.md)); this is
  board-only, with no subcommand.

The earlier proposal used key c, a Claude system-prompt seed and a local
session map. These remain design evidence, not a requirement to duplicate
provider/session handling. Recheck the dated native-resume observations
for each supported version.

## Dependencies

Depends on [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md): it consumes the exercised durable handoff and
provider/session boundaries, including safe recovery. This replaces the
earlier assumption that a separate Claude-only session store could be built
without ordering it against the runner changes.

## Acceptance

1. A user enters either supported interactive harness from live work with
   its current mandate/checkpoint and relevant context, then returns to the
   same refreshed Grove view. Commit-only, stale, ambiguous or actively
   owned sources give an actionable wait or supported safe alternative.
2. Native resume, fresh start and unavailable/missing sessions have clear
   outcomes. Provider/checkout mismatch never silently resumes another
   conversation; fresh continuation uses the durable handoff.
3. Opening or inspecting work does not itself authorize an unattended
   implementation or merge. Any automatically submitted turn is explicit.
4. Local session bookkeeping is not committed or needed by a fresh clone.
   Provider nesting/session variables are handled through the common
   provider boundary without leaking one provider's session into another.
5. Terminal lifecycle tests exercise first launch, resume, failure and
   return. The owner tries both real harnesses in a terminal under an
   explicit mandate, and documentation names capability limitations.

## Next

Needs [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md). Implement the reviewed interactive transition against
its handoff and provider boundaries. The final assembled experience belongs
to [G-260930-r2k4g](G-260930-r2k4g-make-the-complete-grove.md).
