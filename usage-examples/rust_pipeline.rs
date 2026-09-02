//! Lite Prompt Registry - Rust Pipeline Example.
//!
//! Demonstrates loading versioned prompts from YAML (examples/prompts.yaml) and rendering.

use lite_prompt_registry::Registry;
use std::collections::HashMap;
use std::error::Error;

fn main() -> Result<(), Box<dyn Error>> {
    println!("=== Lite Prompt Registry - Rust Example ===");

    // 1. Rust carrega e valida o arquivo de prompts YAML
    let prompts_path = "examples/prompts.yaml";
    println!("\n[1] Carregando prompts de '{}'...", prompts_path);
    let registry = Registry::from_file(prompts_path)?;

    // 2. Renderização de template imutável com variáveis de contexto
    let mut vars = HashMap::new();
    vars.insert("company", "FastCloud");
    vars.insert("user_query", "How do I upgrade to the Enterprise plan?");

    let rendered = registry.render("customer_support", Some(1), &vars)?;
    println!("\n[2] Prompt renderizado a partir da versão v1:");
    println!(" -> \"{}\"", rendered);

    // 3. Enviando prompt formatado para o LLM
    println!("\n[3] Enviando prompt consistente para o LLM e rastreando versão 'customer_support:v1'.");
    Ok(())
}
