# Wire Sink

Sends telemetry frames to another racetelemetry instance over UDP. Each frame is JSON-encoded and sent as a single UDP packet.

## Configuration

```yaml
sinks:
  downstream:
    type: wire
    inputs:
      - my_source
    target_addr: "192.168.1.50:15000"
```

| Field | Required | Description |
|-------|----------|-------------|
| `target_addr` | yes | UDP address to send frames to |

## Protocol

- Transport: UDP
- Payload: JSON-encoded TelemetryFrame (one per packet)
- Preserves `source_name` and `source_type` for downstream routing

## Use Case

Chain multiple racetelemetry instances:

```
[Instance 1: forza source -> wire sink :15000] -> [Instance 2: wire source :15000 -> crewchief sink]
```
