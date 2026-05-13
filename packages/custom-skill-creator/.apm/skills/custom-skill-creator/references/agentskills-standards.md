# AgentSkills.io Specification Standards

This document encapsulates the strict standards from the `agentskills.io` framework required for validating and optimizing `SKILL.md` behavior.

## 1. Specification Validation
The structural integrity of a new `SKILL.md` file MUST adhere to:
*   **YAML Frontmatter:** Must include exactly `name` (kebab-case) and `description`.
*   **Progressive Disclosure:** Keep the main `SKILL.md` lean. Move verbose context, large schemas, or variant-specific patterns into separate markdown files within `references/`.
*   **Directory Schema:** Must isolate logic into `scripts/` (executable code) and `references/` (declarative knowledge).

## 2. Instruction Design
*   **Scoping:** Define clear boundaries. If a workflow fails, define the recovery path explicitly (e.g., "If `X` fails, run `scripts/fallback.sh`").
*   **Formatting:** Use numbered markdown steps for sequential actions. Use imperative verbs ("Run", "Verify", "Create").
*   **Error Handling:** Never instruct the agent to "guess" or "improvise" around fragile operations.

## 3. Trigger Optimization
The YAML `description` functions as API documentation for the agent router.
*   **Specificity:** Include exact trigger conditions (e.g., "Use when handling CSV files for ETL pipeline X").
*   **Explicit Tool Assignment:** Every skill or sub-agent MUST explicitly declare its required tools using the YAML `allowed-tools` field (space-delimited string, per agentskills.io spec). Complex tool schemas must be explained in the Markdown body.
*   **Brevity:** It must be a single-line string.
*   **Slash Command Registration:** If a skill includes a standalone script that should be exposed to the CLI UX, create a corresponding `commands/*.toml` wrapper to register it as a native slash command.
*   **Context:** Tell the router *exactly* when this module is relevant so it isn't loaded accidentally during unrelated tasks.

## 4. Evaluating Skills (Validation)
*   **Test Cases:** Define dry-run instructions within the skill to confirm desired states.
*   **Debugging:** If the skill includes executable tools, they must output clear, concise success/failure messages to standard output.

## 5. Using Scripts (Executable Tool Integration)
*   **Atomic Actions:** When a workflow requires executing code (Python, Shell), place it in `scripts/`.
*   **Ergonomics:** Suppress standard tracebacks. Scripts must paginate or truncate large outputs (e.g., "Success: First 50 lines...") to prevent overflowing the agent's context window.