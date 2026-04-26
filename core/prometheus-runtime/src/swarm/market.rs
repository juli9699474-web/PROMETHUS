use std::{
    collections::{BTreeMap, HashMap, VecDeque},
    sync::Arc,
    time::{Duration, Instant},
};

use atomic_float::AtomicF64;
use tokio::sync::RwLock;

use crate::{AgentId, Task, TaskId};

pub type Priority = u8;

pub struct TaskMarket {
    pub open_tasks: BTreeMap<Priority, VecDeque<Task>>,
    pub active_auctions: HashMap<TaskId, Auction>,
    pub agent_balances: HashMap<AgentId, AtomicF64>,
    pub completed_trades: Arc<RwLock<Vec<Trade>>>,
    pub market_maker: Arc<MarketMaker>,
}

pub struct Auction {
    pub task: Task,
    pub deadline: Instant,
    pub bids: Vec<Bid>,
    pub minimum_bid: f64,
    pub quality_weight: f64,
}

pub struct Bid {
    pub agent_id: AgentId,
    pub price: f64,
    pub estimated_duration: Duration,
    pub confidence: f64,
    pub specialization_match: f64,
}

pub struct Trade {
    pub task_id: TaskId,
    pub winning_agent: AgentId,
    pub winning_bid: f64,
    pub num_bidders: usize,
    pub settled_at: Instant,
}

pub struct MarketMaker {
    pub max_task_share: f64,
}

impl TaskMarket {
    pub fn settle_auction(&mut self, task_id: TaskId, task_share: &HashMap<AgentId, f64>) -> Option<Trade> {
        let auction = self.active_auctions.remove(&task_id)?;
        let mut ranked = auction.bids;
        ranked.sort_by(|a, b| {
            let a_score = self.bid_score(a, auction.quality_weight, task_share);
            let b_score = self.bid_score(b, auction.quality_weight, task_share);
            b_score.partial_cmp(&a_score).unwrap_or(std::cmp::Ordering::Equal)
        });
        let winner = ranked.into_iter().next()?;
        let reward = auction.task.reward.max(auction.minimum_bid);
        if let Some(balance) = self.agent_balances.get(&winner.agent_id) {
            balance.fetch_add(reward, std::sync::atomic::Ordering::Relaxed);
        }
        Some(Trade {
            task_id,
            winning_agent: winner.agent_id,
            winning_bid: winner.price,
            num_bidders: auction.bids.len(),
            settled_at: Instant::now(),
        })
    }

    fn bid_score(&self, bid: &Bid, quality_weight: f64, task_share: &HashMap<AgentId, f64>) -> f64 {
        let anti_monopoly_penalty = if task_share.get(&bid.agent_id).copied().unwrap_or(0.0)
            > self.market_maker.max_task_share
        {
            0.4
        } else {
            0.0
        };
        (bid.specialization_match * quality_weight + bid.confidence * 0.3) - bid.price * 0.05 - anti_monopoly_penalty
    }
}
