---
description: Future architecture reference for migrating from imperative CI/CD to declarative Cloud Native Buildpacks
applyTo: "**"
---

# Future Architecture: CNCF Cloud Native Buildpacks

While the Inversion of Control architecture successfully decouples build logic from the CI platform via the Bootstrapper pattern, the ultimate evolution is to adopt Cloud Native Buildpacks (CNB). Maintained by the CNCF, Buildpacks eliminate the need for custom build scripts entirely by automatically detecting application frameworks and compiling them into secure, optimized, OCI-compliant container images.

## The Paradigm Shift

Moving to CNB changes the CI/CD contract from *imperative* to *declarative*.

- **Current state:** The repository tells the framework exactly *how* to build the app step-by-step using shell scripts.
- **Future state:** The repository hands the source code to a standardized CNB Builder (e.g., Paketo Buildpacks). The Builder detects language and framework files, automatically provisions the correct toolchains, and outputs a deployable artifact.

## Key Considerations for Migration

### Deployment Target Mismatch (Critical Risk)

Cloud Native Buildpacks natively output *container images* (Docker images), not static files. If the deployment target only accepts static assets, either:
1. **Migrate hosting** to a containerized runtime (Google Cloud Run, AWS App Runner, Azure Container Apps, or Kubernetes).
2. **Implement an artifact extraction step:** Use CNB to build the container, run it temporarily, extract compiled static files from the container's `/workspace` directory, and push those files to the static host.

### Pipeline Simplification

Once the deployment target is resolved, the CI orchestrator reduces to a single command:

```bash
pack build my-website --builder paketobuildpacks/builder-jammy-base
```

### Deprecation of Custom Scripts

The `ci/setup-env.sh` and `ci/execute-build.sh` files are deleted. The CNB lifecycle handles all dependency fetching, caching, and compilation automatically.

### Local Developer Experience

Testing the CNB architecture locally requires the `pack` CLI and Docker. The build process replicates exactly between local and cloud environments without relying on CI runners for validation. On Linux, installing and executing the `pack` CLI is a frictionless, native experience.

### Security & Patching

CNBs introduce "rebasing," which allows the underlying OS layers of the container to be patched for security vulnerabilities without requiring a full rebuild of the application code. This significantly reduces long-term maintenance overhead.
