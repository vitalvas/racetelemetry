# racetelemetry

A telemetry pipeline for racing games. Receives telemetry data from game sources, normalizes it, and routes it to output sinks -- like [vector.dev](https://vector.dev) but for racing telemetry.

## Features

- Multiple named source and sink instances
- Input routing: each sink subscribes to specific sources
- Auto-detection of protocol versions
- Unified internal telemetry model across all games
- Each frame carries `source_name` (config key) and `source_type` (e.g. `forza`, `gt7`)
- Pointer-based fields: `nil` means the source does not provide the data, `omitempty` in JSON output

## Supported Sources

| Type | Games | Protocol |
|------|-------|----------|
| `forza` | Forza Motorsport 7, Forza Horizon 4/5 (Xbox/PC/PS5), Forza Motorsport (2023) | UDP, auto-detects V1 (311B), V2 (324B), V3 (331B) |
| `gt7` | Gran Turismo 7, Gran Turismo Sport, Gran Turismo 6 | UDP with Salsa20 encryption, heartbeat keep-alive, format C (368B) |

## Supported Sinks

| Type | Description | Platform |
|------|-------------|----------|
| `pcars1_udp` | Project CARS 1 UDP telemetry | Cross-platform |
| `pcars1_shm` | Project CARS 1 shared memory (`$pcars$`) | Windows |
| `pcars2_udp` | Project CARS 2/3 UDP telemetry | Cross-platform |
| `pcars2_shm` | Project CARS 2/3 shared memory (`$pcars2$`) | Windows |
| `json_stdout` | JSON lines to stdout | Cross-platform |
| `grafana_live` | Grafana Live push API (Influx line protocol) | Cross-platform |

## Configuration

```yaml
log:
  level: info    # debug, info, warn, error
  format: text   # text, json

sources:
  forza_xbox:
    type: forza
    listen_addr: ":5300"

  forza_pc:
    type: forza
    listen_addr: ":5301"

  forza_ps5:
    type: forza
    listen_addr: ":5302"

  gt7_ps5:
    type: gt7
    console_addr: "192.168.1.100"
    listen_addr: ":33740"

sinks:
  crewchief_forza:
    type: pcars2_udp
    inputs:
      - forza_xbox
      - forza_pc
      - forza_ps5
    target_addr: "127.0.0.1:5606"

  crewchief_gt7:
    type: pcars2_udp
    inputs:
      - gt7_ps5
    target_addr: "127.0.0.1:5607"

  debug:
    type: json_stdout
    inputs:
      - forza_xbox

  grafana:
    type: grafana_live
    inputs:
      - forza_xbox
      - gt7_ps5
    endpoint: "http://localhost:3000/api/live/push/race_telemetry"
    api_key: "glsa_xxxxxxxxxxxx"
```

## Usage

```
racetelemetry -config config.yaml
```

## Game Setup

### Forza Horizon 4/5 / Forza Motorsport

1. Go to Settings > HUD and Gameplay > Data Out
2. Set Data Out to **On**
3. Set Data Out IP Address to the machine running racetelemetry
4. Set Data Out IP Port to match `listen_addr` in config

### Gran Turismo 7

1. Set `console_addr` in config to the PlayStation IP address
2. The source sends heartbeat packets automatically to enable telemetry output
