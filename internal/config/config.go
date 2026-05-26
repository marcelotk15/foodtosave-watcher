package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultInterval      = 5 * time.Minute
	DefaultNtfyServerURL = "https://ntfy.sh"
	DefaultConfigPath    = "./config.yaml"
	DefaultCachePath     = "./cache.json"
)

// Config is the root structure loaded from config.yaml.
type Config struct {
	Interval       DurationYAML  `yaml:"interval"`
	Ntfy           NtfyConfig    `yaml:"ntfy"`
	Merchants      []Merchant    `yaml:"merchants"`
	CachePath      string        `yaml:"cache_path"`
	intervalParsed time.Duration `yaml:"-"`
}

// NtfyConfig defines the destination of notifications.
type NtfyConfig struct {
	Topic     string `yaml:"topic"`
	ServerURL string `yaml:"server_url"`
}

// Merchant represents a store in the Food To Save API.
type Merchant struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type DurationYAML time.Duration

func (d *DurationYAML) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		raw := strings.TrimSpace(n.Value)
		if raw == "" {
			return nil
		}
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("campo interval: %w", err)
		}
		*d = DurationYAML(parsed)
		return nil
	}
	var s string
	if err := n.Decode(&s); err != nil {
		return err
	}
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("campo interval: %w", err)
	}
	*d = DurationYAML(parsed)
	return nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ler config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml inválido: %w", err)
	}
	if err := cfg.validateAndNormalize(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validateAndNormalize() error {
	var errs []error

	dur := time.Duration(c.Interval)
	if dur <= 0 {
		dur = DefaultInterval
	}
	c.intervalParsed = dur

	if strings.TrimSpace(c.CachePath) == "" {
		c.CachePath = DefaultCachePath
	}

	for i := range c.Merchants {
		m := &c.Merchants[i]
		if strings.TrimSpace(m.ID) == "" {
			errs = append(errs, fmt.Errorf("merchant[%d]: id is empty", i))
		}
		if strings.TrimSpace(m.Name) == "" {
			errs = append(errs, fmt.Errorf("merchant[%d]: name is empty", i))
		}
	}

	if len(c.Merchants) < 1 {
		errs = append(errs, errors.New("merchants: at least one merchant is required"))
	}

	if strings.TrimSpace(c.Ntfy.Topic) == "" {
		errs = append(errs, errors.New("ntfy.topic: required and cannot be empty"))
	}

	rawURL := strings.TrimSpace(c.Ntfy.ServerURL)
	if rawURL == "" {
		rawURL = DefaultNtfyServerURL
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		errs = append(errs, fmt.Errorf("ntfy.server_url: invalid URL %q", c.Ntfy.ServerURL))
	} else {
		c.Ntfy.ServerURL = strings.TrimRight(rawURL, "/")
	}

	return errors.Join(errs...)
}

func (c *Config) IntervalDuration() time.Duration {
	return c.intervalParsed
}
