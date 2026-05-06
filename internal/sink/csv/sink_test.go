package csv

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New(config.SinkEntry{FilePath: "out.csv"})
	assert.Equal(t, "csv", s.Name())
}

func TestSink_Run(t *testing.T) {
	t.Run("writes header and rows", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "telemetry.csv")

		s := New(config.SinkEntry{FilePath: path})

		in := make(chan model.TelemetryFrame, 2)
		in <- model.TelemetryFrame{
			Timestamp:  time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
			SourceName: "forza_xbox",
			SourceType: "forza",
			IsRaceOn:   model.Ptr(true),
			EngineRPM:  model.Ptr(float32(5000)),
			Speed:      model.Ptr(float32(30.5)),
			Gear:       model.Ptr(int8(3)),
			TireTemp:   &[4]float32{85.0, 86.0, 82.0, 83.0},
		}
		in <- model.TelemetryFrame{
			Timestamp:  time.Date(2026, 5, 6, 12, 0, 1, 0, time.UTC),
			SourceName: "forza_xbox",
			SourceType: "forza",
			IsRaceOn:   model.Ptr(true),
			EngineRPM:  model.Ptr(float32(6000)),
			Speed:      model.Ptr(float32(40.0)),
			Gear:       model.Ptr(int8(4)),
		}
		close(in)

		err := s.Run(context.Background(), in)
		require.NoError(t, err)

		f, err := os.Open(path)
		require.NoError(t, err)
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		require.NoError(t, err)

		// Header + 2 data rows
		assert.Len(t, records, 3)

		// Verify header
		assert.Equal(t, "timestamp", records[0][0])
		assert.Equal(t, "source_name", records[0][1])

		rpmIdx := findColumnIndex("engine_rpm")
		gearIdx := findColumnIndex("gear")
		srcNameIdx := findColumnIndex("source_name")
		srcTypeIdx := findColumnIndex("source_type")
		raceOnIdx := findColumnIndex("is_race_on")

		// Verify first row
		assert.Equal(t, "forza_xbox", records[1][srcNameIdx])
		assert.Equal(t, "forza", records[1][srcTypeIdx])
		assert.Equal(t, "true", records[1][raceOnIdx])
		assert.Equal(t, "5000", records[1][rpmIdx])
		assert.Equal(t, "3", records[1][gearIdx])

		// Verify second row
		assert.Equal(t, "6000", records[2][rpmIdx])
		assert.Equal(t, "4", records[2][gearIdx])
	})

	t.Run("nil fields are empty strings", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "telemetry.csv")

		s := New(config.SinkEntry{FilePath: path})

		in := make(chan model.TelemetryFrame, 1)
		in <- model.TelemetryFrame{
			Timestamp:  time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
			SourceName: "forza_xbox",
			SourceType: "forza",
			EngineRPM:  model.Ptr(float32(5000)),
			// OilTemp, WaterTemp, FuelCapacity etc are nil
		}
		close(in)

		err := s.Run(context.Background(), in)
		require.NoError(t, err)

		f, err := os.Open(path)
		require.NoError(t, err)
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		require.NoError(t, err)

		assert.Len(t, records, 2)

		// All rows have same column count
		assert.Equal(t, len(records[0]), len(records[1]))

		// oil_temp column should be empty
		oilTempIdx := findColumnIndex("oil_temp")
		assert.Equal(t, "", records[1][oilTempIdx])
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "telemetry.csv")

		s := New(config.SinkEntry{FilePath: path})

		in := make(chan model.TelemetryFrame)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("sink did not stop after cancel")
		}
	})

	t.Run("stops on channel close", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "telemetry.csv")

		s := New(config.SinkEntry{FilePath: path})

		in := make(chan model.TelemetryFrame)
		close(in)

		err := s.Run(context.Background(), in)
		assert.NoError(t, err)
	})

	t.Run("default column count matches", func(t *testing.T) {
		frame := model.TelemetryFrame{
			Timestamp:  time.Now(),
			SourceName: "test",
			SourceType: "forza",
		}

		row := defaultFrameToRow(&frame)
		assert.Equal(t, len(defaultColumns), len(row))
	})

	t.Run("forza column count matches", func(t *testing.T) {
		frame := model.TelemetryFrame{
			Timestamp:  time.Now(),
			SourceName: "test",
			SourceType: "forza",
		}

		row := forzaFrameToRow(&frame)
		assert.Equal(t, len(forzaColumns), len(row))
	})

	t.Run("forza format header", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "forza.csv")

		s := New(config.SinkEntry{FilePath: path, CSVFormat: "forza"})

		in := make(chan model.TelemetryFrame, 1)
		in <- model.TelemetryFrame{
			Timestamp:  time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
			SourceName: "forza_xbox",
			SourceType: "forza",
			IsRaceOn:   model.Ptr(true),
			EngineRPM:  model.Ptr(float32(5000)),
			Gear:       model.Ptr(int8(3)),
		}
		close(in)

		err := s.Run(context.Background(), in)
		require.NoError(t, err)

		f, err := os.Open(path)
		require.NoError(t, err)
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		require.NoError(t, err)

		assert.Len(t, records, 2)
		assert.Equal(t, "IsRaceOn", records[0][0])
		assert.Equal(t, "CurrentEngineRpm", records[0][4])
		assert.Equal(t, "true", records[1][0])
		assert.Equal(t, "5000", records[1][4])
	})
}

func BenchmarkDefaultFrameToRow(b *testing.B) {
	frame := &model.TelemetryFrame{
		Timestamp:  time.Now(),
		SourceName: "forza_xbox",
		SourceType: "forza",
		IsRaceOn:   model.Ptr(true),
		EngineRPM:  model.Ptr(float32(5000)),
		Speed:      model.Ptr(float32(30.0)),
		Gear:       model.Ptr(int8(3)),
		TireTemp:   &[4]float32{85.0, 86.0, 82.0, 83.0},
	}

	b.ResetTimer()

	for b.Loop() {
		defaultFrameToRow(frame)
	}
}

func BenchmarkForzaFrameToRow(b *testing.B) {
	frame := &model.TelemetryFrame{
		Timestamp:  time.Now(),
		SourceName: "forza_xbox",
		SourceType: "forza",
		IsRaceOn:   model.Ptr(true),
		EngineRPM:  model.Ptr(float32(5000)),
		Speed:      model.Ptr(float32(30.0)),
		Gear:       model.Ptr(int8(3)),
		TireTemp:   &[4]float32{85.0, 86.0, 82.0, 83.0},
	}

	b.ResetTimer()

	for b.Loop() {
		forzaFrameToRow(frame)
	}
}

func findColumnIndex(name string) int {
	for i, col := range defaultColumns {
		if col == name {
			return i
		}
	}

	return -1
}
