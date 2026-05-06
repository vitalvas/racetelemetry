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

## Mapped Fields

Fields from the unified telemetry model mapped to pCars 1 packet fields:

| pCars 1 Field | Source Field | Notes |
|---------------|-------------|-------|
| Acceleration [3] | `acceleration_x/y/z` | Local and world |
| Angular velocity [3] | `angular_velocity_x/y/z` | |
| Best/Last/Current lap time | `best_lap_time`, `last_lap_time`, `current_lap_time` | In timings packet |
| Boost | `boost` | |
| Brake (uint8 + float32) | `brake` | |
| Car flags | `is_race_on`, `hand_brake` | ENGINE_ACTIVE, HANDBRAKE |
| Clutch (uint8 + float32) | `clutch` | |
| Fuel capacity | `fuel_capacity` | |
| Fuel level | `fuel` | 0.0-1.0 |
| Gear | `gear` | |
| Lap distance | `lap_distance` | In timings packet |
| Lap number | `lap_number` | In timings packet |
| Max RPM | `engine_max_rpm` | |
| Num gears | `num_gears` | |
| Oil pressure | `oil_pressure` | |
| Oil temp | `oil_temp` | |
| Orientation (yaw/pitch/roll) | `yaw`, `pitch`, `roll` | |
| Position [3] | `position_x/y/z` | |
| RPM | `engine_rpm` | |
| Race position | `race_position` | In timings packet |
| Slip speed [4] | `slip_ratio` | |
| Speed | `speed` | m/s |
| Steering (int8 + float32) | `steer` | |
| Suspension travel [4] | `suspension_travel` | |
| Throttle (uint8 + float32) | `throttle` | |
| Tire temp [4] | `tire_temp` | Also used as brake temp |
| Total laps | `total_laps` | |
| Velocity [3] | `velocity_x/y/z` | Local and world |
| Water temp | `water_temp` | |
| Wheel speed [4] | `wheel_speed` | |

## Protocol Details

- Packet version: 1
- Max participants: 56
- Telemetry packet: 538 bytes
- Game state packet: 16 bytes
- Timings packet: 993 bytes
- 12-byte common header on all packets
- Telemetry sent every frame, game state and timings sent every second
