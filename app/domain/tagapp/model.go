package tagapp

import (
	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/foundation/validate"
)

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page     int    `query:"page"`
	Rows     int    `query:"rows"`
	OrderBy  string `query:"orderBy"`
	ID       string `query:"tag_id"`
	Name     string `query:"name"`
}

// Tag represents information about an individual tag.
type Tag struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
}

func toAppTag(tag tagbus.Tag) Tag {
	return Tag{
		ID:          tag.ID.String(),
		Name:        tag.Name,
	}
}

func toAppTags(tags []tagbus.Tag) []Tag {
	items := make([]Tag, len(tags))
	for i, tag := range tags {
		items[i] = toAppTag(tag)
	}

	return items
}

// NewTag defines the data needed to add a new tag.
type NewTag struct {
	Name     string  `json:"name" validate:"required"`
}

func toBusNewTag(app NewTag) (tagbus.NewTag, error) {
	tag := tagbus.NewTag{
		Name:     app.Name,
	}

	return tag, nil
}

// Validate checks the data in the model is considered clean.
func (app NewTag) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}

// UpdateTag defines the data needed to update a tag.
type UpdateTag struct {
	Name     *string  `json:"name"`
}

func toBusUpdateTag(app UpdateTag) tagbus.UpdateTag {
	core := tagbus.UpdateTag{
		Name:     app.Name,
	}

	return core
}

// Validate checks the data in the model is considered clean.
func (app UpdateTag) Validate() error {
	if err := validate.Check(app); err != nil {
		return errs.Newf(errs.FailedPrecondition, "validate: %s", err)
	}

	return nil
}
