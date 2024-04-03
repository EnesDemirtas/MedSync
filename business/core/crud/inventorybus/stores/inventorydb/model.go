package inventorydb

import (
	"database/sql"
	"time"

	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/business/data/sqldb/dbarray"
	"github.com/google/uuid"
)

// TODO: Keep track of number of medicines.

type dbInventory struct {
	ID 			 uuid.UUID		`db:"inventory_id"`
	Name		 string			`db:"name"`
	Description  sql.NullString	`db:"description"`
	Medicines	 dbarray.String	`db:"medicines"`
	DateCreated  time.Time		`db:"date_created"`
	DateUpdated  time.Time		`db:"date_updated"`
}

func toDBInventory(inv inventorybus.Inventory) dbInventory {
	meds := make([]string, len(inv.Medicines))
	for i, med := range inv.Medicines {
		meds[i] = med.String()
	}

	return dbInventory{
		ID:			  inv.ID,
		Name:		  inv.Name,
		Description:  sql.NullString{
			String: inv.Description,
			Valid:	inv.Description != "",
		},
		Medicines: 	  meds,
		DateCreated:  inv.DateCreated,
		DateUpdated:  inv.DateUpdated,
	}
}

func toCoreInventory(dbInventory dbInventory) inventorybus.Inventory {
	meds := make([]uuid.UUID, len(dbInventory.Medicines))
	for i, med := range dbInventory.Medicines {
		meds[i] = uuid.MustParse(med)
	}
	inv := inventorybus.Inventory{
		ID:			  dbInventory.ID,
		Name:		  dbInventory.Name,
		Description:  dbInventory.Description.String,
		Medicines: 	  meds,
		DateCreated:  dbInventory.DateCreated,
		DateUpdated:  dbInventory.DateUpdated,
	}

	return inv
}

func toCoreInventorySlice(dbInventories []dbInventory) []inventorybus.Inventory {
	invs := make([]inventorybus.Inventory, len(dbInventories))

	for i, dbInventory := range dbInventories {
		invs[i] = toCoreInventory(dbInventory)
	}

	return invs
}