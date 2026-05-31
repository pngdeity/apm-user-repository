# Writing SKILL.md Files

## Overview

Skills are self-contained capability bundles deployed to the agent's skills
directory. They are model-invoked — the agent reads the `description` at runtime
and summons the skill when it matches the user's intent.

## File Format

Skills live in `.apm/skills/<name>/SKILL.md`. The directory name must match the
`name` field in the frontmatter.

### Frontmatter

```yaml
---
name: skill-name
description: Precise trigger conditions — functions as API doc for the agent router
allowed-tools: tool1, tool2
metadata:
  tags: "kebab-case tags"
compatibility: Runtime requirements
---
```

| Field           | Required | Notes                                                   |
| --------------- | -------- | ------------------------------------------------------- |
| `name`          | Yes      | Must match the skill directory name                     |
| `description`   | Yes      | Agent router trigger — be specific about when to invoke |
| `allowed-tools` | No       | Comma-separated tool list                               |
| `metadata.tags` | No       | Discovery tags                                          |
| `compatibility` | No       | Runtime, OS, or tool requirements                       |

### Body Structure

```markdown
# Skill: <Title>

## Description

Brief overview of what this skill provides.

## When to Invoke

Explicit trigger phrases and scenarios.

## Instructions (or Phases)

### 1. First step...

### 2. Second step...

Each step should be imperative and actionable.
```

## Description as Trigger

The `description` field is the single most important field — it determines
whether the agent invokes the skill. Write it as a clear set of trigger
conditions:

- **Good**: "Use when asked to create a pull request, draft a PR description, or
  summarize uncommitted changes for review."
- **Good**: "Meta-skill for creating and packaging new AI agent skills. Use when
  requested to create a skill, document a multi-step process, or configure agent
  behavior."
- **Bad**: "PR helper" — too vague
- **Bad**: "Use this skill for everything related to pull requests, branches,
  commits, reviews, merges, GitHub, and git operations." — too broad, causes
  false positives

## Progressive Disclosure with `references/`

Skills can include a `references/` directory for on-demand content. The main
SKILL.md body should stay under 500 lines — forward the agent to references for
depth.

```text
skills/my-skill/
  SKILL.md              # ~100-200 lines: decision tree + reference index
  references/
    detailed-guide.md   # Deep detail loaded only when needed
    examples.md         # Usage examples
    templates/          # Config templates
```

In the SKILL.md body, dispatch to references:

```markdown
## Phase 2: Implementation

- For full API reference: see `references/api-reference.md`
- For deployment patterns: see `references/deployment-patterns.md`
```

## Full Example

```markdown
---
name: pr-description
description: >-
  Activate when the user asks for a pull-request description, a summary of
  uncommitted changes, or release notes. Use when preparing to open a PR.
---

# PR Description Skill

## Description

Produces a structured PR description from uncommitted changes.

## Instructions

### 1. Gather Changes

Run `git diff --stat` and `git log --oneline` to collect the changeset.

### 2. Draft Summary

Write a one-sentence summary of what changed and why.

### 3. Produce Sections

Generate the PR description with these sections:

- Summary, Motivation, Changes, Risk/Rollback, Testing

### 4. Output

Present the draft for user approval before posting.
```

## Validation

After writing, verify the skill structure:

```bash
skills-ref validate ./skills/<skill-name>
```

If `skills-ref` is unavailable, manually verify:

- [ ] `name` matches directory name
- [ ] `description` is specific and actionable
- [ ] Main body is under 500 lines
- [ ] Cross-references to `references/` files are correct paths
