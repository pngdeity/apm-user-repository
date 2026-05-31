---
description: Rust development defaults bundle — this package aggregates the microsoft-rust-guidelines package. Installed as a transitive dependency via rust-defaults.
applyTo: "**/*.rs"
---

# Rust Development Defaults

This bundle installs the `microsoft-rust-guidelines` package, which provides:

- Safety-critical rules enforced implicitly: unsafe code discipline, soundness
  guarantees, panic behavior, lint overrides, Debug/Display implementations
- An invocable skill covering universal Rust conventions, library design
  (interop, UX, resilience, building), applications, FFI, performance,
  documentation, and AI-oriented design

When you install `rust-defaults`, your agent follows Microsoft's Pragmatic Rust
Guidelines automatically for every `.rs` file.
