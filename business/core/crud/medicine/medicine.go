// Package medicine provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package medicine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EnesDemirtas/medisync/business/core/crud/delegate"
	"github.com/EnesDemirtas/medisync/business/core/crud/tag"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/business/web/order"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound 		= errors.New("medicine not found")
)

// Storer interface ddeclares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, med Medicine) error
	Update(ctx context.Context, med Medicine) error
	Delete(ctx context.Context, med Medicine) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Medicine, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, medicineID uuid.UUID) (Medicine, error)
	QueryByIDs(ctx context.Context, medicineIDs []uuid.UUID) ([]Medicine, error)
}

// Core manages the set of APIs for medicine access.
type Core struct {
	log 		*logger.Logger
	tagCore		*tag.Core
	delegate	*delegate.Delegate
	storer 		Storer
}

// NewCore constructs a medicine core API for use.
func NewCore(log *logger.Logger, tagCore *tag.Core, delegate *delegate.Delegate, storer Storer) *Core {
	return &Core{
		log: 		log,
		tagCore:	tagCore,
		delegate:	delegate,
		storer:		storer,
	}
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	trS, err := c.storer.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	tagCore, err := c.tagCore.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	core := Core{
		log:		c.log,
		tagCore: 	tagCore,
		delegate:	c.delegate,
		storer:		trS,
	}

	return &core, nil
}

// Create adds a new medicine to the system.
func (c *Core) Create(ctx context.Context, newMed Medicine) (Medicine, error) {
	_, err := c.tagCore.QueryByIDs(ctx, newMed.Tags)
	if err != nil {
		return Medicine{}, fmt.Errorf("tag.querybyids: %s: %w", newMed.Tags, err)
	}

	now := time.Now()

	med := Medicine{
		ID: 			uuid.New(),
		Name:			newMed.Name,
		Description: 	newMed.Manufacturer,
		Manufacturer: 	newMed.Manufacturer,
		Type:			newMed.Type,
		Tags:			newMed.Tags,
		ExpiryDate: 	newMed.ExpiryDate,
		CreatedDate: 	now,
		UpdatedDate: 	now,
	}

	if err := c.storer.Create(ctx, med); err != nil {
		return Medicine{}, fmt.Errorf("create: %w", err)
	}

	return med, nil
}

// Update modifies information about a medicine.
func (c *Core) Update(ctx context.Context, med Medicine, updatedMed UpdateMedicine) (Medicine, error) {
	if updatedMed.Name != nil {
		med.Name = *updatedMed.Name
	}

	if updatedMed.Description != nil {
		med.Description = *updatedMed.Description
	}

	if updatedMed.Manufacturer != nil {
		med.Manufacturer = *updatedMed.Manufacturer
	}

	if updatedMed.Type != nil {
		med.Type = *updatedMed.Type
	}

	if updatedMed.Tags != nil {
		_, err := c.tagCore.QueryByIDs(ctx, updatedMed.Tags)
		if err != nil {
			return Medicine{}, fmt.Errorf("tag.querybyids: %s: %w", updatedMed.Tags, err)
		}
		
		med.Tags = updatedMed.Tags
	}

	if updatedMed.ExpiryDate != nil {
		med.ExpiryDate = *updatedMed.ExpiryDate
	}

	med.UpdatedDate = time.Now()

	if err := c.storer.Update(ctx, med); err != nil {
		return Medicine{}, fmt.Errorf("update: %w", err)
	}

	// Other domains may need to know when a medicine is updated so business
	// logic can be applied. This represents a delegate call to other domains.
	// if err := c.delegate.Call(ctx, ActionUpdatedData(updatedMed, med.ID)); err != nil {
	// 	return User{}, fmt.Errorf("failed to execute `%s` action: %w", ActionUpdated, err)
	// }

	return med, nil
}

// Delete removes the specified medicine.
func (c *Core) Delete(ctx context.Context, med Medicine) error {
	if err := c.storer.Delete(ctx, med); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing medicines.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Medicine, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	medicines, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return medicines, nil
}

// Count returns the total number of medicines.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the medicine by the specified ID.
func (c *Core) QueryByID(ctx context.Context, medicineID uuid.UUID) (Medicine, error) {
	medicine, err := c.storer.QueryByID(ctx, medicineID)
	if err != nil {
		return Medicine{}, fmt.Errorf("query: medicineID[%s]: %w", medicineID, err)
	}

	return medicine, nil
}

// QueryByIDs finds the medicines by a scpedified Medicine IDs.
func (c *Core) QueryByIDs(ctx context.Context, medicineIDs []uuid.UUID) ([]Medicine, error) {
	medicines, err := c.storer.QueryByIDs(ctx, medicineIDs)
	if err != nil {
		return nil, fmt.Errorf("query: medicineIDs[%s]: %w", medicineIDs, err)
	}

	return medicines, nil
}