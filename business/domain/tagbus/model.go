package tagbus

import "github.com/google/uuid"

// Tag represents a single tag.
type Tag struct {
	ID		uuid.UUID
	Name 	string
}

// NewTag contains information needed to create a new tag.
type NewTag struct {
	Name string
}

// UpdateTag contains information needed to update a tag.
type UpdateTag struct {
	Name *string
}