package config

import (
	"fmt"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", 5432),
		User:     getEnv("DB_USER", "penster"),
		Password: getEnv("DB_PASSWORD", "placeholder"),
		Name:     getEnv("DB_NAME", "penster"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}
