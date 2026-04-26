use std::sync::Arc;

use chrono::Utc;
use tokio::sync::broadcast;
use uuid::Uuid;

use crate::SwarmEvent;

#[derive(Clone)]
pub struct EventBus {
    tx: Arc<broadcast::Sender<SwarmEvent>>,
}

impl EventBus {
    pub fn new(capacity: usize) -> Self {
        let (tx, _) = broadcast::channel(capacity);
        Self { tx: Arc::new(tx) }
    }

    pub fn subscribe(&self) -> broadcast::Receiver<SwarmEvent> {
        self.tx.subscribe()
    }

    pub fn publish(&self, event: SwarmEvent) {
        let _ = self.tx.send(event);
    }

    pub fn publish_heartbeat(&self) {
        self.publish(SwarmEvent {
            event_type: "swarm.heartbeat".to_string(),
            agent_id: Some(Uuid::nil()),
            payload: serde_json::json!({"status":"alive"}),
            ts_unix_ms: Utc::now().timestamp_millis(),
        });
    }
}
