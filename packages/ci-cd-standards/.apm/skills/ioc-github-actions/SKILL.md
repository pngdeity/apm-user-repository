---
name: ioc-github-actions
description: Inversion of Control CI/CD pipeline architecture. Use when designing or refactoring CI/CD pipelines to extract domain logic into testable shell scripts.
---

# Inversion of Control CI/CD Pipeline

An architectural pattern that treats the CI/CD platform as a dumb orchestrator. All domain logic, dependency management, and compilation instructions are extracted into standard shell scripts housed within the repository.

## The Contract

The CI/CD platform runs a single entrypoint script. The framework file (e.g., `.github/workflows/build-deploy.yaml`) defines only triggers and artifact routing. Build logic lives entirely in version-controlled scripts.

## Key Benefits

- **Local Parity:** Because execution logic resides in standard shell scripts, the entire pipeline can be run, tested, and debugged natively on a local machine prior to pushing to the CI runner.
- **Vendor Independence:** Migrating CI providers requires updating only the thin orchestrator file — zero changes to build logic.
- **Testability:** Shell scripts are independently testable without invoking the CI platform.

## Implementation Pattern

### Orchestrator File

The CI platform file defines triggers and delegates execution:

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - run: ./ci/build.sh
      - uses: actions/upload-artifact@v7
        with:
          path: ./ci/out/
```

### Bootstrapper Script Pattern

```
ci/
├── build.sh           # Main entrypoint, sequences the pipeline
├── setup-env.sh        # Bootstraps runner environment (downloads tooling into ./ci/bin/)
├── execute-build.sh    # Handles compilation and standardizes output into ./ci/out/
└── out/                # Standardized artifact output directory
```

### Design Rules

1. **Single entrypoint.** The CI file runs exactly one script.
2. **Isolated tooling.** Fetch build dependencies into a local directory (`./ci/bin/`) rather than assuming global installation.
3. **Standardized output.** All build artifacts land in a single directory (`./ci/out/`) for consistent artifact upload.
4. **Execution permissions.** Ensure all shell scripts have execution permissions (`chmod +x ci/*.sh`) before committing.

## Migration Pattern

When migrating from an imperative CI file to IoC:

1. Identify all inline shell commands in the existing CI YAML
2. Extract each logical group into a named script under `ci/`
3. Replace inline commands with script invocations in the CI YAML
4. Verify local parity by running `./ci/build.sh` from the project root
5. Push and monitor the CI console for any missing dependency errors

## Future Containerization

The IoC contract guarantees that migrating the pipeline to a containerized build (Docker, CNCF Buildpacks) requires zero changes to the CI platform YAML file. The entrypoint script is replaced with a container invocation; the orchestrator remains unchanged.
