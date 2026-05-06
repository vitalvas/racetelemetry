package grafanalive

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
func New(endpoint, apiKey string) *Sink {
	return &Sink{
		endpoint: endpoint,
		apiKey:   apiKey,
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
	b.WriteString("telemetry ")
	fmt.Fprintf(b, "is_race_on=%t,", frame.IsRaceOn)
	fmt.Fprintf(b, "engine_rpm=%.2f,", frame.EngineRPM)
	fmt.Fprintf(b, "engine_max_rpm=%.2f,", frame.EngineMaxRPM)
	fmt.Fprintf(b, "speed=%.4f,", frame.Speed)
	fmt.Fprintf(b, "throttle=%.4f,", frame.Throttle)
	fmt.Fprintf(b, "brake=%.4f,", frame.Brake)
	fmt.Fprintf(b, "clutch=%.4f,", frame.Clutch)
	fmt.Fprintf(b, "steer=%.4f,", frame.Steer)
	fmt.Fprintf(b, "gear=%di,", frame.Gear)
	fmt.Fprintf(b, "boost=%.4f,", frame.Boost)
	fmt.Fprintf(b, "fuel=%.4f,", frame.Fuel)
	fmt.Fprintf(b, "position_x=%.4f,", frame.PositionX)
	fmt.Fprintf(b, "position_y=%.4f,", frame.PositionY)
	fmt.Fprintf(b, "position_z=%.4f,", frame.PositionZ)
	fmt.Fprintf(b, "yaw=%.6f,", frame.Yaw)
	fmt.Fprintf(b, "pitch=%.6f,", frame.Pitch)
	fmt.Fprintf(b, "roll=%.6f,", frame.Roll)
	fmt.Fprintf(b, "lap_number=%di,", frame.LapNumber)
	fmt.Fprintf(b, "race_position=%di,", frame.RacePosition)
	fmt.Fprintf(b, "best_lap_time=%.3f,", frame.BestLapTime)
	fmt.Fprintf(b, "last_lap_time=%.3f,", frame.LastLapTime)
	fmt.Fprintf(b, "current_lap_time=%.3f,", frame.CurrentLapTime)
	fmt.Fprintf(b, "oil_temp=%.1f,", frame.OilTemp)
	fmt.Fprintf(b, "water_temp=%.1f", frame.WaterTemp)
	fmt.Fprintf(b, " %d\n", ts)
}

func writeTireLines(b *strings.Builder, frame *model.TelemetryFrame, ts int64) {
	wheelNames := [4]string{"fl", "fr", "rl", "rr"}

	for i, name := range wheelNames {
		fmt.Fprintf(b, "tire,wheel=%s ", name)
		fmt.Fprintf(b, "temp=%.1f,", frame.TireTemp[i])
		fmt.Fprintf(b, "suspension=%.4f,", frame.SuspensionTravel[i])
		fmt.Fprintf(b, "wheel_speed=%.4f,", frame.WheelSpeed[i])
		fmt.Fprintf(b, "slip_ratio=%.4f,", frame.SlipRatio[i])
		fmt.Fprintf(b, "slip_angle=%.4f", frame.SlipAngle[i])
		fmt.Fprintf(b, " %d\n", ts)
	}
}
