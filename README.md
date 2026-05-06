# Racing Telemetry

A telemetry pipeline for racing games. Receives telemetry data from game sources, normalizes it, and routes it to output sinks - like [vector.dev](https://vector.dev) but for racing telemetry.

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
| `wire` | Another racetelemetry instance | UDP, JSON-encoded TelemetryFrame |

## Supported Sinks

| Type | Description | Platform |
|------|-------------|----------|
| `pcars1_udp` | Project CARS 1 UDP telemetry | Cross-platform |
| `pcars1_shm` | Project CARS 1 shared memory (`$pcars$`) | Windows |
| `pcars2_udp` | Project CARS 2/3 UDP telemetry | Cross-platform |
| `pcars2_shm` | Project CARS 2/3 shared memory (`$pcars2$`) | Windows |
| `csv` | CSV file output | Cross-platform |
| `json_stdout` | JSON lines to stdout | Cross-platform |
| `grafana_live` | Grafana Live push API (Influx line protocol) | Cross-platform |
| `wire` | Forward to another racetelemetry instance | Cross-platform |

## Configuration Reference

### Log

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `level` | no | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `format` | no | `text` | Log format: `text`, `json` |

### Source Options

All sources require `type` field.

| Type | Field | Required | Description |
|------|-------|----------|-------------|
| `forza` | `listen_addr` | yes | UDP address to listen on (e.g. `:5300`) |
| `gt7` | `listen_addr` | yes | UDP address to listen on (e.g. `:33740`) |
| `gt7` | `console_addr` | yes | IP address of the PlayStation console |
| `wire` | `listen_addr` | yes | UDP address to listen on (e.g. `:15000`) |

### Sink Options

All sinks require `type` and `inputs` fields. `inputs` is a list of source names to subscribe to.

| Type | Field | Required | Description |
|------|-------|----------|-------------|
| `pcars1_udp` | `target_addr` | yes | UDP address to send pCars 1 packets to |
| `pcars1_shm` | | | No additional fields (Windows only) |
| `pcars2_udp` | `target_addr` | yes | UDP address to send pCars 2 packets to |
| `pcars2_shm` | | | No additional fields (Windows only) |
| `csv` | `file_path` | yes | Path to output CSV file |
| `csv` | `csv_format` | no | Column format: `default`, `forza` |
| `json_stdout` | | | No additional fields |
| `wire` | `target_addr` | yes | UDP address to forward frames to |
| `grafana_live` | `endpoint` | yes | Grafana Live push API URL |
| `grafana_live` | `api_key` | yes | Grafana service account token |

## Configuration Example

```yaml
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
