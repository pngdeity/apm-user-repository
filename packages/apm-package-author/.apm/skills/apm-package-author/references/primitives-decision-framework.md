# Primitive Selection Decision Framework

## Primitive Types at a Glance

| Primitive        | Invocation                                    | Scoping                                          | Consumer Surface                                                                   |
| ---------------- | --------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------------------------- |
| **instructions** | Implicit — always loaded                      | `applyTo` glob per file pattern                  | Compiled into AGENTS.md, target-specific rules files                               |
| **skill**        | Explicit — agent summons by description match | Natural language trigger in SKILL.md frontmatter | `.agents/skills/<name>/SKILL.md` (cross-client) or `.claude/skills/` (Claude Code) |
| **prompts**      | Manual — user types command name              | Command name                                     | `.github/prompts/` or target-specific paths                                        |
| **agents**       | Explicit — user summons specialist persona    | Agent name                                       | `.github/agents/`, `.claude/agents/`, `.opencode/agents/`                          |
| **hooks**        | Event-driven                                  | Lifecycle event                                  | Per-target hook deployment                                                         |
| **commands**     | Manual — CLI-style invocation                 | Command name                                     | Per-target command deployment                                                      |

## Decision Matrix by Content Nature

| Content nature                                                    | Best primitive |
| ----------------------------------------------------------------- | -------------- |
| Style guides, coding standards, "always do X" rules               | `instructions` |
| Multi-step procedural workflows (how to deploy, how to review)    | `skill`        |
| Saved terminal commands, prompt templates with variables          | `prompts`      |
| Persona with constrained tool access (security auditor, reviewer) | `agents`       |
| Pre-commit validation, post-install hooks                         | `hooks`        |
| Custom CLI commands for APM projects                              | `commands`     |

## Context Budget Analysis

| Approach                                                      | Token cost                         | Risk                                        |
| ------------------------------------------------------------- | ---------------------------------- | ------------------------------------------- |
| All rules as `instructions`                                   | Every file match injects all rules | Context bloat for trivial edits             |
| All rules as `skill` + `references/`                          | Zero until summoned                | Agent may not recognize when to summon      |
| **Hybrid**: critical rules as `instructions`, rest as `skill` | Low implicit cost, rest on demand  | Requires correct critical/optional boundary |

## Invocation Model Comparison

### Instructions (Implicit)

- **Pros**: Always enforced, no summoning required, predictable
- **Cons**: Consumes context tokens on every file touch, no way to opt out
- **Best for**: Safety rules, compliance requirements, formatting conventions

### Skills (Explicit)

- **Pros**: Zero context until needed, progressive disclosure via references/,
  rich content (scripts, assets)
- **Cons**: Agent must correctly match description to user intent, can be missed
- **Best for**: Complex multi-step procedures, infrequent operations, deep
  reference material

## Scoping Mechanics

### Instructions use `applyTo` globs

- `**/*.rs` — all Rust files
- `**/src/**/*.rs` — library/application code only (excludes tests)
- `**/*.{py,js,ts}` — multiple extensions
- Applied at compile time — files outside the glob never see the rule

### Skills use description matching

- The SKILL.md `description` field is the agent router's API
- Must be specific enough to avoid false positives
- Must be broad enough to cover intended use cases
- Example: "Use when designing new Rust APIs, splitting crates, or implementing
  library UX patterns"

## Target Compatibility

| Primitive    | Copilot | Claude   | Cursor   | OpenCode | Codex    | Gemini          | Windsurf |
| ------------ | ------- | -------- | -------- | -------- | -------- | --------------- | -------- |
| instructions | native  | native   | native   | compiled | compiled | compiled        | compiled |
| skills       | native  | native   | native   | native   | native   | native          | native   |
| prompts      | native  | compiled | compiled | native   | native   | native          | compiled |
| agents       | native  | native   | native   | native   | compiled | **unsupported** | compiled |

Skills and instructions deploy universally. Agents do not deploy to Gemini CLI.

## HYBRID Dual Description Model

In `hybrid` packages:

- `apm.yml description` → human-facing marketplace listing text
- `SKILL.md description` → agent-runtime invocation matcher

These are independent and should not be identical. The `apm.yml` description can
be broader (describing the full package), while the SKILL.md description should
be narrower (specific trigger conditions).

## Decision Tree

```
Is the content...
├─ A policy rule that must always apply when editing certain files?
│  └─ instructions (with applyTo glob)
├─ A multi-step procedural workflow invoked by user intent?
│  └─ skill (with SKILL.md description as trigger)
├─ A saved user command or prompt template?
│  └─ prompts
├─ A specialized persona with scoped tool access?
│  └─ agents
├─ A lifecycle event handler?
│  └─ hooks
├─ Both always-on rules AND on-demand workflows in one package?
│  └─ hybrid (instructions + skill co-located)
└─ A custom CLI command for APM projects?
   └─ commands
```
