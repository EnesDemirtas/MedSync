package tag

import "github.com/EnesDemirtas/medisync/business/web/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID			= "tag_id"
	OrderByName			= "name"
)