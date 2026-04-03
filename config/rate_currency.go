package config

import "time"

type RateCurrencyConfig struct {
	Interval time.Duration
	ECBURL   string
}

func LoadRateCurrencyConfig() RateCurrencyConfig {
	return RateCurrencyConfig{
		Interval: time.Hour,
		ECBURL:   getEnv("ECB_URL", "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"),
	}
}
