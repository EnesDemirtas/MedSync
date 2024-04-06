package inventoryapp

import (
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

func parseFilter(qp QueryParams) (inventorybus.QueryFilter, error) {
	var filter inventorybus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return inventorybus.QueryFilter{}, validate.NewFieldsError("inventory_id", err)
		}
		filter.WithInventoryID(id)
	}

	if qp.Name != "" {
		filter.WithName(qp.Name)
	}

	if qp.Description != "" {
		filter.WithDescription(qp.Description)
	}

	// TODO: Add missing filters.

	return filter, nil
}