package config

type AppConfig struct {
	Env          string
	Port         string
	Version      string
	BaseCurrency string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Env:          getEnv("APP_ENV", "local"),
		Port:         getEnv("APP_PORT", "8080"),
		Version:      getEnv("APP_VERSION", "1.0.0"),
		BaseCurrency: getEnv("BASE_CURRENCY", "IDR"),
	}
}
