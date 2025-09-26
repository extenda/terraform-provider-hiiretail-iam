package client

import "fmt"

// ResourceNotFoundError represents an error when a resource is not found
type ResourceNotFoundError struct {
	ResourceType string
	ID          string
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ID)
}

// IsResourceNotFound checks if the given error is a ResourceNotFoundError
func IsResourceNotFound(err error) bool {
	_, ok := err.(*ResourceNotFoundError)
	return ok
}