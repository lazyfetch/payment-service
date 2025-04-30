package config

import "time"

type Config struct {
	ServiceName string
	Insecure    bool
	Traces      TracesConfig
	Metrics     MetricsConfig
}

type TracesConfig struct {
	Endpoint     string
	Timeout      time.Duration
	Sampler      string
	SamplerRatio float64
}

type MetricsConfig struct {
	Endpoint string
	Timeout  time.Duration
	Interval time.Duration
}

type Option func(*Config)

func TracesWithEndpoint(s string) Option {
	return func(c *Config) {
		c.Traces.Endpoint = s
	}
}

func MetricsWithEndpoint(s string) Option {
	return func(c *Config) {
		c.Metrics.Endpoint = s
	}
}

// Setup service name
func WithService(s string) Option {
	return func(c *Config) {
		c.ServiceName = s
	}
}

// true = no SSL certificate
func WithInsecure(b bool) Option {
	return func(c *Config) {
		c.Insecure = true
	}
}

func MetricsWithTimeout(t time.Duration) Option {
	return func(c *Config) {
		c.Metrics.Timeout = t
	}
}

func MetricsWithInterval(t time.Duration) Option {
	return func(c *Config) {
		c.Metrics.Interval = t
	}
}

func TracesWithTimeout(t time.Duration) Option {
	return func(c *Config) {
		c.Traces.Timeout = t
	}
}

func TracesWithSampler(s string) Option {
	return func(c *Config) {
		c.Traces.Sampler = s
	}
}

func TracesWithRatio(f float64) Option {
	return func(c *Config) {
		c.Traces.SamplerRatio = f
	}
}
