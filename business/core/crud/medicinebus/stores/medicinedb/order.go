package medicinedb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
)

var orderByFields = map[string]string{
	medicinebus.OrderByID: 			"medicine_id",
	medicinebus.OrderByName:			"name",
	medicinebus.OrderByDescription:	"description",
	medicinebus.OrderByManufacturer:	"manufacturer",
	medicinebus.OrderByType:			"type",
	medicinebus.OrderByExpiryDate:		"expiry_date",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}