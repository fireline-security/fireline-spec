# AGENTS.md

## Purpose

This repo is the canonical, versioned source of Fireline's wire contracts.
Right now that's one contract: the Observation ("Smoke") a scanner adapter
emits. fireline-core and fireline-adapters each build against what's defined
here. They do not import this repo as a Go module (see their own AGENTS.md
files); they hand-copy the parts of the contract they need and document that
copy at the point it exists. This repo is the thing those copies get checked
against, not a library they depend on.

## What to avoid

- **Don't edit a shipped schema in place for a breaking change.** A field
  removal, a type change, a new required field: anything that would break
  an adapter or core build already targeting `v1` gets a new
  `schemas/v2/` (and `fixtures/v2/`) tree instead. Non-breaking additions
  (a new optional field) can land in `v1` directly.
- **Don't add or change a field without a fixture proving it.** Every
  constraint the schema encodes needs a `valid/` example that satisfies it
  and, where the constraint is worth guarding, an `invalid/` example that
  violates it. `task validate` is the check, not a suggestion.
- **Don't close off `identity_components`.** It's a free-form string map in
  `v1`, not a closed schema per source type: SAST, SCA, and DAST tools have
  genuinely different identity shapes, and over-constraining this now would
  force a schema revision per scanner. Don't "clean this up" into a fixed
  set of fields.
- **Don't add a normalized risk score field.** `severity_raw` stays exactly
  what the source reported. Collapsing severity, exploitability, or
  anything else into one opaque number is the thing this project is
  explicitly not doing (see the project pitch's D-04); this repo is not the
  place to start.
- **Don't couple this repo to fireline-core or fireline-adapters.** No Go
  `replace` directive pointing at a sibling checkout, no relative-path
  assumption about where those repos live on disk. They're separate
  modules; keep it that way.
- **Don't reach for a heavier validation toolchain without discussion.**
  `tools/validate` uses a Go JSON-Schema library on purpose, so this repo
  stays a self-contained Go module with no Node/npm dependency. If a
  non-Go contributor workflow genuinely needs `ajv-cli` or similar, that's
  a conversation, not a silent swap.

## What to ensure

- Every schema change ships with its fixtures in the same PR.
- `task validate` passes: `valid/` fixtures validate cleanly, `invalid/`
  fixtures fail for the reason their filename claims.
- `task lint` (`gofmt`, `go vet`) is clean.
- A breaking change gets a CHANGELOG.md entry explaining what changed and
  why existing `v1` consumers aren't affected.
