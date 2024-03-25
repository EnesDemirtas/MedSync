package medicine

import "github.com/google/uuid"

// Tag represents a single tag.
type Tag struct {
	ID		uuid.UUID
	name 	string
}

// Name returns the name of a tag.
func (t Tag) Name() string {
	return t.name
}