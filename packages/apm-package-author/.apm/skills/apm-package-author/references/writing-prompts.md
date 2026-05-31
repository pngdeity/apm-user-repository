# Writing Prompt Files

## Overview

Prompts are saved user commands — named prompt templates that the user invokes
in their agent harness. They live in `.apm/prompts/` with the `.prompt.md`
extension.

## File Format

### Frontmatter

```yaml
---
description: What this prompt does
input:
  - param_name: "Description of the input parameter"
allowed-tools: [Bash, Read, Grep]
---
```

| Field           | Required | Notes                                                         |
| --------------- | -------- | ------------------------------------------------------------- |
| `description`   | Yes      | Purpose of the prompt                                         |
| `input`         | No       | Named input parameters, referenced in body as `${input:name}` |
| `allowed-tools` | No       | List of tools the prompt can access                           |

### Body

Markdown with `${input:name}` placeholders for dynamic values.

## Full Example

```markdown
---
description: Review a pull request against our coding standards.
input:
  - pr_url: "URL of the PR to review"
  - focus: "Optional focus area (e.g. security, perf)"
allowed-tools: [Bash, Read, Grep]
---

# Review PR ${input:pr_url}

You are reviewing the changes in ${input:pr_url}. Focus on ${input:focus} when
set; otherwise apply the full checklist.

1. Fetch the diff.
2. Flag any deviation from `.github/CONTRIBUTING.md`.
3. Summarize blockers, suggestions, and nits in three sections.
```

## Deployment

Prompts deploy to target-specific directories:

- **Copilot**: `.github/prompts/<name>.prompt.md`
- **Claude**: `.claude/commands/<name>.md` (compiled format)
- **OpenCode**: Native prompt support
- **Windsurf**: Compiled to Cascade format

## When to Use Prompts vs Skills

| Characteristic | Prompts                             | Skills                              |
| -------------- | ----------------------------------- | ----------------------------------- |
| Invocation     | User types command name             | Agent matches description to intent |
| Input          | Explicit `${input:name}` parameters | Implicit from conversation context  |
| Best for       | Repeatable command templates        | Multi-step procedural guides        |
| Complexity     | Low — single invocation             | High — multi-step workflow          |

Use `prompts` for saved command templates (e.g., "review this PR", "generate
release notes"). Use `skills` for procedural guidance that requires multiple
steps and decision points.
