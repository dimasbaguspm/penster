package config

type RateCurrencyConfig struct {
	ECBURL string
}

func LoadRateCurrencyConfig() RateCurrencyConfig {
	return RateCurrencyConfig{
		ECBURL: getEnv("ECB_URL", "http://www.ecb.int/vocabulary/2002-08-01/eurofxref"),
	}
}
