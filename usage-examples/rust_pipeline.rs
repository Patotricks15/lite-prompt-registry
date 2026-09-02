//! Lite Prompt Registry - Rust Pipeline Example.
//!
//! Direct Rust core usage to create, version, and fetch prompts.

use lite_prompt_registry::Registry;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("=== Lite Prompt Registry - Rust Example ===");
    let mut registry = Registry::default();

    // 1. Create prompt entity
    let prompt_id = registry.create_prompt("support_agent", "Primary customer support prompt");
    println!("Created prompt with ID: {}", prompt_id);

    // 2. Add immutable versions
    let v1 = registry.create_version(prompt_id, "You are a support bot. Help with: {q}", "alice")?;
    let v2 = registry.create_version(prompt_id, "You are an empathetic support specialist. Resolve: {q}", "bob")?;

    println!("Registered versions: v{} and v{}", v1, v2);

    // 3. Retrieve prompt for LLM invocation
    let prompt = registry.get(prompt_id).expect("Prompt exists");
    let active_version = &prompt.versions[1]; // v2

    println!("\nPreparing LLM call:");
    println!(" -> Active template: '{}'", active_version.template);
    println!(" -> Author: {}", active_version.author);
    println!(" -> Ready to format and send to LLM.");

    Ok(())
}
