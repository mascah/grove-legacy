---
id: "G-260923-gsthp"
type: work
title: "Prepare Grove for external distribution"
status: proposed
created: "2026-09-23T16:05:09Z"
updated: "2026-09-26T20:02:47Z"
relates_to: ["G-260921-5gz9a", "G-260922-jtsed", "G-260925-358a2"]
depends_on: ["G-260925-3pj9a", "G-260925-ej1xh"]
size: small
---

## Outcome

Anyone with Go installed, the owner first, installs Grove with one
`go install` line, and `grove version` names a release rather than a
pseudo-version or `(devel)`. Owner intent, shaping conversation 2026-09-26:
Grove is a tool for the owner alone right now; do the bare minimum that
leaves distribution in a clean state, and no Homebrew tap, release pipeline
or site.

That narrows the intent of 2026-09-23 (distribution artifacts, GitHub Pages,
release-please and GitHub configuration for a small external preview). The
brief's preview audience stands; what this work delivers toward it is the
install path, and the rest is listed below as not selected, not rejected.

## Constraints

Observed at `main` `2bd50b9`, 2026-09-26:

- [G-260925-358a2](G-260925-358a2-give-every-distributed-b.md) delivered build
  identity, and [Version and guide](../docs/commands.md#version-and-guide) is
  the contract: `version` prints the stamped release, else the module version
  Go recorded. No tag exists (`git tag --list` is empty), so every build
  prints `(devel)` or a pseudo-version.
- In a disposable clone tagged `v0.1.0` at `2bd50b9`, an unstamped
  `go build ./cmd/grove` printed `grove v0.1.0 (2bd50b9271ea…) guides
  sha256:1828d72e34c0 content sha256:07bd3cf86419`; one commit later,
  `grove v0.1.1-0.20260926200057-239604f76daf (239604f7…)`. Go 1.26.8
  derives the module version from the tag, so a tag alone makes `version`
  meaningful for `go build` and for `go install …@v0.1.0`, which the contract
  says records no commit.
- The README's install block says `go install
  github.com/mascah/grove/cmd/grove@COMMIT`, "pushed commits only". With a
  tag, `@v0.1.0` and `@latest` resolve through the Go module proxy, which
  also verifies the module checksum; a Go user needs no artifact, checksum
  file or platform matrix.
- GitHub (`gh repo view`, 2026-09-26): public; description, homepage and
  topics empty; no license, releases or tags; issues, wiki and projects on,
  discussions off. The tree has no release workflow, Pages source or LICENSE.
- The license stays as the owner chose on 2026-09-22
  ([G-260922-jtsed](G-260922-jtsed-run-secure-ci-and-depend.md)): unlicensed,
  an open item and not a question record. A public repository without a
  license grants nobody else the right to use the code; that matters only
  when someone else is meant to, and adding a LICENSE file is one commit then.

Both prerequisites are done and merged, and are why an installed `grove`
works in a project without this checkout present:
[G-260925-ej1xh](G-260925-ej1xh-strip-grove-repository-p.md) (the shipped
documents point at nothing this repository holds) and
[G-260925-3pj9a](G-260925-3pj9a-launch-attempts-only-whe.md) (`run` refuses
uncommitted entrypoints instead of spending on them).

In scope, proposed:

1. The README's install block says `go install
   github.com/mascah/grove/cmd/grove@latest`, `@v0.1.0` to pin, and that
   `grove version` prints `grove v0.1.0 …`; the command reference's "`go
   install …@COMMIT` resolves only a pushed commit" sentence admits tags.
   Upgrade is the same line again, then `init --check` and `init`, as the
   README already says.
2. Tag `v0.1.0` on `main` after that commit, and push the tag. This is the
   owner's act: a tag on the public repository is publication. `v0.1.0` is
   proposed because the build contract's example uses it and it promises
   nothing about stability.
3. A one-line GitHub description, set by the owner. Proposed: "Local project
   workspace for humans and agents: work, decisions and knowledge as Markdown
   in Git, with a CLI and terminal board." Topics and homepage optional.

Not selected now, on the owner's direction of 2026-09-26, each with the
condition that would bring it back: prebuilt binaries, checksums and a
platform matrix (add when a non-Go user asks for a binary); release-please,
GitHub Releases and a changelog (a tag is the release and `git log
v0.1.0..` is the changelog); GitHub Pages, TUI images, troubleshooting and a
feedback route (the README is the onboarding); a Homebrew tap; signing and
notarization (`go install` builds locally, so nothing is signed). Live
GitHub settings other than the description stay as they are.

## Acceptance

1. The README's install block and `docs/commands.md` say the `go install`
   line and what `version` prints, and nothing in either claims a binary,
   checksum, release page or site exists.
2. `v0.1.0` is on `main` and on `origin`, at a commit that holds 1.
3. From a `GOBIN` outside this checkout, `go install
   github.com/mascah/grove/cmd/grove@v0.1.0` installs, `grove version` prints
   `grove v0.1.0` with the same `guides` and `content` digests as a `go build`
   at the tag, and `grove init` then `grove check` work in a disposable Git
   repository with that binary. Rehearsed so far only with a local tag in a
   disposable clone; the install through the proxy is exercised after the
   push.
4. The GitHub description is set, or the owner declined it; either is
   recorded here.

## Next

Small enough for the owner to do by hand from the main checkout with a clean
tree, in this order, or to assign with `/grove-work G-260923-gsthp` for step
one and take the rest after review:

```sh
# 1. edit README.md and docs/commands.md (acceptance 1), commit
# 2. publish the tag; the module proxy can lag a few minutes for @latest
git tag -a v0.1.0 -m "v0.1.0" && git push origin v0.1.0
# 3. exercise the install away from this checkout
GOBIN="$(mktemp -d)" go install github.com/mascah/grove/cmd/grove@v0.1.0 && "$(go env GOBIN)"/grove version
# 4. description
gh repo edit mascah/grove --description "Local project workspace for humans and agents: work, decisions and knowledge as Markdown in Git, with a CLI and terminal board."
```

Record the observed `version` line and the description under Evidence, then
`done` on `main`. Nothing here is published until the owner runs step 2.
