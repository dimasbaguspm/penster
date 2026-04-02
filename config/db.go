package config

import (
	"fmt"
)

type DBConfig struct {
	Primary  string
	MaxConns int
	MinConns int
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Primary:  fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", getEnv("DB_USER", "penster"), getEnv("DB_PASSWORD", "placeholder"), getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", 5432), getEnv("DB_NAME", "penster"), getEnv("DB_SSLMODE", "disable")),
		MaxConns: getEnv("DB_MAX_CONNS", 10),
		MinConns: getEnv("DB_MIN_CONNS", 2),
	}
}
