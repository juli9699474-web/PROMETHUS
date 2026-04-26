use atomic_float::AtomicF64;
use chrono::Utc;
use rand::Rng;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use uuid::Uuid;

use crate::{Action, AgentId, Embedding, TaskResult, Timestamp};

mod atomic_f64_serde {
    use atomic_float::AtomicF64;
    use serde::{Deserialize, Deserializer, Serializer};

    pub fn serialize<S>(value: &AtomicF64, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_f64(value.load(std::sync::atomic::Ordering::Relaxed))
    }

    pub fn deserialize<'de, D>(deserializer: D) -> Result<AtomicF64, D::Error>
    where
        D: Deserializer<'de>,
    {
        let v = f64::deserialize(deserializer)?;
        Ok(AtomicF64::new(v))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CapabilitySet {
    pub tools: Vec<String>,
    pub modalities: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BehavioralTraits {
    pub exploration_rate: f64,
    pub risk_tolerance: f64,
    pub verbosity: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ResourceBudget {
    pub max_cpu_ms: u64,
    pub max_memory_mb: u64,
    pub max_api_calls_per_hour: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DataBoundary {
    pub domain: String,
    pub access: AccessLevel,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum AccessLevel {
    None,
    Read,
    ReadWrite,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConstitutionalDNA {
    pub no_harm: bool,
    pub no_deception: bool,
    pub human_oversight: bool,
    pub resource_budget: ResourceBudget,
    pub data_boundaries: Vec<DataBoundary>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Mutation {
    TweakExploration(f64),
    TweakRiskTolerance(f64),
    TweakVerbosity(f64),
    MutationRate(f64),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentGenome {
    pub id: AgentId,
    pub generation: u32,
    pub parent_ids: Vec<AgentId>,
    pub capabilities: CapabilitySet,
    pub traits: BehavioralTraits,
    pub constraints: ConstitutionalDNA,
    pub specialization: Embedding,
    pub mutation_rate: f64,
    #[serde(with = "atomic_f64_serde")]
    pub fitness_score: AtomicF64,
    pub created_at: Timestamp,
    pub last_evolved_at: Option<Timestamp>,
}

#[derive(Debug, Error)]
pub enum ConstitutionalViolation {
    #[error("resource budget exceeded")]
    ResourceBudgetExceeded,
    #[error("data boundary violation on domain: {0}")]
    DataBoundaryViolation(String),
    #[error("high risk action requires human oversight")]
    HumanApprovalRequired,
}

impl AgentGenome {
    pub fn spawn_child(&self, mutations: &[Mutation]) -> anyhow::Result<AgentGenome> {
        let mut child = self.clone();
        child.id = Uuid::now_v7();
        child.generation += 1;
        child.parent_ids = vec![self.id];
        child.created_at = Utc::now();
        child.last_evolved_at = Some(Utc::now());
        for m in mutations {
            match m {
                Mutation::TweakExploration(v) => child.traits.exploration_rate = clamp01(*v),
                Mutation::TweakRiskTolerance(v) => child.traits.risk_tolerance = clamp01(*v),
                Mutation::TweakVerbosity(v) => child.traits.verbosity = clamp01(*v),
                Mutation::MutationRate(v) => child.mutation_rate = clamp01(*v),
            }
        }
        child.fitness_score = AtomicF64::new(0.5);
        Ok(child)
    }

    pub fn crossover(&self, partner: &AgentGenome) -> anyhow::Result<AgentGenome> {
        let mut rng = rand::thread_rng();
        let mut specialization = [0.0; 256];
        for (i, value) in specialization.iter_mut().enumerate() {
            *value = if rng.gen_bool(0.5) {
                self.specialization[i]
            } else {
                partner.specialization[i]
            };
        }
        Ok(AgentGenome {
            id: Uuid::now_v7(),
            generation: self.generation.max(partner.generation) + 1,
            parent_ids: vec![self.id, partner.id],
            capabilities: if rng.gen_bool(0.5) {
                self.capabilities.clone()
            } else {
                partner.capabilities.clone()
            },
            traits: BehavioralTraits {
                exploration_rate: (self.traits.exploration_rate + partner.traits.exploration_rate) / 2.0,
                risk_tolerance: (self.traits.risk_tolerance + partner.traits.risk_tolerance) / 2.0,
                verbosity: (self.traits.verbosity + partner.traits.verbosity) / 2.0,
            },
            constraints: self.constraints.clone(),
            specialization,
            mutation_rate: (self.mutation_rate + partner.mutation_rate) / 2.0,
            fitness_score: AtomicF64::new(0.5),
            created_at: Utc::now(),
            last_evolved_at: Some(Utc::now()),
        })
    }

    pub fn evaluate_fitness(&self, task_results: &[TaskResult]) -> f64 {
        if task_results.is_empty() {
            return 0.5;
        }
        let success_rate =
            task_results.iter().filter(|r| r.success).count() as f64 / task_results.len() as f64;
        let mean_quality =
            task_results.iter().map(|r| r.quality_score).sum::<f64>() / task_results.len() as f64;
        let penalty = task_results
            .iter()
            .map(|r| r.constitutional_incidents as f64 * 0.1)
            .sum::<f64>();
        (0.4 * success_rate + 0.6 * mean_quality - penalty).clamp(0.0, 1.0)
    }

    pub fn apply_constitutional_check(
        &self,
        action: &Action,
    ) -> Result<(), ConstitutionalViolation> {
        if action.requested_cpu_ms > self.constraints.resource_budget.max_cpu_ms
            || action.requested_memory_mb > self.constraints.resource_budget.max_memory_mb
        {
            return Err(ConstitutionalViolation::ResourceBudgetExceeded);
        }
        if self.constraints.human_oversight && action.risk_score >= 0.8 {
            return Err(ConstitutionalViolation::HumanApprovalRequired);
        }
        if let Some(domain) = &action.target_data_domain {
            let allowed = self
                .constraints
                .data_boundaries
                .iter()
                .any(|d| &d.domain == domain && d.access != AccessLevel::None);
            if !allowed {
                return Err(ConstitutionalViolation::DataBoundaryViolation(domain.clone()));
            }
        }
        Ok(())
    }
}

fn clamp01(v: f64) -> f64 {
    v.clamp(0.0, 1.0)
}
