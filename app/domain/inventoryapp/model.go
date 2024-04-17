package inventoryapp

import (
	"time"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page 		int		`json:"page"`
	Rows		int		`json:"rows"`
	OrderBy 	string  `json:"orderBy"`
	ID 			string 	`json:"inventory_id"`
	Name		string  `json:"name"`
	Description string  `json:"description"`
}

// Inventory represents information about an indiviual inventory.
type Inventory struct {
	ID 				   string 		  `json:"id"`
	Name    		   string 		  `json:"name"`
	Description 	   string 		  `json:"description"`
	MedicineQuantities map[string]int `json:"medicineQuantities"`
	DateCreated 	   string 		  `json:"dateCreated"`
	DateUpdated 	   string 		  `json:"dateUpdated"`
}

func toAppInventory(inv inventorybus.Inventory) Inventory {
	medQua := make(map[string]int, len(inv.MedicineQuantities))
	for k, v := range inv.MedicineQuantities {
		medQua[k.String()] = v
	}

	return Inventory{
		ID:			 inv.ID.String(),
		Name:		 inv.Name,
		Description: inv.Description,
		MedicineQuantities: medQua,
		DateCreated: inv.DateCreated.Format(time.RFC3339),
		DateUpdated: inv.DateUpdated.Format(time.RFC3339),
	}
}

func toAppInventories(invs []inventorybus.Inventory) []Inventory {
	items := make([]Inventory, len(invs))
	for i, inv := range invs {
		items[i] = toAppInventory(inv)
	}

	return items
}

// NewInventory defines the data needed to adda new inventory.
type NewInventory struct {
	Name			   string `json:"name" validate:"required"`
	Description 	   string `json:"description"`
}

func toBusNewInventory(app NewInventory) inventorybus.NewInventory {
	inv := inventorybus.NewInventory{
		Name: 		 app.Name,
		Description: app.Description,
	}

	return inv
}

// Validate checks the data in the model is considered clean.
func (app NewInventory) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}

// UpdateInventory defines the data needed to update an inventory.
type UpdateInventory struct {
	Name 			   *string 		  `json:"name"`
	Description        *string 		  `json:"description"`
	MedicineQuantities map[string]int `json:"medicineQuantities" validate:"omitempty"`
}

func toBusUpdateInventory(app UpdateInventory) (inventorybus.UpdateInventory, error) {
	meds := make(map[uuid.UUID]int, len(app.MedicineQuantities))
	for idStr, qua := range app.MedicineQuantities {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return inventorybus.UpdateInventory{}, errs.Newf(errs.Internal, "parse: medicineID[%s]: %s", idStr, err)
		}
		meds[id] = qua
	}

	// TODO: Check specified medicine ids are valid.

	inv := inventorybus.UpdateInventory{
		Name: 			app.Name,
		Description:    app.Description,
		MedicineQuantities: meds,
	}

	return inv, nil
}

// Validate checks the data in the model is considered clean.
func (app UpdateInventory) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}