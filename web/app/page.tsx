"use client";

import { useEffect, useState } from "react";
import { ConsciousnessMap } from "../components/swarm/ConsciousnessMap";
import { MarketTicker } from "../components/swarm/MarketTicker";
import { SwarmVisualizer } from "../components/swarm/SwarmVisualizer";
import { SwarmEvent } from "../components/swarm/types";

type Agent = {
  id: string;
  name: string;
  status: string;
  fitnessScore: number;
};

type Task = {
  id: string;
  agentId: string;
  title: string;
  status: string;
  reward: number;
};

type ReplayMetrics = {
  totalQueries: number;
  rejectedQueries: number;
  totalReturnedEvents: number;
  avgLatencyMs: number;
};

type ReplayPercentiles = {
  p50: number;
  p95: number;
  p99: number;
};

type ReplayPreset = {
  id: string;
  label: string;
  eventType: string;
  minutesBack: number;
};

const replayPresets: ReplayPreset[] = [
  { id: "all-15", label: "All Events (15m)", eventType: "", minutesBack: 15 },
  { id: "heartbeats-5", label: "Heartbeats (5m)", eventType: "swarm.heartbeat", minutesBack: 5 },
  { id: "tasks-60", label: "Task Activity (60m)", eventType: "task.assigned", minutesBack: 60 },
  { id: "agents-60", label: "Agent Lifecycle (60m)", eventType: "agent.created", minutesBack: 60 }
];

const PRESET_STORAGE_KEY = "prometheus.replayPreset.v1";

export default function HomePage() {
  const [events, setEvents] = useState<SwarmEvent[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [liveMode, setLiveMode] = useState(true);
  const [eventTypeFilter, setEventTypeFilter] = useState("");
  const [filterAgentId, setFilterAgentId] = useState("");
  const [sinceInput, setSinceInput] = useState("");
  const [untilInput, setUntilInput] = useState("");
  const [nextCursorMs, setNextCursorMs] = useState<number>(0);
  const [agentName, setAgentName] = useState("Genesis Agent");
  const [agentId, setAgentId] = useState("");
  const [taskTitle, setTaskTitle] = useState("Investigate market anomaly");
  const [activePresetId, setActivePresetId] = useState("all-15");
  const [replayMetrics, setReplayMetrics] = useState<ReplayMetrics | null>(null);
  const [latencyTrend, setLatencyTrend] = useState<number[]>([]);
  const [percentiles, setPercentiles] = useState<ReplayPercentiles>({ p50: 0, p95: 0, p99: 0 });

  useEffect(() => {
    async function loadInitialState() {
      const params = new URLSearchParams({ limit: "100" });
      if (eventTypeFilter) params.set("eventType", eventTypeFilter);
      if (filterAgentId) params.set("agentId", filterAgentId);
      if (sinceInput) params.set("sinceMs", String(new Date(sinceInput).getTime()));
      if (untilInput) params.set("untilMs", String(new Date(untilInput).getTime()));
      const [eventRes, agentRes, taskRes] = await Promise.all([
        fetch(`http://localhost:3001/api/v1/events?${params.toString()}`),
        fetch("http://localhost:3001/api/v1/agents?limit=50"),
        fetch("http://localhost:3001/api/v1/tasks?limit=50")
      ]);
      if (eventRes.ok) {
        const data = await eventRes.json();
        setEvents((data.items ?? []) as SwarmEvent[]);
        setNextCursorMs((data.nextCursorMs as number) ?? 0);
      }
      if (agentRes.ok) {
        const data = await agentRes.json();
        setAgents((data.items ?? []) as Agent[]);
      }
      if (taskRes.ok) {
        const data = await taskRes.json();
        setTasks((data.items ?? []) as Task[]);
      }
    }
    void loadInitialState();

    if (!liveMode) return;
    const ws = new WebSocket("ws://localhost:3001/ws/swarm");
    ws.onmessage = (message) => {
      try {
        const parsed = JSON.parse(message.data) as SwarmEvent;
        if (eventTypeFilter && parsed.eventType !== eventTypeFilter) return;
        setEvents((prev) => [parsed, ...prev].slice(0, 100));
      } catch {
        // Ignore malformed websocket messages.
      }
    };
    return () => ws.close();
  }, [liveMode, eventTypeFilter, filterAgentId, sinceInput, untilInput]);

  useEffect(() => {
    const timer = setInterval(async () => {
      const [summaryRes, seriesRes, pctRes] = await Promise.all([
        fetch("http://localhost:3001/api/v1/metrics/replay"),
        fetch("http://localhost:3001/api/v1/metrics/replay/timeseries?limit=24"),
        fetch("http://localhost:3001/api/v1/metrics/replay/percentiles?limit=500")
      ]);
      if (summaryRes.ok) {
        const data = (await summaryRes.json()) as ReplayMetrics;
        setReplayMetrics(data);
      }
      if (seriesRes.ok) {
        const series = await seriesRes.json();
        const latencies = ((series.items ?? []) as Array<{ latencyMs: number; rejected: boolean }>)
          .filter((p) => !p.rejected)
          .map((p) => p.latencyMs)
          .slice(0, 24)
          .reverse();
        setLatencyTrend(latencies);
      }
      if (pctRes.ok) {
        const pct = (await pctRes.json()) as ReplayPercentiles;
        setPercentiles(pct);
      }
    }, 3000);
    return () => clearInterval(timer);
  }, []);

  async function loadMoreHistory() {
    if (!nextCursorMs) return;
    const params = new URLSearchParams({ limit: "100", cursorMs: String(nextCursorMs) });
    if (eventTypeFilter) params.set("eventType", eventTypeFilter);
    if (filterAgentId) params.set("agentId", filterAgentId);
    if (sinceInput) params.set("sinceMs", String(new Date(sinceInput).getTime()));
    if (untilInput) params.set("untilMs", String(new Date(untilInput).getTime()));
    const res = await fetch(`http://localhost:3001/api/v1/events?${params.toString()}`);
    if (!res.ok) return;
    const data = await res.json();
    setEvents((prev) => [...prev, ...((data.items ?? []) as SwarmEvent[])]);
    setNextCursorMs((data.nextCursorMs as number) ?? 0);
  }

  function applyPreset(preset: ReplayPreset) {
    setActivePresetId(preset.id);
    setEventTypeFilter(preset.eventType);
    const now = new Date();
    const since = new Date(now.getTime() - preset.minutesBack * 60_000);
    setSinceInput(since.toISOString().slice(0, 16));
    setUntilInput(now.toISOString().slice(0, 16));
    try {
      localStorage.setItem(PRESET_STORAGE_KEY, preset.id);
    } catch {
      // Ignore storage errors.
    }
  }

  useEffect(() => {
    try {
      const saved = localStorage.getItem(PRESET_STORAGE_KEY);
      const preset = replayPresets.find((p) => p.id === saved);
      if (preset) applyPreset(preset);
    } catch {
      // Ignore storage errors.
    }
  }, []);

  async function createAgent() {
    const response = await fetch("http://localhost:3001/api/v1/agents", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: agentName })
    });
    if (!response.ok) return;
    const data = await response.json();
    setAgentId(data.id);
    setAgents((prev) => [data as Agent, ...prev].slice(0, 50));
  }

  async function assignTask() {
    if (!agentId) return;
    await fetch(`http://localhost:3001/api/v1/agents/${agentId}/tasks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        title: taskTitle,
        description: "Phase-3 runtime wiring task",
        reward: 12.5
      })
    });
    const taskRes = await fetch("http://localhost:3001/api/v1/tasks?limit=50");
    if (taskRes.ok) {
      const data = await taskRes.json();
      setTasks((data.items ?? []) as Task[]);
    }
  }

  return (
    <main style={{ padding: 24 }}>
      <h1 style={{ marginTop: 0 }}>PROMETHEUS Mission Control</h1>
      <p>Phase-6 controls: live stream plus filtered replay.</p>
      <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12, marginBottom: 16 }}>
        <h3 style={{ marginTop: 0 }}>Playback Controls</h3>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          <button onClick={() => setLiveMode(true)} disabled={liveMode}>Live Mode</button>
          <button onClick={() => setLiveMode(false)} disabled={!liveMode}>Replay Mode</button>
          <input
            placeholder="eventType filter"
            value={eventTypeFilter}
            onChange={(e) => setEventTypeFilter(e.target.value)}
          />
          <input
            placeholder="agentId filter"
            value={filterAgentId}
            onChange={(e) => setFilterAgentId(e.target.value)}
          />
          <input
            type="datetime-local"
            value={sinceInput}
            onChange={(e) => setSinceInput(e.target.value)}
          />
          <input
            type="datetime-local"
            value={untilInput}
            onChange={(e) => setUntilInput(e.target.value)}
          />
          <button onClick={loadMoreHistory} disabled={liveMode || !nextCursorMs}>Load More History</button>
        </div>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap", marginTop: 8 }}>
          {replayPresets.map((preset) => (
            <button
              key={preset.id}
              onClick={() => applyPreset(preset)}
              style={{
                border: activePresetId === preset.id ? "1px solid #7c3aed" : "1px solid #333",
                background: activePresetId === preset.id ? "#21103d" : "#111",
                color: "#fff"
              }}
            >
              {preset.label}
            </button>
          ))}
        </div>
      </section>
      <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12, marginBottom: 16 }}>
        <h3 style={{ marginTop: 0 }}>Control Panel</h3>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap", marginBottom: 8 }}>
          <input value={agentName} onChange={(e) => setAgentName(e.target.value)} />
          <button onClick={createAgent}>Create Agent</button>
        </div>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          <input value={taskTitle} onChange={(e) => setTaskTitle(e.target.value)} />
          <button onClick={assignTask} disabled={!agentId}>
            Assign Task
          </button>
          <span style={{ color: "#bbb" }}>Agent ID: {agentId || "none"}</span>
        </div>
      </section>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(3, minmax(0, 1fr))", gap: 16, marginBottom: 16 }}>
        <SwarmVisualizer events={events} />
        <MarketTicker events={events} />
        <ConsciousnessMap events={events} />
      </div>
      <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12, marginBottom: 16 }}>
        <h3 style={{ marginTop: 0 }}>Replay Metrics</h3>
        <div style={{ display: "flex", gap: 16, flexWrap: "wrap", marginBottom: 8 }}>
          <div>Total queries: {replayMetrics?.totalQueries ?? 0}</div>
          <div>Rejected: {replayMetrics?.rejectedQueries ?? 0}</div>
          <div>Avg latency: {(replayMetrics?.avgLatencyMs ?? 0).toFixed(1)} ms</div>
          <div>P50: {(percentiles.p50 ?? 0).toFixed(1)} ms</div>
          <div>P95: {(percentiles.p95 ?? 0).toFixed(1)} ms</div>
          <div>P99: {(percentiles.p99 ?? 0).toFixed(1)} ms</div>
        </div>
        <div style={{ display: "flex", alignItems: "flex-end", gap: 3, height: 60 }}>
          {latencyTrend.map((v, i) => (
            <div
              key={i}
              style={{
                width: 8,
                height: `${Math.max(4, Math.min(56, v))}px`,
                background: "#7c3aed"
              }}
            />
          ))}
        </div>
      </section>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(2, minmax(0, 1fr))", gap: 16, marginBottom: 16 }}>
        <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
          <h3 style={{ marginTop: 0 }}>Persisted Agents</h3>
          {agents.length === 0 ? (
            <p style={{ color: "#aaa" }}>No agents yet.</p>
          ) : (
            agents.map((agent) => (
              <div key={agent.id} style={{ padding: "6px 0", borderBottom: "1px solid #222" }}>
                <strong>{agent.name}</strong> - {agent.status}
              </div>
            ))
          )}
        </section>
        <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
          <h3 style={{ marginTop: 0 }}>Persisted Tasks</h3>
          {tasks.length === 0 ? (
            <p style={{ color: "#aaa" }}>No tasks yet.</p>
          ) : (
            tasks.map((task) => (
              <div key={task.id} style={{ padding: "6px 0", borderBottom: "1px solid #222" }}>
                <strong>{task.title}</strong> - {task.status}
              </div>
            ))
          )}
        </section>
      </div>
      <div style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
        <h3 style={{ marginTop: 0 }}>Raw Event Stream</h3>
        {events.length === 0 ? (
          <p style={{ color: "#aaa" }}>Waiting for swarm events...</p>
        ) : (
          events.map((event, idx) => (
            <div key={`${event.tsUnixMs}-${idx}`} style={{ padding: "8px 0", borderBottom: "1px solid #222" }}>
              <strong>{event.eventType}</strong> @ {new Date(event.tsUnixMs).toLocaleTimeString()}
            </div>
          ))
        )}
      </div>
    </main>
  );
}
