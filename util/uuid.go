package util

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// NewUuid Create new Uuid v4
func NewUuid() string {
	id := uuid.NewV4()
	return fmt.Sprintf("%s", id)
}

// IsUuid validate uuid
func IsUuid(text string) bool {
	_, err := uuid.FromString(text)
	if err != nil {
		return false
	}
	return true
}
