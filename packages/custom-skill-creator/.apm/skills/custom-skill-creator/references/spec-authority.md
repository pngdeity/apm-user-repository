# Authoritative Specification References

These are the canonical sources for the Agent Skills specification. All skill validation, creation, and packaging rules derive from these documents.

## Specification Sources

| Source | URL | Authority |
|---|---|---|
| **agentskills.io** | https://agentskills.io/ | Open specification for portable agent skills |
| **Microsoft Agent Framework — Agent Skills** | https://learn.microsoft.com/en-us/agent-framework/agents/skills | Authoritative implementation guide with field definitions, progressive disclosure pattern, security practices |
| **Anthropic Claude — Agent Skills** | https://platform.claude.com/docs/en/agents-and-tools/agent-skills | Claude Code implementation reference |
| **agentskills/agentskills** (GitHub) | https://github.com/agentskills/agentskills | Specification repo, issue tracker, community governance |

## Key Field Definitions (from Microsoft Agent Framework docs)

| Field | Required | Constraints |
|---|---|---|
| `name` | Yes | Max 64 chars. Lowercase letters, numbers, hyphens. Must match directory name. No leading/trailing hyphens. No consecutive hyphens. |
| `description` | Yes | Max 1024 chars. Should include keywords for task identification. |
| `license` | No | License name or reference to bundled license file. |
| `compatibility` | No | Max 500 chars. Environment requirements. |
| `metadata` | No | Arbitrary key-value mapping. |
| `allowed-tools` | No | Space-delimited list of pre-approved tools. Experimental. |

## Progressive Disclosure (from Microsoft Agent Framework docs)

1. **Advertise** (~100 tokens/skill) — Names and descriptions injected into system prompt
2. **Load** (<5,000 tokens recommended) — Full SKILL.md body loaded via `load_skill` tool
3. **Read resources** (as needed) — Supplementary files via `read_skill_resource`
4. **Run scripts** (as needed) — Script execution via `run_skill_script`

## Security Best Practices (from Microsoft Agent Framework docs)

- Review all skill content before deploying
- Verify scripts match stated intent
- Only install skills from trusted authors
- Sandbox script execution
- Maintain audit trails

## Verification Tools

| Tool | Source | Purpose |
|---|---|---|
| `skills-ref validate` | npm `skills-ref` v0.1.5 | Official structural validator per agentskills.io spec |
| `skill-compliance-check.cjs` | This project | Project-specific checks (description quality, placeholder pollution, compatibility) |
