package medicine

import (
	"time"

	"github.com/google/uuid"
)

// Medicine represents information about a single medicine.
type Medicine struct {
	ID 				uuid.UUID
	Name 			string
	Description 	string
	Manufacturer	string
	Type 			string
	ExpiryDate		time.Time
	CreatedDate		time.Time
	UpdatedDate		time.Time
}