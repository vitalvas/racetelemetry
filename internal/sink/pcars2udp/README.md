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

## Mapped Fields

Fields from the unified telemetry model mapped to pCars 2 packet fields:

| pCars 2 Field | Source Field | Notes |
|---------------|-------------|-------|
| Acceleration [3] | `acceleration_x/y/z` | Local and world |
| Angular velocity [3] | `angular_velocity_x/y/z` | |
| Best/Last/Current lap time | `best_lap_time`, `last_lap_time`, `current_lap_time` | In timings packet |
| Boost | `boost` | |
| Brake (uint8 + float32) | `brake` | |
| Car flags | `is_race_on`, `hand_brake` | ENGINE_ACTIVE, HANDBRAKE |
| Clutch (uint8 + float32) | `clutch` | |
| Engine torque | `torque` | |
| Fuel capacity | `fuel_capacity` | |
| Fuel level | `fuel` | 0.0-1.0 |
| Gear | `gear` | |
| HandBrake | `hand_brake` | Car flags bit + float field |
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

- Packet version: 2
- Max participants: 64
- Telemetry packet: 559 bytes
- Game state packet: 16 bytes
- Timings packet: 1063 bytes
- 12-byte common header on all packets
- Sequence number support (even = valid frame)
- Telemetry sent every frame, game state and timings sent every second
