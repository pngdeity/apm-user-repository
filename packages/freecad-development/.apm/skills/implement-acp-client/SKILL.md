---
name: implement-acp-client
description: Procedural guide for implementing the Agent Client Protocol (ACP) in host applications (IDEs, desktop apps) to enable AI agent integration. Use when adding support for ACP-compatible agents, implementing stdio or HTTP/WS transports, or exposing host-specific tools to an agent.
---

# Implement ACP Client

This skill provides a standardized workflow for integrating AI agents into desktop or web applications using the Agent Client Protocol (ACP).

## Prerequisites

- Official Python SDK: `pip install acp-sdk`
- Understanding of the protocol: See [references/acp-spec-summary.md](references/acp-spec-summary.md)

## Implementation Workflow

### 1. Define the Client Interface
Implement the `acp.interfaces.Client` protocol. This is where you map protocol requests (like `read_text_file`) to your application's logic.

- Use [references/python-sdk-patterns.md](references/python-sdk-patterns.md) for code snippets.
- **Extension Methods**: Use `ext_method` to expose custom host tools (e.g., `execute_command`, `get_editor_state`).

### 2. Bridge the Event Loop (GUI Apps)
If your application uses a non-asyncio event loop (like Qt, Tkinter, or FreeCAD), you must run the ACP SDK in a background thread.
- Use a `QThread` (Qt) or `threading.Thread` to run an `asyncio` loop.
- Use thread-safe signals/slots to communicate between the ACP background thread and the UI main thread.

### 3. Setup Transport
- **Local (Subprocess)**: Use `acp.spawn_agent_process(command, client_impl)`. This handles the `stdio` pipe and JSON-RPC lifecycle automatically.
- **Remote**: Use `AsyncSseClient` or similar for HTTP/WebSockets.

### 4. Inject Tool Schemas (The "System Prompt" Trick)
Agents need to know what custom tools are available. Since ACP doesn't have a mandatory tool discovery endpoint for extensions yet, you should:
1. Define your tools in a JSON-schema registry.
2. Generate a Markdown summary of these tools.
3. Prepend this summary as a "System Instruction" or "Preamble" to the first message sent to the agent in a session.

### 5. Permission Handling
For sensitive tools (file writes, script execution), always implement `request_permission`.
- Return `AllowedOutcome.ALLOW` or `DENY`.
- In GUI apps, this should trigger a modal confirmation dialog on the main thread.

## Pitfalls & Failure Shields
- **Mismatched SDKs**: Ensure you are using the protocol-generic `acp-sdk` (`agentclientprotocol/python-sdk`), not platform-specific versions unless intended.
- **Threading Safety**: Never call GUI methods or modify application state directly from the ACP background thread. Always marshal via signals.
- **Zombie Agents**: Ensure the agent subprocess is killed when the client application exits. Use `.stop()` on your thread and `proc.kill()` if necessary.
- **Feedback Loops**: When an agent executes a command, return BOTH stdout and stderr. This allows the agent to self-correct if a command fails.
