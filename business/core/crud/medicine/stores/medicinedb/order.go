package medicinedb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/medicine"
	"github.com/EnesDemirtas/medisync/business/web/order"
)

var orderByFields = map[string]string{
	medicine.OrderByID: 			"medicine_id",
	medicine.OrderByName:			"name",
	medicine.OrderByDescription:	"description",
	medicine.OrderByManufacturer:	"manufacturer",
	medicine.OrderByType:			"type",
	medicine.OrderByExpiryDate:		"expiry_date",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}