package model

import "time"

// TelemetryFrame represents a single frame of normalized telemetry data.
// This is the superset of fields from all supported racing game sources,
// normalized to common units and conventions.
type TelemetryFrame struct {
	Timestamp time.Time

	// Session state
	IsRaceOn bool

	// Engine
	EngineRPM     float32
	EngineMaxRPM  float32
	EngineIdleRPM float32

	// Drivetrain
	// Gear: -1=Reverse, 0=Neutral, 1+=forward gears
	Gear     int8
	NumGears int8
	// Speed in meters per second
	Speed float32
	// Throttle 0.0-1.0
	Throttle float32
	// Brake 0.0-1.0
	Brake float32
	// Clutch 0.0-1.0
	Clutch float32
	// HandBrake 0.0-1.0
	HandBrake float32
	// Steer -1.0 (left) to 1.0 (right)
	Steer float32

	// Physics - car local space
	AccelerationX    float32 // m/s^2
	AccelerationY    float32
	AccelerationZ    float32
	VelocityX        float32 // m/s
	VelocityY        float32
	VelocityZ        float32
	AngularVelocityX float32 // rad/s
	AngularVelocityY float32
	AngularVelocityZ float32

	// Orientation - global space, radians
	Yaw   float32
	Pitch float32
	Roll  float32

	// Position - global space, meters
	PositionX float32
	PositionY float32
	PositionZ float32

	// Power
	Power  float32 // watts
	Torque float32 // Nm
	Boost  float32

	// Fuel - 0.0 to 1.0 (percentage of capacity)
	Fuel         float32
	FuelCapacity float32

	// Temperatures - all Celsius
	OilTemp   float32
	WaterTemp float32

	// Per-wheel data ordered: [FL, FR, RL, RR]
	TireTemp         [4]float32 // Celsius
	SuspensionTravel [4]float32 // meters
	WheelSpeed       [4]float32 // rad/s
	SlipRatio        [4]float32
	SlipAngle        [4]float32

	// Lap and race
	LapNumber      int16
	TotalLaps      int16
	RacePosition   uint8
	BestLapTime    float32 // seconds
	LastLapTime    float32 // seconds
	CurrentLapTime float32 // seconds
	LapDistance    float32 // meters

	// Car identity
	CarClass int32
	CarIndex int32

	// Gear ratios (where available)
	GearRatios [8]float32
}

// Wheel index constants.
const (
	WheelFL = 0
	WheelFR = 1
	WheelRL = 2
	WheelRR = 3
)
