package medicineapp

import (
	"errors"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByID 			= "medicine_id"
		orderByName 		= "name"
		orderByDescription  = "description"
		orderByManufacturer = "manufacturer"
		orderByType			= "type"
		orderByExpiryDate	= "expiry_date"
	)

	var orderByFields = map[string]string{
		orderByID: 			 medicinebus.OrderByID,
		orderByName:		 medicinebus.OrderByName,
		orderByDescription:  medicinebus.OrderByDescription,
		orderByManufacturer: medicinebus.OrderByManufacturer,
		orderByType:		 medicinebus.OrderByType,
		orderByExpiryDate:   medicinebus.OrderByExpiryDate,
	}

	orderBy, err := order.Parse(qp.OrderBy, order.NewBy(orderByID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}