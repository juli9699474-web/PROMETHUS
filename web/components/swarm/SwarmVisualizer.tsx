import { SwarmEvent } from "./types";

export function SwarmVisualizer({ events }: { events: SwarmEvent[] }) {
  const agentCreated = events.filter((e) => e.eventType === "agent.created").length;
  const taskAssigned = events.filter((e) => e.eventType === "task.assigned").length;
  const heartbeats = events.filter((e) => e.eventType === "swarm.heartbeat").length;

  return (
    <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
      <h3 style={{ marginTop: 0 }}>Swarm Visualizer</h3>
      <p style={{ color: "#bbb" }}>Live event topology proxy metrics.</p>
      <div style={{ display: "flex", gap: 16 }}>
        <div>Agents created: {agentCreated}</div>
        <div>Tasks assigned: {taskAssigned}</div>
        <div>Heartbeats: {heartbeats}</div>
      </div>
    </section>
  );
}
