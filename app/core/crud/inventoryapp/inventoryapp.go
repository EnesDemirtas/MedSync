// Package inventoryapp maintains the app layer api for the inventory domain.
package inventoryapp

import (
	"context"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/app/api/page"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
)

// Core manages the set of app layer api functions for the inventory domain.
type Core struct {
	inventoryBus *inventorybus.Core
}

// NewCore constructs an inventory core API for use.
func NewCore(inventoryBus *inventorybus.Core) *Core {
	return &Core {
		inventoryBus: inventoryBus,
	}
}

// Create adds a new inventory to the system.
func (c *Core) Create(ctx context.Context, app NewInventory) (Inventory, error) {
	ni := toBusNewInventory(app)

	inv, err := c.inventoryBus.Create(ctx, ni)
	if err != nil {
		return Inventory{}, errs.Newf(errs.Internal, "create: inv[%+v]: %s", inv, err)
	}

	return toAppInventory(inv), nil
}

// Update updates an existing inventory.
func (c *Core) Update(ctx context.Context, app UpdateInventory) (Inventory, error) {
	inv, err := mid.GetInventory(ctx)
	if err != nil {
		return Inventory{}, errs.Newf(errs.Internal, "inventory missing in context: %s", err)
	}

	busUpdInv, err := toBusUpdateInventory(app)
	if err != nil {
		return Inventory{}, err
	}

	updInv, err := c.inventoryBus.Update(ctx, inv, busUpdInv)
	if err != nil {
		return Inventory{}, errs.Newf(errs.Internal, "update: inventoryID[%s] up[%+v]: %s", inv.ID, app, err)
	}

	return toAppInventory(updInv), nil
}

// Delete removes an inventory from the system.
func (c *Core) Delete(ctx context.Context) error {
	inv, err := mid.GetInventory(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "inventoryID missing in context: %s", err)
	}

	if err := c.inventoryBus.Delete(ctx, inv); err != nil {
		return errs.Newf(errs.Internal, "delete: inventoryID[%s]: %s", inv.ID, err)
	}

	return nil
}

// Query returns a list of inventories with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[Inventory], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[Inventory]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Inventory]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[Inventory]{}, err
	}

	invs, err := c.inventoryBus.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Inventory]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := c.inventoryBus.Count(ctx, filter)
	if err != nil {
		return page.Document[Inventory]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppInventories(invs), total, qp.Page, qp.Rows), nil
}

// QueryByID returns an inventory by its ID.
func (c *Core) QueryByID(ctx context.Context) (Inventory, error) {
	inv, err := mid.GetInventory(ctx)
	if err != nil {
		return Inventory{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppInventory(inv), nil
}