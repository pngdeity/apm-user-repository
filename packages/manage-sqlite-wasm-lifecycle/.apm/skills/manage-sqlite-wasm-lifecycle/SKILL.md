---
license: MIT
name: manage-sqlite-wasm-lifecycle
description: >
  Safely manage SQLite database file lifecycles and DbContext connection states 
  in Blazor WebAssembly sandboxes during database resets, imports, exports, and restores.
  USE FOR: managing virtual sqlite.db files inside browser filesystems (OPFS/IndexedDB),
  draining DbContext connection pools before replacing database files on disk, executing direct
  binary database exports, and configuring connection strings for local browser sandboxes.
  DO NOT USE FOR: standard server-side database migrations or normal EF Core entity schema mapping.
---

# Manage SQLite WebAssembly Lifecycle

In client-side Blazor WebAssembly applications using SQLite (via WebAssembly virtual disk drivers), the SQLite engine binds directly to the virtual filesystem in the browser (e.g., OPFS or memory-backed storage). Overwriting or replacing the database file while active connections are open leads to unrecoverable engine corruption (`SQLITE_CORRUPT` or `SQLITE_BUSY` errors).

## Purpose
This skill provides a safe playbook for importing, exporting, and managing browser-embedded SQLite database files, preventing file-handle lockups and data corruption.

## Inputs
Before executing this workflow, the agent needs:
* An active client-side Blazor project containing an EF Core `DbContext`.
* Access to `IJSRuntime` (for browser-level interop) and the local assembly `System.IO` API.
* An incoming database backup file represented as a `byte[]` array or base64 string.

## Workflow

### Step 1 — Drain the DbContext Connection Pool
Before modifying the SQLite file on the virtual disk, release all active file locks:
1. Ensure all active scoped `DbContext` instances have been fully disposed.
2. Force the SQLite connection pool to clear and release all virtual handles:
   ```csharp
   Microsoft.Data.Sqlite.SqliteConnection.ClearAllPools();
   ```
*Checkpoint: ClearAllPools() completes synchronously, ensuring no active lock remains on `contest.db`.*

### Step 2 — Validate the Database Payload
Never write incoming bytes directly to the virtual disk without confirming file integrity:
1. Verify the byte array length is at least 100 bytes (the minimum size of a valid SQLite header).
2. Check the first 16 bytes for the standard SQLite magic header:
   ```csharp
   var magic = "SQLite format 3\0"u8;
   ```
3. Throw an `ArgumentException` early if validation fails.
*Checkpoint: Mismatched magic headers are caught before writing, preventing file corruption.*

### Step 3 — Write the Database File
Perform the filesystem replacement safely:
1. Write the validated bytes directly to the target database file using asynchronous filesystem operations:
   ```csharp
   await File.WriteAllBytesAsync("contest.db", backupBytes);
   ```
*Checkpoint: The target database file is completely replaced with valid SQLite structures.*

## Validation
To confirm the database was successfully restored:
1. Re-initialize a new scope and resolve your `DbContext` instance.
2. Run a lightweight query to check if the schema is readable:
   ```csharp
   var isAvailable = await context.Categories.AnyAsync();
   ```
3. Verify no `SqliteException` or crash is thrown.

## Common Pitfalls

| Pitfall | Why It Fails | Correct Approach |
| :--- | :--- | :--- |
| Overwriting the file with active `DbContext` instances open | SQLite locks the file and throws `SQLITE_BUSY`. | Close active contexts and call `ClearAllPools()` before writing. |
| Trusting unverified inputs from JS `localStorage` | Corrupt or empty localStorage keys will write blank or malformed files, crashing the app on reload. | Always validate incoming sizes and magic headers prior to a write operation. |
