---
license: MIT
name: implement-offline-first-sync
description: >
  Design offline-first local data storage sync logs and implement optimistic 
  concurrency models to handle conflicts when syncing local SQLite instances back to a central server.
  USE FOR: managing database synchronization states, authoring local sync queues, 
  implementing version-concurrency keys (RowVersion/Timestamp), and designing conflict resolution strategies (Client Wins / Server Wins / Merge).
  DO NOT USE FOR: simple local storage caching or direct non-transactional database updates.
---

# Implement Offline-First Sync

In local-first Blazor WebAssembly architectures, database mutations are committed to the browser's local sandbox SQLite database. When network connectivity is restored, these local operations must be synchronized back to the primary server without silently overwriting concurrent changes made by other users.

## Purpose
This skill provides a programmatic playbook for building an offline sync queue, managing optimistic concurrency tokens, and resolving data synchronization conflicts in client-side .NET applications.

## Inputs
Before executing this workflow, the agent needs:
* An active SQLite-backed client-side EF Core application.
* Entities decorated with a concurrency version property (e.g., `uint ConcurrencyToken` or `long LastModifiedTicks`).
* A central synchronization API endpoint accepting batch payloads.

## Workflow

### Step 1 — Scaffold the Local Sync Queue
To record database mutations while offline, maintain a dedicated sync queue table inside the SQLite schema:
1. Define a `SyncItem` entity to record changes:
   ```csharp
   public class SyncItem
   {
       public int Id { get; set; }
       public string EntityName { get; set; } = string.Empty;
       public string EntityId { get; set; } = string.Empty;
       public string Action { get; set; } = string.Empty; // INSERT, UPDATE, DELETE
       public string Payload { get; set; } = string.Empty; // JSON Serialized Entity DTO
       public long CreatedAt { get; set; }
   }
   ```
2. Set up the DB schema and verify it is generated during context initialization.
*Checkpoint: Ensure `SyncItem` mutations are transactionally saved alongside primary model updates.*

### Step 2 — Implement Optimistic Concurrency Controls
Prevent overwrite races by tracking entity state versions:
1. Ensure the server API rejects any mutation payload where the incoming concurrency token does not match the active server token.
2. In the local database client, keep the current token in-memory and increment the token only when a server synchronization operation succeeds.
*Checkpoint: Client-side writes will trigger an optimistic lock warning if the server's record has evolved.*

### Step 3 — Process the Sync Queue & Resolve Conflicts
When synchronization starts:
1. Query local `SyncQueue` items ordered by `CreatedAt`.
2. Push items sequentially to the backend server.
3. Catch any concurrency conflicts (e.g., HTTP `409 Conflict` or `DbUpdateConcurrencyException`) and apply the conflict resolution policy:
   * **Server Wins:** Discard the local change, query the server's state, and update the local SQLite database.
   * **Client Wins:** Re-fetch the server's record, override its fields, set the client's concurrency token, and force a sync update.
   * **Three-Way Merge:** Compare conflicting fields, merge non-overlapping properties, and prompt the user for manual overrides if categories clash.
4. Remove successfully synchronized items from the local `SyncQueue`.
*Checkpoint: The local sync queue is fully drained and matches the server's primary database state.*

## Validation
To verify the synchronization engine is functioning correctly:
1. Put the browser in offline mode (using DevTools Network throttling) and perform several scoring mutations.
2. Verify that local SQLite updates succeed and the `SyncQueue` table collects matching transaction records.
3. Re-enable network connectivity, run the sync process, and verify that:
   - The central server receives the mutations.
   - The local `SyncQueue` is completely emptied.
   - All synchronized entities have their concurrency tokens successfully updated.

## Common Pitfalls

| Pitfall | Why It Fails | Correct Approach |
| :--- | :--- | :--- |
| Deleting `SyncQueue` items before synchronization succeeds | Network drops will cause permanent data loss. | Delete local sync entries only *after* receiving a successful synchronization response from the server. |
| Synchronizing out of order | Child entities (like `EntryScore`) will fail to insert if parent entities (like `Entry`) have not synced yet. | Always process the queue strictly in the chronological order (`CreatedAt` ascending) in which the changes occurred. |
