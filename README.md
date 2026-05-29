# pngdeity Agent Context

Curated APM marketplace of AI context files for software engineering.

Each package bundles primitives that configure AI coding agents — skills,
prompts, agents, and instructions — following the
[Agent Skills](https://agentskills.io) and [AGENTS.md](https://agents.md) open
standards. Declare what you need in your project's `apm.yml` and APM resolves,
fetches, and wires primitives into Gemini CLI, OpenCode, Copilot, Claude Code,
Cursor, and every other supported harness.

## Quickstart

```bash
apm marketplace add pngdeity/apm-user-repository
apm marketplace update
apm install agent-architect@pngdeity-agent-context
apm compile
```

## Packages

| Package                          | Type         | Description                                                                                                             |
| -------------------------------- | ------------ | ----------------------------------------------------------------------------------------------------------------------- |
| **agent-architect**              | hybrid       | Meta-tooling for AI context files — agent-architect skill, context-refinement, update-agents-md prompt, Confucius agent |
| **ai-agent-orchestration**       | skill        | AI agent orchestration design patterns, production best practices, multi-agent coordination                             |
| **apm-migration**                | skill        | 6-phase procedural skill for migrating consumer projects to APM                                                         |
| **architectural-review**         | skill        | ADR and TDD writing and auditing with structural, logic, and consistency review phases                                  |
| **aur-package-management**       | skill        | Arch Linux AUR package maintenance — update, patch recovery, bootstrap                                                  |
| **automated-modernization**      | instructions | ARME blueprint for intelligent repository modernization                                                                 |
| **ci-cd-standards**              | hybrid       | GitHub Actions best practices, IoC architecture, CNCF Buildpacks, vendor-decoupled CI/CD design                         |
| **code-review-commons**          | skill        | Common guidelines and constraints for performing high-quality code reviews                                              |
| **codecompanion**                | hybrid       | CodeCompanion.nvim quick reference and full documentation reference                                                     |
| **config-maintenance**           | hybrid       | Configuration file auditing, modernization, and Gemini CLI maintenance guide                                            |
| **configure-trim-safe-efcore**   | skill        | Configure EF Core for reflection-free compilation, Native AOT, and aggressive trimming                                  |
| **context-quality-gate**         | skill        | CI/CD quality gate for AI agent context files — structural compliance, trigger accuracy, output quality                 |
| **custom-skill-creator**         | skill        | Meta-skill for creating and packaging new AI agent skills                                                               |
| **design-process**               | hybrid       | Elon Musk's five-step engineering and design process                                                                    |
| **design-systems**               | hybrid       | Design system references — Hyperstudio and 099                                                                          |
| **development-practices**        | hybrid       | Conventional commits, technical integrity, engineering workflow, git and code standards                                 |
| **find-docs**                    | skill        | Retrieve up-to-date documentation for libraries, frameworks, and tools via Context7 API                                 |
| **freecad-development**          | hybrid       | FreeCAD ACP development — venv, toolchain, debugging, quality, ACP client implementation                                |
| **gh-cli-patterns**              | skill        | Patterns for invoking GitHub CLI (gh) from AI agents                                                                    |
| **implement-offline-first-sync** | skill        | Design offline-first sync logs with optimistic concurrency for SQLite-to-server conflicts                               |
| **init-project-guidance**        | skill        | Scaffold GEMINI.md and AGENTS.md for new or existing projects                                                           |
| **local-first**                  | skill        | Discover local information sources (man pages, --help, logs) before reaching for external docs                          |
| **manage-sqlite-wasm-lifecycle** | skill        | Manage SQLite database file lifecycles and DbContext connections in Blazor WASM sandboxes                               |
| **microsoft-rust-guidelines**    | hybrid       | Microsoft's Pragmatic Rust Guidelines — safety-critical rules enforced implicitly, full guidelines on demand            |
| **normalize-dataset**            | skill        | Pre-process raw documents into LLM-friendly Markdown datasets                                                           |
| **pngdeity-defaults**            | instructions | Default package bundle — aggregates development-practices, ci-cd-standards, and code-review-commons                     |
| **project-management**           | skill        | Project management integration — GitHub Projects, TODO.md, issue trackers                                               |
| **skill-eval-agents**            | hybrid       | Seven specialized sub-agent skills and Go CLI tools for autonomous skill evaluation                                     |
| **skill-eval-pipeline**          | hybrid       | Multi-agent orchestration pipeline for evaluating and optimizing SKILL.md files                                         |
| **ssh-gpg-host**                 | skill        | Host-specific SSH and GPG configuration constraints for agent commit signing                                            |
| **technical-documentation**      | skill        | ADR and Technical Design Document authoring with templates and checklists                                               |
| **technical-integrity**          | skill        | Maintain intellectual honesty and technical rigor during software engineering tasks                                     |
