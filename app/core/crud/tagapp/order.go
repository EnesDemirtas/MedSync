package tagapp

import (
	"errors"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByTagID     = "tag_id"
		orderByName      = "name"
	)

	var orderByFields = map[string]string{
		orderByTagID: 		tagbus.OrderByID,
		orderByName:      tagbus.OrderByName,
	}

	orderBy, err := order.Parse(qp.OrderBy, order.NewBy(orderByTagID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
