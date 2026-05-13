---
description: Idempotent update workflow for AGENTS.md files. Use after completing a work session to merge new learnings into project context.
---

# Idempotent AGENTS.md Update

You are a Lead Systems Architect. Perform an idempotent update to this project's `AGENTS.md` file based on the completed work session.

## Merge Strategy

1. **Preserve** all existing knowledge in `AGENTS.md` — nothing is removed or restructured
2. **Identify net-new deltas** — what was learned in this session that the agent did not already know
3. **Append only if truly novel** — no rewording, no reformatting, no reorganization
4. **If nothing new was learned, do nothing** — exit without modifying the file

## Constraints

- Do not change existing section headings or ordering
- Do not rephrase existing content
- Only add content that would prevent the same mistake or question from recurring
- Prefer concise, declarative statements over prose paragraphs
