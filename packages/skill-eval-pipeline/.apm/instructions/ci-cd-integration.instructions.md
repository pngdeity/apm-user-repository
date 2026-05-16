---
description: How to wire the skill eval pipeline into GitHub Actions CI/CD — CLI setup, authentication, Go compilation, and quality gate enforcement
applyTo: "**/.github/workflows/*.yml"
---

# CI/CD Integration for Skill Eval Pipeline

## Prerequisites

- **opencode CLI**: Installed via `curl -fsSL https://opencode.ai/install.sh | sh`
- **Gemini CLI**: Installed via `npm install -g @google/gemini-cli`
- **Go 1.24+**: Required to compile eval tool binaries
- **just**: Command runner installed via `https://just.systems/install.sh`
- **skills-ref**: Installed via `npm install -g skills-ref`

## Authentication

Two secrets must be stored in the GitHub repository:

| Secret | Purpose |
|--------|---------|
| `OPENCODE_API_KEY` | Authenticates opencode CLI for agent invocation |
| `GEMINI_API_KEY` | Authenticates Gemini CLI for eval prompt execution |

Set these via `gh secret set` or the repository Settings > Secrets > Actions page.

## Go Compilation

The CI/CD workflow auto-compiles Go binaries when `.go` files change under `packages/skill-eval-agents/.apm/scripts/`. The following binaries are built:

- `invoke-cli` — Invokes agents with context files
- `parse-session` — Parses agent session output
- `compute-benchmark` — Computes evaluation benchmark metrics
- `select-best` — Selects the best candidate from results

Output binaries are placed in `bin/` at the repo root.

## Path Filtering

The eval pipeline runs on pull request events when any files under `packages/**` change. For manual dispatches, specify comma-separated skill directory paths via the `skill_paths` input.

## Artifacts

Two artifacts are uploaded per eval run:

- `skill-eval-report` — `benchmark.json` with raw evaluation metrics
- `skill-eval-selected` — `selected.json` with the best candidate and composite scores

## Quality Gate

The quality gate fails the PR if:

1. Structural validation fails (`skills-ref validate` returns errors on any changed SKILL.md)
2. Composite score falls below the minimum threshold

## IoC Pattern

GitHub Actions workflows are thin wrappers around `just` recipes. All eval logic lives in the Justfile — not inline in the workflow. This follows the existing ci-cd-standards package conventions.

## Gotchas

- **API keys**: Both `OPENCODE_API_KEY` and `GEMINI_API_KEY` must be valid and present in repo secrets. Missing keys cause silent failures in CLI setup.
- **Workspace cleanup**: The eval workspace must be cleaned between runs to avoid cross-contamination of results. The `just eval-skill` recipe handles this.
- **No-op guard**: If no SKILL.md files changed, skip the eval pipeline entirely to avoid unnecessary cost.
