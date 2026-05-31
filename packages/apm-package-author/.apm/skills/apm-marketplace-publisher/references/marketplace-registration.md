# Marketplace Registration Reference

## Root `apm.yml` Marketplace Block

The marketplace block is a top-level section of the root `apm.yml` that declares
available packages:

```yaml
marketplace:
  owner:
    name: pngdeity
    url: https://github.com/pngdeity
  build:
    tagPattern: "v{version}"
  packages:
    - name: package-name
      description: Short description for consumers
      source: pngdeity/apm-user-repository
      subdir: packages/package-name
      version: ">=0.0.1"
```

### Package Entry Fields

| Field         | Required | Notes                                                           |
| ------------- | -------- | --------------------------------------------------------------- |
| `name`        | Yes      | Package identifier, must match `packages/<name>/apm.yml` `name` |
| `description` | Yes      | Consumer-facing description in marketplace listings             |
| `source`      | Yes      | For this repo: `pngdeity/apm-user-repository`                   |
| `subdir`      | Yes      | Path relative to repo root: `packages/<name>`                   |
| `version`     | Yes      | Semver constraint for consumers — e.g., `">=0.0.1"`, `"^1.0.0"` |

### Alphabetical Order

Package entries are ordered alphabetically by `name`. Insert new entries at the
correct position to avoid CI format failures.

## Version Constraints

### Semver Constraint Types

| Constraint          | Meaning                            | `sync-check`       | `sync-fix`                                                     |
| ------------------- | ---------------------------------- | ------------------ | -------------------------------------------------------------- |
| `">=0.0.1"`         | Any version >= 0.0.1               | Passes             | Bumps to 0.0.1 if below                                        |
| `"^0.4.2"`          | >=0.4.2, <0.5.0                    | Passes             | Bumps to 0.4.2 if below                                        |
| `"~0.4.2"`          | >=0.4.2, <0.5.0 (patch-compatible) | Passes             | Bumps to 0.4.2 if below                                        |
| `"0.4.2"`           | Exact match                        | Passes if == 0.4.2 | Bumps to 0.4.2                                                 |
| `">=0.0.1 <=1.0.0"` | Range                              | Passes             | **Fails** — `extractBaseVersion` can't parse multi-part ranges |

### Current Convention

This repo uses `">=0.0.1"` (lower-bound only). This allows individual packages
to be independently versioned without an upper cap. The `sync-versions` tool's
`extractBaseVersion()` function strips the `>=` prefix to derive the minimum fix
target (`0.0.1`).

## `sync-versions` Tool

Located at `cmd/sync-versions/main.go`. It:

1. Reads each package's constraint from root `apm.yml marketplace.packages`
2. Reads each package's declared version from `packages/<name>/apm.yml`
3. Validates: `semver.NewConstraint(constraint).Check(packageVersion)`
4. In `--check` mode: reports mismatches, exits non-zero
5. In fix mode: writes the minimum valid version into the package's `apm.yml`

### Constraint Enforcement

The constraint enforces two properties:

- **Lower bound**: Package version can't be older than the marketplace promises
- **Upper bound** (if present): Package version can't drift beyond the
  marketplace minor/major line

With `">=0.0.1"` (lower only), packages are free to advance independently. The
marketplace root version (`0.4.2`) still governs the git tag (`v0.4.2`) that
`apm pack` uses for the `ref` in `marketplace.json`.

## Release Process

1. Run `just pre-commit-check` (full pipeline)
2. Commit with conventional format: `feat: add <name> to marketplace`
3. Push to `main`
4. Consumers run `apm marketplace update` to refresh their local cache
5. Consumers install with `apm install <name>@pngdeity-agent-context`

No separate git tag is needed per package — the `tagPattern: "v{version}"` uses
the marketplace root version, and all packages share the same tag.
