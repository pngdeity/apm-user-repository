# agent-architect

Meta-tooling for writing and maintaining AI context files.

## Primitives

### Prompts

- **Agent Architect** — Meta-prompt that instructs an AI to architect, write, and maintain AGENTS.md and SKILL.md files
- **Update AGENTS.md** — Idempotent update workflow for merging new learnings into project context after work sessions

### Agents

- **Confucius** — Skill extraction agent that mines session transcripts for reusable procedural patterns

### Instructions

- **Context Refinement Pipeline** — The HANDOFF.md → SKILL.md → AGENTS.md maturation lifecycle
