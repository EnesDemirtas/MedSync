package userapp

import (
	"errors"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/domain/userbus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByID 		= "user_id"
		orderByName		= "name"
		orderByEmail 	= "email"
		orderByRoles 	= "roles"
		orderByEnabled  = "enabled"
	)

	var orderByFields = map[string]string{
		orderByID: 		userbus.OrderByID,
		orderByName:	userbus.OrderByName,
		orderByEmail:   userbus.OrderByEmail,
		orderByRoles: 	userbus.OrderByRoles,
		orderByEnabled: userbus.OrderByEnabled,
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