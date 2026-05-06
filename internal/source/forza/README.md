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

## Available Fields

All protocol versions (V1, V2, V3) provide the same fields.

| Field | Available | Unit | Notes |
|-------|-----------|------|-------|
| `acceleration_x/y/z` | yes | m/s^2 | Car local space |
| `angular_velocity_x/y/z` | yes | rad/s | Car local space |
| `best_lap_time` | yes | seconds | Populates during race events |
| `boost` | yes | PSI | |
| `brake` | yes | 0.0-1.0 | |
| `car_class` | yes | | 0=D, 1=C, 2=B, 3=A, 4=S1, 5=S2, 6=X |
| `car_index` | yes | | Unique car make/model ID |
| `car_performance_index` | yes | | 100-999 |
| `clutch` | yes | 0.0-1.0 | |
| `current_lap_time` | yes | seconds | Populates during race events |
| `current_race_time` | yes | seconds | Session elapsed time |
| `drivetrain_type` | yes | | 0=FWD, 1=RWD, 2=AWD |
| `engine_idle_rpm` | yes | RPM | |
| `engine_max_rpm` | yes | RPM | |
| `engine_rpm` | yes | RPM | |
| `fuel` | yes | 0.0-1.0 | Percentage |
| `gear` | yes | | -1=R, 0=N, 1+=forward |
| `hand_brake` | yes | 0.0-1.0 | |
| `is_race_on` | yes | bool | |
| `lap_distance` | yes | meters | Per-lap, resets each lap |
| `lap_number` | yes | | Populates during race events |
| `last_lap_time` | yes | seconds | Populates during race events |
| `normalized_ai_brake_diff` | yes | | |
| `normalized_driving_line` | yes | | |
| `normalized_suspension_travel` | yes | 0.0-1.0 | Per wheel |
| `num_cylinders` | yes | | |
| `pitch` | yes | radians | Global space |
| `position_x/y/z` | yes | meters | Global space |
| `power` | yes | watts | |
| `race_position` | yes | | Populates during race events |
| `roll` | yes | radians | Global space |
| `slip_angle` | yes | | Per wheel |
| `slip_ratio` | yes | | Per wheel |
| `speed` | yes | m/s | |
| `steer` | yes | -1.0 to 1.0 | |
| `surface_rumble` | yes | | Per wheel |
| `suspension_travel` | yes | meters | Per wheel |
| `throttle` | yes | 0.0-1.0 | |
| `tire_combined_slip` | yes | | Per wheel |
| `tire_temp` | yes | Celsius | Converted from Fahrenheit, per wheel |
| `torque` | yes | Nm | |
| `velocity_x/y/z` | yes | m/s | Car local space |
| `wheel_in_puddle_depth` | yes | 0.0-1.0 | Per wheel |
| `wheel_on_rumble_strip` | yes | bool | Per wheel |
| `wheel_speed` | yes | rad/s | Per wheel |
| `yaw` | yes | radians | Global space |

## Protocol Details

- Transport: UDP, little-endian
- All three versions share the same fields up to offset 228 (NumCylinders)
- V2 has a 12-byte gap at offset 232 before Position fields
- V3 has the same layout as V1 with 20 extra bytes at the end
- Tire temperatures are in Fahrenheit on the wire, converted to Celsius during normalization
- Gear encoding: 0=Reverse, 1=Neutral, 2=1st, normalized to -1=R, 0=N, 1=1st
