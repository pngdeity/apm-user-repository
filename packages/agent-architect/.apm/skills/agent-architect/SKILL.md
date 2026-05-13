---
name: agent-architect
description: Architect, write, structure, and maintain AI agent context files (AGENTS.md and SKILL.md). Use when creating new agent context configurations, scaffolding AGENTS.md files, designing SKILL.md modules, or determining where content belongs in the agent abstraction stack. Not for writing code or implementing features.
allowed-tools: grep, glob, read_file, write_file, replace, run_shell_command
metadata:
  tags: "agent context architecture scaffolding skill-creation agents-md"
compatibility: Generic — no environment restrictions. Consults agentskills.io, Microsoft Agent Framework, Google ADK, and Anthropic sources for authoritative guidance.
---

# Agent Architect Skill

This skill governs the creation, structuring, and maintenance of AI agent context files. It enforces the "Abstraction Stack" model and provides a decision framework for determining correct file placement.

## Abstraction Stack

Structure agent knowledge across two distinct layers to optimize token context and execution reliability:

1. **Foundation (`AGENTS.md`):** Agent identity, global boundaries, and environment. Always active.
2. **Modules (`SKILL.md`):** Specialized workflow recipes, loaded on demand via progressive disclosure.

## Decision Heuristic

When documenting a process or configuring agent behavior, apply this logic to determine the correct output:

1. **Global constraint, architectural baseline, or repository-wide command?** → Write it in `AGENTS.md`.
2. **Single, atomic function call (e.g., querying a database or an API)?** → This is a *Tool* (a verb), not a context file. Suggest a code implementation. Consult the Google ADK and Microsoft Agent Framework sources (see `references/knowledge-sources.md`) to clarify the boundary between Tools and Skills.
3. **Specialized, multi-step workflow requiring reasoning?** → This is a *Skill* (expertise). Create a new skill directory and write a `SKILL.md` file.
4. **Optimizing repository rules for a specific IDE or LLM?** → Create a model-specific override file (e.g., `CLAUDE.md`, `GEMINI.md`) that inherits the spirit of `AGENTS.md` but translates formatting and tool-calling instructions into that model's optimal syntax.

## Workflow: Creating Context Files

1. **Determine placement:** Apply the Decision Heuristic above.
2. **Consult authoritative sources:** Read `references/file-specifications.md` for required components of each file type (AGENTS.md, SKILL.md, model-specific overrides).
3. **Consult knowledge base:** Read `references/knowledge-sources.md` for the canonical sources to consult based on the task at hand (agentskills.io specification, GitHub Copilot lessons, Google ADK, Microsoft Agent Framework, Anthropic guidance).
4. **Draft:** Write the file with mandatory components and proper formatting.
5. **Validate:** For SKILL.md, verify YAML frontmatter compliance against agentskills.io specification.

## Verification

1. Confirm the file type (AGENTS.md / SKILL.md / model override) matches the Decision Heuristic outcome.
2. For SKILL.md: verify YAML frontmatter has `name` (kebab-case) and `description` (written as an API doc for the agent router).
3. For AGENTS.md: verify it includes Persona, Tech Stack, Executable Commands, and Boundaries sections.
4. For model-specific overrides: verify it translates, not duplicates, the AGENTS.md rules into model-appropriate syntax.
