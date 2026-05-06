# Forza Source

Receives UDP telemetry from Forza games.

## Supported Games

| Game | Protocol Version | Packet Size |
|------|-----------------|-------------|
| Forza Motorsport 7 | V1 | 311 bytes |
| Forza Horizon 4 | V2 | 324 bytes |
| Forza Horizon 5 (Xbox/PC/PS5) | V2 | 324 bytes |
| Forza Motorsport (2023) | V3 | 331 bytes |

The protocol version is auto-detected by packet size.

## Configuration

```yaml
sources:
  my_forza:
    type: forza
    listen_addr: ":5300"
```

| Field | Required | Description |
|-------|----------|-------------|
| `listen_addr` | yes | UDP address to listen on (e.g. `:5300`, `0.0.0.0:5300`) |

## Game Setup

1. Go to Settings > HUD and Gameplay > Data Out
2. Set Data Out to **On**
3. Set Data Out IP Address to the machine running racetelemetry
4. Set Data Out IP Port to match `listen_addr`

## Protocol Details

- Transport: UDP, little-endian
- All three versions share the same fields up to offset 228 (NumCylinders)
- V2 has a 12-byte gap at offset 232 before Position fields
- V3 has the same layout as V1 with 20 extra bytes at the end
- Tire temperatures are in Fahrenheit on the wire, converted to Celsius during normalization
- Gear encoding: 0=Reverse, 1=Neutral, 2=1st, normalized to -1=R, 0=N, 1=1st
