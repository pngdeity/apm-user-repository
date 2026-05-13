---
description: The HANDOFF.md to SKILL.md to AGENTS.md refinement pipeline for AI context files
applyTo: "**"
---

# Context Refinement Pipeline

AI context files follow a maturation lifecycle.

## HANDOFF.md

Created during an interactive session. Captures a specific, single-use task, its context, decisions made, and resolution. Ephemeral by nature.

## SKILL.md

Generalized from similar handoffs. Encodes a repeatable procedural workflow with triggers, steps, and verification. Lives in a skill directory with optional scripts, references, and assets.

## AGENTS.md

Anticipatory entity definition. Rather than responding to new situations by updating a SKILL.md, the agent can make decisions as if a skill had already been created for the situation. Encodes identity, boundaries, conventions, and always-on rules.

## When to Elevate

- A HANDOFF.md has been re-used with minor variations → distill its core pattern into a SKILL.md
- A SKILL.md has been refined through enough use cases that novel situations no longer require its activation → encode the governing principles into AGENTS.md
