package medicinedb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
)

func applyFilter(filter medicinebus.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["medicine_id"] = *filter.ID
		wc = append(wc, "medicine_id = :medicine_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.Description != nil {
		data["description"] = fmt.Sprintf("%%%s%%", *filter.Description)
		wc = append(wc, "description LIKE :description")
	}

	if filter.Manufacturer != nil {
		data["manufacturer"] = fmt.Sprintf("%%%s%%", *filter.Manufacturer)
		wc = append(wc, "manufacturer LIKE :manufacturer")
	}

	if filter.Type != nil {
		data["type"] = fmt.Sprintf("%%%s%%", *filter.Tag)
		wc = append(wc, "type LIKE :type")
	}

	if filter.StartExpiryDate != nil {
		data["start_expiry_date"] = *filter.StartExpiryDate
		wc = append(wc, "expiry_date >= :start_expiry_date")
	}

	if filter.EndExpiryDate != nil {
		data["end_expiry_date"] = *filter.EndExpiryDate
		wc = append(wc, "expiry_date <= :end_expiry_date")
	}


	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}	
}