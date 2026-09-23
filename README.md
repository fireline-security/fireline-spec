# fireline-spec

The canonical, versioned source of Fireline's wire contracts. Right now
that's one contract: the Observation ("Smoke") a scanner adapter emits.
[fireline-adapters](https://github.com/fireline-security/fireline-adapters)
and [fireline-core](https://github.com/fireline-security/fireline-core) each
build against what's defined here.

## Layout

```
schemas/v1/observation.schema.json   the contract: JSON Schema, 2020-12 dialect
fixtures/v1/observation/
  valid/       fixtures that must pass
  invalid/     fixtures that must fail, and why
tools/validate/  checks every fixture against its schema
```

Versioning is directory-based: a breaking change to the Observation contract
gets a new `schemas/v2/` (and `fixtures/v2/`) tree rather than an edit to
`v1` in place, so nothing that already builds against `v1` breaks silently.

## Adding a fixture

Every field the schema newly requires or newly forbids needs a fixture
proving it: a `valid/` example showing the field used correctly, or an
`invalid/` one showing what happens when the constraint is violated. Run
`task validate` before opening a PR; CI runs the same check.

## Testing

```sh
task validate
```

This compiles the schema and checks every fixture under `fixtures/v1/`
against it: `valid/` fixtures must pass, `invalid/` fixtures must fail. No
external services, no Docker.
