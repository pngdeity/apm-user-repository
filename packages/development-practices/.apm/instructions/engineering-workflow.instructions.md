---
description: Engineering workflow mandates — ExecPlan lifecycle, context hygiene, and empirical validation
applyTo: "**"
---

## Engineering Workflow
- **ExecPlan Mandate:** For all complex features or significant refactors, follow the Research -> Strategy -> Execution cycle. Use a dedicated `PLANS.md` (following the schema in `~/.gemini/PLANS.md`) to document and track progress.
- **Context Hygiene:** Always analyze the foundational context files before initiating research to ensure architectural alignment and tool awareness. This includes Gemini-specific files (`AGENTS.md`, `CONTEXT.md`, `GEMINI.md`), Claude Code conventions (`CLAUDE.md`), and Codex directives (`.codex` or `CODEX.md`).
- **Validation:** Never consider a bug fix complete without an empirical reproduction test.
