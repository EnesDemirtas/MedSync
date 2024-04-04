// Package tagapp maintains the app layer api for the tag domain.
package tagapp

import (
	"context"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/app/api/page"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
)

// Core manages the set of app layer api functions for the tag domain.
type Core struct {
	tagBus *tagbus.Core
}

// NewCore constructs a tag core API for use.
func NewCore(tagBus *tagbus.Core) *Core {
	return &Core{
		tagBus: tagBus,
	}
}

// Create adds a new tag to the system.
func (c *Core) Create(ctx context.Context, app NewTag) (Tag, error) {
	nt, err := toBusNewTag(app)
	if err != nil {
		return Tag{}, errs.New(errs.FailedPrecondition, err)
	}

	tag, err := c.tagBus.Create(ctx, nt)
	if err != nil {
		return Tag{}, errs.Newf(errs.Internal, "create: prd[%+v]: %s", tag, err)
	}

	return toAppTag(tag), nil
}

// Update updates an existing tag.
func (c *Core) Update(ctx context.Context, app UpdateTag) (Tag, error) {
	tag, err := mid.GetTag(ctx)
	if err != nil {
		return Tag{}, errs.Newf(errs.Internal, "tag missing in context: %s", err)
	}

	updTag, err := c.tagBus.Update(ctx, tag, toBusUpdateTag(app))
	if err != nil {
		return Tag{}, errs.Newf(errs.Internal, "update: tagID[%s] up[%+v]: %s", tag.ID, app, err)
	}

	return toAppTag(updTag), nil
}

// Delete removes a tag from the system.
func (c *Core) Delete(ctx context.Context) error {
	tag, err := mid.GetTag(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "tagID missing in context: %s", err)
	}

	if err := c.tagBus.Delete(ctx, tag); err != nil {
		return errs.Newf(errs.Internal, "delete: tagID[%s]: %s", tag.ID, err)
	}

	return nil
}

// Query returns a list of tags with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[Tag], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[Tag]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Tag]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[Tag]{}, err
	}

	tags, err := c.tagBus.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Tag]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := c.tagBus.Count(ctx, filter)
	if err != nil {
		return page.Document[Tag]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppTags(tags), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a tag by its ID.
func (c *Core) QueryByID(ctx context.Context) (Tag, error) {
	tag, err := mid.GetTag(ctx)
	if err != nil {
		return Tag{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppTag(tag), nil
}
