package inventorybus

import "github.com/EnesDemirtas/medisync/business/api/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID			= "inventory_id"
	OrderByName			= "name"
	OrderByDescription	= "description"
	OrderByExpiryDate   = "expiry_date"
)