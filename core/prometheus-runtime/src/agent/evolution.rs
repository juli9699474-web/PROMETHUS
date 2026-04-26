use std::{collections::HashMap, sync::Arc};

use tokio::sync::RwLock;

use super::genome::AgentGenome;

pub struct SecureSandbox;
pub struct ConstitutionalValidator;

pub struct EvolutionEngine {
    pub population: Arc<RwLock<Vec<AgentGenome>>>,
    pub fitness_history: HashMap<uuid::Uuid, Vec<f64>>,
    pub evolution_strategy: EvolutionStrategy,
    pub sandbox: Arc<SecureSandbox>,
    pub validator: Arc<ConstitutionalValidator>,
}

pub enum EvolutionStrategy {
    GeneticAlgorithm { crossover_rate: f64, mutation_rate: f64 },
    CmaEs,
    Neat,
    Lamarckian,
}

impl EvolutionEngine {
    pub async fn evolve_once(&self) -> anyhow::Result<()> {
        let mut population = self.population.write().await;
        if population.len() < 10 {
            return Ok(());
        }
        population.sort_by(|a, b| {
            let af = a.fitness_score.load(std::sync::atomic::Ordering::Relaxed);
            let bf = b.fitness_score.load(std::sync::atomic::Ordering::Relaxed);
            bf.partial_cmp(&af).unwrap_or(std::cmp::Ordering::Equal)
        });

        let top_n = ((population.len() as f64) * 0.2).ceil() as usize;
        let bottom_n = ((population.len() as f64) * 0.1).ceil() as usize;

        let elite = population[..top_n].to_vec();
        let mut children = Vec::new();
        for pair in elite.chunks(2) {
            if pair.len() == 2 {
                children.push(pair[0].crossover(&pair[1])?);
            }
        }

        let survivors_len = population.len().saturating_sub(bottom_n);
        population.truncate(survivors_len);
        population.extend(children);
        Ok(())
    }
}
