package model

import "time"

// TelemetryFrame represents a single frame of normalized telemetry data.
// Pointer fields indicate data that may or may not be available from the source.
// A nil pointer means the source does not provide this field.
type TelemetryFrame struct {
	Timestamp  time.Time `json:"timestamp"`
	SourceName string    `json:"source_name"`
	SourceType string    `json:"source_type"`

	IsRaceOn *bool `json:"is_race_on,omitempty"`

	EngineRPM     *float32 `json:"engine_rpm,omitempty"`
	EngineMaxRPM  *float32 `json:"engine_max_rpm,omitempty"`
	EngineIdleRPM *float32 `json:"engine_idle_rpm,omitempty"`

	Gear      *int8    `json:"gear,omitempty"`
	NumGears  *int8    `json:"num_gears,omitempty"`
	Speed     *float32 `json:"speed,omitempty"`
	Throttle  *float32 `json:"throttle,omitempty"`
	Brake     *float32 `json:"brake,omitempty"`
	Clutch    *float32 `json:"clutch,omitempty"`
	HandBrake *float32 `json:"hand_brake,omitempty"`
	Steer     *float32 `json:"steer,omitempty"`

	AccelerationX    *float32 `json:"acceleration_x,omitempty"`
	AccelerationY    *float32 `json:"acceleration_y,omitempty"`
	AccelerationZ    *float32 `json:"acceleration_z,omitempty"`
	VelocityX        *float32 `json:"velocity_x,omitempty"`
	VelocityY        *float32 `json:"velocity_y,omitempty"`
	VelocityZ        *float32 `json:"velocity_z,omitempty"`
	AngularVelocityX *float32 `json:"angular_velocity_x,omitempty"`
	AngularVelocityY *float32 `json:"angular_velocity_y,omitempty"`
	AngularVelocityZ *float32 `json:"angular_velocity_z,omitempty"`

	Yaw   *float32 `json:"yaw,omitempty"`
	Pitch *float32 `json:"pitch,omitempty"`
	Roll  *float32 `json:"roll,omitempty"`

	PositionX *float32 `json:"position_x,omitempty"`
	PositionY *float32 `json:"position_y,omitempty"`
	PositionZ *float32 `json:"position_z,omitempty"`

	Power  *float32 `json:"power,omitempty"`
	Torque *float32 `json:"torque,omitempty"`
	Boost  *float32 `json:"boost,omitempty"`

	Fuel         *float32 `json:"fuel,omitempty"`
	FuelCapacity *float32 `json:"fuel_capacity,omitempty"`

	OilTemp   *float32 `json:"oil_temp,omitempty"`
	WaterTemp *float32 `json:"water_temp,omitempty"`

	TireTemp         *[4]float32 `json:"tire_temp,omitempty"`
	SuspensionTravel *[4]float32 `json:"suspension_travel,omitempty"`
	WheelSpeed       *[4]float32 `json:"wheel_speed,omitempty"`
	SlipRatio        *[4]float32 `json:"slip_ratio,omitempty"`
	SlipAngle        *[4]float32 `json:"slip_angle,omitempty"`

	LapNumber      *int16   `json:"lap_number,omitempty"`
	TotalLaps      *int16   `json:"total_laps,omitempty"`
	RacePosition   *uint8   `json:"race_position,omitempty"`
	BestLapTime    *float32 `json:"best_lap_time,omitempty"`
	LastLapTime    *float32 `json:"last_lap_time,omitempty"`
	CurrentLapTime *float32 `json:"current_lap_time,omitempty"`
	LapDistance    *float32 `json:"lap_distance,omitempty"`

	CarClass *int32 `json:"car_class,omitempty"`
	CarIndex *int32 `json:"car_index,omitempty"`

	GearRatios *[8]float32 `json:"gear_ratios,omitempty"`
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
