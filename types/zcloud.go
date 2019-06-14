package types

type MonitorConfig struct {
	// Monitoring server provider
	MetricsProvider string `yaml:"metrics_provider" json:"metrics_provider,omitempty" norman:"default=metrics-server"`
	// Metrics server options
	MetricsOptions map[string]string `yaml:"metrics_options" json:"metrics_options,omitempty"`
}
