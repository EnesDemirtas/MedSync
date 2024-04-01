// Package medicinedb contains medicine related CRUD functionality.
package medicinedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/medicine"
	"github.com/EnesDemirtas/medisync/business/data/sqldb"
	"github.com/EnesDemirtas/medisync/business/data/sqldb/dbarray"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/business/web/order"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for medicine database access.
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
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (medicine.Storer, error) {
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

// Create inserts a new medicine into the database.
func (s *Store) Create(ctx context.Context, med medicine.Medicine) error {
	const q = `
	INSERT INTO medicines
		(medicine_id, name, description, manufacturer, type, tags, expiry_date, date_created, date_updated)
	VALUES
		(:medicine_id, :name, :description, :manufacturer, :type, :tags, :expiry_date, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBMedicine(med)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", medicine.ErrUniquePK)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a medicine document in the dabase.
func (s *Store) Update(ctx context.Context, med medicine.Medicine) error {
	const q = `
	UPDATE
		medicines
	SET
		"name" = :name,
		"description" = :description,
		"manufacturer" = :manufacturer,
		"type" = :type,
		"tags" = :tags,
		"expiry_date" = :expiry_date
	WHERE
		medicine_id = :medicine_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBMedicine(med)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return medicine.ErrUniquePK
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a medicine from the database.
func (s *Store) Delete(ctx context.Context, med medicine.Medicine) error {
	data := struct {
		ID string `db:"medicine_id"`
	}{
		ID: med.ID.String(),
	}

	const q = `
	DELETE FROM
		medicines
	WHERE
		medicine_id = :medicine_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing medicines from the database.
func (s *Store) Query(ctx context.Context, filter medicine.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]medicine.Medicine, error) {
	data := map[string]interface{}{
		"offset":		 (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		medicine_id, name, description, manufacturer, type, tags, expiry_date, date_created, date_updated
	FROM
		medicines`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbMedicines []dbMedicine
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbMedicines); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreMedicineSlice(dbMedicines), nil
}

// Count returns the total number of medicines in the database.
func (s *Store) Count(ctx context.Context, filter medicine.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		medicines`

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

// QueryByID gets the specified medicine from the database.
func (s *Store) QueryByID(ctx context.Context, medicineID uuid.UUID) (medicine.Medicine, error) {
	data := struct {
		ID string `db:"medicine_id"`
	}{
		ID: medicineID.String(),
	}

	const q = `
	SELECT
		medicine_id, name, description, manufacturer, type, tags, expiry_date, date_created, date_updated
	FROM
		medicines
	WHERE
		medicine_id = :medicine_id`

	var dbMedicine dbMedicine
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbMedicine); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return medicine.Medicine{}, fmt.Errorf("db: %w", medicine.ErrNotFound)
		}
		return medicine.Medicine{}, fmt.Errorf("db: %w", err)
	}

	return toCoreMedicine(dbMedicine), nil
}

// QueryByIDs gets the specified medicines from the database.
func (s *Store) QueryByIDs(ctx context.Context, medicineIDs []uuid.UUID) ([]medicine.Medicine, error) {
	ids := make([]string, len(medicineIDs))
	for i, medicineID := range medicineIDs {
		ids[i] = medicineID.String()
	}
	
	data := struct {
		ID any `db:"medicine_id"`
	}{
		ID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		medicine_id, name, description, manufacturer, type, tags, expiry_date, date_created, date_updated
	FROM
		medicines
	WHERE
		medicine_id = ANY(:medicine_id)`

	var dbMedicines []dbMedicine
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbMedicines); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return nil, medicine.ErrNotFound
		}
		return nil, fmt.Errorf("db: %w", err)
	}

	return toCoreMedicineSlice(dbMedicines), nil
}

// QueryByName gets the specified medicine from the database by name.
func (s *Store) QueryByName(ctx context.Context, name string) (medicine.Medicine, error) {
	data := struct {
		Name string `db:"name"`
	}{
		Name: name,
	}

	const q = `
	SELECT
		medicine_id, name, description, manufacturer, type, tags, expiry_date, date_created, date_updated
	FROM
		medicines
	WHERE
		name = :name`

	var dbMedicine dbMedicine
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbMedicine); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return medicine.Medicine{}, fmt.Errorf("db: %w", medicine.ErrNotFound)
		}
		return medicine.Medicine{}, fmt.Errorf("db: %w", err)
	}

	return toCoreMedicine(dbMedicine), nil
}