# Grafana Live Sink

Pushes telemetry to Grafana Live using the HTTP push API with Influx line protocol format. Enables real-time telemetry dashboards in Grafana.

## Configuration

```yaml
sinks:
  grafana:
    type: grafana_live
    inputs:
      - my_source
    endpoint: "http://localhost:3000/api/live/push/race_telemetry"
    api_key: "glsa_xxxxxxxxxxxx"
```

| Field | Required | Description |
|-------|----------|-------------|
| `endpoint` | yes | Grafana Live push API URL (`/api/live/push/<stream_id>`) |
| `api_key` | yes | Grafana service account token with admin permissions |

## Measurements

Two measurements are sent per frame:

### `telemetry`

Car-level data: RPM, speed, throttle, brake, clutch, steer, gear, boost, fuel, position, orientation, lap info, temperatures.

### `tire`

Per-wheel data tagged by `wheel` (fl, fr, rl, rr): temperature, suspension travel, wheel speed, slip ratio, slip angle.

## Grafana Channels

Data is available on Grafana Live channels:

- `stream/<stream_id>/telemetry`
- `stream/<stream_id>/tire`

Use these channel paths in Grafana Live data source panels for real-time visualization.
