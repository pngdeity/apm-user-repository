---
description: Microsoft Rust guidelines for panic behavior — panics mean immediate program termination, detected programming bugs panic rather than return errors. Enforced on all Rust source files.
applyTo: "**/*.rs"
---

# Rust Panic Guidelines (Microsoft Pragmatic Rust)

## Panic Means 'Stop the Program' (M-PANIC-IS-STOP)

Panics are not exceptions. Instead, they suggest immediate program termination.

Although your code must be
[panic-safe](https://doc.rust-lang.org/nomicon/exception-safety.html) (i.e., a
survived panic may not lead to inconsistent state), invoking a panic means _this
program should stop now_. It is not valid to:

- use panics to communicate (errors) upstream,
- use panics to handle self-inflicted error conditions,
- assume panics will be caught, even by your own code.

For example, if the application calling you is compiled with a `Cargo.toml`
containing

```toml
[profile.release]
panic = "abort"
```

then any invocation of panic will cause an otherwise functioning program to
needlessly abort. Valid reasons to panic are:

- when encountering a programming error, e.g., `x.expect("must never happen")`,
- anything invoked from const contexts, e.g., `const { foo.unwrap() }`,
- when user requested, e.g., providing an `unwrap()` method yourself,
- when encountering a poison, e.g., by calling `unwrap()` on a lock result (a
  poisoned lock signals another thread has panicked already).

Any of those are directly or indirectly linked to programming errors.

## Detected Programming Bugs are Panics, Not Errors (M-PANIC-ON-BUG)

As an extension of M-PANIC-IS-STOP above, when an unrecoverable programming
error has been detected, libraries and applications must panic, i.e., request
program termination.

In these cases, no `Error` type should be introduced or returned, as any such
error could not be acted upon at runtime.

Contract violations, i.e., the breaking of invariants either within a library or
by a caller, are programming errors and must therefore panic.

However, what constitutes a violation is situational. APIs are not expected to
go out of their way to detect them, as such checks can be impossible or
expensive. Encountering `must_be_even == 3` during an already existing check
clearly warrants a panic, while a function `parse(&str)` clearly must return a
`Result`. If in doubt, we recommend you take inspiration from the standard
library.

```rust,ignore
// Generally, a function with bad parameters must either
// - Ignore a parameter and/or return the wrong result
// - Signal an issue via Result or similar
// - Panic
// If in this divide_by we see that y == 0, panicking is the correct approach.
fn divide_by(x: u32, y: u32) -> u32 { ... }

// However, it can also be permissible to omit such checks
// and return an unspecified (but not an undefined) result.
fn divide_by_fast(x: u32, y: u32) -> u32 { ... }

// Here, passing an invalid URI is not a contract violation.
// Since parsing is inherently fallible, a Result must be returned.
fn parse_uri(s: &str) -> Result<Uri, ParseError> { };
```

> While panicking on a detected programming error is the 'least bad option',
> your panic might still ruin someone's day. For any user input or calling
> sequence that would otherwise panic, you should also explore if you can use
> the type system to avoid panicking code paths altogether.
