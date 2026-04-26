use std::sync::Arc;

use atomic_float::AtomicF64;
use dashmap::DashMap;
use tokio::sync::RwLock;

use crate::{AgentId, Embedding};

#[derive(Clone)]
pub struct WorkspaceEntry {
    pub by_agent: AgentId,
    pub note: String,
    pub confidence: f64,
}

#[derive(Clone)]
pub struct EmergentGoal {
    pub description: String,
    pub supporting_agents: Vec<AgentId>,
    pub confidence: f64,
}

pub struct SwarmConsciousness {
    pub collective_beliefs: DashMap<String, (f64, Vec<AgentId>)>,
    pub attention_vector: Arc<RwLock<Embedding>>,
    pub arousal: AtomicF64,
    pub valence: AtomicF64,
    pub shared_workspace: Arc<RwLock<Vec<WorkspaceEntry>>>,
    pub emergent_goals: Arc<RwLock<Vec<EmergentGoal>>>,
}

impl SwarmConsciousness {
    pub async fn ingest_signal(
        &self,
        concept_key: String,
        agent_id: AgentId,
        confidence: f64,
        workspace_note: String,
    ) {
        self.collective_beliefs
            .entry(concept_key)
            .and_modify(|(c, agents)| {
                *c = ((*c + confidence) / 2.0).clamp(0.0, 1.0);
                if !agents.contains(&agent_id) {
                    agents.push(agent_id);
                }
            })
            .or_insert((confidence, vec![agent_id]));

        self.shared_workspace.write().await.push(WorkspaceEntry {
            by_agent: agent_id,
            note: workspace_note,
            confidence,
        });
    }
}
