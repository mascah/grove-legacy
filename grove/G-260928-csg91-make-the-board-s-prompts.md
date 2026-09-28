---
id: "G-260928-csg91"
type: work
title: "Make the board's prompts and refusals legible"
status: proposed
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T19:34:00Z"
kind: fix
size: small
relates_to: ["G-260924-ecs9m", "G-260921-jwk4e", "G-260921-7trd7", "G-260919-k7b8j"]
---

## Outcome

## Constraints

## Acceptance

## Next

## Outcome

When the owner types on a board prompt they see everything they typed, and a
refused launch or action is unmistakable and stays on screen until they
dismiss it, so a launch is never repeated in the belief that it succeeded.

Owner intent, shaping conversation 2026-09-28: "I sometimes can't see what
I'm typing on the bottom input bar if the text is overflowing the width";
"I relaunched the same item 3 times without realizing it was telling me
there was a problem and it couldn't launch."

## Constraints

Observed at main `6fbb888`, 2026-09-28, in `internal/tui`:

- Every prompt is one footer row. `promptRow` (`review.go`, around 779-811)
  builds one string and `line` (`view.go`, around 110) cuts it at the width
  with `…`. The launch line puts the typed text first, so its help is lost
  first; `a` and `f` put the typed text last, after a preamble naming the
  branch and the group, so the typed text is lost first. There is no text
  input component: `promptKey` (`review.go`, around 671) appends runes and
  removes the last on backspace, with no cursor movement.
- A launch refused before its prompt opens (`launch`, `attempts.go`, around
  660-716: work in review, a blocking question, a live attempt, not in this
  checkout) sets `m.notice`, drawn as unstyled text in the banner row under
  the header, and `m.notice` is cleared on the next key (`model.go`, around
  660). A refusal on Enter (an unknown flag, `--branch`, a missing budget or
  mode) does the same with the line left open.
- A refusal from the runner (`attempt.Start`: a wait, the record changed, an
  attempt running) reaches the result screen through `actMsg` (`model.go`,
  around 627-651), whose only difference from a success is a bold
  `NOT DONE:` prefix (`resultRows`, `review.go`, around 813-829) in the same
  colour; the board is re-read either way, and Esc dismisses it.
- The board's text escaping and exact-width rows are constraints
  ([G-260919-k7b8j](G-260919-k7b8j-browse-a-terminal-kanban.md)); the fit
  tests assert every row is exactly the width.

Proposed design, labelled proposed:

- A prompt that takes typed text draws over as many footer rows as the text
  needs, wrapped with the existing `wrap`, typed text first and help last,
  so nothing typed is lost; the y/n prompts stay one row. No modal and no
  text-input component: the footer grows.
- A refusal is drawn in an attention style no column uses, stays until Esc
  or until the action succeeds, and names the next action; the notice's
  clear-on-any-key goes. `NOT DONE` takes the same style on the result
  screen, reason first.
- A launch refused because an attempt is live names that attempt and `A`.

Out of scope: text editing beyond append and backspace; a web or modal form.
The owner deferred a web UI on 2026-09-28 until the board fixes land.

## Acceptance

1. On an 80-column terminal, a verdict, a feedback text or a launch line
   longer than the width shows every typed character, and Enter records
   exactly what was typed.
2. Each refusal path (before the prompt, on Enter, from the runner) is shown
   in the attention style and stays until Esc; a test covers one of each,
   and the pseudo-terminal script covers the live-attempt refusal.
3. No raw escape reaches the terminal and every row stays exactly the
   width, with the existing fit tests extended to the grown footer.
4. `docs/board.md` says what a prompt and a refusal look like; the owner
   judges both in a terminal.

## Next

Assign: `/grove-work G-260928-csg91`. Touches `internal/tui` in places
separable from [G-260928-r1hkh](G-260928-r1hkh-show-approval-and-merge.md),
[G-260928-y50a4](G-260928-y50a4-edit-any-record-in-your.md) and
[G-260928-63124](G-260928-63124-show-each-member-s-state.md), which needs no
edge; the owner chooses whether to select them together.
