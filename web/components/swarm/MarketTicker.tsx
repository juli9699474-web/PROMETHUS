import { SwarmEvent } from "./types";

export function MarketTicker({ events }: { events: SwarmEvent[] }) {
  const taskEvents = events.filter((e) => e.eventType === "task.assigned").slice(0, 10);

  return (
    <section style={{ border: "1px solid #333", borderRadius: 8, padding: 12 }}>
      <h3 style={{ marginTop: 0 }}>Market Ticker</h3>
      {taskEvents.length === 0 ? (
        <p style={{ color: "#aaa" }}>No task auctions yet.</p>
      ) : (
        taskEvents.map((event, i) => (
          <div key={`${event.tsUnixMs}-${i}`} style={{ padding: "6px 0", borderBottom: "1px solid #222" }}>
            {event.eventType} @ {new Date(event.tsUnixMs).toLocaleTimeString()}
          </div>
        ))
      )}
    </section>
  );
}
