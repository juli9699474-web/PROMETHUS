import { SwarmEvent } from "./types";

export function ConsciousnessMap({ events }: { events: SwarmEvent[] }) {
  const recent = events.slice(0, 8);

  return (
    <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
      <h3 style={{ marginTop: 0 }}>Consciousness Map</h3>
      <p style={{ color: "#bbb" }}>Recent swarm thought-stream events.</p>
      {recent.length === 0 ? (
        <p style={{ color: "#aaa" }}>No shared cognition yet.</p>
      ) : (
        recent.map((event, i) => (
          <div key={`${event.tsUnixMs}-${i}`} style={{ padding: "4px 0" }}>
            {event.eventType}
          </div>
        ))
      )}
    </section>
  );
}
