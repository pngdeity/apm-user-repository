---
description: "Run linting, type checking, formatting, and tests via make targets"
---
## Quality checks

```bash
make all          # lint, typecheck, format-check, test
make lint         # ruff check .
make typecheck    # mypy freecad/
make test         # pytest tests/ -v
make format       # ruff format .
make format-check # ruff format --check .
```
