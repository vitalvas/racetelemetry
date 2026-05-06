package grafanalive

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

// Sink pushes telemetry frames to Grafana Live using the HTTP push API
// with Influx line protocol format.
type Sink struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

// New creates a new Grafana Live sink.
func New(entry config.SinkEntry) *Sink {
	return &Sink{
		endpoint: entry.Endpoint,
		apiKey:   entry.APIKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "grafana_live"
}

// Run consumes telemetry frames and pushes them to Grafana Live.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case frame, ok := <-in:
			if !ok {
				return nil
			}

			if err := s.push(ctx, &frame); err != nil {
				slog.Debug("grafana_live: push error", "error", err)
			}
		}
	}
}

func (s *Sink) push(ctx context.Context, frame *model.TelemetryFrame) error {
	line := formatLineProtocol(frame)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, strings.NewReader(line))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func val[T any](p *T) T {
	if p != nil {
		return *p
	}

	var zero T

	return zero
}

func valArr4[T any](p *[4]T) [4]T {
	if p != nil {
		return *p
	}

	var zero [4]T

	return zero
}

func formatLineProtocol(frame *model.TelemetryFrame) string {
	ts := frame.Timestamp.UnixNano()
	if ts <= 0 {
		ts = time.Now().UnixNano()
	}

	var b strings.Builder

	writeTelemetryLine(&b, frame, ts)
	writeTireLines(&b, frame, ts)

	return b.String()
}

func writeTelemetryLine(b *strings.Builder, frame *model.TelemetryFrame, ts int64) {
	fmt.Fprintf(b, "telemetry,source_name=%s,source_type=%s ", frame.SourceName, frame.SourceType)
	fmt.Fprintf(b, "is_race_on=%t,", val(frame.IsRaceOn))
	fmt.Fprintf(b, "engine_rpm=%.2f,", val(frame.EngineRPM))
	fmt.Fprintf(b, "engine_max_rpm=%.2f,", val(frame.EngineMaxRPM))
	fmt.Fprintf(b, "engine_idle_rpm=%.2f,", val(frame.EngineIdleRPM))
	fmt.Fprintf(b, "speed=%.4f,", val(frame.Speed))
	fmt.Fprintf(b, "throttle=%.4f,", val(frame.Throttle))
	fmt.Fprintf(b, "brake=%.4f,", val(frame.Brake))
	fmt.Fprintf(b, "clutch=%.4f,", val(frame.Clutch))
	fmt.Fprintf(b, "hand_brake=%.4f,", val(frame.HandBrake))
	fmt.Fprintf(b, "steer=%.4f,", val(frame.Steer))
	fmt.Fprintf(b, "gear=%di,", val(frame.Gear))
	fmt.Fprintf(b, "suggested_gear=%di,", val(frame.SuggestedGear))
	fmt.Fprintf(b, "power=%.2f,", val(frame.Power))
	fmt.Fprintf(b, "torque=%.2f,", val(frame.Torque))
	fmt.Fprintf(b, "boost=%.4f,", val(frame.Boost))
	fmt.Fprintf(b, "fuel=%.4f,", val(frame.Fuel))
	fmt.Fprintf(b, "fuel_capacity=%.2f,", val(frame.FuelCapacity))
	fmt.Fprintf(b, "position_x=%.4f,", val(frame.PositionX))
	fmt.Fprintf(b, "position_y=%.4f,", val(frame.PositionY))
	fmt.Fprintf(b, "position_z=%.4f,", val(frame.PositionZ))
	fmt.Fprintf(b, "velocity_x=%.4f,", val(frame.VelocityX))
	fmt.Fprintf(b, "velocity_y=%.4f,", val(frame.VelocityY))
	fmt.Fprintf(b, "velocity_z=%.4f,", val(frame.VelocityZ))
	fmt.Fprintf(b, "acceleration_x=%.4f,", val(frame.AccelerationX))
	fmt.Fprintf(b, "acceleration_y=%.4f,", val(frame.AccelerationY))
	fmt.Fprintf(b, "acceleration_z=%.4f,", val(frame.AccelerationZ))
	fmt.Fprintf(b, "yaw=%.6f,", val(frame.Yaw))
	fmt.Fprintf(b, "pitch=%.6f,", val(frame.Pitch))
	fmt.Fprintf(b, "roll=%.6f,", val(frame.Roll))
	fmt.Fprintf(b, "heading=%.6f,", val(frame.Heading))
	fmt.Fprintf(b, "ride_height=%.4f,", val(frame.RideHeight))
	fmt.Fprintf(b, "lap_number=%di,", val(frame.LapNumber))
	fmt.Fprintf(b, "total_laps=%di,", val(frame.TotalLaps))
	fmt.Fprintf(b, "race_position=%di,", val(frame.RacePosition))
	fmt.Fprintf(b, "total_positions=%di,", val(frame.TotalPositions))
	fmt.Fprintf(b, "best_lap_time=%.3f,", val(frame.BestLapTime))
	fmt.Fprintf(b, "last_lap_time=%.3f,", val(frame.LastLapTime))
	fmt.Fprintf(b, "current_lap_time=%.3f,", val(frame.CurrentLapTime))
	fmt.Fprintf(b, "current_race_time=%.3f,", val(frame.CurrentRaceTime))
	fmt.Fprintf(b, "lap_distance=%.2f,", val(frame.LapDistance))
	fmt.Fprintf(b, "oil_temp=%.1f,", val(frame.OilTemp))
	fmt.Fprintf(b, "oil_pressure=%.1f,", val(frame.OilPressure))
	fmt.Fprintf(b, "water_temp=%.1f,", val(frame.WaterTemp))
	fmt.Fprintf(b, "drivetrain_type=%di,", val(frame.DrivetrainType))
	fmt.Fprintf(b, "num_cylinders=%di,", val(frame.NumCylinders))
	fmt.Fprintf(b, "car_class=%di,", val(frame.CarClass))
	fmt.Fprintf(b, "car_index=%di,", val(frame.CarIndex))
	fmt.Fprintf(b, "car_performance_index=%di", val(frame.CarPerformanceIndex))
	fmt.Fprintf(b, " %d\n", ts)
}

func writeTireLines(b *strings.Builder, frame *model.TelemetryFrame, ts int64) {
	wheelNames := [4]string{"fl", "fr", "rl", "rr"}

	tireTemp := valArr4(frame.TireTemp)
	suspTravel := valArr4(frame.SuspensionTravel)
	normSuspTravel := valArr4(frame.NormalizedSuspensionTravel)
	wheelSpeed := valArr4(frame.WheelSpeed)
	wheelRadius := valArr4(frame.WheelRadius)
	slipRatio := valArr4(frame.SlipRatio)
	slipAngle := valArr4(frame.SlipAngle)
	combinedSlip := valArr4(frame.TireCombinedSlip)
	puddleDepth := valArr4(frame.WheelInPuddleDepth)
	surfaceRumble := valArr4(frame.SurfaceRumble)

	for i, name := range wheelNames {
		fmt.Fprintf(b, "tire,source_name=%s,source_type=%s,wheel=%s ", frame.SourceName, frame.SourceType, name)
		fmt.Fprintf(b, "temp=%.1f,", tireTemp[i])
		fmt.Fprintf(b, "suspension=%.4f,", suspTravel[i])
		fmt.Fprintf(b, "normalized_suspension=%.4f,", normSuspTravel[i])
		fmt.Fprintf(b, "wheel_speed=%.4f,", wheelSpeed[i])
		fmt.Fprintf(b, "wheel_radius=%.4f,", wheelRadius[i])
		fmt.Fprintf(b, "slip_ratio=%.4f,", slipRatio[i])
		fmt.Fprintf(b, "slip_angle=%.4f,", slipAngle[i])
		fmt.Fprintf(b, "combined_slip=%.4f,", combinedSlip[i])
		fmt.Fprintf(b, "puddle_depth=%.4f,", puddleDepth[i])
		fmt.Fprintf(b, "surface_rumble=%.4f", surfaceRumble[i])
		fmt.Fprintf(b, " %d\n", ts)
	}
}
