package inventory

import (
	"fmt"
	"time"

	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID					*uuid.UUID
	Name 				*string	`validate:"omitempty,min=3"`
	Description 		*string
	Medicine			*uuid.UUID
	Medicines 			[]uuid.UUID
	StartExpiryDate		*time.Time
	EndExpiryDate		*time.Time
}

// Validate can perform a check of tha data against the validate tags.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// WithInventoryID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithInventoryID(inventoryID uuid.UUID) {
	qf.ID = &inventoryID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithDescription sets the Description field of the QueryFilter value.
func (qf *QueryFilter) WithDescription(description string) {
	qf.Description = &description
}

// WithMedicine sets the Medicine field of the QueryFilter value.
func (qf *QueryFilter) WithMedicine(medicineID uuid.UUID) {
	qf.Medicine = &medicineID
}

// WithMedicines sets the Medicines field of the QueryFilter value.
func (qf *QueryFilter) WithMedicines(medicineIDs []uuid.UUID) {
	qf.Medicines = medicineIDs
}