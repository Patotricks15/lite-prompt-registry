/**
 * Lite Prompt Registry - High-Level TypeScript Pipeline Example.
 *
 * Loads versioned prompt definitions directly via Rust and renders them for LLMs.
 */

import { PromptRegistry } from '../bindings/node/index';

async function main() {
  console.log('=== Lite Prompt Registry - High-Level TypeScript Example ===');

  // 1. Rust carrega e valida o arquivo de prompts YAML
  const promptsPath = 'examples/prompts.yaml';
  console.log(`\n[1] Carregando prompts de '${promptsPath}' via Rust Core...`);
  const registry = PromptRegistry.fromFile(promptsPath);

  // 2. Renderizando template com contexto
  const context = {
    company: 'FastCloud',
    user_query: 'How do I upgrade to Enterprise?',
  };

  const rendered = registry.render('customer_support', context, 1);
  console.log('\n[2] Prompt renderizado (versão v1):');
  console.log(` -> "${rendered}"`);

  // 3. Enviando ao LLM
  console.log('\n[3] Enviando prompt consistente ao LLM e rastreando versão da template.');
}

main().catch(console.error);
