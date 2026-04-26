from __future__ import annotations

from dataclasses import dataclass
from typing import Any, List


@dataclass
class KnowledgeChunk:
    id: str
    content: str
    confidence: float


class KnowledgeGraph:
    """
    Phase-1 knowledge graph interface. Backend adapters are plugged in next.
    """

    async def absorb(self, content: Any, source: str, confidence: float) -> None:
        _ = (content, source, confidence)

    async def query(self, question: str, k: int = 20) -> List[KnowledgeChunk]:
        _ = (question, k)
        return []

    async def consolidate(self) -> None:
        return None
