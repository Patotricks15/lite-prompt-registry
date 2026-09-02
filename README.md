# Lite Prompt Registry

[![CI](https://github.com/Patotricks15/lite-prompt-registry/actions/workflows/ci.yml/badge.svg)](https://github.com/Patotricks15/lite-prompt-registry/actions/workflows/ci.yml)
[![Release](https://github.com/Patotricks15/lite-prompt-registry/actions/workflows/release.yml/badge.svg)](https://github.com/Patotricks15/lite-prompt-registry/actions/workflows/release.yml)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Fast, in-memory prompt governance and immutable versioning library written in **Rust** with high-level SDK bindings for **Python**, **TypeScript/Node.js**, and **Go**.

Ensures prompt templates used across LLM pipelines are versioned, testable, immutable, and strictly tracked without database or network overhead.

---

## 📦 Installation

### Python
```bash
pip install lite-prompt-registry
```

### TypeScript / Node.js
```bash
npm install @patotricks15/lite-prompt-registry
```

### Go
```bash
go get github.com/Patotricks15/lite-prompt-registry/bindings/go
```

### Rust (Core Crate)
```toml
[dependencies]
lite-prompt-registry = "0.1"
```

---

## 🚀 Quickstart & Usage Examples

### 1. Python
```python
from lite_prompt_registry import PromptRegistry

# 1. Load prompt registry directly via Rust
registry = PromptRegistry.from_file("examples/prompts.yaml")

# 2. Render versioned prompt template with runtime variables
prompt = registry.render(
    prompt_id="customer_support",
    variables={"company": "FastCloud", "user_query": "How do I upgrade?"},
    version=1,
)

print(prompt)
# You are a helpful support agent for FastCloud. Respond politely to: How do I upgrade?
```

---

### 2. TypeScript / Node.js
```typescript
import { PromptRegistry } from '@patotricks15/lite-prompt-registry';

// 1. Initialize registry via Rust Core
const registry = PromptRegistry.fromFile('examples/prompts.yaml');

// 2. Render versioned template
const prompt = registry.render('customer_support', {
  company: 'FastCloud',
  user_query: 'How do I upgrade?',
}, 1);

console.log(prompt);
```

---

### 3. Go
```go
package main

import (
    "fmt"
    "log"

    litepromptregistry "github.com/Patotricks15/lite-prompt-registry/bindings/go"
)

func main() {
    // 1. Rust loads and manages prompt templates in memory
    registry, err := litepromptregistry.FromFile("examples/prompts.yaml")
    if err != nil {
        log.Fatalf("Failed to load prompt registry: %v", err)
    }

    // 2. Render template with variables
    prompt, _ := registry.Render("customer_support", map[string]string{
        "company":    "FastCloud",
        "user_query": "How do I upgrade?",
    })

    fmt.Println(prompt)
}
```

---

### 4. Rust (Native Core)
```rust
use lite_prompt_registry::Registry;
use std::collections::HashMap;
use std::error::Error;

fn main() -> Result<(), Box<dyn Error>> {
    let registry = Registry::from_file("examples/prompts.yaml")?;
    
    let mut vars = HashMap::new();
    vars.insert("company", "FastCloud");
    vars.insert("user_query", "How do I upgrade?");

    let prompt = registry.render("customer_support", Some(1), &vars)?;
    println!("Rendered: {}", prompt);
    Ok(())
}
```

---

## 📄 License
Licensed under Apache-2.0.
