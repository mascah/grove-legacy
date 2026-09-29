---
id: "G-260929-04svs"
type: work
title: "Open a record in an interactive Claude session from the board, and resume it"
status: proposed
created: "2026-09-29T03:04:23Z"
updated: "2026-09-29T03:05:42Z"
kind: feature
size: small
relates_to: ["G-260928-y50a4", "G-260923-tnn5e", "G-260928-kehya", "G-260928-y2p5h", "G-260919-k7b8j", "G-260921-sth8q"]
---

## Outcome

From a record's detail on the board, the owner opens an interactive Claude
Code session in the checkout that holds the record, with the record already
in the session's scope, and types what they want; leaving the session
returns the board to that same detail. A later press on the same record
resumes that session instead of starting fresh, and the owner may start a
new one instead.

Owner intent, shaping conversation 2026-09-28: "launch an interactive claude
session from grove with the currently open item pre-populated in the input
field like `@grove/G-260929-z5nec-implementation-writes-denied.md` and then
I can just start typing afterwards to explain. exiting the claude process
would return me to grove on the item that it was opened from"; "Being able
to resume the sessions would be really helpful."

## Constraints

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
- [G-260928-kehya](G-260928-kehya-resume-an-attempt-s-sess.md) resumes an
  attempt's print session with `--resume`, and
  [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) reshapes the
  attempt command for Codex. Both touch `internal/attempt`; this touches
  `internal/tui` and a new session store, so no order is declared.

Proposed design, labelled proposed:

- `c` on a record's detail: the `e` refusals, then `claude` in the checkout
  holding the shown version, on the board's terminal, with `--session-id` a
  UUID Grove generates, `--name ID`, and `--append-system-prompt` saying the
  session was opened from the Grove board on ID at PATH in this checkout,
  to read it before answering, and that the person will say what they want.
  Session variables scrubbed as attempts scrub them. The record is in scope
  without a spent turn. The visible, editable mention the owner described
  is not possible; the positional prompt `@PATH`, which costs one turn
  before the owner types, is the alternative seed.
- Grove keeps the checkout, session id and start time per record under the
  Git common directory beside attempts (`.git/grove/sessions/ID.json`, never
  committed). When one exists for this record in this checkout, `c`
  prompts: Enter resumes it (`--resume UUID`), showing when it started; `n`
  starts a new one and replaces the entry; Esc cancels. A resume Claude
  refuses ends like any exit and the notice says so; the owner starts new.
- After the session exits, the board re-reads as after `e`, and the notice
  names the session and that `c` resumes it.
- Not an [attempt](G-260921-sth8q-attempt.md): nothing under attempts, no
  attempt record, no budget, no permission mode, no result reconciled; the
  owner is in the session. Claude only; Codex's `codex resume` belongs to
  the provider work.
- `docs/board.md`: a section beside "Editing a record", the key table row,
  and where the session store lives.

Out of scope: a subcommand; forking a session (`--fork-session`); choosing
model or effort from the board; a provider choice.

## Acceptance

1. `c` on a live record's detail runs the provider in the checkout holding
   it with `--session-id`, `--name ID` and the seed; the board returns to
   that detail and re-reads; each `e` refusal refuses `c` too and writes
   nothing. A fake provider records its arguments; `terminal.py` covers a
   first launch and a resume end to end, as `edit_record` does.
2. A second `c` on the same record in the same checkout offers the resume
   and passes `--resume` with the kept id; `n` starts fresh and replaces the
   kept id; the store is under the Git common directory and never
   committed.
3. A board started inside a Claude session launches the provider with the
   session variables scrubbed, as attempts do.
4. `docs/board.md` says so; the owner judges in a terminal with the real
   `claude`, including one resume.

## Next

Assign: `/grove-work G-260929-04svs`. At assignment the owner may choose
the positional-prompt seed over the system-prompt seed, or another free key
over `c`; neither changes the acceptance.
