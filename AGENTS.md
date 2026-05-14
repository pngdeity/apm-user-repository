# AGENTS.md

## What this repo is
An APM marketplace of AI agent context files (skills, prompts, agents, instructions) published as 16 packages. No code to build, test, or lint — content-only.

## Prerequisites
All commands below require the `apm` CLI. CI installs it via `microsoft/apm-action@v1` with `setup-only: 'true'`.

## Package structure
Each package lives at `packages/<name>/` and requires:
- `apm.yml` — defines `name`, `version`, `type`, `includes: auto`
- `.apm/` — directories for primitives. The `type` field determines the expected structure:
  - `skill` → `.apm/skills/`
  - `instructions` → `.apm/instructions/`
  - `hybrid` → any combination of `skills/`, `prompts/`, `agents/`, `instructions/`
- Optional `README.md`

## Marketplace-level manifest
Root `apm.yml` declares all 16 packages with `subdir` and `version` constraints. The `version` field here is the **marketplace version** — independent from per-package versions.

## Verification commands
Run before committing package changes:
```bash
apm marketplace check        # validates all refs resolve
apm pack --dry-run           # validates marketplace.json generation
```

## Generated file
`.claude-plugin/marketplace.json` is **auto-generated** by `apm pack`. After changing any package or the root `apm.yml`, regenerate it:
```bash
apm pack
```
Then **commit the updated marketplace.json**. The CI workflow will reject PRs with mismatched versions.

## Release process
1. Update root `apm.yml` `version` field
2. Bump per-package versions in `packages/*/apm.yml` as needed
3. Run `apm pack` to regenerate `.claude-plugin/marketplace.json`
4. Commit, tag as `v{version}` (e.g. `v0.2.0`)

## Ignored directories
`apm_modules/`, `.apm_cache/`, `build/`, `*.tar.gz`, `handoffs/`, `to-import/` — all gitignored. Never commit content from these.
