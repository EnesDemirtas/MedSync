package inventory

import (
	"time"

	"github.com/google/uuid"
)

// Inventory represents a single inventory that keeps medicine(s) in itself.
type Inventory struct {
	ID 			uuid.UUID
	Name		string
	Description string
	Medicines 	[]uuid.UUID
	DateCreated time.Time
	DateUpdated	time.Time
}

// NewInventory contains information needed to create a new inventory.
type NewInventory struct {
	Name		string
	Description	string
	Medicines	[]uuid.UUID
}

// UpdateInventory contains information needed to update an inventory.
type UpdateInventory struct {
	Name 		*string
	Description	*string
	Medicines	[]uuid.UUID
}