# Knowledge Base & Source Materials

When faced with ambiguity, or when determining the structural boundaries of these files, consult the following authoritative sources.

*Important: Do not just read the domain roots; target these specific technical guidelines based on the task at hand:*

## For AGENTS.md formatting and boundary setting

### GitHub Copilot Lessons
`https://github.blog/ai-and-ml/github-copilot/how-to-write-a-great-agents-md-lessons-from-over-2500-repositories/`

- **What it is:** Empirical data on writing repository instructions that LLMs actually obey.
- **When to use it:** When drafting the "Boundaries" section of `AGENTS.md` to learn how to frame effective negative constraints.

## For SKILL.md standardization, schema, and routing logic

### AgentSkills Open Standard

#### Specification Validation
`https://agentskills.io/specification`

- **What it is:** The strict YAML frontmatter and directory schema definitions required for a valid skill package.
- **When to use it:** When validating the structural integrity of a new `SKILL.md` file and ensuring its required metadata fields are correctly mapped.

#### Scaffolding Basics
`https://agentskills.io/skill-creation/quickstart`

- **What it is:** The fundamental boilerplate and minimum viable structure for generating a skill.
- **When to use it:** When generating the initial directory structure and foundational `SKILL.md` template for a brand new workflow.

#### Instruction Design
`https://agentskills.io/skill-creation/best-practices`

- **What it is:** Strategic guidelines on scoping, error handling, and formatting markdown steps for maximum LLM adherence.
- **When to use it:** When refining the step-by-step instructions of a skill to ensure it executes predictably and handles edge cases effectively without hallucinating.

#### Trigger Optimization
`https://agentskills.io/skill-creation/optimizing-descriptions`

- **What it is:** Tactical advice for writing high-fidelity YAML descriptions that function as API documentation for the agent router.
- **When to use it:** When drafting the YAML frontmatter to ensure the router knows exactly when to load this specific module.

#### Validation & Debugging
`https://agentskills.io/skill-creation/evaluating-skills`

- **What it is:** Methodologies for testing and iterating on skill reliability.
- **When to use it:** When defining test cases, dry-run instructions, or validation steps within the skill to confirm the agent achieved the desired state.

#### Executable Tool Integration
`https://agentskills.io/skill-creation/using-scripts`

- **What it is:** Instructions for embedding and invoking executable code (`scripts/`) alongside the declarative markdown instructions.
- **When to use it:** When a workflow requires executing a script to perform an atomic action alongside the agent's reasoning steps.

### Google ADK Guide

#### Architecture Baseline
`https://developers.googleblog.com/developers-guide-to-building-adk-agents-with-skills/`

- **What it is:** Google's architectural approach to separating declarative workflows from underlying code.
- **When to use it:** When evaluating whether a requested behavior is complex enough to warrant a standalone skill document versus a simple global instruction.

### Microsoft Agent Framework

#### Implementation Detail
`https://learn.microsoft.com/en-us/agent-framework/agents/skills?pivots=programming-language-csharp`

- **What it is:** Framework-specific details for packaging agent skills.
- **When to use it:** When formatting a `SKILL.md` workflow so it translates cleanly into actionable execution for an agent operating within a .NET environment.

### Anthropic Guide

#### Reasoning Design
`https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf`

- **What it is:** A deep dive into tool descriptions and multi-step reasoning.
- **When to use it:** When optimizing the step-by-step markdown instructions within a `SKILL.md` file for maximum model adherence.
