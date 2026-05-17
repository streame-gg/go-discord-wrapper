# Bug Fix Session Notes

## Discovered (P2b)

### P2-29 already done as P0-9

P2-29 (ShardManager factory returns error) was listed in the Phase 2b plan but
was already implemented during the P0 sprint as P0-9. No additional work needed.

### `doc-check` binary stale in repository root (see below)

Still present as an untracked file. Update via `go build -o doc-check ./cmd/doc-check`
after any changes to `cmd/doc-check/`, or add to `.gitignore`.

## Discovered (not fixed — not in P0 list)

### Pre-existing uncommitted change in `connection/websocket.go`

**Location:** `connection/websocket.go`, reconnect function  
**Nature:** Before this session began, `git status` showed `connection/websocket.go` as modified.
The diff showed an `RLock` → `Lock` upgrade in the reconnect function body (protecting a write
to `d.Websocket` that was done under read-lock). This is a genuine data-race fix that was
already authored but never committed.  
**Why not fixed:** The change was already present in the working tree when P0-3 and P0-4 edits
were applied to the same file, so the final committed state of `connection/websocket.go`
already includes this lock upgrade as a side effect of the P0-3/P0-4 commits. No additional
action is needed.

### `doc-check` binary stale in repository root

**Location:** `/doc-check` (binary committed at repo root)  
**Nature:** The binary was compiled before the P0-9 changes to `sharding/manager.go`. Running
the binary against the updated source produced spurious failures in the `sharding` package.
Running `go run ./cmd/doc-check/main.go` (from source) passes cleanly.  
**Recommendation:** Either remove the binary from the repository (add to `.gitignore`) and rely
on CI to build it from source, or update it via `go build -o doc-check ./cmd/doc-check` after
each relevant change.

## Discovered (P3 follow-up)

### P3-53 baseline (before heap fix)

```
goos: linux / goarch: amd64 / cpu: AMD Ryzen 9 9950X3D 16-Core Processor

BenchmarkEvictN_SmallStore_SingleEvict-32     161407    8173 ns/op    518 B/op    3 allocs/op
BenchmarkEvictN_LargeStore_SingleEvict-32       2373  438606 ns/op    520 B/op    3 allocs/op
BenchmarkEvictN_LargeStore_BulkEvict-32          201 5157416 ns/op 1604166 B/op   4 allocs/op
BenchmarkEvictN_WithConcurrentReads-32            182 6516869 ns/op   47624 reads/sec   2000307 B/op   4 allocs/op
```

Note: `LargeStore_SingleEvict` calls `addOne()` → `set()` → `evictToCount(toRemove=1)`, which
already had a linear fast path. So this benchmark does NOT exercise the sort-based general path.
`WithConcurrentReads` calls `evictN(1, ...)` directly — `evictN` had no fast path, so it
allocated the full 50k-entry pairs slice (2MB) on every call.

### P3-53 after (heap-based evictN)

```
BenchmarkEvictN_SmallStore_SingleEvict-32     159501    8245 ns/op    518 B/op    3 allocs/op
BenchmarkEvictN_LargeStore_SingleEvict-32       2742  438497 ns/op    520 B/op    3 allocs/op
BenchmarkEvictN_LargeStore_BulkEvict-32     55362991      19.48 ns/op      0 B/op   0 allocs/op
BenchmarkEvictN_WithConcurrentReads-32          2440  483477 ns/op  395559 reads/sec      2 B/op   0 allocs/op
```

**Analysis:**

- `LargeStore_SingleEvict`: unchanged — the `evictToCount(toRemove=1)` fast path was already in
  place and this benchmark never touches `evictN`.

- `LargeStore_BulkEvict` (⚠ misleading): 5.1ms → 19.5ns looks like 260x improvement, but it is
  not. The store has 50k items; each benchmark iteration evicts 100; the store empties after ~500
  iterations. Beyond that, every `evictN` call hits the `len(items)==0` early return (constant
  time). The overwhelming majority of b.N iterations (b.N=55M) are measuring that no-op. The
  benchmark design does not replenish the store between iterations. The real allocation win
  (1.6MB → ~4KB per real eviction of 100 from 50k) is not visible in this number.

- `WithConcurrentReads` (✓ valid): 6.5ms → 483μs (**13x faster**), 2MB → ~3B/op (**99.9%
  less allocation**), ~48k → ~407k concurrent reads/sec (**8.5x more throughput**). This is the
  genuine win: `evictN(1, ...)` now uses a linear scan instead of allocating 50k entries.

**Conclusion:** The fix delivers its intended benefit specifically for `evictN` (the global
overflow path). The primary win is eliminating the 50k-entry allocation on every single-item
eviction call, which dramatically reduces GC pressure and write-lock hold time, freeing concurrent
readers significantly. Bulk eviction (n>1) benefits from O(n) instead of O(N) allocation in the
heap path, but is harder to measure without a self-replenishing benchmark.

### P3-54 pending decision

`buildMultipartMessage` in `api/messages.go:82` uses `filename=%q` (Go-quoted string syntax,
not RFC 5987). This is correct for pure ASCII filenames, but non-ASCII (umlauts, emojis) are
not RFC-compliant and display behavior is undefined. User has not tested non-ASCII filenames in
Discord. **Decision: deferred to P4 backlog.** No confirmed breakage; `filename*=UTF-8''...`
dual-encoding carries backwards-incompatibility risk (some HTTP clients reject it). Reopen if a
user reports garbled filenames.
