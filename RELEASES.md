# RELEASES

## v0.2.0 (2026-05-13)

### Package enhancements

- **agent-architect** — Confucius agent expanded from 36 to 162 lines of source;
  meta-prompt converted to skill with 2 references (knowledge-sources.md,
  file-specifications.md)
- **ci-cd-standards** — GHA best practices merged into gha-standards.instructions.md;
  CNCF Buildpacks addendum merged into cncf-buildpacks.instructions.md;
  vendor-decoupled ci-cd-pipeline skill added
- **All 16 packages** — YAML frontmatter added to all primitives lacking it;
  naming audit passed (16/16 stems validated against first-principles naming)

### Marketplace infrastructure

- CI workflow validates with canonical two-step (`apm marketplace check` +
  `apm pack --dry-run`)
- `microsoft/apm-action@v1` with `setup-only: 'true'` as reference implementation

### Consumer migration

11 consumer projects migrated to APM with `apm.yml`, `apm.lock.yaml`,
`.apm/instructions/`, and `.gitignore` entries. 6 projects deferred with
handoff documents. 15 projects assessed and skipped.

### Consumer upgrade instructions

If your `apm.yml` pins `^0.1.0`, update to `^0.2.0`:

```bash
apm install --update
```

Or manually edit `apm.yml` to change `version: "^0.1.0"` to `version: "^0.2.0"`
for each installed package, then run `apm install`.

The content changes in v0.2.0 are additive — no breaking changes, no removed
primitives, no renamed packages.
