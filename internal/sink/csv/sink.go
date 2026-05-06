package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

// Column name sets for different CSV formats.
var defaultColumns = []string{
	"timestamp", "source_name", "source_type",
	"is_race_on", "current_race_time",
	"engine_rpm", "engine_max_rpm", "engine_idle_rpm", "num_cylinders",
	"gear", "suggested_gear", "num_gears", "speed",
	"throttle", "brake", "clutch", "clutch_engaged", "rpm_after_clutch",
	"hand_brake", "steer", "steer_velocity", "drivetrain_type",
	"acceleration_x", "acceleration_y", "acceleration_z",
	"velocity_x", "velocity_y", "velocity_z",
	"angular_velocity_x", "angular_velocity_y", "angular_velocity_z",
	"yaw", "pitch", "roll", "heading",
	"position_x", "position_y", "position_z",
	"power", "torque", "boost",
	"fuel", "fuel_capacity",
	"oil_temp", "oil_pressure", "water_temp", "ride_height",
	"tire_temp_fl", "tire_temp_fr", "tire_temp_rl", "tire_temp_rr",
	"suspension_travel_fl", "suspension_travel_fr", "suspension_travel_rl", "suspension_travel_rr",
	"normalized_suspension_fl", "normalized_suspension_fr", "normalized_suspension_rl", "normalized_suspension_rr",
	"wheel_speed_fl", "wheel_speed_fr", "wheel_speed_rl", "wheel_speed_rr",
	"wheel_radius_fl", "wheel_radius_fr", "wheel_radius_rl", "wheel_radius_rr",
	"slip_ratio_fl", "slip_ratio_fr", "slip_ratio_rl", "slip_ratio_rr",
	"slip_angle_fl", "slip_angle_fr", "slip_angle_rl", "slip_angle_rr",
	"tire_combined_slip_fl", "tire_combined_slip_fr", "tire_combined_slip_rl", "tire_combined_slip_rr",
	"wheel_in_puddle_fl", "wheel_in_puddle_fr", "wheel_in_puddle_rl", "wheel_in_puddle_rr",
	"surface_rumble_fl", "surface_rumble_fr", "surface_rumble_rl", "surface_rumble_rr",
	"wheel_on_rumble_fl", "wheel_on_rumble_fr", "wheel_on_rumble_rl", "wheel_on_rumble_rr",
	"road_plane_x", "road_plane_y", "road_plane_z", "road_plane_dist",
	"lap_number", "total_laps", "race_position", "total_positions",
	"best_lap_time", "last_lap_time", "current_lap_time", "lap_distance",
	"car_class", "car_index", "car_performance_index", "estimated_top_speed",
	"gear_ratio_1", "gear_ratio_2", "gear_ratio_3", "gear_ratio_4",
	"gear_ratio_5", "gear_ratio_6", "gear_ratio_7", "gear_ratio_8",
	"top_speed_ratio",
	"normalized_driving_line", "normalized_ai_brake_diff",
}

// Forza-compatible column names matching the official Forza packet spec.
// Uses the same column names as richstokes/Forza-data-tools.
var forzaColumns = []string{
	"IsRaceOn", "TimestampMS",
	"EngineMaxRpm", "EngineIdleRpm", "CurrentEngineRpm",
	"AccelerationX", "AccelerationY", "AccelerationZ",
	"VelocityX", "VelocityY", "VelocityZ",
	"AngularVelocityX", "AngularVelocityY", "AngularVelocityZ",
	"Yaw", "Pitch", "Roll",
	"NormalizedSuspensionTravelFrontLeft", "NormalizedSuspensionTravelFrontRight",
	"NormalizedSuspensionTravelRearLeft", "NormalizedSuspensionTravelRearRight",
	"TireSlipRatioFrontLeft", "TireSlipRatioFrontRight",
	"TireSlipRatioRearLeft", "TireSlipRatioRearRight",
	"WheelRotationSpeedFrontLeft", "WheelRotationSpeedFrontRight",
	"WheelRotationSpeedRearLeft", "WheelRotationSpeedRearRight",
	"WheelOnRumbleStripFrontLeft", "WheelOnRumbleStripFrontRight",
	"WheelOnRumbleStripRearLeft", "WheelOnRumbleStripRearRight",
	"WheelInPuddleDepthFrontLeft", "WheelInPuddleDepthFrontRight",
	"WheelInPuddleDepthRearLeft", "WheelInPuddleDepthRearRight",
	"SurfaceRumbleFrontLeft", "SurfaceRumbleFrontRight",
	"SurfaceRumbleRearLeft", "SurfaceRumbleRearRight",
	"TireSlipAngleFrontLeft", "TireSlipAngleFrontRight",
	"TireSlipAngleRearLeft", "TireSlipAngleRearRight",
	"TireCombinedSlipFrontLeft", "TireCombinedSlipFrontRight",
	"TireCombinedSlipRearLeft", "TireCombinedSlipRearRight",
	"SuspensionTravelMetersFrontLeft", "SuspensionTravelMetersFrontRight",
	"SuspensionTravelMetersRearLeft", "SuspensionTravelMetersRearRight",
	"CarOrdinal", "CarClass", "CarPerformanceIndex", "DrivetrainType", "NumCylinders",
	"PositionX", "PositionY", "PositionZ",
	"Speed", "Power", "Torque",
	"TireTempFrontLeft", "TireTempFrontRight", "TireTempRearLeft", "TireTempRearRight",
	"Boost", "Fuel", "DistanceTraveled",
	"BestLap", "LastLap", "CurrentLap", "CurrentRaceTime",
	"LapNumber", "RacePosition",
	"Accel", "Brake", "Clutch", "HandBrake", "Gear", "Steer",
	"NormalizedDrivingLine", "NormalizedAIBrakeDifference",
}

// Sink writes telemetry frames as CSV to a file.
type Sink struct {
	filePath string
	format   string
}

// New creates a new CSV file sink.
func New(entry config.SinkEntry) *Sink {
	format := entry.CSVFormat
	if format == "" {
		format = "default"
	}

	return &Sink{
		filePath: entry.FilePath,
		format:   format,
	}
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "csv"
}

// Run consumes telemetry frames and writes them as CSV rows.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	f, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("csv sink: create file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := s.getColumns()

	if err := w.Write(header); err != nil {
		return fmt.Errorf("csv sink: write header: %w", err)
	}

	slog.Info("csv sink writing", "file", s.filePath, "format", s.format)

	writeRow := s.getRowFunc()

	for {
		select {
		case <-ctx.Done():
			return nil
		case frame, ok := <-in:
			if !ok {
				return nil
			}

			if err := w.Write(writeRow(&frame)); err != nil {
				return fmt.Errorf("csv sink: write row: %w", err)
			}

			w.Flush()
		}
	}
}

func (s *Sink) getColumns() []string {
	if s.format == "forza" {
		return forzaColumns
	}

	return defaultColumns
}

func (s *Sink) getRowFunc() func(*model.TelemetryFrame) []string {
	if s.format == "forza" {
		return forzaFrameToRow
	}

	return defaultFrameToRow
}

func defaultFrameToRow(f *model.TelemetryFrame) []string {
	row := make([]string, 0, len(defaultColumns))

	row = append(row, f.Timestamp.Format("2006-01-02T15:04:05.000000Z07:00"))
	row = append(row, f.SourceName)
	row = append(row, f.SourceType)

	row = append(row, fmtBool(f.IsRaceOn))
	row = append(row, fmtFloat(f.CurrentRaceTime))

	row = append(row, fmtFloat(f.EngineRPM))
	row = append(row, fmtFloat(f.EngineMaxRPM))
	row = append(row, fmtFloat(f.EngineIdleRPM))
	row = append(row, fmtInt32(f.NumCylinders))

	row = append(row, fmtInt8(f.Gear))
	row = append(row, fmtInt8(f.SuggestedGear))
	row = append(row, fmtInt8(f.NumGears))
	row = append(row, fmtFloat(f.Speed))

	row = append(row, fmtFloat(f.Throttle))
	row = append(row, fmtFloat(f.Brake))
	row = append(row, fmtFloat(f.Clutch))
	row = append(row, fmtFloat(f.ClutchEngaged))
	row = append(row, fmtFloat(f.RPMAfterClutch))

	row = append(row, fmtFloat(f.HandBrake))
	row = append(row, fmtFloat(f.Steer))
	row = append(row, fmtFloat(f.SteerVelocity))
	row = append(row, fmtInt32(f.DrivetrainType))

	row = append(row, fmtFloat(f.AccelerationX))
	row = append(row, fmtFloat(f.AccelerationY))
	row = append(row, fmtFloat(f.AccelerationZ))

	row = append(row, fmtFloat(f.VelocityX))
	row = append(row, fmtFloat(f.VelocityY))
	row = append(row, fmtFloat(f.VelocityZ))

	row = append(row, fmtFloat(f.AngularVelocityX))
	row = append(row, fmtFloat(f.AngularVelocityY))
	row = append(row, fmtFloat(f.AngularVelocityZ))

	row = append(row, fmtFloat(f.Yaw))
	row = append(row, fmtFloat(f.Pitch))
	row = append(row, fmtFloat(f.Roll))
	row = append(row, fmtFloat(f.Heading))

	row = append(row, fmtFloat(f.PositionX))
	row = append(row, fmtFloat(f.PositionY))
	row = append(row, fmtFloat(f.PositionZ))

	row = append(row, fmtFloat(f.Power))
	row = append(row, fmtFloat(f.Torque))
	row = append(row, fmtFloat(f.Boost))

	row = append(row, fmtFloat(f.Fuel))
	row = append(row, fmtFloat(f.FuelCapacity))

	row = append(row, fmtFloat(f.OilTemp))
	row = append(row, fmtFloat(f.OilPressure))
	row = append(row, fmtFloat(f.WaterTemp))
	row = append(row, fmtFloat(f.RideHeight))

	row = append(row, fmtFloatArr4(f.TireTemp)...)
	row = append(row, fmtFloatArr4(f.SuspensionTravel)...)
	row = append(row, fmtFloatArr4(f.NormalizedSuspensionTravel)...)
	row = append(row, fmtFloatArr4(f.WheelSpeed)...)
	row = append(row, fmtFloatArr4(f.WheelRadius)...)
	row = append(row, fmtFloatArr4(f.SlipRatio)...)
	row = append(row, fmtFloatArr4(f.SlipAngle)...)
	row = append(row, fmtFloatArr4(f.TireCombinedSlip)...)
	row = append(row, fmtFloatArr4(f.WheelInPuddleDepth)...)
	row = append(row, fmtFloatArr4(f.SurfaceRumble)...)
	row = append(row, fmtBoolArr4(f.WheelOnRumbleStrip)...)

	row = append(row, fmtFloat(f.RoadPlaneX))
	row = append(row, fmtFloat(f.RoadPlaneY))
	row = append(row, fmtFloat(f.RoadPlaneZ))
	row = append(row, fmtFloat(f.RoadPlaneDist))

	row = append(row, fmtInt16(f.LapNumber))
	row = append(row, fmtInt16(f.TotalLaps))
	row = append(row, fmtUint8(f.RacePosition))
	row = append(row, fmtInt16(f.TotalPositions))

	row = append(row, fmtFloat(f.BestLapTime))
	row = append(row, fmtFloat(f.LastLapTime))
	row = append(row, fmtFloat(f.CurrentLapTime))
	row = append(row, fmtFloat(f.LapDistance))

	row = append(row, fmtInt32(f.CarClass))
	row = append(row, fmtInt32(f.CarIndex))
	row = append(row, fmtInt32(f.CarPerformanceIndex))
	row = append(row, fmtInt16(f.EstimatedTopSpeed))

	row = append(row, fmtFloatArr8(f.GearRatios)...)

	row = append(row, fmtFloat(f.TopSpeedRatio))

	row = append(row, fmtInt8(f.NormalizedDrivingLine))
	row = append(row, fmtInt8(f.NormalizedAIBrakeDiff))

	return row
}

// forzaFrameToRow outputs columns in official Forza packet field order,
// compatible with richstokes/Forza-data-tools and austinbaccus/forza-telemetry.
func forzaFrameToRow(f *model.TelemetryFrame) []string {
	row := make([]string, 0, len(forzaColumns))

	row = append(row, fmtBool(f.IsRaceOn))
	row = append(row, f.Timestamp.Format("2006-01-02T15:04:05.000000Z07:00"))

	row = append(row, fmtFloat(f.EngineMaxRPM))
	row = append(row, fmtFloat(f.EngineIdleRPM))
	row = append(row, fmtFloat(f.EngineRPM))

	row = append(row, fmtFloat(f.AccelerationX))
	row = append(row, fmtFloat(f.AccelerationY))
	row = append(row, fmtFloat(f.AccelerationZ))

	row = append(row, fmtFloat(f.VelocityX))
	row = append(row, fmtFloat(f.VelocityY))
	row = append(row, fmtFloat(f.VelocityZ))

	row = append(row, fmtFloat(f.AngularVelocityX))
	row = append(row, fmtFloat(f.AngularVelocityY))
	row = append(row, fmtFloat(f.AngularVelocityZ))

	row = append(row, fmtFloat(f.Yaw))
	row = append(row, fmtFloat(f.Pitch))
	row = append(row, fmtFloat(f.Roll))

	row = append(row, fmtFloatArr4(f.NormalizedSuspensionTravel)...)
	row = append(row, fmtFloatArr4(f.SlipRatio)...)
	row = append(row, fmtFloatArr4(f.WheelSpeed)...)
	row = append(row, fmtBoolArr4(f.WheelOnRumbleStrip)...)
	row = append(row, fmtFloatArr4(f.WheelInPuddleDepth)...)
	row = append(row, fmtFloatArr4(f.SurfaceRumble)...)
	row = append(row, fmtFloatArr4(f.SlipAngle)...)
	row = append(row, fmtFloatArr4(f.TireCombinedSlip)...)
	row = append(row, fmtFloatArr4(f.SuspensionTravel)...)

	row = append(row, fmtInt32(f.CarIndex))
	row = append(row, fmtInt32(f.CarClass))
	row = append(row, fmtInt32(f.CarPerformanceIndex))
	row = append(row, fmtInt32(f.DrivetrainType))
	row = append(row, fmtInt32(f.NumCylinders))

	row = append(row, fmtFloat(f.PositionX))
	row = append(row, fmtFloat(f.PositionY))
	row = append(row, fmtFloat(f.PositionZ))

	row = append(row, fmtFloat(f.Speed))
	row = append(row, fmtFloat(f.Power))
	row = append(row, fmtFloat(f.Torque))

	row = append(row, fmtFloatArr4(f.TireTemp)...)

	row = append(row, fmtFloat(f.Boost))
	row = append(row, fmtFloat(f.Fuel))
	row = append(row, fmtFloat(f.LapDistance))

	row = append(row, fmtFloat(f.BestLapTime))
	row = append(row, fmtFloat(f.LastLapTime))
	row = append(row, fmtFloat(f.CurrentLapTime))
	row = append(row, fmtFloat(f.CurrentRaceTime))

	row = append(row, fmtInt16(f.LapNumber))
	row = append(row, fmtUint8(f.RacePosition))

	row = append(row, fmtFloat(f.Throttle))
	row = append(row, fmtFloat(f.Brake))
	row = append(row, fmtFloat(f.Clutch))
	row = append(row, fmtFloat(f.HandBrake))
	row = append(row, fmtInt8(f.Gear))
	row = append(row, fmtFloat(f.Steer))

	row = append(row, fmtInt8(f.NormalizedDrivingLine))
	row = append(row, fmtInt8(f.NormalizedAIBrakeDiff))

	return row
}

func fmtFloat(p *float32) string {
	if p == nil {
		return ""
	}

	return strconv.FormatFloat(float64(*p), 'f', -1, 32)
}

func fmtBool(p *bool) string {
	if p == nil {
		return ""
	}

	return strconv.FormatBool(*p)
}

func fmtInt8(p *int8) string {
	if p == nil {
		return ""
	}

	return strconv.FormatInt(int64(*p), 10)
}

func fmtInt16(p *int16) string {
	if p == nil {
		return ""
	}

	return strconv.FormatInt(int64(*p), 10)
}

func fmtInt32(p *int32) string {
	if p == nil {
		return ""
	}

	return strconv.FormatInt(int64(*p), 10)
}

func fmtUint8(p *uint8) string {
	if p == nil {
		return ""
	}

	return strconv.FormatUint(uint64(*p), 10)
}

func fmtFloatArr4(p *[4]float32) []string {
	if p == nil {
		return []string{"", "", "", ""}
	}

	return []string{
		strconv.FormatFloat(float64(p[0]), 'f', -1, 32),
		strconv.FormatFloat(float64(p[1]), 'f', -1, 32),
		strconv.FormatFloat(float64(p[2]), 'f', -1, 32),
		strconv.FormatFloat(float64(p[3]), 'f', -1, 32),
	}
}

func fmtFloatArr8(p *[8]float32) []string {
	if p == nil {
		return []string{"", "", "", "", "", "", "", ""}
	}

	out := make([]string, 8)

	for i := range 8 {
		out[i] = strconv.FormatFloat(float64(p[i]), 'f', -1, 32)
	}

	return out
}

func fmtBoolArr4(p *[4]bool) []string {
	if p == nil {
		return []string{"", "", "", ""}
	}

	return []string{
		strconv.FormatBool(p[0]),
		strconv.FormatBool(p[1]),
		strconv.FormatBool(p[2]),
		strconv.FormatBool(p[3]),
	}
}
