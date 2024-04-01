package tagdb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/tag"
	"github.com/EnesDemirtas/medisync/business/web/order"
)

var orderByFields = map[string]string {
	tag.OrderByID: 	 "tag_id",
	tag.OrderByName: "name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}