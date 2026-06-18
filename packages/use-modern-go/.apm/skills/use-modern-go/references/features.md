# Modern Go Guidelines Explained
**Work in progress** — inconsistencies may be present.

This file provides a more detailed description of the features supported in the Go Modern Guidelines.

**Modernizer legend:**
- [x] — available in go fix as modernizer
- [ ] — not implemented

**Impact legend:**
- Critical — found in almost every project, dozens of occurrences
- High — found often, 5-20 occurrences per project
- Medium — found regularly, 1-5 occurrences per project
- Low — found rarely or in specific code

| Category | Diagnostic | Modernizer | Go | Impact |
|----------|------------|------------|-----|--------|
| Collections | `slicescontains` | [x] | 1.21 | Critical |
| Collections | `sortslice` | [x] | 1.21 | High |
| Collections | `minmax` | [x] | 1.21 | High |
| Collections | `mapkeysvalues` | [ ] | 1.23 | High |
| Collections | `mapsloop` | [x] | 1.21 | Medium |
| Collections | `mapsclone` | [ ] | 1.21 | Medium |
| Collections | `rangeoverindex` | [ ] | 1.21 | Medium |
| Collections | `slicesmaxmin` | [ ] | 1.21 | Medium |
| Collections | `slicesreverse` | [ ] | 1.21 | Medium |
| Collections | `clearcollection` | [ ] | 1.21 | Medium |
| Collections | `mapdeletefunc` | [ ] | 1.21 | Low |
| Collections | `slicescompact` | [ ] | 1.21 | Low |
| Strings | `stringscutprefix` | [x] | 1.20 | High |
| Strings | `stringsseq` | [x] | 1.24 | High |
| Strings | `stringsclone` | [ ] | 1.20 | Medium |
| Strings | `bytescut` | [ ] | 1.18 | Medium |
| Loops | `rangeint` | [x] | 1.22 | Critical |
| Loops | `forvar` | [x] | 1.22 | High |
| Types | `efaceany` | [x] | 1.18 | Critical |
| Types | `reflecttypefor` | [ ] | 1.22 | Low |
| Errors | `erris` | [ ] | 1.13 | Critical |
| Errors | `errorsjoin` | [ ] | 1.20 | High |
| Time | `timesince` | [ ] | 1.0 | High |
| Time | `timeuntil` | [ ] | 1.8 | Medium |
| Context | `testingcontext` | [x] | 1.24 | High |
| Context | `contextcause` | [ ] | 1.20 | Medium |
| Context | `contextafterfunc` | [ ] | 1.21 | Medium |
| Context | `contexttimeoutcause` | [ ] | 1.21 | Low |
| Sync | `waitgroup` | [x] | 1.25 | High |
| Sync | `synconcefunc` | [ ] | 1.21 | Medium |
| Sync | `atomicvalue` | [ ] | 1.19 | Medium |
| Fmt | `fmtappendf` | [x] | 1.19 | Medium |
| Testing | `bloop` | [x] | 1.25 | Medium |
| JSON | `omitzero` | [x] | 1.24 | Medium |
| Utilities | `cmpor` | [ ] | 1.22 | High |
| Utilities | `newliteral` | [x] | 1.26 | High |
| HTTP | `httpmux` | [ ] | 1.22 | Medium |

## 1. Collections (slices, maps)

### [x] 1.1 `slicescontains` — replace search loop with `slices.Contains`

**Go 1.21+ | Impact: Critical**

```go
// Before
found := false
for _, v := range s {
    if v == needle {
        found = true
        break
    }
}
// After
found := slices.Contains(s, needle)
```

### [x] 1.2 `sortslice` — replace `sort.Slice` with `slices.Sort`

**Go 1.21+ | Impact: High**

```go
// Before
sort.Slice(users, func(i, j int) bool {
    return users[i].Name < users[j].Name
})
// After
slices.SortFunc(users, func(a, b User) int {
    return cmp.Compare(a.Name, b.Name)
})
```

### [x] 1.3 `minmax` — replace conditional assignment with `min`/`max`

**Go 1.21+ | Impact: High**

```go
// Before
if a < b {
    result = a
} else {
    result = b
}
// After
result := min(a, b)
```

### [ ] 1.4 `mapsclone` — explicit use of `maps.Clone`

**Go 1.21+ | Impact: Medium**

```go
// Before
copy := make(map[string]int, len(m))
for k, v := range m {
    copy[k] = v
}
// After
copy := maps.Clone(m)
```

### [ ] 1.5 `mapkeysvalues` — use `maps.Keys` / `maps.Values`

**Go 1.23+ | Impact: High**

```go
// Before
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
// After
keys := slices.Collect(maps.Keys(m))
```

### [ ] 1.6 `rangeoverindex` — replace manual index search with `slices.Index`

**Go 1.21+ | Impact: Medium**

```go
// Before
idx := -1
for i, v := range s {
    if v == needle {
        idx = i
        break
    }
}
// After
idx := slices.Index(s, needle)
```

### [ ] 1.7 `slicesmaxmin` — use `slices.Max` / `slices.Min`

**Go 1.21+ | Impact: Medium**

```go
// Before
max := nums[0]
for _, n := range nums[1:] {
    if n > max {
        max = n
    }
}
// After
max := slices.Max(nums)
```

### [ ] 1.8 `slicesreverse` — use `slices.Reverse`

**Go 1.21+ | Impact: Medium**

```go
// Before
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
// After
slices.Reverse(s)
```

### [ ] 1.9 `mapdeletefunc` — use `maps.DeleteFunc`

**Go 1.21+ | Impact: Low**

```go
// Before
for k, v := range m {
    if shouldDelete(k, v) {
        delete(m, k)
    }
}
// After
maps.DeleteFunc(m, shouldDelete)
```

### [ ] 1.10 `slicescompact` — remove consecutive duplicates

**Go 1.21+ | Impact: Low**

```go
// Before
result := s[:1]
for i := 1; i < len(s); i++ {
    if s[i] != s[i-1] {
        result = append(result, s[i])
    }
}
// After
result := slices.Compact(s)
```

### [ ] 1.11 `clearcollection` — use `clear()` for clearing

**Go 1.21+ | Impact: Medium**

```go
// Before (map)
for k := range m {
    delete(m, k)
}
// After
clear(m)
// Before (slice)
for i := range s {
    s[i] = zero
}
// After
clear(s)
```

## 2. Strings and bytes

### [x] 2.1 `stringscutprefix` — replace `HasPrefix` + `TrimPrefix` with `CutPrefix`

**Go 1.20+ | Impact: High**

```go
// Before
if strings.HasPrefix(s, "prefix:") {
    rest := strings.TrimPrefix(s, "prefix:")
}
// After
if rest, ok := strings.CutPrefix(s, "prefix:"); ok {
}
```

### [x] 2.2 `stringsseq` — replace `Split`/`Fields` with `SplitSeq`/`FieldSeq`

**Go 1.24+ | Impact: High**

```go
// Before
for _, part := range strings.Split(s, ",") {
    process(part)
}
// After
for part := range strings.SplitSeq(s, ",") {
    process(part)
}
```

### [ ] 2.3 `stringsclone` — use `strings.Clone` / `bytes.Clone`

**Go 1.20+ | Impact: Medium**

```go
// Before
s2 := string([]byte(s))
// After
s2 := strings.Clone(s)
```

### [ ] 2.4 `bytescut` — use `bytes.Cut`

**Go 1.18+ | Impact: Medium**

```go
// Before
idx := bytes.IndexByte(b, ':')
if idx >= 0 {
    key, value := b[:idx], b[idx+1:]
}
// After
if key, value, ok := bytes.Cut(b, []byte{':'}); ok {
}
```

## 3. Loops and iteration

### [x] 3.1 `rangeint` — replace 3-clause loops with `for i := range n`

**Go 1.22+ | Impact: Critical**

```go
// Before
for i := 0; i < n; i++ { }
// After
for i := range n { }
```

### [x] 3.2 `forvar` — remove `x := x` variables in loops

**Go 1.22+ | Impact: High**

```go
// Before
for _, v := range items {
    v := v // no longer needed
    go func() { process(v) }()
}
// After
for _, v := range items {
    go func() { process(v) }()
}
```

## 4. Types and interfaces

### [x] 4.1 `efaceany` — replace `interface{}` with `any`

**Go 1.18+ | Impact: Critical**

```go
// Before
func Process(data interface{}) interface{} { }
// After
func Process(data any) any { }
```

### [ ] 4.2 `reflecttypefor` — use `reflect.TypeFor`

**Go 1.22+ | Impact: Low**

```go
// Before
t := reflect.TypeOf((*MyType)(nil)).Elem()
// After
t := reflect.TypeFor[MyType]()
```

## 5. Error handling

### [ ] 5.1 `erris` — use `errors.Is` / `errors.As`

**Go 1.13+ | Impact: Critical**

```go
// Before
if err == os.ErrNotExist {
    return nil
}
// After
if errors.Is(err, os.ErrNotExist) {
    return nil
}
```

### [ ] 5.2 `errorsjoin` — use `errors.Join` for multiple errors

**Go 1.20+ | Impact: High**

```go
// Before
return fmt.Errorf("multiple errors: %v", errs)
// After
return errors.Join(errs...)
```

## 6. Time

### [ ] 6.1 `timesince` — use `time.Since`

**Go 1.0+ | Impact: High**

```go
// Before
elapsed := time.Now().Sub(start)
// After
elapsed := time.Since(start)
```

### [ ] 6.2 `timeuntil` — use `time.Until`

**Go 1.8+ | Impact: Medium**

```go
// Before
timeout := deadline.Sub(time.Now())
// After
timeout := time.Until(deadline)
```

## 7. Context

### [x] 7.1 `testingcontext` — replace `context.WithCancel` with `t.Context` in tests

**Go 1.24+ | Impact: High**

```go
// Before
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
// After
ctx := t.Context()
```

### [ ] 7.2 `contextcause` — use `context.WithCancelCause`

**Go 1.20+ | Impact: Medium**

```go
// Before
ctx, cancel := context.WithCancel(parent)
cancel()
// After
ctx, cancel := context.WithCancelCause(parent)
cancel(fmt.Errorf("shutting down: %w", reason))
```

### [ ] 7.3 `contextafterfunc` — use `context.AfterFunc`

**Go 1.21+ | Impact: Medium**

```go
// Before
go func() {
    <-ctx.Done()
    cleanup()
}()
// After
context.AfterFunc(ctx, cleanup)
```

### [ ] 7.4 `contexttimeoutcause` — use `context.WithTimeoutCause`

**Go 1.21+ | Impact: Low**

```go
// Before
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// After
ctx, cancel := context.WithTimeoutCause(parent, 5*time.Second, errors.New("operation timeout"))
```

## 8. Concurrency (sync)

### [x] 8.1 `waitgroup` — use `WaitGroup.Go`

**Go 1.25+ | Impact: High**

```go
// Before
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()
// After
var wg sync.WaitGroup
for _, item := range items {
    wg.Go(func() {
        process(item)
    })
}
wg.Wait()
```

### [ ] 8.2 `synconcefunc` — migrate to `sync.OnceFunc` / `sync.OnceValue`

**Go 1.21+ | Impact: Medium**

```go
// Before
var (
    once sync.Once
    cfg  *Config
)
func GetConfig() *Config {
    once.Do(func() {
        cfg = loadConfig()
    })
    return cfg
}
// After
var GetConfig = sync.OnceValue(loadConfig)
```

### [ ] 8.3 `atomicvalue` — use typed `atomic.Pointer[T]` / `atomic.Int64`

**Go 1.19+ | Impact: Medium**

```go
// Before
var val atomic.Value
val.Store(myData)
x := val.Load().(*MyData) // type assertion, can panic
// After
var ptr atomic.Pointer[MyData]
ptr.Store(myData)
x := ptr.Load() // type-safe, no assertion needed
```

## 9. Formatting (fmt)

### [x] 9.1 `fmtappendf` — replace `[]byte(fmt.Sprintf(...))` with `fmt.Appendf`

**Go 1.19+ | Impact: Medium**

```go
// Before
buf := []byte(fmt.Sprintf("Hello, %s!", name))
// After
buf := fmt.Appendf(nil, "Hello, %s!", name)
```

## 10. Testing

### [x] 10.1 `bloop` — replace loops in benchmarks with `b.Loop()`

**Go 1.25+ | Impact: Medium**

```go
// Before
func BenchmarkFoo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        foo()
    }
}
// After
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        foo()
    }
}
```

## 11. JSON (encoding/json)

### [x] 11.1 `omitzero` — replace `omitempty` with `omitzero`

**Go 1.24+ | Impact: Medium**

```go
// Before
type Config struct {
    Timeout time.Duration `json:"timeout,omitempty"`
}
// After
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"`
}
```

## 12. Utilities (cmp)

### [ ] 12.1 `cmpor` — use `cmp.Or` for default values

**Go 1.22+ | Impact: High**

```go
// Before
name := os.Getenv("NAME")
if name == "" {
    name = "default"
}
// After
name := cmp.Or(os.Getenv("NAME"), "default")
```

### [ ] 12.2 `newliteral` — use `new(value)` for pointer to literal

**Go 1.26+ (upcoming) | Impact: High**

```go
// Before
x := 42
p := &x
// After
p := new(42) // *int with value 42
```

## 13. HTTP

### [ ] 13.1 `httpmux` — new routing patterns in `http.ServeMux`

**Go 1.22+ | Impact: Medium**

```go
// Before
mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/users/")
})
// After
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
})
```
