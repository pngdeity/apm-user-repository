# Writing Instruction Files

## Overview

Instructions are long-lived behavior rules compiled into AGENTS.md and
target-specific rules files. They are always active when files matching the
`applyTo` glob are touched.

## File Format

Files live in `.apm/instructions/` with the `.instructions.md` extension. The
basename becomes the deployed filename stem.

### Frontmatter

```yaml
---
description: Rule description — what this instructs the agent to do
applyTo: "**/*.rs"           # Glob pattern binding the rule to files
---
```

- `description` (required): Explains the rule's purpose to the agent
- `applyTo` (required): Glob pattern determining which files this rule affects

### Body

Markdown content with rules, code examples, and checklists. Keep it declarative
and actionable.

## `applyTo` Glob Patterns

| Pattern           | Scope                                     |
| ----------------- | ----------------------------------------- |
| `**/*.rs`         | All Rust source files                     |
| `**/src/**/*.rs`  | Library/application code (excludes tests) |
| `**/*.{py,js,ts}` | Multiple extensions                       |
| `**/*`            | All files                                 |
| `docs/**/*.md`    | Documentation only                        |

Use the narrowest glob that covers the rule's scope. This minimizes context
injection for unrelated file types.

## Full Example

````markdown
---
description: Rust safety guidelines — unsafe implies UB risk, unsafe needs justification, all code must be sound
applyTo: "**/*.rs"
---

# Rust Safety Guidelines

## Unsafe Implies Undefined Behavior

The marker `unsafe` may only be applied to functions and traits if misuse
implies the risk of undefined behavior (UB). It must not be used to mark
functions that are dangerous to call for other reasons.

```rust
// Valid use of unsafe
unsafe fn print_string(x: *const String) { }

// Invalid use of unsafe
unsafe fn delete_database() { }
```
````

## All Code Must be Sound

Unsound code is seemingly _safe_ code that may produce undefined behavior.
Unsound abstractions are never permissible.

```
## Compilation Behavior

Instructions are compiled differently per target:
- **Copilot, Claude**: Written to `.github/instructions/` or `.claude/rules/` as native rules
- **OpenCode, Gemini, Codex**: Folded into AGENTS.md or equivalent context file
- **Windsurf**: Reformatted per Cascade conventions

The `applyTo` field dictates how targets scope the rule — some preserve the glob verbatim, others convert to target-specific scoping.

## Best Practices

1. One rule file per logical domain (e.g., `rust-safety.instructions.md`, not `all-rust-rules.instructions.md`)
2. Keep files focused — if a single instruction file exceeds ~150 lines, consider splitting by sub-topic
3. Include code examples for clarity — agents learn well from patterns
4. Use the narrowest `applyTo` glob possible — `**/src/**/*.rs` not `**/*` unless the rule is truly universal
5. Combine related rules into one file — safety guidelines belong together, not split across 5 tiny files
6. Reference source material in comments when rules derive from external standards
```
