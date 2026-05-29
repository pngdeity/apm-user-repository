---
description: Microsoft Rust guidelines for Debug and Display trait implementations — all public types implement Debug (custom if sensitive), readable types implement Display. Enforced on library source files.
applyTo: "**/src/**/*.rs"
---

# Rust Debug and Display Guidelines (Microsoft Pragmatic Rust)

## Public Types are Debug (M-PUBLIC-DEBUG)

All public types exposed by a crate should implement `Debug`. Most types can do
so via `#[derive(Debug)]`:

```rust
#[derive(Debug)]
struct Endpoint(String);
```

Types designed to hold sensitive data should also implement `Debug`, but do so
via a custom implementation. This implementation must employ unit tests to
ensure sensitive data isn't actually leaked, and will not be in the future.

```rust
use std::fmt::{Debug, Formatter};

struct UserSecret(String);

impl Debug for UserSecret {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "UserSecret(...)")
    }
}

#[test]
fn test() {
    let key = "552d3454-d0d5-445d-ab9f-ef2ae3a8896a";
    let secret = UserSecret(key.to_string());
    let rendered = format!("{:?}", secret);

    assert!(rendered.contains("UserSecret"));
    assert!(!rendered.contains(key));
}
```

## Public Types Meant to be Read are Display (M-PUBLIC-DISPLAY)

If your type is expected to be read by upstream consumers, be it developers or
end users, it should implement `Display`. This in particular includes:

- Error types, which are mandated by `std::error::Error` to implement `Display`
- Wrappers around string-like data

Implementations of `Display` should follow Rust customs; this includes rendering
newlines and escape sequences. The handling of sensitive data outlined in
M-PUBLIC-DEBUG applies analogously.
