# ACP Specification Summary

The Agent Client Protocol (ACP) standardizes how code editors (Clients) and AI agents (Agents) communicate.

## Roles
- **Client**: The host application (IDE, CAD software, etc.) where the user resides.
- **Agent**: The AI entity (local subprocess or remote service) that provides assistance.

## Transports
- **stdio**: Primary for local agents. Agent runs as a subprocess.
- **HTTP/WebSockets**: For remote agents (Sse/SseMcp).

## Message Lifecycle
1. **Initialize**: Exchange capabilities (read-only, terminal, etc.).
2. **Session**: Create a new session with an agent.
3. **Prompt**: Send user input (text/images) to the agent.
4. **Turns**: Agent responds with text chunks or tool calls.
5. **Tool Calls**: Agent requests to use a client capability (e.g., `read_text_file`).
6. **Extensions**: Clients can expose custom methods (`ext_method`).

## Client Responsibilities
- Manage agent lifecycle (spawn/kill).
- Render Markdown/Diffs.
- Provide host context (files, environment).
- Guard sensitive actions (permissions).
