# Project CARS 1 UDP Sink

Sends telemetry as Project CARS 1 UDP packets. Used to feed telemetry data to applications that consume the pCars 1 protocol, such as older versions of CrewChief.

## Configuration

```yaml
sinks:
  my_pcars1:
    type: pcars1_udp
    inputs:
      - my_source
    target_addr: "127.0.0.1:5606"
```

| Field | Required | Description |
|-------|----------|-------------|
| `target_addr` | yes | UDP address to send packets to |

## Protocol Details

- Packet version: 1
- Max participants: 56
- Telemetry packet: 538 bytes
- Game state packet: 16 bytes
- Timings packet: 993 bytes
- 12-byte common header on all packets
- Telemetry sent every frame, game state and timings sent every second
