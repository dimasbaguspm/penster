package conv

import (
	"github.com/google/uuid"
)

func ParseUUID(s string) [16]byte {
	u, err := uuid.Parse(s)
	if err != nil {
		return [16]byte{}
	}
	return u
}

func StringPtrToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringPtrToNull converts a nil string pointer to an empty string
func StringPtrToNull(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int64PtrToInt64 converts a nil int64 pointer to 0
func Int64PtrToInt64(s *int64) int64 {
	if s == nil {
		return 0
	}
	return *s
}
