# CI/CD Pipeline Best Practices

## 1. Vendor Decoupling
- **Principle:** CI/CD YAML files (GitHub Actions, GitLab CI, etc.) should be "dumb" wrappers.
- **Implementation:** The YAML should only call a `Justfile` recipe (e.g., `just ci`) or `Makefile` target (e.g., `make ci`). This ensures the pipeline can be run locally and migrated between vendors with minimal friction.

## 2. Fail-Fast Cascades
- **Principle:** Halt execution at the earliest possible failure point to save compute and time.
- **Ordering:** Linting > Type-Checking > Unit Tests > Integration Tests > Build > Security Scan.

## 3. Verify-then-Publish
- **Principle:** Never trigger a remote build/push without a local verification pass.
- **Implementation:** Use pre-commit or pre-push hooks to run the `just ci` or `make ci` suite.

## 4. HPA Matrix Builds
- **Principle:** Build artifacts for multiple target architectures simultaneously.
- **Implementation:** Use matrix builds to compile generic binaries alongside optimized variants (e.g., `x86-64-v3`, `aarch64`).
