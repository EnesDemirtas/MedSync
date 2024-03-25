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