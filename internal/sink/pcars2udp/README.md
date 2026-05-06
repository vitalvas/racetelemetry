# Project CARS 2/3 UDP Sink

Sends telemetry as Project CARS 2 UDP packets. Used to feed telemetry data to CrewChief V4 and other applications that consume the pCars 2 protocol. Also compatible with Project CARS 3.

## Configuration

```yaml
sinks:
  my_pcars2:
    type: pcars2_udp
    inputs:
      - my_source
    target_addr: "127.0.0.1:5606"
```

| Field | Required | Description |
|-------|----------|-------------|
| `target_addr` | yes | UDP address to send packets to |

## Protocol Details

- Packet version: 2
- Max participants: 64
- Telemetry packet: 559 bytes
- Game state packet: 16 bytes
- Timings packet: 1063 bytes
- 12-byte common header on all packets
- Sequence number support (even = valid frame)
- Telemetry sent every frame, game state and timings sent every second
