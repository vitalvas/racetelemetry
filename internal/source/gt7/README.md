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

## Available Fields

Fields depend on the packet format the console sends. The source requests format C (richest). Older games may respond with smaller formats.

| Field | Standard (A) | Addendum1 (B) | Addendum3 (C) | Unit | Notes |
|-------|:---:|:---:|:---:|------|-------|
| `angular_velocity_x/y/z` | yes | yes | yes | rad/s | |
| `best_lap_time` | yes | yes | yes | seconds | |
| `boost` | yes | yes | yes | bar | Relative to atmospheric |
| `brake` | yes | yes | yes | 0.0-1.0 | |
| `car_index` | yes | yes | yes | | Vehicle ID |
| `clutch` | yes | yes | yes | 0.0-1.0 | |
| `clutch_engaged` | yes | yes | yes | 0.0-1.0 | |
| `current_lap_time` | no | no | yes | seconds | Format C only |
| `engine_max_rpm` | yes | yes | yes | RPM | From rev limiter |
| `engine_rpm` | yes | yes | yes | RPM | |
| `estimated_top_speed` | yes | yes | yes | km/h | |
| `fuel` | conditional | conditional | conditional | 0.0-1.0 | Only when fuel_capacity > 0 |
| `fuel_capacity` | yes | yes | yes | liters | |
| `gear` | yes | yes | yes | | -1=R, 0=N, 1+=forward |
| `gear_ratios` | yes | yes | yes | | 8 gear ratios |
| `heading` | yes | yes | yes | | 0.0=south, 1.0=north |
| `is_race_on` | yes | yes | yes | bool | From status flags bit 0 |
| `lap_number` | yes | yes | yes | | |
| `last_lap_time` | yes | yes | yes | seconds | |
| `oil_pressure` | yes | yes | yes | | |
| `oil_temp` | yes | yes | yes | Celsius | |
| `pitch` | yes | yes | yes | radians | |
| `position_x/y/z` | yes | yes | yes | meters | |
| `race_position` | yes | yes | yes | | |
| `ride_height` | yes | yes | yes | | |
| `road_plane_dist` | yes | yes | yes | | |
| `road_plane_x/y/z` | yes | yes | yes | | Road surface normal |
| `roll` | yes | yes | yes | radians | |
| `rpm_after_clutch` | yes | yes | yes | RPM | |
| `speed` | yes | yes | yes | m/s | |
| `steer` | no | yes | yes | radians | Steering wheel angle |
| `steer_velocity` | no | yes | yes | rad/s | |
| `suggested_gear` | yes | yes | yes | | |
| `suspension_travel` | yes | yes | yes | meters | Per wheel |
| `throttle` | yes | yes | yes | 0.0-1.0 | |
| `tire_temp` | yes | yes | yes | Celsius | Per wheel |
| `top_speed_ratio` | yes | yes | yes | | |
| `total_laps` | yes | yes | yes | | |
| `total_positions` | yes | yes | yes | | Total cars |
| `velocity_x/y/z` | yes | yes | yes | m/s | |
| `water_temp` | yes | yes | yes | Celsius | |
| `wheel_radius` | yes | yes | yes | meters | Per wheel |
| `wheel_speed` | yes | yes | yes | rad/s | Per wheel |
| `yaw` | yes | yes | yes | radians | |

## Packet Formats

| Format | Heartbeat | Size | IV Seed | Extra Fields |
|--------|-----------|------|---------|-------------|
| Standard | `A` | 296B | `0xDEADBEAF` | Base telemetry |
| Addendum1 | `B` | 316B | `0xDEADBEEF` | + steering wheel angle/velocity |
| Addendum2 | `~` | 344B | `0x55FABB4F` | + throttle input, brake output, energy recovery |
| Addendum3 | `C` | 368B | `0xDEADBEEF` | + current lap time, surface type, wheel steering angles |

The source requests format C by default. Smaller packets from older games are handled gracefully with addendum fields left as nil.

## Protocol Details

- Transport: UDP, little-endian, Salsa20 encrypted
- Key: first 32 bytes of `"Simulator Interface Packet GT7 ver 0.0"`
- Magic number after decryption: `0x47375330` ("G7S0")
- Heartbeat: sends `"C"` byte to console on port 33739 every 10 seconds
