"""Lite Prompt Registry - Python LLM Pipeline Example.

Demonstrates loading immutable prompt templates into an LLM call.
"""

def load_prompt_version(prompt_id: int, version: int):
    # Simulated registry retrieval (matches Lite Prompt Registry schema)
    registry = {
        (1, 1): "You are an assistant. Answer concisely: {query}",
        (1, 2): "You are a senior analyst. Provide structured insights for: {query}",
    }
    return registry.get((prompt_id, version))


def run_prompt_governed_chain(prompt_id: int, version: int, query: str):
    template = load_prompt_version(prompt_id, version)
    if not template:
        raise ValueError(f"Prompt {prompt_id} version {version} not found or not approved.")
    
    formatted_prompt = template.format(query=query)
    print(f"\n[1] Resolved immutable prompt (ID: {prompt_id}, v{version})")
    print(f" -> Rendered: '{formatted_prompt}'")
    
    # Invoking LLM
    print(" -> [2] Forwarding to LLM provider with lineage metadata...")
    metadata = {"prompt_id": prompt_id, "prompt_version": version}
    return {"response": "Simulated output", "metadata": metadata}


def main():
    print("=== Lite Prompt Registry - Python Example ===")
    result = run_prompt_governed_chain(1, 2, "Summarize Q3 revenue growth.")
    print(f" -> LLM Result: {result}")


if __name__ == "__main__":
    main()
