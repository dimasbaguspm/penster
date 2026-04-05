package dto

import "github.com/google/uuid"

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
