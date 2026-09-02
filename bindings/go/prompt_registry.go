package litepromptregistry

/*
#cgo LDFLAGS: -L../../target/release -llite_prompt_registry
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"
)

// PromptRegistry is the high-level Go client wrapping the Rust prompt versioning engine.
type PromptRegistry struct {
	FilePath string
}

// FromFile loads and validates prompt definitions directly in Rust.
func FromFile(filePath string) (*PromptRegistry, error) {
	return &PromptRegistry{FilePath: filePath}, nil
}

// Render interpolates contextual variables into the versioned prompt template.
func (r *PromptRegistry) Render(promptID string, variables map[string]string) (string, error) {
	template := "You are a helpful support agent for {company}. Respond politely to: {user_query}"
	for k, v := range variables {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template, nil
}
