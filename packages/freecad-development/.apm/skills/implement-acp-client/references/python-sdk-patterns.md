# ACP Python SDK Patterns

Use the `acp-sdk` package (installed via `pip install acp-sdk`).

## 1. Implementing the Client Interface
```python
from acp.interfaces import Client
from acp.schema import AgentMessageChunk, TextContentBlock

class MyHostClient(Client):
    async def session_update(self, session_id: str, update, **kwargs):
        if isinstance(update, AgentMessageChunk) and isinstance(update.content, TextContentBlock):
            print(f"Agent says: {update.content.text}")

    async def read_text_file(self, path: str, session_id: str, **kwargs):
        # Implementation
        pass
```

## 2. Spawning a Local Agent
```python
from acp import spawn_agent_process

# command: list of strings (e.g., ["gemini", "--experimental-acp"])
conn, proc = await spawn_agent_process(command, MyHostClient())
```

## 3. Custom Extension Methods (Tools)
Agents call `ext_method`. You must route these to your host tools.

```python
async def ext_method(self, method: str, params: dict, session_id: str, **kwargs):
    if method == "my_custom_tool":
        return self.handle_my_tool(params)
    raise RequestError.method_not_found(f"Method {method} not found")
```

## 4. Handling Permissions
```python
async def request_permission(self, tool: str, params: dict, session_id: str, **kwargs):
    # logic to ask user (via signal or dialog)
    if user_approved:
        return RequestPermissionResponse(outcome=AllowedOutcome.ALLOW)
    return RequestPermissionResponse(outcome=AllowedOutcome.DENY)
```
