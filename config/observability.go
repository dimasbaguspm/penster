package config

type ObservabilityConfig struct {
	Enabled         bool
	TracingEndpoint string
	LogsEndpoint    string
	MetricsEndpoint string
}

func LoadObservabilityConfig() ObservabilityConfig {
	return ObservabilityConfig{
		Enabled:         getEnv("OTEL_ENABLED", "true") == "true",
		TracingEndpoint: getEnv("TEMPO_ENDPOINT", "localhost:4317"),
		LogsEndpoint:    getEnv("LOKI_ENDPOINT", "localhost:3100"),
		MetricsEndpoint: getEnv("MIMIR_ENDPOINT", "localhost:9009"),
	}
}