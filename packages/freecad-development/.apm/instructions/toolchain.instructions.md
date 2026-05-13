---
description: "Use uv for package management, install type stubs for mypy"
---
## Toolchain

This project uses **uv** for package management. Install dependencies with
`uv pip install`, not raw `pip`.

Install stubs for type checking (required for mypy to pass):

```bash
uv pip install types-Markdown
```
