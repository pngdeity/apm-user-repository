# Semantic Chunking and Layout Normalization

## Data Normalization Principles

1. **Layout Normalization:**
   - Visual layouts (e.g., PDF text columns, visual tables) MUST be converted into semantic Markdown (e.g., Markdown tables, hierarchical headers).
   - Strip site-wide navigation, footers, and non-semantic fluff (syntactic sugar).

2. **Semantic Chunking:**
   - Ensure facts are not orphaned during chunking.
   - Use explicit Markdown headers to "tether" numerical values or specific properties to their parent objects.
   - Preserve image `alt` text and context around figures.

3. **Compression (Optional):**
   - For high-token payloads, use algorithmic compression (e.g., LLMLingua) to remove stop-words and redundant syntax without losing factual signal.
