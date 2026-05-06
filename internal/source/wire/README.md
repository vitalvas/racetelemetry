# Wire Source

Receives telemetry frames from another racetelemetry instance over UDP. Each UDP packet contains one JSON-encoded TelemetryFrame.

## Configuration

```yaml
sources:
  upstream:
    type: wire
    listen_addr: ":15000"
```

| Field | Required | Description |
|-------|----------|-------------|
| `listen_addr` | yes | UDP address to listen on |

## Protocol

- Transport: UDP
- Payload: JSON-encoded TelemetryFrame (one per packet)
- Preserves `source_name` and `source_type` from the upstream instance

## Use Case

Chain multiple racetelemetry instances:

```
[Game] -> [Instance 1: forza source -> wire sink] -> [Instance 2: wire source -> crewchief sink]
```
