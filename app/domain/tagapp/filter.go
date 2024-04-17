package tagapp

import (
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

func parseFilter(qp QueryParams) (tagbus.QueryFilter, error) {
	var filter tagbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return tagbus.QueryFilter{}, validate.NewFieldsError("tag_id", err)
		}
		filter.WithID(id)
	}

	if qp.Name != "" {
		filter.WithName(qp.Name)
	}

	return filter, nil
}
