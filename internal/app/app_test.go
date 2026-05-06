package app

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/config"
)

func TestRun(t *testing.T) {
	t.Run("missing config file", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := Run(ctx, "/nonexistent/config.yaml")
		assert.Error(t, err)
	})

	t.Run("invalid config no sources", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		err := os.WriteFile(cfgPath, []byte("log:\n  level: info\n  format: text\n"), 0o644)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = Run(ctx, cfgPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no sources configured")
	})

	t.Run("valid config runs and stops", func(t *testing.T) {
		skipIfNoNetwork(t)

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		cfgContent := `
sources:
  my_forza:
    type: forza
    listen_addr: "127.0.0.1:15300"
sinks:
  my_sink:
    type: pcars2_udp
    inputs: [my_forza]
    target_addr: "127.0.0.1:15606"
`
		err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		err = Run(ctx, cfgPath)
		assert.NoError(t, err)
	})

	t.Run("invalid config validation fails", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		cfgContent := `
log:
  level: invalid_level
  format: text
sources:
  s1:
    type: forza
    listen_addr: ":5300"
sinks:
  k1:
    type: pcars2_udp
    inputs: [s1]
    target_addr: "127.0.0.1:5606"
`
		err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = Run(ctx, cfgPath)
		assert.Error(t, err)
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("valid yaml", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		cfgContent := `
sources:
  my_forza:
    type: forza
    listen_addr: ":5300"
sinks:
  my_sink:
    type: pcars2_udp
    inputs: [my_forza]
    target_addr: "127.0.0.1:5606"
`
		err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644)
		require.NoError(t, err)

		cfg, err := loadConfig(cfgPath)
		require.NoError(t, err)

		assert.Len(t, cfg.Sources, 1)
		assert.Equal(t, "forza", cfg.Sources["my_forza"].Type)
		assert.Equal(t, ":5300", cfg.Sources["my_forza"].ListenAddr)
	})

	t.Run("valid json", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.json")

		cfgContent := `{"sources":{"s1":{"type":"forza","listen_addr":":5300"}},"sinks":{"k1":{"type":"pcars2_udp","inputs":["s1"],"target_addr":"127.0.0.1:5606"}}}`

		err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644)
		require.NoError(t, err)

		cfg, err := loadConfig(cfgPath)
		require.NoError(t, err)
		assert.Len(t, cfg.Sources, 1)
	})

	t.Run("invalid yaml content", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		err := os.WriteFile(cfgPath, []byte("{{invalid yaml"), 0o644)
		require.NoError(t, err)

		_, err = loadConfig(cfgPath)
		assert.Error(t, err)
	})
}

func TestBuildSources(t *testing.T) {
	t.Run("forza source", func(t *testing.T) {
		entries := map[string]config.SourceEntry{
			"my_forza": {Type: "forza", ListenAddr: ":5300"},
		}

		sources, err := buildSources(entries)
		require.NoError(t, err)
		require.Len(t, sources, 1)
		assert.Equal(t, "forza", sources["my_forza"].Name())
	})

	t.Run("gt7 source", func(t *testing.T) {
		entries := map[string]config.SourceEntry{
			"my_gt7": {Type: "gt7", ConsoleAddr: "192.168.1.1", ListenAddr: ":33740"},
		}

		sources, err := buildSources(entries)
		require.NoError(t, err)
		require.Len(t, sources, 1)
		assert.Equal(t, "gt7", sources["my_gt7"].Name())
	})

	t.Run("multiple sources", func(t *testing.T) {
		entries := map[string]config.SourceEntry{
			"forza1": {Type: "forza", ListenAddr: ":5300"},
			"forza2": {Type: "forza", ListenAddr: ":5301"},
			"gt7":    {Type: "gt7", ConsoleAddr: "192.168.1.1", ListenAddr: ":33740"},
		}

		sources, err := buildSources(entries)
		require.NoError(t, err)
		assert.Len(t, sources, 3)
	})

	t.Run("unknown type", func(t *testing.T) {
		entries := map[string]config.SourceEntry{
			"bad": {Type: "unknown"},
		}

		_, err := buildSources(entries)
		assert.Error(t, err)
	})
}

func TestBuildSinks(t *testing.T) {
	t.Run("pcars2_udp sink", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"my_sink": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"},
		}

		sinks, err := buildSinks(entries)
		require.NoError(t, err)
		require.Len(t, sinks, 1)
		assert.Equal(t, "pcars2_udp", sinks["my_sink"].Sink.Name())
		assert.Equal(t, []string{"s1"}, sinks["my_sink"].Inputs)
	})

	t.Run("pcars1_udp sink", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"my_sink": {Type: "pcars1_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"},
		}

		sinks, err := buildSinks(entries)
		require.NoError(t, err)
		require.Len(t, sinks, 1)
		assert.Equal(t, "pcars1_udp", sinks["my_sink"].Sink.Name())
	})

	t.Run("json_stdout sink", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"debug": {Type: "json_stdout", Inputs: []string{"s1"}},
		}

		sinks, err := buildSinks(entries)
		require.NoError(t, err)
		require.Len(t, sinks, 1)
		assert.Equal(t, "json_stdout", sinks["debug"].Sink.Name())
	})

	t.Run("grafana_live sink", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"g": {Type: "grafana_live", Inputs: []string{"s1"}, Endpoint: "http://localhost:3000/api/live/push/test", APIKey: "key"},
		}

		sinks, err := buildSinks(entries)
		require.NoError(t, err)
		require.Len(t, sinks, 1)
		assert.Equal(t, "grafana_live", sinks["g"].Sink.Name())
	})

	t.Run("pcars1_shm sink unsupported on non-windows", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"shm": {Type: "pcars1_shm", Inputs: []string{"s1"}},
		}

		_, err := buildSinks(entries)
		assert.Error(t, err)
	})

	t.Run("pcars2_shm sink unsupported on non-windows", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"shm": {Type: "pcars2_shm", Inputs: []string{"s1"}},
		}

		_, err := buildSinks(entries)
		assert.Error(t, err)
	})

	t.Run("unknown type", func(t *testing.T) {
		entries := map[string]config.SinkEntry{
			"bad": {Type: "unknown", Inputs: []string{"s1"}},
		}

		_, err := buildSinks(entries)
		assert.Error(t, err)
	})
}

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
	}{
		{"info text", "info", "text"},
		{"debug json", "debug", "json"},
		{"warn text", "warn", "text"},
		{"error json", "error", "json"},
		{"default", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			setupLogger(config.LogConfig{Level: tt.level, Format: tt.format})
		})
	}
}

func skipIfNoNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("skipping: UDP sockets unavailable: %v", err)
	}

	conn.Close()
}
