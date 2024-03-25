package inventory

import "github.com/google/uuid"

// Inventory represents a single inventory that keeps medicine(s) in itself.
type Inventory struct {
	ID 			uuid.UUID
	Name		string
	Description string
	Medicines 	[]uuid.UUID
}