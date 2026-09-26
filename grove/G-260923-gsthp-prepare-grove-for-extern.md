---
id: "G-260923-gsthp"
type: work
title: "Prepare Grove for external distribution"
status: proposed
created: "2026-09-23T16:05:09Z"
updated: "2026-09-25T19:08:16Z"
relates_to: ["G-260921-5gz9a", "G-260922-jtsed", "G-260923-fwakw", "G-260923-p5pt6", "G-260923-895zb"]
depends_on: ["G-260925-3pj9a", "G-260925-ej1xh"]
---

## Outcome

Prepare reproducible Grove distribution and documentation that someone outside the
owner's projects can install and exercise without the owner reconstructing
the project's history or supplying unwritten setup knowledge.

Owner intent, shaping conversation 2026-09-23: prepare distribution, public
documentation/GitHub Pages, release-please, and GitHub project configuration.
The owner selected the assistant's option "Prepare for a small external preview"
and then asked what "preview" meant. The assistant clarified that it meant
initial use outside the owner's projects. It is not a selected release channel,
launch commitment or support policy. Specific release scope, platform coverage,
licensing and compatibility promises remain later choices.

## Constraints

Observed at main `f27444e`: [G-260921-5gz9a](G-260921-5gz9a-bootstrap-projects-with.md) delivered
`init`, versioned embedded guides and thin adapters; the documented distribution
is a build or Go installation from a named commit.
[CI](../.github/workflows/ci.yml) already checks Linux and macOS and
[Dependabot](../.github/dependabot.yml) covers Go and actions. No release
workflow, Pages source or LICENSE was found in the inspected tracked surfaces.
[G-260922-jtsed](G-260922-jtsed-run-secure-ci-and-depend.md) records the owner's earlier choice to leave licensing
open and keep CI advisory. Live GitHub settings were not rechecked here;
historical settings in that record are not present-state evidence.

Proposed scope:

- Reproducible versioned artifacts for a declared platform matrix, with
  checksums, meaningful `grove version`/guide identity, and install, update and
  recovery instructions exercised outside this development checkout.
- Release-please configuration and a reviewed build/publication workflow.
  [Release Please](https://github.com/googleapis/release-please) supplies release
  PRs, version changes, changelogs and GitHub releases; artifact packaging needs
  its own steps. Keep candidate-preserving local integration intact; release
  tooling must not silently require a different product merge policy.
- Public onboarding, a compact end-to-end example, TUI images, troubleshooting,
  current limitations and a feedback route. Publish maintained documentation
  through GitHub Pages without creating a second editable copy of the guides.
- Reviewable GitHub description/topics/homepage and relevant repository settings,
  based on a fresh inventory. Do not recreate existing CI or add team-oriented
  automation without a concrete preview need.

[G-260923-fwakw](G-260923-fwakw-reconcile-current-docume.md) owns the existing-document reconciliation
and information ownership. This work owns distribution-specific documentation,
site delivery and the assembled newcomer experience; coordinate that boundary.
[G-260923-p5pt6](G-260923-p5pt6-establish-behavioral-eva.md) and [G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) provide
related quality evidence. Their completion is not automatically a prerequisite
or sufficient evidence that a preview is ready; the owner selects release scope.

This proposal prepares distribution artifacts and publication/configuration changes
for review. It does not authorize publishing, changing live GitHub settings or
installing into sibling projects. No broad launch, hosted service, automatic
self-update or expansion to unsupported execution providers is selected.

Owner intent, review conversation 2026-09-25: Grove is used in other
projects as an installed CLI without this repository present. That review
found the binary nearly self-contained, fixed the last path links on `main`
(`001b271`), and shaped two prerequisites of this work's acceptance 2 and
4: [G-260925-3pj9a](G-260925-3pj9a-launch-attempts-only-whe.md) (`run` and `R` depend
silently on `init`'s files being committed) and
[G-260925-ej1xh](G-260925-ej1xh-strip-grove-repository-p.md) (the shipped documents
still name Grove's own records and unshipped documents).

## Acceptance

1. A concrete release proposal identifies intended users, platform coverage,
   versioning, license disposition and upgrade expectations, distinguishing
   selected choices from recommendations and matters awaiting the owner. The
   earlier pre-release policy is not silently converted into a support promise.
2. A clean checkout of the intended release revision produces attributable
   artifacts for the selected matrix; installation and `version`, `init`,
   `check`, guide discovery and upgrade behavior are exercised in disposable
   environments. State platform and signing/notarization limits honestly.
3. Release-please, artifact production and site build have reviewable
   configuration and validation evidence. Version, changelog, executable and
   embedded-guide identity agree. Remote publication behavior that cannot be
   rehearsed without publishing is documented as unverified.
4. Using only preview documentation, a fresh user or fresh-session surrogate
   can initialize a disposable repository, shape work, execute a bounded
   assignment, inspect/review the result, and understand an interrupted attempt.
   Distinguish surrogate evidence from a real external user's judgment. Any
   live agent trial has explicit resource bounds. The owner judges readiness.
5. A handoff identifies exact artifacts and revisions, documentation/site output,
   proposed GitHub settings, observed limits, and concrete publication steps.
   Nothing is presented as publicly released or configured before it is.

## Next

The owner can review this proposal now. Prepare a concrete packaging/site plan
and inspect current GitHub settings read-only. Propose the initial platform
matrix, and surface owner choices when implementation or publication would
commit to them. Release compatibility does not need deciding during this shaping
session or before the documentation, eval and UI work. Commit the agreed proposal before
assigning it through `$grove-work G-260923-gsthp`. Publication remains a later explicit
action on the prepared candidate.
