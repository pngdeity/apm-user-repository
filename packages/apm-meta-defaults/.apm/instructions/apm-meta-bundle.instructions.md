---
description: APM meta-authoring defaults bundle — aggregates packages for APM package creation, skill authoring, migration, evaluation, quality gates, agent orchestration, and project scaffolding. Installed as transitive dependencies via apm-meta-defaults.
applyTo: "**"
---

# APM Meta-Authoring Defaults

This bundle installs the core APM meta-authoring packages:

| Package                     | Purpose                                                       |
| --------------------------- | ------------------------------------------------------------- |
| `apm-package-author`        | Create APM packages (anatomy, primitive selection, authoring) |
| `apm-marketplace-publisher` | Included transitively via `apm-package-author`                |
| `apm-migration`             | Migrate consumer projects to APM                              |
| `custom-skill-creator`      | Create and package AI agent skills                            |
| `agent-architect`           | Meta-tooling for AI context files                             |
| `skill-eval-pipeline`       | Evaluate and optimize SKILL.md files                          |
| `skill-eval-agents`         | Evaluation agents and Go CLI tools                            |
| `context-quality-gate`      | CI/CD quality gate for context files                          |
| `ai-agent-orchestration`    | Agent orchestration design patterns                           |
| `init-project-guidance`     | Scaffold AGENTS.md for new projects                           |

When you install `apm-meta-defaults`, all of these packages become available in
your agent context.
