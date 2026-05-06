package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name:    "no sources",
			config:  Config{Log: LogConfig{Level: "info", Format: "text"}},
			wantErr: "no sources configured",
		},
		{
			name: "no sinks",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
			},
			wantErr: "no sinks configured",
		},
		{
			name: "unknown source type",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "unknown"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "unknown source type",
		},
		{
			name: "forza missing listen_addr",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "listen_addr is required",
		},
		{
			name: "gt7 missing console_addr",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "gt7", ListenAddr: ":33740"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "console_addr is required",
		},
		{
			name: "gt7 missing listen_addr",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "gt7", ConsoleAddr: "192.168.1.1"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "listen_addr is required",
		},
		{
			name: "unknown sink type",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "unknown", Inputs: []string{"s1"}}},
			},
			wantErr: "unknown sink type",
		},
		{
			name: "sink missing inputs",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "inputs is required",
		},
		{
			name: "sink references unknown source",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"nonexistent"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "references unknown source",
		},
		{
			name: "pcars2_udp missing target_addr",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}}},
			},
			wantErr: "target_addr is required",
		},
		{
			name: "invalid log level",
			config: Config{
				Log:     LogConfig{Level: "bad", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "invalid log level",
		},
		{
			name: "invalid log format",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "xml"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
			wantErr: "invalid log format",
		},
		{
			name: "valid single source single sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars2_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
		},
		{
			name: "valid multiple sources multiple sinks with routing",
			config: Config{
				Log: LogConfig{Level: "debug", Format: "json"},
				Sources: map[string]SourceEntry{
					"forza_xbox": {Type: "forza", ListenAddr: ":5300"},
					"forza_pc":   {Type: "forza", ListenAddr: ":5301"},
					"gt7_ps5":    {Type: "gt7", ConsoleAddr: "192.168.1.100", ListenAddr: ":33740"},
				},
				Sinks: map[string]SinkEntry{
					"cc_forza": {Type: "pcars2_udp", Inputs: []string{"forza_xbox", "forza_pc"}, TargetAddr: "127.0.0.1:5606"},
					"cc_gt7":   {Type: "pcars2_udp", Inputs: []string{"gt7_ps5"}, TargetAddr: "127.0.0.1:5607"},
				},
			},
		},
		{
			name: "pcars1_udp missing target_addr",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars1_udp", Inputs: []string{"s1"}}},
			},
			wantErr: "target_addr is required",
		},
		{
			name: "grafana_live missing endpoint",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "grafana_live", Inputs: []string{"s1"}, APIKey: "key"}},
			},
			wantErr: "endpoint is required",
		},
		{
			name: "grafana_live missing api_key",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "grafana_live", Inputs: []string{"s1"}, Endpoint: "http://localhost:3000/api/live/push/test"}},
			},
			wantErr: "api_key is required",
		},
		{
			name: "valid pcars1_udp sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "pcars1_udp", Inputs: []string{"s1"}, TargetAddr: "127.0.0.1:5606"}},
			},
		},
		{
			name: "valid pcars1_shm sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"shm": {Type: "pcars1_shm", Inputs: []string{"s1"}}},
			},
		},
		{
			name: "valid pcars2_shm sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"shm": {Type: "pcars2_shm", Inputs: []string{"s1"}}},
			},
		},
		{
			name: "valid json_stdout sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "json_stdout", Inputs: []string{"s1"}}},
			},
		},
		{
			name: "valid grafana_live sink",
			config: Config{
				Log:     LogConfig{Level: "info", Format: "text"},
				Sources: map[string]SourceEntry{"s1": {Type: "forza", ListenAddr: ":5300"}},
				Sinks:   map[string]SinkEntry{"k1": {Type: "grafana_live", Inputs: []string{"s1"}, Endpoint: "http://localhost:3000/api/live/push/test", APIKey: "key"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_SetDefaults(t *testing.T) {
	t.Run("applies log defaults", func(t *testing.T) {
		cfg := Config{}
		cfg.SetDefaults()

		assert.Equal(t, "info", cfg.Log.Level)
		assert.Equal(t, "text", cfg.Log.Format)
	})

	t.Run("preserves explicit values", func(t *testing.T) {
		cfg := Config{
			Log: LogConfig{Level: "debug", Format: "json"},
		}
		cfg.SetDefaults()

		assert.Equal(t, "debug", cfg.Log.Level)
		assert.Equal(t, "json", cfg.Log.Format)
	})
}
