package config

import "os"

var (
	errInvalidInt    = &errInvalidType{"invalid integer"}
	errInvalidInt64  = &errInvalidType{"invalid int64"}
	errInvalidBool   = &errInvalidType{"invalid boolean"}
)

type errInvalidType struct {
	msg string
}

func (e *errInvalidType) Error() string {
	return e.msg
}

func getEnv[T any](key string, defaultValue T) T {
	if value := os.Getenv(key); value != "" {
		return parseEnv[T](value, defaultValue)
	}
	return defaultValue
}

func parseEnv[T any](value string, defaultValue T) T {
	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T)
	case int:
		if v, err := strconvAtoi(value); err == nil {
			return any(v).(T)
		}
	case int64:
		if v, err := strconvParseInt(value); err == nil {
			return any(v).(T)
		}
	case bool:
		if v, err := strconvParseBool(value); err == nil {
			return any(v).(T)
		}
	}
	return defaultValue
}

func strconvAtoi(s string) (int, error) {
	i := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errInvalidInt
		}
		i = i*10 + int(c-'0')
	}
	return i, nil
}

func strconvParseInt(s string) (int64, error) {
	var v int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errInvalidInt64
		}
		v = v*10 + int64(c-'0')
	}
	return v, nil
}

func strconvParseBool(s string) (bool, error) {
	switch s {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	}
	return false, errInvalidBool
}
