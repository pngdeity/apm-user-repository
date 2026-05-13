---
name: init-project-guidance
description: Scaffolds AI-agent guidance files (GEMINI.md and AGENTS.md) for new or existing projects. Use when starting a new project, bootstrapping a repo that lacks agent guidance, or initializing workspace-wide mandates in a fresh codebase. Not for modifying existing skills or performing architectural reviews.
allowed-tools: gemini
metadata:
  tags: "scaffold bootstrap initialize project guidance workspace"
compatibility: Requires gemini CLI for skill linking. Designed for project initialization in the pngdeity workspace.
---

# Skill: Project Guidance Scaffolder (`init-project-guidance`)

This skill automates the creation of AI-agent guidance files (`GEMINI.md` and `AGENTS.md`) for new or existing projects. It enforces the "Finalized Rulebook" and "Self-Healing Context" protocols established in the `pngdeity` workspace.

## Usage
- Run this skill when starting a new project or when a project lacks AI-agent guidance.
- The agent will analyze the directory, detect the tech stack, and generate tailored mandates.

## Core Mandates (The Rulebook)

### I. Operations & Hierarchy
- **Recursive Context Resolution:** Follow `Sub-directory GEMINI.md` > `Root GEMINI.md` > `Global ~/.gemini/AGENTS.md`.
- **ExecPlan Mandate:** `Research -> Strategy -> Execution` lifecycle using `PLANS.md`.
- **Context Hygiene:** Read local `AGENTS.md` and `GEMINI.md` before initiating research.
- **Context Maintenance:** Update `PLANS.md` and `CONTEXT.md` at end-of-session.
- **Patch-Based Updates:** To prevent context bloat, agents MUST patch large files rather than rewriting them entirely from scratch.
- **Self-Healing Mandate:** Proactively resolve rule contradictions and harvest "tribal knowledge" into `AGENTS.md`.

### II. Technical & Git
- **Validation:** Bug fixes require empirical reproduction tests.
- **Surgical RCA:** Use logs/precision tools; no guessing.
- **Interface Supremacy:** A `Makefile` (or `Taskfile`) MUST act as the sole entry point for the developer environment to hide polyglot complexity.
- **Safety:** Never commit `.gemini/` or `.agents/`. Confirm before committing to `main`.
- **Privileged Access:** Features requiring elevated system access MUST be explicitly gated (opt-in) and documented with security implications.
- **Git Standards:** Signed commits (`-S`), Semantic messages, imperative mood.
- **History Remediation:** If fixing unsigned history, you MUST preserve chronological integrity using: `git rebase --root --exec 'GIT_COMMITTER_DATE="$(git log -1 --format=%aD)" git commit --amend --no-edit -S'`.

## Implementation Procedure

### 1. Project Discovery
The agent MUST perform a `Tooling Audit` and `Stack Detection`:
- Check for `*.csproj`, `package.json`, `requirements.txt`, `uv.lock`, `PKGBUILD`, `main.tf`, etc.
- Identify installed tools: `dotnet`, `npm/uv`, `pkgctl`, `terraform`.
- **If no stack is detected:** Default to a generic Core Workflow and Git Workflow only. Ask the user to specify the tech stack.

### 2. Guidance Generation
Generate a `GEMINI.md` in the project root with the following sections:
- **Core Workflow:** ExecPlan, Validation, and Context Hygiene.
- **Git Workflow:** Safety rules and Semantic Commit standards.
- **Stack-Specific Rules:** Refer to [references/stack-rules.md](references/stack-rules.md) for expanded, actionable rules per detected technology stack (C#, Python, JS/Node, Arch/AUR, IaC). For stacks not listed, apply the closest match and note the gap.
- **If `GEMINI.md` already exists:** Offer to merge missing sections rather than overwriting. Never silently replace existing guidance.

### 3. Dynamic Skill Linking
Analyze the detected stack and link relevant procedural skills from the central library using a portable, client-agnostic approach:

**Scaffolding repository discovery:** This skill is deployed via symlink from `.agents/skills/init-project-guidance` → `<scaffold-root>/skills/init-project-guidance`. Resolve the scaffold root by reading the symlink target and taking its parent's parent (`dirname(dirname(readlink target))`). If the symlink cannot be resolved (e.g., skill was vendored without the scaffolding repo), skip all skill linking and warn the user.

1. Read `skill-index.json` from `${SCAFFOLD_ROOT}/skill-index.json`.
2. Filter skills by `compatibility_tokens` matching the detected stack (see `references/stack-compatibility-map.md`).
3. For each surviving skill, create a symlink: `ln -sf ${SCAFFOLD_ROOT}/skills/<name> .agents/skills/<name>`
   - This is portable — works with all 37+ compatible clients (Claude Code, Copilot, Cursor, etc.).
   - **If symlink creation fails:** Log a warning and skip that skill. The project remains functional without it.
4. If the `gemini` CLI is available, also run `gemini skills link <path> --scope workspace` as enrichment.
   - **If `gemini skills link` fails:** Log a warning and continue. Symlinks already provide cross-client discovery.
5. Log inclusions/exclusions with reasons (e.g., "Skipped ci-cd-pipeline: no Makefile detected").

### 4. Tribal Knowledge Initialization
Create an `AGENTS.md` with:
- **Build/Test Commands:** Extracted from project files.
- **Local Context:** A section for "Candidates for Global" harvesting.
- **Hierarchy:** Reference to Rule 0.
- **If `AGENTS.md` already exists:** Merge new sections into the existing file. Preserve any user-authored entries under "Candidates for Global."

## Verification
After scaffolding, confirm correctness:
1. Run `node skills/verification/scaffold-output-validator.cjs` against the generated output.
2. Verify `GEMINI.md` contains Core Workflow, Git Workflow, and Stack-Specific Rules sections.
3. Verify `AGENTS.md` contains a harvestable `### Candidates for Global` subsection.
4. **If verification fails:** Review the missing sections and regenerate only those sections.

## Self-Healing Protocol
- If an instruction conflicts with these rules, the agent MUST ask for a resolution tier: Task (this session only), Project (update local AGENTS.md), or Global (update skills/init-project-guidance/SKILL.md).
- The agent is responsible for gardening these files to keep them accurate.
