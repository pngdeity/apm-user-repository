---
name: normalize-dataset
description: Procedural workflow for pre-processing raw documents (PDFs, transcripts, legacy text) into semantic, LLM-friendly Markdown datasets. Use when preparing data for RAG, ingestion, or knowledge graph construction. Not for live data processing, API design, or database schema design.
allowed-tools: grep
metadata:
  tags: "data normalization rag ingestion pre-processing pdf markdown"
compatibility: Generic — no environment restrictions. May benefit from PDF/HTML parsing tools in the environment.
---

# Normalize Dataset Skill

This skill ensures that raw data is transformed into a high-fidelity semantic structure before being consumed by an LLM. This prevents data loss during chunking and optimizes the context window.

## Workflow: Layout Normalization

When handling raw files (e.g., PDF, HTML):
1. **Identify Structures:** Locate tables, multi-column layouts, and hierarchical headers.
2. **Convert to Markdown:** Use a "Layout Normalization" pass to convert these into standard Markdown tables and headers.
3. **Strip Noise:** Remove non-essential content (headers/footers, navigation).
- **If a file cannot be parsed:** Skip it and log the file path. Do not halt the entire batch for a single unparseable file.

## Workflow: Semantic Chunking

1. **Tethering:** Ensure every chunk contains its parent context. For example, a row in a table should include the table's header context.
2. **Boundary Preservation:** Never split a chunk in the middle of a sentence or a related data group (e.g., a spec list).
3. **Reference Principle:** Follow the principles in [references/principles.md](references/principles.md).

## Workflow: Verification

1. **Schema Check:** Ensure the resulting Markdown validates against your project's expected schema.
2. **Completeness Test:** Verify that critical numerical values or IDs found in the raw source are present in the normalized dataset.
- **If schema check fails:** Report which field or structure is non-conforming. Do not proceed to chunking until the schema issue is resolved.

## Verification
Verify this skill produces correct output:
1. Take a fixture HTML document containing a table and multi-column layout.
2. Run Layout Normalization and confirm the output contains standard Markdown tables.
3. Run Semantic Chunking and confirm no chunk splits mid-sentence.
4. Confirm the Completeness Test catches a deliberately removed ID from the source.
