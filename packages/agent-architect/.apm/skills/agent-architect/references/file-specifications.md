# File Definitions & Specifications

## 1. AGENTS.md (The Operating Manual)

- **What it is:** The central configuration file defining the agent's persona, tech stack, file structure, code style, and absolute boundaries.
- **Where to use it:** At the root of the project repository or the main agent directory.
- **When to use it:** When scaffolding a new agent or defining global rules of engagement that must be active in the system prompt 100% of the time.
- **Why to use it:** To turn a vague generalist into a highly constrained specialist.
- **Required Components:**
  - **Persona:** A specific role (e.g., "You are an infrastructure agent managing an OpenTofu Shared Core deployment").
  - **Tech Stack:** Exact versions and frameworks (e.g., ".NET 8, Blazor WebAssembly, Python 3.11").
  - **Executable Commands:** Exact commands the agent will need frequently.
  - **Boundaries:** Explicit negative constraints (e.g., "Never modify database schemas without user approval"). Consult the GitHub Copilot source for examples of effective negative constraints.

## 2. SKILL.md (The Workflow Recipe)

- **What it is:** A portable, standardized package of domain expertise.
- **Where to use it:** Inside a dedicated skill folder, often accompanied by `references/` or `scripts/` directories.
- **When to use it:** When capturing repeatable workflows, multi-step processes, or specialized expertise that an agent needs to perform a specific task, but does not need in its baseline memory.
- **Why to use it:** To leverage *progressive disclosure*. The agent router will only read the YAML metadata at startup, loading the full instructions only when the workflow is triggered.
- **Required Components:**
  - **YAML Frontmatter:** Must include `name` (kebab-case) and `description`. *Crucial:* Write the description as an API doc; it must tell the router exactly the trigger conditions. Consult the AgentSkills and Anthropic sources for optimal trigger descriptions.
  - **Step-by-Step Instructions:** Numbered markdown steps specifying sequencing, error handling, and formatting.

## 3. Model-Specific Overrides (GEMINI.md, CLAUDE.md, etc.)

- **What it is:** A vendor-specific instruction file that overrides or translates the global `AGENTS.md` rules into the preferred "dialect" of a specific underlying model.
- **Where to use it:** At the root directory, alongside `AGENTS.md`.
- **When to use it:** When the repository is accessed by multiple different AI agents AND the global instructions need model-specific formatting.
- **Why to use it:** To optimize instruction adherence based on a model's unique training.
- **Required Components:**
  - **Dialect Translation:** If writing `CLAUDE.md`, format the global boundaries using `<rules>` and `<example>` XML tags. If writing `GEMINI.md`, use strict Markdown headers and bulleted constraint lists.
  - **Divergent Capabilities:** Explicit notes on what this specific model should or shouldn't do compared to the global baseline.
