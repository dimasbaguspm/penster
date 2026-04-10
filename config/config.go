package config

import "github.com/joho/godotenv"

type Config struct {
	App          AppConfig
	DB           DBConfig
	Kafka        KafkaConfig
	Migrate      MigrateConfig
	RateCurrency RateCurrencyConfig
	OTEL         OTELConfig
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App:          LoadAppConfig(),
		DB:           LoadDBConfig(),
		Kafka:        LoadKafkaConfig(),
		Migrate:      LoadMigrateConfig(),
		RateCurrency: LoadRateCurrencyConfig(),
		OTEL:         LoadOTELConfig(),
	}
}
