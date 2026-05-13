---
description: Gemini CLI maintenance guide — managing the context library, progressive disclosure model, drop-in workflow for enabling skills in projects
applyTo: "**/.gemini/**"
---
# Gemini CLI: Expertise & Maintenance Guide

This document summarizes the "On-Demand Expertise" architecture and maintenance strategies established during your configuration tune-up.

---

## 1. Managing Your Context Library (Skills & Agents)

To maintain high-quality context without bloating your session, follow the "Progressive Disclosure" model.

### Structure for Local Expertise
Organize your library into discrete folders. Each folder is a "Skill."

```text
/path/to/my-library/
├── rails-audit/
│   └── SKILL.md        # Description: "Expert guidance for auditing Rails applications."
├── dotnet-best-practices/
│   └── SKILL.md        # Description: "Official .NET 9 architectural patterns."
└── agents/
    └── security-pro.md # A specialized sub-agent for deep security reviews.
```

### The "Drop-In" Workflow
When you enter a new project and want to enable a specific part of your library:
1. **Link the Skill:** `gemini skills link /path/to/my-library/rails-audit --scope workspace`
2. **Link the Agents:** `gemini skills link /path/to/my-library/agents --scope workspace`
3. **Set Trust:** Ensure the project is trusted (`gemini trust`).
4. **Use it:** The AI will now "see" these skills in the `/skills list` and can activate them on-demand.

### Remote Expertise (Microsoft .NET)
To pull in external, modular expertise:
`gemini skills install https://github.com/microsoft/dotnet-agents.git --scope workspace`

---

## 2. Maintaining a Forked Extension (e.g., GitHub Extension)

If an extension author neglects their `gemini-extension.json` metadata (causing "v1.0.0 to 1.0.0" update loops), follow this workflow to keep your fork clean and functional.

### The Problem
- **Git** tracks code changes (the "truth").
- **`gemini-extension.json`** tracks version numbers (often stale).
- **`gemini extensions update`** pulls code but reports version based on the JSON.

### The Solution: The "Shadow Manifest" Strategy
1. **Fork the Repo:** Fork `github/github-mcp-server` to your own account.
2. **Local Patch:** 
   - Add your missing `settings` block for `GITHUB_MCP_PAT` (with `sensitive: true`).
   - Manually increment the `version` in `gemini-extension.json` to something higher (e.g., `1.1.0`).
3. **Install from Fork:**
   `gemini extensions install https://github.com/YOUR_USER/github-mcp-server`
4. **Upstream Sync:**
   - Periodically sync your fork with the official `upstream/main`.
   - After syncing, **manually bump your version** in the fork. This forces the CLI to recognize an "actual" update and provides you with a clean update log.

---

## 3. Global Hygiene & Git Protection

### Global Git Ignore
Prevent `.gemini/` from being accidentally committed project-wide:
1. `touch ~/.gitignore_global`
2. `echo ".gemini/" >> ~/.gitignore_global`
3. `git config --global core.excludesfile ~/.gitignore_global`

### Modular Mandates
Your `GEMINI.md` is now a "Header" file. To add new rules:
1. Create `~/.gemini/mandates/new-rule.md`.
2. Add `@./mandates/new-rule.md` to `~/.gemini/GEMINI.md`.
