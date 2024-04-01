package inventorydb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/inventory"
	"github.com/EnesDemirtas/medisync/business/web/order"
)

var orderByFields = map[string]string{
	inventory.OrderByID: 			"inventory_id",
	inventory.OrderByName:			"name",
	inventory.OrderByDescription:	"description",
	inventory.OrderByExpiryDate:	"expiry_date",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}