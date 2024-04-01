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
	Tags 			[]uuid.UUID
	ExpiryDate		time.Time
	DateCreated		time.Time
	DateUpdated		time.Time
}

// NewMedicine contains information needed to create a new medicine.
type NewMedicine struct {
	Name 			string
	Description		string
	Manufacturer	string
	Type 			string
	Tags			[]uuid.UUID
	ExpiryDate		time.Time
}

// UpdateMedicine contains information needed to update a medicine.
type UpdateMedicine struct {
	Name			*string
	Description		*string
	Manufacturer	*string
	Type 			*string
	Tags			[]uuid.UUID
	ExpiryDate		*time.Time
}