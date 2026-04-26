pub mod agent;
pub mod protocol;
pub mod runtime;
pub mod swarm;

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

pub type AgentId = Uuid;
pub type TaskId = Uuid;
pub type Timestamp = DateTime<Utc>;
pub type Embedding = [f32; 256];

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Task {
    pub id: TaskId,
    pub task_type: String,
    pub priority: u8,
    pub payload: serde_json::Value,
    pub reward: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TaskResult {
    pub task_id: TaskId,
    pub success: bool,
    pub quality_score: f64,
    pub duration_ms: u64,
    pub constitutional_incidents: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Action {
    pub name: String,
    pub risk_score: f64,
    pub requested_cpu_ms: u64,
    pub requested_memory_mb: u64,
    pub target_data_domain: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SwarmEvent {
    pub event_type: String,
    pub agent_id: Option<AgentId>,
    pub payload: serde_json::Value,
    pub ts_unix_ms: i64,
}
