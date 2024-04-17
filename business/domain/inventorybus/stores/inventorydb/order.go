package inventorydb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
)

var orderByFields = map[string]string{
	inventorybus.OrderByID: 			"inventory_id",
	inventorybus.OrderByName:			"name",
	inventorybus.OrderByDescription:	"description",
	inventorybus.OrderByExpiryDate:	"expiry_date",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}