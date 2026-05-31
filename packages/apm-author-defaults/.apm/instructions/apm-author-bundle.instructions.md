---
description: APM package authoring bundle — aggregates packages for creating APM packages, designing skills, and managing AI context files. Installed as transitive dependencies via apm-author-defaults.
applyTo: "**"
---

# APM Authoring Defaults

This bundle installs the core APM package authoring toolkit:

| Package                     | Purpose                                                       |
| --------------------------- | ------------------------------------------------------------- |
| `apm-package-author`        | Create APM packages (anatomy, primitive selection, authoring) |
| `apm-marketplace-publisher` | Included transitively via `apm-package-author`                |
| `custom-skill-creator`      | Create and package AI agent skills                            |
| `agent-architect`           | Meta-tooling for AI context files                             |

When you install `apm-author-defaults`, your agent gains the full toolset for
creating APM packages, skills, and agent context files.
