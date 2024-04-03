// Package inventorydb contains inventory related CRUD functionality.
package inventorydb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/business/data/sqldb"
	"github.com/EnesDemirtas/medisync/business/data/sqldb/dbarray"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for inventory database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// ExecuteUnderTransacrion constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (inventorybus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create inserts a new inventory into the database.
func (s *Store) Create(ctx context.Context, inv inventorybus.Inventory) error {
	const q = `
	INSERT INTO inventories
		(inventory_id, name, description, medicines, date_created, date_updated)
	VALUES
		(:inventory_id, :name, :description, :medicines, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBInventory(inv)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", inventorybus.ErrUniquePK)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces an inventory document in the database.
func (s *Store) Update(ctx context.Context, inv inventorybus.Inventory) error {
	const q = `
	UPDATE
		inventories
	SET
		"name" = :name,
		"description" = :description,
		"medicines" = :medicines
	WHERE
		inventory_id = :inventory_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBInventory(inv)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return inventorybus.ErrUniquePK
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes an inventory from the database.
func (s *Store) Delete(ctx context.Context, inv inventorybus.Inventory) error {
	data := struct {
		ID string `db:"inventory_id"`
	}{
		ID: inv.ID.String(),
	}

	const q = `
	DELETE FROM
		inventories
	WHERE
		inventory_id = :inventory_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing inventories from the database.
func (s *Store) Query(ctx context.Context, filter inventorybus.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]inventorybus.Inventory, error) {
	data := map[string]interface{}{
		"offset":		 (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		inventory_id, name, description, medicines, date_created, date_updated
	FROM
		inventories`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbInventories []dbInventory
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbInventories); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreInventorySlice(dbInventories), nil
}

// Count returns the total number of inventories in the database.
func (s *Store) Count(ctx context.Context, filter inventorybus.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		inventories`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified inventory from the database.
func (s *Store) QueryByID(ctx context.Context, inventoryID uuid.UUID) (inventorybus.Inventory, error) {
	data := struct {
		ID string `db:"inventory_id"`
	}{
		ID: inventoryID.String(),
	}

	const q = `
	SELECT
		inventory_id, name, description, medicines, date_created, date_updated
	FROM
		inventories
	WHERE
		inventory_id = :inventory_id`

	var dbInventory dbInventory
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbInventory); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return inventorybus.Inventory{}, fmt.Errorf("db: %w", inventorybus.ErrNotFound)
		}
		return inventorybus.Inventory{}, fmt.Errorf("db: %w", err)
	}

	return toCoreInventory(dbInventory), nil
}

// QueryByIDs gets the specified inventories from the database.
func (s *Store) QueryByIDs(ctx context.Context, inventoryIDs []uuid.UUID) ([]inventorybus.Inventory, error) {
	ids := make([]string, len(inventoryIDs))
	for i, inventoryID := range inventoryIDs {
		ids[i] = inventoryID.String()
	}
	
	data := struct {
		ID any `db:"inventory_id"`
	}{
		ID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		inventory_id, name, description, medicines, date_created, date_updated
	FROM
		inventories
	WHERE
		inventory_id = ANY(:inventory_id)`

	var dbInventories []dbInventory
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbInventories); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return nil, inventorybus.ErrNotFound
		}
		return nil, fmt.Errorf("db: %w", err)
	}

	return toCoreInventorySlice(dbInventories), nil
}

// QueryByName gets the specified inventory from the database by name.
func (s *Store) QueryByName(ctx context.Context, name string) (inventorybus.Inventory, error) {
	data := struct {
		Name string `db:"name"`
	}{
		Name: name,
	}

	const q = `
	SELECT
		inventory_id, name, description, medicines, date_created, date_updated
	FROM
		inventories
	WHERE
		name = :name`

	var dbInventory dbInventory
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbInventory); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return inventorybus.Inventory{}, fmt.Errorf("db: %w", inventorybus.ErrNotFound)
		}
		return inventorybus.Inventory{}, fmt.Errorf("db: %w", err)
	}

	return toCoreInventory(dbInventory), nil
}