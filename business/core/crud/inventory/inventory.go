// Package inventory provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/EnesDemirtas/medisync/business/core/crud/delegate"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicine"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/business/web/order"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
)

// Storer interface ddeclares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, inventory Inventory) error
	Update(ctx context.Context, inventory Inventory) error
	Delete(ctx context.Context, inventory Inventory) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Inventory, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, inventoryID uuid.UUID) (Inventory, error)
	QueryByIDs(ctx context.Context, inventoryIDs []uuid.UUID) ([]Inventory, error)
}

// Core manages the set of APIs for inventory access.
type Core struct {
	log 			*logger.Logger
	medicineCore 	*medicine.Core
	delegate		*delegate.Delegate
	storer			Storer
}

// NewCore constructs an inventory core API for use.
func NewCore(log *logger.Logger, medicineCore *medicine.Core, delegate *delegate.Delegate, storer Storer) *Core {
	return &Core{
		log:      log,
		medicineCore: medicineCore,
		delegate: delegate,
		storer:   storer,
	}
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	storer, err := c.storer.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	medicineCore, err := c.medicineCore.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	core := Core{
		log:      c.log,
		medicineCore: medicineCore,
		delegate: c.delegate,
		storer:   storer,
	}

	return &core, nil
}

// Create adds a new inventory to the system.
func (c *Core) Create(ctx context.Context, newInventory NewInventory) (Inventory, error) {
	_, err := c.medicineCore.QueryByIDs(ctx, newInventory.Medicines)
	if err != nil {
		return Inventory{}, fmt.Errorf("medicine.querybyids: %s: %w", newInventory.Medicines, err)
	}

	now := time.Now()

	inventory := Inventory{
		ID: 			uuid.New(),
		Name:			newInventory.Name,
		Description: 	newInventory.Description,
		Medicines: 		newInventory.Medicines,
		DateCreated: 	now,
		DateUpdated: 	now,
	}

	if err := c.storer.Create(ctx, inventory); err != nil {
		return Inventory{}, fmt.Errorf("create: %w", err)
	}

	return inventory, nil
}

// Update modifies information about an inventory.
func (c *Core) Update(ctx context.Context, inventory Inventory, updatedInventory UpdateInventory) (Inventory, error) {
	if updatedInventory.Name != nil {
		inventory.Name = *updatedInventory.Name
	}

	if updatedInventory.Description != nil {
		inventory.Description = *updatedInventory.Description
	}

	if updatedInventory.Medicines != nil {
		_, err := c.medicineCore.QueryByIDs(ctx, updatedInventory.Medicines)
		if err != nil {
			return Inventory{}, fmt.Errorf("medicine.querybyids: %s: %w", updatedInventory.Medicines, err)
		}

		inventory.Medicines = updatedInventory.Medicines
	}

	inventory.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, inventory); err != nil {
		return Inventory{}, fmt.Errorf("update: %w", err)
	}

	// Other domains may need to know when a medicine is updated so business
	// logic can be applied. This represents a delegate call to other domains.
	// if err := c.delegate.Call(ctx, ActionUpdatedData(updatedMed, med.ID)); err != nil {
	// 	return User{}, fmt.Errorf("failed to execute `%s` action: %w", ActionUpdated, err)
	// }

	return inventory, nil
}

// Delete removes the specified inventory.
func (c *Core) Delete(ctx context.Context, inventory Inventory) error {
	if err := c.storer.Delete(ctx, inventory); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing inventories.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Inventory, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	inventories, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return inventories, nil
}


// Count returns the total number of inventories.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the inventory by the specified ID.
func (c *Core) QueryByID(ctx context.Context, inventoryID uuid.UUID) (Inventory, error) {
	inventory, err := c.storer.QueryByID(ctx, inventoryID)
	if err != nil {
		return Inventory{}, fmt.Errorf("query: inventoryID[%s]: %w", inventoryID, err)
	}

	return inventory, nil
}

// QueryByIDs finds the inventories by a scpedified Inventory IDs.
func (c *Core) QueryByIDs(ctx context.Context, inventoryIDs []uuid.UUID) ([]Inventory, error) {
	inventories, err := c.storer.QueryByIDs(ctx, inventoryIDs)
	if err != nil {
		return nil, fmt.Errorf("query: inventoryIDs[%s]: %w", inventoryIDs, err)
	}

	return inventories, nil
}
