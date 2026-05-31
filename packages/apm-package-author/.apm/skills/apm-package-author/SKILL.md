---
name: apm-package-author
description: Meta-skill for creating standalone APM packages. Use when asked to create an APM package, author a new skill, write instruction files, scaffold an .apm/ directory, or create agent context primitives. Covers package anatomy, primitive selection (instructions vs skill vs hybrid vs prompts), scaffolding, and authoring. For marketplace registration and publishing, use the apm-marketplace-publisher skill.
allowed-tools: mkdir, apm
metadata:
  tags: "meta apm package-creation authoring primitives instructions skill hybrid prompts scaffolding"
compatibility: Requires apm CLI. Works with any APM-compatible project.
---

# Skill: APM Package Author

## When to Invoke

Invoke this skill when the user asks to:

- "create an APM package"
- "author a new skill" or "write instruction files"
- "scaffold an .apm/ directory"
- "create agent context primitives"
- "package a [skill/instruction/prompt/agent] for APM"

If the user asks to "add a package to the marketplace" or "publish an APM
package", redirect to the `apm-marketplace-publisher` skill instead.

## Phase 1: Choose Your Primitive

Before writing any files, classify the content against this decision tree.
Choose the single correct primitive per unit of content.

| Content nature                                              | Primitive        | Reasoning                                                                            |
| ----------------------------------------------------------- | ---------------- | ------------------------------------------------------------------------------------ |
| Always-on policy rules, style guides, conventions           | **instructions** | Implicit — compiled into AGENTS.md, active when files matching `applyTo` are touched |
| Multi-step procedural workflow, task guide                  | **skill**        | Explicit — SKILL.md with `description` as agent router trigger                       |
| Saved user command, prompt template                         | **prompts**      | Manual — user types command name                                                     |
| Both always-on rules AND on-demand workflows in one package | **hybrid**       | Combines instructions + skills (or other primitives)                                 |

**Gate questions** (from the custom-skill-creator heuristic):

1. Is this a global constraint or architecture rule? → Write it in `AGENTS.md`.
   Do not create a skill.
2. Is this a single atomic function? → This is a tool. Do not create a skill.
3. Is this a specialized, multi-step workflow? → This is a skill. Proceed.

For a detailed analysis of each primitive, see
`references/primitives-decision-framework.md`.

## Phase 2: Scaffold the Package

Create the directory with mandatory files:

```text
packages/<name>/
  apm.yml              # Required: package manifest
  .apm/                # Required: primitives directory
    instructions/      # .instructions.md files (if type: instructions or hybrid)
    skills/            # SKILL.md per skill (if type: skill or hybrid)
      <skill-name>/
        SKILL.md
        references/    # Optional: on-demand reference files
    prompts/           # .prompt.md files (if type: prompts)
```

### `apm.yml` fields

See `references/package-anatomy.md` for the complete field reference with types
and defaults.

Minimal valid `apm.yml`:

```yaml
name: my-package
version: 0.4.2
type: skill # or instructions, hybrid, prompts
includes: auto
```

### Type determines `.apm/` contents

| `type`         | Populate `.apm/` with                                                |
| -------------- | -------------------------------------------------------------------- |
| `instructions` | `.apm/instructions/*.instructions.md`                                |
| `skill`        | `.apm/skills/<name>/SKILL.md`                                        |
| `hybrid`       | Any combination of `instructions/`, `skills/`, `prompts/`, `agents/` |
| `prompts`      | `.apm/prompts/*.prompt.md`                                           |

## Phase 3: Write Primitives

Dispatch to the appropriate reference based on your chosen type:

| Primitive                     | Reference                            |
| ----------------------------- | ------------------------------------ |
| Instruction files             | `references/writing-instructions.md` |
| SKILL.md for skills           | `references/writing-skills.md`       |
| Prompt files                  | `references/writing-prompts.md`      |
| HYBRID dual-description model | `references/hybrid-packages.md`      |

### Progressive disclosure

Skills can use `references/` to keep the main SKILL.md body slim (aim for under
500 lines). Load references on demand when the agent needs depth. See
`references/writing-skills.md` for the pattern.

## Phase 4: Local Verification

Before proceeding to marketplace registration, run these two commands:

```bash
just fmt              # yamlfmt — formats apm.yml + packages/*/apm.yml
just yaml-validate    # check-jsonschema against schemas/apm.json
```

Resolve any failures. The package is now structurally valid and ready for
marketplace registration — use the `apm-marketplace-publisher` skill for that
phase.
