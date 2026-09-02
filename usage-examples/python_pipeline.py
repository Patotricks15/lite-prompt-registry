"""Lite Prompt Registry - High-Level Python Pipeline Example.

Loads immutable versioned prompt definitions directly via Rust and renders them for LLMs.
"""

import sys
import os

# Importa o binding local
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../bindings/python")))
from lite_prompt_registry import PromptRegistry


def main():
    print("=== Lite Prompt Registry - High-Level Python Example ===")

    # 1. Uma única linha: Rust carrega e valida o arquivo de prompts YAML
    prompts_path = "examples/prompts.yaml"
    print(f"\n[1] Carregando prompts de '{prompts_path}' via Rust Core...")
    registry = PromptRegistry.from_file(prompts_path)

    # 2. Renderizando template com contexto do usuário
    context_vars = {
        "company": "FastCloud",
        "user_query": "How do I upgrade to the Enterprise plan?",
    }
    rendered_prompt = registry.render("customer_support", context_vars, version=1)

    print("\n[2] Prompt renderizado a partir da versão v1:")
    print(f" -> {rendered_prompt!r}")

    # 3. Enviando para o LLM
    print("\n[3] Invocando LLM com prompt consistente e governado em memória.")


if __name__ == "__main__":
    main()
