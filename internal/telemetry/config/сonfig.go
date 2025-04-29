package config

import "time"

type Config struct {
	Enabled     bool
	ServiceName string
	Insecure    bool
	Trace       TraceConfig
	Metrics     MetricsConfig
}

type TraceConfig struct {
	Endpoint     string
	Headers      map[string]string
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

func WithEndpoint(e string) Option {
	return func(t *Config) {
		t.Endpoint = e
	}
}

func WithService(s string) Option {
	return func(t *Config) {
		t.Service = s
	}
}

func WithInsecure() Option {

}
