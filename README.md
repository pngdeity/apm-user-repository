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

<!-- PACKAGE_TABLE_START -->
| Package | Type | Description |
|---------|------|-------------|
| **agent-architect** | hybrid | Meta-tooling for AI context files — agent-architect skill with knowledge sources and file specifications, context-refinement instruction, update-agents-md prompt, and the Confucius skill extraction agent. |
| **ai-agent-orchestration** | skill | AI agent orchestration design patterns, production best practices, and multi-agent coordination workflows — sequential, concurrent, group chat, handoff, and magentic patterns with decision matrix, production hardening checklist from Google AI Agent Clinic, and practical orchestration workflow from the Cyrus case study. |
| **apm-author-defaults** | instructions | APM package authoring bundle — create packages, design skills, and manage AI context files |
| **apm-eval-defaults** | instructions | APM skill evaluation bundle — evaluate SKILL.md quality, run multi-agent eval pipelines, and enforce CI quality gates |
| **apm-migration** | skill | 6-phase procedural skill for migrating consumer projects to APM — detection, tribal knowledge recovery, scaffolding, install/compile, git operations, and verification |
| **apm-package-author** | skill | Meta-skills for APM package lifecycle — standalone package creation (anatomy, primitive selection, authoring, local verification) and marketplace registration (root apm.yml integration, version constraints, publishing) |
| **architectural-review** | skill | ADR and TDD writing and auditing workflow with structural, logic, and consistency review phases. |
| **aur-package-management** | skill | Arch Linux AUR package maintenance skills — update packages to new upstream versions, recover from patch failures, and bootstrap new packages. |
| **automated-modernization** | instructions | Blueprint for the Automated Repository Modernization Engine (ARME) — intelligent discovery, staging strategy, and implementation framework. |
| **ci-cd-standards** | hybrid | CI/CD pipeline design standards — GitHub Actions best practices, Inversion of Control architecture, CNCF Buildpacks migration reference, and vendor-decoupled CI/CD pipeline design skill. |
| **code-review-commons** | skill | Common guidelines and constraints for performing high-quality code reviews |
| **codecompanion** | hybrid | CodeCompanion.nvim quick reference and full documentation reference for AI agents working with the Neovim plugin. |
| **config-maintenance** | hybrid | Configuration file auditing, modernization, and Gemini CLI maintenance guide. |
| **configure-trim-safe-efcore** | skill | Configure Entity Framework Core for reflection-free compilation, Native AOT compatibility, and aggressive assembly trimming in Blazor WebAssembly apps |
| **context-quality-gate** | skill | CI/CD quality gate skill for AI agent context files — enforces structural compliance via skills-ref, trigger accuracy thresholds, and output quality minimums before allowing skill publication through automated PR checks |
| **custom-skill-creator** | skill | Meta-skill for creating and packaging new AI agent skills. Use when creating a skill, documenting a multi-step process, or configuring agent behavior as a skill module. |
| **design-process** | instructions | Elon Musk's five-step engineering and design process for driving innovation and eliminating inefficiency. |
| **design-systems** | hybrid | Design system references — Hyperstudio (monochrome terminal with amber accents) and 099 (terminal aesthetic digital workbench). |
| **development-practices** | hybrid | Development standards — conventional commits, technical integrity, engineering workflow, git standards, and code style conventions. |
| **find-docs** | skill | Retrieves up-to-date documentation for any library, framework, or tool via Context7 API |
| **freecad-development** | hybrid | FreeCAD ACP development standards — virtual environment, toolchain, debugging, quality checks, and ACP client implementation. |
| **gh-cli-patterns** | skill | Patterns for invoking the GitHub CLI (gh) from agents — structured output, pagination, repo targeting, search vs list, and gh api fallback. |
| **implement-offline-first-sync** | skill | Design offline-first local data storage sync logs with optimistic concurrency to handle SQLite-to-server conflicts |
| **init-project-guidance** | skill | Scaffold GEMINI.md and AGENTS.md for new or existing projects. Detects tech stack and generates tailored AI agent guidance. |
| **local-first** | skill | Discover and use local information sources (man pages, --help, apropos, info, system logs, kernel params) before reaching for external documentation. Applies when the agent queries installed CLI tools, system utilities, or infrastructure commands. |
| **manage-sqlite-wasm-lifecycle** | skill | Safely manage SQLite database file lifecycles and DbContext connections during resets, imports, exports, and restores in Blazor WASM sandboxes |
| **microsoft-rust-guidelines** | hybrid | Microsoft's Pragmatic Rust Guidelines — safety-critical rules enforced implicitly (unsafe code, soundness, panics, lint overrides, Debug/Display), plus an invocable skill covering universal conventions, library design (interop, UX, resilience, building), applications, FFI, performance, documentation, and AI-oriented design |
| **normalize-dataset** | skill | Pre-process raw documents (PDFs, transcripts, legacy text) into semantic, LLM-friendly Markdown datasets for RAG and knowledge ingestion. |
| **pngdeity-defaults** | instructions | Default package bundle for pngdeity projects — aggregates development-practices, ci-cd-standards, and code-review-commons |
| **project-management** | skill | Project management integration skill for AI agents — interface with GitHub Projects, TODO.md, and other project management tools. |
| **rust-defaults** | instructions | Default Rust development bundle — Microsoft Pragmatic Rust Guidelines with safety-critical rules enforced implicitly and design guidelines on demand |
| **skill-eval-agents** | hybrid | Seven specialized sub-agent skills and Go CLI tools for autonomous skill evaluation — struct-validator, trigger-evaluator, trigger-aggregator, quality-evaluator, output-grader, revision-synthesizer, and candidate-selector with bundled Go scripts for CLI invocation, session parsing, benchmark computation, and candidate selection |
| **skill-eval-pipeline** | hybrid | Multi-agent orchestration pipeline for evaluating and optimizing AI agent SKILL.md files — structural validation, trigger testing across opencode and gemini CLI, output quality measurement with graded assertions, automated revision synthesis, and best-candidate selection |
| **ssh-gpg-host** | skill | Host-specific SSH and GPG configuration constraints for agent commit signing and authentication |
| **technical-documentation** | skill | Architecture Decision Record (ADR) and Technical Design Document authoring workflows with templates, verification checklists, and anti-pattern guidance. |
| **technical-integrity** | skill | Maintain intellectual honesty and technical rigor during software engineering tasks. Use when performing code reviews, responding to peer agent feedback, or defending technical decisions. |
<!-- PACKAGE_TABLE_END -->
