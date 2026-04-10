package config

type OTELConfig struct {
	ExporterType string
	Enabled      bool
}

func LoadOTELConfig() OTELConfig {
	return OTELConfig{
		ExporterType: getEnv("OTEL_EXPORTER_TYPE", "stdout"),
		Enabled:      isOTELEnabled(),
	}
}

func isOTELEnabled() bool {
	return getEnv("OTEL_ENABLED", "true") == "true"
}
