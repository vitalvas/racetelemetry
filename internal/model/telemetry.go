package model

import "time"

// TelemetryFrame represents a single frame of normalized telemetry data.
// Pointer fields indicate data that may or may not be available from the source.
// A nil pointer means the source does not provide this field.
type TelemetryFrame struct {
	Timestamp  time.Time `json:"timestamp"`
	SourceName string    `json:"source_name"`
	SourceType string    `json:"source_type"`

	// Session state
	IsRaceOn        *bool    `json:"is_race_on,omitempty"`
	CurrentRaceTime *float32 `json:"current_race_time,omitempty"` // seconds

	// Engine
	EngineRPM     *float32 `json:"engine_rpm,omitempty"`
	EngineMaxRPM  *float32 `json:"engine_max_rpm,omitempty"`
	EngineIdleRPM *float32 `json:"engine_idle_rpm,omitempty"`
	NumCylinders  *int32   `json:"num_cylinders,omitempty"`

	// Drivetrain
	Gear           *int8    `json:"gear,omitempty"`           // -1=R, 0=N, 1+=forward
	SuggestedGear  *int8    `json:"suggested_gear,omitempty"` // GT7 only
	NumGears       *int8    `json:"num_gears,omitempty"`
	Speed          *float32 `json:"speed,omitempty"`            // m/s
	Throttle       *float32 `json:"throttle,omitempty"`         // 0.0-1.0
	Brake          *float32 `json:"brake,omitempty"`            // 0.0-1.0
	Clutch         *float32 `json:"clutch,omitempty"`           // 0.0-1.0
	ClutchEngaged  *float32 `json:"clutch_engaged,omitempty"`   // GT7 only
	RPMAfterClutch *float32 `json:"rpm_after_clutch,omitempty"` // GT7 only
	HandBrake      *float32 `json:"hand_brake,omitempty"`       // 0.0-1.0
	Steer          *float32 `json:"steer,omitempty"`            // -1.0 to 1.0 or radians
	SteerVelocity  *float32 `json:"steer_velocity,omitempty"`   // rad/s, GT7 only
	DrivetrainType *int32   `json:"drivetrain_type,omitempty"`  // 0=FWD, 1=RWD, 2=AWD

	// Physics - car local space
	AccelerationX    *float32 `json:"acceleration_x,omitempty"` // m/s^2
	AccelerationY    *float32 `json:"acceleration_y,omitempty"`
	AccelerationZ    *float32 `json:"acceleration_z,omitempty"`
	VelocityX        *float32 `json:"velocity_x,omitempty"` // m/s
	VelocityY        *float32 `json:"velocity_y,omitempty"`
	VelocityZ        *float32 `json:"velocity_z,omitempty"`
	AngularVelocityX *float32 `json:"angular_velocity_x,omitempty"` // rad/s
	AngularVelocityY *float32 `json:"angular_velocity_y,omitempty"`
	AngularVelocityZ *float32 `json:"angular_velocity_z,omitempty"`

	// Orientation - global space
	Yaw     *float32 `json:"yaw,omitempty"`     // radians
	Pitch   *float32 `json:"pitch,omitempty"`   // radians
	Roll    *float32 `json:"roll,omitempty"`    // radians
	Heading *float32 `json:"heading,omitempty"` // 0.0=south, 1.0=north, GT7 only

	// Position - global space, meters
	PositionX *float32 `json:"position_x,omitempty"`
	PositionY *float32 `json:"position_y,omitempty"`
	PositionZ *float32 `json:"position_z,omitempty"`

	// Power
	Power  *float32 `json:"power,omitempty"`  // watts
	Torque *float32 `json:"torque,omitempty"` // Nm
	Boost  *float32 `json:"boost,omitempty"`

	// Fuel
	Fuel         *float32 `json:"fuel,omitempty"`          // 0.0-1.0 percentage
	FuelCapacity *float32 `json:"fuel_capacity,omitempty"` // liters

	// Temperatures - all Celsius
	OilTemp     *float32 `json:"oil_temp,omitempty"`
	OilPressure *float32 `json:"oil_pressure,omitempty"` // GT7 only
	WaterTemp   *float32 `json:"water_temp,omitempty"`
	RideHeight  *float32 `json:"ride_height,omitempty"` // GT7 only

	// Per-wheel data ordered: [FL, FR, RL, RR]
	TireTemp                   *[4]float32 `json:"tire_temp,omitempty"`                    // Celsius
	SuspensionTravel           *[4]float32 `json:"suspension_travel,omitempty"`            // meters
	NormalizedSuspensionTravel *[4]float32 `json:"normalized_suspension_travel,omitempty"` // 0.0-1.0, Forza only
	WheelSpeed                 *[4]float32 `json:"wheel_speed,omitempty"`                  // rad/s
	WheelRadius                *[4]float32 `json:"wheel_radius,omitempty"`                 // meters, GT7 only
	WheelOnRumbleStrip         *[4]bool    `json:"wheel_on_rumble_strip,omitempty"`        // Forza only
	WheelInPuddleDepth         *[4]float32 `json:"wheel_in_puddle_depth,omitempty"`        // 0.0-1.0, Forza only
	SurfaceRumble              *[4]float32 `json:"surface_rumble,omitempty"`               // Forza only
	SlipRatio                  *[4]float32 `json:"slip_ratio,omitempty"`
	SlipAngle                  *[4]float32 `json:"slip_angle,omitempty"`
	TireCombinedSlip           *[4]float32 `json:"tire_combined_slip,omitempty"` // Forza only

	// Road surface - GT7 only
	RoadPlaneX    *float32 `json:"road_plane_x,omitempty"`
	RoadPlaneY    *float32 `json:"road_plane_y,omitempty"`
	RoadPlaneZ    *float32 `json:"road_plane_z,omitempty"`
	RoadPlaneDist *float32 `json:"road_plane_dist,omitempty"`

	// Lap and race
	LapNumber      *int16   `json:"lap_number,omitempty"`
	TotalLaps      *int16   `json:"total_laps,omitempty"`
	RacePosition   *uint8   `json:"race_position,omitempty"`
	TotalPositions *int16   `json:"total_positions,omitempty"`  // GT7 only
	BestLapTime    *float32 `json:"best_lap_time,omitempty"`    // seconds
	LastLapTime    *float32 `json:"last_lap_time,omitempty"`    // seconds
	CurrentLapTime *float32 `json:"current_lap_time,omitempty"` // seconds
	LapDistance    *float32 `json:"lap_distance,omitempty"`     // meters, resets per lap

	// Car identity
	CarClass            *int32 `json:"car_class,omitempty"`
	CarIndex            *int32 `json:"car_index,omitempty"`
	CarPerformanceIndex *int32 `json:"car_performance_index,omitempty"` // Forza only, 100-999
	EstimatedTopSpeed   *int16 `json:"estimated_top_speed,omitempty"`   // GT7 only

	// Transmission
	GearRatios    *[8]float32 `json:"gear_ratios,omitempty"`
	TopSpeedRatio *float32    `json:"top_speed_ratio,omitempty"` // GT7 only

	// Driving aids
	NormalizedDrivingLine *int8 `json:"normalized_driving_line,omitempty"`  // Forza only
	NormalizedAIBrakeDiff *int8 `json:"normalized_ai_brake_diff,omitempty"` // Forza only
}

// Wheel index constants.
const (
	WheelFL = 0
	WheelFR = 1
	WheelRL = 2
	WheelRR = 3
)

// Ptr creates a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}
