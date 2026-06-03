# TODO

## High Priority

- [ ] Declare `marketplace.versioning.strategy: per_package` in `apm.yml` —
      current default `lockstep` requires all packages to share the root version.
      Mitigated: all package constraints use `>=0.0.1` so root bumps do not force
      reversioning. True per-package independent versioning still desirable.
- [x] Add `apm pack --check-versions --check-clean` to `just ci` release
      pipeline
- [ ] Use `apm marketplace package add` instead of hand-editing `apm.yml` for
      new packages
- [x] Add SARIF output (`--format sarif`) and `codeql-action/upload-sarif` to
      the existing `apm audit --ci --no-drift` step in CI for GitHub Code
      Scanning integration

## Medium Priority

- [x] Add `apm compile --validate` to CI to catch frontmatter/structure errors
      in primitives
- [ ] Create `apm-policy.yml` at org level (`warn` mode first, then `block`
      after triage)
- [ ] Set up `apm marketplace publish` with `consumer-targets.yml` for automated
      fan-out PRs
- [ ] Add SHA256 checksum generation to release pipeline (canonical release
      sequence from APM docs)

## Low Priority

- [ ] Add `apm audit --strip` to check for hidden Unicode in deployed files
- [ ] Run `apm marketplace outdated` periodically to check for stale upstream
      refs
- [ ] Add `apm marketplace doctor` to `pre-commit-check`
- [ ] Add multi-profile marketplace outputs (Codex
      `.agents/plugins/marketplace.json`)
- [ ] Use `apm plugin init` instead of manual `mkdir` scaffolding for new
      packages
- [ ] Resolve dogfood local deps blocking `apm pack` in CI (`devDependencies` or
      remote refs)
- [ ] Add `apm list`/`apm deps` inspection commands to troubleshooting workflow
- [ ] Explore hooks (`.apm/hooks/`) and commands (`.apm/commands/`) primitives
      for CI integration
- [x] Add tamper detection CI pattern (`apm audit --ci --no-drift` with
      `setup-only: true`)
- [x] Add `apm prune` to cleanup workflow for orphaned dependencies
