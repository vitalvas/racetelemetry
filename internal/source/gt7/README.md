# Gran Turismo Source

Receives UDP telemetry from Gran Turismo games.

## Supported Games

- Gran Turismo 7
- Gran Turismo Sport
- Gran Turismo 6

All games use the same Salsa20-encrypted UDP protocol with the same key.

## Configuration

```yaml
sources:
  my_gt7:
    type: gt7
    console_addr: "192.168.1.100"
    listen_addr: ":33740"
```

| Field | Required | Description |
|-------|----------|-------------|
| `console_addr` | yes | IP address of the PlayStation console |
| `listen_addr` | yes | UDP address to listen on (default port is 33740) |

## Protocol Details

- Transport: UDP, little-endian, Salsa20 encrypted
- Key: first 32 bytes of `"Simulator Interface Packet GT7 ver 0.0"`
- Magic number after decryption: `0x47375330` ("G7S0")
- Heartbeat: sends `"C"` byte to console on port 33739 every 10 seconds
- Uses format C (Addendum3, 368 bytes) for the richest telemetry data

## Packet Formats

| Format | Heartbeat | Size | IV Seed | Extra Fields |
|--------|-----------|------|---------|-------------|
| Standard | `A` | 296B | `0xDEADBEAF` | Base telemetry |
| Addendum1 | `B` | 316B | `0xDEADBEEF` | + steering wheel angle/velocity |
| Addendum2 | `~` | 344B | `0x55FABB4F` | + throttle input, brake output, energy recovery |
| Addendum3 | `C` | 368B | `0xDEADBEEF` | + current lap time, surface type, wheel steering angles |

The source requests format C by default. Smaller packets from older games are handled gracefully with addendum fields left as zero.
