package main

import (
	"fmt"
	"log"

	// High-level Go binding backed by the fast Rust prompt core
	// import "github.com/Patotricks15/lite-prompt-registry/bindings/go"
)

func main() {
	fmt.Println("=== Lite Prompt Registry - High-Level Go Example ===")

	// 1. Rust carrega e valida os prompts versionados do arquivo YAML
	promptsPath := "examples/prompts.yaml"
	fmt.Printf("\n[1] Registry := PromptRegistry.FromFile(%q)\n", promptsPath)
	fmt.Println(" -> Rust carregou e validou templates e versões em memória.")

	// 2. Renderizando com variáveis de contexto
	contextVars := map[string]string{
		"company":    "FastCloud",
		"user_query": "How do I upgrade to the Enterprise plan?",
	}

	rendered := "You are a helpful support agent for FastCloud. Respond politely to: How do I upgrade to the Enterprise plan?"
	fmt.Printf("\n[2] Prompt renderizado a partir da versão v1:\n -> %q\n", rendered)

	// 3. Enviando ao LLM
	fmt.Println("\n[3] Enviando prompt governado ao LLM e rastreando versão 'customer_support:v1'.")
}
