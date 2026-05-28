---
license: MIT
name: configure-trim-safe-efcore
description: >
  Configure Entity Framework Core for reflection-free compilation, Native AOT 
  compatibility, and aggressive assembly trimming in Blazor WebAssembly apps.
  USE FOR: resolving trim warnings (IL2026/IL2111), configuring precompiled 
  queries with dotnet ef dbcontext optimize, implementing source-generated JSON 
  contexts, and minimizing WebAssembly build payloads.
  DO NOT USE FOR: standard runtime reflection-based queries or server-side DB seeding.
---

# Configure Trim-Safe EF Core

Entity Framework Core relies heavily on runtime reflection and dynamic code compilation to discover model schemas and execute LINQ queries. In trimmed Blazor WebAssembly applications, reflection causes massive WASM payload sizes and introduces runtime crash risks when the IL Linker strips out internal C# entity properties.

## Purpose
This skill provides a systematic workflow to optimize and configure EF Core for reflection-free builds, ensuring safety and minimal payloads in trimmed Blazor WebAssembly projects.

## Inputs
Before executing this workflow, the agent needs:
* An active Blazor WebAssembly project targeting .NET 9.0/10.0+ using EF Core.
* The `.csproj` configured for publishing with trimming (`<PublishTrimmed>true</PublishTrimmed>`).
* Installation of the `.NET WebAssembly build tools` workload on the system.

## Workflow

### Step 1 — Configure Query Precompilation
Instead of compiling LINQ queries dynamically at runtime, generate static interceptors:
1. Run the EF Core optimizer tool to generate compiled model files:
   ```bash
   dotnet ef dbcontext optimize --output-dir Persistence/CompiledModels --namespace YourApp.CompiledModels
   ```
2. Modify your DbContext options registration in `Program.cs` to use the precompiled model:
   ```csharp
   options.UseModel(YourApp.CompiledModels.ContestDbContextModel.Instance);
   ```
*Checkpoint: Ensure `Persistence/CompiledModels` is populated with CompiledModel files.*

### Step 2 — Implement Source-Generated JSON Serialization
To prevent reflection crashes during Entity serialization (common when saving SQLite backups or parsing JSON payloads):
1. Declare a partial class inheriting from `JsonSerializerContext` and mark all your database entities as serializable:
   ```csharp
   [JsonSerializable(typeof(CategoryEntity))]
   [JsonSerializable(typeof(EntryEntity))]
   [JsonSerializable(typeof(RelationEntity))]
   public partial class EntityJsonContext : JsonSerializerContext { }
   ```
2. Use this context when deserializing payloads to avoid reflection:
   ```csharp
   var rawBase64 = JsonSerializer.Deserialize(raw, typeof(string), EntityJsonContext.Default);
   ```
*Checkpoint: All JSON operations on entities are resolved statically at compile time.*

### Step 3 — Apply Assembly Linker Descriptors
If third-party assemblies throw linker warnings during publish, declare explicit preservation rules:
1. Create a `LinkerConfig.xml` file at the project root.
2. Instruct the trimmer to preserve database entities:
   ```xml
   <linker>
     <assembly fullname="ContestJudging.Core">
       <type fullname="ContestJudging.Core.Entities.*" preserve="all" />
     </assembly>
   </linker>
   ```
3. Include the XML file in your `.csproj`:
   ```xml
   <ItemGroup>
     <IllinkImportDescriptor Include="LinkerConfig.xml" />
   </ItemGroup>
   ```
*Checkpoint: Critical database model structures are protected from being stripped by the trimmer.*

## Validation
To verify the trimming configurations are correct:
1. Publish the application using the Release configuration:
   ```bash
   dotnet publish -c Release
   ```
2. Confirm the compiler outputs zero `IL2026` or `IL2111` warnings during publish evaluation.
3. Launch the published WASM client in a headless browser test and verify the initial database load succeeds.

## Common Pitfalls

| Pitfall | Why It Fails | Correct Approach |
| :--- | :--- | :--- |
| Suppressing warnings using `[UnconditionalSuppressMessage]` | It hides warnings but does not prevent the trimmer from stripping models, leading to `NullReferenceException` at runtime. | Implement compile-time model precompilation and source serialization contexts instead. |
| Using dynamic runtime queries | Trimming strips methods needed to evaluate dynamic query expressions. | Rely solely on statically analyzable LINQ queries. |
