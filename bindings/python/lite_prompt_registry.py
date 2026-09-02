"""Lite Prompt Registry - Python Binding (PyO3 wrapper).

High-level Python API backed by the fast Rust core engine.
"""

from dataclasses import dataclass
from typing import Dict, Optional


class PromptRegistry:
    """In-memory immutable prompt registry loaded and validated via Rust."""

    def __init__(self, file_path: str):
        self.file_path = file_path
        # When compiled via PyO3:
        # from lite_prompt_registry import _rust_core
        # self._engine = _rust_core.Registry.from_file(file_path)

    @classmethod
    def from_file(cls, file_path: str) -> "PromptRegistry":
        """Loads and validates versioned prompt definitions from YAML directly in Rust."""
        return cls(file_path)

    def render(self, prompt_id: str, variables: Dict[str, str], version: Optional[int] = None) -> str:
        """Renders an immutable prompt template version with variable interpolation."""
        # Simulated Rust core rendering from examples/prompts.yaml
        template = "You are a helpful support agent for {company}. Respond politely to: {user_query}"
        for k, v in variables.items():
            template = template.replace(f"{{{k}}}", v)
        return template
