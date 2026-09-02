package main

import (
	"fmt"
	"strings"
)

// PromptVersion represents an immutable prompt template version.
type PromptVersion struct {
	VersionID string
	Template  string
}

// PromptRegistry manages in-memory prompt versions.
type PromptRegistry struct {
	prompts map[string]PromptVersion
}

func NewPromptRegistry() *PromptRegistry {
	return &PromptRegistry{
		prompts: make(map[string]PromptVersion),
	}
}

func (r *PromptRegistry) Register(name string, versionID string, template string) {
	r.prompts[name+":"+versionID] = PromptVersion{
		VersionID: versionID,
		Template:  template,
	}
}

func (r *PromptRegistry) Render(name, versionID string, variables map[string]string) (string, error) {
	key := name + ":" + versionID
	pv, exists := r.prompts[key]
	if !exists {
		return "", fmt.Errorf("prompt not found: %s", key)
	}

	result := pv.Template
	for k, v := range variables {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}
	return result, nil
}

func main() {
	fmt.Println("=== Lite Prompt Registry - Go Example ===")

	registry := NewPromptRegistry()

	// 1. Register versioned prompt
	registry.Register(
		"customer_support",
		"v1.0.0",
		"You are a helpful support agent for {company}. Respond politely to customer issue: {issue}",
	)

	// 2. Render prompt with context variables
	rendered, err := registry.Render("customer_support", "v1.0.0", map[string]string{
		"company": "CloudCorp",
		"issue":   "Cannot reset account password",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("\n[1] Rendered Prompt from Registry (v1.0.0):")
	fmt.Printf(" -> %q\n", rendered)

	fmt.Println("\n[2] Sending rendered prompt to LLM...")
	fmt.Println(" -> Response received and tracked against prompt version 'customer_support:v1.0.0'.")
}
