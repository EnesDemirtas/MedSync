// Package medicineapp maintains the app layer api for the medicine domain.
package medicineapp

import (
	"context"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/app/api/page"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
)

// Core manages the set of app layer api functions for the medicine domain.
type Core struct {
	medicineBus *medicinebus.Core
}

// NewCore constructs a medicine core API for use.
func NewCore(medicineBus *medicinebus.Core) *Core {
	return &Core{
		medicineBus: medicineBus,
	}
}

// Create adds a new medicine to the system.
func (c *Core) Create(ctx context.Context, app NewMedicine) (Medicine, error) {
	nm , err := toBusNewMedicine(app)
	if err != nil {
		return Medicine{}, errs.New(errs.FailedPrecondition, err)
	}

	med, err := c.medicineBus.Create(ctx, nm)
	if err != nil {
		return Medicine{}, errs.Newf(errs.Internal, "create: med[%+v]: %s", med, err)
	}

	return toAppMedicine(med), nil
}

// Update updates an existing medicine.
func (c *Core) Update(ctx context.Context, app UpdateMedicine) (Medicine, error) {
	med, err := mid.GetMedicine(ctx)
	if err != nil {
		return Medicine{}, errs.Newf(errs.Internal, "medicine missing in context: %s", err)
	} 
	
	busUpdMed, err := toBusUpdateMedicine(app)
	if err != nil {
		return Medicine{}, err
	}

	um, err := c.medicineBus.Update(ctx, med, busUpdMed)
	if err != nil {
		return Medicine{}, errs.Newf(errs.Internal, "update: medicineID[%s] um[%+v]: %s", med.ID, app, err)
	}

	return toAppMedicine(um), nil
}

// Delete removes a medicine from the system.
func (c *Core) Delete(ctx context.Context) error {
	med, err := mid.GetMedicine(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "medicineID missing in context: %s", err)
	}

	if err := c.medicineBus.Delete(ctx, med); err != nil {
		return errs.Newf(errs.Internal, "delete: medicineID[%s]: %s", med.ID, err)
	}

	return nil
}

// Query returns a list of medicines with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[Medicine], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[Medicine]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Medicine]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[Medicine]{}, err
	}

	meds, err := c.medicineBus.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Medicine]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := c.medicineBus.Count(ctx, filter)
	if err != nil {
		return page.Document[Medicine]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewDocument[Medicine](toAppMedicines(meds), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a medicine by its ID.
func (c *Core) QueryByID(ctx context.Context) (Medicine, error) {
	med, err := mid.GetMedicine(ctx)
	if err != nil {
		return Medicine{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppMedicine(med), nil
}