package tagdb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/tag"
	"github.com/google/uuid"
)

type dbTag struct {
	ID		uuid.UUID `db:"tag_id"`
	Name	string	  `db:"name"`
}

func toDBTag(tag tag.Tag) dbTag {
	tagDB := dbTag{
		ID:		tag.ID,
		Name:	tag.Name,
	}

	return tagDB
}

func toCoreTag(dbTag dbTag) (tag.Tag, error) {
	tag := tag.Tag{
		ID:		dbTag.ID,
		Name:	dbTag.Name,
	}

	return tag, nil
}

func toCoreTagSlice(dbTags []dbTag) ([]tag.Tag, error) {
	tags := make([]tag.Tag, len(dbTags))

	for i, dbTag := range dbTags {
		var err error
		tags[i], err = toCoreTag(dbTag)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return tags, nil
}
