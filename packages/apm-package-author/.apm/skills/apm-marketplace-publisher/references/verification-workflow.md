# Verification Workflow Reference

## Command Reference

| Command                 | What it does                                                                                   | Exit code                        |
| ----------------------- | ---------------------------------------------------------------------------------------------- | -------------------------------- |
| `just fmt`              | Formats YAML files (`apm.yml` + `packages/*/apm.yml`) with `yamlfmt`                           | 0 on success                     |
| `just sync-check`       | Verifies package versions satisfy root constraints (read-only)                                 | 1 on mismatch                    |
| `just sync-fix`         | Auto-bumps mismatched package versions to satisfy constraints                                  | 0 on success, 1 if unrecoverable |
| `just yaml-validate`    | Validates all apm.yml files against `schemas/apm.json` with `check-jsonschema`                 | 0 if all valid                   |
| `apm marketplace check` | Validates all marketplace refs resolve correctly                                               | 0 on success                     |
| `apm pack`              | Generates `.claude-plugin/marketplace.json`                                                    | 0 on success                     |
| `just pre-commit-check` | Full pipeline: fmt → sync-fix → yaml-validate → marketplace check → pack → stale check         | 1 on any failure                 |
| `just ci`               | CI equivalent: fmt-check → sync-check → yaml-validate → marketplace check → pack → stale check | 1 on any failure                 |

## Pre-Commit Pipeline (`just pre-commit-check`)

```
just fmt                 # Format YAML
just sync-fix            # Align package versions
just yaml-validate       # Schema validation
apm marketplace check    # Resolve refs
apm pack                 # Regenerate marketplace.json
# stale check            # Verify marketplace.json is committed
```

The stale check at the end ensures `.claude-plugin/marketplace.json` matches the
committed state. If `apm pack` changed it, the check fails — you must commit the
regenerated file.

## CI Pipeline (`just ci`)

```
just fmt-check           # Lint YAML (no formatting)
just sync-check          # Check versions (no modification)
just yaml-validate       # Schema validation
apm marketplace check    # Resolve refs
apm pack --dry-run       # Validate packability
# stale check            # Verify committed marketplace.json
```

CI does not modify files — it uses `sync-check` (not `sync-fix`) and `fmt-check`
(not `fmt`). Any failure blocks the PR.

## Common Failures and Fixes

| Error                                | Cause                                             | Fix                                                                  |
| ------------------------------------ | ------------------------------------------------- | -------------------------------------------------------------------- |
| `[mismatch] package: X.X.X -> 0.4.2` | Package version doesn't satisfy root constraint   | `just sync-fix` or manually update `packages/<name>/apm.yml version` |
| `yamlfmt: not found`                 | `yamlfmt` not installed                           | `go install github.com/google/yamlfmt/cmd/yamlfmt@latest`            |
| `check-jsonschema: not found`        | `check-jsonschema` not installed                  | `pip install check-jsonschema`                                       |
| `No such file: schemas/apm.json`     | Schema file missing                               | Ensure `schemas/apm.json` exists in repo root                        |
| `marketplace.json is stale`          | `apm pack` generated changes not committed        | `git add .claude-plugin/marketplace.json && git commit`              |
| `entry not found / ref not OK`       | Package path doesn't exist or ref doesn't resolve | Check `subdir` path and ensure files are committed                   |
| `Found N version mismatch(es)`       | Multiple packages out of sync                     | `just sync-fix` then re-run                                          |

## YAML Formatting

Config: `.yamlfmt` at repo root. Rules:

- 2-space indent
- LF line endings
- Trailing newline
- Trimmed trailing whitespace
- No max line length
- Matches: `apm.yml`, `packages/*/apm.yml`

## Schema Validation

Schema: `schemas/apm.json` (JSON Schema Draft-07). Permissive — accepts unknown
fields, validates structure of known fields. Required fields: `name`, `version`.
Covers both root and package manifests.

Schema ID:
`https://raw.githubusercontent.com/pngdeity/apm-user-repository/main/schemas/apm.json`

## CI Workflow File

`.github/workflows/validate.yml` triggers on push/PR to `main` when
`packages/**`, `apm.yml`, `justfile`, or `schemas/**` change. It:

1. Checks out the repo
2. Installs Go 1.26 + `just`
3. Installs `yamlfmt` (via `go install`) + `check-jsonschema` (via
   `pip install`)
4. Runs `just ci`
