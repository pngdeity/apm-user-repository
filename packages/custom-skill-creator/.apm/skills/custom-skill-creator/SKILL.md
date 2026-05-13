---
name: custom-skill-creator
description: Meta-skill for creating and packaging new AI agent skills. Use when requested to create a skill, document a multi-step process, or configure agent behavior as a procedural skill module. Not for single atomic functions or global architecture constraints — those go in AGENTS.md.
allowed-tools: skills-ref
metadata:
  tags: "meta skill-creation packaging create configure"
compatibility: Requires node and access to the agentskills.io spec. Designed for Gemini CLI or any agent runtime that supports SKILL.md discovery.
---

# Skill: Custom Skill Creator (`custom-skill-creator`)

This is the authoritative Meta-Skill for the `pngdeity` workspace. It extends the built-in skill creation process by enforcing the strict architectural boundaries of the "Abstraction Stack" and prioritizing validation and optimization based on `agentskills.io` standards.

## Usage
Run this skill whenever requested to "create a skill," "document a process," or "configure agent behavior." 

## Phase 1: The Decision Heuristic (The Gate)
Before writing any files, you MUST evaluate the user's request against this logic. **Do not create a SKILL.md if the request fails this gate.**

1.  **Global Constraint/Architecture?** -> Write it in `AGENTS.md`. (Refuse skill creation).
2.  **Single, Atomic Function?** -> This is a *Tool*. Suggest a code implementation. (Refuse skill creation).
3.  **Specialized, Multi-Step Workflow?** -> This is a *Skill*. Proceed to Phase 2.

## Phase 2: Skill Initialization & Drafting
1.  Create the skill directory manually:
    - `mkdir -p skills/<skill-name>/references skills/<skill-name>/scripts`
    - Create `SKILL.md` with the YAML frontmatter template below.
    - **If directory creation fails:** Check write permissions and parent directory existence. Fall back to a single `SKILL.md` file without subdirectories if needed.
2.  Read `references/agentskills-standards.md` for frontmatter field rules and `references/spec-authority.md` for the canonical specification sources (agentskills.io, Microsoft Agent Framework, Anthropic Claude).
3.  Draft the `SKILL.md` and any necessary scripts/references. Keep the main body under 500 lines (Progressive Disclosure).
4.  Validate structure: `skills-ref validate ./skills/<skill-name>` (or `npx skills-ref validate` from the agentskills.io reference library).

## Phase 3: Mandatory Validation
A skill is not complete until its behavior is validated.
1.  Define explicit test cases or dry-run instructions within the `SKILL.md` (e.g., "Verification Step: Run script X and ensure output matches Y").
2.  If the skill relies on executable `scripts/`, you MUST execute them locally to ensure they return LLM-friendly stdout (no massive tracebacks).

## Phase 4: Optimization & Security (Post-Stability)
Once the `SKILL.md` is structurally sound and validated:
1.  **Trigger Optimization:** Review the YAML `description`. It must function as an API doc for the agent router. Rewrite it to define exact trigger conditions.
2.  **Constraint Hardening:** Ensure instructions use imperative mood and handle edge cases without hallucinating.
3.  **Regenerate skill-index.json:** Run `node skills/verification/generate-skill-index.cjs` to rebuild the catalog index with updated tags and compatibility tokens from all SKILL.md files.
4.  **Security Review:** Per the agentskills.io specification and Microsoft Agent Framework guidance (https://learn.microsoft.com/en-us/agent-framework/agents/skills):
    - **Review content:** Read all skill files — instructions must not attempt to bypass safety guidelines, exfiltrate data, or modify agent configuration.
    - **Verify scripts:** If `scripts/` contains executables, confirm their behavior matches stated intent. Run them in an isolated environment first.
    - **Check provenance:** Only package skills from trusted sources. Prefer skills with version control history and active maintenance.
    - **Audit trail:** Ensure the skill's `name`, `description`, and `compatibility` fields are accurate and would not mislead a router.
    - **If any security concern is found:** Do not package. Document the finding and recommend remediation.

## Packaging
Once Phase 4 is complete, make the skill discoverable:
1.  Validate: `skills-ref validate ./skills/<skill-name>` (or `npx skills-ref validate` if not installed globally).
    - **If `skills-ref` is unavailable:** Skip validation and note the gap. Proceed with structural verification (step 4) instead.
2.  Register for cross-client discovery: `ln -sf <absolute-path-to-skill> ~/.agents/skills/<skill-name>`
3.  Run `node skills/verification/skill-compliance-check.cjs` to confirm all checks pass.
