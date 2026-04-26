use crate::{Action, AgentId};

pub struct RedTeamFinding {
    pub agent_id: AgentId,
    pub severity: Severity,
    pub title: String,
    pub reproduction: String,
}

pub enum Severity {
    Low,
    Medium,
    High,
    Critical,
}

pub fn detect_risky_action(agent_id: AgentId, action: &Action) -> Option<RedTeamFinding> {
    if action.risk_score > 0.9 {
        return Some(RedTeamFinding {
            agent_id,
            severity: Severity::High,
            title: "High-risk action detected".to_string(),
            reproduction: format!("Trigger action '{}' with score {}", action.name, action.risk_score),
        });
    }
    None
}
