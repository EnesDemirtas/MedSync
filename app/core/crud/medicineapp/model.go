package medicineapp

import (
	"fmt"
	"time"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
	"github.com/google/uuid"
)

// TODO: Add StartCreatedDate, EndCreatedDate for QueryParams

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page 			 int 	`query:"page"`
	Rows			 int 	`query:"rows"`
	OrderBy 		 string	`query:"orderBy"`
	ID 				 string	`query:"medicine_id"`
	Name			 string	`query:"name"`
	Description      string `query:"desctiption"`
	Manufacturer     string `query:"manufacturer"`
	Type   			 string `query:"type"`
	Tags			 []string `query:"tags"`
	StartExpiryDate  string `query:"start_expiry_date"`
	EndExpiryDate    string `query:"end_expiry_date"`
}

// Medicine represents information about an individual medicine.
type Medicine struct {
	ID			 string   `json:"id"`
	Name		 string	  `json:"name"`
	Description  string   `json:"description"`
	Manufacturer string   `json:"manufacturer"`
	Type   		 string   `json:"type"`
	Tags		 []string `json:"tags"`
	ExpiryDate   string   `json:"expiryDate"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

func toAppMedicine(med medicinebus.Medicine) Medicine {
	tags := make([]string, len(med.Tags))
	for i, tag := range med.Tags {
		tags[i] = tag.String()
	}

	return Medicine{
		ID:			  med.ID.String(),
		Name:		  med.Name,
		Description:  med.Description,
		Manufacturer: med.Manufacturer,
		Type:		  med.Type,
		Tags:		  tags,
		ExpiryDate:   med.ExpiryDate.Format(time.RFC3339),
		DateCreated:  med.DateCreated.Format(time.RFC3339),
		DateUpdated:  med.DateUpdated.Format(time.RFC3339),
	}
}

func toAppMedicines(meds []medicinebus.Medicine) []Medicine {
	items := make([]Medicine, len(meds))
	for i, med := range meds {
		items[i] = toAppMedicine(med)
	}

	return items
}

// NewMedicine defines the data needed to add a new medicine.
type NewMedicine struct {
	Name 		 string   `json:"name" validate:"required"`
	Description  string   `json:"description"`
	Manufacturer string   `json:"manufacturer"`
	Type         string   `json:"type"`
	Tags         []string `json:"tags"`
	ExpiryDate   string   `json:"expiryDate"`
}

func toBusNewMedicine(app NewMedicine) (medicinebus.NewMedicine, error) {
	var tags []uuid.UUID
	if app.Tags != nil {
		tags = make([]uuid.UUID, len(app.Tags))
		for i, tagIDStr := range app.Tags {
			tag, err := uuid.Parse(tagIDStr)
			if err != nil {
				return medicinebus.NewMedicine{}, fmt.Errorf("parse: %w", err)
			}
			tags[i] = tag
		}
	}

	var expiryDate time.Time
	if app.ExpiryDate != "" {
		var err error
		expiryDate, err = time.Parse(time.RFC3339, app.ExpiryDate)
		if err != nil {
			return medicinebus.NewMedicine{}, fmt.Errorf("parse: %w", err)
		}
	}


	med := medicinebus.NewMedicine{
		Name:		  app.Name,
		Description:  app.Description,
		Manufacturer: app.Manufacturer,
		Type:		  app.Type,
		Tags:         tags,
		ExpiryDate:   expiryDate,
	}

	return med, nil
}

// Validate checks the data in the model is considered clean.
func (app NewMedicine) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}

// UpdateMedicine defines the data needed to update a medicine.
type UpdateMedicine struct {
	Name 		 *string  `json:"name"`
	Description  *string  `json:"description"`
	Manufacturer *string  `json:"manufacturer"`
	Type		 *string  `json:"type"`
	Tags 		 []string `json:"tags"`
	ExpiryDate   *string  `json:"expiryDate"`
}

func toBusUpdateMedicine (app UpdateMedicine) (medicinebus.UpdateMedicine, error) {
	var tags []uuid.UUID
	if app.Tags != nil {
		tags = make([]uuid.UUID, len(app.Tags))
		for i, tagIDStr := range app.Tags {
			tag, err := uuid.Parse(tagIDStr)
			if err != nil {
				return medicinebus.UpdateMedicine{}, fmt.Errorf("parse: %w", err)
			}
			tags[i] = tag
		}
	}

	var expiryDate time.Time
	if app.ExpiryDate != nil {
		var err error
		expiryDate, err = time.Parse(time.RFC3339, *app.ExpiryDate)
		if err != nil {
			return medicinebus.UpdateMedicine{}, fmt.Errorf("parse: %w", err)
		}
	}

	um := medicinebus.UpdateMedicine{
		Name:		  app.Name,
		Description:  app.Description,
		Manufacturer: app.Manufacturer,
		Type:		  app.Type,
		Tags:		  tags,
		ExpiryDate:   &expiryDate,
	}

	return um, nil
}

// Validate checks the data in the model is considered clean.
func (app UpdateMedicine) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}