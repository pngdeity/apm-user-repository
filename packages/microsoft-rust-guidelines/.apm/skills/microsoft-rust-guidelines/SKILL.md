---
name: microsoft-rust-guidelines
description: Microsoft's Pragmatic Rust Guidelines for API design, crate architecture, library UX, resilience, FFI, performance optimization, documentation, and AI-oriented code. Invoke this skill when designing new Rust APIs, setting up crate structures, writing library code, implementing FFI, optimizing hot paths, authoring documentation, or making architectural decisions about Rust code organization.
allowed-tools: cargo, rustc, rustfmt, clippy, miri
metadata:
  tags: "rust guidelines api-design library-ux ffi performance documentation crate-architecture builder-pattern error-handling concurrency"
compatibility: Requires Rust toolchain (cargo, rustc). Works with any Rust project.
---

# Microsoft Pragmatic Rust Guidelines

## Description

Microsoft's comprehensive Rust coding guidelines covering library design, API
surface decisions, resilience patterns, application conventions, FFI,
performance, documentation, and AI-oriented code design. Safety-critical rules
(unsafe code, soundness, panic behavior, lint overrides, Debug/Display) are
enforced implicitly via instructions; all other guidelines are available on
demand through this skill.

## When to Invoke

Summon this skill when performing any of the following activities:

- **API Design** — designing public interfaces, choosing between generics/dyn
  traits, error type design, builder pattern
- **Crate Architecture** — splitting crates, feature flag design, dependency
  management
- **Library Resilience** — mocking I/O, avoiding statics, strong typing, test
  utilities
- **Bare-metal / Building** — -sys crates, OOBE compliance, build script design
- **FFI** — cross-DLL state management, native escape hatches, type portability
- **Applications** — allocator choice, application-level error handling
- **Performance** — hot path identification, throughput optimization, yield
  points
- **Documentation** — module docs, canonical doc sections, first sentence
  conventions
- **AI-Oriented Design** — making codebases agent-friendly

## Reference Files

Load the relevant reference for your specific task:

| Task                                                                            | Reference                          |
| ------------------------------------------------------------------------------- | ---------------------------------- |
| API design, traits, builders, AsRef, Sans-IO                                    | `references/library-ux.md`         |
| Mockable I/O, statics, strong types, glob reexports                             | `references/library-resilience.md` |
| Feature flags, OOBE, -sys crates                                                | `references/library-building.md`   |
| Type leakage, Send/Sync, escape hatches                                         | `references/library-interop.md`    |
| Naming conventions, magic values, crate splitting, static verification, logging | `references/universal.md`          |
| Error handling (anyhow/eyre), mimalloc                                          | `references/applications.md`       |
| DLL state isolation, portable types                                             | `references/ffi.md`                |
| Hot path profiling, throughput, yield points                                    | `references/performance.md`        |
| Module docs, canonical sections, doc(inline), first sentence                    | `references/documentation.md`      |
| Designing APIs for AI agent consumption                                         | `references/ai.md`                 |
