---
name: ci-cd-pipeline
description: Workflow for designing and maintaining high-performance, vendor-decoupled CI/CD pipelines. Use when setting up new build systems or optimizing existing automation. Not for configuring specific CI providers or writing pipeline YAML directly — use Makefile targets instead.
allowed-tools: make
metadata:
  tags: "ci cd pipeline build automation make devops"
compatibility: Requires make. Designed for projects with a Makefile as the developer interface entry point.
---

# CI/CD Pipeline Skill

This skill governs the design of resilient and efficient automation pipelines. It prioritizes local reproducibility and fast feedback loops.

## Workflow: Pipeline Construction

1. **Interface with Makefile:** Ensure every CI step corresponds to a `Makefile` target.
2. **Implement Fail-Fast:** Order the jobs so that the fastest, most critical checks (linting) run first.
3. **Configure Matrix:** Use matrix strategies for High-Performance Architecture (HPA) builds.
4. **Reference Best Practices:** Adhere to the standards in [references/best-practices.md](references/best-practices.md).

## Workflow: The "Verify-then-Publish" Gate

- **Mandate:** Before any remote repository push, the agent MUST run the local `make test` or `make ci` suite.
- **Reporting:** If a local check fails, stop and fix the issue before attempting to push.

## Workflow: Infrastructure Simulation

- For complex deployments, use a two-tier simulation (e.g., Local Container -> Staging -> Prod) to verify environment-specific variables before final delivery.

## Workflow: Error Recovery

- **If `make test` fails:** Read the failure output, fix the issue, and re-run. Never skip or silence a failing test.
- **If `make ci` is not defined:** Create it as a target that runs lint, test, and build in sequence.
- **If matrix build fails for a single platform:** Investigate platform-specific issues. Do not disable the failing platform without documenting the rationale.

## Verification
Verify this skill produces correct pipeline configuration:
1. Verify every CI step has a corresponding `make` target.
2. Confirm fail-fast ordering: linting runs before tests, tests before builds.
3. Confirm the pipeline includes a local `make test` gate before any remote push step.
4. For matrix builds, confirm at least two platform variants are configured.
