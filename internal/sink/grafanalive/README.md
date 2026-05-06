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

Two measurements are sent per frame, both tagged by `source_name` and `source_type`.

### `telemetry`

Car-level data fields:

| Field | Type | Description |
|-------|------|-------------|
| `acceleration_x/y/z` | float | m/s^2 |
| `best_lap_time` | float | seconds |
| `boost` | float | |
| `brake` | float | 0.0-1.0 |
| `car_class` | int | |
| `car_index` | int | |
| `car_performance_index` | int | Forza only |
| `clutch` | float | 0.0-1.0 |
| `current_lap_time` | float | seconds |
| `current_race_time` | float | seconds |
| `drivetrain_type` | int | 0=FWD, 1=RWD, 2=AWD |
| `engine_idle_rpm` | float | |
| `engine_max_rpm` | float | |
| `engine_rpm` | float | |
| `fuel` | float | 0.0-1.0 |
| `fuel_capacity` | float | liters |
| `gear` | int | -1=R, 0=N, 1+=forward |
| `hand_brake` | float | 0.0-1.0 |
| `heading` | float | 0.0-1.0, GT7 only |
| `is_race_on` | bool | |
| `lap_distance` | float | meters |
| `lap_number` | int | |
| `last_lap_time` | float | seconds |
| `num_cylinders` | int | |
| `oil_pressure` | float | GT7 only |
| `oil_temp` | float | Celsius |
| `pitch` | float | radians |
| `position_x/y/z` | float | meters |
| `power` | float | watts |
| `race_position` | int | |
| `ride_height` | float | GT7 only |
| `roll` | float | radians |
| `speed` | float | m/s |
| `steer` | float | |
| `suggested_gear` | int | GT7 only |
| `throttle` | float | 0.0-1.0 |
| `torque` | float | Nm |
| `total_laps` | int | |
| `total_positions` | int | GT7 only |
| `velocity_x/y/z` | float | m/s |
| `water_temp` | float | Celsius |
| `yaw` | float | radians |

### `tire`

Per-wheel data, additionally tagged by `wheel` (fl, fr, rl, rr):

| Field | Type | Description |
|-------|------|-------------|
| `combined_slip` | float | Forza only |
| `normalized_suspension` | float | 0.0-1.0, Forza only |
| `puddle_depth` | float | 0.0-1.0, Forza only |
| `slip_angle` | float | |
| `slip_ratio` | float | |
| `surface_rumble` | float | Forza only |
| `suspension` | float | meters |
| `temp` | float | Celsius |
| `wheel_radius` | float | meters, GT7 only |
| `wheel_speed` | float | rad/s |

## Grafana Channels

Data is available on Grafana Live channels:

- `stream/<stream_id>/telemetry`
- `stream/<stream_id>/tire`

Use these channel paths in Grafana Live data source panels for real-time visualization.
