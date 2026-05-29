---
description: Microsoft Rust guideline for lint overrides — use #[expect] not #[allow] to prevent stale lint accumulation. Enforced on all Rust source files.
applyTo: "**/*.rs"
---

# Rust Lint Override Guidelines (Microsoft Pragmatic Rust)

## Lint Overrides Should Use `#[expect]` (M-LINT-OVERRIDE-EXPECT)

When overriding project-global lints inside a submodule or item, you should do
so via `#[expect]`, not `#[allow]`.

Expected lints emit a warning if the marked warning was not encountered, thus
preventing the accumulation of stale lints. That said, `#[allow]` lints are
still useful when applied to generated code, and can appear in macros.

Overrides should be accompanied by a `reason`:

```rust,edition2021
#[expect(clippy::unused_async, reason = "API fixed, will use I/O later")]
pub async fn ping_server() {
  // Stubbed out for now
}
```
