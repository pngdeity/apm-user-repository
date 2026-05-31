---
description: APM skill evaluation bundle — aggregates packages for evaluating SKILL.md quality, running multi-agent eval pipelines, and enforcing CI quality gates. Installed as transitive dependencies via apm-eval-defaults.
applyTo: "**"
---

# APM Evaluation Defaults

This bundle installs the skill evaluation toolkit:

| Package                | Purpose                                                               |
| ---------------------- | --------------------------------------------------------------------- |
| `skill-eval-pipeline`  | Multi-agent orchestration for evaluating SKILL.md files               |
| `skill-eval-agents`    | 7 sub-agent skills + Go CLI tools for evaluation                      |
| `context-quality-gate` | CI/CD quality gate enforcing structural compliance + trigger accuracy |

When you install `apm-eval-defaults`, your agent can evaluate skill quality, run
multi-agent evaluation pipelines, and enforce quality gates in CI.
