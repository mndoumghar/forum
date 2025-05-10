package utils

<<<<<<< HEAD
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
=======
// import (
// )		
// GenerateUUID generates a random UUID (Universally Unique Identifier).
>>>>>>> 45193a583d02665e59fba785b999e8bf16e9d9b3
