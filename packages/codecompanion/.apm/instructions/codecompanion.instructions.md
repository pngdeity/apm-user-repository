---
description: CodeCompanion.nvim quick reference for AI agents. Use when working with or configuring the CodeCompanion.nvim Neovim plugin.
applyTo: "**"
---

# CodeCompanion.nvim — AI Agent Quick Reference

> AI coding plugin for Neovim. Connect to LLMs (OpenAI, Anthropic, Copilot, Gemini, Ollama, etc.) and CLI agents (Claude Code, Codex, OpenCode, etc.) for chat, inline edits, and agentic tool use — all within Neovim buffers.

For full details, see the accompanying `references/full-docs.md` file.

---

## Core Commands

| Command | Purpose |
|---|---|
| `:CodeCompanion <prompt>` | Inline interaction — write/refactor code directly into a buffer |
| `:CodeCompanion /<alias>` | Call a prompt library item by its alias (e.g. `/explain`, `/tests`) |
| `:CodeCompanion adapter=<name> <prompt>` | Inline with a specific adapter |
| `:CodeCompanionChat` | Open a chat buffer for turn-based LLM conversation |
| `:CodeCompanionChat <prompt>` | Open chat and send prompt immediately |
| `:CodeCompanionChat Toggle` | Toggle visibility of the chat buffer |
| `:CodeCompanionChat adapter=<name> model=<model>` | Open chat with specific adapter/model |
| `:CodeCompanionChat Add` | Add visual selection to current chat buffer |
| `:CodeCompanionCLI` | Open a CLI agent interaction (terminal-based agents) |
| `:CodeCompanionCLI <prompt>` | Send prompt to the last CLI interaction |
| `:CodeCompanionCLI! <prompt>` | Auto-submit prompt to CLI agent, stay in current buffer |
| `:CodeCompanionCLI agent=<name> <prompt>` | Use a specific CLI agent |
| `:CodeCompanionCLI Ask` | Open rich prompt input buffer for CLI |
| `:CodeCompanionCmd` | Generate a command in the Neovim command-line |
| `:CodeCompanionActions` | Open the Action Palette (chat, prompts, workflows) |
| `:CodeCompanionActions refresh` | Reload markdown prompts from disk |

---

## Key Concepts

### Interactions
Five ways to interact with LLMs/agents via the plugin:
- **Chat** (`CodeCompanionChat`) — Turn-based buffer conversation. Supports editor context (`#`), slash commands (`/`), and tools (`@`)
- **CLI** (`CodeCompanionCLI`) — Terminal wrapper around agent CLI tools (Claude Code, Codex, etc.). Shares Neovim context via `@path` references
- **Inline** (`CodeCompanion`) — Prompts that write code directly into the current buffer, replacing selections or creating new buffers
- **Cmd** (`CodeCompanionCmd`) — Generate commands in the Neovim command-line
- **Background** — Async tasks like compacting messages or generating chat titles

### Adapters
Adapters connect Neovim to LLMs or agents. Two types:
- **HTTP adapters** — Connect to LLM APIs (OpenAI, Anthropic, Copilot, Gemini, Ollama, etc.)
- **ACP adapters** — Connect to agents via the Agent Client Protocol (Claude Code, Codex, Gemini CLI, OpenCode, etc.)

Each interaction can use a different adapter.

### Editor Context (`#`)
Dynamically inject Neovim state into prompts using `#{context}` syntax. Available in chat buffer, inline, and CLI. Built-in contexts: `#buffer`, `#buffers`, `#selection`, `#diagnostics`, `#diff`, `#messages`, `#quickfix`, `#terminal`, `#viewport`.

### Slash Commands (`/`)
Insert context into the chat buffer by typing `/`. Commands: `/buffer`, `/file`, `/symbols`, `/fetch`, `/compact`, `/fork`, `/help`, `/image`, `/rules`, `/mcp`, `/mode`, `/now`, `/resume`.

### Tools (`@`)
Let LLMs execute functions: read/write files, run commands, search code, etc. Built-in tools: `run_command`, `insert_edit_into_file`, `read_file`, `grep_search`, `file_search`, `create_file`, `delete_file`, `memory`, `web_search`, `fetch_webpage`, `get_diagnostics`, `get_changed_files`, `ask_questions`. Tool groups: `@{agent}`, `@{files}`.

### Rules
Persistent instructions from files like `CLAUDE.md`, `AGENTS.md`, `.cursorrules` that are automatically loaded into chat buffers.

### Prompt Library
Reusable prompt templates (Lua or Markdown) accessible via the Action Palette, keymaps, or slash commands. Built-in: `/explain`, `/fix`, `/tests`, `/commit`, `/lsp`.

### Events/Hooks
42+ autocmd events fired during the plugin lifecycle (e.g., `CodeCompanionChatCreated`, `CodeCompanionRequestStarted`, `CodeCompanionToolsFinished`).

---

## Common Configuration Patterns

### Minimal Setup (copilot adapter — no API key needed)
```lua
{
  "olimorris/codecompanion.nvim",
  version = "^19.0.0",
  dependencies = { "nvim-lua/plenary.nvim" },
  opts = {},
}
```

### Set Adapter and API Key
```lua
require("codecompanion").setup({
  interactions = {
    chat = { adapter = { name = "anthropic", model = "claude-sonnet-4-20250514" } },
    inline = { adapter = "anthropic" },
  },
  adapters = {
    http = {
      anthropic = function()
        return require("codecompanion.adapters").extend("anthropic", {
          env = { api_key = "ANTHROPIC_API_KEY" },
        })
      end,
    },
  },
})
```

### Configure a CLI Agent (Claude Code)
```lua
require("codecompanion").setup({
  interactions = {
    cli = {
      agent = "claude_code",
      agents = {
        claude_code = { cmd = "claude", args = {}, description = "Claude Code CLI", provider = "terminal" },
      },
    },
  },
})
```

### Suggested Keymaps
```lua
vim.keymap.set({ "n", "v" }, "<C-a>", "<cmd>CodeCompanionActions<cr>", { noremap = true, silent = true })
vim.keymap.set({ "n", "v" }, "<LocalLeader>a", "<cmd>CodeCompanionChat Toggle<cr>", { noremap = true, silent = true })
vim.keymap.set("v", "ga", "<cmd>CodeCompanionChat Add<cr>", { noremap = true, silent = true })
vim.cmd([[cab cc CodeCompanion]])
```

### Set up Rules (CLAUDE.md, AGENTS.md, etc.)
```lua
require("codecompanion").setup({
  rules = {
    opts = {
      chat = { autoload = "default", enabled = true },
    },
  },
})
```

### MCP Server Configuration
```lua
require("codecompanion").setup({
  mcp = {
    servers = {
      ["tavily-mcp"] = { cmd = { "npx", "-y", "tavily-mcp@latest" } },
    },
    opts = { default_servers = { "tavily-mcp" } },
  },
})
```

---

## Section Index

For the full reference, load `references/full-docs.md`. Key sections:

- [Features](references/full-docs.md) — LLM support, agents, tools, MCP, editor context
- [Installation](references/full-docs.md) — requirements, lazy.nvim, completion
- [Architecture](references/full-docs.md) — context management in the chat buffer
- [HTTP Adapters](references/full-docs.md) — models, env vars, proxy, community adapters
- [ACP Adapters](references/full-docs.md) — session config, MCP servers, per-agent setup
- [Chat Buffer](references/full-docs.md) — adapter, completion, keymaps, context, diff, tools, UI
- [CLI Configuration](references/full-docs.md) — agents, providers, keymaps
- [MCP Servers](references/full-docs.md) — roots, tool overrides, defaults
- [Prompt Library](references/full-docs.md) — markdown prompts, placeholders, workflows, pickers
- [Rules Configuration](references/full-docs.md) — enabling, groups, parsers, autoload
- [Extending](references/full-docs.md) — adapters, tools, rules parsers, UI integrations
