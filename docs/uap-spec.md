# Universal Agent Protocol (UAP) v1.0

`UAP/1.0` is a transport-agnostic, signed agent-to-agent protocol.

## Envelope

- `version`: fixed string `UAP/1.0`
- `from_agent`: sender ID plus key identity
- `to_agent`: receiver ID or `swarm`
- `type`: message enum
- `payload`: binary protobuf body
- `signature`: Ed25519 signature over envelope with empty signature field
- `timestamp`: unix epoch seconds
- `nonce`: unique anti-replay token

## Message Types

- `CAPABILITY_ANNOUNCE`
- `TASK_REQUEST`
- `TASK_ACCEPT`
- `TASK_RESULT`
- `KNOWLEDGE_SHARE`
- `HELP_REQUEST`
- `EVOLUTION_PROPOSAL`
- `CONSENSUS_VOTE`

## Compatibility Goals

- Works over gRPC, WebSocket, MQTT, HTTP/2, and raw TCP
- Cross-language clients for Rust, Python, Go, and TypeScript
- Self-describing startup via capability announcement handshake
