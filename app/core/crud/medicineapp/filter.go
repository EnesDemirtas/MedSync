package medicineapp

import (
	"time"

	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

func parseFilter(qp QueryParams) (medicinebus.QueryFilter, error) {
	var filter medicinebus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return medicinebus.QueryFilter{}, validate.NewFieldsError("medicine_id", err)
		}
		filter.WithMedicineID(id)
	}

	if qp.Name != "" {
		filter.WithName(qp.Name)
	}

	if qp.Description != "" {
		filter.WithDescription(qp.Description)
	}

	if qp.Manufacturer != "" {
		filter.WithManufacturer(qp.Manufacturer)
	}

	if qp.Type != "" {
		filter.WithType(qp.Type)
	}

	if qp.Tags != nil {
		tags := make([]uuid.UUID, len(qp.Tags))
		for i, tagStr := range qp.Tags {
			tag, err := uuid.Parse(tagStr)
			if err != nil {
				return medicinebus.QueryFilter{}, validate.NewFieldsError("tags", err)
			}
			tags[i] = tag
		}
		filter.WithTags(tags)
	}

	if qp.StartExpiryDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartExpiryDate)
		if err != nil {
			return medicinebus.QueryFilter{}, validate.NewFieldsError("start_expiry_date", err)
		}
		filter.WithStartExpiryDate(t)
	}

	if qp.EndExpiryDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndExpiryDate)
		if err != nil {
			return medicinebus.QueryFilter{}, validate.NewFieldsError("end_expiry_date", err)
		}
		filter.WithEndExpiryDate(t)
	}

	// TODO: Add StartCreatedDate, EndCreatedDate

	return filter, nil
}