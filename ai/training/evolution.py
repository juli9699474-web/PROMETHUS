from __future__ import annotations

from dataclasses import dataclass
from typing import List

import torch
import torch.nn as nn
import torch.optim as optim


@dataclass
class EvolutionSample:
    reward: float
    action_vector: List[float]


class PromptPolicy(nn.Module):
    def __init__(self, dim: int = 64) -> None:
        super().__init__()
        self.prompt_embedding = nn.Parameter(torch.randn(dim) * 0.02)

    def forward(self, action: torch.Tensor) -> torch.Tensor:
        return torch.sigmoid(torch.dot(self.prompt_embedding, action))


class PromptEvolutionTrainer:
    """
    Micro-RL style optimizer that treats prompt parameters as differentiable state.
    """

    def __init__(self, dim: int = 64, lr: float = 1e-2) -> None:
        self.model = PromptPolicy(dim=dim)
        self.optimizer = optim.Adam(self.model.parameters(), lr=lr)

    def train_step(self, sample: EvolutionSample) -> float:
        action = torch.tensor(sample.action_vector, dtype=torch.float32)
        pred = self.model(action)
        target = torch.tensor(sample.reward, dtype=torch.float32).clamp(0.0, 1.0)
        loss = torch.nn.functional.mse_loss(pred, target)
        self.optimizer.zero_grad()
        loss.backward()
        self.optimizer.step()
        return float(loss.detach().item())
