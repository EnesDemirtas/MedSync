package inventoryapp

import (
	"errors"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByInventoryID = "inventory_id"
		orderByName 	   = "name"
		orderByDescription = "description"
	)

	var orderByFields = map[string]string{
		orderByInventoryID: inventorybus.OrderByID,
		orderByName: 		inventorybus.OrderByName,
		orderByDescription: inventorybus.OrderByDescription,
	}

	orderBy, err := order.Parse(qp.OrderBy, order.NewBy(orderByInventoryID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}