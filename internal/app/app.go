package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/vitalvas/gokit/xconfig"
	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/pipeline"
	"github.com/vitalvas/racetelemetry/internal/sink/grafanalive"
	"github.com/vitalvas/racetelemetry/internal/sink/jsonstdout"
	"github.com/vitalvas/racetelemetry/internal/sink/pcars1shm"
	"github.com/vitalvas/racetelemetry/internal/sink/pcars1udp"
	"github.com/vitalvas/racetelemetry/internal/sink/pcars2shm"
	"github.com/vitalvas/racetelemetry/internal/sink/pcars2udp"
	sinkwire "github.com/vitalvas/racetelemetry/internal/sink/wire"
	"github.com/vitalvas/racetelemetry/internal/source/forza"
	"github.com/vitalvas/racetelemetry/internal/source/gt7"
	sourcewire "github.com/vitalvas/racetelemetry/internal/source/wire"
)

// Run loads config, builds the pipeline, and runs it until ctx is done.
func Run(ctx context.Context, configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	cfg.SetDefaults()

	if err := cfg.Validate(); err != nil {
		return err
	}

	setupLogger(cfg.Log)

	sources, err := buildSources(cfg.Sources)
	if err != nil {
		return err
	}

	sinks, err := buildSinks(cfg.Sinks)
	if err != nil {
		return err
	}

	pipe := pipeline.New(sources, sinks, slog.Default())

	slog.Info("starting racetelemetry pipeline")

	return pipe.Run(ctx)
}

func loadConfig(path string) (*config.Config, error) {
	var cfg config.Config

	if err := xconfig.Load(&cfg, xconfig.WithFiles(path)); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setupLogger(cfg config.LogConfig) {
	var level slog.Level

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func buildSources(entries map[string]config.SourceEntry) (map[string]pipeline.Source, error) {
	sources := make(map[string]pipeline.Source, len(entries))

	for name, entry := range entries {
		switch entry.Type {
		case "forza":
			sources[name] = forza.New(entry.ListenAddr)
		case "gt7":
			sources[name] = gt7.New(entry.ConsoleAddr, entry.ListenAddr)
		case "wire":
			sources[name] = sourcewire.New(entry.ListenAddr)
		default:
			return nil, fmt.Errorf("unknown source type %q for %q", entry.Type, name)
		}
	}

	return sources, nil
}

func buildSinks(entries map[string]config.SinkEntry) (map[string]pipeline.SinkRoute, error) {
	sinks := make(map[string]pipeline.SinkRoute, len(entries))

	for name, entry := range entries {
		var s pipeline.Sink

		switch entry.Type {
		case "pcars1_udp":
			s = pcars1udp.New(entry.TargetAddr)
		case "pcars1_shm":
			shm, err := pcars1shm.New()
			if err != nil {
				return nil, fmt.Errorf("sink %q: %w", name, err)
			}

			s = shm
		case "pcars2_udp":
			s = pcars2udp.New(entry.TargetAddr)
		case "pcars2_shm":
			shm, err := pcars2shm.New()
			if err != nil {
				return nil, fmt.Errorf("sink %q: %w", name, err)
			}

			s = shm
		case "json_stdout":
			s = jsonstdout.New()
		case "wire":
			s = sinkwire.New(entry.TargetAddr)
		case "grafana_live":
			s = grafanalive.New(entry.Endpoint, entry.APIKey)
		default:
			return nil, fmt.Errorf("unknown sink type %q for %q", entry.Type, name)
		}

		sinks[name] = pipeline.SinkRoute{
			Sink:   s,
			Inputs: entry.Inputs,
		}
	}

	return sinks, nil
}
