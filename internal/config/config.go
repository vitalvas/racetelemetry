package config

import (
	"errors"
	"fmt"
)

// Config is the top-level configuration for racetelemetry.
type Config struct {
	Log     LogConfig              `yaml:"log" json:"log"`
	Sources map[string]SourceEntry `yaml:"sources" json:"sources"`
	Sinks   map[string]SinkEntry   `yaml:"sinks" json:"sinks"`
}

// LogConfig controls logging behavior.
type LogConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
}

// SourceEntry defines a named source instance.
type SourceEntry struct {
	Type string `yaml:"type" json:"type"`

	// Forza-specific
	ListenAddr string `yaml:"listen_addr,omitempty" json:"listen_addr,omitempty"`

	// GT7-specific
	ConsoleAddr string `yaml:"console_addr,omitempty" json:"console_addr,omitempty"`
}

// SinkEntry defines a named sink instance with input routing.
type SinkEntry struct {
	Type   string   `yaml:"type" json:"type"`
	Inputs []string `yaml:"inputs" json:"inputs"`

	// pcars1_udp / pcars2_udp specific
	TargetAddr string `yaml:"target_addr,omitempty" json:"target_addr,omitempty"`

	// grafana_live specific
	Endpoint string `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	APIKey   string `yaml:"api_key,omitempty" json:"api_key,omitempty"`
}

// Validate checks the configuration for correctness.
func (c *Config) Validate() error {
	if err := c.Log.validate(); err != nil {
		return err
	}

	if len(c.Sources) == 0 {
		return errors.New("config: no sources configured, at least one source is required")
	}

	if len(c.Sinks) == 0 {
		return errors.New("config: no sinks configured, at least one sink is required")
	}

	for name, src := range c.Sources {
		if err := src.validate(); err != nil {
			return fmt.Errorf("config: source %q: %w", name, err)
		}
	}

	for name, sink := range c.Sinks {
		if err := sink.validate(c.Sources); err != nil {
			return fmt.Errorf("config: sink %q: %w", name, err)
		}
	}

	return nil
}

func (l *LogConfig) validate() error {
	switch l.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("config: invalid log level %q, must be one of: debug, info, warn, error", l.Level)
	}

	switch l.Format {
	case "text", "json":
	default:
		return fmt.Errorf("config: invalid log format %q, must be one of: text, json", l.Format)
	}

	return nil
}

func (s *SourceEntry) validate() error {
	switch s.Type {
	case "forza":
		if s.ListenAddr == "" {
			return errors.New("listen_addr is required")
		}
	case "gt7":
		if s.ConsoleAddr == "" {
			return errors.New("console_addr is required")
		}

		if s.ListenAddr == "" {
			return errors.New("listen_addr is required")
		}
	case "wire":
		if s.ListenAddr == "" {
			return errors.New("listen_addr is required")
		}
	default:
		return fmt.Errorf("unknown source type %q", s.Type)
	}

	return nil
}

func (s *SinkEntry) validate(sources map[string]SourceEntry) error {
	switch s.Type {
	case "pcars1_udp":
		if s.TargetAddr == "" {
			return errors.New("target_addr is required")
		}
	case "pcars1_shm":
	case "pcars2_udp":
		if s.TargetAddr == "" {
			return errors.New("target_addr is required")
		}
	case "pcars2_shm":
	case "json_stdout":
	case "wire":
		if s.TargetAddr == "" {
			return errors.New("target_addr is required")
		}
	case "grafana_live":
		if s.Endpoint == "" {
			return errors.New("endpoint is required")
		}

		if s.APIKey == "" {
			return errors.New("api_key is required")
		}
	default:
		return fmt.Errorf("unknown sink type %q", s.Type)
	}

	if len(s.Inputs) == 0 {
		return errors.New("inputs is required, must reference at least one source")
	}

	for _, input := range s.Inputs {
		if _, ok := sources[input]; !ok {
			return fmt.Errorf("input %q references unknown source", input)
		}
	}

	return nil
}

// SetDefaults applies default values to the configuration.
func (c *Config) SetDefaults() {
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}

	if c.Log.Format == "" {
		c.Log.Format = "text"
	}
}
