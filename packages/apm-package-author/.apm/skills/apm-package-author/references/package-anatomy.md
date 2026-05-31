# Package Anatomy Reference

## `apm.yml` Manifest Fields

| Field             | Required | Type           | Notes                                                                                   |
| ----------------- | -------- | -------------- | --------------------------------------------------------------------------------------- |
| `name`            | Yes      | string         | Package identifier, kebab-case                                                          |
| `version`         | Yes      | string         | SemVer (e.g., `0.4.2`)                                                                  |
| `description`     | No       | string         | Human-readable, shown in marketplace listings                                           |
| `author`          | No       | string         | Package author                                                                          |
| `license`         | No       | string         | SPDX identifier recommended                                                             |
| `type`            | No       | enum           | `instructions`, `skill`, `hybrid`, or `prompts`                                         |
| `targets`         | No       | list           | Harness slugs: `opencode`, `claude`, `copilot`, `cursor`, `codex`, `gemini`, `windsurf` |
| `includes`        | No       | string or list | `"auto"` (publish all primitives) or explicit repo paths                                |
| `dependencies`    | No       | map            | `apm:` and/or `mcp:` keys listing dependency specs                                      |
| `devDependencies` | No       | map            | Same shape as `dependencies`, excluded from `apm pack`                                  |
| `scripts`         | No       | map            | Name → shell command, runnable via `apm run <name>`                                     |

### Minimal example

```yaml
name: my-skill
version: 0.4.2
type: skill
includes: auto
```

### Full example

```yaml
name: my-package
version: 1.0.0
description: Code review skills for Python services
author: Jane Doe
license: MIT
type: hybrid
targets:
  - opencode
  - claude
includes: auto
dependencies:
  apm:
    - microsoft/apm-sample-package#v1.0.0
  mcp:
    - microsoft/azure-devops-mcp
devDependencies:
  apm:
    - my-org/internal-test-skills
scripts:
  start: echo "hello"
```

## On-Disk Layout

```text
packages/<name>/
  apm.yml
  .apm/
    instructions/
      style.instructions.md       # Instruction files (frontmatter: description + applyTo)
    skills/
      my-skill/
        SKILL.md                  # Required: skill entry point
        scripts/                  # Optional: executable helpers
        references/               # Optional: on-demand reference files
        assets/                   # Optional: templates, images
    prompts/
      review.prompt.md            # Prompt files (frontmatter: description + optional input)
    agents/
      specialist.agent.md         # Agent persona files
    hooks/
      pre-commit.hook.md          # Lifecycle handlers
    commands/
      deploy.command.md           # CLI-style commands
  README.md                       # Optional
```

## Type Constraints

| `type`         | Allowable `.apm/` subdirectories |
| -------------- | -------------------------------- |
| `instructions` | `instructions/` only             |
| `skill`        | `skills/` only                   |
| `prompts`      | `prompts/` only                  |
| `hybrid`       | Any combination of the above     |

## Version Conventions

- SemVer: `MAJOR.MINOR.PATCH`
- `0.x.y` is pre-stable — breaking changes allowed on minor bumps
- `1.0.0+` — breaking changes require major bump
- This repo uses `0.4.2` as the lockstep version
- Per-package `version` must satisfy the root marketplace constraint (currently
  `">=0.0.1"`)
