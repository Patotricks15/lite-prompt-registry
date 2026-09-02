/**
 * Lite Prompt Registry - TypeScript LLM Pipeline Example.
 *
 * Demonstrates pulling verified prompt templates before LLM execution.
 */

interface PromptMetadata {
  promptId: number;
  version: number;
}

class PromptRegistryClient {
  private templates: Map<string, string> = new Map([
    ["1:1", "Translate to Spanish: {{input}}"],
    ["1:2", "Translate following text into formal European Spanish: {{input}}"],
  ]);

  getTemplate(promptId: number, version: number): string {
    const key = `${promptId}:${version}`;
    const tpl = this.templates.get(key);
    if (!tpl) throw new Error(`Prompt ${key} not found.`);
    return tpl;
  }
}

async function callLlmWithRegistry(promptId: number, version: number, userInput: string) {
  const registry = new PromptRegistryClient();
  const template = registry.getTemplate(promptId, version);
  const prompt = template.replace("{{input}}", userInput);

  console.log(`\n[Registry] Loaded Prompt #${promptId} (v${version})`);
  console.log(` -> Formatted text: "${prompt}"`);
  console.log(" -> [LLM Call]: Generating translation with attached metadata...");
  return { text: "Simulated translation", promptMetadata: { promptId, version } };
}

async function main() {
  console.log("=== Lite Prompt Registry - TypeScript Example ===");
  await callLlmWithRegistry(1, 2, "Good morning, how are you?");
}

main().catch(console.error);
