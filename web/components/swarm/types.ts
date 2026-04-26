export type SwarmEvent = {
  eventType: string;
  payload: Record<string, unknown>;
  tsUnixMs: number;
};
