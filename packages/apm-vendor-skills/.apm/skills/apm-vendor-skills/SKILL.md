---
name: apm-vendor-skills
description: Step-by-step procedural skill for vendoring external SKILL.md files into an APM marketplace. Use when importing skills, prompts, instructions, or agents from a non-APM GitHub repository (skills.sh, Claude Code plugins, or raw SKILL.md collections) into this marketplace.
allowed-tools: bash, Read, Write, Edit, Grep, Glob, webfetch
metadata:
  tags: "apm vendoring vendor import migration marketplace packaging"
compatibility: Requires apm CLI, just, and git. Works with any APM marketplace repo.
---

# APM Skill Vendoring

Vendors external SKILL.md content from a non-APM repository into this Agent
Package Manager marketplace. Assumes the source repo is **not** already set up
as an APM marketplace.

## Phase 1: Discovery

Explore the source repository to identify all vendorable content.

### Step 1.1: Map the repo structure

Fetch the GitHub tree or clone the repo. Identify:
- Every `SKILL.md` file (the primary primitive)
- Every `references/*.md` file (progressive disclosure)
- Every `scripts/*.sh` file (automation)
- Every `assets/*` file (templates)
- Any `prompts/*.prompt.md`, `instructions/*.instructions.md`, or `agents/*.agent.md` files

### Step 1.2: Understand the naming conventions

Note how skills are named, how they cross-reference each other, and what
frontmatter fields they use. Common source layouts:

| Source Layout | Key Files |
|---------------|-----------|
| `skills/<name>/SKILL.md` | skills.sh convention |
| `claude/<name>/skills/<name>/SKILL.md` | Claude Code plugin structure |
| `SKILL.md` at repo root | Single-skill repo |

## Phase 2: Mapping

Map source locations to APM `.apm/` directory conventions.

| Source Path | APM Path |
|-------------|----------|
| `skills/<name>/SKILL.md` | `packages/<name>/.apm/skills/<name>/SKILL.md` |
| `skills/<name>/references/*.md` | `packages/<name>/.apm/skills/<name>/references/*.md` |
| `skills/<name>/scripts/*.sh` | `packages/<name>/.apm/skills/<name>/scripts/*.sh` |
| `skills/<name>/assets/*` | `packages/<name>/.apm/skills/<name>/assets/*` |
| `skills/<name>/evals/` | `packages/<name>/.apm/skills/<name>/evals/` |

For multi-skill packages (hybrid type), group related skills under
`packages/<bundle>/.apm/skills/<skill-name>/`.

## Phase 3: Frontmatter Normalization

Adapt upstream YAML frontmatter to this marketplace's conventions.

### Fields to keep
- `name` — always preserved
- `description` — preserved; may be shortened for the `apm.yml` entry
- `compatibility` — preserved if present
- `metadata` — preserve the block, add a `tags` field

### Fields to add
- `metadata.tags` — space-separated keywords derived from `metadata.sources` or
  skill content. Use lowercase, use spaces as separators.

Example:
```yaml
# From upstream:
metadata:
  sources: "Effective Go, Google Style Guide, Uber Style Guide"

# Normalize to:
metadata:
  tags: "go style effective-go google-style uber-style"
```

### Fields to drop
- `license` — not used in this marketplace's frontmatter convention
- `metadata.sources` — converted to `metadata.tags`

### Tool name conversion
Platform-specific tool syntax must be normalized:

| Upstream Syntax | OpenCode Syntax |
|-----------------|-----------------|
| `Bash(bash:*)` | `bash` |
| `Read, Grep` | `Read, Grep` (already compatible) |

### Platform-specific features
Remove or adapt features that don't exist in the target platform:
- `!`command`` (Claude Code inline shell execution) → Replace with instructions to use the `bash` tool
- `AskUserQuestion` (Claude Code) → Replace with `Use the question tool`
- Junie/Claude Code-specific directives → Remove

## Phase 4: Cross-References

Upstream skills often cross-reference each other with relative paths:

```markdown
See [go-error-handling](../go-error-handling/SKILL.md) when...
```

Options:
1. **Remove them** — Simplest. Skills become self-contained.
2. **Convert to name references** — Change to "See the go-error-handling skill."
3. **Keep as relative links** — Works within multi-skill hybrid packages since
   all skills live under the same `.apm/skills/` directory.

For skills bundled into a single hybrid package, option 3 works. For skills
split across separate packages, options 1 or 2 are safer.

## Phase 5: Packaging

Create the APM package structure.

### Single-skill package (type: skill)
```
packages/<name>/
  apm.yml
  .apm/skills/<name>/
    SKILL.md
    references/  (optional)
    scripts/     (optional)
    assets/      (optional)
```

### Multi-skill bundle (type: hybrid)
```
packages/<bundle>/
  apm.yml
  .apm/skills/
    <skill-1>/
      SKILL.md
      references/
    <skill-2>/
      SKILL.md
      references/
```

### apm.yml template
```yaml
name: <package-name>
version: 0.1.0
description: <one-line description>
type: <skill | hybrid | instructions>
includes: auto
```

## Phase 6: Registration

Add the package to the root `apm.yml` under `marketplace.packages`.

### Entry template
```yaml
    - name: <package-name>
      description: <same as package description>
      source: pngdeity/apm-user-repository
      subdir: packages/<package-name>
      version: ">=0.0.1"
```

### Placement rules
- Insert in alphabetical order among existing packages
- Bundle packages (e.g., `pngdeity-defaults`) always go last

## Phase 7: Validation

Run these commands in order:

```bash
just fmt                  # Auto-format all YAML files
just sync-check           # Verify package versions satisfy root constraints
just sync-fix             # Auto-fix mismatched versions (if sync-check fails)
just yaml-validate        # Validate against JSON Schema
apm marketplace check     # Validate all refs resolve
apm pack --dry-run        # Validate marketplace.json generation
just ci                   # Full CI check (read-only)
```

Fix any failures before committing.

## Phase 8: Commit

```bash
apm pack                           # Regenerate marketplace.json
git add packages/<name>/ apm.yml .claude-plugin/
git commit -S -m "vendor: add <name> from <source-url>"
```
