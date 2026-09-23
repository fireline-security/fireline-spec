# Contributing to fireline-spec

## Dev setup

Requires Go 1.26+. No external services needed.

```sh
task build
task validate
```

## Proposing a schema change

1. Open an issue or PR describing the change and why the current contract
   can't express what you need.
2. A non-breaking addition (a new optional field) can land in `v1` directly,
   with a `valid/` fixture exercising it.
3. A breaking change (removing or narrowing a field, changing a type) needs
   a new `schemas/v2/` tree (see the README's "Layout" section), so
   existing adapters and fireline-core don't break silently.
4. Every new constraint needs both a `valid/` fixture showing correct usage
   and an `invalid/` fixture showing what violates it.

## Code style

`gofmt`, `go vet`, and `golangci-lint` cover `tools/validate`; run `task lint`.

## Commit messages

Describe the *why*, not just the *what*: the diff already shows what
changed.

## Code comments

`docs/` is the source of truth for design rationale, trade-offs, and scope
decisions (use AGENTS.md/README.md for a repo with nothing in `docs/` yet).
A code comment should point to it, not restate it: `// see docs/<file>.md
for why.`

A doc comment says what a function or type does, in one to three lines.
Add local *why* only for something the code can't already show: an edge
case, a library quirk. Never a design decision or scope caveat; that
belongs in docs/. Keep it plain: no "X, not Y" framing, no "deliberately,"
no hedges like "it's worth noting," and no em dashes.
