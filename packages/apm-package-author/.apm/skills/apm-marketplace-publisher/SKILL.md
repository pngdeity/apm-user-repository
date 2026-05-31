---
name: apm-marketplace-publisher
description: Skill for adding an APM package to a marketplace and publishing. Use when asked to register a package in apm.yml, add a package to the marketplace, publish an APM package, or integrate a package into the pngdeity-agent-context marketplace. Assumes the package already exists in packages/. For standalone package creation, use the apm-package-author skill.
allowed-tools: apm, git, just
metadata:
  tags: "meta apm marketplace registration publishing versioning"
compatibility: Requires apm CLI and just. For marketplace repos with sync-versions tooling.
---

# Skill: APM Marketplace Publisher

## When to Invoke

Invoke this skill when the user asks to:

- "add a package to the marketplace"
- "register a package in apm.yml"
- "publish an APM package"
- "integrate [package] into the marketplace"
- "add [package] to pngdeity-agent-context"

If the package does not yet exist, use the `apm-package-author` skill first.

## Phase 5: Register in the Marketplace

Add a new entry to the root `apm.yml` under `marketplace.packages`, maintaining
alphabetical order by `name`:

```yaml
marketplace:
  packages:
    - name: my-package
      description: Short description for marketplace consumers
      source: pngdeity/apm-user-repository
      subdir: packages/my-package
      version: ">=0.0.1"
```

### Version constraint

The `version` field in the marketplace entry is a **semver constraint for
consumers** — not the package's own version. The package's own version lives in
`packages/<name>/apm.yml`.

- `">=0.0.1"` — lower bound only, allows any version >= 0.0.1 (current
  convention in this repo)
- `"^0.4.2"` — caret constraint: `>=0.4.2, <0.5.0`
- `"0.4.2"` — exact match

The `sync-versions` tool verifies that the package's declared version satisfies
this constraint. Lower-bound-only (`>=`) works with both `sync-check` and
`sync-fix`.

See `references/marketplace-registration.md` for the full marketplace block
schema and field descriptions.

## Phase 6: Align Versions

Run the version sync tool to verify the package's version satisfies the root
constraint:

```bash
just sync-check        # read-only check — exits non-zero on mismatch
just sync-fix          # auto-bumps mismatched package versions
```

If `sync-check` passes, proceed. If it fails, either:

- `just sync-fix` to auto-bump, or
- Manually update the package's `apm.yml` version, then re-run `sync-check`

## Phase 7: Publish

Run the full verification pipeline:

```bash
just pre-commit-check   # fmt → sync-fix → yaml-validate → marketplace check → pack → stale check
```

This regenerates `.claude-plugin/marketplace.json`. Commit everything:

```bash
git add apm.yml packages/<name>/ .claude-plugin/marketplace.json
git commit -S -m "feat: add <package-name> to marketplace"
git push
```

The package is now available to consumers.

## Consumer Installation

After publishing, consumers install the package with:

```bash
apm marketplace add pngdeity/apm-user-repository
apm marketplace update
apm install <package-name>@pngdeity-agent-context
apm compile
```

See `references/verification-workflow.md` for the full CI pipeline and
troubleshooting.
