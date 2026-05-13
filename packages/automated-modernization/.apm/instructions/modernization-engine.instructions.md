---
description: Automated Repository Modernization Engine (ARME) blueprint — Phase A discovery, Phase B staging strategy, Phase C+ implementation framework for modernizing arbitrary software repositories
applyTo: "**"
---
# Blueprint: Automated Repository Modernization Engine (ARME)

## Objective
To provide a standardized, machine-executable framework for identifying, planning, and implementing modernization efforts across any arbitrary software repository.

---

## Phase A: Intelligent Discovery (The "Architect")
**Goal**: Generate a comprehensive map of technical debt and linguistic deficiencies for the repository.

### Audit Vectors
- **Linguistic/Hygiene**: Technical American English compliance, tone audit, profanity check, and identification of defunct/obsolete distribution or system references.
- **Standards/Compliance**: Verification of linting (`shellcheck`, `eslint`), formatting (`shfmt`, `prettier`), and commit standards (Conventional Commits).
- **Architecture**: Analysis of modularity (monolithic vs. library-based), configuration handling (hardcoded vs. environment-aware), and build system flexibility.
- **Functional**: Identification of missing platform-native features (e.g., systemd user units, containerization, unprivileged execution modes).

**Output**: `MODERNIZATION_MANIFEST.json` or `DEBT_MAP.md` categorizing issues by Severity (Critical to Low) and Complexity (Simple to Architectural).

---

## Phase B: Staging Strategy (The "Planner")
**Goal**: Partition the manifest into a sequence of atomic, verifiable implementation plans.

### Partitioning Logic
- **Stage 1 (Professionalism)**: Non-breaking changes focusing on linguistics, documentation formatting, and hygiene.
- **Stage 2 (Quality Gates)**: Implementation of automated linting, formatting, and CI/CD pipeline standardization.
- **Stage 3 (Refactor)**: Internal logic modernization, including data structure upgrades and architectural modularization.
- **Stage 4 (Features)**: Major functional expansions and modern platform integrations.

### Staging Document Standard
Each `STAGEX_PLAN.md` must include:
- **Task List**: Atomic, checkbox-style actions for tracking.
- **Context Blocks**: Specific file paths and technical constraints.
- **Verification Logic**: Exact shell commands required to validate the stage (e.g., `make check`).
- **Safety Mandates**: Explicit "Do Not Touch" zones to prevent regression.

---

## Phase C: Execution Orchestration (The "Orchestrator")
**Goal**: Manage the sequential execution loop of implementation agents.

### The Sequential Loop
1. **Bootstrap**: Create a dedicated `fix/modernization-base` branch.
2. **Dispatch**: Spawn an implementation agent (e.g., GitHub Copilot or a sub-agent) with the current `STAGEX_PLAN.md`.
3. **Targeting**: All Pull Requests must target the `modernization-base` branch.
4. **Validation**: The orchestrator verifies that the agent has executed the "Verification Logic" and that CI passes.
5. **Finalization**: Merge the PR, delete the temporary implementation branch, and trigger the next stage.

---

## Automation Constraints & Safety
- **Idempotency**: Agents must ensure that rerunning a stage results in a "no-op" if the work is already complete.
- **Recursive Testing**: Stage $N$ must always pass the verification logic established in Stage $N-1$.
- **Source-Awareness**: Agents must prioritize instructions found in `CONTRIBUTING.md` or `README.md` to ensure modernization respects original project intent.
