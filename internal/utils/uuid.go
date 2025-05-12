package utils

import "github.com/gofrs/uuid"

func GenerateUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
