// Package tag provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package tagbus

import (
	"context"
	"errors"
	"fmt"

	"github.com/EnesDemirtas/medisync/business/api/delegate"
	"github.com/EnesDemirtas/medisync/business/api/order"
	"github.com/EnesDemirtas/medisync/business/data/transaction"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var ErrNotFound = errors.New("tag not found")

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, tag Tag) error
	Update(ctx context.Context, tag Tag) error
	Delete(ctx context.Context, tag Tag) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Tag, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, tagID uuid.UUID) (Tag, error)
	QueryByIDs(ctx context.Context, tagIDs []uuid.UUID) ([]Tag, error)
	// QueryByName(ctx context.Context, tagName string) (Tag, error)
}

// Core manages the set of APIs for tag access.
type Core struct {
	log 		*logger.Logger
	storer 		Storer
	delegate 	*delegate.Delegate
}

// NewCore constructs a tag core API for use.
func NewCore(log *logger.Logger, delegate *delegate.Delegate, storer Storer) *Core {
	return &Core{
		log:		log,
		delegate: 	delegate,
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

	core := Core{
		log:		c.log,
		delegate:	c.delegate,
		storer:		trS,
	}

	return &core, nil
}

// Create adds a new tag to the system.
func (c *Core) Create(ctx context.Context, newTag NewTag) (Tag, error) {
	tag := Tag{
		ID:		uuid.New(),
		Name:	newTag.Name,
	}

	if err := c.storer.Create(ctx, tag); err != nil {
		return Tag{}, fmt.Errorf("create: %w", err)
	}

	return tag, nil
}

// Update modifies information about a tag.
func (c *Core) Update(ctx context.Context, tag Tag, updatedTag UpdateTag) (Tag, error) {
	if updatedTag.Name != nil {
		tag.Name = *updatedTag.Name
	}

	if err := c.storer.Update(ctx, tag); err != nil {
		return Tag{}, fmt.Errorf("update: %w", err)
	}

	return tag, nil
}

// Delete removes the specified tag.
func (c *Core) Delete(ctx context.Context, tag Tag) error {
	if err := c.storer.Delete(ctx, tag); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing tags.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Tag, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	tags, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return tags, nil
}

// Count returns the total number of tags.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the tag by the specified ID.
func (c *Core) QueryByID(ctx context.Context, tagID uuid.UUID) (Tag, error) {
	tag, err := c.storer.QueryByID(ctx, tagID)
	if err != nil {
		return Tag{}, fmt.Errorf("query: tagID[%s]: %w", tagID, err)
	}

	return tag, nil
}

// QueryByIDs finds the tags by a specified Tag IDs.
func (c *Core) QueryByIDs(ctx context.Context, tagIDs []uuid.UUID) ([]Tag, error) {
	tags, err := c.storer.QueryByIDs(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("query: tagIDs[%s]: %w", tagIDs, err)
	}

	return tags, nil
}

// QueryByName finds the tag by a specified tag name.
// func (c *Core) QueryByName(ctx context.Context, tagName string) (Tag, error) {
// 	tag, err := c.storer.QueryByName(ctx, tagName)
// 	if err != nil {
// 		return Tag{}, fmt.Errorf("query: tagName[%s]: %w", tagName, err)
// 	}

// 	return tag, nil
// }