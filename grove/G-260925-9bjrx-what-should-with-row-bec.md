---
id: "G-260925-9bjrx"
type: question
title: "What should G-260925-pbx81's with row become now that G-260925-dzxm6 shipped no agent-facing search?"
status: resolved
created: "2026-09-25T21:19:37Z"
updated: "2026-09-25T21:24:49Z"
blocks: ["G-260925-pbx81"]
relates_to: ["G-260925-dzxm6", "G-260925-gymkr", "G-260925-khwkq"]
---

## Question

[G-260925-pbx81](G-260925-pbx81-evaluate-whether-agents.md) acceptance 2 runs each case "on
its `with` row once G-260925-dzxm6's candidate is integrated", and acceptance 3
compares the rows. G-260925-dzxm6 is integrated on main (`1db5a75`, done at
`47852e3`), but its preparation narrowed it to the board's body search
and the review listing ([G-260925-dzxm6](G-260925-dzxm6-search-record-bodies-and.md), Scope at
preparation): no `grove search` command, no guide or reviewer sentence.
A headless shaping run cannot use the board, so a `with` row at main
would measure the guide edits of G-260925-3pj9a and G-260925-ej1xh that landed on main
meanwhile, not search. G-260925-dzxm6's own record says so ("G-260925-pbx81's `with` row
compares nothing new"), and the owner's verdict on its candidate was
"fine with deferring the grove search command for now". G-260925-gymkr's mandate
("Use the G-260924-frzeg settings") was given for the `without` row; its item 2
names the `with` row only for its run count.

This session did not spend, since which of these G-260925-pbx81 does is a scope
change to its acceptance, not a routine choice. Options:

1. **Drop the `with` row from G-260925-pbx81.** Amend acceptance 2 and 3 to the
   `without` row; the review completes on that row with the recomputed
   facts, and G-260925-pbx81 is handed into Review. A later search command or
   guide sentence, reshaped as its own work, carries its own `with` row
   on these cases, whose runner and `without` row stay on main. Cost:
   nothing further.
2. **Run the `with` row at main anyway**, under G-260925-gymkr's settings (about
   $3.35), as a check that main's guide edits since digest
   `41324c3655a1` neither break these cases nor add cost, then complete
   the review across both rows and hand off. It says nothing about
   search. The G-260923-p5pt6 pair is already the regression rerun for a guide
   edit (G-260923-h9c30 Evidence).
3. **Keep G-260925-pbx81 active**, waiting for a future work that ships `grove
   search` or the guide sentences, and run the `with` row then.

Recommendation: 1. G-260925-khwkq found no miss for search to recover at these
fixture sizes, so neither 2 nor 3 can show what G-260925-pbx81's outcome asks;
evidence for search needs the larger-fixture case G-260925-khwkq's Disposition
names, which is its own owner choice.

Who can answer: the owner. Answer below, set `status=resolved`, then
assign `/grove-work G-260925-pbx81` again on `worktree-G-260925-pbx81`.

## Next

Answered; G-260925-pbx81's outcome and acceptance 2 and 3 are amended to the
`without` row.

## Answer
Drop the with row
