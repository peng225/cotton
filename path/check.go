package path

import "github.com/google/uuid"

const (
	maxPathLength = 1024
)

func Valid(path string) bool {
	if len(path) > maxPathLength {
		return false
	}
	return true
}

func ValidateUUID(value string) (*uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}
	return &parsedUUID, nil
}
