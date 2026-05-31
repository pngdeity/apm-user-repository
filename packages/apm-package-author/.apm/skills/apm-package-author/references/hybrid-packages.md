# Hybrid Package Design

## Overview

Hybrid packages contain multiple primitive types in one `.apm/` directory. The
`type: hybrid` field in `apm.yml` signals that the package ships two or more
primitive categories.

## When to Use Hybrid

Use `hybrid` when a single logical package needs both:

- Always-on behavior rules AND on-demand skills
- Example: `microsoft-rust-guidelines` — safety rules as implicit
  `instructions`, design guidelines as summonable `skill`
- Example: `development-practices` — coding standards as `instructions`,
  conventional commits as `skill`

Use separate single-type packages when the concerns are unrelated. Don't merge a
Python style guide and a Docker deployment skill into one hybrid just for
consolidation.

## Dual Description Model

In hybrid packages, two independently authored descriptions coexist:

| Location               | Purpose                                  | Audience                 |
| ---------------------- | ---------------------------------------- | ------------------------ |
| `apm.yml description`  | Marketplace listing, `apm search` output | Humans browsing packages |
| `SKILL.md description` | Agent router trigger                     | Agent runtime            |

These should be **complementary, not identical**:

- `apm.yml description`: Broader — "Microsoft's Pragmatic Rust Guidelines —
  safety rules plus design conventions"
- `SKILL.md description`: Narrower — "Use when designing new Rust APIs,
  splitting crates, or implementing FFI"

## Directory Structure

```text
packages/my-hybrid/
  apm.yml                          # type: hybrid
  .apm/
    instructions/
      critical-rules.instructions.md   # Always active
    skills/
      detailed-workflow/
        SKILL.md                       # Summoned on demand
        references/                    # Progressive disclosure
```

## Design Pattern: Context-Optimized Hybrid

The strongest use case for hybrid is **context optimization**:

1. Identify the catastrophic-failure rules (violation = UB, corruption, security
   breach)
2. Package those as `instructions` with tight `applyTo` globs
3. Package everything else as `skill` with `references/` for progressive
   disclosure

**Result**: ~2,600 tokens always loaded for safety, ~20,000 tokens deferred
until the agent needs them.

Example from `microsoft-rust-guidelines`:

```
instructions/           # ~290 lines, always active
  rust-safety.md        # M-UNSAFE, M-UNSOUND → applyTo: **/*.rs
  rust-panic.md         # M-PANIC-IS-STOP, M-PANIC-ON-BUG → applyTo: **/*.rs
  rust-lint.md          # M-LINT-OVERRIDE-EXPECT → applyTo: **/*.rs
  rust-debug-display.md # M-PUBLIC-DEBUG, M-PUBLIC-DISPLAY → applyTo: **/src/**/*.rs
skills/
  microsoft-rust-guidelines/
    SKILL.md            # Trigger: "designing Rust APIs, splitting crates..."
    references/         # 10 files, ~1,900 lines of on-demand detail
```

## Anti-Patterns

1. **Thin wrapper skills**: Don't create a SKILL.md that just says "read the
   instruction files." The skill must provide value beyond the instructions.
2. **Redundant content**: Don't duplicate rules between instructions and skill
   references. Pick one home.
3. **Inconsistent globs**: If a safety rule applies to `**/*.rs` in instructions
   but the skill reference only mentions library code, the gaps will cause
   drift.
