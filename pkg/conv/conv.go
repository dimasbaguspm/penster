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
