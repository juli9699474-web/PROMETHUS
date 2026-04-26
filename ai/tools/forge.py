from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass
class ToolSpec:
    description: str
    inputs: dict[str, Any]
    outputs: dict[str, Any]


@dataclass
class ValidationResult:
    passed: bool
    tool: dict[str, Any] | None = None
    errors: list[str] | None = None


class ToolForge:
    def __init__(self, llm: Any, validator: Any, registry: Any) -> None:
        self.llm = llm
        self.validator = validator
        self.registry = registry

    async def generate_tool(self, spec: ToolSpec) -> dict[str, Any]:
        code = await self.llm.generate(
            system="You are a secure MCP tool generator.",
            user=(
                f"Create a tool that: {spec.description}\n"
                f"Inputs: {spec.inputs}\nOutputs: {spec.outputs}"
            ),
        )
        validated: ValidationResult = await self.validator.run_in_sandbox(code)
        if validated.passed and validated.tool is not None:
            return await self.registry.register(validated.tool)
        return await self.repair_and_retry(code, validated.errors or [])

    async def repair_and_retry(self, code: str, errors: list[str]) -> dict[str, Any]:
        patched = await self.llm.generate(
            system="Repair insecure or broken Python tool code.",
            user=f"Fix this code:\n{code}\n\nErrors:\n{errors}",
        )
        validated: ValidationResult = await self.validator.run_in_sandbox(patched)
        if not validated.passed or validated.tool is None:
            raise RuntimeError("tool forge failed after repair attempt")
        return await self.registry.register(validated.tool)
