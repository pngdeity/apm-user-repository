---
description: "Activate the project virtual environment before running any shell commands"
---
## Before any work: activate the venv

Always run this first before any shell command:

```bash
export PATH=".venv/bin:$PATH"
```

This ensures `mypy`, `ruff`, and `pytest` resolve from the project venv.
System-wide binaries may run against wrong site-packages and produce
spurious errors.
