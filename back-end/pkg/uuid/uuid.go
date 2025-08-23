package uuid

import "github.com/google/uuid"

func GenerateUUID() (string, error) {
	uniqueID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uniqueID.String(), nil
}
