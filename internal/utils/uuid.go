package utils

import (
	"github.com/gofrs/uuid"
	"log"
)

// GenerateUUID generates a new UUID.
func NewUUID() string {
	// Generate a new UUID
	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	return newUUID.String()
}
