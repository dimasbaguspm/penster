package config

type MigrateConfig struct {
	AutoMigrate bool
	Path        string
}

func LoadMigrateConfig() MigrateConfig {
	return MigrateConfig{
		AutoMigrate: getEnv("AUTO_MIGRATE", true),
		Path:        getEnv("MIGRATE_PATH", "migrations"),
	}
}
