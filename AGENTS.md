# AGENTS.md

## What this repo is

An Agent Package Manager (APM) marketplace of AI agent context files (skills,
prompts, agents, instructions). Content-only — no application code to build or
test.

## Prerequisites

All commands below require the `apm` CLI, `go`, and `just`. CI installs them via
`microsoft/apm-action@v1`, `actions/setup-go@v6`, and
`extractions/setup-just@v4`.

## Package structure

Each package lives at `packages/<name>/` and requires:

- `apm.yml` — defines `name`, `version`, `type`, `includes: auto`
- `.apm/` — directories for primitives. The `type` field determines the expected
  structure:
  - `skill` → `.apm/skills/`
  - `instructions` → `.apm/instructions/`
  - `hybrid` → any combination of `skills/`, `prompts/`, `agents/`,
    `instructions/`
- Optional `README.md`

## Marketplace-level manifest

Root `apm.yml` declares all packages with `subdir` and `version` constraints.
The `version` field here is the **marketplace version** — independent from
per-package versions. Per-package versions must satisfy the semver constraint
declared in the root manifest.

## Verification commands

Run before committing package changes:

```bash
just fmt                    # auto-format all YAML files (apm.yml + packages/*/apm.yml)
just pre-commit-check        # syncs versions, validates marketplace, checks for staleness
```

CI equivalent (read-only, exits non-zero on mismatch):

```bash
just ci                      # fmt-check → sync-check → yaml-validate → apm marketplace check → apm pack → stale check
```

Individual steps:

```bash
just fmt-check              # verify YAML files are formatted (yamlfmt -lint)
just sync-check              # check package versions satisfy root constraints
just sync-fix                # auto-bump mismatched package versions
just yaml-validate           # validate apm.yml files against JSON Schema (schemas/apm.json)
apm marketplace check        # validates all refs resolve
apm pack --dry-run            # validates marketplace.json generation
```

## YAML formatting

All `apm.yml` files (root + package-level) are formatted with `yamlfmt`
(`.yamlfmt` config at repo root). 2-space indent, trailing newline, trimmed
whitespace. Run `just fmt` to auto-format or `just fmt-check` to verify in CI.

## Generated file

`.claude-plugin/marketplace.json` is **auto-generated** by `apm pack`. After
changing any package or the root `apm.yml`, regenerate it:

```bash
apm pack
```

Then **commit the updated marketplace.json**. The CI workflow will reject PRs
with mismatched versions.

## Release process

1. Update root `apm.yml` `version` field
2. Run `just sync-fix` to align per-package versions with root constraints
3. Run `apm pack` to regenerate `.claude-plugin/marketplace.json`
4. Commit, tag as `v{version}` (e.g. `v0.3.0`)

## Ignored directories

`apm_modules/`, `.apm_cache/`, `build/`, `*.tar.gz`, `handoffs/`, `to-import/` —
all gitignored. Never commit content from these.
