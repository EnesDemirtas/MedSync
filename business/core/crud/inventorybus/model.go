package inventorybus

import (
	"time"

	"github.com/google/uuid"
)

// TODO: Keep track of number of medicines.

// Inventory represents a single inventory that keeps medicine(s) in itself.
type Inventory struct {
	ID 					uuid.UUID
	Name				string
	Description 		string
	MedicineQuantities 	map[uuid.UUID]int
	DateCreated 		time.Time
	DateUpdated			time.Time
}

// NewInventory contains information needed to create a new inventory.
type NewInventory struct {
	Name				string
	Description			string
}

// UpdateInventory contains information needed to update an inventory.
type UpdateInventory struct {
	Name 				*string
	Description			*string
	MedicineQuantities	map[uuid.UUID]int
}