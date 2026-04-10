package config

type OTELConfig struct {
	Enabled  bool
	Endpoint string
}

func LoadOTELConfig() OTELConfig {
	return OTELConfig{
		Enabled:  getEnv("OTEL_ENABLED", "true") == "true",
		Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
	}
}
