---
name: context-quality-gate
description: Enforce quality thresholds for AI agent context files in CI/CD. Use when validating SKILL.md files in pull requests, running automated quality checks on agent skills, or gating skill publication behind measurable quality criteria.
allowed-tools: read_file, run_shell_command, glob, write_file
metadata:
  tags: "ci-cd quality gate validation context testing"
compatibility: Requires opencode CLI, gemini CLI, go 1.24+, skills-ref. Designed for GitHub Actions with justfile orchestration.
---

# Context Quality Gate

This skill gates context file publication behind quality thresholds. When loaded in CI/CD:

## Workflow

### Step 1: Structural Gate

Run `skills-ref validate` on each changed SKILL.md. Block the pipeline if any structural errors are found. This ensures YAML frontmatter, required fields, and format compliance before any semantic evaluation occurs.

Verify that `evals/workspace/struct-validation.json` has `"status": "pass"`. If `"status": "fail"`, fail the gate with a summary of structural errors.

### Step 2: Trigger Gate

Read trigger results from `evals/workspace/trigger-results/opencode.json` and `evals/workspace/trigger-results/gemini.json`. Confirm:

- **Threshold:** trigger rate ≥ 0.5 on all evaluated CLI targets.
- If any CLI's activation rate is below 0.5, the gate fails. A low trigger rate indicates the context file is not reliably activating when it should, which undermines its utility.

### Step 3: Quality Gate

Read `evals/workspace/quality-results/iteration-1/benchmark.json` and verify:

- **Threshold:** pass rate delta (with — without skill) ≥ 0.1.
- If delta ≤ 0.1, the gate fails. A delta below 0.1 suggests the context file provides negligible measurable value over the no-skill baseline.

### Step 4: Selection Gate

Read `evals/workspace/selected.json` and verify:

- A best candidate was selected (`selected` is not null).
- **Threshold:** composite score ≥ 0.5.

The composite score combines pass rate (0.5 weight), trigger rate (0.3 weight), and token efficiency delta (0.2 weight). A composite below 0.5 indicates the skill is not production-ready.

### Step 5: Post PR Report

Using `write_file`, post a PR comment with a Markdown table showing results for each evaluated skill. Source data from `evals/workspace/selected.json`, `evals/workspace/quality-results/iteration-1/benchmark.json`, and `evals/workspace/trigger-aggregation.json`.

| Skill | Structural | Trigger Rate | Pass Delta | Composite | Status |
|-------|------------|-------------|------------|-----------|--------|
| ...   | pass/fail  | 0.XX        | 0.XX       | 0.XX      | pass/fail |

A skill passes only if ALL four gates above pass. Any single gate failure marks the skill as `fail`.

## IoC Pattern

Follow the existing IoC pattern from ci-cd-standards: the GitHub Actions workflow calls `just` recipes rather than inline scripts. Example:

```
just eval-gate-check
```

## Verification

Verify this skill produces correct gating decisions:

1. Set up `evals/workspace/` with fixture files: `struct-validation.json` (status=pass), `trigger-results/opencode.json` (rate=0.72), `trigger-results/gemini.json` (rate=0.68), `quality-results/iteration-1/benchmark.json` (delta=0.15), `selected.json` (selected=revision-B, composite=0.72).
2. Run context-quality-gate. Confirm all gates pass and status is `pass`.
3. Change `trigger-results/opencode.json` rate to 0.4. Re-run. Confirm gate fails with trigger rate below threshold.
4. Change `benchmark.json` delta to 0.05. Re-run. Confirm gate fails with quality delta below threshold.
5. Change `selected.json` composite to 0.4. Re-run. Confirm gate fails with composite below threshold.
6. Change `struct-validation.json` status to fail. Re-run. Confirm gate fails immediately on structural validation.

## Gotchas

- **Valid API keys**: opencode and gemini require valid API keys stored in GitHub secrets (`OPENCODE_API_KEY`, `GEMINI_API_KEY`). Missing keys cause silent failures.
- **Workspace hygiene**: The eval workspace must be cleaned between runs. The `just eval-skill` recipe handles this automatically.
- **Early exit**: If no SKILL.md files changed in the PR, skip the pipeline to avoid unnecessary cost and CI time.
