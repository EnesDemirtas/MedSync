// Package tagdb contains tag related CRUD functionality.
package tagdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/business/data/sqldb"
	"github.com/EnesDemirtas/medisync/business/data/sqldb/dbarray"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for tag database access.
type Store struct {
	log *logger.Logger
	db	sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:	 db,
	}
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (tagbus.Storer, error) {
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


// Create inserts a new tag into the database.
func (s *Store) Create(ctx context.Context, tag tagbus.Tag) error {
	const q = `
	INSERT INTO tags
		(tag_id, name)
	VALUES
		(:tag_id, :name)`
	
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTag(tag)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a tag from the database.
func (s *Store) Delete(ctx context.Context, tag tagbus.Tag) error {
	data := struct {
		ID string `db:"tag_id"`
	}{
		ID: tag.ID.String(),
	}

	const q = `
	DELETE FROM
		tags
	WHERE
		tag_id = :tag_id`
	
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a tag document in the database.
func (s *Store) Update(ctx context.Context, tag tagbus.Tag) error {
	const q = `
	UPDATE
		tags
	SET
		"name"	= :name
	WHERE
		tag_id = :tag_id`
	
	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBTag(tag)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing tags from database.
func (s *Store) Query(ctx context.Context, filter tagbus.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]tagbus.Tag, error) {
	data := map[string]interface{}{
		"offset": 		 (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		tag_id, name
	FROM
		tags`
	
	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbTags []dbTag
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbTags); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	tags, err := toCoreTagSlice(dbTags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// Count returns the total number of tags in the database.
func (s *Store) Count(ctx context.Context, filter tagbus.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		tags`
	
	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified tag from the database.
func (s *Store) QueryByID(ctx context.Context, tagID uuid.UUID) (tagbus.Tag, error) {
	data := struct {
		ID string `db:"tag_id"`
	}{
		ID: tagID.String(),
	}

	const q = `
	SELECT
		tag_id, name
	FROM
		tags
	WHERE
		tag_id = :tag_id`

	var dbTag dbTag
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbTag); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return tagbus.Tag{}, fmt.Errorf("db: %w", tagbus.ErrNotFound)
		}
		return tagbus.Tag{}, fmt.Errorf("db: %w", err)
	}

	return toCoreTag(dbTag)
}

// QueryByIDs gets the specified tags from the database.
func (s *Store) QueryByIDs(ctx context.Context, tagIDs []uuid.UUID) ([]tagbus.Tag, error) {
	ids := make([]string, len(tagIDs))
	for i, tagID := range tagIDs {
		ids[i] = tagID.String()
	}

	data := struct {
		ID any `db:"tag_id"`
	}{
		ID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		tag_id, name
	FROM
		tags
	WHERE
		tag_id = ANY(:tag_id)`

	var dbTags []dbTag
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbTags); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return nil, tagbus.ErrNotFound
		}
		return nil, fmt.Errorf("db: %w", err)
	}

	return toCoreTagSlice(dbTags)
}

// TODO: QueryByName